## 这是使用例子


### 先安装
```bash
go install . 
# 或者 go install github.com/enustah/protoc-gen-gin-gateway
```

### 编写pb文件 (假设文件在example/api/src)
首先是通用的开头
```protobuf
syntax = "proto3";

package TestSvc; 
option go_package = ".;test_svc";
```

复制本项目里的[gin_gateway_option.proto](../gin_gateway_option.proto) 复制到其他目录(这里假设example/api/common),import到pb文件
```protobuf
import "gin_gateway_option.proto";
```
然后是两个fileOption定义
```protobuf
option (extend_import) ={
  import:[
    "github.com/enustah/protoc-gen-gin-gateway/example/go_grpc"
  ]
};

option (gin_gateway_package_name) = "gateway";
```
    extend_import 是生成gin代码时 import块额外添加的内容,可选.默认只有一个github.com/gin-gonic/gin
    一般extend_import 需要添加grpc所在的库
    gin_gateway_package_name 是生成的package名称,这个必须声明



下面是具体服务定义
```protobuf
service TestSvc {
  //服务选项 没有option则跳过整个服务的网关代码生成
  option(svc_gen_info) = {
    //是否生成swag注释, 如果生成,swag_tag不能空
    generate_swag: true,
    // swagger route 里的路径前缀 完整路径是 ${swag_path_prefix}${path}
    swag_path_prefix: "/prefix",
    // swag 注释的tag
    swag_tag: "测试服务", 
    //这个是获取grpc客户端的函数,不能有参数 需要自己实现
    /*
    func getTestSvcCli()test_svc.TestSvcClient{
    
    }
    */
    //如果函数来自其他包,需要吧这个包的import 加入到extend_import
    get_grpc_cli_func: "getTestSvcCli", 
    //这个是具体处理逻辑 需要自己实现
    /*
    func testSvcHandle[REQ, RESP any, handlerF func(context.Context, REQ) (RESP, error)](
    	ctx *gin.Context, req REQ, handler handlerF, path string, err error) {
    	
    }
    */
    //req resp 对应grpc的入参和出參类型,path对应完整grpc的http2路径,err是ShouldBindJSON/ShouldBindQuery 返回的错误
    //如果函数来自其他包,需要吧这个包的import 加入到extend_import
    handler_func: "testSvcHandle",

    //额外注释 添加到每个生成函数的上方注释
    method_extend_comment: [
      "test svc 1 method extend comment", // 这个逗号必须加
      "test svc 2 method extend comment"

    ]
  };


  rpc Method1(Method1Req) returns(Method1Resp){
    //函数选项 没有option则不生成这个函数的网关
    option (method_gen_info) = {
      //对应gin注册的路径
      path:"/test/m1",
      //对应请求方法 目前支持 GET POST HEAD PUT DELETE
      //其中 GET HEAD DELETE 是query请求 POST 和 PUT 是JSON 目前固定不可变
      method:"GET",
      //swag的summary  generate_swag=true时不能空
      swag_summary: "方法1",
      //swag的返回结构  generate_swag=true时不能空 %s 会替换成grpc的响应结构
      swag_resp_obj: "util.Resp{data=%s}",
      //额外注释 这里可以添加额外的swag注释
      extend_comment: [
        "@Param session header string true \"auth header\""
      ]
    };
  }

  rpc Method2(Method2Req) returns(Method2Resp){
    option (method_gen_info) = {
      path:"/test/m2",
      method:"POST",
      swag_summary: "方法2",
      swag_resp_obj: "util.Resp{data=%s}",
    };

  }
}

message Method1Req{

}

message Method1Resp{

}

message Method2Req{

}

message Method2Resp{

}
```

然后生成代码 在api目录执行
```bash
protoc -I./common -I. --gin-gateway_out=../gateway src/*.proto
#如果要生成grpc 和 go 
# protoc -I./common -I. --go_out=../go_grpc --go-grpc_out=require_unimplemented_servers=false:../go_grpc --gin-gateway_out=../gateway src/*.proto
```

其中 getTestSvcCli 和 testSvcHandle 需要自己实现([example](gateway/svc.go))









