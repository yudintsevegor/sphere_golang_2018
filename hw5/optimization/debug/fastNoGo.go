package main

import (
	"bytes"
	"strconv"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	//"regexp"
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
/**/
)
/**/
// suppress unused package warning
type User struct {
	Browsers []string `json:"browsers"`
	Company string `json:"-"`
	Country string `json:"-"`
	Email string `json:"email"`
	Job string `json:"-"`
	Name string `json:"name"`
	Phone string `json:"-"`
	isMSIE bool
	isAndroid bool

}

//const filePath string = "./data/users.txt"
// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	seenBrowsers := make(map[string]string)
	foundUsers := ""

	lines := bytes.Split(fileContents, []byte("\n"))

	users := make([]User, len(lines))
	exists := ""
	//user := new(User)
	for ind, line := range lines {
		erro := users[ind].UnmarshalJSON(line)
		if erro != nil {
			panic(err)
		}
	}
	buf := bytes.Buffer{}

	for i, user := range users {

		browsers := user.Browsers
		for  _, browserRaw := range browsers{
			if strings.Contains(browserRaw, "Android"){
				user.isAndroid = true
				seenBrowsers[browserRaw] = exists
			} else if strings.Contains(browserRaw, "MSIE"){
				user.isMSIE = true
				seenBrowsers[browserRaw] = exists
			}
		}

		if !(user.isAndroid && user.isMSIE) {
			continue
		}

		email := strings.Replace(user.Email, "@", " [at] ", -1)
		buf.WriteByte('[')
		buf.WriteString(strconv.Itoa(i))
		buf.WriteByte(']')
		buf.WriteByte(' ')
		buf.WriteString(user.Name)
		buf.WriteByte(' ')
		buf.WriteByte('<')
		buf.WriteString(email)
		buf.WriteByte('>')
		buf.WriteString("\n")
	}
	foundUsers = buf.String()
	fmt.Fprintln(out, "found users:\n" + foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
	//fmt.Fprintln(out, "Total unique browsers", uniqueBrowsers)

}

var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)
// UnmarshalJSON supports json.Unmarshaler interface
func (v *User) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9f2eff5fDecodeGolang20182599HwOptimization(&r, v)
	return r.Error()
}


func easyjson9f2eff5fDecodeGolang20182599HwOptimization(in *jlexer.Lexer, out *User) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "browsers":
			if in.IsNull() {
				in.Skip()
				out.Browsers = nil
			} else {
				in.Delim('[')
				if out.Browsers == nil {
					if !in.IsDelim(']') {
						out.Browsers = make([]string, 0, 4)
					} else {
						out.Browsers = []string{}
					}
				} else {
					out.Browsers = (out.Browsers)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Browsers = append(out.Browsers, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "Company":
			out.Company = string(in.String())
		case "Country":
			out.Country = string(in.String())
		case "email":
			out.Email = string(in.String())
		case "Job":
			out.Job = string(in.String())
		case "name":
			out.Name = string(in.String())
		case "Phone":
			out.Phone = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}


