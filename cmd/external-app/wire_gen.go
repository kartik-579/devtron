// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/devtron-labs/authenticator/apiToken"
	"github.com/devtron-labs/authenticator/client"
	"github.com/devtron-labs/authenticator/middleware"
	"github.com/devtron-labs/common-lib-private/utils/k8s"
	"github.com/devtron-labs/common-lib/cloud-provider-identifier"
	apiToken2 "github.com/devtron-labs/devtron/api/apiToken"
	chartProvider2 "github.com/devtron-labs/devtron/api/appStore/chartProvider"
	"github.com/devtron-labs/devtron/api/appStore/deployment"
	"github.com/devtron-labs/devtron/api/appStore/discover"
	"github.com/devtron-labs/devtron/api/appStore/values"
	"github.com/devtron-labs/devtron/api/auth/authorisation/globalConfig"
	sso2 "github.com/devtron-labs/devtron/api/auth/sso"
	user2 "github.com/devtron-labs/devtron/api/auth/user"
	chartRepo2 "github.com/devtron-labs/devtron/api/chartRepo"
	cluster2 "github.com/devtron-labs/devtron/api/cluster"
	"github.com/devtron-labs/devtron/api/connector"
	"github.com/devtron-labs/devtron/api/dashboardEvent"
	externalLink2 "github.com/devtron-labs/devtron/api/externalLink"
	client3 "github.com/devtron-labs/devtron/api/helm-app"
	"github.com/devtron-labs/devtron/api/helm-app/gRPC"
	"github.com/devtron-labs/devtron/api/helm-app/service"
	application2 "github.com/devtron-labs/devtron/api/k8s/application"
	capacity2 "github.com/devtron-labs/devtron/api/k8s/capacity"
	module2 "github.com/devtron-labs/devtron/api/module"
	"github.com/devtron-labs/devtron/api/restHandler"
	"github.com/devtron-labs/devtron/api/restHandler/app/appInfo"
	"github.com/devtron-labs/devtron/api/restHandler/app/appList"
	"github.com/devtron-labs/devtron/api/router"
	app3 "github.com/devtron-labs/devtron/api/router/app"
	appInfo2 "github.com/devtron-labs/devtron/api/router/app/appInfo"
	appList2 "github.com/devtron-labs/devtron/api/router/app/appList"
	server2 "github.com/devtron-labs/devtron/api/server"
	team2 "github.com/devtron-labs/devtron/api/team"
	terminal2 "github.com/devtron-labs/devtron/api/terminal"
	webhookHelm2 "github.com/devtron-labs/devtron/api/webhook/helm"
	"github.com/devtron-labs/devtron/client/argocdServer"
	"github.com/devtron-labs/devtron/client/dashboard"
	"github.com/devtron-labs/devtron/client/telemetry"
	"github.com/devtron-labs/devtron/enterprise/pkg/deploymentWindow"
	repository5 "github.com/devtron-labs/devtron/internal/sql/repository"
	"github.com/devtron-labs/devtron/internal/sql/repository/app"
	"github.com/devtron-labs/devtron/internal/sql/repository/appStatus"
	repository7 "github.com/devtron-labs/devtron/internal/sql/repository/dockerRegistry"
	"github.com/devtron-labs/devtron/internal/sql/repository/helper"
	"github.com/devtron-labs/devtron/internal/sql/repository/pipelineConfig"
	"github.com/devtron-labs/devtron/internal/sql/repository/security"
	"github.com/devtron-labs/devtron/internal/util"
	"github.com/devtron-labs/devtron/pkg/apiToken"
	app2 "github.com/devtron-labs/devtron/pkg/app"
	repository9 "github.com/devtron-labs/devtron/pkg/appStore/chartGroup/repository"
	"github.com/devtron-labs/devtron/pkg/appStore/chartProvider"
	"github.com/devtron-labs/devtron/pkg/appStore/discover/repository"
	service3 "github.com/devtron-labs/devtron/pkg/appStore/discover/service"
	repository6 "github.com/devtron-labs/devtron/pkg/appStore/installedApp/repository"
	service2 "github.com/devtron-labs/devtron/pkg/appStore/installedApp/service"
	"github.com/devtron-labs/devtron/pkg/appStore/installedApp/service/EAMode"
	"github.com/devtron-labs/devtron/pkg/appStore/installedApp/service/common"
	"github.com/devtron-labs/devtron/pkg/appStore/values/repository"
	service4 "github.com/devtron-labs/devtron/pkg/appStore/values/service"
	"github.com/devtron-labs/devtron/pkg/argoApplication"
	"github.com/devtron-labs/devtron/pkg/attributes"
	"github.com/devtron-labs/devtron/pkg/auth/authentication"
	"github.com/devtron-labs/devtron/pkg/auth/authorisation/casbin"
	client2 "github.com/devtron-labs/devtron/pkg/auth/authorisation/casbin/client"
	"github.com/devtron-labs/devtron/pkg/auth/authorisation/globalConfig"
	"github.com/devtron-labs/devtron/pkg/auth/authorisation/globalConfig/repository"
	"github.com/devtron-labs/devtron/pkg/auth/sso"
	"github.com/devtron-labs/devtron/pkg/auth/user"
	repository2 "github.com/devtron-labs/devtron/pkg/auth/user/repository"
	"github.com/devtron-labs/devtron/pkg/chartRepo"
	"github.com/devtron-labs/devtron/pkg/chartRepo/repository"
	"github.com/devtron-labs/devtron/pkg/cluster"
	repository4 "github.com/devtron-labs/devtron/pkg/cluster/repository"
	"github.com/devtron-labs/devtron/pkg/clusterTerminalAccess"
	delete2 "github.com/devtron-labs/devtron/pkg/delete"
	"github.com/devtron-labs/devtron/pkg/deployment/gitOps/config"
	"github.com/devtron-labs/devtron/pkg/devtronResource"
	repository10 "github.com/devtron-labs/devtron/pkg/devtronResource/repository"
	"github.com/devtron-labs/devtron/pkg/externalLink"
	"github.com/devtron-labs/devtron/pkg/genericNotes"
	repository8 "github.com/devtron-labs/devtron/pkg/genericNotes/repository"
	"github.com/devtron-labs/devtron/pkg/globalPolicy"
	repository12 "github.com/devtron-labs/devtron/pkg/globalPolicy/repository"
	k8s2 "github.com/devtron-labs/devtron/pkg/k8s"
	"github.com/devtron-labs/devtron/pkg/k8s/application"
	"github.com/devtron-labs/devtron/pkg/k8s/capacity"
	"github.com/devtron-labs/devtron/pkg/k8s/informer"
	"github.com/devtron-labs/devtron/pkg/kubernetesResourceAuditLogs"
	repository11 "github.com/devtron-labs/devtron/pkg/kubernetesResourceAuditLogs/repository"
	"github.com/devtron-labs/devtron/pkg/module"
	"github.com/devtron-labs/devtron/pkg/module/repo"
	"github.com/devtron-labs/devtron/pkg/module/store"
	"github.com/devtron-labs/devtron/pkg/pipeline"
	"github.com/devtron-labs/devtron/pkg/resourceQualifiers"
	"github.com/devtron-labs/devtron/pkg/server"
	"github.com/devtron-labs/devtron/pkg/server/config"
	"github.com/devtron-labs/devtron/pkg/server/store"
	"github.com/devtron-labs/devtron/pkg/sql"
	"github.com/devtron-labs/devtron/pkg/team"
	"github.com/devtron-labs/devtron/pkg/terminal"
	"github.com/devtron-labs/devtron/pkg/timeoutWindow"
	repository3 "github.com/devtron-labs/devtron/pkg/timeoutWindow/repository"
	util3 "github.com/devtron-labs/devtron/pkg/util"
	"github.com/devtron-labs/devtron/pkg/webhook/helm"
	util2 "github.com/devtron-labs/devtron/util"
	"github.com/devtron-labs/devtron/util/argo"
	"github.com/devtron-labs/devtron/util/cron"
	"github.com/devtron-labs/devtron/util/rbac"
)

// Injectors from wire.go:

func InitializeApp() (*App, error) {
	sqlConfig, err := sql.GetConfig()
	if err != nil {
		return nil, err
	}
	sugaredLogger, err := util.NewSugardLogger()
	if err != nil {
		return nil, err
	}
	db, err := sql.NewDbConnection(sqlConfig, sugaredLogger)
	if err != nil {
		return nil, err
	}
	runtimeConfig, err := client.GetRuntimeConfig()
	if err != nil {
		return nil, err
	}
	k8sClient, err := client.NewK8sClient(runtimeConfig)
	if err != nil {
		return nil, err
	}
	dexConfig, err := client.BuildDexConfig(k8sClient)
	if err != nil {
		return nil, err
	}
	settings, err := client.GetSettings(dexConfig)
	if err != nil {
		return nil, err
	}
	apiTokenSecretStore := apiTokenAuth.InitApiTokenSecretStore()
	sessionManager := middleware.NewSessionManager(settings, dexConfig, apiTokenSecretStore)
	validate, err := util.IntValidator()
	if err != nil {
		return nil, err
	}
	syncedEnforcer := casbin.Create()
	casbinSyncedEnforcer := casbin.CreateV2()
	casbinClientConfig, err := client2.GetConfig()
	if err != nil {
		return nil, err
	}
	casbinClientImpl := client2.NewCasbinClientImpl(sugaredLogger, casbinClientConfig)
	casbinServiceImpl := casbin.NewCasbinServiceImpl(sugaredLogger, casbinClientImpl)
	globalAuthorisationConfigRepositoryImpl := repository.NewGlobalAuthorisationConfigRepositoryImpl(sugaredLogger, db)
	globalAuthorisationConfigServiceImpl := auth.NewGlobalAuthorisationConfigServiceImpl(sugaredLogger, globalAuthorisationConfigRepositoryImpl)
	enterpriseEnforcerImpl, err := casbin.NewEnterpriseEnforcerImpl(syncedEnforcer, casbinSyncedEnforcer, sessionManager, sugaredLogger, casbinServiceImpl, globalAuthorisationConfigServiceImpl)
	if err != nil {
		return nil, err
	}
	defaultAuthPolicyRepositoryImpl := repository2.NewDefaultAuthPolicyRepositoryImpl(db, sugaredLogger)
	defaultAuthRoleRepositoryImpl := repository2.NewDefaultAuthRoleRepositoryImpl(db, sugaredLogger)
	userAuthRepositoryImpl := repository2.NewUserAuthRepositoryImpl(db, sugaredLogger, defaultAuthPolicyRepositoryImpl, defaultAuthRoleRepositoryImpl)
	userRepositoryImpl := repository2.NewUserRepositoryImpl(db, sugaredLogger)
	roleGroupRepositoryImpl := repository2.NewRoleGroupRepositoryImpl(db, sugaredLogger)
	rbacPolicyDataRepositoryImpl := repository2.NewRbacPolicyDataRepositoryImpl(sugaredLogger, db)
	rbacRoleDataRepositoryImpl := repository2.NewRbacRoleDataRepositoryImpl(sugaredLogger, db)
	rbacDataCacheFactoryImpl := repository2.NewRbacDataCacheFactoryImpl(sugaredLogger, rbacPolicyDataRepositoryImpl, rbacRoleDataRepositoryImpl)
	userCommonServiceImpl := user.NewUserCommonServiceImpl(userAuthRepositoryImpl, sugaredLogger, userRepositoryImpl, roleGroupRepositoryImpl, sessionManager, rbacDataCacheFactoryImpl)
	userAuditRepositoryImpl := repository2.NewUserAuditRepositoryImpl(db)
	userAuditServiceImpl := user.NewUserAuditServiceImpl(sugaredLogger, userAuditRepositoryImpl)
	roleGroupServiceImpl := user.NewRoleGroupServiceImpl(userAuthRepositoryImpl, sugaredLogger, userRepositoryImpl, roleGroupRepositoryImpl, userCommonServiceImpl)
	userGroupMapRepositoryImpl := repository2.NewUserGroupMapRepositoryImpl(db, sugaredLogger)
	timeWindowRepositoryImpl := repository3.NewTimeWindowRepositoryImpl(db, sugaredLogger)
	timeWindowServiceImpl := timeoutWindow.NewTimeWindowServiceImpl(sugaredLogger, timeWindowRepositoryImpl)
	userServiceImpl := user.NewUserServiceImpl(userAuthRepositoryImpl, sugaredLogger, userRepositoryImpl, roleGroupRepositoryImpl, sessionManager, userCommonServiceImpl, userAuditServiceImpl, globalAuthorisationConfigServiceImpl, roleGroupServiceImpl, userGroupMapRepositoryImpl, enterpriseEnforcerImpl, timeWindowServiceImpl)
	ssoLoginRepositoryImpl := sso.NewSSOLoginRepositoryImpl(db)
	sshTunnelWrapperServiceImpl, err := k8s.NewSSHTunnelWrapperServiceImpl(sugaredLogger)
	if err != nil {
		return nil, err
	}
	k8sUtilExtended := k8s.NewK8sUtilExtended(sugaredLogger, runtimeConfig, sshTunnelWrapperServiceImpl)
	devtronSecretConfig, err := util2.GetDevtronSecretName()
	if err != nil {
		return nil, err
	}
	selfRegistrationRolesRepositoryImpl := repository2.NewSelfRegistrationRolesRepositoryImpl(db, sugaredLogger)
	userSelfRegistrationServiceImpl := user.NewUserSelfRegistrationServiceImpl(sugaredLogger, selfRegistrationRolesRepositoryImpl, userServiceImpl, globalAuthorisationConfigServiceImpl)
	userAuthOidcHelperImpl, err := authentication.NewUserAuthOidcHelperImpl(sugaredLogger, userSelfRegistrationServiceImpl, dexConfig, settings, sessionManager)
	if err != nil {
		return nil, err
	}
	ssoLoginServiceImpl := sso.NewSSOLoginServiceImpl(sugaredLogger, ssoLoginRepositoryImpl, k8sUtilExtended, devtronSecretConfig, userAuthOidcHelperImpl, globalAuthorisationConfigServiceImpl)
	ssoLoginRestHandlerImpl := sso2.NewSsoLoginRestHandlerImpl(validate, sugaredLogger, enterpriseEnforcerImpl, userServiceImpl, ssoLoginServiceImpl)
	ssoLoginRouterImpl := sso2.NewSsoLoginRouterImpl(ssoLoginRestHandlerImpl)
	teamRepositoryImpl := team.NewTeamRepositoryImpl(db)
	loginService := middleware.NewUserLogin(sessionManager, k8sClient)
	userAuthServiceImpl := user.NewUserAuthServiceImpl(userAuthRepositoryImpl, sessionManager, loginService, sugaredLogger, userRepositoryImpl, roleGroupRepositoryImpl, userServiceImpl)
	teamServiceImpl := team.NewTeamServiceImpl(sugaredLogger, teamRepositoryImpl, userAuthServiceImpl)
	clusterRepositoryImpl := repository4.NewClusterRepositoryImpl(db, sugaredLogger)
	v := informer.NewGlobalMapClusterNamespace()
	k8sInformerFactoryImpl := informer.NewK8sInformerFactoryImpl(sugaredLogger, v, runtimeConfig, k8sUtilExtended)
	clusterServiceImpl := cluster.NewClusterServiceImpl(clusterRepositoryImpl, sugaredLogger, k8sUtilExtended, k8sInformerFactoryImpl, userAuthRepositoryImpl, userRepositoryImpl, roleGroupRepositoryImpl, globalAuthorisationConfigServiceImpl, userServiceImpl)
	appStatusRepositoryImpl := appStatus.NewAppStatusRepositoryImpl(db, sugaredLogger)
	environmentRepositoryImpl := repository4.NewEnvironmentRepositoryImpl(db, sugaredLogger, appStatusRepositoryImpl)
	attributesRepositoryImpl := repository5.NewAttributesRepositoryImpl(db)
	environmentServiceImpl := cluster.NewEnvironmentServiceImpl(environmentRepositoryImpl, clusterServiceImpl, sugaredLogger, k8sUtilExtended, k8sInformerFactoryImpl, userAuthServiceImpl, attributesRepositoryImpl)
	chartRepoRepositoryImpl := chartRepoRepository.NewChartRepoRepositoryImpl(db)
	acdAuthConfig, err := util3.GetACDAuthConfig()
	if err != nil {
		return nil, err
	}
	httpClient := util.NewHttpClient()
	serverEnvConfigServerEnvConfig, err := serverEnvConfig.ParseServerEnvConfig()
	if err != nil {
		return nil, err
	}
	chartRepositoryServiceImpl := chartRepo.NewChartRepositoryServiceImpl(sugaredLogger, chartRepoRepositoryImpl, k8sUtilExtended, clusterServiceImpl, acdAuthConfig, httpClient, serverEnvConfigServerEnvConfig)
	installedAppRepositoryImpl := repository6.NewInstalledAppRepositoryImpl(sugaredLogger, db)
	helmClientConfig, err := gRPC.GetConfig()
	if err != nil {
		return nil, err
	}
	helmAppClientImpl := gRPC.NewHelmAppClientImpl(sugaredLogger, helmClientConfig)
	pumpImpl := connector.NewPumpImpl(sugaredLogger)
	appRepositoryImpl := app.NewAppRepositoryImpl(db, sugaredLogger)
	enforcerUtilHelmImpl := rbac.NewEnforcerUtilHelmImpl(sugaredLogger, clusterRepositoryImpl, teamRepositoryImpl, appRepositoryImpl, environmentRepositoryImpl, installedAppRepositoryImpl)
	serverDataStoreServerDataStore := serverDataStore.InitServerDataStore()
	appStoreApplicationVersionRepositoryImpl := appStoreDiscoverRepository.NewAppStoreApplicationVersionRepositoryImpl(sugaredLogger, db)
	pipelineRepositoryImpl := pipelineConfig.NewPipelineRepositoryImpl(db, sugaredLogger)
	helmReleaseConfig, err := service.GetHelmReleaseConfig()
	if err != nil {
		return nil, err
	}
	helmAppServiceImpl := service.NewHelmAppServiceImpl(sugaredLogger, clusterServiceImpl, helmAppClientImpl, pumpImpl, enforcerUtilHelmImpl, serverDataStoreServerDataStore, serverEnvConfigServerEnvConfig, appStoreApplicationVersionRepositoryImpl, environmentServiceImpl, pipelineRepositoryImpl, installedAppRepositoryImpl, appRepositoryImpl, clusterRepositoryImpl, k8sUtilExtended, helmReleaseConfig)
	dockerArtifactStoreRepositoryImpl := repository7.NewDockerArtifactStoreRepositoryImpl(db)
	dockerRegistryIpsConfigRepositoryImpl := repository7.NewDockerRegistryIpsConfigRepositoryImpl(db)
	ociRegistryConfigRepositoryImpl := repository7.NewOCIRegistryConfigRepositoryImpl(db)
	dockerRegistryConfigImpl := pipeline.NewDockerRegistryConfigImpl(sugaredLogger, helmAppServiceImpl, dockerArtifactStoreRepositoryImpl, dockerRegistryIpsConfigRepositoryImpl, ociRegistryConfigRepositoryImpl)
	deleteServiceImpl := delete2.NewDeleteServiceImpl(sugaredLogger, teamServiceImpl, clusterServiceImpl, environmentServiceImpl, chartRepositoryServiceImpl, installedAppRepositoryImpl, dockerRegistryConfigImpl, dockerArtifactStoreRepositoryImpl)
	teamRestHandlerImpl := team2.NewTeamRestHandlerImpl(sugaredLogger, teamServiceImpl, userServiceImpl, enterpriseEnforcerImpl, validate, userAuthServiceImpl, deleteServiceImpl)
	teamRouterImpl := team2.NewTeamRouterImpl(teamRestHandlerImpl)
	userAuthHandlerImpl := user2.NewUserAuthHandlerImpl(userAuthServiceImpl, validate, sugaredLogger, enterpriseEnforcerImpl)
	userAuthRouterImpl := user2.NewUserAuthRouterImpl(sugaredLogger, userAuthHandlerImpl, userAuthOidcHelperImpl)
	policiesCleanUpRepositoryImpl := repository2.NewPoliciesCleanUpRepositoryImpl(db, sugaredLogger)
	cronLoggerImpl := cron.NewCronLoggerImpl(sugaredLogger)
	cleanUpPoliciesServiceImpl := user.NewCleanUpPoliciesServiceImpl(userAuthRepositoryImpl, sugaredLogger, userRepositoryImpl, roleGroupRepositoryImpl, policiesCleanUpRepositoryImpl, cronLoggerImpl)
	userRestHandlerImpl := user2.NewUserRestHandlerImpl(userServiceImpl, validate, sugaredLogger, enterpriseEnforcerImpl, roleGroupServiceImpl, userCommonServiceImpl, cleanUpPoliciesServiceImpl)
	userRouterImpl := user2.NewUserRouterImpl(userRestHandlerImpl)
	genericNoteRepositoryImpl := repository8.NewGenericNoteRepositoryImpl(db)
	genericNoteHistoryRepositoryImpl := repository8.NewGenericNoteHistoryRepositoryImpl(db)
	genericNoteHistoryServiceImpl := genericNotes.NewGenericNoteHistoryServiceImpl(genericNoteHistoryRepositoryImpl, sugaredLogger)
	genericNoteServiceImpl := genericNotes.NewGenericNoteServiceImpl(genericNoteRepositoryImpl, genericNoteHistoryServiceImpl, userRepositoryImpl, sugaredLogger)
	clusterDescriptionRepositoryImpl := repository4.NewClusterDescriptionRepositoryImpl(db, sugaredLogger)
	clusterDescriptionServiceImpl := cluster.NewClusterDescriptionServiceImpl(clusterDescriptionRepositoryImpl, userRepositoryImpl, sugaredLogger)
	helmUserServiceImpl, err := argo.NewHelmUserServiceImpl(sugaredLogger)
	if err != nil {
		return nil, err
	}
	clusterRbacServiceImpl := cluster.NewClusterRbacServiceImpl(environmentServiceImpl, enterpriseEnforcerImpl, clusterServiceImpl, sugaredLogger, userServiceImpl)
	clusterRestHandlerImpl := cluster2.NewClusterRestHandlerImpl(clusterServiceImpl, genericNoteServiceImpl, clusterDescriptionServiceImpl, sugaredLogger, userServiceImpl, validate, enterpriseEnforcerImpl, deleteServiceImpl, helmUserServiceImpl, environmentServiceImpl, clusterRbacServiceImpl)
	clusterRouterImpl := cluster2.NewClusterRouterImpl(clusterRestHandlerImpl)
	dashboardConfig, err := dashboard.GetConfig()
	if err != nil {
		return nil, err
	}
	dashboardRouterImpl := dashboard.NewDashboardRouterImpl(sugaredLogger, dashboardConfig)
	chartGroupDeploymentRepositoryImpl := repository9.NewChartGroupDeploymentRepositoryImpl(db, sugaredLogger)
	clusterInstalledAppsRepositoryImpl := repository6.NewClusterInstalledAppsRepositoryImpl(db, sugaredLogger)
	eaModeDeploymentServiceImpl := EAMode.NewEAModeDeploymentServiceImpl(sugaredLogger, helmAppServiceImpl, appStoreApplicationVersionRepositoryImpl, helmAppClientImpl, installedAppRepositoryImpl, ociRegistryConfigRepositoryImpl)
	chartTemplateServiceImpl := util.NewChartTemplateServiceImpl(sugaredLogger)
	appStoreDeploymentCommonServiceImpl := appStoreDeploymentCommon.NewAppStoreDeploymentCommonServiceImpl(sugaredLogger, appStoreApplicationVersionRepositoryImpl, chartTemplateServiceImpl)
	installedAppVersionHistoryRepositoryImpl := repository6.NewInstalledAppVersionHistoryRepositoryImpl(sugaredLogger, db)
	deploymentServiceTypeConfig, err := service2.GetDeploymentServiceTypeConfig()
	if err != nil {
		return nil, err
	}
	devtronResourceRepositoryImpl := repository10.NewDevtronResourceRepositoryImpl(db, sugaredLogger)
	devtronResourceSchemaRepositoryImpl := repository10.NewDevtronResourceSchemaRepositoryImpl(db, sugaredLogger)
	devtronResourceObjectRepositoryImpl := repository10.NewDevtronResourceObjectRepositoryImpl(sugaredLogger, db)
	devtronResourceSchemaAuditRepositoryImpl := repository10.NewDevtronResourceSchemaAuditRepositoryImpl(sugaredLogger, db)
	devtronResourceObjectAuditRepositoryImpl := repository10.NewDevtronResourceObjectAuditRepositoryImpl(sugaredLogger, db)
	appListingRepositoryQueryBuilder := helper.NewAppListingRepositoryQueryBuilder(sugaredLogger)
	appListingRepositoryImpl := repository5.NewAppListingRepositoryImpl(sugaredLogger, db, appListingRepositoryQueryBuilder, environmentRepositoryImpl)
	devtronResourceServiceImpl, err := devtronResource.NewDevtronResourceServiceImpl(sugaredLogger, devtronResourceRepositoryImpl, devtronResourceSchemaRepositoryImpl, devtronResourceObjectRepositoryImpl, devtronResourceSchemaAuditRepositoryImpl, devtronResourceObjectAuditRepositoryImpl, appRepositoryImpl, pipelineRepositoryImpl, appListingRepositoryImpl, userRepositoryImpl)
	if err != nil {
		return nil, err
	}
	acdConfig, err := argocdServer.GetACDDeploymentConfig()
	if err != nil {
		return nil, err
	}
	gitOpsConfigRepositoryImpl := repository5.NewGitOpsConfigRepositoryImpl(sugaredLogger, db)
	globalEnvVariables, err := util2.GetGlobalEnvVariables()
	if err != nil {
		return nil, err
	}
	gitOpsConfigReadServiceImpl := config.NewGitOpsConfigReadServiceImpl(sugaredLogger, gitOpsConfigRepositoryImpl, userServiceImpl, globalEnvVariables)
	appStoreDeploymentServiceImpl := service2.NewAppStoreDeploymentServiceImpl(sugaredLogger, installedAppRepositoryImpl, chartGroupDeploymentRepositoryImpl, appStoreApplicationVersionRepositoryImpl, environmentRepositoryImpl, clusterInstalledAppsRepositoryImpl, appRepositoryImpl, eaModeDeploymentServiceImpl, eaModeDeploymentServiceImpl, environmentServiceImpl, clusterServiceImpl, helmAppServiceImpl, appStoreDeploymentCommonServiceImpl, installedAppVersionHistoryRepositoryImpl, deploymentServiceTypeConfig, devtronResourceServiceImpl, acdConfig, gitOpsConfigReadServiceImpl)
	installedAppDBServiceImpl := EAMode.NewInstalledAppDBServiceImpl(sugaredLogger, installedAppRepositoryImpl, appRepositoryImpl, userServiceImpl, installedAppVersionHistoryRepositoryImpl)
	attributesServiceImpl := attributes.NewAttributesServiceImpl(sugaredLogger, attributesRepositoryImpl)
	helmAppRestHandlerImpl := client3.NewHelmAppRestHandlerImpl(sugaredLogger, helmAppServiceImpl, enterpriseEnforcerImpl, clusterServiceImpl, enforcerUtilHelmImpl, appStoreDeploymentServiceImpl, installedAppDBServiceImpl, userServiceImpl, attributesServiceImpl, serverEnvConfigServerEnvConfig)
	helmAppRouterImpl := client3.NewHelmAppRouterImpl(helmAppRestHandlerImpl)
	k8sResourceHistoryRepositoryImpl := repository11.NewK8sResourceHistoryRepositoryImpl(db, sugaredLogger)
	k8sResourceHistoryServiceImpl := kubernetesResourceAuditLogs.Newk8sResourceHistoryServiceImpl(k8sResourceHistoryRepositoryImpl, sugaredLogger, appRepositoryImpl, environmentRepositoryImpl)
	argoApplicationServiceImpl := argoApplication.NewArgoApplicationServiceImpl(sugaredLogger, clusterRepositoryImpl, k8sUtilExtended, helmUserServiceImpl, helmAppServiceImpl)
	k8sCommonServiceImpl := k8s2.NewK8sCommonServiceImpl(sugaredLogger, k8sUtilExtended, helmAppServiceImpl, k8sResourceHistoryServiceImpl, clusterServiceImpl, argoApplicationServiceImpl)
	environmentRestHandlerImpl := cluster2.NewEnvironmentRestHandlerImpl(environmentServiceImpl, sugaredLogger, userServiceImpl, validate, enterpriseEnforcerImpl, deleteServiceImpl, k8sUtilExtended, k8sCommonServiceImpl)
	environmentRouterImpl := cluster2.NewEnvironmentRouterImpl(environmentRestHandlerImpl)
	ephemeralContainersRepositoryImpl := repository4.NewEphemeralContainersRepositoryImpl(db)
	ephemeralContainerServiceImpl := cluster.NewEphemeralContainerServiceImpl(ephemeralContainersRepositoryImpl, sugaredLogger)
	terminalSessionHandlerImpl := terminal.NewTerminalSessionHandlerImpl(environmentServiceImpl, clusterServiceImpl, sugaredLogger, k8sUtilExtended, ephemeralContainerServiceImpl, argoApplicationServiceImpl)
	qualifiersMappingRepositoryImpl, err := resourceQualifiers.NewQualifiersMappingRepositoryImpl(db, sugaredLogger)
	if err != nil {
		return nil, err
	}
	devtronResourceSearchableKeyRepositoryImpl := repository10.NewDevtronResourceSearchableKeyRepositoryImpl(sugaredLogger, db)
	devtronResourceSearchableKeyServiceImpl, err := devtronResource.NewDevtronResourceSearchableKeyServiceImpl(sugaredLogger, devtronResourceSearchableKeyRepositoryImpl)
	if err != nil {
		return nil, err
	}
	qualifierMappingServiceImpl, err := resourceQualifiers.NewQualifierMappingServiceImpl(sugaredLogger, qualifiersMappingRepositoryImpl, devtronResourceSearchableKeyServiceImpl)
	if err != nil {
		return nil, err
	}
	timeoutWindowResourceMappingRepositoryImpl := repository3.NewTimeoutWindowResourceMappingRepositoryImpl(db, sugaredLogger)
	timeoutWindowResourceMappingServiceImpl := timeoutWindow.NewTimeoutWindowResourceMappingServiceImpl(sugaredLogger, timeoutWindowResourceMappingRepositoryImpl, timeWindowServiceImpl)
	globalPolicyRepositoryImpl := repository12.NewGlobalPolicyRepositoryImpl(sugaredLogger, db)
	globalPolicySearchableFieldRepositoryImpl := repository12.NewGlobalPolicySearchableFieldRepositoryImpl(sugaredLogger, db)
	globalPolicyDataManagerImpl := globalPolicy.NewGlobalPolicyDataManagerImpl(sugaredLogger, globalPolicyRepositoryImpl, globalPolicySearchableFieldRepositoryImpl)
	deploymentWindowServiceImpl, err := deploymentWindow.NewDeploymentWindowServiceImpl(sugaredLogger, qualifierMappingServiceImpl, timeoutWindowResourceMappingServiceImpl, globalPolicyDataManagerImpl, db, userServiceImpl)
	if err != nil {
		return nil, err
	}
	k8sApplicationServiceImpl, err := application.NewK8sApplicationServiceImpl(sugaredLogger, clusterServiceImpl, pumpImpl, helmAppServiceImpl, k8sUtilExtended, acdAuthConfig, k8sResourceHistoryServiceImpl, k8sCommonServiceImpl, terminalSessionHandlerImpl, ephemeralContainerServiceImpl, ephemeralContainersRepositoryImpl, argoApplicationServiceImpl, deploymentWindowServiceImpl)
	if err != nil {
		return nil, err
	}
	ciPipelineRepositoryImpl := pipelineConfig.NewCiPipelineRepositoryImpl(db, sugaredLogger)
	enforcerUtilImpl := rbac.NewEnforcerUtilImpl(sugaredLogger, teamRepositoryImpl, appRepositoryImpl, environmentRepositoryImpl, pipelineRepositoryImpl, ciPipelineRepositoryImpl, clusterRepositoryImpl, enterpriseEnforcerImpl)
	k8sApplicationRestHandlerImpl := application2.NewK8sApplicationRestHandlerImpl(sugaredLogger, k8sApplicationServiceImpl, pumpImpl, terminalSessionHandlerImpl, enterpriseEnforcerImpl, enforcerUtilHelmImpl, enforcerUtilImpl, helmAppServiceImpl, userServiceImpl, k8sCommonServiceImpl, validate)
	k8sApplicationRouterImpl := application2.NewK8sApplicationRouterImpl(k8sApplicationRestHandlerImpl)
	chartRepositoryRestHandlerImpl := chartRepo2.NewChartRepositoryRestHandlerImpl(sugaredLogger, userServiceImpl, chartRepositoryServiceImpl, enterpriseEnforcerImpl, validate, deleteServiceImpl, attributesServiceImpl)
	chartRepositoryRouterImpl := chartRepo2.NewChartRepositoryRouterImpl(chartRepositoryRestHandlerImpl)
	appStoreServiceImpl := service3.NewAppStoreServiceImpl(sugaredLogger, appStoreApplicationVersionRepositoryImpl)
	appStoreRestHandlerImpl := appStoreDiscover.NewAppStoreRestHandlerImpl(sugaredLogger, userServiceImpl, appStoreServiceImpl, enterpriseEnforcerImpl)
	appStoreDiscoverRouterImpl := appStoreDiscover.NewAppStoreDiscoverRouterImpl(appStoreRestHandlerImpl)
	appStoreVersionValuesRepositoryImpl := appStoreValuesRepository.NewAppStoreVersionValuesRepositoryImpl(sugaredLogger, db)
	appStoreValuesServiceImpl := service4.NewAppStoreValuesServiceImpl(sugaredLogger, appStoreApplicationVersionRepositoryImpl, installedAppRepositoryImpl, appStoreVersionValuesRepositoryImpl, userServiceImpl)
	appStoreValuesRestHandlerImpl := appStoreValues.NewAppStoreValuesRestHandlerImpl(sugaredLogger, userServiceImpl, appStoreValuesServiceImpl)
	appStoreValuesRouterImpl := appStoreValues.NewAppStoreValuesRouterImpl(appStoreValuesRestHandlerImpl)
	appStoreDeploymentRestHandlerImpl := appStoreDeployment.NewAppStoreDeploymentRestHandlerImpl(sugaredLogger, userServiceImpl, enterpriseEnforcerImpl, enforcerUtilImpl, enforcerUtilHelmImpl, appStoreDeploymentServiceImpl, validate, helmAppServiceImpl, appStoreDeploymentCommonServiceImpl, helmUserServiceImpl, attributesServiceImpl)
	appStoreDeploymentRouterImpl := appStoreDeployment.NewAppStoreDeploymentRouterImpl(appStoreDeploymentRestHandlerImpl)
	chartProviderServiceImpl := chartProvider.NewChartProviderServiceImpl(sugaredLogger, chartRepoRepositoryImpl, chartRepositoryServiceImpl, dockerArtifactStoreRepositoryImpl, ociRegistryConfigRepositoryImpl)
	chartProviderRestHandlerImpl := chartProvider2.NewChartProviderRestHandlerImpl(sugaredLogger, userServiceImpl, validate, chartProviderServiceImpl, enterpriseEnforcerImpl)
	chartProviderRouterImpl := chartProvider2.NewChartProviderRouterImpl(chartProviderRestHandlerImpl)
	dockerRegRestHandlerImpl := restHandler.NewDockerRegRestHandlerImpl(dockerRegistryConfigImpl, sugaredLogger, chartProviderServiceImpl, userServiceImpl, validate, enterpriseEnforcerImpl, teamServiceImpl, deleteServiceImpl)
	dockerRegRouterImpl := router.NewDockerRegRouterImpl(dockerRegRestHandlerImpl)
	posthogClient, err := telemetry.NewPosthogClient(sugaredLogger)
	if err != nil {
		return nil, err
	}
	moduleRepositoryImpl := moduleRepo.NewModuleRepositoryImpl(db)
	providerIdentifierServiceImpl := providerIdentifier.NewProviderIdentifierServiceImpl(sugaredLogger)
	telemetryEventClientImpl, err := telemetry.NewTelemetryEventClientImpl(sugaredLogger, httpClient, clusterServiceImpl, k8sUtilExtended, acdAuthConfig, userServiceImpl, attributesRepositoryImpl, ssoLoginServiceImpl, posthogClient, moduleRepositoryImpl, serverDataStoreServerDataStore, userAuditServiceImpl, helmAppClientImpl, installedAppRepositoryImpl, providerIdentifierServiceImpl, cronLoggerImpl)
	if err != nil {
		return nil, err
	}
	dashboardTelemetryRestHandlerImpl := dashboardEvent.NewDashboardTelemetryRestHandlerImpl(sugaredLogger, telemetryEventClientImpl)
	dashboardTelemetryRouterImpl := dashboardEvent.NewDashboardTelemetryRouterImpl(dashboardTelemetryRestHandlerImpl)
	commonDeploymentRestHandlerImpl := appStoreDeployment.NewCommonDeploymentRestHandlerImpl(sugaredLogger, userServiceImpl, enterpriseEnforcerImpl, enforcerUtilImpl, enforcerUtilHelmImpl, appStoreDeploymentServiceImpl, installedAppDBServiceImpl, validate, helmAppServiceImpl, helmAppRestHandlerImpl, helmUserServiceImpl)
	commonDeploymentRouterImpl := appStoreDeployment.NewCommonDeploymentRouterImpl(commonDeploymentRestHandlerImpl)
	externalLinkMonitoringToolRepositoryImpl := externalLink.NewExternalLinkMonitoringToolRepositoryImpl(db)
	externalLinkIdentifierMappingRepositoryImpl := externalLink.NewExternalLinkIdentifierMappingRepositoryImpl(db)
	externalLinkRepositoryImpl := externalLink.NewExternalLinkRepositoryImpl(db)
	externalLinkServiceImpl := externalLink.NewExternalLinkServiceImpl(sugaredLogger, externalLinkMonitoringToolRepositoryImpl, externalLinkIdentifierMappingRepositoryImpl, externalLinkRepositoryImpl)
	externalLinkRestHandlerImpl := externalLink2.NewExternalLinkRestHandlerImpl(sugaredLogger, externalLinkServiceImpl, userServiceImpl, enterpriseEnforcerImpl, enforcerUtilImpl)
	externalLinkRouterImpl := externalLink2.NewExternalLinkRouterImpl(externalLinkRestHandlerImpl)
	moduleActionAuditLogRepositoryImpl := module.NewModuleActionAuditLogRepositoryImpl(db)
	serverCacheServiceImpl := server.NewServerCacheServiceImpl(sugaredLogger, serverEnvConfigServerEnvConfig, serverDataStoreServerDataStore, helmAppServiceImpl)
	moduleEnvConfig, err := module.ParseModuleEnvConfig()
	if err != nil {
		return nil, err
	}
	moduleCacheServiceImpl := module.NewModuleCacheServiceImpl(sugaredLogger, k8sUtilExtended, moduleEnvConfig, serverEnvConfigServerEnvConfig, serverDataStoreServerDataStore, moduleRepositoryImpl, teamServiceImpl)
	moduleServiceHelperImpl := module.NewModuleServiceHelperImpl(serverEnvConfigServerEnvConfig)
	moduleResourceStatusRepositoryImpl := moduleRepo.NewModuleResourceStatusRepositoryImpl(db)
	moduleDataStoreModuleDataStore := moduleDataStore.InitModuleDataStore()
	moduleCronServiceImpl, err := module.NewModuleCronServiceImpl(sugaredLogger, moduleEnvConfig, moduleRepositoryImpl, serverEnvConfigServerEnvConfig, helmAppServiceImpl, moduleServiceHelperImpl, moduleResourceStatusRepositoryImpl, moduleDataStoreModuleDataStore, cronLoggerImpl)
	if err != nil {
		return nil, err
	}
	scanToolMetadataRepositoryImpl := security.NewScanToolMetadataRepositoryImpl(db, sugaredLogger)
	moduleServiceImpl := module.NewModuleServiceImpl(sugaredLogger, serverEnvConfigServerEnvConfig, moduleRepositoryImpl, moduleActionAuditLogRepositoryImpl, helmAppServiceImpl, serverDataStoreServerDataStore, serverCacheServiceImpl, moduleCacheServiceImpl, moduleCronServiceImpl, moduleServiceHelperImpl, moduleResourceStatusRepositoryImpl, scanToolMetadataRepositoryImpl)
	moduleRestHandlerImpl := module2.NewModuleRestHandlerImpl(sugaredLogger, moduleServiceImpl, userServiceImpl, enterpriseEnforcerImpl, validate)
	moduleRouterImpl := module2.NewModuleRouterImpl(moduleRestHandlerImpl)
	serverActionAuditLogRepositoryImpl := server.NewServerActionAuditLogRepositoryImpl(db)
	serverServiceImpl := server.NewServerServiceImpl(sugaredLogger, serverActionAuditLogRepositoryImpl, serverDataStoreServerDataStore, serverEnvConfigServerEnvConfig, helmAppServiceImpl, moduleRepositoryImpl)
	serverRestHandlerImpl := server2.NewServerRestHandlerImpl(sugaredLogger, serverServiceImpl, userServiceImpl, enterpriseEnforcerImpl, validate)
	serverRouterImpl := server2.NewServerRouterImpl(serverRestHandlerImpl)
	apiTokenSecretServiceImpl, err := apiToken.NewApiTokenSecretServiceImpl(sugaredLogger, attributesServiceImpl, apiTokenSecretStore)
	if err != nil {
		return nil, err
	}
	apiTokenRepositoryImpl := apiToken.NewApiTokenRepositoryImpl(db)
	apiTokenServiceImpl, err := apiToken.NewApiTokenServiceImpl(sugaredLogger, apiTokenSecretServiceImpl, userServiceImpl, userAuditServiceImpl, apiTokenRepositoryImpl)
	if err != nil {
		return nil, err
	}
	apiTokenRestHandlerImpl := apiToken2.NewApiTokenRestHandlerImpl(sugaredLogger, apiTokenServiceImpl, userServiceImpl, enterpriseEnforcerImpl, validate)
	apiTokenRouterImpl := apiToken2.NewApiTokenRouterImpl(apiTokenRestHandlerImpl)
	clusterCronServiceImpl, err := cluster.NewClusterCronServiceImpl(sugaredLogger, clusterServiceImpl, cronLoggerImpl)
	if err != nil {
		return nil, err
	}
	k8sCapacityServiceImpl := capacity.NewK8sCapacityServiceImpl(sugaredLogger, clusterServiceImpl, k8sApplicationServiceImpl, k8sUtilExtended, k8sCommonServiceImpl, clusterCronServiceImpl)
	k8sCapacityRestHandlerImpl := capacity2.NewK8sCapacityRestHandlerImpl(sugaredLogger, k8sCapacityServiceImpl, userServiceImpl, enterpriseEnforcerImpl, clusterServiceImpl, environmentServiceImpl, clusterRbacServiceImpl)
	k8sCapacityRouterImpl := capacity2.NewK8sCapacityRouterImpl(k8sCapacityRestHandlerImpl)
	webhookHelmServiceImpl := webhookHelm.NewWebhookHelmServiceImpl(sugaredLogger, helmAppServiceImpl, clusterServiceImpl, chartRepositoryServiceImpl, attributesServiceImpl)
	webhookHelmRestHandlerImpl := webhookHelm2.NewWebhookHelmRestHandlerImpl(sugaredLogger, webhookHelmServiceImpl, userServiceImpl, enterpriseEnforcerImpl, validate)
	webhookHelmRouterImpl := webhookHelm2.NewWebhookHelmRouterImpl(webhookHelmRestHandlerImpl)
	userAttributesRepositoryImpl := repository5.NewUserAttributesRepositoryImpl(db)
	userAttributesServiceImpl := attributes.NewUserAttributesServiceImpl(sugaredLogger, userAttributesRepositoryImpl)
	userAttributesRestHandlerImpl := restHandler.NewUserAttributesRestHandlerImpl(sugaredLogger, enterpriseEnforcerImpl, userServiceImpl, userAttributesServiceImpl)
	userAttributesRouterImpl := router.NewUserAttributesRouterImpl(userAttributesRestHandlerImpl)
	telemetryRestHandlerImpl := restHandler.NewTelemetryRestHandlerImpl(sugaredLogger, telemetryEventClientImpl, enterpriseEnforcerImpl, userServiceImpl)
	telemetryRouterImpl := router.NewTelemetryRouterImpl(sugaredLogger, telemetryRestHandlerImpl)
	terminalAccessRepositoryImpl := repository5.NewTerminalAccessRepositoryImpl(db, sugaredLogger)
	userTerminalSessionConfig, err := clusterTerminalAccess.GetTerminalAccessConfig()
	if err != nil {
		return nil, err
	}
	userTerminalAccessServiceImpl, err := clusterTerminalAccess.NewUserTerminalAccessServiceImpl(sugaredLogger, terminalAccessRepositoryImpl, userTerminalSessionConfig, k8sCommonServiceImpl, terminalSessionHandlerImpl, k8sCapacityServiceImpl, k8sUtilExtended, cronLoggerImpl)
	if err != nil {
		return nil, err
	}
	userTerminalAccessRestHandlerImpl := terminal2.NewUserTerminalAccessRestHandlerImpl(sugaredLogger, userTerminalAccessServiceImpl, enterpriseEnforcerImpl, userServiceImpl, validate)
	userTerminalAccessRouterImpl := terminal2.NewUserTerminalAccessRouterImpl(userTerminalAccessRestHandlerImpl)
	attributesRestHandlerImpl := restHandler.NewAttributesRestHandlerImpl(sugaredLogger, enterpriseEnforcerImpl, userServiceImpl, attributesServiceImpl)
	attributesRouterImpl := router.NewAttributesRouterImpl(attributesRestHandlerImpl)
	appLabelRepositoryImpl := pipelineConfig.NewAppLabelRepositoryImpl(db)
	materialRepositoryImpl := pipelineConfig.NewMaterialRepositoryImpl(db)
	appCrudOperationServiceImpl := app2.NewAppCrudOperationServiceImpl(appLabelRepositoryImpl, sugaredLogger, appRepositoryImpl, userRepositoryImpl, installedAppRepositoryImpl, teamRepositoryImpl, genericNoteServiceImpl, materialRepositoryImpl)
	appInfoRestHandlerImpl := appInfo.NewAppInfoRestHandlerImpl(sugaredLogger, appCrudOperationServiceImpl, userServiceImpl, validate, enforcerUtilImpl, enterpriseEnforcerImpl, helmAppServiceImpl, enforcerUtilHelmImpl, genericNoteServiceImpl)
	appInfoRouterImpl := appInfo2.NewAppInfoRouterImpl(sugaredLogger, appInfoRestHandlerImpl)
	appFilteringRestHandlerImpl := appList.NewAppFilteringRestHandlerImpl(sugaredLogger, teamServiceImpl, enterpriseEnforcerImpl, userServiceImpl, clusterServiceImpl, environmentServiceImpl)
	appFilteringRouterImpl := appList2.NewAppFilteringRouterImpl(appFilteringRestHandlerImpl)
	appRouterEAModeImpl := app3.NewAppRouterEAModeImpl(appInfoRouterImpl, appFilteringRouterImpl)
	rbacPolicyResourceDetailRepositoryImpl := repository2.NewRbacPolicyResourceDetailRepositoryImpl(sugaredLogger, db)
	rbacRoleResourceDetailRepositoryImpl := repository2.NewRbacRoleResourceDetailRepositoryImpl(sugaredLogger, db)
	rbacRoleServiceImpl := user.NewRbacRoleServiceImpl(sugaredLogger, rbacPolicyResourceDetailRepositoryImpl, rbacRoleResourceDetailRepositoryImpl, rbacRoleDataRepositoryImpl, rbacPolicyDataRepositoryImpl, rbacDataCacheFactoryImpl, userAuthRepositoryImpl, userCommonServiceImpl)
	defaultRbacRoleDataRepositoryImpl := repository2.NewDefaultRbacRoleDataRepositoryImpl(sugaredLogger, db)
	defaultRbacRoleServiceImpl := user.NewDefaultRbacRoleServiceImpl(sugaredLogger, defaultRbacRoleDataRepositoryImpl, rbacRoleServiceImpl)
	rbacRoleRestHandlerImpl := user2.NewRbacRoleHandlerImpl(sugaredLogger, validate, rbacRoleServiceImpl, userServiceImpl, enterpriseEnforcerImpl, enforcerUtilImpl, defaultRbacRoleServiceImpl)
	rbacRoleRouterImpl := user2.NewRbacRoleRouterImpl(sugaredLogger, validate, rbacRoleRestHandlerImpl)
	authorisationConfigRestHandlerImpl := globalConfig.NewGlobalAuthorisationConfigRestHandlerImpl(validate, sugaredLogger, enterpriseEnforcerImpl, userServiceImpl, globalAuthorisationConfigServiceImpl, userCommonServiceImpl)
	authorisationConfigRouterImpl := globalConfig.NewGlobalConfigAuthorisationConfigRouterImpl(authorisationConfigRestHandlerImpl)
	muxRouter := NewMuxRouter(sugaredLogger, ssoLoginRouterImpl, teamRouterImpl, userAuthRouterImpl, userRouterImpl, clusterRouterImpl, dashboardRouterImpl, helmAppRouterImpl, environmentRouterImpl, k8sApplicationRouterImpl, chartRepositoryRouterImpl, appStoreDiscoverRouterImpl, appStoreValuesRouterImpl, appStoreDeploymentRouterImpl, chartProviderRouterImpl, dockerRegRouterImpl, dashboardTelemetryRouterImpl, commonDeploymentRouterImpl, externalLinkRouterImpl, moduleRouterImpl, serverRouterImpl, apiTokenRouterImpl, k8sCapacityRouterImpl, webhookHelmRouterImpl, userAttributesRouterImpl, telemetryRouterImpl, userTerminalAccessRouterImpl, attributesRouterImpl, appRouterEAModeImpl, rbacRoleRouterImpl, authorisationConfigRouterImpl)
	mainApp := NewApp(db, sessionManager, muxRouter, telemetryEventClientImpl, posthogClient, sugaredLogger, userServiceImpl)
	return mainApp, nil
}
