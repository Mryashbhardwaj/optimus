package v1beta1

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	pb "github.com/odpf/optimus/api/proto/odpf/optimus/core/v1beta1"
	"github.com/odpf/optimus/core/progress"
	"github.com/odpf/optimus/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (sv *RuntimeServiceServer) DeployJobSpecification(stream pb.RuntimeService_DeployJobSpecificationServer) error {
	startTime := time.Now()
	errNamespaces := []string{}

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			stream.Send(&pb.DeployJobSpecificationResponse{
				Success: false,
				Ack:     true,
				Message: err.Error(),
			})
			return err // immediate error returned (grpc error level)
		}
		namespaceSpec, err := sv.namespaceService.Get(stream.Context(), req.GetProjectName(), req.GetNamespaceName())
		if err != nil {
			stream.Send(&pb.DeployJobSpecificationResponse{
				Success: false,
				Ack:     true,
				Message: err.Error(),
			})
			errNamespaces = append(errNamespaces, req.NamespaceName)
			continue
		}

		jobsToKeep, err := sv.getJobsToKeep(stream.Context(), namespaceSpec, req)
		if err != nil {
			stream.Send(&pb.DeployJobSpecificationResponse{
				Success: false,
				Ack:     true,
				Message: err.Error(),
			})
			errNamespaces = append(errNamespaces, req.NamespaceName)
			continue
		}

		observers := new(progress.ObserverChain)
		observers.Join(sv.progressObserver)
		observers.Join(&jobSyncObserver{
			stream: stream,
			log:    sv.l,
			mu:     new(sync.Mutex),
		})

		// delete specs not sent for deployment from internal repository
		if err := sv.jobSvc.KeepOnly(stream.Context(), namespaceSpec, jobsToKeep, observers); err != nil {
			stream.Send(&pb.DeployJobSpecificationResponse{
				Success: false,
				Ack:     true,
				Message: fmt.Sprintf("failed to delete jobs: \n%s", err.Error()),
			})
			errNamespaces = append(errNamespaces, req.NamespaceName)
			continue
		}
		if err := sv.jobSvc.Sync(stream.Context(), namespaceSpec, observers); err != nil {
			stream.Send(&pb.DeployJobSpecificationResponse{
				Success: false,
				Ack:     true,
				Message: fmt.Sprintf("failed to sync jobs: \n%s", err.Error()),
			})
			errNamespaces = append(errNamespaces, req.NamespaceName)
			continue
		}
		runtimeDeployJobSpecificationCounter.Add(float64(len(req.Jobs)))
		stream.Send(&pb.DeployJobSpecificationResponse{
			Success: true,
			Ack:     true,
			Message: "success",
		})
	}
	sv.l.Info("finished job deployment", "time", time.Since(startTime))
	if len(errNamespaces) > 0 {
		sv.l.Warn("there's error while deploying namespaces: %v", errNamespaces)
		return fmt.Errorf("error when deploying: %v", errNamespaces)
	}
	return nil
}

func (sv *RuntimeServiceServer) getJobsToKeep(ctx context.Context, namespaceSpec models.NamespaceSpec, req *pb.DeployJobSpecificationRequest) ([]models.JobSpec, error) {
	jobs := req.GetJobs()
	if len(jobs) == 0 {
		return []models.JobSpec{}, nil
	}

	var jobsToKeep []models.JobSpec
	for _, reqJob := range jobs {
		adaptJob, err := sv.adapter.FromJobProto(reqJob)
		if err != nil {
			sv.l.Error(fmt.Sprintf("%s: cannot adapt job %s", err.Error(), reqJob.GetName()))
			continue
		}

		err = sv.jobSvc.Create(ctx, namespaceSpec, adaptJob)
		if err != nil {
			sv.l.Error(fmt.Sprintf("%s: failed to save %s", err.Error(), adaptJob.Name))
			continue
		}
		jobsToKeep = append(jobsToKeep, adaptJob)
	}

	if jobsToKeep == nil {
		return nil, errors.New("job spec creation is failed")
	}

	return jobsToKeep, nil
}

func (sv *RuntimeServiceServer) ListJobSpecification(ctx context.Context, req *pb.ListJobSpecificationRequest) (*pb.ListJobSpecificationResponse, error) {
	namespaceSpec, err := sv.namespaceService.Get(ctx, req.GetProjectName(), req.GetNamespaceName())
	if err != nil {
		return nil, mapToGRPCErr(sv.l, err, "unable to get namespace")
	}

	jobSpecs, err := sv.jobSvc.GetAll(ctx, namespaceSpec)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s: failed to retrieve jobs for project %s", err.Error(), req.GetProjectName())
	}

	jobProtos := []*pb.JobSpecification{}
	for _, jobSpec := range jobSpecs {
		jobProto, err := sv.adapter.ToJobProto(jobSpec)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "%s: failed to parse job spec %s", err.Error(), jobSpec.Name)
		}
		jobProtos = append(jobProtos, jobProto)
	}
	return &pb.ListJobSpecificationResponse{
		Jobs: jobProtos,
	}, nil
}

func (sv *RuntimeServiceServer) CheckJobSpecification(ctx context.Context, req *pb.CheckJobSpecificationRequest) (*pb.CheckJobSpecificationResponse, error) {
	namespaceSpec, err := sv.namespaceService.Get(ctx, req.GetProjectName(), req.GetNamespaceName())
	if err != nil {
		return nil, mapToGRPCErr(sv.l, err, "unable to get namespace")
	}

	j, err := sv.adapter.FromJobProto(req.GetJob())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s: failed to adapt job %s", err.Error(), req.GetJob().Name)
	}
	reqJobs := []models.JobSpec{j}

	if err = sv.jobSvc.Check(ctx, namespaceSpec, reqJobs, nil); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to compile jobs\n%s", err.Error())
	}
	return &pb.CheckJobSpecificationResponse{Success: true}, nil
}

func (sv *RuntimeServiceServer) CheckJobSpecifications(req *pb.CheckJobSpecificationsRequest, respStream pb.RuntimeService_CheckJobSpecificationsServer) error {
	namespaceSpec, err := sv.namespaceService.Get(respStream.Context(), req.GetProjectName(), req.GetNamespaceName())
	if err != nil {
		return mapToGRPCErr(sv.l, err, "unable to get namespace")
	}

	observers := new(progress.ObserverChain)
	observers.Join(sv.progressObserver)
	observers.Join(&jobCheckObserver{
		stream: respStream,
		log:    sv.l,
		mu:     new(sync.Mutex),
	})

	var reqJobs []models.JobSpec
	for _, jobProto := range req.GetJobs() {
		j, err := sv.adapter.FromJobProto(jobProto)
		if err != nil {
			return status.Errorf(codes.Internal, "%s: failed to adapt job %s", err.Error(), jobProto.Name)
		}
		reqJobs = append(reqJobs, j)
	}

	if err = sv.jobSvc.Check(respStream.Context(), namespaceSpec, reqJobs, observers); err != nil {
		return status.Errorf(codes.Internal, "failed to compile jobs\n%s", err.Error())
	}
	return nil
}

func (sv *RuntimeServiceServer) CreateJobSpecification(ctx context.Context, req *pb.CreateJobSpecificationRequest) (*pb.CreateJobSpecificationResponse, error) {
	namespaceSpec, err := sv.namespaceService.Get(ctx, req.GetProjectName(), req.GetNamespaceName())
	if err != nil {
		return nil, mapToGRPCErr(sv.l, err, "unable to get namespace")
	}

	jobSpec, err := sv.adapter.FromJobProto(req.GetSpec())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot deserialize job: \n%s", err.Error())
	}

	// validate job spec
	if err = sv.jobSvc.Check(ctx, namespaceSpec, []models.JobSpec{jobSpec}, sv.progressObserver); err != nil {
		return nil, status.Errorf(codes.Internal, "spec validation failed\n%s", err.Error())
	}

	err = sv.jobSvc.Create(ctx, namespaceSpec, jobSpec)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s: failed to save job %s", err.Error(), jobSpec.Name)
	}

	if err := sv.jobSvc.Sync(ctx, namespaceSpec, sv.progressObserver); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to sync jobs: \n%s", err.Error())
	}

	runtimeDeployJobSpecificationCounter.Inc()
	return &pb.CreateJobSpecificationResponse{
		Success: true,
		Message: fmt.Sprintf("job %s is created and deployed successfully on project %s", jobSpec.Name, req.GetProjectName()),
	}, nil
}

func (sv *RuntimeServiceServer) GetJobSpecification(ctx context.Context, req *pb.GetJobSpecificationRequest) (*pb.GetJobSpecificationResponse, error) {
	namespaceSpec, err := sv.namespaceService.Get(ctx, req.GetProjectName(), req.GetNamespaceName())
	if err != nil {
		return nil, mapToGRPCErr(sv.l, err, "unable to get namespace")
	}

	jobSpec, err := sv.jobSvc.GetByName(ctx, req.GetJobName(), namespaceSpec)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%s: error while finding the job %s", err.Error(), req.GetJobName())
	}

	jobSpecAdapt, err := sv.adapter.ToJobProto(jobSpec)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot serialize job: \n%s", err.Error())
	}

	return &pb.GetJobSpecificationResponse{
		Spec: jobSpecAdapt,
	}, nil
}

func (sv *RuntimeServiceServer) DeleteJobSpecification(ctx context.Context, req *pb.DeleteJobSpecificationRequest) (*pb.DeleteJobSpecificationResponse, error) {
	namespaceSpec, err := sv.namespaceService.Get(ctx, req.GetProjectName(), req.GetNamespaceName())
	if err != nil {
		return nil, mapToGRPCErr(sv.l, err, "unable to get namespace")
	}

	jobSpecToDelete, err := sv.jobSvc.GetByName(ctx, req.GetJobName(), namespaceSpec)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%s: job %s does not exist", err.Error(), req.GetJobName())
	}

	if err := sv.jobSvc.Delete(ctx, namespaceSpec, jobSpecToDelete); err != nil {
		return nil, status.Errorf(codes.Internal, "%s: failed to delete job %s", err.Error(), req.GetJobName())
	}

	return &pb.DeleteJobSpecificationResponse{
		Success: true,
		Message: fmt.Sprintf("job %s has been deleted", jobSpecToDelete.Name),
	}, nil
}

func (sv *RuntimeServiceServer) RefreshJobs(req *pb.RefreshJobsRequest, respStream pb.RuntimeService_RefreshJobsServer) error {
	startTime := time.Now()

	var namespaceJobNamePairs []models.NamespaceJobNamePair
	for _, namespaceJobs := range req.NamespaceJobs {
		namespaceSpec, err := sv.namespaceService.Get(respStream.Context(), req.GetProjectName(), namespaceJobs.NamespaceName)
		if err != nil {
			return mapToGRPCErr(err, "unable to get namespace")
		}

		namespaceJobNamePairs = append(namespaceJobNamePairs, models.NamespaceJobNamePair{
			Namespace: namespaceSpec,
			JobNames:  namespaceJobs.JobNames,
		})
	}

	observers := new(progress.ObserverChain)
	observers.Join(sv.progressObserver)
	observers.Join(&jobRefreshObserver{
		stream: respStream,
		log:    sv.l,
		mu:     new(sync.Mutex),
	})

	projectSpec, err := sv.projectService.GetByName(respStream.Context(), req.ProjectName)
	if err != nil {
		return mapToGRPCErr(err, "unable to get project")
	}

	if err = sv.jobSvc.Refresh(respStream.Context(), projectSpec, namespaceJobNamePairs, observers); err != nil {
		return status.Errorf(codes.Internal, "failed to refresh jobs: \n%s", err.Error())
	}

	sv.l.Info("finished job refresh", "time", time.Since(startTime))
	return nil
}
