// Package gobber provides functions for serializing and deserializing Go objects using the Gob encoding format.
// Gob encoding is a way to convert Go data structures into a stream of bytes for transmission or storage,
// and then back into data structures again.
// This package offers two main functions: Serialize, which serializes any Go object into a byte slice,
// and Deserialize, which deserializes a byte slice back into a Go object.
// These functions are useful for efficiently transmitting and storing data in Go applications.
package gobber

import (
	"bytes"
	"encoding/gob"
)

// Serialize takes any object and returns a byte slice.
func Serialize(data interface{}) ([]byte, error) {
	// Create a buffer to store the serialized data
	var buffer bytes.Buffer

	// Initialize a new encoder with the buffer
	encoder := gob.NewEncoder(&buffer)

	// Encode the data into the buffer
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}

	// Return the content of the buffer as a byte slice
	return buffer.Bytes(), nil
}

// Deserialize takes a byte slice and a pointer to an object,
// and fills the object with the deserialized data.
func Deserialize(data []byte, obj interface{}) error {
	// Create a buffer and fill it with the serialized data
	buffer := bytes.NewBuffer(data)

	// Initialize a new decoder with the buffer
	decoder := gob.NewDecoder(buffer)

	// Decode the data from the buffer into the provided object
	err := decoder.Decode(obj)
	if err != nil {
		return err
	}

	return nil
}
