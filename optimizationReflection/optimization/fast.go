package main

import (
	"bytes"
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
	"io"
	"strconv"
	"strings"
	"os"
	"bufio"
)

type User struct {
	Browsers  []string `json:"browsers"`
	Company   string   `json:"-"`
	Country   string   `json:"-"`
	Email     string   `json:"email"`
	Job       string   `json:"-"`
	Name      string   `json:"name"`
	Phone     string   `json:"-"`
}

const (
	exists = ""
)

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buf := bytes.Buffer{}
	buf.WriteString("found users:\n")
	seenBrowsers := make(map[string]string)
	scanner := bufio.NewScanner(file)

	ind := 0
	user := new(User)
	for scanner.Scan() {
		isMSIE := false
		isAndroid := false

		line := scanner.Bytes()
		err := user.UnmarshalJSON(line)
		if err != nil {
			panic(err)
		}
		for _, browserRaw := range user.Browsers {
			if strings.Contains(browserRaw, "Android") {
				isAndroid = true
				seenBrowsers[browserRaw] = exists
			} else if strings.Contains(browserRaw, "MSIE") {
				isMSIE = true
				seenBrowsers[browserRaw] = exists
			}
		}

		if !(isAndroid && isMSIE) {
			ind++
			continue

		}

		user.Email = strings.Replace(user.Email, "@", " [at] ", -1)
		buf.WriteByte('[')
		buf.WriteString(strconv.Itoa(ind))
		buf.WriteByte(']')
		buf.WriteByte(' ')
		buf.WriteString(user.Name)
		buf.WriteByte(' ')
		buf.WriteByte('<')
		buf.WriteString(user.Email)
		buf.WriteByte('>')
		buf.WriteString("\n")
		buf.WriteTo(out)
		buf.Reset()
		ind++
	}

	buf.WriteString("\nTotal unique browsers ")
	buf.WriteString(strconv.Itoa(len(seenBrowsers)))
	buf.WriteString("\n")
	buf.WriteTo(out)
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
