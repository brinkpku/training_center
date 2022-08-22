//+build !tag

package go_build_constraint

import (
	"fmt"
	"testing"
)

func TestNoTagCase(t *testing.T) {
	fmt.Println("testing:", t.Name())
}
