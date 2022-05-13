package generate

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"google.golang.org/protobuf/compiler/protogen"
	"io/ioutil"
)

type FileGenerator struct {
	packageName      string // proto 文件的 package 声明
	ginGWPackageName string // gin gateway option 的 package_name
	goPackageName    string // proto 文件的 go_package 声明
	extImport        []string
	svc              []*protogen.Service
}

func NewFileGenerator(service []*protogen.Service, extImport []string, pkgName, ginGWPkgName, goPkgName string) *FileGenerator {
	g := &FileGenerator{
		packageName:      pkgName,
		ginGWPackageName: ginGWPkgName,
		goPackageName:    goPkgName,
		extImport:        extImport,
		svc:              service,
	}

	return g
}

func (this *FileGenerator) GetImportNode() []ast.Spec {
	n := make([]ast.Spec, 0, len(this.extImport))
	for _, v := range this.extImport {
		n = append(n, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `"` + v + `"`,
			},
		})
	}
	return n
}

// 根据定义的svc构建 []ast.Spec
func (this *FileGenerator) buildSvcNode() {

}

func (this *FileGenerator) BuildAst() (*ast.File, error) {
	extImportNode := this.GetImportNode()
	// 构建语法树
	root := &ast.File{
		Name: ast.NewIdent(this.ginGWPackageName),
		Decls: []ast.Decl{
			// ==========import 声明块================
			// import 声明会根据extImport 拓展
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: append([]ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"github.com/gin-gonic/gin"`,
						},
					},
				}, extImportNode...),
			},
			// ==========import 声明块 end================
		},
	}
	declLen := len(root.Decls)
	for _, svc := range this.svc {
		if getSvcGenInfo(svc) == nil { // svc 没有设置option 跳过
			continue
		}
		decl, err := NewSvcGenerator(svc, this.packageName, this.goPackageName).GetSvcNode()
		if err != nil {
			return nil, err
		}
		root.Decls = append(root.Decls, decl...)
	}
	if len(root.Decls) == declLen {
		return nil, nil
	}
	return root, nil
}

func (this *FileGenerator) Generate() (string, error) {
	f := token.NewFileSet()
	root, err := this.BuildAst()
	if err != nil || root == nil {
		return "", err
	}
	b := bytes.NewBuffer(nil)
	if err := printer.Fprint(b, f, root); err != nil {
		panic(err)
	}
	s, _ := ioutil.ReadAll(b)
	return string(s), nil
}
