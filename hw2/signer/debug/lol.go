package main

// сюда писать код
import (
	//"runtime"
	"strconv"
	//"sort"
	"fmt"
	//"strings"
	"time"
	"sync"
)

func SingleHash(in chan int, out chan string) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	for value := range in{
		wg.Add(2)
		itoa := strconv.Itoa(value)

			defer wg.Done()
		go func(itoa string) {
			res_1 := make(chan string)
			res_2 := make(chan string)
			defer wg.Done()
			res_1 <- DataSignerCrc32(itoa)

		go func(itoa string) {
			defer wg.Done()
			mu.Lock()
			md5 := DataSignerMd5(itoa)
			mu.Unlock()

			res_2 <- DataSignerCrc32(md5)

		}(itoa)

			out <- <-res_1 + "~" + <-res_2
		}(itoa)
	//	go func(){
	//		defer wg.Done()
	//	}()

	}
	wg.Wait()


	/*hash :=  make(map[int]string)
	for value := res_1{
		//out <- <-res_1 + "~" + <-res_2
	}

	*/
	close(out)
	for value := range out {
		fmt.Println(value)
	}

}

func main() {
	//data := []int{0,1}
	start := time.Now()
	data := []int{0, 1, 1, 2, 3, 5, 8}
	//data := []int{0}
	//var slc []string
	ch_data := make(chan int, len(data))
	ch_ans := make(chan string, len(data))
	//ch_data := make(chan int, len(data))
	for _, fibNum := range data {
		ch_data <- fibNum
	}
	close(ch_data)
	SingleHash(ch_data, ch_ans)
//	for val := range ch_ans{
	/*for i := 0; i < len(ch_ans); i++ {
		fmt.Println(len(ch_ans))
		fmt.Println("LOL",<-ch_ans)
	}*/
	//MultiHash
	//CombineResults
	/*for i := 0; i < 7; i++ {
		fmt.Println(<-ch_data)
		//runtime.Gosched()
	}
	*/

	//close(ch_data)
	end := time.Since(start)
	fmt.Println(end)
}


