package transport_grpc

import (
	"context"
	"reflect"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/profzone/eden-framework/pkg/courier"
	"github.com/profzone/eden-framework/pkg/courier/httpx"
	"github.com/profzone/eden-framework/pkg/courier/status_error"
	"github.com/profzone/eden-framework/pkg/duration"
	logContext "github.com/profzone/eden-framework/pkg/log/context"
)

type DecodeStreamFunc func(c context.Context, data []byte) (request interface{}, err error)

var (
	ContextKeyServerName = uuid.New().String()
)

func ContextWithServiceName(ctx context.Context, serverName string) context.Context {
	return context.WithValue(ctx, ContextKeyServerName, serverName)
}

func CreateStreamHandler(s *ServeGRPC, ops ...courier.IOperator) grpc.StreamHandler {
	opMetas := courier.ToOperatorMetaList(ops...)

	return func(_ interface{}, stream grpc.ServerStream) (err error) {
		ctx := stream.Context()
		ctx = ContextWithServiceName(ctx, s.Name)

		reqID := getRequestID(ctx)

		if reqID == "" {
			reqID = uuid.New().String()
		}

		logContext.SetLogID(reqID)
		defer logContext.Close()

		d := duration.NewDuration()

		defer func() {
			fields := logrus.Fields{
				"tag":       "access",
				"log_id":    reqID,
				"remote_ip": ClientIP(ctx),
				"method":    "/" + opMetas[len(opMetas)-1].Type.Name(),
			}

			fields["request_time"] = d.Get()

			logger := logrus.WithFields(fields)

			if err != nil {
				statusErr := status_error.FromError(err)
				if statusErr.Status() >= 500 {
					logger.Errorf(err.Error())
				} else {
					logger.Warnf(err.Error())
				}
			} else {
				logger.Infof("")
			}
		}()

		opDecode := createGRPCStreamDecoder(receiveMsgData(stream))

		for _, opMeta := range opMetas {
			op, decodeErr := courier.NewOperatorBy(opMeta.Type, opMeta.Operator, opDecode)
			if decodeErr != nil {
				err = passErr(ctx, decodeErr)
				return
			}

			response, endpointErr := op.Output(ctx)
			if endpointErr != nil {
				err = passErr(ctx, endpointErr)
				return
			}

			if !opMeta.IsLast {
				// set result in context with key of operator name
				ctx = context.WithValue(ctx, opMeta.ContextKey, response)
				continue
			}

			encodeErr := sendMsg(ctx, stream, response)
			if encodeErr != nil {
				err = passErr(ctx, encodeErr)
				return
			}
		}
		return
	}
}

func createGRPCStreamDecoder(data []byte) courier.OperatorDecoder {
	return func(op courier.IOperator, rv reflect.Value) (err error) {
		err = msgpack.Unmarshal(data, op)
		if err != nil {
			err = status_error.InvalidStruct.StatusError().WithDesc(err.Error())
			return
		}
		return
	}
}

func getRequestID(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if values, ok := md[httpx.HeaderRequestID]; ok {
			if len(values) > 0 {
				return values[0]
			}
		}
	}
	return ""
}

func GetFieldDisplayName(field reflect.StructField) string {
	pathName := field.Name
	jsonName, _ := field.Tag.Lookup("json")
	if jsonName != "" {
		pathName = jsonName
	}
	return pathName
}

func sendMsg(_ context.Context, stream grpc.ServerStream, response interface{}) (err error) {
	md := metadata.Pairs(httpx.HeaderRequestID, logContext.GetLogID())
	if canMeta, ok := response.(courier.IMeta); ok {
		md = metadata.Join(md, metadata.MD(canMeta.Meta()))
	}
	if err = stream.SetHeader(md); err != nil {
		return
	}
	return stream.SendMsg(response)
}

func passErr(ctx context.Context, err error) error {
	if err == nil {
		return err
	}
	if _, ok := status.FromError(err); !ok {
		finalStatusErr := status_error.FromError(err)
		err = status.Error(CodeFromHTTPStatus(finalStatusErr.Status()), finalStatusErr.WithSource(ctx.Value(ContextKeyServerName).(string)).String())
	}

	return err
}

func receiveMsgData(stream grpc.ServerStream) (data []byte) {
	stream.RecvMsg(&data)
	return
}

func MarshalOperator(stream grpc.ServerStream, operator courier.IOperator) error {
	opDecode := createGRPCStreamDecoder(receiveMsgData(stream))
	rv := reflect.Indirect(reflect.ValueOf(operator))
	op, err := courier.NewOperatorBy(rv.Type(), operator, opDecode)
	if err != nil {
		return err
	}
	rv.Set(reflect.ValueOf(op).Elem())
	return nil
}
