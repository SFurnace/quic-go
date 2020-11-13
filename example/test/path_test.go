package example_test

import (
	"fmt"
	"os"
	"testing"
)

func TestArgs0(t *testing.T) {
	fmt.Println(os.Args)
}

func TestWd(t *testing.T) {
	fmt.Println(os.Getwd())
}
