//package main

// это программа для которой ваш кодогенератор будет писать код
// запускать через go test -v, как обычно

// этот код закомментирован чтобы он не светился в тестовом покрытии

import (
	//"fmt"
	"context"
	"encoding/json"
	"net/http"
	//"reflect"
	"strconv"
	//"errors"
)

func (srv *OtherApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	switch r.URL.Path {
	case "/user/create":
		srv.wrapperCreateOtherApi(ctx, w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		answer := map[string]interface{}{
			"error": "unknown method",
		}
		jsonRes, _ := json.Marshal(answer)
		w.Write(jsonRes)
	}
}

func (srv *OtherApi) wrapperCreateOtherApi(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	answer := make(map[string]interface{})
	check := map[string]string{
		"warrior":  "",
		"sorcerer": "",
		"rouge":    "",
	}
	if r.Method == http.MethodPost {
		var auth string
		if len(r.Header["X-Auth"]) != 0 {
			auth = r.Header["X-Auth"][0]
		}

		if auth != "100500" {
			w.WriteHeader(http.StatusForbidden)
			answer = map[string]interface{}{
				"error": "unauthorized",
			}
		//} else if errAtoi != nil {
		//	w.WriteHeader(http.StatusBadRequest)
		//	answer = map[string]interface{}{
		//		"error": "age must be int",
		//	}
		} else {
			level, errAtoi := strconv.Atoi(r.FormValue("level"))
			params := OtherCreateParams{
				Username: r.FormValue("username"),
				Name:     r.FormValue("account_name"),
				Class:    r.FormValue("class"),
				Level:    level,
			}
			if params.Class == "" {
				params.Class = "warrior"
			}
			_, ok := check[params.Class]
			switch {
			//if params.Username == ""{
			case errAtoi != nil:
				w.WriteHeader(http.StatusBadRequest)
				answer = map[string]interface{}{
					"error": "age must be int",
				}
			case params.Username == "":
				w.WriteHeader(http.StatusBadRequest)
				answer = map[string]interface{}{
					"error": "username must be not empty",
				}
			//} else if len(params.Username) < 3 {
			case len(params.Username) < 3:
				w.WriteHeader(http.StatusBadRequest)
				answer = map[string]interface{}{
					"error": "username len must be >= 3",
				}
			//} else if params.Level < 1 {
			case params.Level < 1:
				w.WriteHeader(http.StatusBadRequest)
				answer = map[string]interface{}{
					"error": "level must be >= 1",
				}
			//} else if params.Level > 50{
			case params.Level > 50:
				w.WriteHeader(http.StatusBadRequest)
				answer = map[string]interface{}{
					"error": "level must be <= 50",
				}
			//} else if _, ok := check[params.Class]; !ok{
			case !ok:
				w.WriteHeader(http.StatusBadRequest)
				answer = map[string]interface{}{
					"error": "class must be one of [warrior, sorcerer, rouge]",
				}
			//} else {
			default:
				res, err := srv.Create(ctx, params)
				if err == nil {
					answer = map[string]interface{}{
						"error":    "",
						"response": res,
					}
				}
			}
		}
	} else {
		w.WriteHeader(http.StatusNotAcceptable)
		answer = map[string]interface{}{
			"error": "bad method",
		}
	}
	jsonRes, _ := json.Marshal(answer)
	w.Write(jsonRes)
}

/**/
func (srv *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	switch r.URL.Path {
	case "/user/profile":
		srv.wrapperProfileMyApi(ctx, w, r)
	case "/user/create":
		srv.wrapperCreateMyApi(ctx, w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		answer := map[string]interface{}{
			"error": "unknown method",
		}
		jsonRes, _ := json.Marshal(answer)
		w.Write(jsonRes)
	}
}

func (srv *MyApi) wrapperCreateMyApi(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	answer := make(map[string]interface{})
	var ansError string
	check := map[string]string{
		"admin":     "",
		"moderator": "",
		"user":      "",
	}
	if r.Method == http.MethodPost {
		var auth string
		if len(r.Header["X-Auth"]) != 0 {
			auth = r.Header["X-Auth"][0]
		}
		if auth != "100500" {
			w.WriteHeader(http.StatusForbidden)
			//answer = map[string]interface{}{
			//	"error": "unauthorized",
			//}
			ansError =  "unauthorized"
		//} else if err != nil {
		//	w.WriteHeader(http.StatusBadRequest)
		//	answer = map[string]interface{}{
		//		"error": "age must be int",
		//	}
		} else {
			age, errAtoi := strconv.Atoi(r.FormValue("age"))
			params := CreateParams{
				Login:  r.FormValue("login"),
				Name:   r.FormValue("full_name"),
				Status: r.FormValue("status"),
				Age:    age,
			}
			if params.Status == "" {
				params.Status = "user"
			}
			_, ok := check[params.Status]
			switch {
			//if params.Login == ""{
			case errAtoi != nil:
				w.WriteHeader(http.StatusBadRequest)
			//	answer = map[string]interface{}{
			//		"error": "age must be int",
			//	}
				ansError =  "age must be int"
			case params.Login == "":
				w.WriteHeader(http.StatusBadRequest)
			//	answer = map[string]interface{}{
			//		"error": "login must be not empty",
				//}
				ansError ="login must be not empty"
			//} else if len(params.Login) < 10 {
			case len(params.Login) < 10:
				w.WriteHeader(http.StatusBadRequest)
				//answer = map[string]interface{}{
				//	"error": "login len must be >= 10",
				//}
				ansError = "login len must be >= 10"
			//} else if params.Age < 0 {
			case params.Age < 0:
				w.WriteHeader(http.StatusBadRequest)
			//	answer = map[string]interface{}{
			//		"error": "age must be >= 0",
				//}
					ansError = "age must be >= 0"
			//} else if params.Age > 128{
			case params.Age > 128:
				w.WriteHeader(http.StatusBadRequest)
				//answer = map[string]interface{}{
				//	"error": "age must be <= 128",
				//}
					ansError = "age must be <= 128"
			//} else if _, ok := check[params.Status]; !ok{
			case !ok:
				w.WriteHeader(http.StatusBadRequest)
				//answer = map[string]interface{}{
				//	"error": "status must be one of [user, moderator, admin]",
			//	}
					ansError ="status must be one of [user, moderator, admin]"
			//} else {
			default:
				res, err := srv.Create(ctx, params)
				if err != nil {
					switch err.Error() {
					case "bad user":
						w.WriteHeader(http.StatusInternalServerError)
					default:
						w.WriteHeader(http.StatusConflict)
					}
					//answer = map[string]interface{}{
					//	"error": err.Error(),
					//}
					ansError = err.Error()
				} else {
					answer = map[string]interface{}{
						"error":    "",
						"response": res,
					}
				}
			}
		}
	} else {
		w.WriteHeader(http.StatusNotAcceptable)
		//answer = map[string]interface{}{
		//	"error": "bad method",
		//}
			ansError = "bad method"
	}

	if ansError != ""{
		answer = map[string]interface{}{
				"error":    ansError,
		}
	}// else {
	//	answer = map[string]interface{}{
	//			"error":    "",
	//			"response": res,
	//	}
	//}

	jsonRes, _ := json.Marshal(answer)
	w.Write(jsonRes)
}

func (srv *MyApi) wrapperProfileMyApi(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	params := ProfileParams{
		Login: r.FormValue("login"),
	}
	answer := make(map[string]interface{})
	if params.Login != "" {
		res, err := srv.Profile(ctx, params)
		if err == nil {
			answer = map[string]interface{}{
				"error":    "",
				"response": res,
			}
		} else {
			switch err.Error() {
			case "bad user":
				w.WriteHeader(http.StatusInternalServerError)
				answer = map[string]interface{}{
					"error": "bad user",
				}
			case "user not exist":
				w.WriteHeader(http.StatusNotFound)
				answer = map[string]interface{}{
					"error": "user not exist",
				}
			}
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		err := "login must be not empty"
		answer = map[string]interface{}{
			"error": err,
		}
	}
	jsonRes, _ := json.Marshal(answer)
	w.Write(jsonRes)

}

/**
func main() {
	// будет вызван метод ServeHTTP у структуры MyApi
	http.Handle("/user/", NewMyApi())
	//http.Handle("/user/", NewOtherApi())

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
/**/
