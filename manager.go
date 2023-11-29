package refuel

import (
	"fmt"
	"reflect"
)

const (
	provideMethodName = "Provide"
	backendsName      = "Backends"
)

// Manager is a dependency injection controller.
type Manager struct {
	providers map[reflect.Type]interface{}
}

// NewManager builds a new Manager.
func NewManager() *Manager {
	return &Manager{
		providers: make(map[reflect.Type]interface{}),
	}
}

// Add adds a new interface provider to the set managed by the Manger. It also
// injects any required dependency into the provided instance. For the instance
// to be accepted, it must obey the following rules:
//  1. It must itself be a struct.
//  2. It must implement some defined interface which it will then be a provider
//     of, let's call it InterfaceA.
//  3. It must have a `Provide() InterfaceA` method.
//  4. The struct must have a member called `Backends`. This must be a struct
//     (not a pointer) which holds all the interfaces it wants the manage to
//     inject implementations into.
func (m *Manager) Add(instance interface{}) error {
	err := m.injectDependencies(instance)
	if err != nil {
		return err
	}

	instanceType := reflect.TypeOf(instance)
	provideMethod, exists := instanceType.MethodByName(provideMethodName)
	if !exists {
		// If it does not have a Provide method then there is no
		// provider to register.
		return nil
	}

	if provideMethod.Type.NumOut() != 1 {
		return fmt.Errorf("the given object's %s method must "+
			"have exactly one return value", provideMethodName)

	}

	providedType := provideMethod.Type.Out(0)
	if providedType.Kind() != reflect.Interface {
		return fmt.Errorf("the %s method does not return an "+
			"interface", provideMethodName)
	}

	m.providers[providedType] = instance

	return nil
}

func (m *Manager) injectDependencies(instance interface{}) error {
	val := reflect.ValueOf(instance).Elem()
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("must be a struct type")
	}

	backends := val.FieldByName(backendsName)
	if !backends.IsValid() {
		// If there is no dependency struct, then there is nothing to
		// inject.
		return nil
	}

	if backends.Kind() == reflect.Pointer {
		return fmt.Errorf("the %s field must be a struct type, "+
			"not a pointer", backends)
	}

	if backends.Kind() != reflect.Struct {
		return fmt.Errorf("the Backends field must be a struct type")
	}

	// Iterate through the fields of the Backends struct and inject
	// dependencies.
	for i := 0; i < backends.NumField(); i++ {
		field := backends.Field(i)

		if !field.CanSet() {
			return fmt.Errorf("field %s is not settable",
				field.String())
		}

		if field.Kind() != reflect.Interface {
			return fmt.Errorf("all members of a Backends struct " +
				"must be interfaces")
		}

		provider := m.providers[field.Type()]
		if provider == nil {
			return fmt.Errorf("no implementation found of %s",
				field.Type())
		}

		providerInstance := reflect.ValueOf(provider)
		provide := providerInstance.MethodByName(provideMethodName)

		field.Set(provide.Call(nil)[0])
	}

	return nil
}
