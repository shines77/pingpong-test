package main

import (
	"bufio"
	"fmt"
	"os"
)

func Readin() {
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}

func ReadinEcho() {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	fmt.Printf("Input First Char Is : %v\n", string([]byte(input)[0]))
}

func WaitEnter() {
	fmt.Print("Press [Enter] key to continue ...")
	Readin()
}

func WaitEnterln() {
	fmt.Print("Press [Enter] key to continue ...")
	Readin()
	fmt.Print("\n")
}
