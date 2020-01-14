/*
  author='du'
  date='2020/1/14 23:01'
*/
package main

import (
	"fmt"
	"strings"
	"time"
)

type Reader interface {
	Read(rc chan string)
}

type Writer interface {
	Write(wc chan string)
}

//这里里面包括读和写的结构体
type Process struct {
	rc    chan string
	wc    chan string
	read  Reader
	write Writer
}

type ReadFromFile struct {
	path string //读取文件的路径
}

type Write2InfluxDB struct {
	influxDBDsn string //influx data source
}

func (r *ReadFromFile) Read(rc chan string) {
	//读取模块
	line := "dumessage"
	rc <- line
}

func (w *Write2InfluxDB) Write(wc chan string) {
	//写入模块
	fmt.Println(<-wc)
}

func (l *Process) Analysis() {
	//解析模板
	data := <-l.rc
	l.wc <- strings.ToUpper(data)
}

func main() {
	r := &ReadFromFile{path: "tmp/access.log"}
	w := &Write2InfluxDB{influxDBDsn: "username&password..."}
	lp := &Process{
		rc:    make(chan string),
		wc:    make(chan string),
		read:  r,
		write: w,
	}
	go lp.read.Read(lp.rc)
	go lp.Analysis()
	go lp.write.Write(lp.wc)
	time.Sleep(1 * time.Second)
}
