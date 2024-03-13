package agentExec

import (
	"crypto/tls"
	"github.com/aaronchen2k/deeptest/internal/performance/runner/metrics"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/kataras/iris/v12/websocket"
	"golang.org/x/net/http2"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type ExecSession struct {
	/** quick info */
	ExecUuid string // PerformanceRoom_RunnerId_VuId

	/** exec status */
	IsRunning bool

	// used to exchange request and response data between goja and go
	CurrRequest  domain.BaseRequest
	CurrResponse domain.DebugResponse

	CurrDebugInterfaceId uint
	CurrEnvironmentId    int

	ScopedVariables map[uint][]domain.ExecVariable
	ExecScene       domain.ExecScene

	/** communication related */
	ServerUrl   string
	ServerToken string

	GojaRuntime   *goja.Runtime
	GojaRequire   *require.RequireModule
	GojaVariables *[]domain.ExecVariable
	GojaLogs      *[]string

	/** below for scenario */
	/** quick info */
	ScenarioId uint
	ProjectId  uint

	/** exec data */
	RootProcessor *Processor
	Report        *Report

	/** exec status */
	ForceStopExec bool

	CurrScenarioProcessorId uint // for interface, is is an empty placeholder used in variable opt methods
	CurrScenarioProcessor   *Processor

	ScopedCookies  map[uint][]domain.ExecCookie // only for scenario
	ScopeHierarchy map[uint]*[]uint             // only for scenario (processId -> ancestorProcessIds)

	DatapoolCursor map[string]int // only for scenario

	/** communication related */
	WsMsg *websocket.Message

	ConductorGrpcAddress string
	InfluxdbSender       metrics.MessageSender

	HttpClient  *http.Client
	Http2Client *http.Client
}

func NewInterfaceExecSession(call domain.InterfaceCall) (ret *ExecSession) {
	session := ExecSession{
		ExecUuid: call.ExecUuid,

		ScopedVariables: map[uint][]domain.ExecVariable{},
		ExecScene:       call.ExecScene,

		GojaVariables: &[]domain.ExecVariable{},
		GojaLogs:      &[]string{},

		ServerUrl:   call.ServerUrl,
		ServerToken: call.Token,

		CurrDebugInterfaceId: call.Data.DebugInterfaceId,
		CurrRequest:          domain.BaseRequest{},
		CurrResponse:         domain.DebugResponse{},
	}

	session.GojaRuntime, session.GojaRequire = InitGojaRuntime(&session)

	ret = &session

	return
}

func NewScenarioExecSession(req *ScenarioExecObj, environmentId int, wsMsg *websocket.Message) (ret *ExecSession) {
	session := ExecSession{
		ExecUuid:   req.ExecUuid,
		ScenarioId: req.ScenarioId,
		ProjectId:  req.RootProcessor.ProjectId,

		RootProcessor: req.RootProcessor,
		WsMsg:         wsMsg,

		ScopedVariables: map[uint][]domain.ExecVariable{},
		ScopedCookies:   map[uint][]domain.ExecCookie{},
		ScopeHierarchy:  map[uint]*[]uint{},

		ExecScene:         req.ExecScene,
		DatapoolCursor:    map[string]int{},
		CurrEnvironmentId: environmentId,

		GojaVariables: &[]domain.ExecVariable{},
		GojaLogs:      &[]string{},

		ServerUrl:   req.ServerUrl,
		ServerToken: req.Token,
	}

	ComputerScopeHierarchy(req.RootProcessor, &session.ScopeHierarchy)
	session.GojaRuntime, session.GojaRequire = InitGojaRuntime(&session)

	jar, _ := cookiejar.New(nil)
	session.HttpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Jar:     jar, // insert response cookies into request
		Timeout: 120 * time.Second,
	}

	session.Http2Client = &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 120 * time.Second,
	}

	ret = &session

	return
}

func (s *ExecSession) Run() {
	s.RootProcessor.Run(s)
}
