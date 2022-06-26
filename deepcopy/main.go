package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type ASet struct {
	UID   string   `json:"uid,omitempty"`
	Items []*ItemA `json:"member,omitempty"`
}

type ItemA struct {
	T int
}

func main() {
	a := &ASet{
		UID:   "aaa",
		Items: []*ItemA{{T: 0}},
	}
	b := new(ASet)
	DeepCopy(a, b)
	b.UID = "bbb"
	b.Items = append(b.Items, &ItemA{T: 1})
	fmt.Printf("%#v\n%#v\n", a, b)
}

func DeepCopy(src, dst interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
