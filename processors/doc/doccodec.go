package doc

import (
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"regexp"
	"strings"

	"github.com/k0kubun/pp"
)

func NewCodec(pkgPath string) (*Codec, error) {
	dp := &Codec{}

	fset := token.NewFileSet()
	filter := isGoFile
	pkgs, e := parser.ParseDir(fset, pkgPath, filter, parser.ParseComments)
	if e != nil {
		return nil, e
	}

	// astf := make([]*ast.File, 0)
	// for _, pkg := range pkgs {
	// 	fmt.Printf("package %v\n", pkg.Name)
	// 	for fn, f := range pkg.Files {
	// 		fmt.Printf("file %v\n", fn)

	// 		astf = append(astf, f)
	// 	}
	// }
	var astPkg *ast.Package

	for _, pkg := range pkgs {
		astPkg = pkg
	}

	docPkg := doc.New(astPkg, pkgPath, doc.AllDecls)
	dp.PkgName = docPkg.Name

	dp.ImportPath = docPkg.ImportPath

	dp.Doc = removeSpecialComment(docPkg.Doc)

	for _, typ := range docPkg.Types {
		docPkg.Consts = append(docPkg.Consts, typ.Consts...)
		docPkg.Vars = append(docPkg.Vars, typ.Vars...)
		docPkg.Funcs = append(docPkg.Funcs, typ.Funcs...)
	}

	for _, v := range docPkg.Types {
		if v.Name == "encoder" {
			dp.Encoder = &Encoder{
				Doc: removeSpecialComment(v.Doc),
			}
			continue
		}

		if v.Name == "decoder" {
			dp.Decoder = &Decoder{
				Doc: removeSpecialComment(v.Doc),
			}
			continue
		}

		if v.Name == "encoderOptions" || v.Name == "decoderOptions" {
			options := &CodecOptions{}

			options.Doc = removeSpecialComment(v.Doc)
			options.Options = []*CodecOption{}
			for _, si := range v.Decl.Specs {
				s := si.(*ast.TypeSpec)

				typ := s.Type.(*ast.StructType)
				for _, field := range typ.Fields.List {
					dpo := &CodecOption{}

					var fieldType string
					fieldTags := map[string]string{}

					customType := ""
					if field.Doc != nil {
						for _, c := range field.Doc.List {
							if strings.HasPrefix(c.Text, "// @Default ") {
								dpo.DefaultValue = strings.TrimPrefix(c.Text, "// @Default ")
							}
							if strings.HasPrefix(c.Text, "// @ExampleLS ") {
								dpo.ExampleLS = strings.TrimPrefix(c.Text, "// @ExampleLS ")
							}
							if strings.HasPrefix(c.Text, "// @Type ") {
								customType = strings.ToLower(strings.TrimPrefix(c.Text, "// @Type "))
							}

							if strings.HasPrefix(c.Text, "// @Enum ") {
								st := strings.ToLower(strings.TrimPrefix(c.Text, "// @Enum "))
								dpo.PossibleValues = strings.Split(st, ",")
							}
						}
					}

					dpo.Doc = removeSpecialComment(field.Doc.Text())
					dpo.Name = field.Names[0].String()

					switch t := field.Type.(type) {
					case *ast.MapType:
						fieldType = "map"
						keyKind := t.Key.(*ast.Ident).Name
						valueKind := "string"
						fieldType = "map[" + keyKind + "]" + valueKind
						fieldType = "hash"
					case *ast.ArrayType:
						fieldType = "array of " + t.Elt.(*ast.Ident).Name
						fieldType = "array"
					case *ast.Ident:
						fieldType = t.Name
					case *ast.SelectorExpr:
						xKind := t.X.(*ast.Ident).Name
						selKind := t.Sel.String()
						fieldType = xKind + "." + selKind
					default:
						fieldType = "unknow"
						pp.Println("field-->", field.Type)
					}
					if field.Tag != nil {
						if field.Tag.Value != "" {
							r, _ := regexp.Compile(`([a-z]*):"([a-z_0-9,]*)"`)
							for _, match := range r.FindAllStringSubmatch(field.Tag.Value, 5) {
								fieldTags[match[1]] = match[2]
							}
						}
					}
					// pp.Println("field tag-->", field.Tag.Value)
					dpo.Type = fieldType
					if customType != "" {
						dpo.Type = customType
					}

					if _, ok := fieldTags["mapstructure"]; ok {
						dpo.Alias = fieldTags["mapstructure"]
					}
					if _, ok := fieldTags["validate"]; ok {
						validationTagValues := strings.Split(fieldTags["validate"], ",")
						for _, validationTagValue := range validationTagValues {
							if validationTagValue == "required" {
								dpo.Required = true
							}
						}
					}

					options.Options = append(options.Options, dpo)
				}
			}

			switch v.Name {
			case "encoderOptions":
				dp.Encoder.Options = options
			case "decoderOptions":
				dp.Decoder.Options = options
			}

			continue
		}
	}
	return dp, nil

}
