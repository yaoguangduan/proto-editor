package generate

import (
	"fmt"
	"github.com/yaoguangduan/proto-editor/internal/generate/util"
	"go/ast"
	"go/parser"
	"go/token"
	"google.golang.org/protobuf/cmd/protoc-gen-go/internal_gengo"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"unicode"
	"unicode/utf8"
)

func genImport(gen *protogen.Plugin, g *protogen.GeneratedFile, f *util.FileInfo, imp protoreflect.FileImport) {
	impFile, ok := gen.FilesByPath[imp.Path()]
	if !ok {
		return
	}
	if impFile.GoImportPath == f.GoImportPath {
		// Don't generate imports or aliases for types in the same Go package.
		return
	}
	// Generate imports for all non-weak dependencies, even if they are not
	// referenced, because other code and tools depend on having the
	// full transitive closure of protocol buffer types in the binary.
	if !imp.IsWeak {
		g.Import(impFile.GoImportPath)
	}
	if !imp.IsPublic {
		return
	}

	// Generate public imports by generating the imported file, parsing it,
	// and extracting every symbol that should receive a forwarding declaration.
	impGen := internal_gengo.GenerateFile(gen, impFile)
	impGen.Skip()
	b, err := impGen.Content()
	if err != nil {
		gen.Error(err)
		return
	}
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, "", b, parser.ParseComments)
	if err != nil {
		gen.Error(err)
		return
	}
	genForward := func(tok token.Token, name string, expr ast.Expr) {
		// Don't import unexported symbols.
		r, _ := utf8.DecodeRuneInString(name)
		if !unicode.IsUpper(r) {
			return
		}
		// Don't import the FileDescriptor.
		if name == impFile.GoDescriptorIdent.GoName {
			return
		}
		// Don't import decls referencing a symbol defined in another package.
		// i.e., don't import decls which are themselves public imports:
		//
		//	type T = somepackage.T
		if _, ok := expr.(*ast.SelectorExpr); ok {
			return
		}
		g.P(tok, " ", name, " = ", impFile.GoImportPath.Ident(name))
	}
	g.P("// Symbols defined in public import of ", imp.Path(), ".")
	g.P()
	for _, decl := range astFile.Decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				switch spec := spec.(type) {
				case *ast.TypeSpec:
					genForward(decl.Tok, spec.Name.Name, spec.Type)
				case *ast.ValueSpec:
					for i, name := range spec.Names {
						var expr ast.Expr
						if i < len(spec.Values) {
							expr = spec.Values[i]
						}
						genForward(decl.Tok, name.Name, expr)
					}
				case *ast.ImportSpec:
				default:
					panic(fmt.Sprintf("can't generate forward for spec type %T", spec))
				}
			}
		}
	}
	g.P()
}
