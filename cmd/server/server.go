package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpctags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/odpf/optimus/api/proto/odpf/optimus/core/v1beta1"
	"github.com/odpf/optimus/config"
	"github.com/odpf/optimus/core/gossip"
	_ "github.com/odpf/optimus/ext/datastore"
	"github.com/odpf/optimus/ext/executor/noop"
	"github.com/odpf/optimus/ext/scheduler/airflow"
	"github.com/odpf/optimus/ext/scheduler/airflow2"
	"github.com/odpf/optimus/ext/scheduler/airflow2/compiler"
	"github.com/odpf/optimus/ext/scheduler/prime"
	"github.com/odpf/optimus/models"
	_ "github.com/odpf/optimus/plugin"
	"github.com/odpf/optimus/store/postgres"
	"github.com/odpf/optimus/utils"
	"github.com/odpf/salt/log"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

const (
	shutdownWait       = 30 * time.Second
	GRPCMaxRecvMsgSize = 64 << 20 // 64MB
	GRPCMaxSendMsgSize = 64 << 20 // 64MB

	DialTimeout      = time.Second * 5
	BootstrapTimeout = time.Second * 10
)

func checkRequiredConfigs(conf config.ServerConfig) error {
	errRequiredMissing := errors.New("required config missing")
	if conf.IngressHost == "" {
		return fmt.Errorf("serve.ingress_host: %w", errRequiredMissing)
	}
	if conf.ReplayNumWorkers < 1 {
		return fmt.Errorf("%s should be greater than 0", config.KeyServeReplayNumWorkers)
	}
	if conf.DB.DSN == "" {
		return fmt.Errorf("serve.db.dsn: %w", errRequiredMissing)
	}
	if parsed, err := url.Parse(conf.DB.DSN); err != nil {
		return fmt.Errorf("failed to parse serve.db.dsn: %w", err)
	} else {
		if parsed.Scheme != "postgres" {
			return errors.New("unsupported database scheme, use 'postgres'")
		}
	}
	return nil
}

func setupDB(l log.Logger, conf config.Optimus) (*gorm.DB, error) {
	// setup db
	if err := postgres.Migrate(conf.Server.DB.DSN); err != nil {
		return nil, fmt.Errorf("postgres.Migrate: %w", err)
	}
	dbConn, err := postgres.Connect(conf.Server.DB.DSN, conf.Server.DB.MaxIdleConnection,
		conf.Server.DB.MaxOpenConnection, l.Writer())
	if err != nil {
		return nil, fmt.Errorf("postgres.Connect: %w", err)
	}
	return dbConn, nil
}

func setupGRPCServer(l log.Logger) (*grpc.Server, error) {
	// Logrus entry is used, allowing pre-definition of certain fields by the user.
	grpcLogLevel, err := logrus.ParseLevel(l.Level())
	if err != nil {
		return nil, err
	}
	grpcLogrus := logrus.New()
	grpcLogrus.SetLevel(grpcLogLevel)
	grpcLogrusEntry := logrus.NewEntry(grpcLogrus)
	// Shared options for the logger, with a custom gRPC code to log level function.
	opts := []grpc_logrus.Option{
		grpc_logrus.WithLevels(grpc_logrus.DefaultCodeToLevel),
	}
	// Make sure that log statements internal to gRPC library are logged using the logrus logger as well.
	grpc_logrus.ReplaceGrpcLogger(grpcLogrusEntry)

	grpcOpts := []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(
			grpctags.UnaryServerInterceptor(grpctags.WithFieldExtractor(grpctags.CodeGenRequestFieldExtractor)),
			grpc_logrus.UnaryServerInterceptor(grpcLogrusEntry, opts...),
			otelgrpc.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
		),
		grpc_middleware.WithStreamServerChain(
			otelgrpc.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
		),
		grpc.MaxRecvMsgSize(GRPCMaxRecvMsgSize),
		grpc.MaxSendMsgSize(GRPCMaxSendMsgSize),
	}
	grpcServer := grpc.NewServer(grpcOpts...)
	reflection.Register(grpcServer)
	return grpcServer, nil
}

func prepareHTTPProxy(grpcAddr string, grpcServer *grpc.Server) (*http.Server, func(), error) {
	timeoutGrpcDialCtx, grpcDialCancel := context.WithTimeout(context.Background(), DialTimeout)
	defer grpcDialCancel()

	// prepare http proxy
	gwmux := runtime.NewServeMux(
		runtime.WithErrorHandler(runtime.DefaultHTTPErrorHandler),
	)
	// gRPC dialup options to proxy http connections
	grpcConn, err := grpc.DialContext(timeoutGrpcDialCtx, grpcAddr, []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(GRPCMaxRecvMsgSize),
			grpc.MaxCallSendMsgSize(GRPCMaxSendMsgSize),
		),
	}...)
	if err != nil {
		return nil, func() {}, fmt.Errorf("grpc.DialContext: %w", err)
	}
	runtimeCtx, runtimeCancel := context.WithCancel(context.Background())
	if err := pb.RegisterRuntimeServiceHandler(runtimeCtx, gwmux, grpcConn); err != nil {
		runtimeCancel()
		return nil, func() {}, fmt.Errorf("RegisterRuntimeServiceHandler: %w", err)
	}

	// base router
	baseMux := http.NewServeMux()
	baseMux.HandleFunc("/ping", otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	}), "Ping").ServeHTTP)
	baseMux.Handle("/api/", otelhttp.NewHandler(http.StripPrefix("/api", gwmux), "api"))

	//nolint: gomnd
	srv := &http.Server{
		Handler:      grpcHandlerFunc(grpcServer, baseMux),
		Addr:         grpcAddr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	cleanup := func() {
		runtimeCancel()
	}

	return srv, cleanup, nil
}

func initScheduler(l log.Logger, conf config.Optimus, projectRepoFac *projectRepoFactory) (models.SchedulerUnit, error) {
	jobCompiler := compiler.NewCompiler(conf.Server.IngressHost)
	// init default scheduler
	var scheduler models.SchedulerUnit
	switch conf.Scheduler.Name {
	case "airflow":
		scheduler = airflow.NewScheduler(
			&airflowBucketFactory{},
			&http.Client{},
			jobCompiler,
		)
	case "airflow2":
		scheduler = airflow2.NewScheduler(
			&airflowBucketFactory{},
			&http.Client{},
			jobCompiler,
		)
	default:
		return nil, fmt.Errorf("unsupported scheduler: %s", conf.Scheduler.Name)
	}

	if !conf.Scheduler.SkipInit { // TODO: This should not be required
		registeredProjects, err := projectRepoFac.New().GetAll(context.Background())
		if err != nil {
			return nil, fmt.Errorf("projectRepoFactory.GetAll(): %w", err)
		}
		// bootstrap scheduler for registered projects
		for _, proj := range registeredProjects {
			bootstrapCtx, cancel := context.WithTimeout(context.Background(), BootstrapTimeout)
			l.Info("bootstrapping project", "project name", proj.Name)
			if err := scheduler.Bootstrap(bootstrapCtx, proj); err != nil {
				// Major ERROR, but we can't make this fatal
				// other projects might be working fine
				l.Error("no bootstrapping project", "error", err)
			}
			l.Info("bootstrapped project", "project name", proj.Name)
			cancel()
		}
	}
	return scheduler, nil
}

func initPrimeCluster(l log.Logger, conf config.Optimus, jobrunRepoFac *jobRunRepoFactory, dbConn *gorm.DB) (context.CancelFunc, error) {
	models.ManualScheduler = prime.NewScheduler( // careful global variable
		jobrunRepoFac,
		func() time.Time {
			return time.Now().UTC()
		},
	)

	clusterCtx, clusterCancel := context.WithCancel(context.Background())
	clusterServer := gossip.NewServer(l)
	clusterPlanner := prime.NewPlanner(
		l,
		clusterServer, jobrunRepoFac, &instanceRepoFactory{
			db: dbConn,
		},
		utils.NewUUIDProvider(), noop.NewExecutor(), func() time.Time {
			return time.Now().UTC()
		},
	)
	cleanup := func() {
		// shutdown cluster
		clusterCancel()
		if clusterPlanner != nil {
			clusterPlanner.Close() // err is nil
		}
		if clusterServer != nil {
			clusterServer.Shutdown() // TODO: log error
		}
	}

	if conf.Scheduler.NodeID != "" {
		// start optimus cluster
		if err := clusterServer.Init(clusterCtx, conf.Scheduler); err != nil {
			return cleanup, err
		}

		clusterPlanner.Init(clusterCtx)
	}

	return cleanup, nil
}

// grpcHandlerFunc routes http1 calls to baseMux and http2 with grpc header to grpcServer.
// Using a single port for proxying both http1 & 2 protocols will degrade http performance
// but for our use-case the convenience per performance tradeoff is better suited
// if in the future, this does become a bottleneck(which I highly doubt), we can break the service
// into two ports, default port for grpc and default+1 for grpc-gateway proxy.
// We can also use something like a connection multiplexer
// https://github.com/soheilhy/cmux to achieve the same.
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}
