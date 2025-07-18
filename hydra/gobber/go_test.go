package gobber

import (
	"reflect"
	"testing"
)

// Serialize and Deserialize functions from the previous example

// A test struct to demonstrate nested serialization
type TestStruct struct {
	Name   string
	Values []int
	Nested struct {
		Flag    bool
		Message string
	}
}

func TestSerialization(t *testing.T) {
	// Test cases
	cases := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "Slice of ints",
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name: "Nested struct",
			input: TestStruct{
				Name:   "Test",
				Values: []int{5, 4, 3, 2, 1},
				Nested: struct {
					Flag    bool
					Message string
				}{Flag: true, Message: "Nested message"},
			},
			expected: TestStruct{
				Name:   "Test",
				Values: []int{5, 4, 3, 2, 1},
				Nested: struct {
					Flag    bool
					Message string
				}{Flag: true, Message: "Nested message"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Serialize the input
			serializedData, err := Serialize(tc.input)
			if err != nil {
				t.Fatalf("Failed to serialize %v: %v", tc.input, err)
			}

			// Prepare a variable of the expected type
			resultPtr := reflect.New(reflect.TypeOf(tc.expected)).Interface()

			// Deserialize into the new variable
			if err := Deserialize(serializedData, resultPtr); err != nil {
				t.Fatalf("Failed to deserialize into %T: %v", resultPtr, err)
			}

			// Dereference the pointer to get the result
			result := reflect.ValueOf(resultPtr).Elem().Interface()

			// Compare the deserialized data to the expected data
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Deserialized data (%v) does not match expected data (%v)", result, tc.expected)
			}
		})
	}
}
