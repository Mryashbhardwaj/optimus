package compiler_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/odpf/optimus/compiler"
	"github.com/odpf/optimus/mock"
	"github.com/odpf/optimus/models"
)

func TestJobRunInputCompiler(t *testing.T) {
	ctx := context.Background()
	projectName := "humara-projectSpec"
	projectSpec := models.ProjectSpec{
		ID:   uuid.New(),
		Name: projectName,
		Config: map[string]string{
			"bucket":                 "gs://some_folder",
			"transporterKafkaBroker": "0.0.0.0:9092",
		},
	}
	namespaceSpec := models.NamespaceSpec{
		ID:          uuid.New(),
		Name:        "namespace-1",
		Config:      map[string]string{},
		ProjectSpec: projectSpec,
	}
	execUnit := new(mock.BasePlugin)
	execUnit.On("PluginInfo").Return(&models.PluginInfoResponse{
		Name: "bq",
	}, nil)
	cliMod := new(mock.CLIMod)
	plugin := &models.Plugin{Base: execUnit, CLIMod: cliMod}

	behavior := models.JobSpecBehavior{
		CatchUp:       true,
		DependsOnPast: false,
	}

	schedule := models.JobSpecSchedule{
		StartDate: time.Date(2000, 11, 11, 0, 0, 0, 0, time.UTC),
		Interval:  "* * * * *",
	}
	window := models.JobSpecTaskWindow{
		Size:       time.Hour,
		Offset:     0,
		TruncateTo: "d",
	}

	mockedTimeNow := time.Now()
	scheduledAt := time.Date(2020, 11, 11, 0, 0, 0, 0, time.UTC)
	instanceSpecData := []models.InstanceSpecData{
		{
			Name:  models.ConfigKeyExecutionTime,
			Value: mockedTimeNow.Format(models.InstanceScheduledAtTimeLayout),
			Type:  models.InstanceDataTypeEnv,
		},
		{
			Name:  models.ConfigKeyDstart,
			Value: window.GetStart(scheduledAt).Format(models.InstanceScheduledAtTimeLayout),
			Type:  models.InstanceDataTypeEnv,
		},
		{
			Name:  models.ConfigKeyDend,
			Value: window.GetEnd(scheduledAt).Format(models.InstanceScheduledAtTimeLayout),
			Type:  models.InstanceDataTypeEnv,
		},
	}

	secrets := []models.ProjectSecretItem{
		{
			ID:    uuid.New(),
			Name:  "table_name",
			Value: "secret_table",
			Type:  models.SecretTypeUserDefined,
		},
		{
			ID:    uuid.New(),
			Name:  "bucket",
			Value: "gs://some_secret_bucket",
			Type:  models.SecretTypeUserDefined,
		},
	}
	secretService := new(mock.SecretService)

	pluginRepo := new(mock.SupportedPluginRepo)

	createJobRunInputCompiler := func() compiler.JobRunInputCompiler {
		engine := compiler.NewGoEngine()
		jobConfigCompiler := compiler.NewJobConfigCompiler(engine)
		assetCompiler := compiler.NewJobAssetsCompiler(engine, pluginRepo)
		runInputCompiler := compiler.NewJobRunInputCompiler(secretService, jobConfigCompiler, assetCompiler)
		return runInputCompiler
	}

	t.Run("Compile", func(t *testing.T) {
		t.Run("returns compiled task config for task", func(t *testing.T) {
			jobSpec := models.JobSpec{
				Name:     "foo",
				Owner:    "mee@mee",
				Behavior: behavior,
				Schedule: schedule,
				Task: models.JobSpecTask{
					Unit:     plugin,
					Priority: 2000,
					Window:   window,
					Config: models.JobSpecConfigs{
						{
							Name:  "BQ_VAL",
							Value: "22",
						},
						{
							Name:  "EXECT",
							Value: "{{.EXECUTION_TIME}}",
						},
						{
							Name:  "BUCKET",
							Value: "{{.GLOBAL__bucket}}",
						},
						{
							Name:  "BUCKETX",
							Value: "{{.proj.bucket}}",
						},
						{
							Name:  "LOAD_METHOD",
							Value: "MERGE",
						},
					},
				},
				Dependencies: map[string]models.JobSpecDependency{},
				Assets: *models.JobAssets{}.New(
					[]models.JobSpecAsset{
						{
							Name:  "query.sql",
							Value: "select * from table WHERE event_timestamp > '{{.EXECUTION_TIME}}'",
						},
					},
				),
			}

			jobRun := models.JobRun{
				Spec:        jobSpec,
				Trigger:     models.TriggerSchedule,
				Status:      models.RunStateAccepted,
				ScheduledAt: scheduledAt,
			}
			instanceSpec := models.InstanceSpec{
				Name:   "bq",
				Type:   models.InstanceTypeTask,
				Status: models.RunStateRunning,
				Data:   instanceSpecData,
			}

			cliMod.On("CompileAssets", context.TODO(), models.CompileAssetsRequest{
				Window:           window,
				Config:           models.PluginConfigs{}.FromJobSpec(jobSpec.Task.Config),
				Assets:           models.PluginAssets{}.FromJobSpec(jobSpec.Assets),
				InstanceSchedule: scheduledAt,
				InstanceData:     instanceSpecData,
			}).Return(&models.CompileAssetsResponse{Assets: models.PluginAssets{
				models.PluginAsset{
					Name:  "query.sql",
					Value: "select * from table WHERE event_timestamp > '{{.EXECUTION_TIME}}'",
				},
			}}, nil)
			defer cliMod.AssertExpectations(t)

			pluginRepo.On("GetByName", "bq").Return(plugin, nil)
			defer pluginRepo.AssertExpectations(t)

			secretService.On("GetSecrets", ctx, namespaceSpec).Return(secrets, nil)
			defer secretService.AssertExpectations(t)

			jobRunInputCompiler := createJobRunInputCompiler()
			jobRunInput, err := jobRunInputCompiler.Compile(ctx, namespaceSpec, jobRun, instanceSpec)
			assert.Nil(t, err)

			assert.Equal(t, "2020-11-11T00:00:00Z", jobRunInput.ConfigMap["DEND"])
			assert.Equal(t, "2020-11-10T23:00:00Z", jobRunInput.ConfigMap["DSTART"])
			assert.Equal(t, mockedTimeNow.Format(models.InstanceScheduledAtTimeLayout), jobRunInput.ConfigMap["EXECUTION_TIME"])

			assert.Equal(t, "22", jobRunInput.ConfigMap["BQ_VAL"])
			assert.Equal(t, mockedTimeNow.Format(models.InstanceScheduledAtTimeLayout), jobRunInput.ConfigMap["EXECT"])
			assert.Equal(t, projectSpec.Config["bucket"], jobRunInput.ConfigMap["BUCKET"])
			assert.Equal(t, projectSpec.Config["bucket"], jobRunInput.ConfigMap["BUCKETX"])

			assert.Equal(t,
				fmt.Sprintf("select * from table WHERE event_timestamp > '%s'", mockedTimeNow.Format(models.InstanceScheduledAtTimeLayout)),
				jobRunInput.FileMap["query.sql"],
			)
		})
		t.Run("should return valid compiled instanceSpec config for task type hook", func(t *testing.T) {
			transporterHook := "transporter"
			hookUnit := new(mock.BasePlugin)
			hookUnit.On("PluginInfo").Return(&models.PluginInfoResponse{
				Name:       transporterHook,
				PluginType: models.PluginTypeHook,
			}, nil)

			jobSpec := models.JobSpec{
				Name:     "foo",
				Owner:    "mee@mee",
				Behavior: behavior,
				Schedule: schedule,
				Task: models.JobSpecTask{
					Unit:     &models.Plugin{Base: execUnit, CLIMod: cliMod},
					Priority: 2000,
					Window:   window,
					Config: models.JobSpecConfigs{
						{
							Name:  "BQ_VAL",
							Value: "22",
						},
					},
				},
				Dependencies: map[string]models.JobSpecDependency{},
				Assets: *models.JobAssets{}.New(
					[]models.JobSpecAsset{
						{
							Name:  "query.sql",
							Value: "select * from table WHERE event_timestamp > '{{.EXECUTION_TIME}}'",
						},
					},
				),
				Hooks: []models.JobSpecHook{
					{
						Config: models.JobSpecConfigs{
							{
								Name:  "SAMPLE_CONFIG",
								Value: "200",
							},
							{
								Name:  "INHERIT_CONFIG",
								Value: "{{.TASK__BQ_VAL}}",
							},
							{
								Name:  "INHERIT_CONFIG_AS_WELL",
								Value: "{{.task.BQ_VAL}}",
							},
							{
								Name:  "UNKNOWN",
								Value: "{{.task.TT}}",
							},
							{
								Name:  "FILTER_EXPRESSION",
								Value: "event_timestamp >= '{{.DSTART}}' AND event_timestamp < '{{.DEND}}'",
							},
							{
								Name:  "PRODUCER_CONFIG_BOOTSTRAP_SERVERS",
								Value: `{{.GLOBAL__transporterKafkaBroker}}`,
							},
						},
						Unit: &models.Plugin{Base: hookUnit},
					},
				},
			}
			jobRun := models.JobRun{
				Spec:        jobSpec,
				Trigger:     models.TriggerSchedule,
				Status:      models.RunStateAccepted,
				ScheduledAt: scheduledAt,
			}
			instanceSpec := models.InstanceSpec{
				Name:   transporterHook,
				Type:   models.InstanceTypeHook,
				Status: models.RunStateRunning,
				Data:   instanceSpecData,
			}
			cliMod.On("CompileAssets", context.TODO(), models.CompileAssetsRequest{
				Window:           window,
				Config:           models.PluginConfigs{}.FromJobSpec(jobSpec.Task.Config),
				Assets:           models.PluginAssets{}.FromJobSpec(jobSpec.Assets),
				InstanceSchedule: scheduledAt,
				InstanceData:     instanceSpecData,
			}).Return(&models.CompileAssetsResponse{Assets: models.PluginAssets{
				models.PluginAsset{
					Name:  "query.sql",
					Value: "select * from table WHERE event_timestamp > '{{.EXECUTION_TIME}}'",
				},
			}}, nil)
			defer cliMod.AssertExpectations(t)

			pluginRepo.On("GetByName", "bq").Return(plugin, nil)
			defer pluginRepo.AssertExpectations(t)

			secretService.On("GetSecrets", ctx, namespaceSpec).Return(secrets, nil)
			defer secretService.AssertExpectations(t)

			jobRunInputCompiler := createJobRunInputCompiler()
			jobRunInput, err := jobRunInputCompiler.Compile(ctx, namespaceSpec, jobRun, instanceSpec)
			assert.Nil(t, err)

			assert.Equal(t, "2020-11-11T00:00:00Z", jobRunInput.ConfigMap["DEND"])
			assert.Equal(t, "2020-11-10T23:00:00Z", jobRunInput.ConfigMap["DSTART"])
			assert.Equal(t, mockedTimeNow.Format(models.InstanceScheduledAtTimeLayout), jobRunInput.ConfigMap["EXECUTION_TIME"])

			assert.Equal(t, "0.0.0.0:9092", jobRunInput.ConfigMap["PRODUCER_CONFIG_BOOTSTRAP_SERVERS"])
			assert.Equal(t, "200", jobRunInput.ConfigMap["SAMPLE_CONFIG"])
			assert.Equal(t, "22", jobRunInput.ConfigMap["INHERIT_CONFIG"])
			assert.Equal(t, "22", jobRunInput.ConfigMap["INHERIT_CONFIG_AS_WELL"])
			assert.Equal(t, "<no value>", jobRunInput.ConfigMap["UNKNOWN"])

			assert.Equal(t, "event_timestamp >= '2020-11-10T23:00:00Z' AND event_timestamp < '2020-11-11T00:00:00Z'", jobRunInput.ConfigMap["FILTER_EXPRESSION"])

			assert.Equal(t,
				fmt.Sprintf("select * from table WHERE event_timestamp > '%s'", mockedTimeNow.Format(models.InstanceScheduledAtTimeLayout)),
				jobRunInput.FileMap["query.sql"],
			)
		})
		t.Run("should return compiled instanceSpec config with overridden config provided in NamespaceSpec", func(t *testing.T) {
			projectName := "humara-projectSpec"
			projectSpec := models.ProjectSpec{
				ID:   uuid.Must(uuid.NewRandom()),
				Name: projectName,
				Config: map[string]string{
					"bucket":              "gs://some_folder",
					"transporter_brokers": "129.3.34.1:9092",
				},
			}
			namespaceSpec := models.NamespaceSpec{
				ID:   uuid.Must(uuid.NewRandom()),
				Name: projectName,
				Config: map[string]string{
					"transporter_brokers": "129.3.34.1:9092-overridden",
				},
				ProjectSpec: projectSpec,
			}

			jobSpec := models.JobSpec{
				Name:     "foo",
				Owner:    "mee@mee",
				Behavior: behavior,
				Schedule: schedule,
				Task: models.JobSpecTask{
					Unit:     &models.Plugin{Base: execUnit, CLIMod: cliMod},
					Priority: 2000,
					Window:   window,
					Config: models.JobSpecConfigs{
						{
							Name:  "BQ_VAL",
							Value: "22",
						},
						{
							Name:  "EXECT",
							Value: "{{.EXECUTION_TIME}}",
						},
						{
							Name:  "BUCKET",
							Value: "{{.GLOBAL__bucket}}",
						},
						{
							Name:  "TRANSPORTER_BROKERS",
							Value: "{{.GLOBAL__transporter_brokers}}",
						},
						{
							Name:  "LOAD_METHOD",
							Value: "MERGE",
						},
					},
				},
				Dependencies: map[string]models.JobSpecDependency{},
				Assets: *models.JobAssets{}.New(
					[]models.JobSpecAsset{
						{
							Name:  "query.sql",
							Value: "select * from table WHERE event_timestamp > '{{.EXECUTION_TIME}}'",
						},
					},
				),
			}

			jobRun := models.JobRun{
				Spec:        jobSpec,
				Trigger:     models.TriggerSchedule,
				Status:      models.RunStateAccepted,
				ScheduledAt: scheduledAt,
			}
			instanceSpec := models.InstanceSpec{
				Name:   "bq",
				Type:   models.InstanceTypeTask,
				Status: models.RunStateRunning,
				Data:   instanceSpecData,
			}

			cliMod.On("CompileAssets", context.TODO(), models.CompileAssetsRequest{
				Window:           window,
				Config:           models.PluginConfigs{}.FromJobSpec(jobSpec.Task.Config),
				Assets:           models.PluginAssets{}.FromJobSpec(jobSpec.Assets),
				InstanceSchedule: scheduledAt,
				InstanceData:     instanceSpecData,
			}).Return(&models.CompileAssetsResponse{Assets: models.PluginAssets{
				models.PluginAsset{
					Name:  "query.sql",
					Value: "select * from table WHERE event_timestamp > '{{.EXECUTION_TIME}}'",
				},
			}}, nil)
			defer cliMod.AssertExpectations(t)

			pluginRepo.On("GetByName", "bq").Return(plugin, nil)
			defer pluginRepo.AssertExpectations(t)

			secretService.On("GetSecrets", ctx, namespaceSpec).Return(secrets, nil)
			defer secretService.AssertExpectations(t)

			jobRunInputCompiler := createJobRunInputCompiler()
			jobRunInput, err := jobRunInputCompiler.Compile(ctx, namespaceSpec, jobRun, instanceSpec)
			assert.Nil(t, err)

			assert.Equal(t, "2020-11-11T00:00:00Z", jobRunInput.ConfigMap["DEND"])
			assert.Equal(t, "2020-11-10T23:00:00Z", jobRunInput.ConfigMap["DSTART"])
			assert.Equal(t, mockedTimeNow.Format(models.InstanceScheduledAtTimeLayout), jobRunInput.ConfigMap["EXECUTION_TIME"])

			assert.Equal(t, "22", jobRunInput.ConfigMap["BQ_VAL"])
			assert.Equal(t, mockedTimeNow.Format(models.InstanceScheduledAtTimeLayout), jobRunInput.ConfigMap["EXECT"])
			assert.Equal(t, projectSpec.Config["bucket"], jobRunInput.ConfigMap["BUCKET"])
			assert.Equal(t, namespaceSpec.Config["transporter_brokers"], jobRunInput.ConfigMap["TRANSPORTER_BROKERS"])

			assert.Equal(t,
				fmt.Sprintf("select * from table WHERE event_timestamp > '%s'", mockedTimeNow.Format(models.InstanceScheduledAtTimeLayout)),
				jobRunInput.FileMap["query.sql"],
			)
		})
		t.Run("returns compiled instance spec with secrets", func(t *testing.T) {
			jobSpec := models.JobSpec{
				Name:     "foo",
				Owner:    "mee@mee",
				Behavior: behavior,
				Schedule: schedule,
				Task: models.JobSpecTask{
					Unit:     &models.Plugin{Base: execUnit, CLIMod: cliMod},
					Priority: 2000,
					Window:   window,
					Config: models.JobSpecConfigs{
						{
							Name:  "BQ_VAL",
							Value: "22",
						},
						{
							Name:  "EXECT",
							Value: "{{.EXECUTION_TIME}}",
						},
						{
							Name:  "BUCKET",
							Value: "{{.secret.bucket}}",
						},
						{
							Name:  "LOAD_METHOD",
							Value: "MERGE",
						},
					},
				},
				Dependencies: map[string]models.JobSpecDependency{},
				Assets: *models.JobAssets{}.New(
					[]models.JobSpecAsset{
						{
							Name:  "query.sql",
							Value: "select * from table WHERE event_timestamp > '{{.EXECUTION_TIME}}' and name = '{{.secret.table_name}}'",
						},
					},
				),
			}

			jobRun := models.JobRun{
				Spec:        jobSpec,
				Trigger:     models.TriggerSchedule,
				Status:      models.RunStateAccepted,
				ScheduledAt: scheduledAt,
			}
			instanceSpec := models.InstanceSpec{
				Name:       "bq",
				Type:       models.InstanceTypeTask,
				Status:     models.RunStateRunning,
				Data:       instanceSpecData,
				ExecutedAt: time.Time{},
				UpdatedAt:  time.Time{},
			}

			cliMod.On("CompileAssets", context.TODO(), models.CompileAssetsRequest{
				Window:           window,
				Config:           models.PluginConfigs{}.FromJobSpec(jobSpec.Task.Config),
				Assets:           models.PluginAssets{}.FromJobSpec(jobSpec.Assets),
				InstanceSchedule: scheduledAt,
				InstanceData:     instanceSpecData,
			}).Return(&models.CompileAssetsResponse{Assets: models.PluginAssets{
				models.PluginAsset{
					Name:  "query.sql",
					Value: "select * from table WHERE event_timestamp > '{{.EXECUTION_TIME}}' and name = '{{.secret.table_name}}'",
				},
			}}, nil)
			defer cliMod.AssertExpectations(t)

			pluginRepo.On("GetByName", "bq").Return(plugin, nil)
			defer pluginRepo.AssertExpectations(t)

			secretService.On("GetSecrets", ctx, namespaceSpec).Return(secrets, nil)
			defer secretService.AssertExpectations(t)

			jobRunInputCompiler := createJobRunInputCompiler()
			jobRunInput, err := jobRunInputCompiler.Compile(ctx, namespaceSpec, jobRun, instanceSpec)
			assert.Nil(t, err)

			assert.Equal(t, "2020-11-11T00:00:00Z", jobRunInput.ConfigMap["DEND"])
			assert.Equal(t, "2020-11-10T23:00:00Z", jobRunInput.ConfigMap["DSTART"])
			assert.Equal(t, mockedTimeNow.Format(models.InstanceScheduledAtTimeLayout), jobRunInput.ConfigMap["EXECUTION_TIME"])

			assert.Equal(t, "22", jobRunInput.ConfigMap["BQ_VAL"])
			assert.Equal(t, mockedTimeNow.Format(models.InstanceScheduledAtTimeLayout), jobRunInput.ConfigMap["EXECT"])
			_, ok := jobRunInput.ConfigMap["BUCKET"]
			assert.Equal(t, false, ok)
			assert.Equal(t, "gs://some_secret_bucket", jobRunInput.SecretsMap["BUCKET"])

			assert.Equal(t,
				fmt.Sprintf("select * from table WHERE event_timestamp > '%s' and name = '%s'",
					mockedTimeNow.Format(models.InstanceScheduledAtTimeLayout), secrets[0].Value),
				jobRunInput.FileMap["query.sql"],
			)
		})
	})
}
