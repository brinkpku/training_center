package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

type Monit struct {
	Name string
}

func main() {
	tmplStr, err := readFile("template.monit")
	if err != nil {
		log.Fatalf("read file error: %v", err)
	}
	log.Println(tmplStr)
	m := Monit{Name: "test"}
	tmpl, err := template.New("monit").Parse(tmplStr)
	if err != nil {
		log.Fatal(err)
	}
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}
	for _, fi := range files {
		log.Println("file name:", fi.Name(), fi.IsDir())
	}
	if err = os.Symlink("../encode", "slink"); err != nil {
		log.Fatal(err)
	}
	f, err := os.OpenFile("/etc/monit.d/test.monit", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777) //linux 路径
	if err != nil {
		log.Fatal(err)
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

//复制整个文件夹或单个文件

func Copy(from, to string) (err error) {
	f, err := os.Stat(from)
	if err != nil {
		return
	}
	if f.IsDir() { //from是文件夹，那么定义to也是文件夹
		if list, e := ioutil.ReadDir(from); e == nil {
			for _, item := range list {
				if e = Copy(filepath.Join(from, item.Name()), filepath.Join(to, item.Name())); e != nil {
					return e
				}
			}
		}
	} else { //from是文件，那么创建to的文件夹
		p := filepath.Dir(to)
		if _, e := os.Stat(p); e != nil {
			if e = os.MkdirAll(p, 0777); e != nil {
				return e
			}
		}
		//读取源文件
		file, e := os.Open(from)
		if e != nil {
			return e
		}
		defer file.Close()
		bufReader := bufio.NewReader(file)
		// 创建一个文件用于保存
		out, e := os.Create(to)
		if e != nil {
			return e
		}
		defer out.Close()
		// 然后将文件流和文件流对接起来
		_, e = io.Copy(out, bufReader)
		if e != nil {
			return e
		}
	}
	return
}
