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

package restHandler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/devtron-labs/devtron/pkg/deployment/deployedApp"
	bean2 "github.com/devtron-labs/devtron/pkg/deployment/deployedApp/bean"
	"github.com/devtron-labs/devtron/pkg/deployment/trigger/devtronApps"
	bean3 "github.com/devtron-labs/devtron/pkg/deployment/trigger/devtronApps/bean"
	"github.com/devtron-labs/devtron/pkg/eventProcessor/out"
	bean4 "github.com/devtron-labs/devtron/pkg/eventProcessor/out/bean"
	"github.com/devtron-labs/devtron/pkg/workflow/dag"
	"net/http"
	"strconv"

	"github.com/devtron-labs/devtron/pkg/app"
	"github.com/devtron-labs/devtron/pkg/auth/authorisation/casbin"
	"github.com/devtron-labs/devtron/pkg/auth/user"
	"github.com/devtron-labs/devtron/util"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"

	"github.com/devtron-labs/devtron/api/bean"
	"github.com/devtron-labs/devtron/api/restHandler/common"
	"github.com/devtron-labs/devtron/pkg/deploymentGroup"
	"github.com/devtron-labs/devtron/pkg/pipeline"
	"github.com/devtron-labs/devtron/pkg/team"
	"github.com/devtron-labs/devtron/util/argo"
	"github.com/devtron-labs/devtron/util/rbac"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

type PipelineTriggerRestHandler interface {
	OverrideConfig(w http.ResponseWriter, r *http.Request)
	ReleaseStatusUpdate(w http.ResponseWriter, r *http.Request)
	StartStopApp(w http.ResponseWriter, r *http.Request)
	StartStopDeploymentGroup(w http.ResponseWriter, r *http.Request)
	GetAllLatestDeploymentConfiguration(w http.ResponseWriter, r *http.Request)
	RotatePods(w http.ResponseWriter, r *http.Request)
	DownloadManifest(w http.ResponseWriter, r *http.Request)
	DownloadManifestForSpecificTrigger(w http.ResponseWriter, r *http.Request)
}

type PipelineTriggerRestHandlerImpl struct {
	appService                  app.AppService
	userAuthService             user.UserService
	validator                   *validator.Validate
	enforcer                    casbin.Enforcer
	teamService                 team.TeamService
	logger                      *zap.SugaredLogger
	workflowDagExecutor         dag.WorkflowDagExecutor
	enforcerUtil                rbac.EnforcerUtil
	deploymentGroupService      deploymentGroup.DeploymentGroupService
	argoUserService             argo.ArgoUserService
	deploymentConfigService     pipeline.DeploymentConfigService
	deployedAppService          deployedApp.DeployedAppService
	cdTriggerService            devtronApps.TriggerService
	workflowEventPublishService out.WorkflowEventPublishService
}

func NewPipelineRestHandler(appService app.AppService, userAuthService user.UserService, validator *validator.Validate,
	enforcer casbin.Enforcer, teamService team.TeamService, logger *zap.SugaredLogger, enforcerUtil rbac.EnforcerUtil,
	workflowDagExecutor dag.WorkflowDagExecutor, deploymentGroupService deploymentGroup.DeploymentGroupService,
	argoUserService argo.ArgoUserService, deploymentConfigService pipeline.DeploymentConfigService,
	deployedAppService deployedApp.DeployedAppService,
	cdTriggerService devtronApps.TriggerService,
	workflowEventPublishService out.WorkflowEventPublishService) *PipelineTriggerRestHandlerImpl {
	pipelineHandler := &PipelineTriggerRestHandlerImpl{
		appService:                  appService,
		userAuthService:             userAuthService,
		validator:                   validator,
		enforcer:                    enforcer,
		teamService:                 teamService,
		logger:                      logger,
		workflowDagExecutor:         workflowDagExecutor,
		enforcerUtil:                enforcerUtil,
		deploymentGroupService:      deploymentGroupService,
		argoUserService:             argoUserService,
		deploymentConfigService:     deploymentConfigService,
		deployedAppService:          deployedAppService,
		cdTriggerService:            cdTriggerService,
		workflowEventPublishService: workflowEventPublishService,
	}
	return pipelineHandler
}

func (handler PipelineTriggerRestHandlerImpl) OverrideConfig(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	userId, err := handler.userAuthService.GetLoggedInUser(r)
	if userId == 0 || err != nil {
		common.WriteJsonResp(w, err, "Unauthorized User", http.StatusUnauthorized)
		return
	}
	var overrideRequest bean.ValuesOverrideRequest
	err = decoder.Decode(&overrideRequest)
	if err != nil {
		handler.logger.Errorw("request err, OverrideConfig", "err", err, "payload", overrideRequest)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	overrideRequest.UserId = userId
	handler.logger.Infow("request for OverrideConfig", "payload", overrideRequest)
	err = handler.validator.Struct(overrideRequest)
	if err != nil {
		handler.logger.Errorw("request err, OverrideConfig", "err", err, "payload", overrideRequest)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	token := r.Header.Get("token")

	//rbac block starts from here
	object := handler.enforcerUtil.GetAppRBACNameByAppId(overrideRequest.AppId)
	if ok := handler.enforcer.Enforce(token, casbin.ResourceApplications, casbin.ActionTrigger, object); !ok {
		common.WriteJsonResp(w, fmt.Errorf("unauthorized user"), "Unauthorized User", http.StatusForbidden)
		return
	}
	object = handler.enforcerUtil.GetAppRBACByAppIdAndPipelineId(overrideRequest.AppId, overrideRequest.PipelineId)
	if ok := handler.enforcer.Enforce(token, casbin.ResourceEnvironment, casbin.ActionTrigger, object); !ok {
		common.WriteJsonResp(w, fmt.Errorf("unauthorized user"), "Unauthorized User", http.StatusForbidden)
		return
	}
	//rback block ends here
	acdToken, err := handler.argoUserService.GetLatestDevtronArgoCdUserToken()
	if err != nil {
		handler.logger.Errorw("error in getting acd token", "err", err)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	ctx := context.WithValue(r.Context(), "token", acdToken)
	_, span := otel.Tracer("orchestrator").Start(ctx, "workflowDagExecutor.ManualCdTrigger")
	triggerContext := bean3.TriggerContext{
		Context: ctx,
	}
	mergeResp, helmPackageName, err := handler.cdTriggerService.ManualCdTrigger(triggerContext, &overrideRequest)
	span.End()
	if err != nil {
		handler.logger.Errorw("request err, OverrideConfig", "err", err, "payload", overrideRequest)
		common.WriteJsonResp(w, err, err.Error(), http.StatusInternalServerError)
		return
	}
	res := map[string]interface{}{"releaseId": mergeResp, "helmPackageName": helmPackageName}
	common.WriteJsonResp(w, err, res, http.StatusOK)
}

func (handler PipelineTriggerRestHandlerImpl) RotatePods(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	userId, err := handler.userAuthService.GetLoggedInUser(r)
	if userId == 0 || err != nil {
		common.WriteJsonResp(w, err, "Unauthorized User", http.StatusUnauthorized)
		return
	}
	var podRotateRequest bean2.PodRotateRequest
	err = decoder.Decode(&podRotateRequest)
	if err != nil {
		handler.logger.Errorw("request err, RotatePods", "err", err, "payload", podRotateRequest)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	podRotateRequest.UserId = userId
	handler.logger.Infow("request payload, RotatePods", "err", err, "payload", podRotateRequest)
	err = handler.validator.Struct(podRotateRequest)
	if err != nil {
		handler.logger.Errorw("validation err, RotatePods", "err", err, "payload", podRotateRequest)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	token := r.Header.Get("token")
	object := handler.enforcerUtil.GetAppRBACNameByAppId(podRotateRequest.AppId)
	if ok := handler.enforcer.Enforce(token, casbin.ResourceApplications, casbin.ActionTrigger, object); !ok {
		common.WriteJsonResp(w, fmt.Errorf("unauthorized user"), "Unauthorized User", http.StatusForbidden)
		return
	}
	object = handler.enforcerUtil.GetEnvRBACNameByAppId(podRotateRequest.AppId, podRotateRequest.EnvironmentId)
	if ok := handler.enforcer.Enforce(token, casbin.ResourceEnvironment, casbin.ActionTrigger, object); !ok {
		common.WriteJsonResp(w, fmt.Errorf("unauthorized user"), "Unauthorized User", http.StatusForbidden)
		return
	}
	rotatePodResponse, err := handler.deployedAppService.RotatePods(r.Context(), &podRotateRequest)
	if err != nil {
		handler.logger.Errorw("service err, RotatePods", "err", err, "payload", podRotateRequest)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	common.WriteJsonResp(w, nil, rotatePodResponse, http.StatusOK)
}

func (handler PipelineTriggerRestHandlerImpl) StartStopApp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	userId, err := handler.userAuthService.GetLoggedInUser(r)
	if userId == 0 || err != nil {
		common.WriteJsonResp(w, err, "Unauthorized User", http.StatusUnauthorized)
		return
	}
	var overrideRequest bean2.StopAppRequest
	err = decoder.Decode(&overrideRequest)
	if err != nil {
		handler.logger.Errorw("request err, StartStopApp", "err", err, "payload", overrideRequest)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	overrideRequest.UserId = userId
	handler.logger.Infow("request payload, StartStopApp", "err", err, "payload", overrideRequest)
	err = handler.validator.Struct(overrideRequest)
	if err != nil {
		handler.logger.Errorw("validation err, StartStopApp", "err", err, "payload", overrideRequest)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	token := r.Header.Get("token")
	//rbac block starts from here
	object := handler.enforcerUtil.GetAppRBACNameByAppId(overrideRequest.AppId)
	if ok := handler.enforcer.Enforce(token, casbin.ResourceApplications, casbin.ActionTrigger, object); !ok {
		common.WriteJsonResp(w, fmt.Errorf("unauthorized user"), "Unauthorized User", http.StatusForbidden)
		return
	}
	object = handler.enforcerUtil.GetEnvRBACNameByAppId(overrideRequest.AppId, overrideRequest.EnvironmentId)
	if ok := handler.enforcer.Enforce(token, casbin.ResourceEnvironment, casbin.ActionTrigger, object); !ok {
		common.WriteJsonResp(w, fmt.Errorf("unauthorized user"), "Unauthorized User", http.StatusForbidden)
		return
	}
	//rback block ends here
	acdToken, err := handler.argoUserService.GetLatestDevtronArgoCdUserToken()
	if err != nil {
		handler.logger.Errorw("error in getting acd token", "err", err)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	ctx := context.WithValue(r.Context(), "token", acdToken)
	mergeResp, err := handler.deployedAppService.StopStartApp(ctx, &overrideRequest)
	if err != nil {
		handler.logger.Errorw("service err, StartStopApp", "err", err, "payload", overrideRequest)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	res := map[string]interface{}{"releaseId": mergeResp}
	common.WriteJsonResp(w, err, res, http.StatusOK)
}

func (handler PipelineTriggerRestHandlerImpl) StartStopDeploymentGroup(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	userId, err := handler.userAuthService.GetLoggedInUser(r)
	if userId == 0 || err != nil {
		common.WriteJsonResp(w, err, "Unauthorized User", http.StatusUnauthorized)
		return
	}
	var stopDeploymentGroupRequest bean4.StopDeploymentGroupRequest
	err = decoder.Decode(&stopDeploymentGroupRequest)
	if err != nil {
		handler.logger.Errorw("request err, StartStopDeploymentGroup", "err", err, "payload", stopDeploymentGroupRequest)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	stopDeploymentGroupRequest.UserId = userId
	err = handler.validator.Struct(stopDeploymentGroupRequest)
	if err != nil {
		handler.logger.Errorw("validation err, StartStopDeploymentGroup", "err", err, "payload", stopDeploymentGroupRequest)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	handler.logger.Infow("request payload, StartStopDeploymentGroup", "err", err, "payload", stopDeploymentGroupRequest)

	//rbac block starts from here
	dg, err := handler.deploymentGroupService.GetDeploymentGroupById(stopDeploymentGroupRequest.DeploymentGroupId)
	if err != nil {
		handler.logger.Errorw("request err, StartStopDeploymentGroup", "err", err, "payload", stopDeploymentGroupRequest)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	token := r.Header.Get("token")
	// RBAC enforcer applying
	object := handler.enforcerUtil.GetTeamRBACByCiPipelineId(dg.CiPipelineId)
	if ok := handler.enforcer.Enforce(token, casbin.ResourceApplications, casbin.ActionTrigger, object); !ok {
		common.WriteJsonResp(w, fmt.Errorf("unauthorized user"), "Unauthorized User", http.StatusForbidden)
		return
	}
	object = handler.enforcerUtil.GetEnvRBACNameByCiPipelineIdAndEnvId(dg.CiPipelineId, dg.EnvironmentId)
	if ok := handler.enforcer.Enforce(token, casbin.ResourceEnvironment, casbin.ActionTrigger, object); !ok {
		common.WriteJsonResp(w, fmt.Errorf("unauthorized user"), "Unauthorized User", http.StatusForbidden)
		return
	}
	//rback block ends here
	res, err := handler.workflowEventPublishService.TriggerBulkHibernateAsync(stopDeploymentGroupRequest)
	if err != nil {
		handler.logger.Errorw("service err, StartStopDeploymentGroup", "err", err, "payload", stopDeploymentGroupRequest)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	common.WriteJsonResp(w, err, res, http.StatusOK)
}

func (handler PipelineTriggerRestHandlerImpl) ReleaseStatusUpdate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	userId, err := handler.userAuthService.GetLoggedInUser(r)
	if userId == 0 || err != nil {
		common.WriteJsonResp(w, err, "Unauthorized User", http.StatusUnauthorized)
		return
	}
	var releaseStatusUpdateRequest bean.ReleaseStatusUpdateRequest
	err = decoder.Decode(&releaseStatusUpdateRequest)
	if err != nil {
		handler.logger.Errorw("request err, ReleaseStatusUpdate", "err", err, "payload", releaseStatusUpdateRequest)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	handler.logger.Infow("request payload, ReleaseStatusUpdate, override request ----", "err", err, "payload", releaseStatusUpdateRequest)
	res, err := handler.appService.UpdateReleaseStatus(&releaseStatusUpdateRequest)
	if err != nil {
		handler.logger.Errorw("service err, ReleaseStatusUpdate", "err", err, "payload", releaseStatusUpdateRequest)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	m := map[string]bool{
		"status": res}
	resJson, err := json.Marshal(m)
	if err != nil {
		handler.logger.Errorw("marshal err, ReleaseStatusUpdate", "err", err, "payload", m)
	}
	common.WriteJsonResp(w, err, resJson, http.StatusOK)
}

func (handler PipelineTriggerRestHandlerImpl) GetAllLatestDeploymentConfiguration(w http.ResponseWriter, r *http.Request) {
	userId, err := handler.userAuthService.GetLoggedInUser(r)
	if userId == 0 || err != nil {
		common.WriteJsonResp(w, err, "Unauthorized User", http.StatusUnauthorized)
		return
	}
	handler.logger.Infow("request payload, GetAllLatestDeploymentConfiguration")

	vars := mux.Vars(r)
	appId, err := strconv.Atoi(vars["appId"])
	if err != nil {
		handler.logger.Errorw("request err, GetAllDeployedConfigurationHistoryForSpecificWfrIdForPipeline", "err", err, "appId", appId)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	pipelineId, err := strconv.Atoi(vars["pipelineId"])
	if err != nil {
		handler.logger.Errorw("request err, GetAllDeployedConfigurationHistoryForSpecificWfrIdForPipeline", "err", err, "pipelineId", pipelineId)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	//RBAC START
	token := r.Header.Get("token")
	resourceName := handler.enforcerUtil.GetAppRBACNameByAppId(appId)
	if ok := handler.enforcer.Enforce(token, casbin.ResourceApplications, casbin.ActionGet, resourceName); !ok {
		common.WriteJsonResp(w, fmt.Errorf("unauthorized user"), "Unauthorized User", http.StatusForbidden)
		return
	}
	//RBAC END
	isSuperAdmin := handler.enforcer.Enforce(token, casbin.ResourceGlobal, casbin.ActionGet, "*")
	ctx := r.Context()
	ctx = util.SetSuperAdminInContext(ctx, isSuperAdmin)
	//checking if user has admin access
	userHasAdminAccess := handler.enforcer.Enforce(token, casbin.ResourceApplications, casbin.ActionUpdate, resourceName)
	allDeploymentconfig, err := handler.deploymentConfigService.GetLatestDeploymentConfigurationByPipelineId(ctx, pipelineId, userHasAdminAccess)
	if err != nil {
		handler.logger.Errorw("error in getting latest deployment config, GetAllDeployedConfigurationHistoryForSpecificWfrIdForPipeline", "err", err, "pipelineId", pipelineId)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	common.WriteJsonResp(w, nil, allDeploymentconfig, http.StatusOK)
}

func (handler PipelineTriggerRestHandlerImpl) DownloadManifest(w http.ResponseWriter, r *http.Request) {

	userId, err := handler.userAuthService.GetLoggedInUser(r)
	if userId == 0 || err != nil {
		common.WriteJsonResp(w, err, "Unauthorized User", http.StatusUnauthorized)
		return
	}
	token := r.Header.Get("token")

	vars := mux.Vars(r)
	appId, err := strconv.Atoi(vars["appId"])
	if err != nil {
		handler.logger.Errorw("request err, DownloadManifest", "err", err, "appId", appId)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	envId, err := strconv.Atoi(vars["envId"])
	if err != nil {
		handler.logger.Errorw("request err, DownloadManifest", "err", err, "envId", envId)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	runner := vars["runner"]

	object := handler.enforcerUtil.GetEnvRBACNameByAppId(appId, envId)
	if ok := handler.enforcer.Enforce(token, casbin.ResourceEnvironment, casbin.ActionTrigger, object); !ok {
		common.WriteJsonResp(w, fmt.Errorf("unauthorized user"), "Unauthorized User", http.StatusForbidden)
		return
	}

	manifestByteArr, err := handler.appService.GetLatestDeployedManifestByPipelineId(appId, envId, runner, context.Background())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(manifestByteArr)
	w.Header().Set("Content-Type", "application/octet-stream")
	return
}

func (handler PipelineTriggerRestHandlerImpl) DownloadManifestForSpecificTrigger(w http.ResponseWriter, r *http.Request) {

	userId, err := handler.userAuthService.GetLoggedInUser(r)
	if userId == 0 || err != nil {
		common.WriteJsonResp(w, err, "Unauthorized User", http.StatusUnauthorized)
		return
	}
	token := r.Header.Get("token")

	vars := mux.Vars(r)
	appId, err := strconv.Atoi(vars["appId"])
	if err != nil {
		handler.logger.Errorw("request err, DownloadManifest", "err", err, "appId", appId)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	envId, err := strconv.Atoi(vars["envId"])
	if err != nil {
		handler.logger.Errorw("request err, DownloadManifest", "err", err, "envId", envId)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	cdWorkflowId, err := strconv.Atoi(vars["cd_workflow_id"])
	if err != nil {
		handler.logger.Errorw("request err, DownloadManifest", "err", err, "envId", envId)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	object := handler.enforcerUtil.GetEnvRBACNameByAppId(appId, envId)
	if ok := handler.enforcer.Enforce(token, casbin.ResourceEnvironment, casbin.ActionTrigger, object); !ok {
		common.WriteJsonResp(w, fmt.Errorf("unauthorized user"), "Unauthorized User", http.StatusForbidden)
		return
	}
	manifestByteArr, err := handler.appService.GetDeployedManifestByPipelineIdAndCDWorkflowId(cdWorkflowId, context.Background())
	w.WriteHeader(http.StatusOK)
	w.Write(manifestByteArr)
	w.Header().Set("Content-Type", "application/octet-stream")
}
