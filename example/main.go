package main

import (
	"github.com/enustah/protoc-gen-gin-gateway/example/gateway"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	gateway.RegisterTestSvcTestSvcHandler(r.Group("/"))
	panic(r.Run(":8080"))
}
