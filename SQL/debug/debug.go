/*
	q := []string{`CREATE TABLE items (
  	id int(11) NOT NULL AUTO_INCREMENT,
  	title varchar(255) NOT NULL,
  	description text NOT NULL,
  	updated varchar(255) DEFAULT NULL,
  	PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;`,

		`INSERT INTO items (id, title, description, updated) VALUES
(1,	'database/sql',	'Рассказать про базы данных',	'rvasily'),
(2,	'memcache',	'Рассказать про мемкеш с примером использования',	NULL);`}
	for _, val := range q{
		 db.Exec(val)
	}
	if err != nil {
		panic(err)
	}
	*/

/*switch reflect.TypeOf(val).Kind() {
		case reflect.Float64:
			strVal := strconv.Itoa(int(val.(float64)))
			valuesToRequest = append(valuesToRequest, strVal)
			continue
		case reflect.String:
			valuesToRequest = append(valuesToRequest, val.(string))
			continue
		case reflect.Int:
			strVal := strconv.Itoa(val.(int))
			valuesToRequest = append(valuesToRequest, strVal)
			continue
		}*/
package main

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные

import (
	"database/sql"
	"fmt"
	"strconv"
	//"bytes"
	"encoding/json"
	//"io/ioutil"
	//"reflect"
	"net/http"
	//"net/url"
	//"time"
	//"context"
	_ "github.com/go-sql-driver/mysql"
	"regexp"
)

type Handler struct {
	DB *sql.DB
}

type Table struct{
	Length int
	Columns []string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db := h.DB
	mapTable := map[string]Table{}

	//tableNames := map[string][]string{}
	//tableLength := map[string]int{}
	tables := map[string][]string{}
	names := make([]string, 0, 1)
	mapka := make(map[string]interface{})

	rowsTb, err := db.Query("SHOW TABLES")
	defer rowsTb.Close()
	var name string
	var leng int
	for rowsTb.Next() {
		err = rowsTb.Scan(&name)
		if err != nil {
			panic(err)
		}
		/**
		rows, err := db.Query("SELECT * FROM " + name)
		rows.Close()
		if err != nil {
			panic(err)
		}
		tableNames[name], _ = rows.Columns()
		/**/
		names = append(names, name)
	}
	/**/
	table := Table{}
	for _, name := range names {
		rows, err := db.Query("SELECT * FROM " + name)
		if err != nil {
			panic(err)
		}
		//tableNames[name], _ = rows.Columns()
		table.Columns, _ = rows.Columns()
		rows.Close()

		rowsLen, err := db.Query("SELECT COUNT(1) FROM " + name)
		if err != nil {
			panic(err)
		}
		for rowsLen.Next() {
			err = rowsLen.Scan(&leng)
			if err != nil {
				panic(err)
			}
			//tableLength[name] = leng
			table.Length = leng
			mapTable[name] = table
		}
		rowsLen.Close()
	}
	//fmt.Println(mapTable)
	/**/

	reg := regexp.MustCompile("/([^/?]*)")
	allSubmatches := reg.FindAllStringSubmatch(r.URL.Path, -1)
	tableFromRequest := ""
	IDFromRequest := ""
	for innInd, innSlice := range allSubmatches {
		for ind, value := range innSlice {
			if innInd == 0 && ind == 1 {
				tableFromRequest = value
			}
			if innInd == 1 && ind == 1 {
				IDFromRequest = value
			}

		}
	}
	/**/
	if tableFromRequest == "" {
		tables["tables"] = names
		mapka["response"] = tables
		jsonAns, _ := json.Marshal(mapka)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonAns)
		return
	}
	if _, ok := mapTable[tableFromRequest]; !ok {
	//if _, ok := tableNames[tableFromRequest]; !ok {
		mapka["error"] = "unknown table"
		jsonAns, _ := json.Marshal(mapka)
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonAns)
		return
	}

	fmt.Println(tableFromRequest)
	fmt.Println(IDFromRequest)
	records := make(map[string]interface{})
	sliceRecords := make([]interface{}, 0, 1)

	for ind := 1; ind <= mapTable[tableFromRequest].Length; ind++ {
		row, err := db.Query("SELECT * FROM " + tableFromRequest + " WHERE " + mapTable[tableFromRequest].Columns[0] + "=?", ind)
		defer row.Close()
		fmt.Println(err)

		length := len(mapTable[tableFromRequest].Columns)
		vals := make([]interface{}, length)
		rawResult := make([]sql.RawBytes, length)

		for i, _ := range vals {
			vals[i] = &rawResult[i]
		}
		for row.Next() {
			errr := row.Scan(vals...)
			fmt.Println(errr)
			data := make(map[string]interface{})
			for i, raw := range rawResult {
				if string(raw) == "" {
					data[mapTable[tableFromRequest].Columns[i]] = nil
				} else {
					if i == 0 {
						toData, _ :=  strconv.Atoi(string(raw))
						data[mapTable[tableFromRequest].Columns[i]] = toData

					} else  {
						data[mapTable[tableFromRequest].Columns[i]] = string(raw)
					}
				}
			}
			sliceRecords = append(sliceRecords, data)
		}
	}
	//fmt.Println(sliceRecords)

	records["records"] = sliceRecords
	mapka["response"] = records

	jsonAns, _ := json.Marshal(mapka)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonAns)
}

func NewDbExplorer(db *sql.DB) (http.Handler, error) {
	handler := &Handler{DB: db}
	return handler, nil
}
