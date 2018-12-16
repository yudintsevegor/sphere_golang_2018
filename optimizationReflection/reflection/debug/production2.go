package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Simple struct {
	ID       int
	Username string
	Active   bool
}
type ErrorCase struct {
	Result   interface{}
	JsonData string
}

type IDBlock struct {
	ID int
}

type Complex struct {
	SubSimple  Simple
	ManySimple []Simple
	Blocks     []IDBlock
}

func i2s(data interface{}, out interface{}) error {
	// todo

	outType := reflect.TypeOf(out)
	val := reflect.ValueOf(data)
	outVal := reflect.ValueOf(out)

	fmt.Println("OUT: ",out)
	fmt.Println("OotVal: ",outVal)
	fmt.Println("OutTypeElem: ",outType.Elem())
	fmt.Println("DATA: ", data)
	fmt.Println("value: ", val)
	fmt.Println("type: ", val.Type())
	fmt.Println("kind: ", val.Kind())

	if val.Kind() == reflect.Slice {
		fmt.Println("Slice")
		outVal.Elem().Set(reflect.MakeSlice(outType.Elem(), val.Len(), val.Len()))

		for i := 0; i < val.Len(); i++ {
			elemFromSlice := val.Index(i)

			mapStringInterface := elemFromSlice.Interface().(map[string]interface{})
			lenMap := len(mapStringInterface)

			element := outVal.Elem().Index(i)
			for j := 0; j < lenMap; j++ {
				outField := element.Field(j)
				typeField := element.Type().Field(j)
				k := mapStringInterface[typeField.Name]

				switch t := k.(type) {
				case bool:
					fmt.Println("BOOL: ", reflect.TypeOf(t), t)
					outField.SetBool(t)
				case string:
					fmt.Println("STRING: ", reflect.TypeOf(t), t)
					outField.SetString(t)
				case float64:
					fmt.Println("FLOAT: ", reflect.TypeOf(t), t)
					outField.SetInt(int64(t))
				}

			}
		}

	} else if val.Kind() == reflect.Map {
		fmt.Println("Map")
		outVal = outVal.Elem()
		sliceOfKeys := val.MapKeys()
		mapka := make(map[string]reflect.Value, 1)

		for _, key := range sliceOfKeys {
			for i := 0; i < outVal.NumField(); i++ {
				typeField := outVal.Type().Field(i)
				if key.Interface().(string) == typeField.Name {
					mapka[typeField.Name] = key
				}
			}
		}

		fmt.Println("MAP: ", mapka)

		for i := 0; i < outVal.NumField(); i++ {
			valueField := outVal.Field(i)
			typeField := outVal.Type().Field(i)
			k := val.MapIndex(mapka[typeField.Name])
			fmt.Println("FieldName: ", typeField.Name)
			fmt.Println("Field: ", valueField)

			switch t := k.Interface().(type) {
			case bool:
				fmt.Println("BOOL: ", reflect.TypeOf(t), t)
				valueField.SetBool(t)
			case string:
				fmt.Println("STRING: ", reflect.TypeOf(t), t)
				valueField.SetString(t)
			case float64:
				fmt.Println("INT: ", reflect.TypeOf(t), t)
				valueField.SetInt(int64(t))
			}
		}

	}

/**/
	fmt.Println("OUT: ", out)

	return nil
}

func main() {

	/**
	expected := &Simple{
		ID:       42,
		Username: "rvasily",
		Active:   true,
	}
	result := new(Simple)
	/**/
	smpl := Simple{
		ID:       42,
		Username: "rvasily",
		Active:   true,
	}
	expected := []Simple{smpl, smpl}
	result := []Simple{}

	/**
	smpl := Simple{
		ID:       42,
		Username: "rvasily",
		Active:   true,
	}
	expected := &Complex{
		SubSimple:  smpl,
		ManySimple: []Simple{smpl, smpl},
		Blocks:     []IDBlock{IDBlock{42}, IDBlock{42}},
	}
	result := new(Complex)
	/**/

	jsonRaw, _ := json.Marshal(expected)
	fmt.Println("jsonStr: ", string(jsonRaw))

	var tmpData interface{}
	json.Unmarshal(jsonRaw, &tmpData)

	err := i2s(tmpData, &result)
	//err := i2s(tmpData, result)

	fmt.Println("ERROR: ", err)
	fmt.Println("RESULT: ", result)
	//fmt.Println("DATAStruct: ", tmpData)
	fmt.Println("EXPECTED: ", expected)

	fmt.Println(reflect.DeepEqual(expected, result))

	/*
		// аккуратно в этом тесте
		// писать надо именно в то что пришло
		cases := []ErrorCase{
			// "Active":"DA" - string вместо bool
			ErrorCase{
				&Simple{},
				`{"ID":42,"Username":"rvasily","Active":"DA"}`,
			},
		}
		for _, item := range cases {
			var tmpData interface{}
			fmt.Println(item.JsonData)
			json.Unmarshal([]byte(item.JsonData), &tmpData)
			inType := reflect.ValueOf(item.Result).Type()

			fmt.Println(inType)
		}
	/**/
}
