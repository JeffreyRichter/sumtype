package sumtype

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"reflect"
	"unsafe"
)

// Cast casts From a caster To another sum type projection type.
func Cast[To any, From any](caster *Caster[From]) *To { return (*To)(unsafe.Pointer(caster)) }

// Caster provides methods to cast between a sum type *Json and its variants (which all
// have the same fields as Json). The 1st field of Json and all its variants must be a non-exported
// xxxCaster type whose underlying type is sumtype.Caster[Json]. All methods that cast
// a pointer from one projection type to another, require by-ref receivers.
type Caster[Json any] struct{}

// RULES: Methods that cast a pointer from 1 type to another, require by-ref receiver (Caster[Json] methods)

// Json casts c to *Json where Json is the JSONable struct (ALL JSON fields are exported).
func (c *Caster[Json]) Json() *Json { return Cast[Json](c) }

// MarshalJSON marshals the json struct instance to JSON
func (c *Caster[Json]) MarshalJSON() ([]byte, error) { return json.Marshal(c.Json()) }

// UnmarshalJSON unmarshals JSON data to the Json struct instance
func (c *Caster[Json]) UnmarshalJSON(data []byte) error { return json.Unmarshal(data, c.Json()) }

// String returns a readable JSON representation of the Json struct instance
func (c *Caster[Json]) String() string {
	j, _ := json.Marshal(c.Json(), jsontext.Multiline(true))
	return string(j)
}

// ZeroNonKindFields sets all fields not relevant to "Kind" to their zero value
func (c *Caster[Json]) ZeroNonKindFields(ptrToKindStruct any) {
	jsonFields := reflect.ValueOf(c.Json()).Elem() // Dereference the pointer to get the struct value
	kindFields := reflect.TypeOf(ptrToKindStruct).Elem()
	for f := range kindFields.NumField() {
		if kindField := kindFields.Field(f); !kindField.IsExported() {
			// kindField's unexported fields are zero'd from jsonFields' equivalent field
			if fieldToZero := jsonFields.Field(f); fieldToZero.CanSet() {
				fieldToZero.Set(reflect.Zero(fieldToZero.Type())) // Requires exported field
			}
		}
	}
}

// ValidateStructFields ensures that Json and all the specific projection types have struct fields
// in the same order and same type. If panicOnError is true, ValidateStructFields panics if
// there is an error, otherwise it returns the error (or nil if no error).
func (c Caster[Json]) ValidateStructFields(panicOnError bool, structs ...any) error {
	err := c.validateStructFields(structs...)
	if panicOnError && err != nil {
		panic(err)
	}
	return err
}

// validateStructFields ensures that Json and all the specific projection types have struct fields
// in the same order and same type. It returns nil or an error.
func (c Caster[Json]) validateStructFields(structs ...any) error {
	mainStruct := reflect.TypeFor[Json]()

	// The 1st field's underlying type must be sumtype.Converter for the unsafe casts to work
	underlyingType, firstFieldType := reflect.TypeOf(c), mainStruct.Field(0).Type
	if !firstFieldType.ConvertibleTo(underlyingType) || !underlyingType.ConvertibleTo(firstFieldType) {
		return fmt.Errorf("first field of struct %s must be a type whose underlying type is %T", mainStruct.Name(), c)
	}

	for _, otherStruct := range structs {
		otherStructType := reflect.TypeOf(otherStruct)
		if mainStruct.NumField() != otherStructType.NumField() {
			return fmt.Errorf("structs have different number of fields: %s=%d vs %s=%d",
				mainStruct.Name(), mainStruct.NumField(), otherStructType.Name(), otherStructType.NumField())
		}

		for f := range mainStruct.NumField() {
			// Struct fields must be in same order and same type
			mf, of := mainStruct.Field(f), otherStructType.Field(f)
			if mf.Type != of.Type {
				return fmt.Errorf("structs have incompatible field #%d: %s.%s (%s) vs %s.%s (%s)",
					f, mainStruct.Name(), mf.Name, mf.Type.String(),
					otherStructType.Name(), of.Name, of.Type.String())
			}
		}
	}
	return nil
}
