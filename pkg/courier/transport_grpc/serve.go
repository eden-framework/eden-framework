package transport_grpc

import (
	"fmt"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"

	"github.com/profzone/eden-framework/pkg/conf"
	"github.com/profzone/eden-framework/pkg/courier"
)

type ServeGRPC struct {
	IP           string
	Port         int
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
	Name         string
}

func (s ServeGRPC) DockerDefaults() conf.DockerDefaults {
	return conf.DockerDefaults{
		"Port": 999,
	}
}

func (s ServeGRPC) MarshalDefaults(v interface{}) {
	if grpc, ok := v.(*ServeGRPC); ok {
		if grpc.Name == "" {
			grpc.Name = os.Getenv("PROJECT_NAME")
		}

		if grpc.Port == 0 {
			grpc.Port = 777
		}

		if grpc.ReadTimeout == 0 {
			grpc.ReadTimeout = 15 * time.Second
		}

		if grpc.WriteTimeout == 0 {
			grpc.WriteTimeout = 15 * time.Second
		}
	}
}

func (s *ServeGRPC) Serve(router *courier.Router) error {
	s.MarshalDefaults(s)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		panic(err)
	}

	gs := grpc.NewServer(grpc.CustomCodec(&MsgPackCodec{}))

	serviceDesc := s.convertRouterToServiceDesc(router)
	serviceDesc.ServiceName = s.Name

	gs.RegisterService(serviceDesc, &mockServer{})

	fmt.Printf("[Courier] listen on %s\n", lis.Addr().String())

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
