/*
 * Copyright (c) 2024. Devtron Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package argoApplication

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/devtron-labs/common-lib/utils/k8s"
	k8sCommonBean "github.com/devtron-labs/common-lib/utils/k8s/commonBean"
	"github.com/devtron-labs/devtron/api/helm-app/gRPC"
	openapi "github.com/devtron-labs/devtron/api/helm-app/openapiClient"
	"github.com/devtron-labs/devtron/api/helm-app/service"
	"github.com/devtron-labs/devtron/pkg/argoApplication/bean"
	"github.com/devtron-labs/devtron/pkg/argoApplication/helper"
	"github.com/devtron-labs/devtron/pkg/argoApplication/read"
	cluster2 "github.com/devtron-labs/devtron/pkg/cluster"
	clusterRepository "github.com/devtron-labs/devtron/pkg/cluster/repository"
	k8s2 "github.com/devtron-labs/devtron/pkg/k8s"
	"github.com/devtron-labs/devtron/pkg/k8s/application"
	"github.com/devtron-labs/devtron/util/argo"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type ArgoApplicationService interface {
	ListApplications(clusterIds []int) ([]*bean.ArgoApplicationListDto, error)
	GetAppDetail(resourceName, resourceNamespace string, clusterId int) (*bean.ArgoApplicationDetailDto, error)
	HibernateArgoApplication(ctx context.Context, app *bean.ArgoAppIdentifier, hibernateRequest *openapi.HibernateRequest) ([]*openapi.HibernateStatus, error)
	UnHibernateArgoApplication(ctx context.Context, app *bean.ArgoAppIdentifier, hibernateRequest *openapi.HibernateRequest) ([]*openapi.HibernateStatus, error)
}

type ArgoApplicationServiceImpl struct {
	logger                *zap.SugaredLogger
	clusterRepository     clusterRepository.ClusterRepository
	k8sUtil               *k8s.K8sServiceImpl
	argoUserService       argo.ArgoUserService
	helmAppClient         gRPC.HelmAppClient
	helmAppService        service.HelmAppService
	k8sApplicationService application.K8sApplicationService
	readService           read.ArgoApplicationReadService
}

func NewArgoApplicationServiceImpl(logger *zap.SugaredLogger,
	clusterRepository clusterRepository.ClusterRepository,
	k8sUtil *k8s.K8sServiceImpl,
	argoUserService argo.ArgoUserService, helmAppClient gRPC.HelmAppClient,
	helmAppService service.HelmAppService,
	k8sApplicationService application.K8sApplicationService,
	readService read.ArgoApplicationReadService) *ArgoApplicationServiceImpl {
	return &ArgoApplicationServiceImpl{
		logger:                logger,
		clusterRepository:     clusterRepository,
		k8sUtil:               k8sUtil,
		argoUserService:       argoUserService,
		helmAppService:        helmAppService,
		helmAppClient:         helmAppClient,
		k8sApplicationService: k8sApplicationService,
		readService:           readService,
	}

}

func (impl *ArgoApplicationServiceImpl) ListApplications(clusterIds []int) ([]*bean.ArgoApplicationListDto, error) {
	var clusters []clusterRepository.Cluster
	var err error
	if len(clusterIds) > 0 {
		// getting cluster details by ids
		clusters, err = impl.clusterRepository.FindByIds(clusterIds)
		if err != nil {
			impl.logger.Errorw("error in getting clusters by ids", "err", err, "clusterIds", clusterIds)
			return nil, err
		}
	} else {
		clusters, err = impl.clusterRepository.FindAllActive()
		if err != nil {
			impl.logger.Errorw("error in getting all active clusters", "err", err)
			return nil, err
		}
	}

	listReq := &k8s2.ResourceRequestBean{
		K8sRequest: &k8s.K8sRequestBean{
			ResourceIdentifier: k8s.ResourceIdentifier{
				Namespace:        bean.AllNamespaces,
				GroupVersionKind: bean.GvkForArgoApplication,
			},
		},
	}
	// TODO: make goroutine and channel for optimization
	appListFinal := make([]*bean.ArgoApplicationListDto, 0)
	for _, cluster := range clusters {
		clusterObj := cluster
		if clusterObj.IsVirtualCluster || len(clusterObj.ErrorInConnecting) != 0 {
			continue
		}
		clusterBean := cluster2.GetClusterBean(clusterObj)
		clusterConfig := clusterBean.GetClusterConfig()
		restConfig, err := impl.k8sUtil.GetRestConfigByCluster(clusterConfig)
		if err != nil {
			impl.logger.Errorw("error in getting rest config by cluster Id", "err", err, "clusterId", clusterObj.Id)
			return nil, err
		}
		resp, err := impl.k8sApplicationService.GetResourceListWithRestConfig(context.Background(), "", listReq, nil, restConfig, clusterObj.ClusterName)
		if err != nil {
			if errStatus, ok := err.(*errors.StatusError); ok {
				if errStatus.Status().Code == 404 {
					// no argo apps found, not sending error
					impl.logger.Warnw("error in getting external argo app list, no apps found", "err", err, "clusterId", clusterObj.Id)
					continue
				}
			}
			impl.logger.Errorw("error in getting resource list", "err", err)
			return nil, err
		}
		appLists := getApplicationListDtos(resp, clusterObj.ClusterName, clusterObj.Id)
		appListFinal = append(appListFinal, appLists...)
	}
	return appListFinal, nil
}

func (impl *ArgoApplicationServiceImpl) GetAppDetail(resourceName, resourceNamespace string, clusterId int) (*bean.ArgoApplicationDetailDto, error) {
	appDetail := &bean.ArgoApplicationDetailDto{
		ArgoApplicationListDto: &bean.ArgoApplicationListDto{
			Name:      resourceName,
			Namespace: resourceNamespace,
			ClusterId: clusterId,
		},
	}
	clusters, err := impl.clusterRepository.FindAllActive()
	if err != nil {
		impl.logger.Errorw("error in getting all active clusters", "err", err)
		return nil, err
	}
	var clusterWithApplicationObject clusterRepository.Cluster
	clusterServerUrlIdMap := make(map[string]int, len(clusters))
	for _, cluster := range clusters {
		if cluster.Id == clusterId {
			clusterWithApplicationObject = cluster
		}
		clusterServerUrlIdMap[cluster.ServerUrl] = cluster.Id
	}
	if clusterWithApplicationObject.Id > 0 {
		appDetail.ClusterName = clusterWithApplicationObject.ClusterName
	}
	if clusterWithApplicationObject.IsVirtualCluster {
		return appDetail, nil
	} else if len(clusterWithApplicationObject.ErrorInConnecting) != 0 {
		return nil, fmt.Errorf("error in connecting to cluster")
	}
	clusterBean := cluster2.GetClusterBean(clusterWithApplicationObject)
	clusterConfig := clusterBean.GetClusterConfig()
	restConfig, err := impl.k8sUtil.GetRestConfigByCluster(clusterConfig)
	if err != nil {
		impl.logger.Errorw("error in getting rest config by cluster Id", "err", err, "clusterId", clusterWithApplicationObject.Id)
		return nil, err
	}
	resp, err := impl.k8sUtil.GetResource(context.Background(), resourceNamespace, resourceName, bean.GvkForArgoApplication, restConfig)
	if err != nil {
		impl.logger.Errorw("error in getting resource list", "err", err)
		return nil, err
	}
	var destinationServer string
	var argoManagedResources []*bean.ArgoManagedResource
	if resp != nil && resp.Manifest.Object != nil {
		appDetail.Manifest = resp.Manifest.Object
		appDetail.HealthStatus, appDetail.SyncStatus, destinationServer, argoManagedResources =
			getHealthSyncStatusDestinationServerAndManagedResourcesForArgoK8sRawObject(resp.Manifest.Object)
	}
	appDeployedOnClusterId := 0
	if destinationServer == k8s.DefaultClusterUrl {
		appDeployedOnClusterId = clusterWithApplicationObject.Id
	} else if clusterIdFromMap, ok := clusterServerUrlIdMap[destinationServer]; ok {
		appDeployedOnClusterId = clusterIdFromMap
	}
	var configOfClusterWhereAppIsDeployed bean.ArgoClusterConfigObj
	if appDeployedOnClusterId < 1 {
		// cluster is not added on devtron, need to get server config from secret which argo-cd saved
		coreV1Client, err := impl.k8sUtil.GetCoreV1ClientByRestConfig(restConfig)
		secrets, err := coreV1Client.Secrets(bean.AllNamespaces).List(context.Background(), v1.ListOptions{
			LabelSelector: labels.SelectorFromSet(labels.Set{"argocd.argoproj.io/secret-type": "cluster"}).String(),
		})
		if err != nil {
			impl.logger.Errorw("error in getting resource list, secrets", "err", err)
			return nil, err
		}
		for _, secret := range secrets.Items {
			if secret.Data != nil {
				if val, ok := secret.Data["server"]; ok {
					if string(val) == destinationServer {
						if config, ok := secret.Data["config"]; ok {
							err = json.Unmarshal(config, &configOfClusterWhereAppIsDeployed)
							if err != nil {
								impl.logger.Errorw("error in unmarshaling", "err", err)
								return nil, err
							}
							break
						}
					}
				}
			}
		}
	}
	resourceTreeResp, err := impl.getResourceTreeForExternalCluster(appDeployedOnClusterId, destinationServer, configOfClusterWhereAppIsDeployed, argoManagedResources)
	if err != nil {
		impl.logger.Errorw("error in getting resource tree response", "err", err)
		return nil, err
	}
	appDetail.ResourceTree = resourceTreeResp
	return appDetail, nil
}

func (impl *ArgoApplicationServiceImpl) getResourceTreeForExternalCluster(clusterId int, destinationServer string,
	configOfClusterWhereAppIsDeployed bean.ArgoClusterConfigObj, argoManagedResources []*bean.ArgoManagedResource) (*gRPC.ResourceTreeResponse, error) {
	var resources []*gRPC.ExternalResourceDetail
	for _, argoManagedResource := range argoManagedResources {
		resources = append(resources, &gRPC.ExternalResourceDetail{
			Group:     argoManagedResource.Group,
			Kind:      argoManagedResource.Kind,
			Version:   argoManagedResource.Version,
			Name:      argoManagedResource.Name,
			Namespace: argoManagedResource.Namespace,
		})
	}
	var clusterConfigOfClusterWhereAppIsDeployed *gRPC.ClusterConfig
	if len(configOfClusterWhereAppIsDeployed.BearerToken) > 0 {
		clusterConfigOfClusterWhereAppIsDeployed = &gRPC.ClusterConfig{
			ApiServerUrl:          destinationServer,
			Token:                 configOfClusterWhereAppIsDeployed.BearerToken,
			InsecureSkipTLSVerify: configOfClusterWhereAppIsDeployed.TlsClientConfig.Insecure,
			KeyData:               configOfClusterWhereAppIsDeployed.TlsClientConfig.KeyData,
			CaData:                configOfClusterWhereAppIsDeployed.TlsClientConfig.CaData,
			CertData:              configOfClusterWhereAppIsDeployed.TlsClientConfig.CertData,
		}
	}
	resourceTreeResp, err := impl.helmAppService.GetResourceTreeForExternalResources(context.Background(), clusterId, clusterConfigOfClusterWhereAppIsDeployed, resources)
	if err != nil {
		impl.logger.Errorw("error in getting resource tree for external resources", "err", err)
		return nil, err
	}
	return resourceTreeResp, nil
}

func getApplicationListDtos(resp *k8s.ClusterResourceListMap, clusterName string, clusterId int) []*bean.ArgoApplicationListDto {
	appLists := make([]*bean.ArgoApplicationListDto, 0)
	if resp != nil {
		appLists = make([]*bean.ArgoApplicationListDto, len(resp.Data))
		for i, rowData := range resp.Data {
			if rowData == nil {
				continue
			}
			appListDto := &bean.ArgoApplicationListDto{
				ClusterId:    clusterId,
				ClusterName:  clusterName,
				Name:         rowData[k8sCommonBean.K8sResourceColumnDefinitionName].(string),
				SyncStatus:   rowData[k8sCommonBean.K8sResourceColumnDefinitionSyncStatus].(string),
				HealthStatus: rowData[k8sCommonBean.K8sResourceColumnDefinitionHealthStatus].(string),
				Namespace:    rowData[k8sCommonBean.K8sClusterResourceNamespaceKey].(string),
			}
			appLists[i] = appListDto
		}
	}
	return appLists
}

func getHealthSyncStatusDestinationServerAndManagedResourcesForArgoK8sRawObject(obj map[string]interface{}) (string,
	string, string, []*bean.ArgoManagedResource) {
	var healthStatus, syncStatus, destinationServer string
	argoManagedResources := make([]*bean.ArgoManagedResource, 0)
	if specObjRaw, ok := obj[k8sCommonBean.Spec]; ok {
		specObj := specObjRaw.(map[string]interface{})
		if destinationObjRaw, ok2 := specObj[bean.Destination]; ok2 {
			destinationObj := destinationObjRaw.(map[string]interface{})
			if destinationServerIf, ok3 := destinationObj[bean.Server]; ok3 {
				destinationServer = destinationServerIf.(string)
			}
		}
	}
	if statusObjRaw, ok := obj[k8sCommonBean.K8sClusterResourceStatusKey]; ok {
		statusObj := statusObjRaw.(map[string]interface{})
		if healthObjRaw, ok2 := statusObj[k8sCommonBean.K8sClusterResourceHealthKey]; ok2 {
			healthObj := healthObjRaw.(map[string]interface{})
			if healthStatusIf, ok3 := healthObj[k8sCommonBean.K8sClusterResourceStatusKey]; ok3 {
				healthStatus = healthStatusIf.(string)
			}
		}
		if syncObjRaw, ok2 := statusObj[k8sCommonBean.K8sClusterResourceSyncKey]; ok2 {
			syncObj := syncObjRaw.(map[string]interface{})
			if syncStatusIf, ok3 := syncObj[k8sCommonBean.K8sClusterResourceStatusKey]; ok3 {
				syncStatus = syncStatusIf.(string)
			}
		}
		if resourceObjsRaw, ok2 := statusObj[k8sCommonBean.K8sClusterResourceResourcesKey]; ok2 {
			resourceObjs := resourceObjsRaw.([]interface{})
			argoManagedResources = make([]*bean.ArgoManagedResource, 0, len(resourceObjs))
			for _, resourceObjRaw := range resourceObjs {
				argoManagedResource := &bean.ArgoManagedResource{}
				resourceObj := resourceObjRaw.(map[string]interface{})
				if groupRaw, ok := resourceObj[k8sCommonBean.K8sClusterResourceGroupKey]; ok {
					argoManagedResource.Group = groupRaw.(string)
				}
				if kindRaw, ok := resourceObj[k8sCommonBean.K8sClusterResourceKindKey]; ok {
					argoManagedResource.Kind = kindRaw.(string)
				}
				if versionRaw, ok := resourceObj[k8sCommonBean.K8sClusterResourceVersionKey]; ok {
					argoManagedResource.Version = versionRaw.(string)
				}
				if nameRaw, ok := resourceObj[k8sCommonBean.K8sClusterResourceMetadataNameKey]; ok {
					argoManagedResource.Name = nameRaw.(string)
				}
				if namespaceRaw, ok := resourceObj[k8sCommonBean.K8sClusterResourceNamespaceKey]; ok {
					argoManagedResource.Namespace = namespaceRaw.(string)
				}
				argoManagedResources = append(argoManagedResources, argoManagedResource)
			}
		}
	}
	return healthStatus, syncStatus, destinationServer, argoManagedResources
}

func (impl *ArgoApplicationServiceImpl) HibernateArgoApplication(ctx context.Context, app *bean.ArgoAppIdentifier, hibernateRequest *openapi.HibernateRequest) ([]*openapi.HibernateStatus, error) {
	_, clusterBean, _, err := impl.readService.GetClusterConfigFromAllClusters(app.ClusterId)
	if err != nil {
		impl.logger.Errorw("HibernateArgoApplication", "error in getting the cluster config", err, "clusterId", app.ClusterId, "appName", app.AppName)
		return nil, err
	}
	conf := helper.ConvertClusterBeanToGrpcConfig(clusterBean)

	req := service.HibernateReqAdaptor(hibernateRequest)
	req.ClusterConfig = conf
	res, err := impl.helmAppClient.Hibernate(ctx, req)
	if err != nil {
		impl.logger.Errorw("HibernateArgoApplication", "error in hibernating the requested resource", err, "clusterId", app.ClusterId, "appName", app.AppName)
		return nil, err
	}
	response := service.HibernateResponseAdaptor(res.Status)
	return response, nil
}

func (impl *ArgoApplicationServiceImpl) UnHibernateArgoApplication(ctx context.Context, app *bean.ArgoAppIdentifier, hibernateRequest *openapi.HibernateRequest) ([]*openapi.HibernateStatus, error) {
	_, clusterBean, _, err := impl.readService.GetClusterConfigFromAllClusters(app.ClusterId)
	if err != nil {
		impl.logger.Errorw("HibernateArgoApplication", "error in getting the cluster config", err, "clusterId", app.ClusterId, "appName", app.AppName)
		return nil, err
	}
	conf := helper.ConvertClusterBeanToGrpcConfig(clusterBean)

	req := service.HibernateReqAdaptor(hibernateRequest)
	req.ClusterConfig = conf
	res, err := impl.helmAppClient.UnHibernate(ctx, req)
	if err != nil {
		impl.logger.Errorw("UnHibernateArgoApplication", "error in unHibernating the requested resources", err, "clusterId", app.ClusterId, "appName", app.AppName)
		return nil, err
	}
	response := service.HibernateResponseAdaptor(res.Status)
	return response, nil
}
