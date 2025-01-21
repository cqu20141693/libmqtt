package examples

import (
	"fmt"
	"testing"
)

type Person struct {
	Name string
	Age  int
}

func NewPerson(name string, age int) Person {

	return Person{Name: name, Age: age}
}
func TestMapBasic(t *testing.T) {
	deviceMap := make(map[string]string)
	deviceMap["name"] = "wtieam"

	m := map[string]Person{}
	person := m["witeam"]
	fmt.Println(person)
	delete(m, "witeam")
}
