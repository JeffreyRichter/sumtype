package sumtype_test

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"

	"github.com/JeffreyRichter/sumtype"
)

// ********** THE CODE BELOW SHOWS HOW TO USE A SUM TYPE ********** //

// ptr returns a pointer to the given value. 
// Delete this function and use 'new' when using Go 1.26.
func ptr[T any](v T) *T { return &v }

// jsonFromWebService shows how to build an array of Shape objects of various kinds.
func jsonFromWebService() []byte {
	shapes := []*Shape{
		(&CircleShape{
			Kind:   ptr(CircleShapeKind),
			Color:  ptr("red"),
			Radius: ptr(1),
		}).Shape(),

		(&RectangleShape{
			Kind:   ptr(RectangleShapeKind),
			Color:  ptr("green"),
			Width:  ptr(15),
			Height: ptr(15),
		}).Shape(),

		(&RectangleShape{
			Kind:   ptr(RectangleShapeKind),
			Color:  ptr("blue"),
			Width:  ptr(5),
			Height: ptr(5),
		}).Shape(),
	}

	// Marshal the array to JSON
	incomingJson, _ := json.Marshal(shapes, jsontext.WithIndent("  "))
	fmt.Println(string(incomingJson) + "\n----------") // Show what the JSON looks like
	return ([]byte)(incomingJson)	// Return JSON to caller/client
}

func Example() {
	// Simulate getting a JSON array of shape objects from a Web Service
	incomingJson := jsonFromWebService()

	// Unmarshal JSON to an array of Shape objects
	var shapes []*Shape
	_ = json.Unmarshal(incomingJson, &shapes)

	// Process the shape objects
	for _, s := range shapes {
		switch *s.Kind {
		case CircleShapeKind:
			c := s.Circle()                          // Cast to Circle for code completion/type-safety
			c.Color, c.Radius = ptr("white"), ptr(2) // Demo changing values

		case RectangleShapeKind:
			r := s.Rectangle()                 // Cast to Rectangle for code completion/type-safety
			r.Height = ptr(min(*r.Height, 10)) // Forbid a Height > 10
			r.Width = ptr(min(*r.Width, 10))   // Forbid a Width > 10
			if *r.Height < 10 && *r.Width < 10 {
				// Demo: Convert any Rectangle whose Width/Height < 10 to a Circle
				c := s.SetCircle()
				c.Color, c.Radius = ptr("blue"), ptr(5)
			}

		default:
			// This can happen if the Web Service returns a new kind (perhaps
			// in a new version) that this client code never knew about.
			fmt.Printf("Unrecognized shape kind: %s\n", *s.Kind)
		}
	}

	// Marshal the modified array of Shape objects back to JSON
	outgoingJson, _ := json.Marshal(shapes, jsontext.WithIndent("  "))
	fmt.Println(string(outgoingJson)) // Show what the JSON looks like
	// Not shown: Send outgoingJson back to Web Service

	// Output:
	// [
	//   {
	//     "color": "red",
	//     "kind": "circle",
	//     "radius": 1
	//   },
	//   {
	//     "color": "green",
	//     "kind": "rectangle",
	//     "width": 15,
	//     "height": 15
	//   },
	//   {
	//     "color": "blue",
	//     "kind": "rectangle",
	//     "width": 5,
	//     "height": 5
	//   }
	// ]
	// ----------
	// [
	//   {
	//     "color": "white",
	//     "kind": "circle",
	//     "radius": 2
	//   },
	//   {
	//     "color": "green",
	//     "kind": "rectangle",
	//     "width": 10,
	//     "height": 10
	//   },
	//   {
	//     "color": "blue",
	//     "kind": "circle",
	//     "radius": 5
	//   }
	// ]
}

// ********** THE CODE BELOW SHOWS HOW TO DEFINE A SUM TYPE ********** //

// At app initialization, panic if any of shape's projection structs don't match
var _ = sumtype.Caster[shape]{}.ValidateStructFields(true, Shape{}, CircleShape{}, RectangleShape{})

const (
	// CircleShapeKind is the kind for circle shapes
	CircleShapeKind ShapeKind = "circle"

	// RectangleShapeKind is the kind for rectangle shapes
	RectangleShapeKind ShapeKind = "rectangle"
)

type (
	// ShapeKind is the discriminator indicating which type of Shape
	ShapeKind string

	// shape is package-private and used for (un)marshaling (all data fields are public).
	shape struct {
		// shapeCaster MUST be 1st field, unexported & embedded for method "inheritance"
		shapeCaster

		// Color is the color of the shape (shared by all shapes)
		Color *string `json:"color,omitempty"`

		// Kind is the discriminator indicating which type of Shape (shared by all shapes)
		Kind *ShapeKind `json:"kind,omitempty"`

		// Radius is the radius of a circle shape
		Radius *int `json:"radius,omitempty"`

		// Width is the width of a rectangle shape
		Width *int `json:"width,omitempty"`

		// Height is the height of a rectangle shape
		Height *int `json:"height,omitempty"`
	}

	// Shape is public and exposes fields common to all shape kinds
	Shape struct {
		// shapeCaster MUST be 1st field, unexported & embedded for method "inheritance"
		shapeCaster

		// Color is the color of the shape (shared by all shapes)
		Color *string

		// Kind is the discriminator indicating which type of Shape (shared by all shapes)
		Kind *ShapeKind

		// radius is the radius of a circle shape
		_ *int

		// width is the width of a rectangle shape
		_ *int

		// height is the height of a rectangle shape
		_ *int
	}

	// CircleShape is public and exposes fields related to a circle kind.
	CircleShape struct {
		// shapeCaster MUST be 1st field, unexported & embedded for method "inheritance"
		shapeCaster

		// Color is the color of the shape (shared by all shapes)
		Color *string

		// Kind is the discriminator indicating which type of Shape (shared by all shapes)
		Kind *ShapeKind

		// Radius is the radius of a circle shape
		Radius *int

		// width is the width of a rectangle shape
		_ *int

		// height is the height of a rectangle shape
		_ *int
	}

	// RectangleShape is public and exposes fields related to a rectangle kind.
	RectangleShape struct {
		// shapeCaster MUST be 1st field, unexported & embedded for method "inheritance"
		shapeCaster

		// Color is the color of the shape (shared by all shapes)
		Color *string

		// Kind is the discriminator indicating which type of Shape (shared by all shapes)
		Kind *ShapeKind

		// radius is the radius of a circle shape
		_ *int

		// Width is the width of a rectangle shape
		Width *int

		// Height is the height of a rectangle shape
		Height *int
	}

	// shapeCaster provides methods to cast between *shape and its variants. The 1st field of shape
	// and all its variants is an unexported shapeCaster whose underlying type is sumtype.Caster[shape].
	// NOTE: This also hides sumtypes.Caster's MarshalJSON/UnmarshalJSON/String methods so they
	// cannot be called directly on shape variants.
	shapeCaster sumtype.Caster[shape]
)

// RULES: String & MarshalJSON require by-val receiver, UnmarshalJSON requires by-ref receiver

// String returns a readable JSON representation of the shape
func (s Shape) String() string { return (&s).caster().String() }

// MarshalJSON marshals the shape to JSON
func (s Shape) MarshalJSON() ([]byte, error) { return (&s).caster().MarshalJSON() }

// UnmarshalJSON unmarshals JSON data to the shape
func (s *Shape) UnmarshalJSON(data []byte) error { return s.caster().UnmarshalJSON(data) }


// String returns a readable JSON representation of the shape
func (s CircleShape) String() string { return (&s).caster().String() }

// MarshalJSON marshals the CircleShape to JSON
func (s CircleShape) MarshalJSON() ([]byte, error) { return (&s).caster().MarshalJSON() }

// UnmarshalJSON unmarshals JSON data to the CircleShape
func (s *CircleShape) UnmarshalJSON(data []byte) error { return s.caster().UnmarshalJSON(data) }


// String returns a readable JSON representation of the shape
func (s RectangleShape) String() string { return (&s).caster().String() }

// MarshalJSON marshals the RectangleShape to JSON
func (s RectangleShape) MarshalJSON() ([]byte, error) { return (&s).caster().MarshalJSON() }

// UnmarshalJSON unmarshals JSON data to the RectangleShape
func (s *RectangleShape) UnmarshalJSON(data []byte) error { return s.caster().UnmarshalJSON(data) }

// RULES: Methods that cast a pointer from 1 type to another, require by-ref receiver (XxxCaster methods).

// caster returns shapeCaster's underlyting sumtype.Caster to access its helper methods.
func (c *shapeCaster) caster() *sumtype.Caster[shape] { return (*sumtype.Caster[shape])(c) }

// json casts the pointer *c to *shape, the JSONable type (ALL JSON fields are public).
func (c *shapeCaster) json() *shape { return c.caster().Json() }

// ensureKind ensures that the current shape kind matches the specified kind; it panics if not.
func (c *shapeCaster) ensureKind(kind ShapeKind) {
	if c.json().Kind == nil {
		panic(fmt.Sprintf("can't cast shape from Kind=nil to Kind=%s", kind))
	}
	if *c.json().Kind != kind {
		panic(fmt.Sprintf("can't cast shape from Kind=%v to Kind=%s", *c.json().Kind, kind))
	}
}

// Shape casts a *XxxShape to the common *Shape
func (c *shapeCaster) Shape() *Shape { return sumtype.Cast[Shape](c.caster()) }

// Circle casts any *XxxShape to a *CircleShape; it panics if Kind != CircleShapeKind.
func (c *shapeCaster) Circle() *CircleShape {
	c.ensureKind(CircleShapeKind)
	return sumtype.Cast[CircleShape](c.caster())
}

// Rectangle casts any *XxxShape to a *RectangleShape; it panics if Kind != RectangleShapeKind.
func (c *shapeCaster) Rectangle() *RectangleShape {
	c.ensureKind(RectangleShapeKind)
	return sumtype.Cast[RectangleShape](c.caster())
}

// SetCircle casts any *XxxShape to a *CircleShape
func (c *shapeCaster) SetCircle() *CircleShape {
	s := c.Shape()
	*s.Kind = CircleShapeKind
	c.caster().ZeroNonKindFields(s)
	return s.Circle()
}

// SetRectangle casts any *XxxShape to a *RectangleShape
func (c *shapeCaster) SetRectangle() *RectangleShape {
	s := c.Shape()
	*s.Kind = RectangleShapeKind
	c.caster().ZeroNonKindFields(s)
	return s.Rectangle()
}

// String returns a readable JSON representation of the shape
func (c *shapeCaster) String() string {
	j, _ := json.Marshal(c.json(), jsontext.WithIndent("  "))
	return string(j)
}
