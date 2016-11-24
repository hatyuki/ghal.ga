package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/apex/go-apex"
	"github.com/fujiwara/ridge"
	"github.com/gorilla/mux"
	"github.com/hatyuki/go-ghl"
)

const (
	Home = "https://hatyuki.github.io/ghal.ga/"
)


var router = mux.NewRouter()
var token = os.Getenv("GITHUB_TOKEN")

func init () {
	router.HandleFunc("/{owner}/{repo}", handleRequest)
	router.HandleFunc("/", handleRoot)
}

func main () {
	if os.Getenv("APEX_FUNCTION_NAME") == "" {
		fmt.Println("Accepting connections at https://localhost:8080/")
		fmt.Fprintf(os.Stderr, "%v\n", http.ListenAndServe(":8080", router))
	}

	apex.HandleFunc(handler)
}

func handler (event json.RawMessage, ctx *apex.Context) (interface{}, error) {
	fmt.Errorf("%v\n", event)
	fmt.Errorf("%v\n", ctx)

	request, err := ridge.NewRequest(event)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	response := ridge.NewResponseWriter()
	router.ServeHTTP(response, request)

	return response.Response(), nil
}

func handleRoot(response http.ResponseWriter, request *http.Request) {
    http.Redirect(response, request, Home, http.StatusMovedPermanently)
}

func handleRequest(response http.ResponseWriter, request *http.Request) {
	target := request.URL.Path[1:]
	query := request.URL.Query()
	os := query.Get("os")
	arch := query.Get("arch")

	if os == "" || arch == "" {
		http.Error(response, "OS and Arch is required", http.StatusBadRequest)
		return
	}

	locator := ghl.NewClient(token)
	location, err := locator.GetDownloadURL(target, os, arch)
	if err != nil {
		http.NotFound(response, request)
		return
	}

	var status int
	if strings.Contains(target, "@") {
		status = http.StatusMovedPermanently
	} else {
		status = http.StatusFound
	}

	http.Redirect(response, request, location, status)
}
