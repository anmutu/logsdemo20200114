/*
  author='du'
  date='2020/1/15 0:18'
*/
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Reader interface {
	Read(rc chan []byte)
}

type Writer interface {
	Write(wc chan string)
}

//这里里面包括读和写的结构体
type Process struct {
	rc    chan []byte
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

func (r *ReadFromFile) Read(rc chan []byte) {
	//读取模块
	//第一步，打开文件
	f, err := os.Open(r.path)
	if err != nil {
		panic(fmt.Sprintf("open file error %s", err.Error()))
	}

	//第二步，从文件末尾逐行开始读取文件内容
	f.Seek(0, 2)
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadBytes('\n')
		if err == io.EOF {
			time.Sleep(500 * time.Millisecond)
			continue
		} else if err != nil {
			panic(fmt.Sprintf("ReadBytes error:%s", err.Error()))
		}
		rc <- line[:len(line)-1]
	}
}

func (l *Process) Analysis() {
	//解析模板
	for data := range l.rc {
		l.wc <- strings.ToUpper(string(data))
	}
}

func (w *Write2InfluxDB) Write(wc chan string) {
	//写入模块
	for data := range wc {
		fmt.Println(data)
	}
}

func main() {
	r := &ReadFromFile{path: "./access.log"}
	w := &Write2InfluxDB{influxDBDsn: "username&password..."}
	lp := &Process{
		rc:    make(chan []byte),
		wc:    make(chan string),
		read:  r,
		write: w,
	}
	lp.read.Read(lp.rc)
	//lp.Analysis()
	//lp.write.Write(lp.wc)
	time.Sleep(1 * time.Second)
}
