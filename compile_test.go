package pi2go

import (
	"fmt"
	"os"
	"testing"
)

func TestCompile(t *testing.T) {
	//	for _, t_ := range text {
	Compile(text[0], os.Stdout)
	fmt.Println("--------------------------")
	//	}
}
