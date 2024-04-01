//go:build wireinject
// +build wireinject

/*
 * Copyright (c) 2020 Devtron Labs
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"github.com/devtron-labs/authenticator/middleware"
	util4 "github.com/devtron-labs/common-lib-private/utils/k8s"
	cloudProviderIdentifier "github.com/devtron-labs/common-lib/cloud-provider-identifier"
	pubsub1 "github.com/devtron-labs/common-lib/pubsub-lib"
	"github.com/devtron-labs/devtron/api/apiToken"
	appStoreRestHandler "github.com/devtron-labs/devtron/api/appStore"
	chartGroup2 "github.com/devtron-labs/devtron/api/appStore/chartGroup"
	chartProvider "github.com/devtron-labs/devtron/api/appStore/chartProvider"
	appStoreDeployment "github.com/devtron-labs/devtron/api/appStore/deployment"
	appStoreDiscover "github.com/devtron-labs/devtron/api/appStore/discover"
	appStoreValues "github.com/devtron-labs/devtron/api/appStore/values"
	"github.com/devtron-labs/devtron/api/argoApplication"
	"github.com/devtron-labs/devtron/api/auth/authorisation/globalConfig"
	"github.com/devtron-labs/devtron/api/auth/sso"
	"github.com/devtron-labs/devtron/api/auth/user"
	chartRepo "github.com/devtron-labs/devtron/api/chartRepo"
	"github.com/devtron-labs/devtron/api/cluster"
	"github.com/devtron-labs/devtron/api/connector"
	"github.com/devtron-labs/devtron/api/dashboardEvent"
	"github.com/devtron-labs/devtron/api/deployment"
	"github.com/devtron-labs/devtron/api/devtronResource"
	"github.com/devtron-labs/devtron/api/externalLink"
	"github.com/devtron-labs/devtron/api/globalPolicy"
	client "github.com/devtron-labs/devtron/api/helm-app"
	"github.com/devtron-labs/devtron/api/infraConfig"
	"github.com/devtron-labs/devtron/api/k8s"
	"github.com/devtron-labs/devtron/api/module"
	"github.com/devtron-labs/devtron/api/restHandler"
	"github.com/devtron-labs/devtron/api/restHandler/app/appInfo"
	appList2 "github.com/devtron-labs/devtron/api/restHandler/app/appList"
	"github.com/devtron-labs/devtron/api/restHandler/app/pipeline"
	"github.com/devtron-labs/devtron/api/restHandler/app/pipeline/configure"
	"github.com/devtron-labs/devtron/api/restHandler/app/pipeline/history"
	status2 "github.com/devtron-labs/devtron/api/restHandler/app/pipeline/status"
	"github.com/devtron-labs/devtron/api/restHandler/app/pipeline/trigger"
	"github.com/devtron-labs/devtron/api/restHandler/app/pipeline/webhook"
	"github.com/devtron-labs/devtron/api/restHandler/app/workflow"
	imageDigestPolicy2 "github.com/devtron-labs/devtron/api/restHandler/imageDigestPolicy"
	resourceFilter2 "github.com/devtron-labs/devtron/api/restHandler/resourceFilter"
	"github.com/devtron-labs/devtron/api/restHandler/scopedVariable"
	"github.com/devtron-labs/devtron/api/router"
	app2 "github.com/devtron-labs/devtron/api/router/app"
	appInfo2 "github.com/devtron-labs/devtron/api/router/app/appInfo"
	"github.com/devtron-labs/devtron/api/router/app/appList"
	pipeline2 "github.com/devtron-labs/devtron/api/router/app/pipeline"
	configure2 "github.com/devtron-labs/devtron/api/router/app/pipeline/configure"
	history2 "github.com/devtron-labs/devtron/api/router/app/pipeline/history"
	status3 "github.com/devtron-labs/devtron/api/router/app/pipeline/status"
	trigger2 "github.com/devtron-labs/devtron/api/router/app/pipeline/trigger"
	workflow2 "github.com/devtron-labs/devtron/api/router/app/workflow"
	"github.com/devtron-labs/devtron/api/server"
	"github.com/devtron-labs/devtron/api/sse"
	"github.com/devtron-labs/devtron/api/team"
	"github.com/devtron-labs/devtron/api/terminal"
	util5 "github.com/devtron-labs/devtron/api/util"
	webhookHelm "github.com/devtron-labs/devtron/api/webhook/helm"
	"github.com/devtron-labs/devtron/client/argocdServer"
	"github.com/devtron-labs/devtron/client/argocdServer/application"
	cluster2 "github.com/devtron-labs/devtron/client/argocdServer/cluster"
	"github.com/devtron-labs/devtron/client/argocdServer/connection"
	repository2 "github.com/devtron-labs/devtron/client/argocdServer/repository"
	session2 "github.com/devtron-labs/devtron/client/argocdServer/session"
	"github.com/devtron-labs/devtron/client/cron"
	"github.com/devtron-labs/devtron/client/dashboard"
	eClient "github.com/devtron-labs/devtron/client/events"
	"github.com/devtron-labs/devtron/client/gitSensor"
	"github.com/devtron-labs/devtron/client/grafana"
	"github.com/devtron-labs/devtron/client/lens"
	"github.com/devtron-labs/devtron/client/proxy"
	"github.com/devtron-labs/devtron/client/telemetry"
	"github.com/devtron-labs/devtron/enterprise/api/artifactPromotionApprovalRequest"
	"github.com/devtron-labs/devtron/enterprise/api/artifactPromotionPolicy"
	"github.com/devtron-labs/devtron/enterprise/api/commonPolicyActions"
	"github.com/devtron-labs/devtron/enterprise/api/drafts"
	"github.com/devtron-labs/devtron/enterprise/api/globalTag"
	"github.com/devtron-labs/devtron/enterprise/api/lockConfiguation"
	"github.com/devtron-labs/devtron/enterprise/api/protect"
	app3 "github.com/devtron-labs/devtron/enterprise/pkg/app"
	pipeline3 "github.com/devtron-labs/devtron/enterprise/pkg/pipeline"
	"github.com/devtron-labs/devtron/enterprise/pkg/resourceFilter"
	"github.com/devtron-labs/devtron/internal/sql/repository"
	app4 "github.com/devtron-labs/devtron/internal/sql/repository/app"
	appStatusRepo "github.com/devtron-labs/devtron/internal/sql/repository/appStatus"
	appWorkflow2 "github.com/devtron-labs/devtron/internal/sql/repository/appWorkflow"
	"github.com/devtron-labs/devtron/internal/sql/repository/bulkUpdate"
	"github.com/devtron-labs/devtron/internal/sql/repository/chartConfig"
	dockerRegistryRepository "github.com/devtron-labs/devtron/internal/sql/repository/dockerRegistry"
	"github.com/devtron-labs/devtron/internal/sql/repository/helper"
	repository8 "github.com/devtron-labs/devtron/internal/sql/repository/imageTagging"
	"github.com/devtron-labs/devtron/internal/sql/repository/pipelineConfig"
	resourceGroup "github.com/devtron-labs/devtron/internal/sql/repository/resourceGroup"
	security2 "github.com/devtron-labs/devtron/internal/sql/repository/security"
	"github.com/devtron-labs/devtron/internal/util"
	"github.com/devtron-labs/devtron/pkg/app"
	"github.com/devtron-labs/devtron/pkg/app/status"
	"github.com/devtron-labs/devtron/pkg/appClone"
	"github.com/devtron-labs/devtron/pkg/appClone/batch"
	"github.com/devtron-labs/devtron/pkg/appStatus"
	"github.com/devtron-labs/devtron/pkg/appStore/chartGroup"
	repository4 "github.com/devtron-labs/devtron/pkg/appStore/chartGroup/repository"
	"github.com/devtron-labs/devtron/pkg/appStore/installedApp/service/FullMode"
	deployment3 "github.com/devtron-labs/devtron/pkg/appStore/installedApp/service/FullMode/deployment"
	"github.com/devtron-labs/devtron/pkg/appStore/installedApp/service/FullMode/deploymentTypeChange"
	"github.com/devtron-labs/devtron/pkg/appStore/installedApp/service/FullMode/resource"
	"github.com/devtron-labs/devtron/pkg/appWorkflow"
	"github.com/devtron-labs/devtron/pkg/appWorkflow/read"
	"github.com/devtron-labs/devtron/pkg/attributes"
	client2 "github.com/devtron-labs/devtron/pkg/auth/authorisation/casbin"
	"github.com/devtron-labs/devtron/pkg/build"
	"github.com/devtron-labs/devtron/pkg/bulkAction"
	"github.com/devtron-labs/devtron/pkg/chart"
	"github.com/devtron-labs/devtron/pkg/chart/gitOpsConfig"
	chartRepoRepository "github.com/devtron-labs/devtron/pkg/chartRepo/repository"
	"github.com/devtron-labs/devtron/pkg/commonService"
	delete2 "github.com/devtron-labs/devtron/pkg/delete"
	deployment2 "github.com/devtron-labs/devtron/pkg/deployment"
	git2 "github.com/devtron-labs/devtron/pkg/deployment/gitOps/git"
	"github.com/devtron-labs/devtron/pkg/deploymentGroup"
	"github.com/devtron-labs/devtron/pkg/dockerRegistry"
	"github.com/devtron-labs/devtron/pkg/eventProcessor"
	"github.com/devtron-labs/devtron/pkg/generateManifest"
	"github.com/devtron-labs/devtron/pkg/git"
	"github.com/devtron-labs/devtron/pkg/gitops"
	globalPolicy2 "github.com/devtron-labs/devtron/pkg/globalPolicy"
	"github.com/devtron-labs/devtron/pkg/imageDigestPolicy"
	infraConfigService "github.com/devtron-labs/devtron/pkg/infraConfig"
	"github.com/devtron-labs/devtron/pkg/infraConfig/units"
	"github.com/devtron-labs/devtron/pkg/kubernetesResourceAuditLogs"
	repository7 "github.com/devtron-labs/devtron/pkg/kubernetesResourceAuditLogs/repository"
	"github.com/devtron-labs/devtron/pkg/notifier"
	pipeline4 "github.com/devtron-labs/devtron/pkg/pipeline"
	"github.com/devtron-labs/devtron/pkg/pipeline/executors"
	history3 "github.com/devtron-labs/devtron/pkg/pipeline/history"
	repository3 "github.com/devtron-labs/devtron/pkg/pipeline/history/repository"
	"github.com/devtron-labs/devtron/pkg/pipeline/infraProviders"
	repository5 "github.com/devtron-labs/devtron/pkg/pipeline/repository"
	"github.com/devtron-labs/devtron/pkg/pipeline/types"
	"github.com/devtron-labs/devtron/pkg/plugin"
	repository6 "github.com/devtron-labs/devtron/pkg/plugin/repository"
	"github.com/devtron-labs/devtron/pkg/policyGovernance/artifactApproval"
	artifactPromotion2 "github.com/devtron-labs/devtron/pkg/policyGovernance/artifactPromotion"
	resourceGroup2 "github.com/devtron-labs/devtron/pkg/resourceGroup"
	"github.com/devtron-labs/devtron/pkg/resourceQualifiers"
	"github.com/devtron-labs/devtron/pkg/security"
	"github.com/devtron-labs/devtron/pkg/sql"
	"github.com/devtron-labs/devtron/pkg/timeoutWindow"
	repository9 "github.com/devtron-labs/devtron/pkg/timeoutWindow/repository"
	util3 "github.com/devtron-labs/devtron/pkg/util"
	"github.com/devtron-labs/devtron/pkg/variables"
	"github.com/devtron-labs/devtron/pkg/variables/parsers"
	repository10 "github.com/devtron-labs/devtron/pkg/variables/repository"
	workflow3 "github.com/devtron-labs/devtron/pkg/workflow"
	"github.com/devtron-labs/devtron/pkg/workflow/dag"
	util2 "github.com/devtron-labs/devtron/util"
	"github.com/devtron-labs/devtron/util/argo"
	cron2 "github.com/devtron-labs/devtron/util/cron"
	"github.com/devtron-labs/devtron/util/rbac"
	"github.com/google/wire"
)

func InitializeApp() (*App, error) {

	wire.Build(
		// ----- wireset start
		commonPolicyActions.CommonPolicyActionWireSet,
		sql.PgSqlWireSet,
		user.SelfRegistrationWireSet,
		externalLink.ExternalLinkWireSet,
		team.TeamsWireSet,
		AuthWireSet,
		util4.NewK8sUtilExtended,
		user.UserWireSet,
		sso.SsoConfigWireSet,
		cluster.ClusterWireSet,
		dashboard.DashboardWireSet,
		proxy.ProxyWireSet,
		client.HelmAppWireSet,
		k8s.K8sApplicationWireSet,
		chartRepo.ChartRepositoryWireSet,
		appStoreDiscover.AppStoreDiscoverWireSet,
		chartProvider.AppStoreChartProviderWireSet,
		appStoreValues.AppStoreValuesWireSet,
		appStoreDeployment.AppStoreDeploymentWireSet,
		server.ServerWireSet,
		module.ModuleWireSet,
		apiToken.ApiTokenWireSet,
		webhookHelm.WebhookHelmWireSet,
		terminal.TerminalWireSet,
		client2.CasbinWireSet,
		globalTag.GlobalTagWireSet,
		globalPolicy.GlobalPolicyWireSet,
		drafts.DraftsWireSet,
		protect.ProtectWireSet,
		devtronResource.DevtronResourceWireSet,
		globalConfig.GlobalConfigWireSet,
		lockConfiguation.LockConfigWireSet,
		build.BuildWireSet,
		deployment2.DeploymentWireSet,
		argoApplication.ArgoApplicationWireSet,

		eventProcessor.EventProcessorWireSet,
		workflow3.WorkflowWireSet,
		artifactApproval.ArtifactApprovalWireSet,
		artifactPromotion2.ArtifactPromotionWireSet,
		// -------wireset end ----------
		// -------
		gitSensor.GetConfig,
		gitSensor.NewGitSensorClient,
		wire.Bind(new(gitSensor.Client), new(*gitSensor.ClientImpl)),
		// -------
		helper.NewAppListingRepositoryQueryBuilder,
		// sql.GetConfig,
		eClient.GetEventClientConfig,
		util2.GetGlobalEnvVariables,
		// sql.NewDbConnection,
		// app.GetACDAuthConfig,
		util3.GetACDAuthConfig,
		connection.SettingsManager,
		// auth.GetConfig,

		connection.GetConfig,
		wire.Bind(new(session2.ServiceClient), new(*middleware.LoginService)),

		sse.NewSSE,
		trigger2.NewPipelineTriggerRouter,
		wire.Bind(new(trigger2.PipelineTriggerRouter), new(*trigger2.PipelineTriggerRouterImpl)),

		// ---- pprof start ----
		restHandler.NewPProfRestHandler,
		wire.Bind(new(restHandler.PProfRestHandler), new(*restHandler.PProfRestHandlerImpl)),

		router.NewPProfRouter,
		wire.Bind(new(router.PProfRouter), new(*router.PProfRouterImpl)),
		// ---- pprof end ----

		trigger.NewPipelineRestHandler,
		wire.Bind(new(trigger.PipelineTriggerRestHandler), new(*trigger.PipelineTriggerRestHandlerImpl)),
		app.GetAppServiceConfig,
		app.NewAppService,
		wire.Bind(new(app.AppService), new(*app.AppServiceImpl)),

		bulkUpdate.NewBulkUpdateRepository,
		wire.Bind(new(bulkUpdate.BulkUpdateRepository), new(*bulkUpdate.BulkUpdateRepositoryImpl)),

		chartConfig.NewEnvConfigOverrideRepository,
		wire.Bind(new(chartConfig.EnvConfigOverrideRepository), new(*chartConfig.EnvConfigOverrideRepositoryImpl)),
		chartConfig.NewPipelineOverrideRepository,
		wire.Bind(new(chartConfig.PipelineOverrideRepository), new(*chartConfig.PipelineOverrideRepositoryImpl)),
		wire.Struct(new(util.MergeUtil), "*"),
		util.NewSugardLogger,

		deployment.NewDeploymentConfigRestHandlerImpl,
		wire.Bind(new(deployment.DeploymentConfigRestHandler), new(*deployment.DeploymentConfigRestHandlerImpl)),
		deployment.NewDeploymentRouterImpl,
		wire.Bind(new(deployment.DeploymentConfigRouter), new(*deployment.DeploymentConfigRouterImpl)),

		dashboardEvent.NewDashboardTelemetryRestHandlerImpl,
		wire.Bind(new(dashboardEvent.DashboardTelemetryRestHandler), new(*dashboardEvent.DashboardTelemetryRestHandlerImpl)),
		dashboardEvent.NewDashboardTelemetryRouterImpl,
		wire.Bind(new(dashboardEvent.DashboardTelemetryRouter),
			new(*dashboardEvent.DashboardTelemetryRouterImpl)),

		infraConfigService.NewInfraProfileRepositoryImpl,
		wire.Bind(new(infraConfigService.InfraConfigRepository), new(*infraConfigService.InfraConfigRepositoryImpl)),

		units.NewUnits,
		infraConfigService.NewValidatorImpl,
		wire.Bind(new(infraConfigService.Validator), new(*infraConfigService.ValidatorImpl)),
		infraConfigService.NewInfraConfigServiceImpl,
		wire.Bind(new(infraConfigService.InfraConfigService), new(*infraConfigService.InfraConfigServiceImpl)),
		infraProviders.NewInfraProviderImpl,
		wire.Bind(new(infraProviders.InfraProvider), new(*infraProviders.InfraProviderImpl)),

		infraConfig.NewInfraConfigRestHandlerImpl,
		wire.Bind(new(infraConfig.InfraConfigRestHandler), new(*infraConfig.InfraConfigRestHandlerImpl)),

		infraConfig.NewInfraProfileRouterImpl,
		wire.Bind(new(infraConfig.InfraConfigRouter), new(*infraConfig.InfraConfigRouterImpl)),

		router.NewMuxRouter,

		app4.NewAppRepositoryImpl,
		wire.Bind(new(app4.AppRepository), new(*app4.AppRepositoryImpl)),

		pipeline4.GetDeploymentServiceTypeConfig,

		pipeline4.NewPipelineBuilderImpl,
		wire.Bind(new(pipeline4.PipelineBuilder), new(*pipeline4.PipelineBuilderImpl)),
		pipeline4.NewBuildPipelineSwitchServiceImpl,
		wire.Bind(new(pipeline4.BuildPipelineSwitchService), new(*pipeline4.BuildPipelineSwitchServiceImpl)),
		pipeline4.NewCiPipelineConfigServiceImpl,
		wire.Bind(new(pipeline4.CiPipelineConfigService), new(*pipeline4.CiPipelineConfigServiceImpl)),
		pipeline4.NewCiMaterialConfigServiceImpl,
		wire.Bind(new(pipeline4.CiMaterialConfigService), new(*pipeline4.CiMaterialConfigServiceImpl)),

		pipeline4.NewAppArtifactManagerImpl,
		wire.Bind(new(pipeline4.AppArtifactManager), new(*pipeline4.AppArtifactManagerImpl)),
		pipeline4.NewDevtronAppCMCSServiceImpl,
		wire.Bind(new(pipeline4.DevtronAppCMCSService), new(*pipeline4.DevtronAppCMCSServiceImpl)),
		pipeline4.NewDevtronAppStrategyServiceImpl,
		wire.Bind(new(pipeline4.DevtronAppStrategyService), new(*pipeline4.DevtronAppStrategyServiceImpl)),
		pipeline4.NewAppDeploymentTypeChangeManagerImpl,
		wire.Bind(new(pipeline4.AppDeploymentTypeChangeManager), new(*pipeline4.AppDeploymentTypeChangeManagerImpl)),
		pipeline4.NewCdPipelineConfigServiceImpl,
		wire.Bind(new(pipeline4.CdPipelineConfigService), new(*pipeline4.CdPipelineConfigServiceImpl)),
		pipeline4.NewDevtronAppConfigServiceImpl,
		wire.Bind(new(pipeline4.DevtronAppConfigService), new(*pipeline4.DevtronAppConfigServiceImpl)),
		pipeline.NewDevtronAppAutoCompleteRestHandlerImpl,
		wire.Bind(new(pipeline.DevtronAppAutoCompleteRestHandler), new(*pipeline.DevtronAppAutoCompleteRestHandlerImpl)),

		util5.NewLoggingMiddlewareImpl,
		wire.Bind(new(util5.LoggingMiddleware), new(*util5.LoggingMiddlewareImpl)),
		configure.NewPipelineRestHandlerImpl,
		wire.Bind(new(configure.PipelineConfigRestHandler), new(*configure.PipelineConfigRestHandlerImpl)),
		// -----------
		configure2.NewPipelineRouterImpl,
		wire.Bind(new(configure2.PipelineConfigRouter), new(*configure2.PipelineConfigRouterImpl)),
		history2.NewPipelineHistoryRouterImpl,
		wire.Bind(new(history2.PipelineHistoryRouter), new(*history2.PipelineHistoryRouterImpl)),
		status3.NewPipelineStatusRouterImpl,
		wire.Bind(new(status3.PipelineStatusRouter), new(*status3.PipelineStatusRouterImpl)),
		pipeline2.NewDevtronAppAutoCompleteRouterImpl,
		wire.Bind(new(pipeline2.DevtronAppAutoCompleteRouter), new(*pipeline2.DevtronAppAutoCompleteRouterImpl)),
		workflow2.NewAppWorkflowRouterImpl,
		wire.Bind(new(workflow2.AppWorkflowRouter), new(*workflow2.AppWorkflowRouterImpl)),

		pipeline4.NewCiCdPipelineOrchestrator,
		pipeline3.NewCiCdPipelineOrchestratorEnterpriseImpl,
		wire.Bind(new(pipeline4.CiCdPipelineOrchestrator), new(*pipeline3.CiCdPipelineOrchestratorEnterpriseImpl)),
		// ------------
		pipelineConfig.NewMaterialRepositoryImpl,
		wire.Bind(new(pipelineConfig.MaterialRepository), new(*pipelineConfig.MaterialRepositoryImpl)),

		util.NewChartTemplateServiceImpl,
		wire.Bind(new(util.ChartTemplateService), new(*util.ChartTemplateServiceImpl)),

		// scoped variables start
		variables.NewScopedVariableServiceImpl,
		wire.Bind(new(variables.ScopedVariableService), new(*variables.ScopedVariableServiceImpl)),

		parsers.NewVariableTemplateParserImpl,
		wire.Bind(new(parsers.VariableTemplateParser), new(*parsers.VariableTemplateParserImpl)),
		repository10.NewVariableEntityMappingRepository,
		wire.Bind(new(repository10.VariableEntityMappingRepository), new(*repository10.VariableEntityMappingRepositoryImpl)),

		repository10.NewVariableSnapshotHistoryRepository,
		wire.Bind(new(repository10.VariableSnapshotHistoryRepository), new(*repository10.VariableSnapshotHistoryRepositoryImpl)),
		variables.NewVariableEntityMappingServiceImpl,
		wire.Bind(new(variables.VariableEntityMappingService), new(*variables.VariableEntityMappingServiceImpl)),
		variables.NewVariableSnapshotHistoryServiceImpl,
		wire.Bind(new(variables.VariableSnapshotHistoryService), new(*variables.VariableSnapshotHistoryServiceImpl)),

		variables.NewScopedVariableManagerImpl,
		wire.Bind(new(variables.ScopedVariableManager), new(*variables.ScopedVariableManagerImpl)),

		variables.NewScopedVariableCMCSManagerImpl,
		wire.Bind(new(variables.ScopedVariableCMCSManager), new(*variables.ScopedVariableCMCSManagerImpl)),

		// end

		gitOpsConfig.NewDevtronAppGitOpConfigServiceImpl,
		wire.Bind(new(gitOpsConfig.DevtronAppGitOpConfigService), new(*gitOpsConfig.DevtronAppGitOpConfigServiceImpl)),
		chart.NewChartServiceImpl,
		wire.Bind(new(chart.ChartService), new(*chart.ChartServiceImpl)),
		bulkAction.NewBulkUpdateServiceImpl,
		wire.Bind(new(bulkAction.BulkUpdateService), new(*bulkAction.BulkUpdateServiceImpl)),

		repository.NewImageTagRepository,
		wire.Bind(new(repository.ImageTagRepository), new(*repository.ImageTagRepositoryImpl)),

		pipeline4.NewCustomTagService,
		wire.Bind(new(pipeline4.CustomTagService), new(*pipeline4.CustomTagServiceImpl)),

		repository.NewGitProviderRepositoryImpl,
		wire.Bind(new(repository.GitProviderRepository), new(*repository.GitProviderRepositoryImpl)),
		pipeline4.NewGitRegistryConfigImpl,
		wire.Bind(new(pipeline4.GitRegistryConfig), new(*pipeline4.GitRegistryConfigImpl)),

		appList.NewAppFilteringRouterImpl,
		wire.Bind(new(appList.AppFilteringRouter), new(*appList.AppFilteringRouterImpl)),
		appList2.NewAppFilteringRestHandlerImpl,
		wire.Bind(new(appList2.AppFilteringRestHandler), new(*appList2.AppFilteringRestHandlerImpl)),

		appList.NewAppListingRouterImpl,
		wire.Bind(new(appList.AppListingRouter), new(*appList.AppListingRouterImpl)),
		appList2.NewAppListingRestHandlerImpl,
		wire.Bind(new(appList2.AppListingRestHandler), new(*appList2.AppListingRestHandlerImpl)),
		app.NewAppListingServiceImpl,
		wire.Bind(new(app.AppListingService), new(*app.AppListingServiceImpl)),
		repository.NewAppListingRepositoryImpl,
		wire.Bind(new(repository.AppListingRepository), new(*repository.AppListingRepositoryImpl)),

		repository.NewDeploymentTemplateRepositoryImpl,
		wire.Bind(new(repository.DeploymentTemplateRepository), new(*repository.DeploymentTemplateRepositoryImpl)),
		generateManifest.NewDeploymentTemplateServiceImpl,
		wire.Bind(new(generateManifest.DeploymentTemplateService), new(*generateManifest.DeploymentTemplateServiceImpl)),

		router.NewJobRouterImpl,
		wire.Bind(new(router.JobRouter), new(*router.JobRouterImpl)),

		pipelineConfig.NewPipelineRepositoryImpl,
		wire.Bind(new(pipelineConfig.PipelineRepository), new(*pipelineConfig.PipelineRepositoryImpl)),
		pipeline4.NewPropertiesConfigServiceImpl,
		wire.Bind(new(pipeline4.PropertiesConfigService), new(*pipeline4.PropertiesConfigServiceImpl)),

		util.NewHttpClient,

		eClient.NewEventRESTClientImpl,
		wire.Bind(new(eClient.EventClient), new(*eClient.EventRESTClientImpl)),

		eClient.NewEventSimpleFactoryImpl,
		wire.Bind(new(eClient.EventFactory), new(*eClient.EventSimpleFactoryImpl)),

		repository.NewCiArtifactRepositoryImpl,
		wire.Bind(new(repository.CiArtifactRepository), new(*repository.CiArtifactRepositoryImpl)),
		pipeline4.NewWebhookServiceImpl,
		wire.Bind(new(pipeline4.WebhookService), new(*pipeline4.WebhookServiceImpl)),

		router.NewWebhookRouterImpl,
		wire.Bind(new(router.WebhookRouter), new(*router.WebhookRouterImpl)),
		pipelineConfig.NewCiTemplateRepositoryImpl,
		wire.Bind(new(pipelineConfig.CiTemplateRepository), new(*pipelineConfig.CiTemplateRepositoryImpl)),
		pipelineConfig.NewCiPipelineRepositoryImpl,
		wire.Bind(new(pipelineConfig.CiPipelineRepository), new(*pipelineConfig.CiPipelineRepositoryImpl)),
		pipelineConfig.NewCiPipelineMaterialRepositoryImpl,
		wire.Bind(new(pipelineConfig.CiPipelineMaterialRepository), new(*pipelineConfig.CiPipelineMaterialRepositoryImpl)),
		git2.NewGitFactory,

		application.NewApplicationClientImpl,
		wire.Bind(new(application.ServiceClient), new(*application.ServiceClientImpl)),
		cluster2.NewServiceClientImpl,
		wire.Bind(new(cluster2.ServiceClient), new(*cluster2.ServiceClientImpl)),
		connector.NewPumpImpl,
		repository2.NewServiceClientImpl,
		wire.Bind(new(repository2.ServiceClient), new(*repository2.ServiceClientImpl)),
		wire.Bind(new(connector.Pump), new(*connector.PumpImpl)),

		// app.GetConfig,

		pipeline4.GetEcrConfig,
		// otel.NewOtelTracingServiceImpl,
		// wire.Bind(new(otel.OtelTracingService), new(*otel.OtelTracingServiceImpl)),
		NewApp,
		// session.NewK8sClient,
		repository8.NewImageTaggingRepositoryImpl,
		wire.Bind(new(repository8.ImageTaggingRepository), new(*repository8.ImageTaggingRepositoryImpl)),
		pipeline4.NewImageTaggingServiceImpl,
		wire.Bind(new(pipeline4.ImageTaggingService), new(*pipeline4.ImageTaggingServiceImpl)),
		argocdServer.NewVersionServiceImpl,
		wire.Bind(new(argocdServer.VersionService), new(*argocdServer.VersionServiceImpl)),

		router.NewGitProviderRouterImpl,
		wire.Bind(new(router.GitProviderRouter), new(*router.GitProviderRouterImpl)),
		restHandler.NewGitProviderRestHandlerImpl,
		wire.Bind(new(restHandler.GitProviderRestHandler), new(*restHandler.GitProviderRestHandlerImpl)),

		router.NewNotificationRouterImpl,
		wire.Bind(new(router.NotificationRouter), new(*router.NotificationRouterImpl)),
		restHandler.NewNotificationRestHandlerImpl,
		wire.Bind(new(restHandler.NotificationRestHandler), new(*restHandler.NotificationRestHandlerImpl)),

		notifier.NewSlackNotificationServiceImpl,
		wire.Bind(new(notifier.SlackNotificationService), new(*notifier.SlackNotificationServiceImpl)),
		repository.NewSlackNotificationRepositoryImpl,
		wire.Bind(new(repository.SlackNotificationRepository), new(*repository.SlackNotificationRepositoryImpl)),
		notifier.NewWebhookNotificationServiceImpl,
		wire.Bind(new(notifier.WebhookNotificationService), new(*notifier.WebhookNotificationServiceImpl)),
		repository.NewWebhookNotificationRepositoryImpl,
		wire.Bind(new(repository.WebhookNotificationRepository), new(*repository.WebhookNotificationRepositoryImpl)),

		notifier.NewNotificationConfigServiceImpl,
		wire.Bind(new(notifier.NotificationConfigService), new(*notifier.NotificationConfigServiceImpl)),
		app.NewAppListingViewBuilderImpl,
		wire.Bind(new(app.AppListingViewBuilder), new(*app.AppListingViewBuilderImpl)),
		repository.NewNotificationSettingsRepositoryImpl,
		wire.Bind(new(repository.NotificationSettingsRepository), new(*repository.NotificationSettingsRepositoryImpl)),
		util.IntValidator,
		types.GetCiCdConfig,

		pipeline4.NewWorkflowServiceImpl,
		wire.Bind(new(pipeline4.WorkflowService), new(*pipeline4.WorkflowServiceImpl)),

		pipeline4.NewCiServiceImpl,
		wire.Bind(new(pipeline4.CiService), new(*pipeline4.CiServiceImpl)),

		pipelineConfig.NewCiWorkflowRepositoryImpl,
		wire.Bind(new(pipelineConfig.CiWorkflowRepository), new(*pipelineConfig.CiWorkflowRepositoryImpl)),

		restHandler.NewGitWebhookRestHandlerImpl,
		wire.Bind(new(restHandler.GitWebhookRestHandler), new(*restHandler.GitWebhookRestHandlerImpl)),

		git.NewGitWebhookServiceImpl,
		wire.Bind(new(git.GitWebhookService), new(*git.GitWebhookServiceImpl)),

		repository.NewGitWebhookRepositoryImpl,
		wire.Bind(new(repository.GitWebhookRepository), new(*repository.GitWebhookRepositoryImpl)),

		pipeline4.NewCiHandlerImpl,
		wire.Bind(new(pipeline4.CiHandler), new(*pipeline4.CiHandlerImpl)),

		pipeline4.NewCiLogServiceImpl,
		wire.Bind(new(pipeline4.CiLogService), new(*pipeline4.CiLogServiceImpl)),

		pubsub1.NewPubSubClientServiceImpl,

		rbac.NewEnforcerUtilImpl,
		wire.Bind(new(rbac.EnforcerUtil), new(*rbac.EnforcerUtilImpl)),

		chartConfig.NewPipelineConfigRepository,
		wire.Bind(new(chartConfig.PipelineConfigRepository), new(*chartConfig.PipelineConfigRepositoryImpl)),

		repository10.NewScopedVariableRepository,
		wire.Bind(new(repository10.ScopedVariableRepository), new(*repository10.ScopedVariableRepositoryImpl)),

		repository.NewLinkoutsRepositoryImpl,
		wire.Bind(new(repository.LinkoutsRepository), new(*repository.LinkoutsRepositoryImpl)),

		router.NewChartRefRouterImpl,
		wire.Bind(new(router.ChartRefRouter), new(*router.ChartRefRouterImpl)),
		restHandler.NewChartRefRestHandlerImpl,
		wire.Bind(new(restHandler.ChartRefRestHandler), new(*restHandler.ChartRefRestHandlerImpl)),

		router.NewConfigMapRouterImpl,
		wire.Bind(new(router.ConfigMapRouter), new(*router.ConfigMapRouterImpl)),
		restHandler.NewConfigMapRestHandlerImpl,
		wire.Bind(new(restHandler.ConfigMapRestHandler), new(*restHandler.ConfigMapRestHandlerImpl)),
		pipeline4.NewConfigMapServiceImpl,
		wire.Bind(new(pipeline4.ConfigMapService), new(*pipeline4.ConfigMapServiceImpl)),
		chartConfig.NewConfigMapRepositoryImpl,
		wire.Bind(new(chartConfig.ConfigMapRepository), new(*chartConfig.ConfigMapRepositoryImpl)),

		notifier.NewSESNotificationServiceImpl,
		wire.Bind(new(notifier.SESNotificationService), new(*notifier.SESNotificationServiceImpl)),

		repository.NewSESNotificationRepositoryImpl,
		wire.Bind(new(repository.SESNotificationRepository), new(*repository.SESNotificationRepositoryImpl)),

		notifier.NewSMTPNotificationServiceImpl,
		wire.Bind(new(notifier.SMTPNotificationService), new(*notifier.SMTPNotificationServiceImpl)),

		repository.NewSMTPNotificationRepositoryImpl,
		wire.Bind(new(repository.SMTPNotificationRepository), new(*repository.SMTPNotificationRepositoryImpl)),

		notifier.NewNotificationConfigBuilderImpl,
		wire.Bind(new(notifier.NotificationConfigBuilder), new(*notifier.NotificationConfigBuilderImpl)),
		appStoreRestHandler.NewAppStoreStatusTimelineRestHandlerImpl,
		wire.Bind(new(appStoreRestHandler.AppStoreStatusTimelineRestHandler), new(*appStoreRestHandler.AppStoreStatusTimelineRestHandlerImpl)),
		appStoreRestHandler.NewInstalledAppRestHandlerImpl,
		wire.Bind(new(appStoreRestHandler.InstalledAppRestHandler), new(*appStoreRestHandler.InstalledAppRestHandlerImpl)),
		FullMode.NewInstalledAppDBExtendedServiceImpl,
		wire.Bind(new(FullMode.InstalledAppDBExtendedService), new(*FullMode.InstalledAppDBExtendedServiceImpl)),
		resource.NewInstalledAppResourceServiceImpl,
		wire.Bind(new(resource.InstalledAppResourceService), new(*resource.InstalledAppResourceServiceImpl)),
		deploymentTypeChange.NewInstalledAppDeploymentTypeChangeServiceImpl,
		wire.Bind(new(deploymentTypeChange.InstalledAppDeploymentTypeChangeService), new(*deploymentTypeChange.InstalledAppDeploymentTypeChangeServiceImpl)),

		appStoreRestHandler.NewAppStoreRouterImpl,
		wire.Bind(new(appStoreRestHandler.AppStoreRouter), new(*appStoreRestHandler.AppStoreRouterImpl)),

		workflow.NewAppWorkflowRestHandlerImpl,
		wire.Bind(new(workflow.AppWorkflowRestHandler), new(*workflow.AppWorkflowRestHandlerImpl)),

		appWorkflow.NewAppWorkflowServiceImpl,
		wire.Bind(new(appWorkflow.AppWorkflowService), new(*appWorkflow.AppWorkflowServiceImpl)),

		read.NewAppWorkflowDataReadServiceImpl,
		wire.Bind(new(read.AppWorkflowDataReadService), new(*read.AppWorkflowDataReadServiceImpl)),

		appWorkflow2.NewAppWorkflowRepositoryImpl,
		wire.Bind(new(appWorkflow2.AppWorkflowRepository), new(*appWorkflow2.AppWorkflowRepositoryImpl)),

		restHandler.NewExternalCiRestHandlerImpl,
		wire.Bind(new(restHandler.ExternalCiRestHandler), new(*restHandler.ExternalCiRestHandlerImpl)),

		grafana.GetGrafanaClientConfig,
		grafana.NewGrafanaClientImpl,
		wire.Bind(new(grafana.GrafanaClient), new(*grafana.GrafanaClientImpl)),

		app.NewReleaseDataServiceImpl,
		wire.Bind(new(app.ReleaseDataService), new(*app.ReleaseDataServiceImpl)),
		restHandler.NewReleaseMetricsRestHandlerImpl,
		wire.Bind(new(restHandler.ReleaseMetricsRestHandler), new(*restHandler.ReleaseMetricsRestHandlerImpl)),
		router.NewReleaseMetricsRouterImpl,
		wire.Bind(new(router.ReleaseMetricsRouter), new(*router.ReleaseMetricsRouterImpl)),
		lens.GetLensConfig,
		lens.NewLensClientImpl,
		wire.Bind(new(lens.LensClient), new(*lens.LensClientImpl)),

		pipelineConfig.NewCdWorkflowRepositoryImpl,
		wire.Bind(new(pipelineConfig.CdWorkflowRepository), new(*pipelineConfig.CdWorkflowRepositoryImpl)),

		pipeline4.NewCdHandlerImpl,
		wire.Bind(new(pipeline4.CdHandler), new(*pipeline4.CdHandlerImpl)),

		pipeline4.NewBlobStorageConfigServiceImpl,
		wire.Bind(new(pipeline4.BlobStorageConfigService), new(*pipeline4.BlobStorageConfigServiceImpl)),

		dag.NewWorkflowDagExecutorImpl,
		wire.Bind(new(dag.WorkflowDagExecutor), new(*dag.WorkflowDagExecutorImpl)),
		appClone.NewAppCloneServiceImpl,
		wire.Bind(new(appClone.AppCloneService), new(*appClone.AppCloneServiceImpl)),

		router.NewDeploymentGroupRouterImpl,
		wire.Bind(new(router.DeploymentGroupRouter), new(*router.DeploymentGroupRouterImpl)),
		restHandler.NewDeploymentGroupRestHandlerImpl,
		wire.Bind(new(restHandler.DeploymentGroupRestHandler), new(*restHandler.DeploymentGroupRestHandlerImpl)),
		deploymentGroup.NewDeploymentGroupServiceImpl,
		wire.Bind(new(deploymentGroup.DeploymentGroupService), new(*deploymentGroup.DeploymentGroupServiceImpl)),
		repository.NewDeploymentGroupRepositoryImpl,
		wire.Bind(new(repository.DeploymentGroupRepository), new(*repository.DeploymentGroupRepositoryImpl)),

		repository.NewDeploymentGroupAppRepositoryImpl,
		wire.Bind(new(repository.DeploymentGroupAppRepository), new(*repository.DeploymentGroupAppRepositoryImpl)),
		restHandler.NewPubSubClientRestHandlerImpl,
		wire.Bind(new(restHandler.PubSubClientRestHandler), new(*restHandler.PubSubClientRestHandlerImpl)),

		// Batch actions
		batch.NewWorkflowActionImpl,
		wire.Bind(new(batch.WorkflowAction), new(*batch.WorkflowActionImpl)),
		batch.NewDeploymentActionImpl,
		wire.Bind(new(batch.DeploymentAction), new(*batch.DeploymentActionImpl)),
		batch.NewBuildActionImpl,
		wire.Bind(new(batch.BuildAction), new(*batch.BuildActionImpl)),
		batch.NewDataHolderActionImpl,
		wire.Bind(new(batch.DataHolderAction), new(*batch.DataHolderActionImpl)),
		batch.NewDeploymentTemplateActionImpl,
		wire.Bind(new(batch.DeploymentTemplateAction), new(*batch.DeploymentTemplateActionImpl)),
		restHandler.NewBatchOperationRestHandlerImpl,
		wire.Bind(new(restHandler.BatchOperationRestHandler), new(*restHandler.BatchOperationRestHandlerImpl)),
		router.NewBatchOperationRouterImpl,
		wire.Bind(new(router.BatchOperationRouter), new(*router.BatchOperationRouterImpl)),

		repository4.NewChartGroupReposotoryImpl,
		wire.Bind(new(repository4.ChartGroupReposotory), new(*repository4.ChartGroupReposotoryImpl)),
		repository4.NewChartGroupEntriesRepositoryImpl,
		wire.Bind(new(repository4.ChartGroupEntriesRepository), new(*repository4.ChartGroupEntriesRepositoryImpl)),
		chartGroup.NewChartGroupServiceImpl,
		wire.Bind(new(chartGroup.ChartGroupService), new(*chartGroup.ChartGroupServiceImpl)),
		chartGroup2.NewChartGroupRestHandlerImpl,
		wire.Bind(new(chartGroup2.ChartGroupRestHandler), new(*chartGroup2.ChartGroupRestHandlerImpl)),
		chartGroup2.NewChartGroupRouterImpl,
		wire.Bind(new(chartGroup2.ChartGroupRouter), new(*chartGroup2.ChartGroupRouterImpl)),
		repository4.NewChartGroupDeploymentRepositoryImpl,
		wire.Bind(new(repository4.ChartGroupDeploymentRepository), new(*repository4.ChartGroupDeploymentRepositoryImpl)),

		commonService.NewCommonServiceImpl,
		wire.Bind(new(commonService.CommonService), new(*commonService.CommonServiceImpl)),

		router.NewImageScanRouterImpl,
		wire.Bind(new(router.ImageScanRouter), new(*router.ImageScanRouterImpl)),
		restHandler.NewImageScanRestHandlerImpl,
		wire.Bind(new(restHandler.ImageScanRestHandler), new(*restHandler.ImageScanRestHandlerImpl)),
		security.NewImageScanServiceImpl,
		wire.Bind(new(security.ImageScanService), new(*security.ImageScanServiceImpl)),
		security2.NewImageScanHistoryRepositoryImpl,
		wire.Bind(new(security2.ImageScanHistoryRepository), new(*security2.ImageScanHistoryRepositoryImpl)),
		security2.NewImageScanResultRepositoryImpl,
		wire.Bind(new(security2.ImageScanResultRepository), new(*security2.ImageScanResultRepositoryImpl)),
		security2.NewImageScanObjectMetaRepositoryImpl,
		wire.Bind(new(security2.ImageScanObjectMetaRepository), new(*security2.ImageScanObjectMetaRepositoryImpl)),
		security2.NewCveStoreRepositoryImpl,
		wire.Bind(new(security2.CveStoreRepository), new(*security2.CveStoreRepositoryImpl)),
		security2.NewImageScanDeployInfoRepositoryImpl,
		wire.Bind(new(security2.ImageScanDeployInfoRepository), new(*security2.ImageScanDeployInfoRepositoryImpl)),
		security2.NewScanToolMetadataRepositoryImpl,
		wire.Bind(new(security2.ScanToolMetadataRepository), new(*security2.ScanToolMetadataRepositoryImpl)),
		router.NewPolicyRouterImpl,
		wire.Bind(new(router.PolicyRouter), new(*router.PolicyRouterImpl)),
		restHandler.NewPolicyRestHandlerImpl,
		wire.Bind(new(restHandler.PolicyRestHandler), new(*restHandler.PolicyRestHandlerImpl)),
		security.NewPolicyServiceImpl,
		wire.Bind(new(security.PolicyService), new(*security.PolicyServiceImpl)),
		security2.NewPolicyRepositoryImpl,
		wire.Bind(new(security2.CvePolicyRepository), new(*security2.CvePolicyRepositoryImpl)),
		security2.NewScanToolExecutionHistoryMappingRepositoryImpl,
		wire.Bind(new(security2.ScanToolExecutionHistoryMappingRepository), new(*security2.ScanToolExecutionHistoryMappingRepositoryImpl)),

		argocdServer.NewArgoK8sClientImpl,
		wire.Bind(new(argocdServer.ArgoK8sClient), new(*argocdServer.ArgoK8sClientImpl)),

		grafana.GetConfig,
		router.NewGrafanaRouterImpl,
		wire.Bind(new(router.GrafanaRouter), new(*router.GrafanaRouterImpl)),

		router.NewGitOpsConfigRouterImpl,
		wire.Bind(new(router.GitOpsConfigRouter), new(*router.GitOpsConfigRouterImpl)),
		restHandler.NewGitOpsConfigRestHandlerImpl,
		wire.Bind(new(restHandler.GitOpsConfigRestHandler), new(*restHandler.GitOpsConfigRestHandlerImpl)),
		gitops.NewGitOpsConfigServiceImpl,
		wire.Bind(new(gitops.GitOpsConfigService), new(*gitops.GitOpsConfigServiceImpl)),

		router.NewAttributesRouterImpl,
		wire.Bind(new(router.AttributesRouter), new(*router.AttributesRouterImpl)),
		restHandler.NewAttributesRestHandlerImpl,
		wire.Bind(new(restHandler.AttributesRestHandler), new(*restHandler.AttributesRestHandlerImpl)),
		attributes.NewAttributesServiceImpl,
		wire.Bind(new(attributes.AttributesService), new(*attributes.AttributesServiceImpl)),
		repository.NewAttributesRepositoryImpl,
		wire.Bind(new(repository.AttributesRepository), new(*repository.AttributesRepositoryImpl)),

		router.NewCommonRouterImpl,
		wire.Bind(new(router.CommonRouter), new(*router.CommonRouterImpl)),
		restHandler.NewCommonRestHandlerImpl,
		wire.Bind(new(restHandler.CommonRestHandler), new(*restHandler.CommonRestHandlerImpl)),

		router.NewScopedVariableRouterImpl,
		wire.Bind(new(router.ScopedVariableRouter), new(*router.ScopedVariableRouterImpl)),
		scopedVariable.NewScopedVariableRestHandlerImpl,
		wire.Bind(new(scopedVariable.ScopedVariableRestHandler), new(*scopedVariable.ScopedVariableRestHandlerImpl)),

		router.NewTelemetryRouterImpl,
		wire.Bind(new(router.TelemetryRouter), new(*router.TelemetryRouterImpl)),
		restHandler.NewTelemetryRestHandlerImpl,
		wire.Bind(new(restHandler.TelemetryRestHandler), new(*restHandler.TelemetryRestHandlerImpl)),
		telemetry.NewPosthogClient,

		cloudProviderIdentifier.NewProviderIdentifierServiceImpl,
		wire.Bind(new(cloudProviderIdentifier.ProviderIdentifierService), new(*cloudProviderIdentifier.ProviderIdentifierServiceImpl)),

		telemetry.NewTelemetryEventClientImplExtended,
		wire.Bind(new(telemetry.TelemetryEventClient), new(*telemetry.TelemetryEventClientImplExtended)),

		router.NewBulkUpdateRouterImpl,
		wire.Bind(new(router.BulkUpdateRouter), new(*router.BulkUpdateRouterImpl)),
		restHandler.NewBulkUpdateRestHandlerImpl,
		wire.Bind(new(restHandler.BulkUpdateRestHandler), new(*restHandler.BulkUpdateRestHandlerImpl)),

		router.NewCoreAppRouterImpl,
		wire.Bind(new(router.CoreAppRouter), new(*router.CoreAppRouterImpl)),
		restHandler.NewCoreAppRestHandlerImpl,
		wire.Bind(new(restHandler.CoreAppRestHandler), new(*restHandler.CoreAppRestHandlerImpl)),

		app3.NewAppCrudOperationServiceEnterpriseImpl,
		wire.Bind(new(app.AppCrudOperationService), new(*app3.AppCrudOperationServiceEnterpriseImpl)),
		pipelineConfig.NewAppLabelRepositoryImpl,
		wire.Bind(new(pipelineConfig.AppLabelRepository), new(*pipelineConfig.AppLabelRepositoryImpl)),

		// Webhook
		repository.NewGitHostRepositoryImpl,
		wire.Bind(new(repository.GitHostRepository), new(*repository.GitHostRepositoryImpl)),
		restHandler.NewGitHostRestHandlerImpl,
		wire.Bind(new(restHandler.GitHostRestHandler), new(*restHandler.GitHostRestHandlerImpl)),
		restHandler.NewWebhookEventHandlerImpl,
		wire.Bind(new(restHandler.WebhookEventHandler), new(*restHandler.WebhookEventHandlerImpl)),
		router.NewGitHostRouterImpl,
		wire.Bind(new(router.GitHostRouter), new(*router.GitHostRouterImpl)),
		router.NewWebhookListenerRouterImpl,
		wire.Bind(new(router.WebhookListenerRouter), new(*router.WebhookListenerRouterImpl)),
		git.NewWebhookSecretValidatorImpl,
		wire.Bind(new(git.WebhookSecretValidator), new(*git.WebhookSecretValidatorImpl)),
		pipeline4.NewGitHostConfigImpl,
		wire.Bind(new(pipeline4.GitHostConfig), new(*pipeline4.GitHostConfigImpl)),
		repository.NewWebhookEventDataRepositoryImpl,
		wire.Bind(new(repository.WebhookEventDataRepository), new(*repository.WebhookEventDataRepositoryImpl)),
		pipeline4.NewWebhookEventDataConfigImpl,
		wire.Bind(new(pipeline4.WebhookEventDataConfig), new(*pipeline4.WebhookEventDataConfigImpl)),
		webhook.NewWebhookDataRestHandlerImpl,
		wire.Bind(new(webhook.WebhookDataRestHandler), new(*webhook.WebhookDataRestHandlerImpl)),

		app2.NewAppRouterImpl,
		wire.Bind(new(app2.AppRouter), new(*app2.AppRouterImpl)),
		appInfo2.NewAppInfoRouterImpl,
		wire.Bind(new(appInfo2.AppInfoRouter), new(*appInfo2.AppInfoRouterImpl)),
		appInfo.NewAppInfoRestHandlerImpl,
		wire.Bind(new(appInfo.AppInfoRestHandler), new(*appInfo.AppInfoRestHandlerImpl)),

		delete2.NewDeleteServiceExtendedImpl,
		wire.Bind(new(delete2.DeleteService), new(*delete2.DeleteServiceExtendedImpl)),
		delete2.NewDeleteServiceFullModeImpl,
		wire.Bind(new(delete2.DeleteServiceFullMode), new(*delete2.DeleteServiceFullModeImpl)),

		deployment3.NewFullModeDeploymentServiceImpl,
		wire.Bind(new(deployment3.FullModeDeploymentService), new(*deployment3.FullModeDeploymentServiceImpl)),
		//	util2.NewGoJsonSchemaCustomFormatChecker,

		// history starts
		history.NewPipelineHistoryRestHandlerImpl,
		wire.Bind(new(history.PipelineHistoryRestHandler), new(*history.PipelineHistoryRestHandlerImpl)),

		repository3.NewConfigMapHistoryRepositoryImpl,
		wire.Bind(new(repository3.ConfigMapHistoryRepository), new(*repository3.ConfigMapHistoryRepositoryImpl)),
		repository3.NewDeploymentTemplateHistoryRepositoryImpl,
		wire.Bind(new(repository3.DeploymentTemplateHistoryRepository), new(*repository3.DeploymentTemplateHistoryRepositoryImpl)),
		repository3.NewPrePostCiScriptHistoryRepositoryImpl,
		wire.Bind(new(repository3.PrePostCiScriptHistoryRepository), new(*repository3.PrePostCiScriptHistoryRepositoryImpl)),
		repository3.NewPrePostCdScriptHistoryRepositoryImpl,
		wire.Bind(new(repository3.PrePostCdScriptHistoryRepository), new(*repository3.PrePostCdScriptHistoryRepositoryImpl)),
		repository3.NewPipelineStrategyHistoryRepositoryImpl,
		wire.Bind(new(repository3.PipelineStrategyHistoryRepository), new(*repository3.PipelineStrategyHistoryRepositoryImpl)),
		repository3.NewGitMaterialHistoryRepositoyImpl,
		wire.Bind(new(repository3.GitMaterialHistoryRepository), new(*repository3.GitMaterialHistoryRepositoryImpl)),

		history3.NewCiTemplateHistoryServiceImpl,
		wire.Bind(new(history3.CiTemplateHistoryService), new(*history3.CiTemplateHistoryServiceImpl)),

		repository3.NewCiTemplateHistoryRepositoryImpl,
		wire.Bind(new(repository3.CiTemplateHistoryRepository), new(*repository3.CiTemplateHistoryRepositoryImpl)),

		history3.NewCiPipelineHistoryServiceImpl,
		wire.Bind(new(history3.CiPipelineHistoryService), new(*history3.CiPipelineHistoryServiceImpl)),

		repository3.NewCiPipelineHistoryRepositoryImpl,
		wire.Bind(new(repository3.CiPipelineHistoryRepository), new(*repository3.CiPipelineHistoryRepositoryImpl)),

		history3.NewPrePostCdScriptHistoryServiceImpl,
		wire.Bind(new(history3.PrePostCdScriptHistoryService), new(*history3.PrePostCdScriptHistoryServiceImpl)),
		history3.NewPrePostCiScriptHistoryServiceImpl,
		wire.Bind(new(history3.PrePostCiScriptHistoryService), new(*history3.PrePostCiScriptHistoryServiceImpl)),
		history3.NewDeploymentTemplateHistoryServiceImpl,
		wire.Bind(new(history3.DeploymentTemplateHistoryService), new(*history3.DeploymentTemplateHistoryServiceImpl)),
		history3.NewConfigMapHistoryServiceImpl,
		wire.Bind(new(history3.ConfigMapHistoryService), new(*history3.ConfigMapHistoryServiceImpl)),
		history3.NewPipelineStrategyHistoryServiceImpl,
		wire.Bind(new(history3.PipelineStrategyHistoryService), new(*history3.PipelineStrategyHistoryServiceImpl)),
		history3.NewGitMaterialHistoryServiceImpl,
		wire.Bind(new(history3.GitMaterialHistoryService), new(*history3.GitMaterialHistoryServiceImpl)),

		history3.NewDeployedConfigurationHistoryServiceImpl,
		wire.Bind(new(history3.DeployedConfigurationHistoryService), new(*history3.DeployedConfigurationHistoryServiceImpl)),
		// history ends

		// plugin starts
		repository6.NewGlobalPluginRepository,
		wire.Bind(new(repository6.GlobalPluginRepository), new(*repository6.GlobalPluginRepositoryImpl)),

		plugin.NewGlobalPluginService,
		wire.Bind(new(plugin.GlobalPluginService), new(*plugin.GlobalPluginServiceImpl)),

		restHandler.NewGlobalPluginRestHandler,
		wire.Bind(new(restHandler.GlobalPluginRestHandler), new(*restHandler.GlobalPluginRestHandlerImpl)),

		router.NewGlobalPluginRouter,
		wire.Bind(new(router.GlobalPluginRouter), new(*router.GlobalPluginRouterImpl)),

		repository5.NewPipelineStageRepository,
		wire.Bind(new(repository5.PipelineStageRepository), new(*repository5.PipelineStageRepositoryImpl)),

		pipeline4.NewPipelineStageService,
		wire.Bind(new(pipeline4.PipelineStageService), new(*pipeline4.PipelineStageServiceImpl)),
		// plugin ends

		connection.NewArgoCDConnectionManagerImpl,
		wire.Bind(new(connection.ArgoCDConnectionManager), new(*connection.ArgoCDConnectionManagerImpl)),
		argo.NewArgoUserServiceImpl,
		wire.Bind(new(argo.ArgoUserService), new(*argo.ArgoUserServiceImpl)),
		util2.GetDevtronSecretName,
		//	AuthWireSet,

		cron.NewCdApplicationStatusUpdateHandlerImpl,
		wire.Bind(new(cron.CdApplicationStatusUpdateHandler), new(*cron.CdApplicationStatusUpdateHandlerImpl)),

		// app_status
		appStatusRepo.NewAppStatusRepositoryImpl,
		wire.Bind(new(appStatusRepo.AppStatusRepository), new(*appStatusRepo.AppStatusRepositoryImpl)),
		appStatus.NewAppStatusServiceImpl,
		wire.Bind(new(appStatus.AppStatusService), new(*appStatus.AppStatusServiceImpl)),
		// app_status ends

		cron.GetCiWorkflowStatusUpdateConfig,
		cron.NewCiStatusUpdateCronImpl,
		wire.Bind(new(cron.CiStatusUpdateCron), new(*cron.CiStatusUpdateCronImpl)),

		cron.GetCiTriggerCronConfig,
		cron.NewCiTriggerCronImpl,
		wire.Bind(new(cron.CiTriggerCron), new(*cron.CiTriggerCronImpl)),

		status2.NewPipelineStatusTimelineRestHandlerImpl,
		wire.Bind(new(status2.PipelineStatusTimelineRestHandler), new(*status2.PipelineStatusTimelineRestHandlerImpl)),

		status.NewPipelineStatusTimelineServiceImpl,
		wire.Bind(new(status.PipelineStatusTimelineService), new(*status.PipelineStatusTimelineServiceImpl)),

		router.NewUserAttributesRouterImpl,
		wire.Bind(new(router.UserAttributesRouter), new(*router.UserAttributesRouterImpl)),
		restHandler.NewUserAttributesRestHandlerImpl,
		wire.Bind(new(restHandler.UserAttributesRestHandler), new(*restHandler.UserAttributesRestHandlerImpl)),
		attributes.NewUserAttributesServiceImpl,
		wire.Bind(new(attributes.UserAttributesService), new(*attributes.UserAttributesServiceImpl)),
		repository.NewUserAttributesRepositoryImpl,
		wire.Bind(new(repository.UserAttributesRepository), new(*repository.UserAttributesRepositoryImpl)),
		pipelineConfig.NewPipelineStatusTimelineRepositoryImpl,
		wire.Bind(new(pipelineConfig.PipelineStatusTimelineRepository), new(*pipelineConfig.PipelineStatusTimelineRepositoryImpl)),
		wire.Bind(new(pipeline4.DeploymentConfigService), new(*pipeline4.DeploymentConfigServiceImpl)),
		pipeline4.NewDeploymentConfigServiceImpl,
		pipelineConfig.NewCiTemplateOverrideRepositoryImpl,
		wire.Bind(new(pipelineConfig.CiTemplateOverrideRepository), new(*pipelineConfig.CiTemplateOverrideRepositoryImpl)),
		pipelineConfig.NewCiBuildConfigRepositoryImpl,
		wire.Bind(new(pipelineConfig.CiBuildConfigRepository), new(*pipelineConfig.CiBuildConfigRepositoryImpl)),
		pipeline4.NewCiBuildConfigServiceImpl,
		wire.Bind(new(pipeline4.CiBuildConfigService), new(*pipeline4.CiBuildConfigServiceImpl)),
		resourceFilter.NewCELServiceImpl,
		wire.Bind(new(resourceFilter.CELEvaluatorService), new(*resourceFilter.CELServiceImpl)),
		resourceFilter.NewResourceFilterRepositoryImpl,
		wire.Bind(new(resourceFilter.ResourceFilterRepository), new(*resourceFilter.ResourceFilterRepositoryImpl)),
		resourceFilter.NewFilterAuditRepositoryImpl,
		wire.Bind(new(resourceFilter.FilterAuditRepository), new(*resourceFilter.FilterAuditRepositoryImpl)),
		resourceFilter.NewFilterEvaluationAuditRepositoryImpl,
		wire.Bind(new(resourceFilter.FilterEvaluationAuditRepository), new(*resourceFilter.FilterEvaluationAuditRepositoryImpl)),
		resourceFilter.NewFilterEvaluationAuditServiceImpl,
		wire.Bind(new(resourceFilter.FilterEvaluationAuditService), new(*resourceFilter.FilterEvaluationAuditServiceImpl)),
		resourceFilter.NewResourceFilterServiceImpl,
		wire.Bind(new(resourceFilter.ResourceFilterService), new(*resourceFilter.ResourceFilterServiceImpl)),
		resourceFilter.NewResourceFilterEvaluatorImpl,
		wire.Bind(new(resourceFilter.ResourceFilterEvaluator), new(*resourceFilter.ResourceFilterEvaluatorImpl)),
		resourceFilter2.NewResourceFilterRestHandlerImpl,
		wire.Bind(new(resourceFilter2.ResourceFilterRestHandler), new(*resourceFilter2.ResourceFilterRestHandlerImpl)),
		router.NewResourceFilterRouterImpl,
		wire.Bind(new(router.ResourceFilterRouter), new(*router.ResourceFilterRouterImpl)),
		pipeline4.NewCiTemplateServiceImpl,
		wire.Bind(new(pipeline4.CiTemplateService), new(*pipeline4.CiTemplateServiceImpl)),
		router.NewGlobalCMCSRouterImpl,
		wire.Bind(new(router.GlobalCMCSRouter), new(*router.GlobalCMCSRouterImpl)),
		restHandler.NewGlobalCMCSRestHandlerImpl,
		wire.Bind(new(restHandler.GlobalCMCSRestHandler), new(*restHandler.GlobalCMCSRestHandlerImpl)),
		pipeline4.NewGlobalCMCSServiceImpl,
		wire.Bind(new(pipeline4.GlobalCMCSService), new(*pipeline4.GlobalCMCSServiceImpl)),
		repository.NewGlobalCMCSRepositoryImpl,
		wire.Bind(new(repository.GlobalCMCSRepository), new(*repository.GlobalCMCSRepositoryImpl)),

		// chartRepoRepository.NewGlobalStrategyMetadataRepositoryImpl,
		// wire.Bind(new(chartRepoRepository.GlobalStrategyMetadataRepository), new(*chartRepoRepository.GlobalStrategyMetadataRepositoryImpl)),
		chartRepoRepository.NewGlobalStrategyMetadataChartRefMappingRepositoryImpl,
		wire.Bind(new(chartRepoRepository.GlobalStrategyMetadataChartRefMappingRepository), new(*chartRepoRepository.GlobalStrategyMetadataChartRefMappingRepositoryImpl)),

		status.NewPipelineStatusTimelineResourcesServiceImpl,
		wire.Bind(new(status.PipelineStatusTimelineResourcesService), new(*status.PipelineStatusTimelineResourcesServiceImpl)),
		pipelineConfig.NewPipelineStatusTimelineResourcesRepositoryImpl,
		wire.Bind(new(pipelineConfig.PipelineStatusTimelineResourcesRepository), new(*pipelineConfig.PipelineStatusTimelineResourcesRepositoryImpl)),

		status.NewPipelineStatusSyncDetailServiceImpl,
		wire.Bind(new(status.PipelineStatusSyncDetailService), new(*status.PipelineStatusSyncDetailServiceImpl)),
		pipelineConfig.NewPipelineStatusSyncDetailRepositoryImpl,
		wire.Bind(new(pipelineConfig.PipelineStatusSyncDetailRepository), new(*pipelineConfig.PipelineStatusSyncDetailRepositoryImpl)),

		repository7.NewK8sResourceHistoryRepositoryImpl,
		wire.Bind(new(repository7.K8sResourceHistoryRepository), new(*repository7.K8sResourceHistoryRepositoryImpl)),

		kubernetesResourceAuditLogs.Newk8sResourceHistoryServiceImpl,
		wire.Bind(new(kubernetesResourceAuditLogs.K8sResourceHistoryService), new(*kubernetesResourceAuditLogs.K8sResourceHistoryServiceImpl)),
		pipelineConfig.NewRequestApprovalUserDataRepositoryImpl,
		wire.Bind(new(pipelineConfig.RequestApprovalUserdataRepository), new(*pipelineConfig.RequestApprovalUserDataRepositoryImpl)),
		pipelineConfig.NewDeploymentApprovalRepositoryImpl,
		wire.Bind(new(pipelineConfig.DeploymentApprovalRepository), new(*pipelineConfig.DeploymentApprovalRepositoryImpl)),
		router.NewResourceGroupingRouterImpl,
		wire.Bind(new(router.ResourceGroupingRouter), new(*router.ResourceGroupingRouterImpl)),
		restHandler.NewResourceGroupRestHandlerImpl,
		wire.Bind(new(restHandler.ResourceGroupRestHandler), new(*restHandler.ResourceGroupRestHandlerImpl)),
		resourceGroup2.NewResourceGroupServiceImpl,
		wire.Bind(new(resourceGroup2.ResourceGroupService), new(*resourceGroup2.ResourceGroupServiceImpl)),
		resourceGroup.NewResourceGroupRepositoryImpl,
		wire.Bind(new(resourceGroup.ResourceGroupRepository), new(*resourceGroup.ResourceGroupRepositoryImpl)),
		resourceGroup.NewResourceGroupMappingRepositoryImpl,
		wire.Bind(new(resourceGroup.ResourceGroupMappingRepository), new(*resourceGroup.ResourceGroupMappingRepositoryImpl)),
		executors.NewArgoWorkflowExecutorImpl,
		wire.Bind(new(executors.ArgoWorkflowExecutor), new(*executors.ArgoWorkflowExecutorImpl)),
		executors.NewSystemWorkflowExecutorImpl,
		wire.Bind(new(executors.SystemWorkflowExecutor), new(*executors.SystemWorkflowExecutorImpl)),
		repository5.NewManifestPushConfigRepository,
		wire.Bind(new(repository5.ManifestPushConfigRepository), new(*repository5.ManifestPushConfigRepositoryImpl)),
		app.NewGitOpsManifestPushServiceImpl,
		wire.Bind(new(app.GitOpsPushService), new(*app.GitOpsManifestPushServiceImpl)),

		app.NewHelmRepoPushServiceImpl,
		wire.Bind(new(app.HelmRepoPushService), new(*app.HelmRepoPushServiceImpl)),

		// start: docker registry wire set injection
		router.NewDockerRegRouterImpl,
		wire.Bind(new(router.DockerRegRouter), new(*router.DockerRegRouterImpl)),
		restHandler.NewDockerRegRestHandlerExtendedImpl,
		wire.Bind(new(restHandler.DockerRegRestHandler), new(*restHandler.DockerRegRestHandlerExtendedImpl)),
		pipeline4.NewDockerRegistryConfigImpl,
		wire.Bind(new(pipeline4.DockerRegistryConfig), new(*pipeline4.DockerRegistryConfigImpl)),
		dockerRegistry.NewDockerRegistryIpsConfigServiceImpl,
		wire.Bind(new(dockerRegistry.DockerRegistryIpsConfigService), new(*dockerRegistry.DockerRegistryIpsConfigServiceImpl)),
		dockerRegistryRepository.NewDockerArtifactStoreRepositoryImpl,
		wire.Bind(new(dockerRegistryRepository.DockerArtifactStoreRepository), new(*dockerRegistryRepository.DockerArtifactStoreRepositoryImpl)),
		dockerRegistryRepository.NewDockerRegistryIpsConfigRepositoryImpl,
		wire.Bind(new(dockerRegistryRepository.DockerRegistryIpsConfigRepository), new(*dockerRegistryRepository.DockerRegistryIpsConfigRepositoryImpl)),
		dockerRegistryRepository.NewOCIRegistryConfigRepositoryImpl,
		wire.Bind(new(dockerRegistryRepository.OCIRegistryConfigRepository), new(*dockerRegistryRepository.OCIRegistryConfigRepositoryImpl)),
		// end: docker registry wire set injection
		util4.NewSSHTunnelWrapperServiceImpl,
		wire.Bind(new(util4.SSHTunnelWrapperService), new(*util4.SSHTunnelWrapperServiceImpl)),

		resourceQualifiers.NewQualifiersMappingRepositoryImpl,
		wire.Bind(new(resourceQualifiers.QualifiersMappingRepository), new(*resourceQualifiers.QualifiersMappingRepositoryImpl)),

		resourceQualifiers.NewQualifierMappingServiceImpl,
		wire.Bind(new(resourceQualifiers.QualifierMappingService), new(*resourceQualifiers.QualifierMappingServiceImpl)),

		argocdServer.NewArgoClientWrapperServiceImpl,
		wire.Bind(new(argocdServer.ArgoClientWrapperService), new(*argocdServer.ArgoClientWrapperServiceImpl)),

		pipeline4.NewPluginInputVariableParserImpl,
		wire.Bind(new(pipeline4.PluginInputVariableParser), new(*pipeline4.PluginInputVariableParserImpl)),

		imageDigestPolicy.NewImageDigestPolicyServiceImpl,
		wire.Bind(new(imageDigestPolicy.ImageDigestPolicyService), new(*imageDigestPolicy.ImageDigestPolicyServiceImpl)),

		router.NewImageDigestPolicyRouterImpl,
		wire.Bind(new(router.ImageDigestPolicyRouter), new(*router.ImageDigestPolicyRouterImpl)),

		imageDigestPolicy2.NewImageDigestPolicyRestHandlerImpl,
		wire.Bind(new(imageDigestPolicy2.ImageDigestPolicyRestHandler), new(*imageDigestPolicy2.ImageDigestPolicyRestHandlerImpl)),

		cron2.NewCronLoggerImpl,

		timeoutWindow.NewTimeWindowServiceImpl,
		wire.Bind(new(timeoutWindow.TimeoutWindowService), new(*timeoutWindow.TimeWindowServiceImpl)),

		repository9.NewTimeWindowRepositoryImpl,
		wire.Bind(new(repository9.TimeWindowRepository), new(*repository9.TimeWindowRepositoryImpl)),

		artifactPromotionApprovalRequest.NewRouterImpl,
		wire.Bind(new(artifactPromotionApprovalRequest.Router), new(*artifactPromotionApprovalRequest.RouterImpl)),

		artifactPromotionApprovalRequest.NewRestHandlerImpl,
		wire.Bind(new(artifactPromotionApprovalRequest.RestHandler), new(*artifactPromotionApprovalRequest.RestHandlerImpl)),
		wire.Bind(new(artifactPromotionApprovalRequest.MaterialRestHandler), new(*artifactPromotionApprovalRequest.RestHandlerImpl)),

		artifactPromotionPolicy.NewCommonPolicyRouterImpl,
		wire.Bind(new(artifactPromotionPolicy.Router), new(*artifactPromotionPolicy.RouterImpl)),

		artifactPromotionPolicy.NewArtifactPromotionPolicyRestHandlerImpl,
		wire.Bind(new(artifactPromotionPolicy.RestHandler), new(*artifactPromotionPolicy.RestHandlerImpl)),

		globalPolicy2.NewGlobalPolicyDataManagerImpl,
		wire.Bind(new(globalPolicy2.GlobalPolicyDataManager), new(*globalPolicy2.GlobalPolicyDataManagerImpl)),
	)
	return &App{}, nil
}
