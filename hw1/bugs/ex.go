package main

import "fmt"
import "sort"

func main() {
	var m = map[int]string{
	    9:  "Сентябрь",
	    1:  "Январь",
	    2:  "Февраль",
	    10: "Октябрь",
	    5:  "Май",
	    7:  "Июль",
	    8:  "Август",
	    12: "Декарь",
	    3:  "Март",
	    6:  "Июнь",
	    4:  "Апрель",
	    11: "Ноябрь",
	}

	fmt.Println(m)

	var ans  []string
	var arr  []int

	for key, _ := range m {
		arr = append(arr,key)
	}

	sort.Ints(arr)
	fmt.Println(arr)

	for ind, item := range arr {
		fmt.Println(ind)
		fmt.Println(m[item])

		ans = append(ans, m[item])
		//ans[ind] = m[item]
	}
		fmt.Println(ans)


}
