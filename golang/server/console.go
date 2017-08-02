package main

import (
	"bufio"
	"fmt"
	"os"
)

// Readin : Read a string input
func Readin() {
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}

// ReadinEcho : Read a string input and echo to terminator
func ReadinEcho() {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	fmt.Printf("Input First Char Is : %v\n", string([]byte(input)[0]))
}

// WaitEnter : wait a [ENTER] key
func WaitEnter() {
	fmt.Print("Press [Enter] key to continue ...")
	Readin()
}

// WaitEnterln : wait a [ENTER] key and new line
func WaitEnterln() {
	fmt.Print("Press [Enter] key to continue ...")
	Readin()
	fmt.Print("\n")
}
