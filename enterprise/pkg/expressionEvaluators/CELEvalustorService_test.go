package expressionEvaluators

import (
	"fmt"
	"github.com/devtron-labs/devtron/internal/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEvaluateCELRequest(t *testing.T) {
	logger, _ := util.NewSugardLogger()
	celService := NewCELServiceImpl(logger)
	t.Run("diff output", func(tt *testing.T) {
		labelMapValue := map[string]string{}
		labelMapValue["devtron.ai/language"] = "python"
		params := []ExpressionParam{
			//{
			//	ParamName: ContainerRepo,
			//	Value:     "helloRepo",
			//	Type:      ParamTypeString,
			//},
			{
				ParamName: AppLabels,
				Value:     labelMapValue,
				Type:      ParamTypeStringMap,
			},
		}
		evalReq := CELRequest{
			Expression: "appLabels['devtron.ai/language'] == 'python' ? 'python_pvc' : 'java_pvc'",
			ExpressionMetadata: ExpressionMetadata{
				Params: params,
			},
		}
		evaluateCELRequest1, err := celService.EvaluateCELForObject(evalReq)
		assert.Nil(tt, err)
		fmt.Println("", evaluateCELRequest1)
	})
	//t.Run("valid release tags list", func(tt *testing.T) {
	//	artifact := "devtron/test:v1beta1"
	//	releaseTags := []string{"tag1", "latest"}
	//	params, err := resourceFilter.GetParamsFromArtifact(artifact, releaseTags, nil)
	//	assert.Nil(tt, err)
	//	evalReq := CELRequest{
	//		Expression: "'latest' in imageLabels",
	//		ExpressionMetadata: ExpressionMetadata{
	//			Params: params,
	//		},
	//	}
	//	res, err := celService.EvaluateCELForBool(evalReq)
	//	assert.Equal(t, true, res)
	//	assert.Equal(t, nil, err)
	//})
	//
	//t.Run("empty release tags list", func(tt *testing.T) {
	//	artifact := "devtron/test:v1beta1"
	//	releaseTags := []string{}
	//	params, err := resourceFilter.GetParamsFromArtifact(artifact, releaseTags, nil)
	//	assert.Nil(tt, err)
	//	evalReq := CELRequest{
	//		Expression: "'latest' in imageLabels",
	//		ExpressionMetadata: ExpressionMetadata{
	//			Params: params,
	//		},
	//	}
	//	res, err := celService.EvaluateCELForBool(evalReq)
	//	assert.Equal(t, false, res)
	//	assert.Equal(t, nil, err)
	//})
	//
	//t.Run("nil release tags list", func(tt *testing.T) {
	//	artifact := "devtron/test:v1beta1"
	//	var releaseTags []string
	//	params, err := resourceFilter.GetParamsFromArtifact(artifact, releaseTags, nil)
	//	assert.Nil(tt, err)
	//	evalReq := CELRequest{
	//		Expression: "'latest' in imageLabels",
	//		ExpressionMetadata: ExpressionMetadata{
	//			Params: params,
	//		},
	//	}
	//	res, err := celService.EvaluateCELForBool(evalReq)
	//	assert.Equal(t, false, res)
	//	assert.Equal(t, nil, err)
	//})
	//
	//t.Run("test commitDetails", func(tt *testing.T) {
	//	artifact := "devtron/test:v1beta1"
	//	var releaseTags []string
	//	minfo := repository.CiMaterialInfo{
	//		Material: repository.Material{
	//			Type:             "git",
	//			GitConfiguration: repository.GitConfiguration{URL: "github.com/test"},
	//		},
	//		Modifications: []repository.Modification{
	//			{
	//				Message: "test commit",
	//				Branch:  "test",
	//			},
	//		},
	//	}
	//	params, err := resourceFilter.GetParamsFromArtifact(artifact, releaseTags, []repository.CiMaterialInfo{minfo})
	//	assert.Nil(tt, err)
	//
	//	evalReq := CELRequest{
	//		Expression: "gitCommitDetails['github.com/test'].branch == 'test'",
	//		ExpressionMetadata: ExpressionMetadata{
	//			Params: params,
	//		},
	//	}
	//	b, err := celService.EvaluateCELForBool(evalReq)
	//	fmt.Println(b, err)
	//	assert.Equal(t, nil, err)
	//})

}
