package astparser

import (
	"github.com/fatih/structtag"
)

// Parser ast parser
type Parser struct {
	// Imports import list
	Imports []string
	// Structs struct list
	Structs []*Struct
	// Values variable list and constant list
	Values []*Value
	// Interfaces interface list
	Interfaces []*Interface

	valueMap     map[string]*Value
	structMap    map[string]*Struct
	interfaceMap map[string]*Interface
}

func (parser *Parser) GetStruct(name string) *Struct {
	return parser.structMap[name]
}

func (parser *Parser) GetInterface(name string) *Interface {
	return parser.interfaceMap[name]
}

func (parser *Parser) GetValue(name string) *Value {
	return parser.valueMap[name]
}

func (parser *Parser) init() {
	parser.valueMap = make(map[string]*Value)
	for _, v := range parser.Values {
		parser.valueMap[v.Name] = v
	}

	parser.structMap = make(map[string]*Struct)
	for _, s := range parser.Structs {
		parser.structMap[s.Name] = s
	}

	parser.interfaceMap = make(map[string]*Interface)
	for _, i := range parser.Interfaces {
		parser.interfaceMap[i.Name] = i
	}
}

// Interface go interface
type Interface struct {
	Name, Doc string
	Functions []*Function
}

// Function go func
type Function struct {
	Name, Doc string
	Anonymous bool
	// InputParams input params
	InputParams []string
	// ReturnParams return values
	ReturnParams []string
}

// Struct go struct
type Struct struct {
	Name, Doc string
	// Fields struct field list
	Fields []*StructField
}

func (st *Struct) RangeFieldTags(fn func(tags *structtag.Tags, field *StructField, anonymous bool) bool) {
	for _, v := range st.Fields {
		if !fn(v.Tags, v, v.Anonymous) {
			break
		}
		v.Tag = v.Tags.String()
	}
}

// StructField go struct field
type StructField struct {
	Name, Doc string
	Anonymous bool
	Type      string
	// Tag fully tag content
	Tag string
	// Tags parsed tag list
	Tags *structtag.Tags
}

// Value go value, contains constants, variables
type Value struct {
	Name, Doc string
	// Kind constant or variable
	Kind string
	Val  string
}
