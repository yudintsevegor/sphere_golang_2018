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
/**/
func simple(outVal reflect.Value, mapStringInterface map[string]interface{}, val reflect.Value) error{
	fmt.Println("SIMPLE")
	for i := 0; i < outVal.NumField(); i++ {
		valueField := outVal.Field(i)
		typeField := outVal.Type().Field(i)
		k := mapStringInterface[typeField.Name]

		fmt.Println("valueField: ", valueField)
		fmt.Println(typeField.Name)
		fmt.Println("K: ", k)
		fmt.Println(reflect.TypeOf(k))

		switch t := k.(type) {
		case bool:
			fmt.Println("BOOL: ", reflect.TypeOf(t), t)
			if typeField.Type != reflect.TypeOf(t) {
				return fmt.Errorf("not bool")
			}
			valueField.SetBool(t)
		case string:
			fmt.Println("STRING: ", reflect.TypeOf(t), t)
			if typeField.Type != reflect.TypeOf(t) {
				return fmt.Errorf("not string")
			}
			valueField.SetString(t)
		case float64:
			fmt.Println("INT: ", reflect.TypeOf(t), t)
			//NEED TO FIXING
			//if typeField.Type != reflect.TypeOf(t) {
			//	return fmt.Errorf("not int")
			//}
			valueField.SetInt(int64(t))
		case map[string]interface{}:
			fmt.Println("refl type: ", reflect.TypeOf(t))
			err := simple(valueField, k.(map[string]interface{}), val)
			fmt.Println(err)
		case []interface{}:
			fmt.Println("You got a Slice")
			/**
			fmt.Println("valueField: ", valueField)
			//val := reflect.ValueOf(valueField)
			fmt.Println("outVal: ", outVal.Field(i))
			outType := reflect.TypeOf(outVal)
			fmt.Println("outType: ", outType)
			err := sliceHandler(valueField, outType, val)
			fmt.Println(err)
			/**/
			err := i2s(k, valueField)
			fmt.Println(err)

		default:
			return fmt.Errorf("Smth went wrong!")
		}
	}
	return nil
}

func sliceHandler(outVal reflect.Value, outType reflect.Type, val reflect.Value) error{
	lenth := val.Len()
	//outVal.Elem().Set(reflect.MakeSlice(outType.Elem(), lenth, lenth))
	if outVal.Kind() == reflect.Ptr{
		outVal = outVal.Elem()
	}
	fmt.Println(outVal, outType, val)

	outVal.Set(reflect.MakeSlice(outType.Elem(), lenth, lenth))
	for i := 0; i < lenth; i++ {
			elemFromSlice := val.Index(i)

			element := outVal.Index(i)
			mapStringInterface := elemFromSlice.Interface().(map[string]interface{})
			err := simple(element, mapStringInterface, val)
			if err != nil{
				return err
			}
		}

	return nil
}
/**/
func i2s(data interface{}, out interface{}) error {
	// todo
	outType := reflect.TypeOf(out)
	val := reflect.ValueOf(data)
	refl := reflect.TypeOf(data)
	outVal := reflect.ValueOf(out)

	fmt.Println("OutVal: ", outVal)
	fmt.Println("OutType: ", outType)
	fmt.Println("value: ", val)
	fmt.Println("valueType: ", refl)

	fmt.Println("OUT: ", out)
	fmt.Println("type: ", reflect.TypeOf(outVal))
	fmt.Println("OutTypeElem: ", outType.Elem())
	//fmt.Println("DATA: ", data)
	fmt.Println("deeper: ", reflect.TypeOf(out))
	fmt.Println("type val: ", reflect.ValueOf(val).Kind())
	fmt.Println("kind val: ", val.Kind())
/**	mapStringInterface := val.Interface().(map[string]interface{})
	fmt.Println(mapStringInterface["SubSimple"])
	fmt.Println(reflect.ValueOf(mapStringInterface["ManySimple"]).Kind())
	fmt.Println("type map: ",reflect.ValueOf(mapStringInterface).Kind())
	t1 := outVal.Elem().Field(0)
	t2 := outVal.Elem().Field(1)
	t3 := outVal.Elem().Field(2)
	fmt.Println(t1, t2, t3)

	f1 := outVal.Elem().Type().Field(0).Name
	f2 := outVal.Elem().Type().Field(1).Name
	f3 := outVal.Elem().Type().Field(2).Name
	fmt.Println(f1, f2, f3)

	fmt.Println("type f1: ", reflect.ValueOf(mapStringInterface[f1]).Kind())
	fmt.Println("type f2: ", reflect.ValueOf(mapStringInterface[f2]).Kind())
	fmt.Println("type f3: ",reflect.ValueOf(mapStringInterface[f3]).Kind())

	mapka := mapStringInterface[f1].(map[string]interface{})
	fmt.Println(mapka["ID"])
	fmt.Println(reflect.ValueOf(mapka["ID"]).Kind())
/**/
	switch val.Kind(){
	case reflect.Slice:
		fmt.Println("Slice")
		err := sliceHandler(outVal, outType, val)
		if err != nil{
			return err
		}
	case reflect.Map:
		fmt.Println("Map")
		mapStringInterface := val.Interface().(map[string]interface{})
		lenth := val.Len()
		fmt.Println("Lenth: ",lenth)
		//for i := 0; i < lenth; i++{
		//	fmt.Println("I: ", i)
			//key := mapStringInterface[outVal.Elem().Type().Field(i).Name]
			//refl := reflect.ValueOf(key).Kind()
			//fmt.Println("KEY: ", key)
			//fmt.Println("REFL: ", refl)
			//if refl == reflect.Slice{
				/*err := sliceHandle()
				if err != nil{
					return err
				}*/
			//}
			outVal = outVal.Elem()
			err := simple(outVal, mapStringInterface, val)
			if err != nil{
				return err
			}
		//}
		default:
		return fmt.Errorf("haha, harcode!")
	}

	fmt.Println("OUT: ", out)
	return nil
}

func main() {
	/**
	js := `{"ID":"42","Username":"rvasily","Active":true}`
	result := &Simple{}
	/**/
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
	//json.Unmarshal([]byte(js), &tmpData)

	err := i2s(tmpData, &result)
	//err := i2s(tmpData, result)

	fmt.Println("ERROR: ", err)
	fmt.Println("RESULT: ", result)
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
