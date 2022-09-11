# astparser

Easily parse go code, and you can get imports, structs, interfaces, functions, variables, constants.

## Installation

Using this package requires a working Go environment. [See the install instructions for Go](http://golang.org/doc/install.html).

This package requires a modern version of Go supporting modules: [see the go blog guide on using Go Modules](https://blog.golang.org/using-go-modules).

### Using package

```bash
go get github.com/kevinma2010/astparser
```

```go
...
import (
 "github.com/kevinma2010/astparser"
)
...
```

## Usage/Examples

```go
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
if err != nil {
  log.Fatalf("parse code failure, reason: %s", err)
}

log.Printf("%+v\n", parser.Values[0])
log.Printf("%+v\n", parser.Structs[0])
log.Printf("%+v\n", parser.Interfaces[0])
log.Printf("%+v\n", parser.Imports)
```
