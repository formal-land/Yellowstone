package main

import (
	"fmt"
	"os"
)

func main() {
	// Computes (2 + 3) * 4 = 20
	program := []Instruction{
		Push{Value: 2},
		Push{Value: 3},
		Add{},
		Push{Value: 4},
		Mul{},
		Print{},
		Halt{},
	}

	vm := NewVM()
	if err := vm.Run(program); err != nil {
		fmt.Fprintln(os.Stderr, "vm error:", err)
		os.Exit(1)
	}
}
