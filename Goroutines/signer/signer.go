package main

// сюда писать код
import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	for value := range in {
		wg.Add(1)
		itoa := strconv.Itoa(value.(int))

		go func(in chan interface{}) {
			defer wg.Done()
			resultCrc := make(chan interface{})
			resultCrcMd5 := make(chan interface{})

			go func(itoa string) {
				resultCrc <- DataSignerCrc32(itoa)
			}(itoa)

			go func(itoa string) {
				mu.Lock()
				md5 := DataSignerMd5(itoa)
				mu.Unlock()

				resultCrcMd5 <- DataSignerCrc32(md5)
			}(itoa)

			out <- (<-resultCrc).(string) + "~" + (<-resultCrcMd5).(string)
		}(in)

	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	th := []int{0, 1, 2, 3, 4, 5}
	hash := make(map[string]string)
	ind := 0
	for value := range in {
		for _, element := range th {
			wg.Add(1)
			go func(ind int, element int, value interface{}) {
				defer wg.Done()
				itoa := strconv.Itoa(element)
				crc := DataSignerCrc32(itoa + value.(string))
				mu.Lock()
				hash[strconv.Itoa(ind)+strconv.Itoa(element)] = crc
				mu.Unlock()
			}(ind, element, value)
		}
		ind++
	}
	wg.Wait()

	arr := []string{}
	for key, _ := range hash {
		arr = append(arr, key)
	}

	sort.Strings(arr)

	for i := 0; i < ind; i++ {
		data := []string{}
		for j := i * len(th); j < (i+1)*(len(th)); j++ {
			data = append(data, hash[arr[j]])
		}
		answer := strings.Join(data, "")
		out <- answer
	}
}

func CombineResults(in, out chan interface{}) {

	data := []string{}
	for value := range in {
		data = append(data, value.(string))
	}

	sort.Strings(data)
	answer := strings.Join(data, "_")
	out <- answer

}

func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}
	chanIn := make(chan interface{})

	for _, yX := range jobs {
		chanOut := make(chan interface{})
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
	end := time.Since(start)
	fmt.Println(end)
}
