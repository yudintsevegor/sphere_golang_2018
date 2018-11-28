package main

// сюда писать код
import (
	"strconv"
	"sort"
	"fmt"
	"strings"
	"time"
	"sync"
)

func SingleHash(in, out chan interface{}) {
//func SingleHash(in, out chan interface{}) chan interface{}{
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	for value := range in {
		res_1 := make(chan interface{})
		res_2 := make(chan interface{})
		wg.Add(3)
		itoa := strconv.Itoa(value.(int))

	go func(in chan interface{}){
		defer wg.Done()

		go func(itoa string) {
			defer wg.Done()
			res_1 <- DataSignerCrc32(itoa)
		}(itoa)

		go func(itoa string) {
			defer wg.Done()
			mu.Lock()
			md5 := DataSignerMd5(itoa)
			mu.Unlock()

			res_2 <- DataSignerCrc32(md5)
		}(itoa)

		//val_1 := <-res_1
		//val_2 := <-res_2

		//fmt.Println(val_1, "~", val_2)
		//out <- val_1.(string) + "~" + val_2.(string)
		out <- (<-res_1).(string) + "~" + (<-res_2).(string)
	}(in)

	}
	wg.Wait()
	//fmt.Println(len(out))
	/*
	for value := range out{
		fmt.Println(value)
	}*/
	//return out

}

//func MultiHash(in, out chan interface{}) chan interface{}{
func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	th := []int{0, 1, 2, 3, 4, 5}
	hash := make(map[string]string)
	ind := 0
	//lenth := len(in)
	for value := range in {
		for _, element := range th {
			wg.Add(1)
			go func(ind int, element int, value interface{}) {
				defer  wg.Done()
				itoa := strconv.Itoa(element)
				crc := DataSignerCrc32(itoa + value.(string))
				mu.Lock()
				hash[strconv.Itoa(ind) + strconv.Itoa(element)] = crc
				mu.Unlock()
			}(ind, element, value)
		}
		ind++
	}
	wg.Wait()

	arr := []string{}
	for key, _ := range hash{
		arr = append(arr, key)
	}

	sort.Strings(arr)

	for i := 0; i < ind; i++{
		data := []string{}
		for j := i*len(th); j < (i + 1)*(len(th)); j++{
			data =  append(data, hash[arr[j]])
		}
		answer := strings.Join(data, "")
		out <- answer
	}

	//return out
}


func CombineResults(in, out chan interface{}){

	data := []string{}
	for value := range in {
		data = append(data, value.(string))
	}

	sort.Strings(data)
	answer := strings.Join(data, "_")
	out <- answer
	//fmt.Println(answer)

}

func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}

	chanIn := make(chan interface{})

	for _, yX := range jobs {
		chanOut := make(chan interface{})
		//chan_1 := make(chan interface{})
		//chan_2 := make(chan interface{})
		wg.Add(1)
		go func(inPut chan interface{}, outPut chan interface{}, sub job) {
			defer wg.Done()
			defer close(outPut)
			sub(inPut, outPut)

		}(chanIn, chanOut, yX)

		chanIn = chanOut

	}
	wg.Wait()

}

func main() {
	start := time.Now()
	/*data := []int{0,1}
	data := []int{0, 1, 1, 2, 3, 5, 8}
	/**/
	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			//inputData := []int{0,1}
			inputData := []int{0, 1, 1, 2, 3, 5, 8}
			for _, fibNum := range inputData {
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			dataRaw := <-in
			data, ok := dataRaw.(string)
			if !ok {
				fmt.Println("cant convert result data to string")
			} else {
				fmt.Println(data)
			}
		}),

	}

	ExecutePipeline(hashSignJobs...)
	/*
	for value := range hashSignJobs{
		ExecutePipeline(value)
	}
	/*
	ch_data := make(chan interface{}, len(data))
	ch_ans := make(chan interface{}, len(data))
	out_1 := make(chan interface{}, len(data))
	out_2 := make(chan interface{}, len(data))
	//out_3 := make(chan string, len(data))
	for _, fibNum := range data {
		ch_data <- fibNum
	}
	close(ch_data)
	out_1 = SingleHash(ch_data, ch_ans)
	////close(out_1)
	out_2 = MultiHash(out_1, ch_ans)
	//close(out_2)
	CombineResults(out_2, ch_ans)
	/**/
	end := time.Since(start)
	fmt.Println(end)
}


