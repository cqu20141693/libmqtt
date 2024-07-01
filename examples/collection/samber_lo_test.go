package collection

import (
	"fmt"
	"github.com/samber/lo"
	"sort"
	"strings"
	"testing"
)

func TestLoSlice(t *testing.T) {

}

func TestMap(t *testing.T) {
	fmt.Println("start test lo map")
	ExampleKeys()
	ExampleValues()
}
func ExampleKeys() {
	kv := map[string]int{"foo": 1, "bar": 2}

	result := lo.Keys(kv)

	sort.StringSlice(result).Sort()
	fmt.Printf("%v \n", result)
	// Output: [bar foo]
}

func ExampleValues() {
	kv := map[string]int{"foo": 1, "bar": 2}

	result := lo.Values(kv)

	sort.IntSlice(result).Sort()
	fmt.Printf("%v\n", result)
	// Output: [1 2]
}

func TestFilter(t *testing.T) {
	strs := lo.Filter[string]([]string{"hello", "good bye", "world", "fuck", "fuck who"}, func(s string, _ int) bool {
		return !strings.Contains(s, "fuck")
	})
	fmt.Println(strs) //[hello good bye world]
}
