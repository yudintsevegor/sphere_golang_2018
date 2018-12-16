package main

// тут писать SearchServer
import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

type XmlUser struct {
	Id        int    `xml:"id"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Age       int    `xml:"age"`
	About     string `xml:"about"`
	Gender    string `xml:"gender"`
}

type toClient struct {
	Id     int
	Name   string
	Age    int
	About  string
	Gender string
}

type XmlUsers struct {
	List []XmlUser `xml:"row"`
}

type SearchServerErrorResponse struct {
	Error string
}

func (usr *XmlUsers) XmlSorting(orderField string, orderBy int, query string) ([]XmlUser, error) {
	var dataBase []XmlUser

	if query != "" {
		for _, value := range usr.List {
			if strings.Contains(value.About, query) || strings.Contains(value.FirstName+" "+value.LastName, query) {
				dataBase = append(dataBase, value)
			}
		}

		if len(dataBase) == 0 {
			for _, value := range usr.List {
				dataBase = append(dataBase, value)
			}
		}
	} else {
		for _, value := range usr.List {
			dataBase = append(dataBase, value)
		}
	}

	if orderBy == 0 || orderBy == 1 || orderBy == -1 {
		switch orderField {
		case "":
			fallthrough
		case "Name":
			if orderBy == -1 {
				sort.Slice(dataBase, func(i, j int) bool {
					return dataBase[i].FirstName+dataBase[i].LastName > dataBase[j].FirstName+dataBase[j].LastName
				})
				return dataBase, nil
			} else if orderBy == 1 {
				sort.Slice(dataBase, func(i, j int) bool {
					return dataBase[i].FirstName+dataBase[i].LastName < dataBase[j].FirstName+dataBase[j].LastName
				})
				return dataBase, nil
			}
		case "Id":
			if orderBy == -1 {
				sort.Slice(dataBase, func(i, j int) bool { return dataBase[i].Id > dataBase[j].Id })
				return dataBase, nil
			} else if orderBy == 1 {
				sort.Slice(dataBase, func(i, j int) bool { return dataBase[i].Id < dataBase[j].Id })
				return dataBase, nil
			}
		case "Age":
			if orderBy == -1 {
				sort.Slice(dataBase, func(i, j int) bool { return dataBase[i].Age > dataBase[j].Age })
				return dataBase, nil
			} else if orderBy == 1 {
				sort.Slice(dataBase, func(i, j int) bool { return dataBase[i].Age < dataBase[j].Age })
				return dataBase, nil
			}
		}
		return dataBase, nil
	}

	return dataBase, fmt.Errorf("smth went wrong")

}

//func SearchServer(w http.ResponseWriter, req *http.Request) {
func SearchServer(fileName string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Header["Accesstoken"][0] == "UseOrDie" {
			xmlData, err := os.Open(fileName)
			users := new(XmlUsers)
			defer xmlData.Close()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			bytes, _ := ioutil.ReadAll(xmlData)
			xml.Unmarshal(bytes, users)

			query := req.FormValue("query")
			orderField := req.FormValue("order_field")
			orderBy := req.FormValue("order_by")
			offset := req.FormValue("offset")
			limit := req.FormValue("limit")

			ord, _ := strconv.Atoi(orderBy)
			lim, _ := strconv.Atoi(limit)
			off, _ := strconv.Atoi(offset)

			var xmlSorted []XmlUser
			var answer []toClient

			if orderField == "Name" || orderField == "Id" || orderField == "Age" || orderField == "" {
				xmlSorted, err = users.XmlSorting(orderField, ord, query)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Header().Set("Content-Type", "application/json")
					jsonError, _ := json.Marshal(SearchServerErrorResponse{Error: "BadOrderBy"})
					w.Write(jsonError)
					return
				}

				if len(xmlSorted) < lim {
					lim = len(xmlSorted)
				}
				xmlSorted = xmlSorted[off:lim]
				answerStruct := new(toClient)

				for _, value := range xmlSorted {
					answerStruct.Id = value.Id
					answerStruct.Name = value.FirstName + " " + value.LastName
					answerStruct.Age = value.Age
					answerStruct.About = value.About
					answerStruct.Gender = value.Gender

					answer = append(answer, *answerStruct)
				}

				jsonAnswer, _ := json.Marshal(answer)

				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonAnswer)
			} else {
				w.WriteHeader(http.StatusBadRequest)
				w.Header().Set("Content-Type", "application/json")
				jsonError, _ := json.Marshal(SearchServerErrorResponse{Error: "ErrorBadOrderField"})
				w.Write(jsonError)
				return
			}

		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	})
}

/**
func main() {
	http.HandleFunc("/", SearchServer)
	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
/**/
