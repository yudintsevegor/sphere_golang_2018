package main

import (
	"fmt"
	"reflect"
)

func simple(outVal reflect.Value, mapStringInterface map[string]interface{}) error {
	if outVal.Kind() != reflect.Struct {
		return fmt.Errorf("not a struct")
	}

	for i := 0; i < outVal.NumField(); i++ {
		valueField := outVal.Field(i)
		typeField := outVal.Type().Field(i)
		valueFromMap := mapStringInterface[typeField.Name]

		switch switchValue := valueFromMap.(type) {
		case bool:
			if typeField.Type != reflect.TypeOf(switchValue) {
				return fmt.Errorf("bad type of value")
			}
			valueField.SetBool(switchValue)
		case string:
			if typeField.Type != reflect.TypeOf(switchValue) {
				return fmt.Errorf("bad type of value")
			}
			valueField.SetString(switchValue)
		case float64:
			intValue := reflect.TypeOf(int(valueFromMap.(float64)))
			if typeField.Type != intValue {
				return fmt.Errorf("bad type of value")
			}
			valueField.SetInt(int64(switchValue))
		case map[string]interface{}:
			err := simple(valueField, valueFromMap.(map[string]interface{}))
			if err != nil {
				return err
			}
		case []interface{}:
			err := sliceHandler(valueField, typeField.Type, reflect.ValueOf(valueFromMap))
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("Smth went wrong!")
		}
	}
	return nil
}

func sliceHandler(outVal reflect.Value, outType reflect.Type, val reflect.Value) error {
	lenth := val.Len()
	if outVal.Kind() == reflect.Ptr {
		outVal = outVal.Elem()
	}
	if outType.Kind() == reflect.Ptr {
		outType = outType.Elem()
	}
	if outType.Kind() != reflect.Slice {
		return fmt.Errorf("not a slice")
	}
	outVal.Set(reflect.MakeSlice(outType, lenth, lenth))

	for i := 0; i < lenth; i++ {
		elemFromSlice := val.Index(i)
		element := outVal.Index(i)
		mapStringInterface := elemFromSlice.Interface().(map[string]interface{})
		err := simple(element, mapStringInterface)
		if err != nil {
			return err
		}
	}
	return nil
}

func i2s(data interface{}, out interface{}) error {
	outType := reflect.TypeOf(out)
	val := reflect.ValueOf(data)
	outVal := reflect.ValueOf(out)
	if outVal.Kind() != reflect.Ptr {
		return fmt.Errorf("not a ptr type, we cant to return result")
	}

	switch val.Kind() {
	case reflect.Slice:
		err := sliceHandler(outVal, outType, val)
		if err != nil {
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
