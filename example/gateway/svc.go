package gateway

import (
	"context"
	test_svc "github.com/enustah/protoc-gen-gin-gateway/example/go_grpc"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type fakeCli struct {
}

func (f fakeCli) Method1(ctx context.Context, in *test_svc.Method1Req, opts ...grpc.CallOption) (*test_svc.Method1Resp, error) {
	return &test_svc.Method1Resp{}, nil
}

func (f fakeCli) Method2(ctx context.Context, in *test_svc.Method2Req, opts ...grpc.CallOption) (*test_svc.Method2Resp, error) {
	return &test_svc.Method2Resp{}, nil
}

func (f fakeCli) Method3(ctx context.Context, in *test_svc.Method2Req, opts ...grpc.CallOption) (*test_svc.Method2Resp, error) {
	return &test_svc.Method2Resp{}, nil
}

func getTestSvcCli() test_svc.TestSvcClient {
	return fakeCli{}
}

func testSvcHandle[REQ, RESP any, handlerF func(context.Context, REQ, ...grpc.CallOption) (RESP, error)](
	ctx *gin.Context, req REQ, handler handlerF, path string, err error) {
	if err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	md := metadata.MD{}
	for k, v := range ctx.Request.Header {
		md[k] = v
	}
	grpcCtx := metadata.NewOutgoingContext(context.TODO(), md)
	resp, err := handler(grpcCtx, req)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	ctx.JSON(200, resp)
}
