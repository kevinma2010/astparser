# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go library called `astparser` that parses Go source code and extracts structural information including imports, structs, interfaces, functions, variables, and constants. It's a pure Go library with minimal dependencies, designed to provide an easy-to-use API for analyzing Go code structure.

## Common Development Commands

### Building and Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v

# Run a specific test
go test -v -run TestParse

# Run tests with coverage
go test -cover

# Generate coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Dependency Management
```bash
# Download dependencies
go mod download

# Tidy up dependencies
go mod tidy

# Verify dependencies
go mod verify
```

### Development Workflow
```bash
# Format code
go fmt ./...

# Run static analysis
go vet ./...

# Install the library locally for testing
go install
```

## Architecture Overview

The library consists of three main files:

1. **parser.go**: Contains the main parsing logic using Go's `go/ast` and `go/parser` packages. Key functions:
   - `Parse(code string)`: Parses Go code from a string
   - `ParseBytes(code []byte)`: Parses Go code from bytes
   - Internal parsing functions for different AST node types

2. **core.go**: Defines the core data structures:
   - `Parser`: Main struct containing parsed results with arrays for imports, structs, values, and interfaces
   - `Struct`: Represents a Go struct with fields and documentation
   - `Interface`: Represents a Go interface with method signatures
   - `Value`: Represents variables and constants
   - `Field`: Represents struct fields with tags support
   - `Function`: Represents functions and methods

3. **parse_test.go**: Contains unit tests using the testify framework

## Key Design Patterns

1. **AST Visitor Pattern**: The parser walks the Go AST tree and extracts information from different node types
2. **Data Collection**: All parsed elements are collected into arrays and also indexed in maps for quick lookup
3. **Documentation Preservation**: The parser preserves Go doc comments for all elements
4. **Struct Tag Parsing**: Uses the `fatih/structtag` library to parse struct field tags

## Testing Approach

- Tests use the standard Go testing framework with testify assertions
- Main test file `parse_test.go` contains comprehensive tests for parsing various Go constructs
- Tests verify both the parsing functionality and the data structure population

## API Usage Pattern

```go
// Parse from string
parser, err := astparser.Parse(goCode)

// Parse from bytes
parser, err := astparser.ParseBytes([]byte(goCode))

// Access parsed data
imports := parser.Imports
structs := parser.Structs
interfaces := parser.Interfaces
values := parser.Values

// Lookup by name
myStruct := parser.GetStruct("MyStruct")
myInterface := parser.GetInterface("MyInterface")
myValue := parser.GetValue("MyConstant")
```

## Important Notes

- The library only parses syntactically valid Go code
- It focuses on type definitions and declarations, not implementation details
- Function bodies are not analyzed, only signatures are extracted
- The parser preserves the original documentation comments from the source code