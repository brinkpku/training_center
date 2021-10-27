package main

import (
	"fmt"
	"net/http"

	"github.com/beevik/etree"
)

func main() {
	doc := etree.NewDocument()
	res, err := http.Get("http://admin:monit@monithttpd:2812/_status?format=xml")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	// bytes, err:=ioutil.ReadAll(res.Body)
	// if err!=nil
	if _, err := doc.ReadFrom(res.Body); err != nil {
		panic(err)
	}
	ele := doc.FindElement("//service[@type='3']/name[text()='operate-manage-service']/../status")
	if ele != nil {
		fmt.Println("status:", ele.Text())
	} else {
		fmt.Println()
	}
}
