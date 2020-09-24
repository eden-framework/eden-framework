package transport_grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/eden-framework/eden-framework/pkg/courier"
	"github.com/eden-framework/eden-framework/pkg/courier/httpx"
	"github.com/eden-framework/eden-framework/pkg/courier/status_error"
	"github.com/eden-framework/eden-framework/pkg/duration"
)

type GRPCRequest struct {
	RequestID  string
	BaseURL    string
	ServerName string
	Method     string
	Timeout    time.Duration
	Req        interface{}
	Metadata   courier.Metadata
}

func (grpcRequest *GRPCRequest) Do() (result courier.Result) {
	result = courier.Result{}
	d := duration.NewDuration()
	defer func() {
		logger := d.ToLogger().WithFields(logrus.Fields{
			"grpc":     grpcRequest.BaseURL,
			"service":  grpcRequest.ServerName,
			"method":   grpcRequest.Method,
			"metadata": grpcRequest.Metadata,
		})

		if result.Err == nil {
			logger.Infof("success")
		} else {
			logger.Warnf("do grpc request failed %s", result.Err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), grpcRequest.Timeout)
	defer cancel()

	conn, dialErr := grpc.DialContext(ctx, grpcRequest.BaseURL, grpc.WithInsecure(), grpc.WithCodec(&MsgPackCodec{}))
	if dialErr != nil {
		result.Err = status_error.RequestTimeout
		return
	}
	defer conn.Close()

	md := metadata.Pairs(httpx.HeaderRequestID, grpcRequest.RequestID)

	stream, errForInitial := grpc.NewClientStream(
		metadata.NewOutgoingContext(context.Background(), md),
		&grpc.StreamDesc{
			StreamName: grpcRequest.Method,
		}, conn,
		fmt.Sprintf("/%s/%s", grpcRequest.ServerName, grpcRequest.Method),
	)

	if errForInitial != nil {
		result.Err = status_error.InvalidStruct.StatusError().WithDesc(errForInitial.Error())
		return
	}

	if errForSend := stream.SendMsg(grpcRequest.Req); errForSend != nil {
		result.Err = status_error.RequestTimeout.StatusError().WithDesc(errForSend.Error())
		return
	}

	if errForClose := stream.CloseSend(); errForClose != nil {
		result.Err = status_error.RequestTimeout.StatusError().WithDesc(errForClose.Error())
		return
	}

	err := stream.RecvMsg(&result.Data)
	md, errForGetHeader := stream.Header()
	if errForGetHeader == nil {
		result.Meta = courier.Metadata(md)
	}

	if err != nil {
		if s, ok := status.FromError(err); ok {
			result.Err = status_error.ParseString(s.Message())
			return
		}
		result.Err = status_error.UnknownError
		return
	}

	result.Unmarshal = MsgPackUnmarshal
	return
}
