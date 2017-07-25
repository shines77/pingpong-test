package main

import (
	"bufio"
	"fmt"
	"os"
)

func readin() {
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}

func readin_echo() {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	fmt.Printf("Input First Char Is : %v\n", string([]byte(input)[0]))
}

func waitEnter() {
	fmt.Print("Press [Enter] key to continue ...")
	readin()
}

func waitEnterLn() {
	fmt.Print("Press [Enter] key to continue ...")
	readin()
	fmt.Print("\n")
}
