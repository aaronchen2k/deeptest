package openapi

import (
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/getkin/kin-openapi/openapi3"
	"strings"
)

func ConvertServersToEnvironments(servers openapi3.Servers) (vars []model.EnvironmentVar, err error) {
	for _, server := range servers {
		vari := model.EnvironmentVar{
			Name: "server",
		}

		vari.Value = genUrl(server.URL, server.Variables)

		vars = append(vars, vari)
	}

	return
}

func genUrl(url string, variables map[string]*openapi3.ServerVariable) (ret string) {
	ret = strings.TrimRight(url, "/")

	for name, svar := range variables {
		ret = strings.ReplaceAll(ret, "${"+name+"}", svar.Default)
	}

	return
}
