package wraplog

import (
	"fmt"
	"net/url"
	"testing"
)

func TestGetCtx(t *testing.T) {
	r1 := url.QueryEscape("a=c")
	fmt.Printf("r1: %+v\n", r1)

	r2 := GetCtx(r1, "d=e")
	fmt.Printf("r2: %+v\n", r2)

	r3 := r2
	// r3 := r1 + r2
	// r3 := GetCtx(r2)
	// fmt.Printf("r3: %+v\n", r3)
	decodeR3, err := url.QueryUnescape(r3)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("decodeR3: %+v\n", decodeR3)
}
