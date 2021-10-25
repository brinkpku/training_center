package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
	"time"

	"github.com/micro/go-micro/util/log"
)

type Monit struct {
	Name string
}

func main() {
	tmplStr, err := readFile("template.monit")
	if err != nil {
		log.Fatalf("read file error: %v", err)
	}
	log.Info(tmplStr)
	m := Monit{Name: "test"}
	tmpl, err := template.New("monit").Parse(tmplStr)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.OpenFile("/etc/monit.d/test.monit", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777) //linux 路径
	if err != nil {
		fmt.Printf("open err%s", err)
		return
	}
	if err := tmpl.Execute(f, m); err != nil {
		log.Fatal(err)
	}
}

func readFile(fName string) (tmpl string, err error) {
	f, err := ioutil.ReadFile(fName)
	if err != nil {
		return
	}
	tmpl = string(f)
	return
}

// WriteFile ...
func WriteFile() {
	f, err := os.OpenFile("/etc/monit.d", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0750) //linux 路径
	if err != nil {
		fmt.Printf("open err%s", err)
		return
	}
	defer f.Close() //资源必须释放,函数刚要返回之前延迟执行
	for i := 0; i < 100; i++ {
		ret, err2 := f.WriteString("s")
		if err2 != nil {
			fmt.Printf("write err:\n%s", err2)
			return
		}
		fmt.Println(ret)
		time.Sleep(2 * time.Second) //延时2秒
	}

}
