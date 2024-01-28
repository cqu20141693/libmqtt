package examples

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/decoder"
	"strings"
	"testing"
)

func TestSonic(t *testing.T) {
	toMap()
	toStruct()
	toNestedStruct()
	toArrayStruct()
	encodeWriter()
	decodeReader()
	decodeNumber()
}

func decodeNumber() {

	var input = `1`
	var data interface{}

	// default float64
	dc := decoder.NewDecoder(input)
	_ = dc.Decode(&data) // data == float64(1),默认
	// use json.Number
	dc = decoder.NewDecoder(input)
	dc.UseNumber()
	_ = dc.Decode(&data) // data == json.Number("1")
	// use int64
	dc = decoder.NewDecoder(input)
	dc.UseInt64()
	_ = dc.Decode(&data) // data == int64(1)

	root, _ := sonic.GetFromString(input)
	// Get json.Number
	jn, _ := root.Number()
	jm, _ := root.InterfaceUseNumber() // jn == jm
	jm = jm.(json.Number)
	// Get float64
	fn, _ := root.Float64()
	i, _ := root.Interface()
	fm := i.(float64) // jn == jm
	fmt.Println(jn, jm, fn, fm)

}

func decodeReader() {
	fmt.Println("test sonic decodeReader")
	var o = map[string]interface{}{}
	var r = strings.NewReader(`{"a":"b"}{"1":"2"}`)
	var dec = sonic.ConfigDefault.NewDecoder(r)
	_ = dec.Decode(&o)
	_ = dec.Decode(&o)
	fmt.Printf("%+v", o)
	// Output:
	// map[1:2 a:b]
}

func encodeWriter() {
	fmt.Println("test sonic encodeWriter")
	var o1 = map[string]interface{}{
		"a": "b",
	}
	var o2 = 1
	var w = bytes.NewBuffer(nil)
	var enc = sonic.ConfigDefault.NewEncoder(w)
	_ = enc.Encode(o1)
	_ = enc.Encode(o2)
	fmt.Println(w.String())
	// Output:
	// {"a":"b"}
	// 1
}

func toArrayStruct() {
	fmt.Println("test sonic toArrayStruct")
	var jsonStr = `[{"name": "Harper", "age": 33}, {"name": "Bella", "age": 34}]`

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	var persons []Person
	err := sonic.Unmarshal([]byte(jsonStr), &persons)
	if err != nil {
		fmt.Println("解析失败：", err)
	}

	for _, person := range persons {
		fmt.Println(person.Name, person.Age)
	}
}

func toNestedStruct() {
	fmt.Println("test sonic toNestedStruct")
	var jsonStr = `{"name": "Harper", "age": 30, "address": {"city": "HaiDian Beijing", "country": "China"}}`

	type Address struct {
		City    string `json:"city"`
		Country string `json:"country"`
	}

	type Person struct {
		Name    string  `json:"name"`
		Age     int     `json:"age"`
		Address Address `json:"address"`
	}

	var person Person
	err := sonic.Unmarshal([]byte(jsonStr), &person)
	if err != nil {
		fmt.Println("解析失败：", err)
	}

	fmt.Println(person.Name, person.Age, person.Address.City, person.Address.Country)
}

func toStruct() {
	fmt.Println("test sonic toStruct")
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	var jsonStr = `{"name": "Harper", "age": 20}`
	var person Person

	err := sonic.Unmarshal([]byte(jsonStr), &person)
	if err != nil {
		fmt.Println("解析失败：", err)
	}

	fmt.Println(person.Name, person.Age)
}

func toMap() {
	var jsonStr = `{"name": "Harper", "age": 18}`
	var data map[string]any

	err := sonic.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		fmt.Println("解析失败：", err)
	}

	fmt.Println(data["name"], data["age"])
}
