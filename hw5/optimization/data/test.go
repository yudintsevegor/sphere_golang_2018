package main

import "fmt"
import "bufio"
import "os"

func main() {
    file, _ := os.Open("users.txt")
    f := bufio.NewReader(file)
    for {
        fmt.Print("I: ")
        read_line, _ := f.ReadString('\n')
        fmt.Print(read_line)
	if read_line == ""{
		break
	}
    }
    file.Close()
}
