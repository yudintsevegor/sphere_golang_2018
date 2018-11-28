package main

// сюда писать код
// фукция main тоже будет тшут

import "fmt"
import "reflect"
import "strconv"

type Stack struct {
	size    int
	data    []float64
	note    error
	element float64
}

func (st *Stack) StackPush(elmnt float64) {
	st.data = append(st.data, elmnt)
	st.size = len(st.data)
}

func (st *Stack) StackPop() {
	if st.size == 0 {
		st.note = fmt.Errorf("Stack is empty")
	} else if st.note == fmt.Errorf("Stack is empty") {
		st.element = st.data[st.size-1]
		st.data = st.data[:st.size-1]
		st.size = len(st.data)

	} else {
		st.note = nil
		st.element = st.data[st.size-1]
		st.data = st.data[:st.size-1]
		st.size = len(st.data)
	}
}

func calculator(arr_str []float64) (float64, error) {

	var st Stack
	for _, item := range arr_str {
		switch item {
		case '\n':
			fallthrough

		case '+':
			st.StackPop()
			val_1 := st.element
			st.StackPop()
			val_2 := st.element
			st.StackPush(val_1 + val_2)

		case '-':
			st.StackPop()
			val_1 := st.element
			st.StackPop()
			val_2 := st.element
			st.StackPush(-val_1 + val_2)

		case '*':
			st.StackPop()
			val_1 := st.element
			st.StackPop()
			val_2 := st.element
			st.StackPush(val_1 * val_2)

		case '/':
			st.StackPop()
			val_1 := st.element
			st.StackPop()
			val_2 := st.element
			st.StackPush(val_2 / val_1)

		case '=':
			break
		default:
			var helper float64
			if reflect.TypeOf(item) == reflect.TypeOf(helper) {
				st.StackPush(item)
			} else {
				fmt.Println("Its not an integer: ")
				fmt.Println(item)
			}
		}
	}

	if st.size > 1 {
		return 0, fmt.Errorf("Incorrect Expression. A lot of components")
	} else if reflect.DeepEqual(st.note, fmt.Errorf("Stack is empty")) {
		return 0,  fmt.Errorf("Incorrect Expression. Lack of components")
	} else if strconv.FormatFloat(st.data[0], 'f', 6, 64) == "-Inf" || strconv.FormatFloat(st.data[0], 'f', 6, 64) == "+Inf" {
		return 0, fmt.Errorf("Division by 0")
	} // else {

	//	fmt.Println(st.note)
	//	}
	fmt.Println(st.data[0], st.note)
	return st.data[0], st.note
}

func main() {
/**/
	arr_str := []float64{11, 2, '*', '+'}
	//arr_str := []float66{2, 3, '+', 3, 5, '*', '='}
	//arr_str := []float64{2, 3, 4, 5, 6, '*', '+', '-', '/'}
	el, er := calculator(arr_str)
	fmt.Println(arr_str)
	fmt.Println(el)
	fmt.Println(er)
	asciiNum := 42
	character := string(asciiNum)
	fmt.Println(asciiNum, " : ", character)
}
