package main

import (
	"bytes"
	json "encoding/json"
	"fmt"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
)

type User struct {
	Browsers  []string `json:"browsers"`
	Company   string   `json:"-"`
	Country   string   `json:"-"`
	Email     string   `json:"email"`
	Job       string   `json:"-"`
	Name      string   `json:"name"`
	Phone     string   `json:"-"`
	isMSIE    bool
	isAndroid bool
}

const (
	exists = ""
)

//const filePath string = "./data/users.txt"
// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	buf := bytes.Buffer{}
	seenBrowsers := make(map[string]string)
	foundUsers := ""
	lines := bytes.Split(fileContents, []byte("\n"))
	users := make([]User, len(lines))

	for i, line := range lines {
		wg.Add(1)
		go func(line []byte, user *User) {
			defer wg.Done()
			err := user.UnmarshalJSON(line)
			if err != nil {
				panic(err)
			}
			for _, browserRaw := range user.Browsers {
			//if ok, err := regexp.MatchString("Android", browser); ok && err == nil {
				if strings.Contains(browserRaw, "Android") {
					user.isAndroid = true
					mu.Lock()
					seenBrowsers[browserRaw] = exists
					mu.Unlock()
				} else if strings.Contains(browserRaw, "MSIE") {
					user.isMSIE = true
					mu.Lock()
					seenBrowsers[browserRaw] = exists
					mu.Unlock()
				}
			}
		}(line, &users[i])
	}
	wg.Wait()

	for i := 0; i < len(users); i++ {
		if !(users[i].isAndroid && users[i].isMSIE) {
			continue
		}
		users[i].Email = strings.Replace(users[i].Email, "@", " [at] ", -1)
		/**/buf.WriteByte('[')
		buf.WriteString(strconv.Itoa(i))
		buf.WriteByte(']')
		buf.WriteByte(' ')
		buf.WriteString(users[i].Name)
		buf.WriteByte(' ')
		buf.WriteByte('<')
		buf.WriteString(users[i].Email)
		buf.WriteByte('>')
		buf.WriteString("\n")
		/**/
		//foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, users[i].Name, users[i].Email)
	}
	foundUsers = buf.String()
	fmt.Fprintln(out, "found users:\n"+foundUsers)
	//io.WriteString(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
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
