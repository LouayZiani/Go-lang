// anonymous and empty structs
package main

import "fmt"

func main() {
	point := struct{ X, Y int }{1, 2} // one-off struct without a type name
	seen := map[string]struct{}{}     // a "set" (remember Week 3 maps)
	seen["go"] = struct{}{}
	_, ok := seen["go"] // true

	fmt.Println(point, ok) // {1 2} true
	_, ok = seen["rust"]
	fmt.Println(ok) // false
}

// Without the Println the file wouldn't compile: Go refuses unused variables (point and ok were "declared and not used").
// struct{} takes 0 bytes, so the map stores only keys. That's why map[string]struct{} is the idiomatic set
// (map[string]bool works too but wastes a byte per entry and lets you store false, which is confusing).
