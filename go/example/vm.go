package main

import "fmt"

type VM struct {
	stack []int64
	pc    int
	halt  bool
}

func NewVM() *VM {
	return &VM{stack: make([]int64, 0, 32)}
}

func (vm *VM) push(v int64) {
	vm.stack = append(vm.stack, v)
}

func (vm *VM) pop() (int64, error) {
	n := len(vm.stack)
	if n == 0 {
		return 0, fmt.Errorf("stack underflow at pc=%d", vm.pc)
	}
	v := vm.stack[n-1]
	vm.stack = vm.stack[:n-1]
	return v, nil
}

func (vm *VM) Run(program []Instruction) error {
	for vm.pc = 0; vm.pc < len(program) && !vm.halt; vm.pc++ {
		if err := program[vm.pc].Execute(vm); err != nil {
			return err
		}
	}
	return nil
}

func (vm *VM) Top() (int64, bool) {
	if len(vm.stack) == 0 {
		return 0, false
	}
	return vm.stack[len(vm.stack)-1], true
}
