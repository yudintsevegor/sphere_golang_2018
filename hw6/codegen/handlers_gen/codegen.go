package main

// код писать тут

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"
	"text/template"
)

type ApiGenLabel struct {
	URL    string
	Auth   bool
	Method string
}

type FieldUrlMethod struct {
	Field string
	URL string
	Method string
}

type tplStructMethodURL struct {
	StructName string
	MethodName string
	URL string
}

type tplAtoi struct {
	FieldName string
}

type tplFieldValue struct {
	FieldName  string
	Value      string
	LowerField string
}

type tplCheck struct {
	Value string
}

var returnError string = `
func returnError(ansError string, StatusCode int, w http.ResponseWriter) {
	w.WriteHeader(StatusCode)
	answer := map[string]interface{}{
		"error": ansError,
	}
	jsonRes, _ := json.Marshal(answer)
	w.Write(jsonRes)
}
`

var checkMethod string = `
	if r.Method != http.MethodPost {
		ansError = "bad method"
		returnError(ansError, http.StatusNotAcceptable, w)
		return
	}
`
var checkAuth string = `
	var auth string
	if len(r.Header["X-Auth"]) != 0 {
		auth = r.Header["X-Auth"][0]
	}

	if auth != "100500" {
		ansError = "unauthorized"
		returnError(ansError, http.StatusForbidden, w)
		return
	}
`
/*
var callReturnError string = `
		returnError(ansError, BadReq, w)
		return
	}
`
*/
var badReq string = `
	BadReq := http.StatusBadRequest
`

var middleServ string = `
	ctx := r.Context()
	switch r.URL.Path {
`

var endServ string = `
	default:
		ansError := "unknown method"
		returnError(ansError, http.StatusNotFound, w)
	}
}
`

var (
	headFunc = template.Must(template.New("").Parse(`
func (srv *{{.StructName}}) wrapper{{.MethodName}}{{.StructName}}(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var ansError string

`))
	atoi = template.Must(template.New("").Parse(`
	{{.FieldName}}, errAtoi := strconv.Atoi(r.FormValue("{{.FieldName}}"))
`))
	defaultValue = template.Must(template.New("").Parse(`
	{{.LowerField}} := r.FormValue("{{.LowerField}}")
	if {{.LowerField}} == "" {
		{{.LowerField}} = "{{.Value}}"
	}
`))
	minValueLenth = template.Must(template.New("").Parse(`
	if len({{.LowerField}}) < {{.Value}} {
		ansError = "{{.LowerField}} len must be >= {{.Value}}"
		returnError(ansError, BadReq, w)
		return
	}
`))
	minValue = template.Must(template.New("").Parse(`
	if {{.LowerField}} < 1 {
		ansError = "{{.LowerField}} must be >= {{.Value}}"
		returnError(ansError, BadReq, w)
		return
	}
`))
	maxValue = template.Must(template.New("").Parse(`
	if {{.LowerField}} > {{.Value}} {
		ansError = "{{.LowerField}} must be <= {{.Value}}"
		returnError(ansError, BadReq, w)
		return
	}
`))
	reqValue = template.Must(template.New("").Parse(`
	{{.LowerField}} := r.FormValue("{{.LowerField}}")
	if {{.LowerField}} == "" {
		ansError = "{{.LowerField}} must be not empty"
		returnError(ansError, BadReq, w)
		return
	}

`))
	enumValue = template.Must(template.New("").Parse(`
	_, ok := check[{{.LowerField}}]
	if !ok {
		ansError = "{{.LowerField}} must be one of [{{.Value}}]"
		returnError(ansError, BadReq, w)
		return
	}
`))
	enumCheck = template.Must(template.New("").Parse(`
		"{{.Value}}": true,
`))
	intValue = template.Must(template.New("").Parse(`
	{{.LowerField}}, errAtoi := strconv.Atoi(r.FormValue("{{.LowerField}}"))
	if errAtoi != nil {
		ansError = "{{.LowerField}} must be int"
		returnError(ansError, BadReq, w)
		return
	}
`))
	paramValue = template.Must(template.New("").Parse(`
	{{.LowerField}} := r.FormValue("{{.Value}}")
`))

	structValue = template.Must(template.New("").Parse(`
		{{.FieldName}}: {{.LowerField}},
`))

	endFunc = template.Must(template.New("").Parse(`
	res, err := srv.{{.MethodName}}(ctx, params)
	if err != nil {
		switch err.Error() {
		case "bad user":
			returnError(err.Error(), http.StatusInternalServerError, w)
			return
		default:
			returnError(err.Error(), err.(ApiError).HTTPStatus, w)
			return
		}
	}
	answer := map[string]interface{}{
		"error":    ansError,
		"response": res,
	}
	jsonRes, _ := json.Marshal(answer)
	w.Write(jsonRes)
}
`))
	serv = template.Must(template.New("").Parse(`
func (srv *{{.StructName}}) ServeHTTP(w http.ResponseWriter, r *http.Request) {
`))

	url = template.Must(template.New("").Parse(`
	case "{{.URL}}":
		srv.wrapper{{.MethodName}}{{.StructName}}(ctx, w, r)
`))

)

var FUM = make([]FieldUrlMethod, 0)

func main() {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, os.Args[1], nil, parser.ParseComments)

	if err != nil {
		log.Fatal(err)
	}
	out, _ := os.Create(os.Args[2])

	fmt.Fprintln(out, `package `+node.Name.Name)
	fmt.Fprintln(out) // empty line
	fmt.Fprintln(out, `import "context"`)
	fmt.Fprintln(out, `import "encoding/json"`)
	fmt.Fprintln(out, `import "net/http"`)
	fmt.Fprintln(out, `import "strconv"`)
	fmt.Fprintln(out) // empty line
	fmt.Fprintln(out, returnError)

	regMin := regexp.MustCompile("min=([0-9]+)")
	regMax := regexp.MustCompile("max=([0-9]+)")
	regEnum := regexp.MustCompile("enum=([^,]*)")
	regDefault := regexp.MustCompile("default=([^,]*)")
	regParamname := regexp.MustCompile("paramname=([^,]*)")
	for _, d := range node.Decls {
		apiLab := &ApiGenLabel{}
		fDecl, okFunc := d.(*ast.FuncDecl)
		if !okFunc {
			continue
		}
		if fDecl.Doc == nil {
			continue
		}
		needCodegen := false
		for _, comment := range fDecl.Doc.List {
			needCodegen = needCodegen || strings.HasPrefix(comment.Text, "// apigen:api")
			splitComment := strings.Split(comment.Text, "")
			toJSON := splitComment[14:len(splitComment)]
			jsonStr := strings.Join(toJSON, "")
			data := []byte(jsonStr)
			json.Unmarshal(data, apiLab)

			structName := fDecl.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name
			methodName := fDecl.Name.Name
			headFunc.Execute(out, tplStructMethodURL{structName, methodName,""})
			fum := FieldUrlMethod{structName, apiLab.URL, methodName}
			FUM = append(FUM ,fum)

			if apiLab.Method == "POST" {
				fmt.Fprintln(out, checkMethod)
			}
			if apiLab.Auth == true {
				fmt.Fprintln(out, checkAuth)
			}
		}
		paramsStruct := fDecl.Type.Params.List[1].Type.(*ast.Ident).Name

		for _, d2 := range node.Decls {
			gDecl, okGen := d2.(*ast.GenDecl)
			if !okGen {
				continue
			}
			for _, spec := range gDecl.Specs {
				currType, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				currStruct, ok := currType.Type.(*ast.StructType)
				if !ok {
					continue
				}
				if currType.Name.Name != paramsStruct {
					continue
				}
				fmt.Fprintf(out, badReq)

				sliceOfFieldsName:= make([]string, 0)
			FIELDS_LOOP:
				for _, field := range currStruct.Fields.List {
					if field.Tag != nil {
						tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
						fieldName := field.Names[0].Name
						sliceOfFieldsName = append(sliceOfFieldsName, fieldName)
						fileType := field.Type.(*ast.Ident).Name

						if tag.Get("apivalidator") == "-" {
							continue FIELDS_LOOP
						}
						apiVal := tag.Get("apivalidator")

						ansDef := regDefault.FindStringSubmatch(apiVal)
						ansMin := regMin.FindStringSubmatch(apiVal)
						ansMax := regMax.FindStringSubmatch(apiVal)
						ansParam := regParamname.FindStringSubmatch(apiVal)
						ansEnum := regEnum.FindStringSubmatch(apiVal)

						lowerField := strings.ToLower(fieldName)
						if fileType == "int" {
							intValue.Execute(out, tplFieldValue{"", "", lowerField})
							//fmt.Fprintf(out, callReturnError)
						}

						if ansDef != nil {
							defaultValue.Execute(out, tplFieldValue{"", ansDef[1], lowerField})
						}
						if strings.Contains(apiVal, "required") {
							reqValue.Execute(out, tplFieldValue{"", "", lowerField})
							//fmt.Fprintf(out, callReturnError)
						}

						if ansMin != nil {
							if fileType == "string" {
								minValueLenth.Execute(out, tplFieldValue{"", ansMin[1], lowerField})
							} else if fileType == "int" {
								minValue.Execute(out, tplFieldValue{"", ansMin[1], lowerField})
							}
							//fmt.Fprintf(out, callReturnError)
						}
						if ansMax != nil {
							maxValue.Execute(out, tplFieldValue{"", ansMax[1], lowerField})
							//fmt.Fprintf(out, callReturnError)
						}
						if ansEnum != nil {
							sliceOfEnum := strings.Split(ansEnum[1], "|")
							fmt.Fprintf(out, "\tcheck := map[string]bool{")
							toTpl := strings.Join(sliceOfEnum, ", ")
							for _, value := range sliceOfEnum {
								enumCheck.Execute(out, tplCheck{value})
							}
							fmt.Fprintf(out, "\t}")
							enumValue.Execute(out, tplFieldValue{"", toTpl, lowerField})
							//fmt.Fprintf(out, callReturnError)
						}
						if ansParam != nil {
							paramValue.Execute(out, tplFieldValue{"", ansParam[1], lowerField})
						}
					}
				} // END FIELD LOOP
				fmt.Fprintf(out, "\tparams := " + paramsStruct + "{")
				for _, value := range sliceOfFieldsName{
					structValue.Execute(out, tplFieldValue{value, "", strings.ToLower(value)})
				}
				fmt.Fprintf(out, "\t}")
				endFunc.Execute(out, tplStructMethodURL{"", fDecl.Name.Name, ""})
			} // END Spec LOOP
		} //END Inner LOOP node.Decl
	}
	mapField := make(map[string][]string, 0)
	for _, value := range FUM{
		sliceField := make([]string, 0, 1)
		_, ok := mapField[value.Field]
		if  ok {
			mapField[value.Field] = append(mapField[value.Field], value.URL)
			mapField[value.Field] = append(mapField[value.Field], value.Method)
			continue
		}
		sliceField = append(sliceField, value.URL)
		sliceField = append(sliceField, value.Method)
		mapField[value.Field] = sliceField
	}
	for key, value := range mapField {
		serv.Execute(out, tplStructMethodURL{key, "", ""})
		fmt.Fprintf(out, middleServ)
		for i := 0; i < len(value)/2; i++  {
			url.Execute(out, tplStructMethodURL{key, value[2*i + 1], value[2*i]})
		}
		fmt.Fprintf(out, endServ)
	}
}


