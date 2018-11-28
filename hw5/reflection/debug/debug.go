package main

import(
	"fmt"
	"reflect"
	"encoding/json"
)
func simple(outVal reflect.Value, mapStringInterface map[string]interface{}) error{
	if outVal.Kind() != reflect.Struct{
		return fmt.Errorf("not a struct")
	}

	for i := 0; i < outVal.NumField(); i++ {
		valueField := outVal.Field(i)
		typeField := outVal.Type().Field(i)
		k := mapStringInterface[typeField.Name]

		switch t := k.(type) {
		case bool:
			if typeField.Type != reflect.TypeOf(t) {
				return fmt.Errorf("bad type of value")
			}
			valueField.SetBool(t)
		case string:
			if typeField.Type != reflect.TypeOf(t) {
				return fmt.Errorf("bad type of value")
			}
			valueField.SetString(t)
		case float64:
			intValue := reflect.TypeOf(int(k.(float64)))
			if typeField.Type != intValue {
				return fmt.Errorf("bad type of value")
			}
			valueField.SetInt(int64(t))
		case map[string]interface{}:
			err := simple(valueField, k.(map[string]interface{}))
			if err != nil {
				return err
			}
		case []interface{}:
			err := sliceHandler(valueField, typeField.Type, reflect.ValueOf(k))
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("Smth went wrong!")
		}
	}
	return nil
}

func sliceHandler(outVal reflect.Value, outType reflect.Type, val reflect.Value) error{
	lenth := val.Len()
	if outVal.Kind() == reflect.Ptr{
		outVal = outVal.Elem()
	}
	if outType.Kind() == reflect.Ptr{
		outType = outType.Elem()
	}
	if outType.Kind() != reflect.Slice{
		return fmt.Errorf("not a slice")
	}
	outVal.Set(reflect.MakeSlice(outType, lenth, lenth))

	for i := 0; i < lenth; i++ {
			elemFromSlice := val.Index(i)
			element := outVal.Index(i)
			mapStringInterface := elemFromSlice.Interface().(map[string]interface{})
			err := simple(element, mapStringInterface)
			if err != nil{
				return err
			}
		}
	return nil
}

func i2s(data interface{}, out interface{}) error {
	// todo
	outType := reflect.TypeOf(out)
	val := reflect.ValueOf(data)
	outVal := reflect.ValueOf(out)
	if outVal.Kind() != reflect.Ptr{
		return fmt.Errorf("not a ptr type, we cant to return result")
	}

	switch val.Kind(){
	case reflect.Slice:
		err := sliceHandler(outVal, outType, val)
		if err != nil{
			return err
		}
	case reflect.Map:
		mapStringInterface := val.Interface().(map[string]interface{})
		outVal = outVal.Elem()
		err := simple(outVal, mapStringInterface)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("haha, classic!")
	}
	return nil
}

/**
type IDBlock struct {
	ID int
}

type Complex struct {
	SubSimple  Simple
	ManySimple []Simple
	Blocks     []IDBlock
}
type Simple struct {
	ID       int
	Username string
	Active   bool
}
/**/

func main() {
	//jsonData := `{"ID":42,"Username":"rvasily","Active":"DA"}`
	//jsonData := `{"ID":"42","Username":"rvasily","Active":true}`
	//jsonData := `{"ID":42,"Username":100500,"Active":true}`
	//result := &Simple{}
	//jsonData := `{"SubSimple":true,"ManySimple":[{"ID":42,"Username":"rvasily","Active":true}]}`
	//jsonData := `[{"ID":42,"Username":"rvasily","Active":true}]`
	//jsonData := `{"SubSimple":{"ID":42,"Username":"rvasily","Active":true},"ManySimple":{}}`
	//result := &Complex{}
	//jsonData := `{"ID":42,"Username":"rvasily","Active":true}`
	//result := Simple{}
	/**
	smpl := Simple{
		ID:       42,
		Username: "rvasily",
		Active:   true,
	}
	expected := &Complex{
		SubSimple:  smpl,
		ManySimple: []Simple{ smpl, smpl},
		Blocks:     []IDBlock{IDBlock{42}, IDBlock{42}},
	}

	result := new(Complex)
	/**/
	//jsonRaw, _ := json.Marshal(expected)
	//fmt.Println("jsonStr: ", string(jsonRaw))

	var tmpData interface{}
	//json.Unmarshal(jsonRaw, &tmpData)
	json.Unmarshal([]byte(jsonData), &tmpData)
	inType := reflect.ValueOf(result).Type()
	i2s(tmpData, result)
	outType := reflect.ValueOf(result).Type()
/*
	if err == nil {
		fmt.Println("ERROR MUST BE HERE: ", err)
	} else {
		fmt.Println("ERROR: ", err)
	}
*/
	if inType != outType {
		fmt.Errorf("results type not match\nGot:\n%#v\nExpected:\n%#v", outType, inType)
	}

	/*
	fmt.Println("RESULT: ", result)
	fmt.Println("EXPECTED: ", expected)
	fmt.Println(reflect.DeepEqual(expected, result))
	*/

}

