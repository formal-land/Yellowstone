package main

import "fmt"

type Instruction interface {
	Execute(vm *VM) error
}

type Push struct{ Value int64 }

func (i Push) Execute(vm *VM) error {
	vm.push(i.Value)
	return nil
}

type Add struct{}

func (Add) Execute(vm *VM) error {
	b, err := vm.pop()
	if err != nil {
		return err
	}
	a, err := vm.pop()
	if err != nil {
		return err
	}
	vm.push(a + b)
	return nil
}

type Mul struct{}

func (Mul) Execute(vm *VM) error {
	b, err := vm.pop()
	if err != nil {
		return err
	}
	a, err := vm.pop()
	if err != nil {
		return err
	}
	vm.push(a * b)
	return nil
}

type Print struct{}

func (Print) Execute(vm *VM) error {
	v, ok := vm.Top()
	if !ok {
		return fmt.Errorf("print on empty stack at pc=%d", vm.pc)
	}
	fmt.Println(v)
	return nil
}

type Halt struct{}

func (Halt) Execute(vm *VM) error {
	vm.halt = true
	return nil
}
