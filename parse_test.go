package astparser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	var code = `package main

import "fmt"

// Greeting greeting interface
type Greeting interface {
	// SayHello say hello
	SayHello()
}

// Dog dog type struct
type Dog struct {
	Name string
}

// SayHello dog say hello
func (dog *Dog) SayHello() {
	fmt.Println("Bow-wow! I am", dog.Name)
}

// FirstName first name
var FirstName = "Kevin"

// LastName last name
const LastName = "Ma"

func main() {
	dog := &Dog{Name: "Buddy"}
	dog.SayHello()
	fmt.Println("The author is", FirstName, LastName)
}`
	parser, err := Parse(code)
	assert.Nil(t, err)
	assert.NotNil(t, parser)
	assert.Equal(t, 2, len(parser.Values))
	assert.Equal(t, 1, len(parser.Structs))
	assert.Equal(t, 1, len(parser.Interfaces))
	assert.Equal(t, []string{"fmt"}, parser.Imports)

	t.Logf("%+v", parser.Values[0])
	t.Logf("%+v", parser.Structs[0])
	t.Logf("%+v", parser.Interfaces[0])
	t.Logf("%+v", parser.Imports)
}
