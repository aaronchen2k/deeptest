package performance

import (
	controllerService "github.com/aaronchen2k/deeptest/internal/agent/performance/conductor/exec"
	ptProto "github.com/aaronchen2k/deeptest/internal/agent/performance/proto"
	"github.com/aaronchen2k/deeptest/internal/pkg/config"
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"google.golang.org/grpc"
	"log"
	"net"
)

func StartGrpcServe() {
	server := grpc.NewServer()
	ptProto.RegisterPerformanceServiceServer(server, &controllerService.GrpcService{})

	lis, err := net.Listen("tcp", config.CONFIG.System.GrpcAddress)
	if err != nil {
		log.Fatalf("grpc net.Listen err: %v", err)
	}
	server.Serve(lis)

	logUtils.Info("grpc")
}
