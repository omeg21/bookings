package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M)  {
	
	
	
	os.Exit(m.Run())
}

type myController struct{

}

func (mh *myController) ServeHTTP(w http.ResponseWriter,r *http.Request) {
	
}