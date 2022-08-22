package assert

import (
	"testing"
)

type Person struct {
	name string
	age  int
}

var paul = Person{
	name: "Paul",
	age:  32,
}

var peter = Person{
	name: "Peter",
	age:  21,
}

func TestEqual(t *testing.T) {
	Equal(t, paul, paul, "Paul equals paul")
}

func TestNotEqual(t *testing.T) {
	NotEqual(t, paul, peter, "Paul does not equal peter")
}

func TestNotEqualNil(t *testing.T) {
	NotEqual(t, paul, nil, "Paul is not nil")
}

func TestNilEqualNil(t *testing.T) {
	Equal(t, nil, nil, "Nil is nil")
}
