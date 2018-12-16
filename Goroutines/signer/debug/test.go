package main

// сюда писать код
import (
	//"runtime"
	"strconv"
	//"sort"
	"fmt"
	//"strings"
	"time"
)

/**/
func SingleHash(input int) string{

	res_1 := DataSignerCrc32(strconv.Itoa(input))
	res_2 := DataSignerCrc32(DataSignerMd5(strconv.Itoa(input)))
	return res_1 + "~" + res_2
}

func MultiHash(data string) string{

	var answer string
	var th = []int{0, 1, 2, 3, 4, 5}
	for _, value := range th {
			answer += DataSignerCrc32(strconv.Itoa(value) + data)
	}
	return answer
}
/*
func CombineResults(data []string) {

	sort.Strings(data)
	answer := strings.Join(data, "_")

	//return answer
}
*/

func main() {
	data := []int{0,1}
	start := time.Now()
	//data := []int{0, 1, 1, 2, 3, 5, 8}
	var slc []string
	for _, fibNum := range data {
		val := SingleHash(fibNum)
		fmt.Println(val)
		slc = append(slc, val)
	}
	for _, val := range slc{
		res := MultiHash(val)
		fmt.Println(res)
	}
	//CombineResults

	end := time.Since(start)
	fmt.Println(end)

}


