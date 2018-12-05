package main

// тут вы пишете код
// обращаю ваше внимание - в этом задании запрещены глобальные переменные

import (
	"database/sql"
	"fmt"
	"net/http"

	"encoding/json"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Handler struct {
	DB          *sql.DB
	TablesNames []string

	TablesInfo map[string]Table
}

type Table struct {
	Columns     []string
	ColumnsType map[string]reflect.Kind
	CheckNull   map[string]bool
	PrimaryKey  string
}

func resultToUser(StatusCode int, resultsFromServer map[string]interface{}, w http.ResponseWriter) {
	jsonAns, _ := json.Marshal(resultsFromServer)
	w.WriteHeader(StatusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonAns)
}

func requestHandler(r *http.Request, fromBody *map[string]interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &*fromBody)
	if err != nil {
		return err
	}
	return nil
}

func (t Table) idHandler(idColumn *string) {
	for _, val := range t.Columns {
		ok := strings.Contains(val, "id")
		if ok {
			*idColumn = val
			continue
		}
	}
}

func checkId(tableID []interface{}, idFromRequest int) bool {
	check := make(map[int]bool)
	for _, idMap := range tableID {
		for _, id := range idMap.(map[string]interface{}) {
			check[id.(int)] = true
		}
	}
	if _, ok := check[idFromRequest]; ok {
		return true
	}
	return false
}

func (t Table) insertMissingColumns(columnsFromRequest []string, values []interface{}) ([]string, []interface{}) {
	check := make(map[string]bool)
	for _, colR := range columnsFromRequest {
		check[colR] = true
	}
	for _, colT := range t.Columns {
		if _, ok := check[colT]; !ok {
			if !t.CheckNull[colT] {
				columnsFromRequest = append(columnsFromRequest, colT)
				switch t.ColumnsType[colT] {
				case reflect.String:
					values = append(values, "")
				case reflect.Int:
					values = append(values, 0)
				}
			}
		}
	}
	return columnsFromRequest, values
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resultsFromServer := make(map[string]interface{})
	tables := make(map[string][]string)

	mapTable := h.TablesInfo
	tablesNames := h.TablesNames

	fromSplitUrl := strings.Split(r.URL.Path, "/")
	tableNameFromRequest := ""
	IDFromRequest := ""

	for ind, value := range fromSplitUrl {
		if ind == 1 {
			tableNameFromRequest = value
		}
		if ind == 2 {
			IDFromRequest = value
		}
	}

	if tableNameFromRequest == "" {
		tables["tables"] = tablesNames
		resultsFromServer["response"] = tables
		resultToUser(http.StatusOK, resultsFromServer, w)
		return
	}
	if _, ok := mapTable[tableNameFromRequest]; !ok {
		resultsFromServer["error"] = "unknown table"
		resultToUser(http.StatusNotFound, resultsFromServer, w)
		return
	}

	records := make(map[string]interface{})

	switch r.Method {
	case "GET":
		answerFromGET, StatusCode, err := h.Get(tableNameFromRequest, IDFromRequest, r)
		if StatusCode == http.StatusNotFound {
			resultsFromServer["error"] = "record not found"
			resultToUser(StatusCode, resultsFromServer, w)
			return
		}
		if err != nil {
			resultsFromServer["error"] = err
			resultToUser(StatusCode, resultsFromServer, w)
			return
		}
		if IDFromRequest == "" {
			records["records"] = answerFromGET
			resultsFromServer["response"] = records
			resultToUser(StatusCode, resultsFromServer, w)
			return
		}
		records["record"] = answerFromGET[0]
		resultsFromServer["response"] = records
		resultToUser(StatusCode, resultsFromServer, w)
		return
	case "POST":
		records, StatusCode, err := h.Post(tableNameFromRequest, IDFromRequest, r)
		if err == nil {
			resultsFromServer["response"] = records
			resultToUser(StatusCode, resultsFromServer, w)
			return
		}
		resultToUser(StatusCode, records, w)
		return
	case "PUT":
		records, StatusCode, err := h.Put(tableNameFromRequest, r)
		if err != nil {
			resultsFromServer["error"] = records
			resultToUser(StatusCode, resultsFromServer, w)
			return
		}
		resultsFromServer["response"] = records
		resultToUser(StatusCode, resultsFromServer, w)
		return
	case "DELETE":
		records, StatusCode, err := h.Delete(tableNameFromRequest, IDFromRequest, r)
		if err != nil {
			resultsFromServer["error"] = records
			resultToUser(StatusCode, resultsFromServer, w)
			return
		}
		resultsFromServer["response"] = records
		resultToUser(StatusCode, resultsFromServer, w)
		return
	}
}

func (h *Handler) Delete(tableNameFromRequest string, IDFromRequest string, r *http.Request) (map[string]interface{}, int, error) {
	db := h.DB
	var fromBody map[string]interface{}
	answerFromDELETE := make(map[string]interface{})
	table := h.TablesInfo[tableNameFromRequest]

	err := requestHandler(r, &fromBody)
	if err != nil {
		answerFromDELETE["error"] = err
		return answerFromDELETE, http.StatusInternalServerError, err
	}

	res, err := db.Exec("DELETE FROM "+tableNameFromRequest+" WHERE "+table.PrimaryKey+"=?", IDFromRequest)
	if err != nil {
		answerFromDELETE["error"] = err
		return answerFromDELETE, http.StatusInternalServerError, err
	}
	affectedRows, err := res.RowsAffected()
	if err != nil {
		answerFromDELETE["error"] = err
		return answerFromDELETE, http.StatusInternalServerError, err
	}
	if affectedRows == 0 {
		answerFromDELETE["deleted"] = 0
	} else {
		answerFromDELETE["deleted"] = 1
	}

	return answerFromDELETE, http.StatusOK, err
}

func (h *Handler) Post(tableNameFromRequest string, IDFromRequest string, r *http.Request) (map[string]interface{}, int, error) {
	db := h.DB
	table := h.TablesInfo[tableNameFromRequest]

	var fromBody map[string]interface{}
	answerFromPOST := make(map[string]interface{})
	requestToSQL := make([]string, 0, 1)
	valuesToRequest := make([]interface{}, 0, 1)

	err := requestHandler(r, &fromBody)
	if err != nil {
		answerFromPOST["error"] = err
		return answerFromPOST, http.StatusInternalServerError, err
	}

	for col, val := range fromBody {
		if col == table.PrimaryKey {
			answerFromPOST["error"] = "field " + table.PrimaryKey + " have invalid type"
			return answerFromPOST, http.StatusBadRequest, fmt.Errorf("field " + table.PrimaryKey + " have invalid type")
		}

		if table.CheckNull[col] && val == nil {
			req := col + "=" + "NULL"
			requestToSQL = append(requestToSQL, req)
			continue
		} else if !table.CheckNull[col] && val == nil {
			answerFromPOST["error"] = "field " + col + " have invalid type"
			return answerFromPOST, http.StatusBadRequest, fmt.Errorf("field " + col + " have invalid type")
		}

		if table.ColumnsType[col] != reflect.TypeOf(val).Kind() {
			answerFromPOST["error"] = "field " + col + " have invalid type"
			return answerFromPOST, http.StatusBadRequest, fmt.Errorf("field " + col + " have invalid type")
		}

		req := col + "=?"
		requestToSQL = append(requestToSQL, req)
		valuesToRequest = append(valuesToRequest, val)
	}

	requestString := strings.Join(requestToSQL, ", ")
	valuesToRequest = append(valuesToRequest, IDFromRequest)
	_, err = db.Exec("UPDATE "+tableNameFromRequest+" SET "+requestString+" WHERE "+table.PrimaryKey+"=?", valuesToRequest...)
	if err != nil {
		answerFromPOST["error"] = err
		return answerFromPOST, http.StatusInternalServerError, err
	}
	answerFromPOST["updated"] = 1
	return answerFromPOST, http.StatusOK, err
}

func (h *Handler) Put(tableNameFromRequest string, r *http.Request) (map[string]interface{}, int, error) {
	db := h.DB
	table := h.TablesInfo[tableNameFromRequest]

	var fromRequestBody map[string]interface{}
	answerFromPUT := make(map[string]interface{})
	columns := make([]string, 0, 1)
	values := make([]interface{}, 0, 1)

	err := requestHandler(r, &fromRequestBody)
	if err != nil {
		answerFromPUT["error"] = err
		return answerFromPUT, http.StatusInternalServerError, err
	}

	for col, val := range fromRequestBody {
		if table.PrimaryKey == col {
			continue
		}
		if _, ok := table.CheckNull[col]; !ok {
			continue
		}
		columns = append(columns, col)
		values = append(values, val)
	}
	columnsToRequest, valuesToRequest := table.insertMissingColumns(columns, values)

	prepareColumns := strings.Join(columnsToRequest, ", ")
	quesString := strings.Repeat("?, ", len(valuesToRequest)-1) + "?"

	statement, err := db.Prepare("INSERT INTO " + tableNameFromRequest + " ( " + prepareColumns + " ) " + "VALUES (" + quesString + ")")
	if err != nil {
		answerFromPUT["error"] = err
		return answerFromPUT, http.StatusInternalServerError, err
	}
	resultFromExec, err := statement.Exec(valuesToRequest...)
	if err != nil {
		answerFromPUT["error"] = err
		return answerFromPUT, http.StatusInternalServerError, err
	}

	insertedId, err := resultFromExec.LastInsertId()
	if err != nil {
		answerFromPUT["error"] = err
		return answerFromPUT, http.StatusInternalServerError, err
	}
	answerFromPUT[table.PrimaryKey] = insertedId

	return answerFromPUT, http.StatusOK, nil
}

func (h *Handler) Get(tableNameFromRequest string, IDFromRequest string, r *http.Request) ([]interface{}, int, error) {
	db := h.DB
	table := h.TablesInfo[tableNameFromRequest]

	answerFromGET := make([]interface{}, 0, 1)
	offset := "0"
	limit := "5"

	ok := strings.Contains(r.URL.RawQuery, "limit")
	if ok {
		limit = r.FormValue("limit")
	}

	ok = strings.Contains(r.URL.RawQuery, "offset")
	if ok {
		offset = r.FormValue("offset")
	}

	intOffset, err := strconv.Atoi(offset)
	if err != nil {
		intOffset = 0
	}
	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		intLimit = 5
	}

	var row *sql.Rows
	switch IDFromRequest {
	case "":
		row, err = db.Query("SELECT * FROM "+tableNameFromRequest+" LIMIT ? OFFSET ?", intLimit, intOffset)
		defer row.Close()
		if err != nil {
			return answerFromGET, http.StatusInternalServerError, err
		}

	default:
		row, err = db.Query("SELECT * FROM "+tableNameFromRequest+" WHERE "+table.PrimaryKey+"=?", IDFromRequest)
		defer row.Close()
		if err != nil {
			return answerFromGET, http.StatusInternalServerError, err
		}
	}

	answerFromGET, err = getHandler(answerFromGET, table, row)
	if err != nil {
		return answerFromGET, http.StatusInternalServerError, err
	}
	if len(answerFromGET) == 0 {
		return answerFromGET, http.StatusNotFound, fmt.Errorf("record not found")
	}

	return answerFromGET, http.StatusOK, nil
}

func getHandler(answerFromGET []interface{}, table Table, row *sql.Rows) ([]interface{}, error) {
	numbersOfColumns := len(table.Columns)
	vals := make([]interface{}, numbersOfColumns)
	rawResult := make([]sql.RawBytes, numbersOfColumns)

	for i, _ := range vals {
		vals[i] = &rawResult[i]
	}

	for row.Next() {
		err := row.Scan(vals...)
		if err != nil {
			return answerFromGET, err
		}
		result := make(map[string]interface{})
		for i, raw := range rawResult {
			if string(raw) == "" {
				if table.CheckNull[table.Columns[i]] == true {
					result[table.Columns[i]] = nil
				} else {
					result[table.Columns[i]] = ""
				}
			} else {
				if i == 0 {
					toData, err := strconv.Atoi(string(raw))
					if err != nil {
						return answerFromGET, err
					}
					result[table.Columns[i]] = toData
				} else {
					result[table.Columns[i]] = string(raw)
				}
			}
		}
		answerFromGET = append(answerFromGET, result)
	}
	return answerFromGET, nil
}

func NewDbExplorer(db *sql.DB) (http.Handler, error) {
	tablesNames := make([]string, 0, 1)
	mapCheckNull := make(map[string]bool)
	mapColumnsType := make(map[string]reflect.Kind)
	mapTable := make(map[string]Table)

	handler := &Handler{}
	tableStruct := Table{}

	var tableName string
	var idColumn string

	rowsTb, err := db.Query("SHOW TABLES")
	for rowsTb.Next() {
		err = rowsTb.Scan(&tableName)
		if err != nil {
			return handler, nil
		}
		tablesNames = append(tablesNames, tableName)
	}
	rowsTb.Close()

	for _, tableName := range tablesNames {
		rows, err := db.Query("SELECT * FROM " + tableName)
		if err != nil {
			return handler, fmt.Errorf("InternalDBError")
		}
		tableStruct.Columns, _ = rows.Columns()
		types, err := rows.ColumnTypes()
		if err != nil {
			return handler, fmt.Errorf("InternalDBError")
		}
		for ind, val := range types {
			null, _ := val.Nullable()
			mapCheckNull[tableStruct.Columns[ind]] = null
			switch val.DatabaseTypeName() {
			case "VARCHAR", "TEXT":
				mapColumnsType[tableStruct.Columns[ind]] = reflect.String
			case "INT":
				mapColumnsType[tableStruct.Columns[ind]] = reflect.Int
			}
		}

		tableStruct.ColumnsType = mapColumnsType
		tableStruct.CheckNull = mapCheckNull
		tableStruct.idHandler(&idColumn)
		tableStruct.PrimaryKey = idColumn
		mapTable[tableName] = tableStruct
		rows.Close()
	}
	handler = &Handler{DB: db, TablesNames: tablesNames, TablesInfo: mapTable}

	return handler, nil
}
