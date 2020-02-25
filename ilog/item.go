package ilog

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

var rotateItems []*item

type item struct {
	filename string
	prefix string
	log *log.Logger
	file *os.File
}

func NewItem(p string, o option) *item {
	i := &item{filename: o.logfile, prefix: p}
	if o.divide {
		o.logfile = i.formatFilename()
	}
	i.logPrefix()
	if o.rotate {
		rotateFile(i)
		o.logfile = i.rotateName()
	}
	file, err := os.OpenFile(o.logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND,0)
	if err != nil {
		panic("open log file error !")
	}
	i.file = file
	i.log = log.New(file, i.prefix, log.LstdFlags)
	return i
}

func (i *item) rotateName() string {
	l := len(i.filename)
	if i.filename[l-4:] == ".log" {
		i.filename = i.filename[:l-4]
	}
	var s bytes.Buffer
	s.WriteString(i.filename)
	s.WriteString("-")
	//s.WriteString(time.Now().Format("15-04")) //code for test
	s.WriteString(time.Now().Format("2006-01-02"))
	s.WriteString(".log")
	return s.String()
}

func rotateFile(i *item) {
	rotateItems = append(rotateItems, i)
	if len(rotateItems) > 1 {
		return
	}
	go func() {
		for {
			now := time.Now()
			// 计算下一个零点
			//next := now.Add(time.Second * 60) //code for test
			//next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), 0, 0, next.Location()) //code for test
			next := now.Add(time.Second * 60)
			next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), 0, 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			<-t.C
			for _, i := range rotateItems {
				go fileRotate(i)
			}
			go fileClean(clean)
		}
	}()
}

func fileRotate(i *item) {
	name := i.rotateName()
	for {
		file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND,0)
		if err == nil {
			i.log = log.New(file, i.prefix, log.LstdFlags)
			err = i.file.Close()
			if err != nil {
				log.Println("former log file close error : ", err)
			}
			i.file = file
			break
		}
	}
}

func fileClean(n int) {
	rd, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println("log dir scan error : ", err)
		return
	}
	for _, fi := range rd {
		if fi.IsDir() {
			continue
		} else {
			//log.Println(fi.Name(), "-", fi.ModTime().Format("2006-01-02 15:04:05"), "-", time.Since(fi.ModTime()).Seconds()) //code for test
			//if time.Since(fi.ModTime()).Seconds() >= (float64(n) * time.Minute.Seconds() - 1) && isLog(fi.Name()) { //code for test
			if time.Since(fi.ModTime()).Seconds() >= (float64(n) * 24 * time.Hour.Seconds() - 1)  && isLog(fi.Name()) {
				err = os.Remove(dir + fi.Name())
				if err != nil {
					log.Println("file remove error : ", err)
				}
			}
		}
	}
}

func isLog(f string) bool {
	l := len(f)
	return l > 4 && f[l-4:] == ".log"
}

func (i *item) formatFilename() string {
	arr := strings.Split(i.filename, "/")
	l := len(arr)
	if l <= 0 {
		return ""
	}
	arr[l-1] = strings.ToLower(i.prefix) + "-" + arr[l-1]
	i.filename = strings.Join(arr, "/")
	return i.filename
}

func (i *item) logPrefix() {
	i.prefix = fmt.Sprintf("[%s] ", i.prefix)
}

func (i *item) print(v ...interface{}) {
	i.log.Println(v...)
}

func (i *item) close() {
	err := i.file.Close()
	if err != nil {
		log.Println(err)
	}
}