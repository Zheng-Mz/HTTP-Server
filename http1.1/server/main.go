package main

import (
    "flag"
    "fmt"
    "github.com/gorilla/mux"
    "log"
    "net/http"
    "time"
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
    var dir, port, certFile, keyFile string
    /*./exe -dir value*/
    flag.StringVar(&dir, "d", ".", "the directory to serve files from. Defaults to the current dir")
    flag.StringVar(&port, "p", "80", "set Http-Server port.")
    flag.StringVar(&certFile, "crt", "", "set cert file.")
    flag.StringVar(&keyFile, "key", "", "set key file.")
    flag.Parse()
    fmt.Println("dir: ", dir)
    fmt.Println("port: ", port)

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
    r.Walk(walkFunc)

    srv := &http.Server{
        Handler:      r,
        Addr:         fmt.Sprintf(":%s", port),
        // Good practice: enforce timeouts for servers you create!
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }

    if certFile==""||keyFile=="" {
        log.Fatal(srv.ListenAndServe())
    } else {
        fmt.Println("certFile: ", certFile)
        fmt.Println("keyFile : ", keyFile)
        log.Fatal(srv.ListenAndServeTLS(certFile, keyFile))
    }
}
