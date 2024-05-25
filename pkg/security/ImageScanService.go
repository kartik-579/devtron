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

package security

import (
	"context"
	"github.com/devtron-labs/devtron/pkg/cluster/repository/bean"
	bean2 "github.com/devtron-labs/devtron/pkg/deployment/trigger/devtronApps/bean"
	"go.opentelemetry.io/otel"
	"time"

	repository1 "github.com/devtron-labs/devtron/internal/sql/repository/app"
	"github.com/devtron-labs/devtron/internal/sql/repository/helper"
	"github.com/devtron-labs/devtron/pkg/auth/user"
	repository2 "github.com/devtron-labs/devtron/pkg/team"

	"github.com/devtron-labs/devtron/internal/sql/repository"
	"github.com/devtron-labs/devtron/internal/sql/repository/pipelineConfig"
	"github.com/devtron-labs/devtron/internal/sql/repository/security"
	"github.com/devtron-labs/devtron/pkg/cluster"
	"github.com/go-pg/pg"
	"go.uber.org/zap"
)

type ImageScanService interface {
	FetchAllDeployInfo(request *ImageScanRequest) ([]*security.ImageScanDeployInfo, error)
	FetchScanExecutionListing(request *ImageScanRequest, ids []int) (*ImageScanHistoryListingResponse, error)
	FetchExecutionDetailResult(request *ImageScanRequest) (*ImageScanExecutionDetail, error)
	FetchMinScanResultByAppIdAndEnvId(request *ImageScanRequest) (*ImageScanExecutionDetail, error)
	VulnerabilityExposure(request *security.VulnerabilityRequest) (*security.VulnerabilityExposureListingResponse, error)
	GetArtifactVulnerabilityStatus(ctx context.Context, request *bean2.VulnerabilityCheckRequest) (bool, error)
}

type ImageScanServiceImpl struct {
	Logger                                    *zap.SugaredLogger
	scanHistoryRepository                     security.ImageScanHistoryRepository
	scanResultRepository                      security.ImageScanResultRepository
	scanObjectMetaRepository                  security.ImageScanObjectMetaRepository
	cveStoreRepository                        security.CveStoreRepository
	imageScanDeployInfoRepository             security.ImageScanDeployInfoRepository
	userService                               user.UserService
	teamRepository                            repository2.TeamRepository
	appRepository                             repository1.AppRepository
	envService                                cluster.EnvironmentService
	ciArtifactRepository                      repository.CiArtifactRepository
	policyService                             PolicyService
	pipelineRepository                        pipelineConfig.PipelineRepository
	ciPipelineRepository                      pipelineConfig.CiPipelineRepository
	scanToolMetaDataRepository                security.ScanToolMetadataRepository
	scanToolExecutionHistoryMappingRepository security.ScanToolExecutionHistoryMappingRepository
	cvePolicyRepository                       security.CvePolicyRepository
}

type ImageScanRequest struct {
	ScanExecutionId       int    `json:"ScanExecutionId"`
	ImageScanDeployInfoId int    `json:"imageScanDeployInfo"`
	AppId                 int    `json:"appId"`
	EnvId                 int    `json:"envId"`
	ObjectId              int    `json:"objectId"`
	ArtifactId            int    `json:"artifactId"`
	Image                 string `json:"image"`
	security.ImageScanFilter
}

type ImageScanHistoryListingResponse struct {
	Offset                   int                         `json:"offset"`
	Size                     int                         `json:"size"`
	Total                    int                         `json:"total"`
	ImageScanHistoryResponse []*ImageScanHistoryResponse `json:"scanList"`
}

type ImageScanHistoryResponse struct {
	ImageScanDeployInfoId int            `json:"imageScanDeployInfoId"`
	AppId                 int            `json:"appId"`
	EnvId                 int            `json:"envId"`
	Name                  string         `json:"name"`
	Type                  string         `json:"type"`
	Environment           string         `json:"environment"`
	LastChecked           *time.Time     `json:"lastChecked"`
	Image                 string         `json:"image,omitempty"`
	SeverityCount         *SeverityCount `json:"severityCount,omitempty"`
}

type ImageScanExecutionDetail struct {
	ImageScanDeployInfoId int                `json:"imageScanDeployInfoId"`
	AppId                 int                `json:"appId,omitempty"`
	EnvId                 int                `json:"envId,omitempty"`
	AppName               string             `json:"appName,omitempty"`
	EnvName               string             `json:"envName,omitempty"`
	ArtifactId            int                `json:"artifactId,omitempty"`
	Image                 string             `json:"image,omitempty"`
	PodName               string             `json:"podName,omitempty"`
	ReplicaSet            string             `json:"replicaSet,omitempty"`
	Vulnerabilities       []*Vulnerabilities `json:"vulnerabilities,omitempty"`
	SeverityCount         *SeverityCount     `json:"severityCount,omitempty"`
	ExecutionTime         time.Time          `json:"executionTime,omitempty"`
	ScanEnabled           bool               `json:"scanEnabled,notnull"`
	Scanned               bool               `json:"scanned,notnull"`
	ObjectType            string             `json:"objectType,notnull"`
	ScanToolId            int                `json:"scanToolId,omitempty""`
}

type Vulnerabilities struct {
	CVEName    string `json:"cveName"`
	Severity   string `json:"severity"`
	Package    string `json:"package,omitempty"`
	CVersion   string `json:"currentVersion"`
	FVersion   string `json:"fixedVersion"`
	Permission string `json:"permission"`
}

type SeverityCount struct {
	High     int `json:"high"`
	Moderate int `json:"moderate"`
	Low      int `json:"low"`
}

func NewImageScanServiceImpl(Logger *zap.SugaredLogger, scanHistoryRepository security.ImageScanHistoryRepository,
	scanResultRepository security.ImageScanResultRepository, scanObjectMetaRepository security.ImageScanObjectMetaRepository,
	cveStoreRepository security.CveStoreRepository, imageScanDeployInfoRepository security.ImageScanDeployInfoRepository,
	userService user.UserService, teamRepository repository2.TeamRepository,
	appRepository repository1.AppRepository,
	envService cluster.EnvironmentService, ciArtifactRepository repository.CiArtifactRepository, policyService PolicyService,
	pipelineRepository pipelineConfig.PipelineRepository, ciPipelineRepository pipelineConfig.CiPipelineRepository, scanToolMetaDataRepository security.ScanToolMetadataRepository, scanToolExecutionHistoryMappingRepository security.ScanToolExecutionHistoryMappingRepository,
	cvePolicyRepository security.CvePolicyRepository) *ImageScanServiceImpl {
	return &ImageScanServiceImpl{Logger: Logger, scanHistoryRepository: scanHistoryRepository, scanResultRepository: scanResultRepository,
		scanObjectMetaRepository: scanObjectMetaRepository, cveStoreRepository: cveStoreRepository,
		imageScanDeployInfoRepository:             imageScanDeployInfoRepository,
		userService:                               userService,
		teamRepository:                            teamRepository,
		appRepository:                             appRepository,
		envService:                                envService,
		ciArtifactRepository:                      ciArtifactRepository,
		policyService:                             policyService,
		pipelineRepository:                        pipelineRepository,
		ciPipelineRepository:                      ciPipelineRepository,
		scanToolMetaDataRepository:                scanToolMetaDataRepository,
		scanToolExecutionHistoryMappingRepository: scanToolExecutionHistoryMappingRepository,
		cvePolicyRepository:                       cvePolicyRepository,
	}
}

func (impl ImageScanServiceImpl) FetchAllDeployInfo(request *ImageScanRequest) ([]*security.ImageScanDeployInfo, error) {
	deployedList, err := impl.imageScanDeployInfoRepository.FindAll()
	if err != nil {
		impl.Logger.Errorw("error while fetching scan execution result", "err", err)
		return nil, err
	}
	return deployedList, nil
}

func (impl ImageScanServiceImpl) FetchScanExecutionListing(request *ImageScanRequest, deployInfoIds []int) (*ImageScanHistoryListingResponse, error) {
	size := request.Size
	request.Size = 0
	groupByListCount, err := impl.imageScanDeployInfoRepository.ScanListingWithFilter(&request.ImageScanFilter, request.Size, request.Offset, deployInfoIds)
	if err != nil {
		impl.Logger.Errorw("error while fetching scan execution result", "err", err)
		return nil, err
	}
	request.Size = size
	groupByList, err := impl.imageScanDeployInfoRepository.ScanListingWithFilter(&request.ImageScanFilter, request.Size, request.Offset, deployInfoIds)
	if err != nil {
		impl.Logger.Errorw("error while fetching scan execution result", "err", err)
		return nil, err
	}
	var ids []int
	for _, item := range groupByList {
		ids = append(ids, item.Id)
	}
	if len(ids) == 0 {
		impl.Logger.Debugw("no image scan deploy info exists", "err", err)
		responseList := make([]*ImageScanHistoryResponse, 0)
		return &ImageScanHistoryListingResponse{ImageScanHistoryResponse: responseList}, nil
	}
	deployedList, err := impl.imageScanDeployInfoRepository.FindByIds(ids)
	if err != nil {
		impl.Logger.Errorw("error while fetching scan execution result", "err", err)
		return nil, err
	}

	groupByListMap := make(map[int]*security.ImageScanDeployInfo)
	for _, item := range deployedList {
		groupByListMap[item.Id] = item
	}

	var finalResponseList []*ImageScanHistoryResponse
	for _, item := range groupByList {
		imageScanHistoryResponse := &ImageScanHistoryResponse{}
		var lastChecked time.Time

		highCount := 0
		moderateCount := 0
		lowCount := 0
		imageScanDeployInfo := groupByListMap[item.Id]
		if imageScanDeployInfo != nil {
			scanResultList, err := impl.scanResultRepository.FetchByScanExecutionIds(imageScanDeployInfo.ImageScanExecutionHistoryId)
			if err != nil && err != pg.ErrNoRows {
				impl.Logger.Errorw("error while fetching scan execution result", "err", err)
				//return nil, err
			}
			if err == pg.ErrNoRows {
				impl.Logger.Errorw("no scan execution data found, but has image scan deployed info", "err", err)
				return nil, err
			}

			for _, item := range scanResultList {
				lastChecked = item.ImageScanExecutionHistory.ExecutionTime
				if item.CveStore.Severity == security.Critical {
					highCount = highCount + 1
				} else if item.CveStore.Severity == security.Medium {
					moderateCount = moderateCount + 1
				} else if item.CveStore.Severity == security.Low {
					lowCount = lowCount + 1
				}
			}
		}
		severityCount := &SeverityCount{
			High:     highCount,
			Moderate: moderateCount,
			Low:      lowCount,
		}
		imageScanHistoryResponse.ImageScanDeployInfoId = item.Id
		if imageScanDeployInfo != nil {
			imageScanHistoryResponse.LastChecked = &lastChecked
		}
		imageScanHistoryResponse.SeverityCount = severityCount
		if imageScanDeployInfo != nil {
			imageScanHistoryResponse.EnvId = imageScanDeployInfo.EnvId
		}
		imageScanHistoryResponse.Environment = item.EnvironmentName
		if len(request.AppName) > 0 || len(request.ObjectName) > 0 {
			imageScanHistoryResponse.Name = item.ObjectName
			imageScanHistoryResponse.Type = item.ObjectType
			if len(request.AppName) > 0 {
				imageScanHistoryResponse.AppId = item.ScanObjectMetaId
			}
		} else {
			if item.ObjectType == security.ScanObjectType_APP || item.ObjectType == security.ScanObjectType_CHART {
				app, err := impl.appRepository.FindById(item.ScanObjectMetaId)
				if err != nil && err != pg.ErrNoRows {
					return nil, err
				}
				if err == pg.ErrNoRows {
					continue
				}
				imageScanHistoryResponse.AppId = app.Id
				imageScanHistoryResponse.Name = app.AppName
				imageScanHistoryResponse.Type = item.ObjectType
			} else if item.ObjectType == security.ScanObjectType_POD {
				scanObjectMeta, err := impl.scanObjectMetaRepository.FindOne(item.ScanObjectMetaId)
				if err != nil && err != pg.ErrNoRows {
					return nil, err
				}
				if err == pg.ErrNoRows {
					continue
				}
				imageScanHistoryResponse.Name = scanObjectMeta.Name
				imageScanHistoryResponse.Type = item.ObjectType
			}
		}
		finalResponseList = append(finalResponseList, imageScanHistoryResponse)
	}

	finalResponse := &ImageScanHistoryListingResponse{
		Offset:                   request.Offset,
		Size:                     request.Size,
		ImageScanHistoryResponse: finalResponseList,
		Total:                    len(groupByListCount),
	}

	/*
		1) fetch from image_deployment_info, group by object id and type,
		2) on iteration collect image scan execution id and fetch its result and put in map
		3) merge 1st result with map
	*/

	return finalResponse, err
}

func (impl ImageScanServiceImpl) FetchExecutionDetailResult(request *ImageScanRequest) (*ImageScanExecutionDetail, error) {
	//var scanExecution *security.ImageScanExecutionHistory
	var scanExecutionIds []int
	var executionTime time.Time
	imageScanResponse := &ImageScanExecutionDetail{}
	isRegularApp := false
	if request.ImageScanDeployInfoId > 0 {
		// scan detail for deployed images
		scanDeployInfo, err := impl.imageScanDeployInfoRepository.FindOne(request.ImageScanDeployInfoId)
		if err != nil {
			impl.Logger.Errorw("error while fetching scan execution result", "err", err)
			return nil, err
		}

		scanExecutionIds = append(scanExecutionIds, scanDeployInfo.ImageScanExecutionHistoryId...)

		if scanDeployInfo.ObjectType == security.ScanObjectType_APP || scanDeployInfo.ObjectType == security.ScanObjectType_CHART {
			request.AppId = scanDeployInfo.ScanObjectMetaId
		} else if scanDeployInfo.ObjectType == security.ScanObjectType_POD {
			request.ObjectId = scanDeployInfo.ScanObjectMetaId
		}
		request.EnvId = scanDeployInfo.EnvId
		if scanDeployInfo.ObjectType == security.ScanObjectType_APP {
			isRegularApp = true
		}
		imageScanResponse.ObjectType = scanDeployInfo.ObjectType
		if !isRegularApp {
			imageScanResponse.ScanEnabled = true
			imageScanResponse.Scanned = true
		}
	} else if request.ArtifactId > 0 {
		// scan detail for artifact weather it is deployed or not, used here for ci build history
		ciArtifact, err := impl.ciArtifactRepository.Get(request.ArtifactId)
		if err != nil {
			impl.Logger.Errorw("error while fetching scan execution result", "err", err)
			return nil, err
		}
		scanExecution, err := impl.scanHistoryRepository.FindByImageAndDigest(ciArtifact.ImageDigest, ciArtifact.Image)
		if err != nil {
			impl.Logger.Errorw("error while fetching scan execution result", "err", err)
			return nil, err
		}
		ciPipeline, err := impl.ciPipelineRepository.FindByIdIncludingInActive(ciArtifact.PipelineId)
		if err != nil {
			impl.Logger.Errorw("error while fetching scan execution result", "err", err)
			return nil, err
		}
		imageScanResponse.AppId = ciPipeline.AppId

		scanExecutionIds = append(scanExecutionIds, scanExecution.Id)
		executionTime = scanExecution.ExecutionTime
		imageScanResponse.ScanEnabled = ciArtifact.ScanEnabled
		imageScanResponse.Scanned = ciArtifact.Scanned
		if ciArtifact.ScanEnabled == false {
			impl.Logger.Debugw("returning without result as scan disabled for this artifact", "ciArtifact", ciArtifact)
			return imageScanResponse, nil
		}
		imageScanResponse.ObjectType = security.ScanObjectType_APP
	}

	var vulnerabilities []*Vulnerabilities
	var highCount, moderateCount, lowCount int
	var cveStores []*security.CveStore
	imageDigests := make(map[string]string)
	if len(scanExecutionIds) > 0 {
		//var imageScanResultFinal []*security.ImageScanExecutionResult
		imageScanResult, err := impl.scanResultRepository.FetchByScanExecutionIds(scanExecutionIds)
		if err != nil {
			impl.Logger.Errorw("error while fetching scan execution result", "err", err)
			return nil, err
		}

		for _, item := range imageScanResult {
			vulnerability := &Vulnerabilities{
				CVEName:  item.CveStore.Name,
				CVersion: item.CveStore.Version,
				FVersion: item.CveStore.FixedVersion,
				Package:  item.CveStore.Package,
				Severity: item.CveStore.Severity.String(),
				//Permission: "BLOCK", TODO
			}
			if len(item.Package) > 0 {
				// data already migrated hence get package from image_scan_execution_result
				vulnerability.Package = item.Package
			}
			if item.CveStore.Severity == security.Critical {
				highCount = highCount + 1
			} else if item.CveStore.Severity == security.Medium {
				moderateCount = moderateCount + 1
			} else if item.CveStore.Severity == security.Low {
				lowCount = lowCount + 1
			}
			vulnerabilities = append(vulnerabilities, vulnerability)
			cveStores = append(cveStores, &item.CveStore)
			if _, ok := imageDigests[item.ImageScanExecutionHistory.ImageHash]; !ok {
				imageDigests[item.ImageScanExecutionHistory.ImageHash] = item.ImageScanExecutionHistory.ImageHash
			}
			executionTime = item.ImageScanExecutionHistory.ExecutionTime
		}
		if len(imageScanResult) > 0 {
			imageScanResponse.ScanToolId = imageScanResult[0].ScanToolId
		} else {
			toolIdFromExecutionHistory, err := impl.getScanToolIdFromExecutionHistory(scanExecutionIds)
			if err != nil || toolIdFromExecutionHistory == -1 {
				impl.Logger.Errorw("error in getting scan tool id from exection history", "err", err, "")
				return nil, err
			}
			imageScanResponse.ScanToolId = toolIdFromExecutionHistory
		}
	}
	severityCount := &SeverityCount{
		High:     highCount,
		Moderate: moderateCount,
		Low:      lowCount,
	}
	imageScanResponse.ImageScanDeployInfoId = request.ImageScanDeployInfoId
	if len(vulnerabilities) == 0 {
		vulnerabilities = make([]*Vulnerabilities, 0)
	}
	imageScanResponse.Vulnerabilities = vulnerabilities
	imageScanResponse.SeverityCount = severityCount
	imageScanResponse.ExecutionTime = executionTime

	// scanned enabled or not only for when we don't ask for direct artifact id
	if request.ImageScanDeployInfoId > 0 && request.ArtifactId == 0 && isRegularApp {
		for _, v := range imageDigests {
			ciArtifact, err := impl.ciArtifactRepository.GetByImageDigest(v)
			if err != nil {
				impl.Logger.Errorw("error while fetching scan execution result", "err", err)
				return nil, err
			}
			imageScanResponse.ScanEnabled = ciArtifact.ScanEnabled
			imageScanResponse.Scanned = ciArtifact.Scanned
		}
	}

	if request.AppId > 0 && request.EnvId > 0 {
		app, err := impl.appRepository.FindById(request.AppId)
		if err != nil {
			impl.Logger.Errorw("error while fetching env", "err", err)
			return nil, err
		}
		imageScanResponse.AppId = request.AppId
		imageScanResponse.AppName = app.AppName
		env, err := impl.envService.FindById(request.EnvId)
		if err != nil {
			impl.Logger.Errorw("error while fetching env", "err", err)
			return nil, err
		}
		imageScanResponse.EnvId = request.EnvId
		imageScanResponse.EnvName = env.Environment

		blockCveList, err := impl.policyService.GetBlockedCVEList(cveStores, env.ClusterId, env.Id, request.AppId, app.AppType == helper.ChartStoreApp)
		if err != nil {
			impl.Logger.Errorw("error while fetching env", "err", err)
			//return nil, err
			//TODO - review @nishant
		}
		if blockCveList != nil {
			vulnerabilityPermissionMap := make(map[string]string)
			for _, cve := range blockCveList {
				vulnerabilityPermissionMap[cve.Name] = "BLOCK"
			}
			var updatedVulnerabilities []*Vulnerabilities
			for _, vulnerability := range imageScanResponse.Vulnerabilities {
				if _, ok := vulnerabilityPermissionMap[vulnerability.CVEName]; ok {
					vulnerability.Permission = "BLOCK"
				} else {
					vulnerability.Permission = "WHITELISTED"
				}
				updatedVulnerabilities = append(updatedVulnerabilities, vulnerability)
			}
			if len(updatedVulnerabilities) == 0 {
				updatedVulnerabilities = make([]*Vulnerabilities, 0)
			}
			imageScanResponse.Vulnerabilities = updatedVulnerabilities
		} else {
			for _, vulnerability := range imageScanResponse.Vulnerabilities {
				vulnerability.Permission = "WHITELISTED"
			}
		}
	}
	return imageScanResponse, nil
}

func (impl ImageScanServiceImpl) FetchMinScanResultByAppIdAndEnvId(request *ImageScanRequest) (*ImageScanExecutionDetail, error) {
	//var scanExecution *security.ImageScanExecutionHistory
	var scanExecutionIds []int
	var executionTime time.Time

	var objectType []string
	objectType = append(objectType, security.ScanObjectType_APP, security.ScanObjectType_CHART)
	scanDeployInfo, err := impl.imageScanDeployInfoRepository.FetchByAppIdAndEnvId(request.AppId, request.EnvId, objectType)
	if err != nil && pg.ErrNoRows != err {
		impl.Logger.Errorw("error while fetching scan execution result", "err", err)
		return nil, err
	}
	if scanDeployInfo == nil || scanDeployInfo.Id == 0 || err == pg.ErrNoRows {
		return nil, err
	}
	scanExecutionIds = append(scanExecutionIds, scanDeployInfo.ImageScanExecutionHistoryId...)

	var highCount, moderateCount, lowCount, scantoolId int
	if len(scanExecutionIds) > 0 {
		imageScanResult, err := impl.scanResultRepository.FetchByScanExecutionIds(scanExecutionIds)
		if err != nil {
			impl.Logger.Errorw("error while fetching scan execution result", "err", err)
			return nil, err
		}
		for _, item := range imageScanResult {
			executionTime = item.ImageScanExecutionHistory.ExecutionTime
			if item.CveStore.Severity == security.Critical {
				highCount = highCount + 1
			} else if item.CveStore.Severity == security.Medium {
				moderateCount = moderateCount + 1
			} else if item.CveStore.Severity == security.Low {
				lowCount = lowCount + 1
			}
		}
		if len(imageScanResult) > 0 {
			scantoolId = imageScanResult[0].ScanToolId
		} else {
			toolIdFromExecutionHistory, err := impl.getScanToolIdFromExecutionHistory(scanExecutionIds)
			if err != nil || toolIdFromExecutionHistory == -1 {
				impl.Logger.Errorw("error in getting scan tool id from exection history", "err", err, "")
				return nil, err
			}
			scantoolId = toolIdFromExecutionHistory
		}
	}
	severityCount := &SeverityCount{
		High:     highCount,
		Moderate: moderateCount,
		Low:      lowCount,
	}
	imageScanResponse := &ImageScanExecutionDetail{
		ImageScanDeployInfoId: scanDeployInfo.Id,
		SeverityCount:         severityCount,
		ExecutionTime:         executionTime,
		ObjectType:            scanDeployInfo.ObjectType,
		ScanEnabled:           true,
		Scanned:               true,
		ScanToolId:            scantoolId,
	}
	return imageScanResponse, nil
}
func (impl ImageScanServiceImpl) getScanToolIdFromExecutionHistory(scanExecutionIds []int) (int, error) {
	scanToolHistoryMappings, err := impl.scanToolExecutionHistoryMappingRepository.GetAllScanHistoriesByExecutionHistoryIds(scanExecutionIds)
	if err != nil {
		if err == pg.ErrNoRows {
			impl.Logger.Errorw("got no rows for scanToolHistoryMappings", "err", err)
		} else {
			impl.Logger.Errorw("error in getting scanToolHistoryMappings", "err", err)
			return -1, err
		}
	}
	if len(scanToolHistoryMappings) > 0 {
		return scanToolHistoryMappings[0].ScanToolId, nil
	}
	return -1, err
}

func (impl ImageScanServiceImpl) VulnerabilityExposure(request *security.VulnerabilityRequest) (*security.VulnerabilityExposureListingResponse, error) {
	vulnerabilityExposureListingResponse := &security.VulnerabilityExposureListingResponse{
		Offset: request.Offset,
		Size:   request.Size,
	}
	size := request.Size
	request.Size = 0
	count, err := impl.cveStoreRepository.VulnerabilityExposure(request)
	if err != nil {
		impl.Logger.Errorw("error while fetching vulnerability exposure", "err", err)
		return nil, err
	}
	request.Size = size
	vulnerabilityExposureListingResponse.Total = len(count)
	vulnerabilityExposureList, err := impl.cveStoreRepository.VulnerabilityExposure(request)
	if err != nil {
		impl.Logger.Errorw("error while fetching vulnerability exposure", "err", err)
		return nil, err
	}

	var cveStores []*security.CveStore
	cveStore, err := impl.cveStoreRepository.FindByName(request.CveName)
	if err != nil {
		impl.Logger.Errorw("error while fetching cve store", "err", err)
		return nil, err
	}

	envMap := make(map[int]bean.EnvironmentBean)
	environments, err := impl.envService.GetAllActive()
	if err != nil {
		impl.Logger.Errorw("error while fetching vulnerability exposure", "err", err)
		return nil, err
	}
	for _, item := range environments {
		envMap[item.Id] = item
	}

	cveStores = append(cveStores, cveStore)
	for _, item := range vulnerabilityExposureList {
		envId := 0
		if item.AppType == helper.ChartStoreApp {
			envId = item.ChartEnvId
		} else if item.AppType == helper.CustomApp {
			envId = item.PipelineEnvId
		}
		env := envMap[envId]
		item.EnvId = envId
		item.EnvName = env.Environment
		var appStore bool
		appStore = item.AppType == helper.ChartStoreApp
		blockCveList, err := impl.policyService.GetBlockedCVEList(cveStores, env.ClusterId, envId, item.AppId, appStore)
		if err != nil {
			impl.Logger.Errorw("error while fetching blocked list", "err", err)
			return nil, err
		}
		if len(blockCveList) > 0 {
			item.Blocked = true
		}
	}
	vulnerabilityExposureListingResponse.VulnerabilityExposure = vulnerabilityExposureList
	return vulnerabilityExposureListingResponse, nil
}

func (impl ImageScanServiceImpl) GetArtifactVulnerabilityStatus(ctx context.Context, request *bean2.VulnerabilityCheckRequest) (bool, error) {
	isVulnerable := false
	if len(request.ImageDigest) > 0 {
		var cveStores []*security.CveStore
		_, span := otel.Tracer("orchestrator").Start(ctx, "scanResultRepository.FindByImageDigest")
		imageScanResult, err := impl.scanResultRepository.FindByImageDigest(request.ImageDigest)
		span.End()
		if err != nil && err != pg.ErrNoRows {
			impl.Logger.Errorw("error fetching image digest", "digest", request.ImageDigest, "err", err)
			return false, err
		}
		for _, item := range imageScanResult {
			cveStores = append(cveStores, &item.CveStore)
		}
		_, span = otel.Tracer("orchestrator").Start(ctx, "cvePolicyRepository.GetBlockedCVEList")
		if request.CdPipeline.Environment.ClusterId == 0 {
			envDetails, err := impl.envService.GetDetailsById(request.CdPipeline.EnvironmentId)
			if err != nil {
				impl.Logger.Errorw("error fetching cluster details by env, GetArtifactVulnerabilityStatus", "envId", request.CdPipeline.EnvironmentId, "err", err)
				return false, err
			}
			request.CdPipeline.Environment = *envDetails
		}
		blockCveList, err := impl.cvePolicyRepository.GetBlockedCVEList(cveStores, request.CdPipeline.Environment.ClusterId, request.CdPipeline.EnvironmentId, request.CdPipeline.AppId, false)
		span.End()
		if err != nil {
			impl.Logger.Errorw("error encountered in GetArtifactVulnerabilityStatus", "clusterId", request.CdPipeline.Environment.ClusterId, "envId", request.CdPipeline.EnvironmentId, "appId", request.CdPipeline.AppId, "err", err)
			return false, err
		}
		if len(blockCveList) > 0 {
			isVulnerable = true
		}
	}
	return isVulnerable, nil
}
