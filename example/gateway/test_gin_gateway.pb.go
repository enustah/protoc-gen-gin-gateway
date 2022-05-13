// 该文件由protoc-gen-gin-gateway生成

package gateway

import (
	"github.com/enustah/protoc-gen-gin-gateway/example/go_grpc"
	"github.com/gin-gonic/gin"
)

const (
	TestSvcTestSvcMethod1Path = "/TestSvc.TestSvc/Method1"
	TestSvcTestSvcMethod2Path = "/TestSvc.TestSvc/Method2"
	TestSvcTestSvcMethod3Path = "/TestSvc.TestSvc/Method3"
)

func RegisterTestSvcTestSvcHandler(group *gin.RouterGroup) {
	group.GET("/test/m1", TestSvcTestSvcMethod1GinHandler)
	group.POST("/test/m2", TestSvcTestSvcMethod2GinHandler)
}

// @Tags 测试服务
// @Summary 方法1
// @accept application/json
// @Produce application/json
// @Param data query go_grpc.Method1Req true "go_grpc.Method1Req"
// @success 200 {object} util.Resp{data=go_grpc.Method1Resp} "返回结果"
// @Router /test/m1 [GET]
// test svc 1 method extend comment
// test svc 2 method extend comment
// test 1 extend comment
// test 2 extend comment
func TestSvcTestSvcMethod1GinHandler(ctx *gin.Context) {
	req := &go_grpc.Method1Req{}
	err := ctx.ShouldBindQuery(req)
	testSvcHandle(ctx, req, getTestSvcCli().Method1, TestSvcTestSvcMethod1Path, err)
}

// @Tags 测试服务
// @Summary 方法2
// @accept application/json
// @Produce application/json
// @Param data body go_grpc.Method2Req true "go_grpc.Method2Req"
// @success 200 {object} util.Resp{data=go_grpc.Method2Resp} "返回结果"
// @Router /test/m2 [POST]
// test svc 1 method extend comment
// test svc 2 method extend comment
// aaaaaffff
// @Param session header string true "auth header"
func TestSvcTestSvcMethod2GinHandler(ctx *gin.Context) {
	req := &go_grpc.Method2Req{}
	err := ctx.ShouldBindJSON(req)
	testSvcHandle(ctx, req, getTestSvcCli().Method2, TestSvcTestSvcMethod2Path, err)
}
