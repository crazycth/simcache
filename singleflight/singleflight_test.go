package singleflight

import (
	"fmt"
	"testing"
)

func TestGroup(t *testing.T) {
	g := Group{}
	v, err := g.Do("key", func() (interface{}, error) {
		return "bar", nil
	})
	fmt.Println(v, err)
}
