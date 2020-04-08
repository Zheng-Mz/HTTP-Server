package main

import (
    "crypto/tls"
    "flag"
    "fmt"
    "github.com/gorilla/mux"
    "golang.org/x/net/http2"
    "golang.org/x/net/http2/h2c"
    "log"
    "net/http"
    "time"
    "strings"
)

var count int

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Do stuff here
        log.Println(r.RequestURI)
        // Call the next handler, which can be another middleware in the chain, or the final handler.
        next.ServeHTTP(w, r)
    })
}

func testHandle(w http.ResponseWriter, r *http.Request) {
    count++
    fmt.Printf("Get test request, count: %d\n", count)
    time.Sleep(time.Second*1)
    //time.Sleep(time.Millisecond*500)
    //fmt.Println("test return.")
    w.WriteHeader(http.StatusOK)
    return
}

func hostHandle(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("other host: %v\n", r.Host)
    w.WriteHeader(http.StatusOK)
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

func main() {
    //init
    count = 0
    var dir string
    /*./exe -dir value*/
    flag.StringVar(&dir, "dir", ".", "the directory to serve files from. Defaults to the current dir")
    flag.Parse()
    fmt.Println("dir: ", dir)

    r := mux.NewRouter()
    r.Use(loggingMiddleware)
    //authentication
    amw := authenticationMiddleware{}
    amw.Populate()
    r.Use(amw.Middleware)

    r.HandleFunc("/test", testHandle).Methods("GET")
    r.HandleFunc("/test/sub", hostHandle).Methods("GET").Schemes("http").Host("127.0.0.1")
    r.HandleFunc("/test/pathVar/{key}", pathVarHandle).Methods("GET")

    s := r.Host("172.0.10.91").Subrouter()  //bind local host
    s.HandleFunc("/test/sub", hostSubHandle).Methods("GET")

    s1 := r.PathPrefix("/other").Subrouter()
    s1.HandleFunc("/test", otherTestHandle).Methods("GET")

    // This will serve files under http://localhost:80/static/<filename>
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

    //http://localhost:80/test/quer?filter=file1
    r.HandleFunc("/test/quer", testQuerHandle).Methods("GET").Queries("filter", "{filter}")

    //The Walk function on mux.Router can be used to visit all of the routes that are registered on a router.
    r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
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
    })

    h2 := &http2.Server{}
    handler := h2c.NewHandler(r, h2)
    srv := &http.Server{Addr: ":80",
        TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
        Handler:      handler}
    srv.ListenAndServe()
    fmt.Println("http server stopped listening")
}
