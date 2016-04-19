package main

import (
	"fmt"
	"os"
)

func main() {
	var cmd string
	fmt.Println(os.Args)
	for {
		fmt.Scanln(&cmd)
		fmt.Println(cmd)
	}
}
