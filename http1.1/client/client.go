package main

import (
    //"encoding/json"
    "fmt"
    //"io/ioutil"
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
    TEST_REQUEST_NUM = 5000
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
    }
    req, err := http.NewRequest("GET", url, strings.NewReader(parama))
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    resp, err := client.Do(req)
    if err != nil {
        fmt.Printf("err 1: %v\n", err)
        return 0
    }
    defer resp.Body.Close()
    return resp.StatusCode
    /*
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Printf("err 2: %v\n", err)
        return 0
    }

    jsonStr := string(body)
    fmt.Println("jsonStr", jsonStr)
    var dat map[string]string
    if err := json.Unmarshal([]byte(jsonStr), &dat); err == nil {
        fmt.Println("token", dat["token"])
    } else {
        fmt.Printf("err 3: json str to struct error")
        return nil
    }
    return dat
    */
}
