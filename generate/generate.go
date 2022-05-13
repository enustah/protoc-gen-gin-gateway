package generate

import (
	"errors"
	"github.com/enustah/protoc-gen-gin-gateway/gin_gateway"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

func getExtendImport(file *protogen.File) []string {
	extImport := proto.GetExtension(file.Desc.Options(), gin_gateway.E_ExtendImport)
	if extImport == nil {
		return nil
	}
	return extImport.(*gin_gateway.ExtendImport).Import
}

func getGinGWPackageName(file *protogen.File) string {
	fileExt := proto.GetExtension(file.Desc.Options(), gin_gateway.E_GinGatewayPackageName)
	if fileExt == nil {
		return ""
	}
	return fileExt.(string)
}

func getSvcGenInfo(service *protogen.Service) *gin_gateway.ServiceGenInfo {
	_svcGenInfo := proto.GetExtension(service.Desc.Options(), gin_gateway.E_SvcGenInfo)
	if _svcGenInfo == nil {
		return nil
	}
	svcGenInfo := _svcGenInfo.(*gin_gateway.ServiceGenInfo)
	return svcGenInfo
}

func getMethodGenInfo(method *protogen.Method) *gin_gateway.MethodGenInfo {
	_mGenInfo := proto.GetExtension(method.Desc.Options(), gin_gateway.E_MethodGenInfo)
	if _mGenInfo == nil {
		return nil
	}
	return _mGenInfo.(*gin_gateway.MethodGenInfo)
}

func GenerateFile(gen *protogen.Plugin, file *protogen.File) error {
	ginGWPkgName := getGinGWPackageName(file)
	if ginGWPkgName == "" {
		return errors.New("option gin_gateway_package_name can not nil")
	}
	fg := NewFileGenerator(file.Services, getExtendImport(file), string(file.Desc.Package().Name()), ginGWPkgName, string(file.GoPackageName))
	c, err := fg.Generate()
	if err != nil {
		return err
	}
	if c == "" {
		return nil
	}
	filename := file.GeneratedFilenamePrefix + "_gin_gateway.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("// 该文件由protoc-gen-gin-gateway生成")
	g.P("")
	g.P(c)
	return nil
}
