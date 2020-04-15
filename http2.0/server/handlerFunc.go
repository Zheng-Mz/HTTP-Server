package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"time"
)

func testHandle(w http.ResponseWriter, r *http.Request) {
	count++
	fmt.Printf("Get test request, count: %d\n", count)
	time.Sleep(time.Second*1)
	//time.Sleep(time.Millisecond*500)
	//fmt.Println("test return.")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "count=%v", count)
	return
}

func hostHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("other host: %v\n", r.Host)
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "host: %v", r.Host)
	return
}

func pathVarHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Printf("Key: %s\n", vars["key"])

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Key: %v\n", vars["key"])
}

func hostSubHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("sub host: %s\n", r.Host)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "host: %v\n", r.Host)
}

func otherTestHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("other path: %s\n", r.URL.Path)
	w.WriteHeader(http.StatusOK)
}

func testQuerHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("queries: %v\n", r.URL.RawQuery)
	w.WriteHeader(http.StatusOK)
}

func walkFunc(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	pathTemplate, err := route.GetPathTemplate()
	if err == nil {
		fmt.Println("ROUTE:", pathTemplate)
	}
	pathRegexp, err := route.GetPathRegexp()
	if err == nil {
		fmt.Println("Path regexp:", pathRegexp)
	}
	queriesTemplates, err := route.GetQueriesTemplates()
	if err == nil {
		fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
	}
	queriesRegexps, err := route.GetQueriesRegexp()
	if err == nil {
		fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
	}
	methods, err := route.GetMethods()
	if err == nil {
		fmt.Println("Methods:", strings.Join(methods, ","))
	}
	fmt.Println()
	return nil
}