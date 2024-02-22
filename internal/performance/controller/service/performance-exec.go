package controllerService

import (
	"context"
	"github.com/aaronchen2k/deeptest/internal/performance/controller/dao"
	"github.com/aaronchen2k/deeptest/internal/performance/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/performance/pkg/domain"
	ptlog "github.com/aaronchen2k/deeptest/internal/performance/pkg/log"
	websocketHelper "github.com/aaronchen2k/deeptest/internal/performance/pkg/websocket"
	ptProto "github.com/aaronchen2k/deeptest/internal/performance/proto"
	_logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"github.com/facebookgo/inject"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
	"github.com/nxadm/tail"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"sync"
)

type PerformanceTestService struct {
	req ptdomain.PerformanceTestReq

	execCtx    context.Context
	execCancel context.CancelFunc
	client     *ptProto.PerformanceServiceClient

	GrpcService         *GrpcService         `inject:"private"`
	ScheduleService     *ScheduleService     `inject:"private"`
	RemoteRunnerService *RemoteRunnerService `inject:"private"`
}

func NewPerformanceTestServiceRef(req ptdomain.PerformanceTestReq) *PerformanceTestService {
	service := &PerformanceTestService{
		req: req,
	}

	var g inject.Graph
	g.Logger = logrus.StandardLogger()

	if err := g.Provide(
		&inject.Object{Value: service},
	); err != nil {
		logrus.Fatalf("provide usecase objects to the Graph: %v", err)
	}

	err := g.Populate()
	if err != nil {
		logrus.Fatalf("populate the incomplete Objects: %v", err)
	}

	return service
}

func (s *PerformanceTestService) Connect(runner *ptProto.Runner) (client ptProto.PerformanceServiceClient) {
	connect, err := grpc.Dial(runner.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}

	client = ptProto.NewPerformanceServiceClient(connect)

	return
}

func (s *PerformanceTestService) ExecStart(req ptdomain.PerformanceTestReq, wsMsg *websocket.Message) (err error) {
	s.execCtx, s.execCancel = context.WithCancel(context.Background())

	dao.ClearData(req.Room)
	dao.ResetInfluxdb(req.Room, req.InfluxdbAddress, req.InfluxdbOrg, req.InfluxdbToken)
	s.ScheduleService.Reset(req.Scenarios)
	s.GrpcService.ClearAllGlobalVar(context.Background(), &ptProto.GlobalVarRequest{})

	// stop execution in 2 ways:
	// 1. call cancel in this method by websocket request OR after all runners completed
	// 2. sub cancel instruction from runner via grpc

	go s.ScheduleService.SendMetricsToClient(s.execCtx, s.execCancel, req, wsMsg)

	var wgRunners sync.WaitGroup
	for _, runner := range req.Runners {
		client := s.Connect(runner)

		stream, err := s.CallRunnerExecStartByGrpc(client, req, runner.Id, runner.Name, runner.Weight)
		if err != nil {
			continue
		}

		wgRunners.Add(1)

		s.HandleGrpcMsgAsync(stream, wgRunners)
	}

	wgRunners.Wait()

	s.execCancel()

	websocketHelper.SendExecInstructionToClient("", "", ptconsts.MsgInstructionEnd, req.Room, wsMsg)

	return
}

func (s *PerformanceTestService) ExecStop(wsMsg *websocket.Message) (err error) {
	// stop server execution
	if s.execCancel != nil {
		s.execCancel()
	}

	// exec by runners
	s.RemoteRunnerService.CallStop(s.req)

	// send end msg to websocket client
	websocketHelper.SendExecInstructionToClient("", "", ptconsts.MsgInstructionEnd, s.req.Room, wsMsg)

	return
}

func (s *PerformanceTestService) SendLogAsync(req ptdomain.PerformanceTestReq, wsMsg *websocket.Message) {
	go func() {
		s.SendLog(req, wsMsg)
	}()
}

func (s *PerformanceTestService) SendLog(req ptdomain.PerformanceTestReq, wsMsg *websocket.Message) (err error) {
	room := req.Room
	logPath := ptlog.GetLogPath(room)

	t, err := tail.TailFile(logPath, tail.Config{Follow: true, ReOpen: true})
	if err != nil {
		return
	}

	var arr []string

	for line := range t.Lines {
		arr = append(arr, line.Text)

		if len(arr) > 100 {
			data := iris.Map{
				"log": line.Text,
			}
			websocketHelper.SendExecResultToClient(data, ptconsts.MsgResultRecord, req.Room, wsMsg)

			arr = make([]string, 0)
		}

		select {
		case <-s.execCtx.Done():
			_logUtils.Debug("<<<<<<< stop sendLog job")
			return

		default:
		}
	}

	s.execCancel()

	return
}

func (s *PerformanceTestService) CallRunnerExecStartByGrpc(
	client ptProto.PerformanceServiceClient, req ptdomain.PerformanceTestReq, runnerId int32, runnerName string, weight int32) (stream ptProto.PerformanceService_ExecStartClient, err error) {

	stream, err = client.ExecStart(context.Background())
	if err != nil {
		ptlog.Logf(err.Error())
		return
	}

	runnerExecScenarios := s.getRunnerExecScenarios(req, runnerId)

	err = stream.Send(&ptProto.PerformanceExecStartReq{
		Room:       req.Room,
		RunnerId:   runnerId,
		RunnerName: runnerName,
		Title:      req.Title,

		Mode:      req.Mode.String(),
		Scenarios: runnerExecScenarios,
		Weight:    weight,

		ServerAddress:   req.ServerAddress,
		InfluxdbAddress: req.InfluxdbAddress,
		InfluxdbOrg:     req.InfluxdbOrg,
		InfluxdbToken:   req.InfluxdbToken,
	})

	if err != nil {
		return
	}

	return
}

func (s *PerformanceTestService) getRunnerExecScenarios(req ptdomain.PerformanceTestReq, runnerId int32) (
	ret []*ptProto.Scenario) {

	notSet := false
	scenarioIdsMap := map[int32]bool{}
	for _, runner := range req.Runners {
		if runner.Id != runnerId {
			continue
		}

		if runner.Scenarios == nil {
			notSet = true

			break
		}

		for _, scId := range runner.Scenarios {
			scenarioIdsMap[scId] = true
		}

		break
	}

	for _, scenario := range req.Scenarios {
		if notSet || scenarioIdsMap[scenario.Id] {
			ret = append(ret, scenario)
		}
	}

	return
}

func (s *PerformanceTestService) HandleGrpcMsgAsync(stream ptProto.PerformanceService_ExecStartClient,
	wgRunners sync.WaitGroup) {

	go func() {
		defer wgRunners.Done()

		s.HandleGrpcMsg(stream)

		stream.CloseSend()
	}()
}

func (s *PerformanceTestService) HandleGrpcMsg(stream ptProto.PerformanceService_ExecStartClient) (err error) {
	for true {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		ptlog.Logf("get grpc msg from runner: %v", resp)

		// dealwith Instruction from agent
		if resp.Instruction == ptconsts.MsgInstructionRunnerFinish.String() {
			break
		}

		select {
		case <-s.execCtx.Done():
			_logUtils.Debug("<<<<<<< stop sendLog job")
			break

		default:
		}
	}

	return
}
