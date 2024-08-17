package tinymapper

import (
	errors2 "errors"
	"fmt"
	"reflect"
)

type FieldMapping map[string]string

type Mapper struct {
	registry map[reflect.Type]map[reflect.Type]func(a any, b any)
}

func New() *Mapper {
	return &Mapper{
		registry: make(map[reflect.Type]map[reflect.Type]func(a any, b any)),
	}
}

// Register adds a new mapping definition to the mapper.
// In the `mapping` func, modify the B pointer's properties to
// ensure your mapping is applied.
// Do not modify the A object, as that is the source, and will be discarded.
func Register[A any, B any](m *Mapper, mapping func(A, *B)) {
	var to B
	var from A

	fromType := reflect.TypeOf(from)
	toType := reflect.TypeOf(to)

	if m.registry[fromType] == nil {
		m.registry[fromType] = make(map[reflect.Type]func(src any, dest any))
	}

	m.registry[fromType][toType] = func(a any, b any) {
		mapping(a.(A), b.(*B))
	}
}

// To maps between two objects, while creating a new T to return.
func To[T any, S any](m *Mapper, source S) (T, error) {
	var to T
	fromType := reflect.TypeOf(source)
	toType := reflect.TypeOf(to)

	f := m.registry[fromType][toType]
	if f == nil {
		return to, fmt.Errorf("mapped type [%s -> %s] was not registered", fromType.String(), toType.String())
	}
	patchStruct(&to, source)
	f(source, &to)

	return to, nil
}

// ToWith uses an existing pointer as the source object, instead of creating a new one.
// Useful if you already have a object to modify.
func ToWith[T any, S any](m *Mapper, source S, existing *T) error {
	fromType := reflect.TypeOf(source)
	toType := reflect.TypeOf(*existing)

	f := m.registry[fromType][toType]
	if f == nil {
		return fmt.Errorf("mapped type [%s -> %s] was not registered", fromType.String(), toType.String())
	}
	patchStruct(existing, source)
	f(source, existing)

	return nil
}

// ArrayTo maps between lists of objects, using registered single definitions.
func ArrayTo[T any, S any](m *Mapper, source []S) ([]T, error) {
	var dests []T
	var errors []error
	for _, s := range source {
		mapped, err := To[T, S](m, s)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		dests = append(dests, mapped)
	}
	return dests, errors2.Join(errors...)
}
