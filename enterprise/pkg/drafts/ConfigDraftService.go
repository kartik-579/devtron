package drafts

import (
	"context"
	"encoding/json"
	"errors"
	bean2 "github.com/devtron-labs/devtron/api/bean"
	"github.com/devtron-labs/devtron/pkg/chart"
	"github.com/devtron-labs/devtron/pkg/pipeline"
	"go.uber.org/zap"
	"time"
)

type ConfigDraftService interface {
	CreateDraft(request ConfigDraftRequest) (*ConfigDraftResponse, error)
	AddDraftVersion(request ConfigDraftVersionRequest) (int, error)
	UpdateDraftState(draftId int, draftVersionId int, toUpdateDraftState DraftState, userId int32) (*DraftVersionDto, error)
	GetDraftVersionMetadata(draftId int) (*DraftVersionMetadataResponse, error) // would return version timestamp and user email id
	GetDraftComments(draftId int) (*DraftVersionCommentResponse, error)
	GetDrafts(appId int, envId int, resourceType DraftResourceType) ([]AppConfigDraft, error) // need to take care of secret data
	GetDraftById(draftId int) (*ConfigDraftResponse, error)                                   //  need to send ** in case of view only user for Secret data
	GetDraftByDraftVersionId(draftVersionId int) (*ConfigDraftResponse, error)
	ApproveDraft(draftId int, draftVersionId int, userId int32) error
	DeleteComment(draftId int, draftCommentId int, userId int32) error
}

type ConfigDraftServiceImpl struct {
	logger                  *zap.SugaredLogger
	configDraftRepository   ConfigDraftRepository
	configMapService        pipeline.ConfigMapService
	chartService            chart.ChartService
	propertiesConfigService pipeline.PropertiesConfigService
}

func NewConfigDraftServiceImpl(logger *zap.SugaredLogger, configDraftRepository ConfigDraftRepository, chartService chart.ChartService, propertiesConfigService pipeline.PropertiesConfigService) *ConfigDraftServiceImpl {
	return &ConfigDraftServiceImpl{
		logger:                  logger,
		configDraftRepository:   configDraftRepository,
		chartService:            chartService,
		propertiesConfigService: propertiesConfigService,
	}
}

func (impl ConfigDraftServiceImpl) CreateDraft(request ConfigDraftRequest) (*ConfigDraftResponse, error) {
	return impl.configDraftRepository.CreateConfigDraft(request)
}

func (impl ConfigDraftServiceImpl) AddDraftVersion(request ConfigDraftVersionRequest) (int, error) {
	latestDraftVersion, err := impl.configDraftRepository.GetLatestDraftVersionId(request.DraftId)
	if err != nil {
		return 0, err
	}
	lastDraftVersionId := request.LastDraftVersionId
	if latestDraftVersion > lastDraftVersionId {
		return 0, errors.New("last-version-outdated")
	}

	currentTime := time.Now()
	if len(request.Data) > 0 {
		draftVersionDto := DraftVersionDto{}
		draftVersionDto.DraftId = request.DraftId
		draftVersionDto.Data = request.Data
		draftVersionDto.Action = request.Action
		draftVersionDto.UserId = request.UserId
		draftVersionDto.CreatedOn = currentTime
		draftVersionId, err := impl.configDraftRepository.SaveDraftVersion(draftVersionDto)
		if err != nil {
			return 0, err
		}
		lastDraftVersionId = draftVersionId
	}

	if len(request.UserComment) > 0 {
		draftVersionCommentDto := &DraftVersionCommentDto{}
		draftVersionCommentDto.DraftId = request.DraftId
		draftVersionCommentDto.DraftVersionId = lastDraftVersionId
		draftVersionCommentDto.Comment = request.UserComment
		draftVersionCommentDto.Active = true
		draftVersionCommentDto.CreatedBy = request.UserId
		draftVersionCommentDto.UpdatedBy = request.UserId
		draftVersionCommentDto.CreatedOn = currentTime
		draftVersionCommentDto.UpdatedOn = currentTime
		err := impl.configDraftRepository.SaveDraftVersionComment(draftVersionCommentDto)
		if err != nil {
			return 0, err
		}
	}
	return lastDraftVersionId, nil
}

func (impl ConfigDraftServiceImpl) UpdateDraftState(draftId int, draftVersionId int, toUpdateDraftState DraftState, userId int32) (*DraftVersionDto, error) {
	impl.logger.Infow("updating draft state", "draftId", draftId, "toUpdateDraftState", toUpdateDraftState, "userId", userId)
	// check app config draft is enabled or not ??
	latestDraftVersion, err := impl.validateDraftAction(draftId, draftVersionId, toUpdateDraftState)
	if err != nil {
		return nil, err
	}
	err = impl.configDraftRepository.UpdateDraftState(draftId, toUpdateDraftState, userId)
	return latestDraftVersion, err
}

func (impl ConfigDraftServiceImpl) validateDraftAction(draftId int, draftVersionId int, toUpdateDraftState DraftState) (*DraftVersionDto, error) {
	latestDraftVersion, err := impl.configDraftRepository.GetLatestConfigDraft(draftId)
	if err != nil {
		return nil, err
	}
	if latestDraftVersion.Id != draftVersionId { // needed for current scope
		return nil, errors.New("last-version-outdated")
	}
	draftMetadataDto, err := impl.configDraftRepository.GetDraftMetadataById(draftId)
	if err != nil {
		return nil, err
	}
	draftCurrentState := draftMetadataDto.DraftState
	if draftCurrentState.IsTerminal() {
		impl.logger.Errorw("draft is already in terminal state", "draftId", draftId, "draftCurrentState", draftCurrentState)
		return nil, errors.New("already-in-terminal-state")
	}
	if toUpdateDraftState == PublishedDraftState && draftCurrentState != AwaitApprovalDraftState {
		impl.logger.Errorw("draft is not in await Approval state", "draftId", draftId, "draftCurrentState", draftCurrentState)
		return nil, errors.New("approval-request-not-raised")
	}
	return latestDraftVersion, nil
}

func (impl ConfigDraftServiceImpl) GetDraftVersionMetadata(draftId int) (*DraftVersionMetadataResponse, error) {
	draftVersionDtos, err := impl.configDraftRepository.GetDraftVersionsMetadata(draftId)
	if err != nil {
		return nil, err
	}
	var draftVersions []DraftVersionMetadata
	for _, draftVersionDto := range draftVersionDtos {
		versionMetadata := draftVersionDto.ConvertToDraftVersionMetadata()
		draftVersions = append(draftVersions, versionMetadata)
	}

	//TODO need to set email id of user
	response := &DraftVersionMetadataResponse{}
	response.DraftId = draftId
	response.DraftVersions = draftVersions
	return response, nil
}

func (impl ConfigDraftServiceImpl) GetDraftComments(draftId int) (*DraftVersionCommentResponse, error) {
	draftComments, err := impl.configDraftRepository.GetDraftVersionComments(draftId)
	if err != nil {
		return nil, err
	}
	var draftVersionVsComments map[int][]UserCommentMetadata
	for _, draftComment := range draftComments {
		draftVersionId := draftComment.DraftVersionId
		userComment := draftComment.ConvertToDraftVersionComment() //TODO email id is not set
		commentMetadataArray, ok := draftVersionVsComments[draftVersionId]
		commentMetadataArray = append(commentMetadataArray, userComment)
		if !ok {
			draftVersionVsComments[draftVersionId] = commentMetadataArray
		}
	}
	var draftVersionComments []DraftVersionComment
	for draftVersionId, userComments := range draftVersionVsComments {
		versionComment := DraftVersionComment{
			DraftVersionId: draftVersionId,
			UserComments:   userComments,
		}
		draftVersionComments = append(draftVersionComments, versionComment)
	}
	response := &DraftVersionCommentResponse{
		DraftId:              draftId,
		DraftVersionComments: draftVersionComments,
	}
	return response, nil
}

func (impl ConfigDraftServiceImpl) GetDrafts(appId int, envId int, resourceType DraftResourceType) ([]AppConfigDraft, error) {
	draftMetadataDtos, err := impl.configDraftRepository.GetDraftMetadata(appId, envId, resourceType)
	if err != nil {
		return nil, err
	}
	var appConfigDrafts []AppConfigDraft
	for _, draftMetadataDto := range draftMetadataDtos {
		appConfigDraft := draftMetadataDto.ConvertToAppConfigDraft()
		appConfigDrafts = append(appConfigDrafts, appConfigDraft)
	}
	return appConfigDrafts, nil
}

func (impl ConfigDraftServiceImpl) GetDraftById(draftId int) (*ConfigDraftResponse, error) {
	configDraft, err := impl.configDraftRepository.GetLatestConfigDraft(draftId)
	if err != nil {
		return nil, err
	}
	draftResponse := configDraft.ConvertToConfigDraft()
	return draftResponse, nil
}

func (impl ConfigDraftServiceImpl) GetDraftByDraftVersionId(draftVersionId int) (*ConfigDraftResponse, error) {
	draftVersionDto, err := impl.configDraftRepository.GetDraftVersionById(draftVersionId)
	if err != nil {
		return nil, err
	}
	return draftVersionDto.ConvertToConfigDraft(), nil
}

func (impl ConfigDraftServiceImpl) DeleteComment(draftId int, draftCommentId int, userId int32) error {
	deletedCount, err := impl.configDraftRepository.DeleteComment(draftId, draftCommentId, userId)
	if err != nil {
		return err
	}
	if deletedCount == 0 {
		return errors.New("failed to delete comment")
	}
	return nil
}

func (impl ConfigDraftServiceImpl) ApproveDraft(draftId int, draftVersionId int, userId int32) error {
	toUpdateDraftState := PublishedDraftState
	draftVersion, err := impl.validateDraftAction(draftId, draftVersionId, toUpdateDraftState)
	if err != nil {
		return err
	}
	// once it is approved, we need to create actual resources
	draftData := draftVersion.Data
	draftsDto := draftVersion.DraftsDto
	draftResourceType := draftsDto.Resource
	if draftResourceType == CmDraftResource || draftResourceType == CsDraftResource {
		err = impl.handleCmCsData(draftResourceType, draftsDto.AppId, draftsDto.EnvId, draftData, draftVersion.UserId)
	} else {
		err = impl.handleDeploymentTemplate(draftsDto.AppId, draftsDto.EnvId, draftData, draftVersion.UserId, draftVersion.Action)
	}
	if err != nil {
		return err
	}
	err = impl.configDraftRepository.UpdateDraftState(draftId, toUpdateDraftState, userId)
	return err
}

func (impl ConfigDraftServiceImpl) handleCmCsData(draftResource DraftResourceType, appId int, envId int, draftData string, userId int32) error {
	// if envId is -1 then it is base Configuration else Env level config
	var configDataRequest *pipeline.ConfigDataRequest
	err := json.Unmarshal([]byte(draftData), configDataRequest)
	if err != nil {
		impl.logger.Errorw("error occurred while unmarshalling draftData of CM/CS", "appId", appId, "envId", envId, "err", err)
		return err
	}
	configDataRequest.UserId = userId // setting draftVersion userId
	isCm := draftResource == CmDraftResource
	if isCm {
		if envId == -1 {
			_, err = impl.configMapService.CMGlobalAddUpdate(configDataRequest)
		} else {
			_, err = impl.configMapService.CMEnvironmentAddUpdate(configDataRequest)
		}
	} else {
		if envId == -1 {
			_, err = impl.configMapService.CSGlobalAddUpdate(configDataRequest)
		} else {
			_, err = impl.configMapService.CSEnvironmentAddUpdate(configDataRequest)
		}
	}
	if err != nil {
		impl.logger.Errorw("error occurred while adding/updating config", "isCm", isCm, "appId", appId, "envId", envId, "err", err)
	}
	return err
}

func (impl ConfigDraftServiceImpl) handleDeploymentTemplate(appId int, envId int, draftData string, userId int32, action ResourceAction) error {

	ctx := context.Background()
	var err error
	var templateValidated bool
	if envId == -1 {
		var templateRequest *chart.TemplateRequest
		err = json.Unmarshal([]byte(draftData), templateRequest)
		if err != nil {
			impl.logger.Errorw("error occurred while unmarshalling draftData of deployment template", "appId", appId, "envId", envId, "err", err)
			return err
		}
		templateValidated, err = impl.chartService.DeploymentTemplateValidate(ctx, draftData, templateRequest.ChartRefId)
		if err != nil {
			return err
		}
		if !templateValidated {
			return errors.New("template-outdated")
		}
		templateRequest.UserId = userId
		_, err = impl.chartService.UpdateAppOverride(ctx, templateRequest)
	} else {
		var envConfigProperties *pipeline.EnvironmentProperties
		err = json.Unmarshal([]byte(draftData), envConfigProperties)
		if err != nil {
			impl.logger.Errorw("error occurred while unmarshalling draftData of env deployment template", "appId", appId, "envId", envId, "err", err)
			return err
		}
		envConfigProperties.UserId = userId
		envConfigProperties.EnvironmentId = envId
		chartRefId := envConfigProperties.ChartRefId
		templateValidated, err = impl.chartService.DeploymentTemplateValidate(ctx, envConfigProperties.EnvOverrideValues, chartRefId)
		if err != nil {
			return err
		}
		if !templateValidated {
			return errors.New("template-outdated")
		}
		if action == AddResourceAction {
			_, err = impl.propertiesConfigService.CreateEnvironmentProperties(appId, envConfigProperties)
			if err != nil {
				if err.Error() == bean2.NOCHARTEXIST {
					appMetrics := false
					if envConfigProperties.AppMetrics != nil {
						appMetrics = *envConfigProperties.AppMetrics
					}
					templateRequest := chart.TemplateRequest{
						AppId:               appId,
						ChartRefId:          envConfigProperties.ChartRefId,
						ValuesOverride:      []byte("{}"),
						UserId:              userId,
						IsAppMetricsEnabled: appMetrics,
					}
					_, err = impl.chartService.CreateChartFromEnvOverride(templateRequest, ctx)
					if err != nil {
						impl.logger.Errorw("service err, EnvConfigOverrideCreate from draft", "appId", appId, "envId", envId, "err", err, "payload", envConfigProperties)
						return err
					}
					_, err = impl.propertiesConfigService.CreateEnvironmentProperties(appId, envConfigProperties)
					if err != nil {
						impl.logger.Errorw("service err, EnvConfigOverrideCreate", "appId", appId, "envId", envId, "err", err, "payload", envConfigProperties)
						return err
					}
				}
			} else {
				return err
			}
		} else {
			_, err = impl.propertiesConfigService.UpdateEnvironmentProperties(appId, envConfigProperties, userId)
			if err != nil {
				impl.logger.Errorw("service err, EnvConfigOverrideUpdate", "appId", appId, "envId", envId, "err", err, "payload", envConfigProperties)
				return err
			}
		}
	}
	return nil
}


