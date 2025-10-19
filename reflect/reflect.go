package reflect

import "reflect"

// TypeFor returns the Type that represents the type argument T.
func TypeFor[T any]() Type { return reflect.TypeFor[T]() }

// Type exposes all the methods common to ALL Kinds of types.
func TypeOf(i any) Type { return reflect.TypeOf(i) }

type Type interface {

	// Align returns the alignment in bytes of a value of
	// this type when allocated in memory.
	Align() int

	// FieldAlign returns the alignment in bytes of a value of
	// this type when used as a field in a struct.
	FieldAlign() int

	// Method returns the i'th method in the type's method set.
	// It panics if i is not in the range [0, NumMethod()).
	//
	// For a non-interface type T or *T, the returned Method's Type and Func
	// fields describe a function whose first argument is the receiver,
	// and only exported methods are accessible.
	//
	// For an interface type, the returned Method's Type field gives the
	// method signature, without a receiver, and the Func field is nil.
	//
	// Methods are sorted in lexicographic order.
	//
	// Calling this method will force the linker to retain all exported methods in all packages.
	// This may make the executable binary larger but will not affect execution time.
	Method(int) reflect.Method

	// MethodByName returns the method with that name in the type's
	// method set and a boolean indicating if the method was found.
	//
	// For a non-interface type T or *T, the returned Method's Type and Func
	// fields describe a function whose first argument is the receiver.
	//
	// For an interface type, the returned Method's Type field gives the
	// method signature, without a receiver, and the Func field is nil.
	//
	// Calling this method will cause the linker to retain all methods with this name in all packages.
	// If the linker can't determine the name, it will retain all exported methods.
	// This may make the executable binary larger but will not affect execution time.
	MethodByName(string) (reflect.Method, bool)

	// NumMethod returns the number of methods accessible using Method.
	//
	// For a non-interface type, it returns the number of exported methods.
	//
	// For an interface type, it returns the number of exported and unexported methods.
	NumMethod() int

	// Name returns the type's name within its package for a defined type.
	// For other (non-defined) types it returns the empty string.
	Name() string

	// PkgPath returns a defined type's package path, that is, the import path
	// that uniquely identifies the package, such as "encoding/base64".
	// If the type was predeclared (string, error) or not defined (*T, struct{},
	// []int, or A where A is an alias for a non-defined type), the package path
	// will be the empty string.
	PkgPath() string

	// Size returns the number of bytes needed to store
	// a value of the given type; it is analogous to unsafe.Sizeof.
	Size() uintptr

	// String returns a string representation of the type.
	// The string representation may use shortened package names
	// (e.g., base64 instead of "encoding/base64") and is not
	// guaranteed to be unique among types. To test for type identity,
	// compare the Types directly.
	String() string

	// Kind returns the specific kind of this type.
	Kind() reflect.Kind

	// Implements reports whether the type implements the interface type u.
	Implements(u reflect.Type) bool

	// AssignableTo reports whether a value of the type is assignable to type u.
	AssignableTo(u reflect.Type) bool

	// ConvertibleTo reports whether a value of the type is convertible to type u.
	// Even if ConvertibleTo returns true, the conversion may still panic.
	// For example, a slice of type []T is convertible to *[N]T,
	// but the conversion will panic if its length is less than N.
	ConvertibleTo(u reflect.Type) bool

	// Comparable reports whether values of this type are comparable.
	// Even if Comparable returns true, the comparison may still panic.
	// For example, values of interface type are comparable,
	// but the comparison will panic if their dynamic type is not comparable.
	Comparable() bool
	/*
		// Bits returns the size of the type in bits.
		// It panics if the type's Kind is not one of the
		// sized or unsized Int, Uint, Float, or Complex kinds.
		Bits() int

		// ChanDir returns a channel type's direction.
		// It panics if the type's Kind is not Chan.
		ChanDir() ChanDir

		// IsVariadic reports whether a function type's final input parameter
		// is a "..." parameter. If so, t.In(t.NumIn() - 1) returns the parameter's
		// implicit actual type []T.
		//
		// For concreteness, if t represents func(x int, y ... float64), then
		//
		//	t.NumIn() == 2
		//	t.In(0) is the reflect.Type for "int"
		//	t.In(1) is the reflect.Type for "[]float64"
		//	t.IsVariadic() == true
		//
		// IsVariadic panics if the type's Kind is not Func.
		IsVariadic() bool

		// Elem returns a type's element type.
		// It panics if the type's Kind is not Array, Chan, Map, Pointer, or Slice.
		Elem() Type

		// Field returns a struct type's i'th field.
		// It panics if the type's Kind is not Struct.
		// It panics if i is not in the range [0, NumField()).
		Field(i int) StructField

		// FieldByIndex returns the nested field corresponding
		// to the index sequence. It is equivalent to calling Field
		// successively for each index i.
		// It panics if the type's Kind is not Struct.
		FieldByIndex(index []int) StructField

		// FieldByName returns the struct field with the given name
		// and a boolean indicating if the field was found.
		// If the returned field is promoted from an embedded struct,
		// then Offset in the returned StructField is the offset in
		// the embedded struct.
		FieldByName(name string) (StructField, bool)

		// FieldByNameFunc returns the struct field with a name
		// that satisfies the match function and a boolean indicating if
		// the field was found.
		//
		// FieldByNameFunc considers the fields in the struct itself
		// and then the fields in any embedded structs, in breadth first order,
		// stopping at the shallowest nesting depth containing one or more
		// fields satisfying the match function. If multiple fields at that depth
		// satisfy the match function, they cancel each other
		// and FieldByNameFunc returns no match.
		// This behavior mirrors Go's handling of name lookup in
		// structs containing embedded fields.
		//
		// If the returned field is promoted from an embedded struct,
		// then Offset in the returned StructField is the offset in
		// the embedded struct.
		FieldByNameFunc(match func(string) bool) (StructField, bool)

		// In returns the type of a function type's i'th input parameter.
		// It panics if the type's Kind is not Func.
		// It panics if i is not in the range [0, NumIn()).
		In(i int) Type

		// Key returns a map type's key type.
		// It panics if the type's Kind is not Map.
		Key() Type

		// Len returns an array type's length.
		// It panics if the type's Kind is not Array.
		Len() int

		// NumField returns a struct type's field count.
		// It panics if the type's Kind is not Struct.
		NumField() int

		// NumIn returns a function type's input parameter count.
		// It panics if the type's Kind is not Func.
		NumIn() int

		// NumOut returns a function type's output parameter count.
		// It panics if the type's Kind is not Func.
		NumOut() int

		// Out returns the type of a function type's i'th output parameter.
		// It panics if the type's Kind is not Func.
		// It panics if i is not in the range [0, NumOut()).
		Out(i int) Type

		// OverflowComplex reports whether the complex128 x cannot be represented by type t.
		// It panics if t's Kind is not Complex64 or Complex128.
		OverflowComplex(x complex128) bool

		// OverflowFloat reports whether the float64 x cannot be represented by type t.
		// It panics if t's Kind is not Float32 or Float64.
		OverflowFloat(x float64) bool

		// OverflowInt reports whether the int64 x cannot be represented by type t.
		// It panics if t's Kind is not Int, Int8, Int16, Int32, or Int64.
		OverflowInt(x int64) bool

		// OverflowUint reports whether the uint64 x cannot be represented by type t.
		// It panics if t's Kind is not Uint, Uintptr, Uint8, Uint16, Uint32, or Uint64.
		OverflowUint(x uint64) bool
	*/
	// CanSeq reports whether a [Value] with this type can be iterated over using [Value.Seq].
	CanSeq() bool

	// CanSeq2 reports whether a [Value] with this type can be iterated over using [Value.Seq2].
	CanSeq2() bool
	// contains filtered or unexported methods
}

// Method represents a single method.
type Method struct {
	Name    string
	PkgPath string
	Type    Type
	Func    any
	Index   int
}

// StructField represents a single field in a struct.
type StructField struct {
	Name      string
	PkgPath   string
	Type      Type
	Tag       string
	Offset    uintptr
	Index     []int
	Anonymous bool
}

// Specialized type structs for each Kind, containing only methods that are safe for that kind.

// BoolType represents a boolean type.
type BoolType struct {
	Type
}

// IntType represents a signed integer type.
type IntType struct {
	Type
}

// Methods safe for IntType: Bits, OverflowInt (plus common methods)
func (t *IntType) Bits() int                { return t.Type.(reflect.Type).Bits() }
func (t *IntType) OverflowInt(x int64) bool { return t.Type.(reflect.Type).OverflowInt(x) }

// Int8Type represents an 8-bit signed integer type.
type Int8Type struct {
	Type
}

func (t *Int8Type) Bits() int                { return t.Type.(reflect.Type).Bits() }
func (t *Int8Type) OverflowInt(x int64) bool { return t.Type.(reflect.Type).OverflowInt(x) }

// Int16Type represents a 16-bit signed integer type.
type Int16Type struct {
	Type
}

func (t *Int16Type) Bits() int                { return t.Type.(reflect.Type).Bits() }
func (t *Int16Type) OverflowInt(x int64) bool { return t.Type.(reflect.Type).OverflowInt(x) }

// Int32Type represents a 32-bit signed integer type.
type Int32Type struct {
	Type
}

func (t *Int32Type) Bits() int                { return t.Type.(reflect.Type).Bits() }
func (t *Int32Type) OverflowInt(x int64) bool { return t.Type.(reflect.Type).OverflowInt(x) }

// Int64Type represents a 64-bit signed integer type.
type Int64Type struct {
	Type
}

func (t *Int64Type) Bits() int                { return t.Type.(reflect.Type).Bits() }
func (t *Int64Type) OverflowInt(x int64) bool { return t.Type.(reflect.Type).OverflowInt(x) }

// UintType represents an unsigned integer type.
type UintType struct {
	Type
}

func (t *UintType) Bits() int                  { return t.Type.(reflect.Type).Bits() }
func (t *UintType) OverflowUint(x uint64) bool { return t.Type.(reflect.Type).OverflowUint(x) }

// Uint8Type represents an 8-bit unsigned integer type.
type Uint8Type struct {
	Type
}

func (t *Uint8Type) Bits() int                  { return t.Type.(reflect.Type).Bits() }
func (t *Uint8Type) OverflowUint(x uint64) bool { return t.Type.(reflect.Type).OverflowUint(x) }

// Uint16Type represents a 16-bit unsigned integer type.
type Uint16Type struct {
	Type
}

func (t *Uint16Type) Bits() int                  { return t.Type.(reflect.Type).Bits() }
func (t *Uint16Type) OverflowUint(x uint64) bool { return t.Type.(reflect.Type).OverflowUint(x) }

// Uint32Type represents a 32-bit unsigned integer type.
type Uint32Type struct {
	Type
}

func (t *Uint32Type) Bits() int                  { return t.Type.(reflect.Type).Bits() }
func (t *Uint32Type) OverflowUint(x uint64) bool { return t.Type.(reflect.Type).OverflowUint(x) }

// Uint64Type represents a 64-bit unsigned integer type.
type Uint64Type struct {
	Type
}

func (t *Uint64Type) Bits() int                  { return t.Type.(reflect.Type).Bits() }
func (t *Uint64Type) OverflowUint(x uint64) bool { return t.Type.(reflect.Type).OverflowUint(x) }

// UintptrType represents a uintptr type.
type UintptrType struct {
	Type
}

func (t *UintptrType) Bits() int                  { return t.Type.(reflect.Type).Bits() }
func (t *UintptrType) OverflowUint(x uint64) bool { return t.Type.(reflect.Type).OverflowUint(x) }

// Float32Type represents a 32-bit floating point type.
type Float32Type struct {
	Type
}

func (t *Float32Type) Bits() int                    { return t.Type.(reflect.Type).Bits() }
func (t *Float32Type) OverflowFloat(x float64) bool { return t.Type.(reflect.Type).OverflowFloat(x) }

// Float64Type represents a 64-bit floating point type.
type Float64Type struct {
	Type
}

func (t *Float64Type) Bits() int                    { return t.Type.(reflect.Type).Bits() }
func (t *Float64Type) OverflowFloat(x float64) bool { return t.Type.(reflect.Type).OverflowFloat(x) }

// Complex64Type represents a 64-bit complex type.
type Complex64Type struct {
	Type
}

func (t *Complex64Type) Bits() int { return t.Type.(reflect.Type).Bits() }
func (t *Complex64Type) OverflowComplex(x complex128) bool {
	return t.Type.(reflect.Type).OverflowComplex(x)
}

// Complex128Type represents a 128-bit complex type.
type Complex128Type struct {
	Type
}

func (t *Complex128Type) Bits() int { return t.Type.(reflect.Type).Bits() }
func (t *Complex128Type) OverflowComplex(x complex128) bool {
	return t.Type.(reflect.Type).OverflowComplex(x)
}

// ArrayType represents an array type.
type ArrayType struct {
	Type
}

func (t *ArrayType) Elem() Type { return t.Type.(reflect.Type).Elem() }
func (t *ArrayType) Len() int   { return t.Type.(reflect.Type).Len() }

// ChanType represents a channel type.
type ChanType struct {
	Type
}

func (t *ChanType) ChanDir() reflect.ChanDir { return t.Type.(reflect.Type).ChanDir() }
func (t *ChanType) Elem() Type               { return t.Type.(reflect.Type).Elem() }

// FuncType represents a function type.
type FuncType struct {
	Type
}

func (t *FuncType) IsVariadic() bool { return t.Type.(reflect.Type).IsVariadic() }
func (t *FuncType) NumIn() int       { return t.Type.(reflect.Type).NumIn() }
func (t *FuncType) NumOut() int      { return t.Type.(reflect.Type).NumOut() }
func (t *FuncType) In(i int) Type    { return t.Type.(reflect.Type).In(i) }
func (t *FuncType) Out(i int) Type   { return t.Type.(reflect.Type).Out(i) }

// InterfaceType represents an interface type.
type InterfaceType struct {
	Type
}

// MapType represents a map type.
type MapType struct {
	Type
}

func (t *MapType) Key() Type  { return t.Type.(reflect.Type).Key() }
func (t *MapType) Elem() Type { return t.Type.(reflect.Type).Elem() }

// PointerType represents a pointer type.
type PointerType struct {
	Type
}

func (t *PointerType) Elem() Type { return t.Type.(reflect.Type).Elem() }

// SliceType represents a slice type.
type SliceType struct {
	Type
}

func (t *SliceType) Elem() Type { return t.Type.(reflect.Type).Elem() }

// StringType represents a string type.
type StringType struct {
	Type
}

// StructType represents a struct type.
type StructType struct {
	Type
}

func (t *StructType) Field(i int) reflect.StructField { return t.Type.(reflect.Type).Field(i) }
func (t *StructType) FieldByIndex(index []int) reflect.StructField {
	return t.Type.(reflect.Type).FieldByIndex(index)
}
func (t *StructType) FieldByName(name string) (reflect.StructField, bool) {
	return t.Type.(reflect.Type).FieldByName(name)
}
func (t *StructType) FieldByNameFunc(match func(string) bool) (reflect.StructField, bool) {
	return t.Type.(reflect.Type).FieldByNameFunc(match)
}
func (t *StructType) NumField() int { return t.Type.(reflect.Type).NumField() }

// UnsafePointerType represents an unsafe.Pointer type.
type UnsafePointerType struct {
	Type
}

/*
AI Prompt: This file has Go's reflect.Type interface defined in it. Produce
for me a 1 struct type for each Kind value. The name of each struct
type should of the form XxxType where Xxx is the kind value with
first letter uppercase so the struct is exported from this package.
Then for each type add only the methods that would not panic if
called for that kind. Do not add methods to each struct if the
method is applicable to all kind values.
*/
