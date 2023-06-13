package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurve(t *testing.T)  {
	var myH myController
	
	h := Nosurf(&myH)

	switch v := h.(type) {
	case http.Handler:


		default :
	t.Error(fmt.Sprintf("Type is not HTTP handler but is %t", v) )
	}
	

}

func TestSessionLoad(t *testing.T)  {
	var myH myController
	
	h := SessionLoad(&myH)

	switch v := h.(type) {
	case http.Handler:


		default :
	t.Error(fmt.Sprintf("Type is not HTTP handler but is %t", v) )
	}
	

}