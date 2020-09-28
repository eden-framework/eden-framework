package transport_grpc

import (
	"fmt"
	ctx "github.com/eden-framework/eden-framework/pkg/context"
	"github.com/profzone/envconfig"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"

	"github.com/eden-framework/eden-framework/pkg/courier"
)

type ServeGRPC struct {
	IP           string
	Port         int
	WriteTimeout envconfig.Duration
	ReadTimeout  envconfig.Duration
	Name         string
}

func (s ServeGRPC) MarshalDefaults(v interface{}) {
	if srv, ok := v.(*ServeGRPC); ok {
		if srv.Name == "" {
			srv.Name = os.Getenv("PROJECT_NAME")
		}

		if srv.Port == 0 {
			srv.Port = 9000
		}

		if srv.ReadTimeout == 0 {
			srv.ReadTimeout = envconfig.Duration(15 * time.Second)
		}

		if srv.WriteTimeout == 0 {
			srv.WriteTimeout = envconfig.Duration(15 * time.Second)
		}
	}
}

func (s *ServeGRPC) Serve(wsCtx *ctx.WaitStopContext, router *courier.Router) error {
	s.MarshalDefaults(s)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		panic(err)
	}

	gs := grpc.NewServer(grpc.CustomCodec(&MsgPackCodec{}))

	serviceDesc := s.convertRouterToServiceDesc(router)
	serviceDesc.ServiceName = s.Name

	gs.RegisterService(serviceDesc, &mockServer{})

	wsCtx.Add(1)
	go func() {
		<-wsCtx.Done()
		fmt.Println("GRPC server shutdown...")
		gs.GracefulStop()
		fmt.Println("GRPC server shutdown complete.")
		wsCtx.Finish()
	}()

	fmt.Printf("[Courier] GRPC listen on %s\n", lis.Addr().String())
	return gs.Serve(lis)
}

func (s *ServeGRPC) convertRouterToServiceDesc(router *courier.Router) *grpc.ServiceDesc {
	routes := router.Routes()

	if len(routes) == 0 {
		panic(fmt.Sprintf("need to register operation before Listion"))
	}

	serviceDesc := grpc.ServiceDesc{
		HandlerType: (*MockServer)(nil),
		Methods:     []grpc.MethodDesc{},
		Streams:     []grpc.StreamDesc{},
	}

	for _, route := range routes {
		operators, operatorTypeNames := route.EffectiveOperators()

		streamDesc := grpc.StreamDesc{
			StreamName:    operatorTypeNames[len(operatorTypeNames)-1],
			Handler:       CreateStreamHandler(s, operators...),
			ServerStreams: true,
		}

		serviceDesc.Streams = append(serviceDesc.Streams, streamDesc)
	}

	return &serviceDesc
}

type MockServer interface {
}

type mockServer struct {
}
