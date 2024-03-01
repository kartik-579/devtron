package policyGovernance

import (
	"github.com/devtron-labs/devtron/pkg/globalPolicy/bean"
	"github.com/devtron-labs/devtron/pkg/resourceQualifiers"
)

const NO_POLICY = "NA"

type PathVariablePolicyType string

const PathVariablePolicyTypeVariable string = "policyType"
const (
	ImagePromotion   PathVariablePolicyType = "artifact-promotion"
	DeploymentWindow PathVariablePolicyType = "deployment-window"
)

var ExistingPolicyTypes = []PathVariablePolicyType{ImagePromotion, DeploymentWindow}
var PathPolicyTypeToGlobalPolicyTypeMap = map[PathVariablePolicyType]bean.GlobalPolicyType{
	//ImagePromotion:   bean.GLOBAL_POLICY_TYPE_IMAGE_PROMOTION_POLICY,
	DeploymentWindow: bean.GLOBAL_POLICY_TYPE_DEPLOYMENT_WINDOW,
}

var GlobalPolicyTypeToResourceTypeMap = map[bean.GlobalPolicyType]resourceQualifiers.ResourceType{
	bean.GLOBAL_POLICY_TYPE_DEPLOYMENT_WINDOW: resourceQualifiers.DeploymentWindowProfile,
	// todo
	// bean.GLOBAL_POLICY_TYPE_DEPLOYMENT_WINDOW: resourceQualifiers.,
}

type AppEnvPolicyContainer struct {
	AppId      int    `json:"appId"`
	EnvId      int    `json:"envId"`
	PolicyId   int    `json:"policyId"`
	AppName    string `json:"appName"`
	EnvName    string `json:"envName"`
	PolicyName string `json:"policyName,omitempty"`
}

type AppEnvPolicyMappingsListFilter struct {
	PolicyType  bean.GlobalPolicyType `json:"-"`
	AppNames    []string              `json:"appNames"`
	EnvNames    []string              `json:"envNames"`
	PolicyNames []string              `json:"policyNames"`
	SortBy      string                `json:"sortBy,omitempty" validate:"omitempty,oneof=appName environmentName"`
	SortOrder   string                `json:"sortOrder,omitempty" validate:"omitempty,oneof=ASC DESC"`
	Offset      int                   `json:"offset,omitempty" validate:"omitempty,min=0"`
	Size        int                   `json:"size,omitempty" validate:"omitempty,min=0"`
}

type BulkPromotionPolicyApplyRequest struct {
	PolicyType              bean.GlobalPolicyType   `json:"-"`
	ApplicationEnvironments []AppEnvPolicyContainer `json:"applicationEnvironments"`
	ApplyToPolicyName       string                  `json:"applyToPolicyName"`
	ApplyToPolicyId         int                     `json:"applyToPolicyId"`
	//AppEnvPolicyListFilter  AppEnvPolicyMappingsListFilter `json:"appEnvPolicyListFilter" validate:"dive"`
}
