package astparser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strings"

	"github.com/fatih/structtag"
)

// Parse parses the code and return a parser
func Parse(code string) (*Parser, error) {
	return ParseBytes([]byte(code))
}

// ParseBytes parses the code and return a parser
func ParseBytes(code []byte) (psr *Parser, err error) {
	psr = new(Parser)
	fileSet := token.NewFileSet()
	var file *ast.File
	if file, err = parser.ParseFile(fileSet, "", code, parser.ParseComments); err != nil {
		return
	}

	// get import codes
	if psr.Imports, err = parseImports(file, fileSet); err != nil {
		return
	}

	// get struct codes
	if psr.Structs, err = parseStruts(file, fileSet); err != nil {
		return
	}

	// get value codes
	if psr.Values, err = parseValues(file); err != nil {
		return
	}

	// get interface codes
	if psr.Interfaces, err = parseInterfaces(file, fileSet); err != nil {
		return
	}

	psr.init()
	return
}

// parseImports get import list
func parseImports(file *ast.File, fileSet *token.FileSet) ([]string, error) {
	imports := make([]string, len(file.Imports))
	for i, spec := range file.Imports {
		// get code block
		var buffer bytes.Buffer
		if err := format.Node(&buffer, fileSet, spec); err != nil {
			return nil, err
		}
		str := buffer.String()
		imports[i] = str[1 : len(str)-1]
	}
	return imports, nil
}

// parseStructs get structs
func parseStruts(file *ast.File, fileSet *token.FileSet) (structs []*Struct, err error) {
	defer func() {
		if err != nil {
			structs = make([]*Struct, 0)
		}
	}()
	inspectTypeSpec(file, func(spec *ast.TypeSpec) (bool, error) {
		typ, ok := spec.Type.(*ast.StructType)
		if !ok {
			return false, nil
		}
		var st = new(Struct)
		// set struct name
		st.Name = spec.Name.Name

		// set struct doc
		if spec.Doc != nil {
			st.Doc = spec.Doc.Text()
		}
		// parse field list
		if typ.Fields != nil {
			st.Fields, err = parseStructFields(st.Name, typ.Fields, fileSet)
			if err != nil {
				return false, err
			}
		}
		structs = append(structs, st)
		return true, nil
	})
	return structs, nil
}

// parseStrutFields get struct field list
func parseStructFields(structName string, fields *ast.FieldList, fileSet *token.FileSet) (structFields []*StructField, err error) {
	for _, fieldSpec := range fields.List {
		var structField = new(StructField)
		// is it anonymous field
		structField.Anonymous = len(fieldSpec.Names) == 0
		if !structField.Anonymous {
			structField.Name = fieldSpec.Names[0].Name
		} else {
			structField.Name, err = getCodeBlock(fileSet, fieldSpec.Type)
			if err != nil {
				return
			}
		}

		// set field type
		if structField.Type, err = getCodeBlock(fileSet, fieldSpec.Type); err != nil {
			return
		}

		// set field doc
		if fieldSpec.Doc != nil {
			structField.Doc = fieldSpec.Doc.Text()
		}

		// set field tags
		if fieldSpec.Tag != nil && len(fieldSpec.Tag.Value) > 2 {
			var (
				tagValue = fieldSpec.Tag.Value[1 : len(fieldSpec.Tag.Value)-1]
			)
			structField.Tag = tagValue
			structField.Tags, err = structtag.Parse(tagValue)
			if err != nil {
				return
			}
		}
		structFields = append(structFields, structField)
	}
	return
}

// parseValues get variables and constants
func parseValues(file *ast.File) (values []*Value, err error) {
	ast.Inspect(file, func(node ast.Node) bool {
		genDecl, ok := node.(*ast.GenDecl)
		if !ok {
			return true
		}
		for _, spec := range genDecl.Specs {
			var val = new(Value)
			switch t := spec.(type) {
			case *ast.ValueSpec:
				val.Name = t.Names[0].Name

				if t.Doc != nil {
					val.Doc = t.Doc.Text()
				}

				var v *ast.BasicLit
				v, ok = t.Values[0].(*ast.BasicLit)
				if !ok {
					continue
				}
				val.Kind = strings.ToLower(v.Kind.String())
				val.Val = v.Value
				if val.Kind == "string" && len(val.Val) > 2 {
					val.Val = val.Val[1 : len(val.Val)-1]
				}
			default:
				continue
			}
			values = append(values, val)
		}
		return true
	})

	return values, nil
}

// parseInterfaces get interfaces
func parseInterfaces(file *ast.File, fileSet *token.FileSet) (interfaces []*Interface, err error) {
	defer func() {
		if err != nil {
			interfaces = make([]*Interface, 0)
		}
	}()
	inspectTypeSpec(file, func(spec *ast.TypeSpec) (bool, error) {
		// parse interface's functions
		typ, ok := spec.Type.(*ast.InterfaceType)
		if !ok {
			return false, nil
		}
		var itf = new(Interface)
		// set interface name
		itf.Name = spec.Name.Name

		// set interface doc
		if spec.Doc != nil {
			itf.Doc = spec.Doc.Text()
		}

		if typ.Methods != nil {
			itf.Functions, err = parseInterfaceFunctions(typ.Methods, fileSet)
			if err != nil {
				return false, err
			}
		}
		interfaces = append(interfaces, itf)
		return true, nil
	})
	return interfaces, nil
}

// parseInterfaceFunctions get interface's function list
func parseInterfaceFunctions(list *ast.FieldList, fileSet *token.FileSet) (methods []*Function, err error) {
	for _, fnSpec := range list.List {
		var fn = new(Function)
		// is it anonymous
		fn.Anonymous = len(fnSpec.Names) == 0
		if !fn.Anonymous {
			fn.Name = fnSpec.Names[0].Name
		} else {
			fn.Name, err = getCodeBlock(fileSet, fnSpec.Type)
			if err != nil {
				return
			}
		}

		// set function doc
		if fnSpec.Doc != nil {
			fn.Doc = fnSpec.Doc.Text()
		}
		switch f := fnSpec.Type.(type) {
		case *ast.FuncType:
			// set function input params
			if f.Params != nil {
				fn.InputParams = parseFunctionFields(f.Params)
			}

			// set function return params
			if f.Results != nil {
				fn.ReturnParams = parseFunctionFields(f.Results)
			}
		default:
			continue
		}
		methods = append(methods, fn)
	}
	return
}

// parseFunctionFields get function's field list
func parseFunctionFields(list *ast.FieldList) []string {
	var results []string
	for _, t := range list.List {
		var expr = t.Type
		if e, ok := t.Type.(*ast.StarExpr); ok {
			expr = e.X
		}
		switch ft := expr.(type) {
		case *ast.Ident:
			results = append(results, ft.Name)
		default:
			continue
		}
	}
	return results
}

func inspectTypeSpec(file *ast.File, fn func(spec *ast.TypeSpec) (bool, error)) {
	ast.Inspect(file, func(node ast.Node) bool {
		genDecl, ok := node.(*ast.GenDecl)
		if !ok {
			return true
		}
		var err error
		for _, spec := range genDecl.Specs {
			switch t := spec.(type) {
			case *ast.TypeSpec:
				if t.Type == nil {
					continue
				}
				if ok, err = fn(t); err != nil {
					return true
				} else if !ok {
					continue
				}
			default:
				continue
			}
		}
		return true
	})
}

func getCodeBlock(fileSet *token.FileSet, i interface{}) (string, error) {
	var dst bytes.Buffer
	err := format.Node(&dst, fileSet, i)
	if err != nil {
		return "", fmt.Errorf("getCodeBlock failure, error: %s", err.Error())
	}
	return dst.String(), nil
}
