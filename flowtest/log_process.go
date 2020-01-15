/*
  author='du'
  date='2020/1/14 21:51'
*/
package flowtest

import (
	"fmt"
	"strings"
	"time"
)

type LogProcess struct {
	path        string      //读取文件的路径
	influxDBDsn string      //influx data source
	rc          chan string //从读取模块到解析模板用来传递数据
	wc          chan string //从解析模板到写入模板用来传递数据
}

func (l *LogProcess) ReadFromFile() {
	//读取模块
	line := "dumessage"
	l.rc <- line
}

func (l *LogProcess) Analysis() {
	//解析模板
	data := <-l.rc
	l.wc <- strings.ToUpper(data)
}

func (l *LogProcess) Write2InfluxDB() {
	//写入模块
	fmt.Println(<-l.wc)
}

func main() {
	//实例化结构体
	lp := &LogProcess{
		path:        "tmp/access.log",
		influxDBDsn: "username&password...",
		rc:          make(chan string),
		wc:          make(chan string),
	}
	go lp.ReadFromFile()
	go lp.Analysis()
	go lp.Write2InfluxDB()
	time.Sleep(1 * time.Second)
}
