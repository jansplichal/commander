//usr/local/go/bin/go run $0 $@ ; exit
package main

import (
	"fmt"
	"sync"
)

/*Name ..*/
type Name struct {
	a string
	b string
}

var hits struct {
	sync.Mutex
	n int
}

func main() {
	fmt.Println("Starting a timer")

	fmt.Println(hits)

	var wait string
	fmt.Scanln(&wait)
}
