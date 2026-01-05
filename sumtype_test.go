package sumtype_test

import (
	"encoding/json/v2"
	"testing"
)

// TestShapeJSONMarshalUnmarshal tests JSON marshaling and unmarshaling for all shape types
func TestShapeJSONMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name         string
		shape        any
		expectedKind ShapeKind
		expectedJSON string
		validateFunc func(t *testing.T, unmarshaled *shape)
	}{
		{
			name: "Circle",
			shape: CircleShape{
				Color:  ptr("red"),
				Kind:   ptr(CircleShapeKind),
				Radius: ptr(50),
			},
			expectedKind: CircleShapeKind,
			validateFunc: func(t *testing.T, unmarshaled *shape) {
				c := unmarshaled.Circle()
				if *c.Radius != 50 {
					t.Errorf("Circle coordinates/radius mismatch: Radius=%d", *c.Radius)
				}
				if *c.Color != "red" {
					t.Errorf("Colors mismatch: Color=%s", *c.Color)
				}
			},
		},
		{
			name: "Rectangle",
			shape: RectangleShape{
				Color:  ptr("green"),
				Kind:   ptr(RectangleShapeKind),
				Width:  ptr(300),
				Height: ptr(150),
			},
			expectedKind: RectangleShapeKind,
			validateFunc: func(t *testing.T, unmarshaled *shape) {
				r := unmarshaled.Rectangle()
				if *r.Width != 300 || *r.Height != 150 {
					t.Errorf("Rectangle dimensions mismatch:  Width=%d, Height=%d", *r.Width, *r.Height)
				}
			},
		},
		{
			name: "Shape",
			shape: Shape{
				Color: ptr("purple"),
				Kind:  ptr(CircleShapeKind),
			},
			expectedKind: CircleShapeKind,
			validateFunc: func(t *testing.T, unmarshaled *shape) {
				s := unmarshaled.Shape()
				if *s.Color != "purple" {
					t.Errorf("FillColor mismatch: expected purple, got %s", *s.Color)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			jsonData, err := json.Marshal(tt.shape)
			if err != nil {
				t.Fatalf("Failed to marshal %s: %v", tt.name, err)
			}

			t.Logf("Marshaled JSON: %s", string(jsonData))

			// Test unmarshaling
			var unmarshaled shape
			err = json.Unmarshal(jsonData, &unmarshaled)
			if err != nil {
				t.Fatalf("Failed to unmarshal %s: %v", tt.name, err)
			}

			// Validate kind
			if *unmarshaled.Kind != tt.expectedKind {
				t.Errorf("Kind mismatch: expected %s, got %s", tt.expectedKind, *unmarshaled.Kind)
			}

			// Run custom validation
			if tt.validateFunc != nil {
				tt.validateFunc(t, &unmarshaled)
			}
		})
	}
}

// TestShapeTypeConversion tests converting between different shape kinds
func TestShapeTypeConversion(t *testing.T) {
	// Start with a rectangle
	rectangle := RectangleShape{
		Color:  ptr("red"),
		Kind:   ptr(RectangleShapeKind),
		Width:  ptr(100),
		Height: ptr(50),
	}

	t.Run("Rectangle to Circle conversion", func(t *testing.T) {
		// Convert to circle
		circle := rectangle.SetCircle()

		// Verify kind changed
		if *circle.Kind != CircleShapeKind {
			t.Errorf("Expected kind %s, got %s", CircleShapeKind, *circle.Kind)
		}

		// Verify common fields preserved
		if *circle.Color != "red" {
			t.Error("Common fields not preserved during conversion")
		}

		// Set circle-specific fields
		circle.Radius = ptr(25)

		// Verify shape-specific fields are accessible
		if *circle.Radius != 25 {
			t.Error("Failed to set circle radius")
		}
	})

	t.Run("Circle to Rectangle conversion", func(t *testing.T) {
		// Start with the circle from previous test
		circle := rectangle.SetCircle()
		circle.Radius = ptr(25)

		// Convert back to rectangle
		newRectangle := circle.SetRectangle()

		// Verify kind changed back
		if *newRectangle.Kind != RectangleShapeKind {
			t.Errorf("Expected kind %s, got %s", RectangleShapeKind, *newRectangle.Kind)
		}

		// Verify common fields preserved
		if *newRectangle.Color != "red" {
			t.Error("Common fields not preserved during conversion")
		}

		// Set rectangle-specific fields
		newRectangle.Width = ptr(200)
		newRectangle.Height = ptr(100)

		// Verify rectangle-specific fields are accessible
		if *newRectangle.Width != 200 || *newRectangle.Height != 100 {
			t.Error("Failed to set rectangle dimensions")
		}
	})
}

// TestDiscriminatedUnionReading tests reading and interpreting discriminated union types
func TestDiscriminatedUnionReading(t *testing.T) {
	testCases := []struct {
		name     string
		jsonData string
		testFunc func(t *testing.T, s *shape)
	}{
		{
			name: "Read Circle JSON",
			jsonData: `{
				"kind": "circle",
				"color": "blue",
				"radius": 30
			}`,
			testFunc: func(t *testing.T, s *shape) {
				if *s.Kind != CircleShapeKind {
					t.Errorf("Expected Circle kind, got %s", *s.Kind)
				}

				// Use the discriminator to access the correct type
				switch *s.Kind {
				case CircleShapeKind:
					circle := s.Circle()
					if *circle.Radius != 30 {
						t.Errorf("Circle properties mismatch:  Radius=%d", *circle.Radius)
					}
				default:
					t.Error("Unexpected shape kind")
				}
			},
		},
		{
			name: "Read Rectangle JSON",
			jsonData: `{
				"kind": "rectangle",
				"color": "green",
				"width": 150,
				"height": 75
			}`,
			testFunc: func(t *testing.T, s *shape) {
				if *s.Kind != RectangleShapeKind {
					t.Errorf("Expected Rectangle kind, got %s", *s.Kind)
				}

				// Use the discriminator to access the correct type
				switch *s.Kind {
				case RectangleShapeKind:
					rectangle := s.Rectangle()
					if *rectangle.Width != 150 || *rectangle.Height != 75 {
						t.Errorf("Rectangle properties mismatch:  Width=%d, Height=%d", *rectangle.Width, *rectangle.Height)
					}
				default:
					t.Error("Unexpected shape kind")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var s shape
			err := json.Unmarshal([]byte(tc.jsonData), &s)
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			tc.testFunc(t, &s)
		})
	}
}

// TestUnsafePointerConversions tests that unsafe pointer conversions work correctly
func TestUnsafePointerConversions(t *testing.T) {
	rectangle := RectangleShape{
		Color:  ptr("red"),
		Kind:   ptr(RectangleShapeKind),
		Width:  ptr(100),
		Height: ptr(50),
	}
	_ = rectangle.Kind
	_ = rectangle.Width
	_ = rectangle.Height

	// Test all conversion methods
	shape := rectangle.Shape()

	// Verify they all point to the same memory
	if shape.Color != rectangle.Color {
		t.Error("Shape conversion failed - FillColor pointers don't match")
	}

	circle := rectangle.SetCircle()
	if circle.Color != rectangle.Color {
		t.Error("Circle conversion failed - FillColor pointers don't match")
	}

	rect := rectangle.SetRectangle()
	if rect.Color != rectangle.Color {
		t.Error("Rectangle conversion failed - FillColor pointers don't match")
	}

	// Test that modifications through one view affect all views
	*rectangle.Color = "green"
	if *shape.Color != "green" {
		t.Error("Modification through rectangle not reflected in shape view")
	}
	if *circle.Color != "green" {
		t.Error("Modification through rectangle not reflected in circle view")
	}
}
