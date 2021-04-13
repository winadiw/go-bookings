package main

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Sprintln("Here!")
	os.Exit(m.Run())
}

type myHandler struct{}

func (mh *myHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

}
