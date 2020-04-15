package main

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/net/http2"
	"io/ioutil"
	"log"
	"net"
	"net/http"
    "strings"
    "sync"
    "time"
)

var lock_1 sync.RWMutex

var count int
var successCt int
var sleep_cnt int
const (
    TEST_REQUEST_NUM = 1
    HTTP_REQ_TIMEOUT_SECOND = 3
)

func main() {
    count = 0
    successCt = 0
    sleep_cnt = 0

    url := "http://127.0.0.1:80/test"
    for i:=0; i<TEST_REQUEST_NUM; i++ {
        go func(){
            prama := fmt.Sprintf("data=%d", i)
            respCode := GetTokenReq(url, prama)
	    lock_1.RLock()
	    count++
	    if respCode != 0 {
		successCt++
	        fmt.Printf("count: %d, resp: %v\n", count, respCode)
	    }
	    lock_1.RUnlock()
	}()
    }

    for (count != TEST_REQUEST_NUM) && (sleep_cnt != 6) {
        time.Sleep(time.Second)
	sleep_cnt++
    }
    fmt.Printf("count: %d, success: %d, %d; Sleep_cnt: %d\n", count, successCt, (successCt*100)/count, sleep_cnt)
}

func GetTokenReq(url string, parama string) int {
    client := &http.Client{
        Timeout: time.Second*HTTP_REQ_TIMEOUT_SECOND,
		Transport: &http2.Transport{
			AllowHTTP: true, //充许非加密的链接
			// TLSClientConfig: &tls.Config{
			//     InsecureSkipVerify: true,
			// },
			DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(netw, addr)
			},
		},
    }
    req, err := http.NewRequest("GET", url, strings.NewReader(parama))
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    resp, err := client.Do(req)
    if err != nil {
		log.Fatal(err)
    }
    defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("resp StatusCode:", resp.StatusCode)
		return resp.StatusCode
	}

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
		log.Fatal(err)
    }

	fmt.Println("resp.Body: ", string(body))
    return resp.StatusCode
}
