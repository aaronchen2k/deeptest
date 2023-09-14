package mockGenerator

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
)

type Request struct {
	Method consts.HttpMethod `json:"method"`
}

type Response struct {
	StatusCode  consts.HttpRespCode    `json:"statusCode"`
	ContentType consts.HttpContentType `json:"contentType"`
	Data        interface{}            `json:"data"`

	// used by adv mock
	Headers []model.EndpointMockExpectResponseHeader `json:"headers"`

	UseAdvMock bool   `json:"useAdvMock"`
	Content    string `json:"content"`
	DelayTime  uint   `json:"delayTime"`
}
