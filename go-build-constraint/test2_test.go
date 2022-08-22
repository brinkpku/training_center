// +build test_all,!t1

package go_build_constraint

import (
	"fmt"
	"testing"
)

func TestOtherCase(t *testing.T) {
	fmt.Println("testing:", t.Name())
}
