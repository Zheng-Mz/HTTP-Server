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

func main() {
    //init
    count = 0
    var dir string
    var port string
    /*./exe -dir value*/
    flag.StringVar(&dir, "d", ".", "the directory to serve files from. Defaults to the current dir")
    flag.StringVar(&port, "p", "80", "set Http-Server port.")
    flag.Parse()
    fmt.Println("dir: ", dir)
    fmt.Println("port: ", port)

    r := mux.NewRouter()
    r.Use(loggingMiddleware)

    r.HandleFunc("/test", testHandle).Methods("GET")
    r.HandleFunc("/test/pathVar/{key}", pathVarHandle).Methods("GET")

    h2 := &http2.Server{}
    handler := h2c.NewHandler(r, h2)
    srv := &http.Server{
        Addr: fmt.Sprintf(":%s", port),
        TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
        Handler:      handler}
    log.Fatal(srv.ListenAndServe())
}
