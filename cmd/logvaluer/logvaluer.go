package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/types"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"
)

// The implementation heavily references the Stringer interface.
// https://cs.opensource.google/go/x/tools/+/refs/tags/v0.12.0:cmd/stringer/stringer.go

var (
	typeName   = flag.String("type", "", "type name must be set")
	outConsole = flag.Bool("console", false, "output console instead of file")
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of logvaluer:\n")
	fmt.Fprintf(os.Stderr, "\tlogvaluer [flag] -type T\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	flag.Usage = Usage
	flag.Parse()
	if len(*typeName) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	g := &Generator{targetName: *typeName}
	if err := g.parsePackage(); err != nil {
		log.Fatal(err)
	}

	g.Printf("// Code generated by \"logvaluer\"; DO NOT EDIT.\n")
	g.Printf("\n")
	g.Printf("package %s\n", g.pkgName)
	g.Printf("\n")
	g.Printf("import \"log/slog\"")
	g.Printf("\n")

	g.generate()

	out := g.format()
	if *outConsole {
		fmt.Println(string(out))
	}

	outputName := fmt.Sprintf("%s_logvalue.go", strings.ToLower(g.targetName))
	err := os.WriteFile(outputName, out, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

type Generator struct {
	buf     bytes.Buffer
	pkgName string

	targetName string
	target     types.Type
}

func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

func (g *Generator) parsePackage() error {
	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax,
		Tests: false,
	}

	pkgs, err := packages.Load(cfg, []string{"."}...) // only current directory
	if err != nil {
		return err
	}
	g.pkgName = pkgs[0].Name

	for _, astFile := range pkgs[0].Syntax {
		g.extractTarget(astFile, pkgs[0].TypesInfo)
	}

	return nil
}

func (g *Generator) extractTarget(file *ast.File, info *types.Info) {
	ast.Inspect(file, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			if typeSpec.Name.String() == g.targetName {
				switch typeSpec.Type.(type) {
				case *ast.StructType:
					typ := info.TypeOf(typeSpec.Name)
					if structTyp, ok := typ.Underlying().(*types.Struct); ok {
						g.target = structTyp
					}
				default:
					// *ast.Ident
					log.Printf("%v %s\n", typeSpec.Name, "not implement")
				}
			}
		}
		return true
	})
}

type field struct {
	name string
	typ  string
	mask bool
}

func (g *Generator) generate() {
	switch t := g.target.(type) {
	case *types.Struct:
		g.generateForStruct(t)
	}
}

func (g *Generator) generateForStruct(st *types.Struct) {
	var fields []field
	for i := 0; i < st.NumFields(); i++ {
		ignoredTag := reflect.StructTag(st.Tag(i)).Get("ignored")
		if isTrue(ignoredTag) {
			continue
		}

		f := field{
			name: st.Field(i).Name(),
		}

		maskTag := reflect.StructTag(st.Tag(i)).Get("mask")
		if isTrue(maskTag) {
			f.mask = true
		}

		switch st.Field(i).Type().String() {
		case "bool":
			f.typ = "Bool"
		case "float64":
			f.typ = "Float64"
		case "string":
			f.typ = "String"
		case "int":
			f.typ = "Int"
		case "int64":
			f.typ = "Int64"
		case "Uint64":
			f.typ = "Uint64"
		case "time.Time":
			f.typ = "Time"
		case "time.Duration":
			f.typ = "Duration"
		default:
			f.typ = "Any"
		}
		fields = append(fields, f)
	}

	r := strings.ToLower(g.targetName)[0]
	g.Printf("func (%c %s) LogValue() slog.Value {\n", r, g.targetName)
	g.Printf("return slog.GroupValue(\n")
	for _, field := range fields {
		g.Printf("slog.%s(\"%s\", ", field.typ, field.name)
		if field.mask {
			g.Printf("\"MASK\"")
		} else {
			g.Printf("%c.%s", r, field.name)
		}
		g.Printf("),\n")
	}
	g.Printf(")")
	g.Printf("}\n")
}

func (g *Generator) format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return g.buf.Bytes()
	}
	return src
}

func isTrue(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}
