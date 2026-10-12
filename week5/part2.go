package main

import (
	"errors"
	"fmt"
)

type Celsius float64

func (c Celsius) String() string { return fmt.Sprintf("%.1f°C", float64(c)) }

// task2: we gotta describe function with bool and map support
func describe(v any) string {
	switch x := v.(type) {
	case nil:
		return "nil"
	case int, int64:
		return fmt.Sprintf("integer %v", x)
	case string:
		return fmt.Sprintf("string of %d bytes", len(x))
	case []int:
		return fmt.Sprintf("slice with %d ints", len(x))
	case error:
		return "error: " + x.Error()
	case fmt.Stringer:
		return "Stringer: " + x.String()
	case bool:
		return fmt.Sprintf("boolean %v", x)
	case map[string]int:
		return fmt.Sprintf("map with %d elements", len(x))
	default:
		return fmt.Sprintf("unknown type %T", x)
	}
}

// here we gotta describe without type switch using comma-ok assertions
func describeWithoutSwitch(i any) string {
	if i == nil {
		return "nil"
	}
	if x, ok := i.(int); ok {
		return fmt.Sprintf("integer %v", x)
	}
	if x, ok := i.(int64); ok {
		return fmt.Sprintf("integer %v", x)
	}
	if x, ok := i.(string); ok {
		return fmt.Sprintf("string of %d bytes", len(x))
	}
	if x, ok := i.([]int); ok {
		return fmt.Sprintf("slice with %d ints", len(x))
	}
	if x, ok := i.(error); ok {
		return "error: " + x.Error()
	}
	if x, ok := i.(fmt.Stringer); ok {
		return "Stringer: " + x.String()
	}
	if x, ok := i.(bool); ok {
		return fmt.Sprintf("boolean %v", x)
	}
	if x, ok := i.(map[string]int); ok {
		return fmt.Sprintf("map with %d elements", len(x))
	}
	return fmt.Sprintf("unknown type %T", i)
}

func main222() {
	values := []any{
		nil,
		7,
		int64(7),
		"привет", // i used here cyrillic chars cuz they take 2 bytes each so "привет" is 12 bytes
		[]int{1, 2, 3},
		errors.New("disk full"),
		Celsius(21.5),
		3.14,
		struct{}{},
		true,
		map[string]int{"a": 1},
	}

	fmt.Println("Type Switch Version")
	for _, val := range values {
		fmt.Println(describe(val))
	}

	fmt.Println("\nComma-Ok Assertions Version")
	for _, val := range values {
		fmt.Println(describeWithoutSwitch(val))
	}
}
