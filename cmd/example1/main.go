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

	// Now add ProcessC which is a provider of InterfaceC and is dependent
	// on InterfaceC and InterfaceD.
	processC := newProcC()
	err = m.Add(processC)
	if err != nil {
		return err
	}

	fmt.Println(processC.GetID())

	return nil
}

type InterfaceA interface {
	GetID() (string, error)
	GetThat() (int, error)
}

// ProcessA is a provider of InterfaceA but has no dependencies.
type ProcessA struct {
	id string
}

func newProcA() *ProcessA {
	return &ProcessA{
		id: "A",
	}
}

func (p *ProcessA) GetID() (string, error) {
	return p.id, nil
}

func (p *ProcessA) GetThat() (int, error) {
	return 5, nil
}

// Provide is a method required by the refueler so that ProcessA can be
// registered as a provider of InterfaceA.
func (p *ProcessA) Provide() InterfaceA {
	return p
}

type InterfaceB interface {
	GetID() (string, error)
	GetSomethingElse() (int, error)
}

// ProcessBDeps defines the dependencies of ProcessB.
type ProcessBDeps struct {
	A InterfaceA
}

// ProcessB is a provider of InterfaceB and depends on InterfaceA.
type ProcessB struct {
	Backends ProcessBDeps
	id       string
}

func newProcB() *ProcessB {
	return &ProcessB{
		id: "B",
	}
}

func (p *ProcessB) GetID() (string, error) {
	aID, err := p.Backends.A.GetID()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("my ID: %s, A's ID: %s", p.id, aID), nil
}

func (p *ProcessB) GetSomethingElse() (int, error) {
	return 4, nil
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
	id       string
}

func newProcC() *ProcessC {
	return &ProcessC{
		id: "C",
	}
}

func (p *ProcessC) GetID() (string, error) {
	b := p.Backends

	aID, err := b.A.GetID()
	if err != nil {
		return "", err
	}

	bID, err := b.B.GetID()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("my ID: %s, A's ID: %s, B's ID %s", p.id, aID, bID),
		nil
}
