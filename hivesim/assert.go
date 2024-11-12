package hivesim

import "fmt"

func (t *T) Nil(err error, msgAndArgs ...interface{}) {
	if err != nil {
		t.Fatal(fmt.Sprintf("Expected nil, but got: %#v", err), msgAndArgs)
	}
}

func (t *T) True(ok bool, msgAndArgs ...interface{}) {
	if !ok {
		t.Fatal("Should be true", msgAndArgs)
	}
}
