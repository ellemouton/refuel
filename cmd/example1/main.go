package main

import (
	"fmt"
	"github.com/ellemouton/refuel"
	"log"
)

func main() {
	if err := doThings(); err != nil {
		log.Fatalln(err)
	}
}

func doThings() error {
	m := refuel.NewManager()

	processA := newProcA()

	// Register ProcessA as a provider of InterfaceA.
	err := m.Add(processA)
	if err != nil {
		return err
	}

	// Now add ProcessB which is a provider of InterfaceB and is dependent
	// on InterfaceA.
	processB := newProcB()
	err = m.Add(processB)
	if err != nil {
		return err
	}

	processB.Print()

	// Now add ProcessC which is a provider of InterfaceC and is dependent
	// on InterfaceC and InterfaceD.
	processC := newProcC()
	err = m.Add(processC)
	if err != nil {
		return err
	}

	processC.Print()

	return nil
}

type InterfaceA interface {
	GetX() string
}

// ProcessA is a provider of InterfaceA but has no dependencies.
type ProcessA struct {
	x string
}

func newProcA() *ProcessA {
	return &ProcessA{
		x: "A",
	}
}

func (p *ProcessA) GetX() string {
	return p.x
}

// Provide is a method required by the refueler so that ProcessA can be
// registered as a provider of InterfaceA.
func (p *ProcessA) Provide() InterfaceA {
	return p
}

type InterfaceB interface {
	GetY() string
}

// ProcessBDeps defines the dependencies of ProcessB.
type ProcessBDeps struct {
	A InterfaceA
}

// ProcessB is a provider of InterfaceB and depends on InterfaceA.
type ProcessB struct {
	Backends ProcessBDeps
	y        string
}

func newProcB() *ProcessB {
	return &ProcessB{
		y: "B",
	}
}

func (p *ProcessB) GetY() string {
	return p.y
}

func (p *ProcessB) Print() {
	x := p.Backends.A.GetX()

	fmt.Printf("x: %s, y: %s\n", x, p.y)
}

// Provide is a method required by the refueler so that ProcessB can be
// registered as a provider of InterfaceB.
func (p *ProcessB) Provide() InterfaceB {
	// Return the interface implementation
	return p
}

// ProcessCDeps holds the dependencies of ProcessC.
type ProcessCDeps struct {
	A InterfaceA
	B InterfaceB
}

// ProcessC does not provide any interface, but it depends on InterfaceA and
// InterfaceB.
type ProcessC struct {
	Backends ProcessCDeps
	z        string
}

func newProcC() *ProcessC {
	return &ProcessC{
		z: "C",
	}
}

func (p *ProcessC) Print() {
	b := p.Backends

	x := b.A.GetX()
	y := b.B.GetY()

	fmt.Printf("x: %s, y: %s, z: %s\n", x, y, p.z)
}
