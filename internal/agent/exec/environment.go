package agentExec

import (
	"fmt"
	execUtils "github.com/aaronchen2k/deeptest/internal/agent/exec/utils/exec"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	_httpUtils "github.com/aaronchen2k/deeptest/pkg/lib/http"
	"strings"
)

func GenRequestUrlWithBaseUrlAndPathParam(session *ExecSession, req *domain.BaseRequest, debugInterfaceId uint, baseUrl string) {
	execScene := session.ExecScene

	// get base url by key consts.KEY_BASE_URL in Environment Variables from server
	envId := execScene.DebugInterfaceToEnvMap[debugInterfaceId]
	vars := execScene.EnvToVariables[envId]
	if baseUrl == "" {
		vari, _ := GetVariableFromList(consts.KEY_BASE_URL, vars)
		baseUrl = fmt.Sprintf("%v", vari.Value)
	}

	if req.PathParams != nil {
		req.Url = ReplacePathParams(req.Url, *req.PathParams)
	}

	notUseBaseUrl := execUtils.IsNotUseBaseUrl(consts.ScenarioDebug, req.ProcessorInterfaceSrc)
	if !notUseBaseUrl {
		req.Url = _httpUtils.CombineUrls(baseUrl, req.Url)
	}
}

func ReplacePathParams(uri string, pathParams []domain.Param) (ret string) {
	ret = uri
	for _, param := range pathParams {
		if param.ParamIn != consts.ParamInPath {
			continue
		}

		vari := fmt.Sprintf("{%v}", param.Name)

		ret = strings.ReplaceAll(uri, vari, param.Value)
	}

	return
}
