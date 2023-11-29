# Refuel

Playing around with dependency injection. 

### Example:

Let's say I have the following dependency tree: 

- Process A has no dependencies. 
- Process B depends on Process A. 
- Process C depends on Process A and B. 

To add something to the refueler, itt can be a Provider, have dependencies, 
or both. In our example above:
  - Process A is used by other processes, so it needs to be a Provider. 
  - Process B needs to be a provider _and_ it must indicate that it depends on 
    Process A. 
  - Process C only needs to indicate that it depends on Process A and B. 

Process A can thus be set up as follows:
```
// InterfaceA is the interface that Process A will provide.
type InterfaceA interface {
	GetX() string
}

// ProcessA is a provider of InterfaceA and has no dependencies.
type ProcessA struct {
	x string
}

// newProcA constructs ProcessA.
func newProcA() *ProcessA {
	return &ProcessA{
		x: "A",
	}
}

// GetX ensures that ProcessA implements the InterfaceA interface
func (p *ProcessA) GetX() string {
	return p.x
}

// Provide is a method required by the refueler so that ProcessA can be
// registered as a provider of InterfaceA.
func (p *ProcessA) Provide() InterfaceA {
	return p
}

```

Process B can be set up as follows:

```
// InterfaceB is the interface that Process B will provide.
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
	y       string
}

// newProcB constructs ProcessB.
func newProcB() *ProcessB {
	return &ProcessB{
		y: "B",
	}
}

// GetY ensures that ProcessB implements the InterfaceB interface.
func (p *ProcessB) GetY() string {
	return p.y
}

// Print demonstrates ProcessB using InterfaceA. 
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

```

ProcessC can be set up as follows:

```
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

// Print demonstrates ProcessC using InterfaceA and InterfaceB. 
func (p *ProcessC) Print() {
	b := p.Backends

	x := b.A.GetX()
	y := b.B.GetY()

	fmt.Printf("x: %s, y: %s, z: %s\n", x, y, p.z)
}
```

```
// Create a new manager.
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
```