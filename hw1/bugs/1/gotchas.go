package main

import "strconv"
import "sort"
/*
	сюда вам надо писать функции, которых не хватает, чтобы проходили тесты в gotchas_test.go

	IntSliceToString
	MergeSlices
	GetMapValuesSortedByKey
*/

func IntSliceToString(arr []int) string{

	var  str string
	for  _, item := range arr {
		str += strconv.Itoa(item)
	}
	return  str
}


func MergeSlices(arr_1 []float32, arr_2 []int32) []int{


	var arr []int

	for _, item := range arr_1 {
		arr = append(arr, int(item))
	}

	for _, item := range arr_2 {
		arr = append( arr, int(item))
	}

	return arr

}


func GetMapValuesSortedByKey(m map[int]string) []string{

	var ans  []string
	var arr  []int

	for key, _ := range m {
		arr = append(arr,key)
	}

	sort.Ints(arr)

	for _, item := range arr {
		ans = append(ans,m[item])
	}


	return ans
}




