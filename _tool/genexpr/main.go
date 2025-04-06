package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type AstField struct {
	Type      string
	Name      string
	ClassName string
}

type AstClass struct {
	Name   string
	Fields []AstField
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run astgen.go <output_directory>")
		os.Exit(1)
	}

	outputDir := os.Args[1]
	packageName := "ast"

	exprs := []string{
		"Binary   : Expr left, Token operator, Expr right",
		"Grouping : Expr expression",
		"Literal  : any value",
		"Unary    : Token operator, Expr right",
	}

	err := defineAst(outputDir, packageName, "Expr", exprs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func defineAst(outputDir, packageName, baseName string, types []string) error {
	outputPath := filepath.Join(outputDir, strings.ToLower(baseName)+".go")
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	classes := parseTypes(types, []string{})

	// First collect all unique fields across all types
	allFields := make(map[string]string) // map[fieldName]fieldType
	for _, class := range classes {
		for _, field := range class.Fields {
			allFields[field.Name] = field.Type
		}
	}

	// Template for the generated code
	tmpl := `package {{.PackageName}}

{{makeImports .Imports}}

// {{.BaseName}} is the interface for all expression types
type {{.BaseName}} interface {
{{range .InterfaceMethods}}	{{.}}
{{end}}}

{{range .Classes}}// {{.Name}}Expr represents a {{.Name}} expression
type {{.Name}}Expr struct {
{{range .Fields}}	{{.Name}} {{.Type}}
{{end}}}

{{end}}// Method implementations for all expression types
{{range $className, $classFields := .ClassFieldMap}}{{range $fieldName, $fieldType := $.AllFields}}func (e *{{$className}}Expr) {{makeGetter $fieldName}}() {{$fieldType}} {
{{if hasField $classFields $fieldName}}	return e.{{$fieldName}}
{{else}}	return nil // This expression type doesn't have this field
{{end}}}

{{end}}{{end}}// New{{.BaseName}} functions for creating expression instances
{{range .Classes}}func New{{.Name}}Expr({{makeParams .Fields}}) *{{.Name}}Expr {
	return &{{.Name}}Expr{
{{makeInitializers .Fields}}
	}
}

{{end}}`

	// Determine interface methods based on all unique fields
	methods := make([]string, 0, len(allFields))
	for fieldName, fieldType := range allFields {
		methodName := makeGetterName(fieldName)
		methods = append(methods, methodName+"() "+fieldType)
	}

	// Create a map for class fields lookup
	classFieldMap := make(map[string][]AstField)
	for _, class := range classes {
		classFieldMap[class.Name] = class.Fields
	}

	// Remove imports if we're in the same package
	imports := []string{}
	if packageName != "ast" {
		imports = append(imports, "github.com/bagaswh/rottenlang/pkg/ast")
	}

	data := struct {
		PackageName      string
		BaseName         string
		Classes          []AstClass
		InterfaceMethods []string
		Imports          []string
		AllFields        map[string]string
		ClassFieldMap    map[string][]AstField
	}{
		PackageName:      packageName,
		BaseName:         baseName,
		Classes:          classes,
		InterfaceMethods: methods,
		Imports:          imports,
		AllFields:        allFields,
		ClassFieldMap:    classFieldMap,
	}

	funcMap := template.FuncMap{
		"makeGetter": makeGetterName,
		"makeParams": func(fields []AstField) string {
			params := make([]string, len(fields))
			for i, field := range fields {
				params[i] = field.Name + " " + field.Type
			}
			return strings.Join(params, ", ")
		},
		"makeInitializers": func(fields []AstField) string {
			inits := make([]string, len(fields))
			for i, field := range fields {
				inits[i] = "\t\t" + field.Name + ": " + field.Name + ","
			}
			return strings.Join(inits, "\n")
		},
		"makeImports": func(imports []string) string {
			if len(imports) == 0 {
				return ""
			}

			str := "import (\n"
			for _, imp := range imports {
				str += fmt.Sprintf("\t\"%s\"\n", imp)
			}
			return str + ")"
		},
		"hasField": func(fields []AstField, fieldName string) bool {
			for _, field := range fields {
				if field.Name == fieldName {
					return true
				}
			}
			return false
		},
	}

	t, err := template.New("ast").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return err
	}

	return t.Execute(file, data)
}

func parseTypes(types []string, interfaceMethods []string) []AstClass {
	var classes []AstClass

	for _, typeSpec := range types {
		// Split class name and fields
		parts := strings.Split(typeSpec, ":")
		className := strings.TrimSpace(parts[0])

		var fields []AstField
		if len(parts) > 1 {
			fieldSpecs := strings.Split(strings.TrimSpace(parts[1]), ",")

			for _, fieldSpec := range fieldSpecs {
				fieldSpec = strings.TrimSpace(fieldSpec)
				if fieldSpec == "" {
					continue
				}

				fieldParts := strings.Fields(fieldSpec)

				// Ensure we have type and name
				if len(fieldParts) >= 2 {
					fieldType := strings.TrimSpace(fieldParts[0])
					fieldName := strings.TrimSpace(fieldParts[1])

					// Handle pointer types like *Token
					if fieldType != "Expr" && fieldType != "any" {
						// Don't add * if it's already a pointer type
						if !strings.HasPrefix(fieldType, "*") {
							fieldType = "*" + fieldType
						}
					}

					fields = append(fields, AstField{
						Type:      fieldType,
						Name:      fieldName,
						ClassName: className,
					})
				}
			}
		}

		classes = append(classes, AstClass{
			Name:   className,
			Fields: fields,
		})
	}

	return classes
}

func makeGetterName(fieldName string) string {
	return strings.ToUpper(fieldName[:1]) + fieldName[1:]
}
