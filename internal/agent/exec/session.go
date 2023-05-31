package agentExec

import (
	"crypto/tls"
	"github.com/kataras/iris/v12/websocket"
	"golang.org/x/net/http2"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type Session struct {
	ScenarioId uint
	Name       string

	HttpClient  *http.Client
	Http2Client *http.Client
	Failfast    bool

	RootProcessor *Processor
	Report        *Report

	WsMsg *websocket.Message

	Step *step
}

func NewSession(req *ScenarioExecObj, failfast bool, wsMsg *websocket.Message) (ret *Session) {
	root := req.RootProcessor
	//root.Result.ScenarioId = req.RootProcessor.ScenarioId
	// TODO: now, interfaces use variables in its own serve's env
	//ImportVariables(root.ID, req.Variables, consts.Public)

	session := Session{
		ScenarioId:    root.ScenarioId,
		Name:          req.Name,
		RootProcessor: root,
		Failfast:      failfast,
		WsMsg:         wsMsg,
		Step:          &step{},
	}

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

func (s *Session) Run() {
	s.RootProcessor.Run(s)
}

type step struct {
	Id int
}

func (s *step) GetId() int {
	s.Id = s.Id + 1
	return s.Id
}
