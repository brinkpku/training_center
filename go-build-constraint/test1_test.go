//+build test_all t1

package go_build_constraint

import (
	"fmt"
	"testing"
)

func TestNormal(t *testing.T) {
	fmt.Println("testing:", t.Name())
}
