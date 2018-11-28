package main

// сюда писать тесты
import "testing"
import "fmt"
import "reflect"

func TestCalculator(t *testing.T) {

	var cases = []struct {
		input        []float64
		expected_val float64
		expected_err error
	}{
		{
			input:        []float64{2, 3, '+'},
			expected_val: 5,
			expected_err: nil,
		},

		{
			input:        []float64{2, 3, '-', 4, 5, '*', '+'},
			expected_val: 19,
			expected_err: nil,
		},

		{
			input:        []float64{11, 2, '*', 56},
			expected_val: 0,
			expected_err: fmt.Errorf("Incorrect Expression. A lot of components"),
		},

		{
			input:        []float64{15, 2, '*', 30, 2, '/'},
			expected_val: 0,
			expected_err: fmt.Errorf("Incorrect Expression. A lot of components"),
		},

		{
			input:        []float64{12, 6, '/'},
			expected_val: 2,
			expected_err: nil,
		},

		{
			input:        []float64{11, 2, '*', '+'},
			expected_val: 0,
			expected_err: fmt.Errorf("Incorrect Expression. Lack of components"),
		},

		{
			input:        []float64{11, '*'},
			expected_val: 0,
			expected_err: fmt.Errorf("Incorrect Expression. Lack of components"),
		},

		{
			input:        []float64{7, 0, '/'},
			expected_val: 0,
			expected_err: fmt.Errorf("Division by 0"),
		},

		{
			input:        []float64{1, 5, -5, 5, 2, '*', '+', '-', '/'},
			expected_val: 0,
			expected_err: fmt.Errorf("Division by 0"),
		},

		{
			input:        []float64{100, 5, 0, 5, 3, '*', '+', '-', '/'},
			expected_val: -10,
			expected_err: nil,
		},
	}

	for _, item := range cases {
		result_val, result_err := calculator(item.input)
		if !reflect.DeepEqual(result_val, item.expected_val) && !reflect.DeepEqual(result_err, item.expected_err) {
			t.Error("Expected: ", item.expected_val, "Have: ", result_val)
			t.Error("Expected: ", item.expected_err, "Have: ", result_err)
		}
	}
}
