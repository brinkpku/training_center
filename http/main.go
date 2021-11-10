package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	resp, err := HTTPGet(nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp)
}

// HTTPGet ...
func HTTPGet(req interface{}) (resp interface{}, err error) {
	// var data []byte
	// data, err = json.Marshal(req)
	// if err != nil {
	// 	return
	// }
	client := &http.Client{
		Timeout: time.Second * time.Duration(5),
	}
	for i := 0; i < 3; i++ {
		var r *http.Request
		if r, err = http.NewRequest("GET", "http://localhost:8888/v1/status/memory", nil); err != nil {
			return
		}
		var res *http.Response
		// r.Header.Add("Content-Type", "application/json")
		r.Header.Add("accept", "application/json")
		r.Header.Add("Authorization", "Basic YWRtaW4xMjM0OmFkbWluMTIzNA==")
		if res, err = client.Do(r); err != nil {
			return
		}
		defer res.Body.Close()
		var body []byte
		if body, err = ioutil.ReadAll(res.Body); err != nil {
			return
		}
		log.Println(string(body))
		if err = json.Unmarshal(body, resp); err != nil {
			return
		}
		if res.StatusCode != 200 {
			continue
		}
		return
	}
	return
}
