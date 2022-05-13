package generate

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"google.golang.org/protobuf/compiler/protogen"
	"html/template"
	"io"
	"net/http"
	"strings"
)

type funcParam struct {
	name []string
	typ  string
}

/*
	构建函数块 没有body

	func (structTyp)name(param)retType{
	}
*/
func funcBlockNode(name, structTyp string, param []funcParam, retType []string) *ast.FuncDecl {
	var recv *ast.FieldList
	if structTyp != "" {
		recv = &ast.FieldList{
			List: []*ast.Field{
				{
					Type: ast.NewIdent(structTyp),
				},
			},
		}
	}

	paramList := make([]*ast.Field, 0, len(param))
	for _, p := range param {
		names := make([]*ast.Ident, 0, len(p.name))
		for _, n := range p.name {
			names = append(names, ast.NewIdent(n))
		}
		paramList = append(paramList, &ast.Field{
			Names: names,
			Type:  ast.NewIdent(p.typ),
		})
	}
	result := make([]*ast.Field, 0, len(retType))
	for _, v := range retType {
		result = append(result, &ast.Field{
			Type: ast.NewIdent(v),
		})
	}
	return &ast.FuncDecl{
		Recv: recv,
		Name: &ast.Ident{
			Name: name,
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: paramList,
			},
			Results: &ast.FieldList{
				List: result,
			},
		},
		Body: nil,
	}
}

// 单个服务构建ast.Spec
type SvcGenerator struct {
	svc       *protogen.Service
	pkgName   string
	goPkgName string
}

func NewSvcGenerator(svc *protogen.Service, pkgName, goPkgName string) *SvcGenerator {
	return &SvcGenerator{
		svc:       svc,
		pkgName:   pkgName,
		goPkgName: goPkgName,
	}
}

func (this *SvcGenerator) getConstPathVarName(method *protogen.Method) string {
	return this.pkgName + string(this.svc.Desc.Name()) + string(method.Desc.Name()) + "Path"
}

func (this *SvcGenerator) getConstPathVal(method *protogen.Method) string {
	return "/" + this.pkgName + "." + string(this.svc.Desc.Name()) + "/" + string(method.Desc.Name())
}

/*
常量块

const(
	...
)
*/
func (this *SvcGenerator) getSvcConstNode() []ast.Spec {
	var (
		n = make([]ast.Spec, 0)
	)
	for _, m := range this.svc.Methods {
		fullDec := this.getConstPathVarName(m) // 完整变量声明 ${package}${svcName}${methodName}Path
		fullVal := this.getConstPathVal(m)     // 完整值 /${package}.${svcName}/${methodName}
		n = append(n, &ast.ValueSpec{
			Names: []*ast.Ident{
				{
					Name: fullDec,
				},
			},
			Values: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"" + fullVal + "\"",
				},
			},
		})

	}
	return n
}

// gin handler 函数名
func (this *SvcGenerator) getGinHandlerFuncName(method *protogen.Method) string {
	return this.pkgName + string(this.svc.Desc.Name()) + string(method.Desc.Name()) + "GinHandler"
}

const swagComment = `// @Tags {{.Tag}}
// @Summary {{.Sum}}
// @accept application/json
// @Produce application/json
// @Param data {{.ParamIn}} {{.ReqStruct}} true "{{.ReqStruct}}"
// @success 200 {object} {{.RespObj}} "返回结果"
// @Router {{.Path}} [{{.Method}}]
`

func (this *SvcGenerator) genSwagComment(m *protogen.Method) ([]string, error) {
	type swagRender struct {
		Tag       string
		Sum       string
		ParamIn   string
		ReqStruct string
		RespObj   string
		Path      string
		Method    string
	}
	var (
		svcGenInfo    = getSvcGenInfo(this.svc)
		methodGenInfo = getMethodGenInfo(m)
		r             = &swagRender{
			Tag:       svcGenInfo.SwagTag,
			Sum:       methodGenInfo.SwagSummary,
			ParamIn:   "",
			ReqStruct: fmt.Sprintf("%s.%s", this.goPkgName, m.Input.Desc.Name()),
			RespObj:   "",
			Path:      methodGenInfo.Path,
			Method:    methodGenInfo.Method,
		}
	)
	if r.Tag == "" || r.Sum == "" || methodGenInfo.SwagRespObj == "" {
		return nil, fmt.Errorf("%s", "tag,summary,respObj can not nil when generate_swag is true")
	}
	if strings.Count(methodGenInfo.SwagRespObj, "%s") != 1 {
		return nil, fmt.Errorf("RespObj format %s illegal. example is %s", methodGenInfo.SwagRespObj, "util.CommResp{data=%s} where %s will replace to resp struct")
	}
	r.RespObj = fmt.Sprintf(methodGenInfo.SwagRespObj, fmt.Sprintf("%s.%s", this.goPkgName, m.Output.Desc.Name()))
	switch methodGenInfo.Method {
	case http.MethodGet, http.MethodHead, http.MethodDelete:
		r.ParamIn = "query"
	case http.MethodPost, http.MethodPut:
		r.ParamIn = "body"
	default:
		return nil, fmt.Errorf("unknown http method %s", methodGenInfo.Method)
	}

	tmp, err := template.New("").Parse(swagComment)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(nil)
	if err := tmp.Execute(buf, r); err != nil {
		panic(err)
	}
	var (
		bufRd = bufio.NewReader(buf)
		ret   []string
	)
	for {
		s, _, err := bufRd.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		ret = append(ret, string(s))

	}
	return ret, nil
}

/*
生成gin handler完整函数块

func ${packageName}${svcName}${methodName}(ctx *gin.Context){
	req:=&xxx{}
	err:=ctx.ShouldBindxxx(req)
	handleFunc(ctx,req,getGrpcCliFunc().MethodName,err)
}
*/
func (this *SvcGenerator) getGinHandlerFuncNode(method *protogen.Method) (*ast.FuncDecl, error) {
	// 函数声明
	methodGenInfo := getMethodGenInfo(method)
	if methodGenInfo == nil {
		return nil, nil
	}
	svcGenInfo := getSvcGenInfo(this.svc)
	var (
		ginBindFunc = ""
	)

	switch methodGenInfo.Method {
	case http.MethodGet, http.MethodHead, http.MethodDelete:
		ginBindFunc = "ShouldBindQuery"
	case http.MethodPost, http.MethodPut:
		ginBindFunc = "ShouldBindJSON"
	default:
		return nil, fmt.Errorf("unknown http method %s", methodGenInfo.Method)
	}

	f := funcBlockNode(this.getGinHandlerFuncName(method), "", []funcParam{
		{
			name: []string{"ctx"},
			typ:  "*gin.Context",
		},
	}, nil)
	// 函数逻辑
	f.Body = &ast.BlockStmt{
		List: []ast.Stmt{
			// req:=&${goPackageName}{}
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent("req"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					ast.NewIdent(fmt.Sprintf("&%s.%s{}", this.goPkgName, method.Input.Desc.Name())),
				},
			},

			// err:=ctx.ShouldBindxxx(req)
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent("err"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					ast.NewIdent(fmt.Sprintf("ctx.%s(req)", ginBindFunc)),
				},
			},
			// handleFunc(ctx,req,getGrpcCliFunc().MethodName,path,err)
			&ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: ast.NewIdent(svcGenInfo.HandlerFunc),
					Args: []ast.Expr{
						ast.NewIdent("ctx"),
						ast.NewIdent("req"),
						&ast.SelectorExpr{
							X: &ast.CallExpr{
								Fun: ast.NewIdent(svcGenInfo.GetGrpcCliFunc),
							},
							Sel: ast.NewIdent(string(method.Desc.Name())),
						},
						// ast.NewIdent(fmt.Sprintf("%s.%s", svcGenInfo.GetGrpcCliFunc, method.Desc.Name())),
						ast.NewIdent(this.getConstPathVarName(method)),
						ast.NewIdent("err"),
					},
				},
			},
		},
	}
	// 生成swag 注释
	if svcGenInfo.GenerateSwag {
		swag, err := this.genSwagComment(method)
		if err != nil {
			return nil, err
		}
		f.Doc = &ast.CommentGroup{}
		for _, v := range swag {
			f.Doc.List = append(f.Doc.List, &ast.Comment{
				Text: v,
			})
		}
	}
	// 生成额外注释
	for _, v := range svcGenInfo.MethodExtendComment {
		f.Doc.List = append(f.Doc.List, &ast.Comment{
			Text: "// " + v,
		})
	}

	for _, v := range methodGenInfo.ExtendComment {
		f.Doc.List = append(f.Doc.List, &ast.Comment{
			Text: "// " + v,
		})
	}
	return f, nil
}

/*
生成gin路由注册函数

func registerXXXHandler(group *gin.RouterGroup){
	group.XXX(path,handlerFunc)
	...
	...
	...
}
*/
func (this *SvcGenerator) getGinHandlerRegisterNode() *ast.FuncDecl {
	f := funcBlockNode(fmt.Sprintf("Register%s%sHandler", this.pkgName, this.svc.Desc.Name()), "", []funcParam{
		{
			name: []string{"group"},
			typ:  "*gin.RouterGroup",
		},
	}, nil)
	f.Body = &ast.BlockStmt{}
	for _, m := range this.svc.Methods {
		methodGenInfo := getMethodGenInfo(m)
		if methodGenInfo == nil {
			continue
		}
		// group.XXX(path,handlerFunc)
		stmt := &ast.ExprStmt{
			X: ast.NewIdent(fmt.Sprintf("group.%s(%s,%s)", methodGenInfo.Method, `"`+methodGenInfo.Path+`"`, this.getGinHandlerFuncName(m))),
		}
		f.Body.List = append(f.Body.List, stmt)
	}

	return f

}

/*
根据服务构建[]ast.Spec 对于每个服务 都生成以下代码

const(
	//grpc 的路径
)

func Register${package}${svcName}Handler(group *gin.RouterGroup){
	// 例如 group.POST("/login", login)
    group.${requestMethod}(${path}, handlerFunc)
	...
}

// method名称对应上面注册的handlerFunc
func ${method}(ctx *gin.Context){
	req := &${goPkgName}.${reqStruct}{}
	err:= ctx.ShouldBindxxx(req)

	//这个函数需要另外实现
	${handler_func_name}(ctx,req,${get_grpc_cli_func_name}.${Method},${GrpcMethodPath},err)
}

...

*/
func (this *SvcGenerator) GetSvcNode() ([]ast.Decl, error) {
	consNode := this.getSvcConstNode()
	regNode := this.getGinHandlerRegisterNode()

	decls := []ast.Decl{
		&ast.GenDecl{
			Tok:   token.CONST,
			Specs: consNode,
		},
		regNode,
	}
	for _, v := range this.svc.Methods {
		funcNode, err := this.getGinHandlerFuncNode(v)
		if err != nil {
			return nil, err
		}
		if funcNode == nil {
			continue
		}
		decls = append(decls, funcNode)
	}
	return decls, nil
}
