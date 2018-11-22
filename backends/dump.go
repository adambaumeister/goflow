package backends

import (
	"fmt"
	"github.com/adamb/goflow/fields"
	"strings"
)

type Dump struct{}

func (b Dump) Init() {}

func (b Dump) Add(values map[uint16]fields.Value) {
	var sl []string
	for t, val := range values {
		sl = append(sl, fmt.Sprintf("(%v)%v", t, val.ToString()))
	}
	fmt.Printf("%v", strings.Join(sl, " : ")+"\n")
}
