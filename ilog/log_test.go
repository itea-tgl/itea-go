package ilog

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)


//func Test_1(t *testing.T) {
//	Init(LogFile,
//		WithFile("aaaa.log"),
//		EnableDivide(),
//		EnableRotate(),
//		FileKeep(1),
//	)
//
//	i := 0
//
//	for {
//		if i == 180 {
//			break
//		}
//		Info("a", "b", "c")
//		Error("1", "2", "3")
//		//Fatal("1", "1", "1")
//		//Debug(123)
//		time.Sleep(1 * time.Second)
//		i++
//	}
//
//	for {
//		if Done() {
//			break
//		}
//	}
//}

func Test_ConsoleLog(t *testing.T) {
	Init(nil)
	Info("a", "b", "c")
}

func Test_FileLog(t *testing.T) {
	Init(LogFile)
	Info("a", "b", "c")
	for {
		if Done() {
			break
		}
	}
	checkFile(t, "itea.log")
}

func Test_FileLog_Rotate(t *testing.T) {
	Init(LogFile, EnableRotate())
	Info("a", "b", "c")
	for {
		if Done() {
			break
		}
	}
	checkFile(t, fmt.Sprintf("itea-%s.log", time.Now().Format("2006-01-02")))
}

func Test_FileLog_Rotate_File(t *testing.T) {
	Init(LogFile, EnableRotate(), WithFile("aaa.log"))
	Info("a", "b", "c")
	for {
		if Done() {
			break
		}
	}
	checkFile(t, fmt.Sprintf("aaa-%s.log", time.Now().Format("2006-01-02")))
}

func checkFile(t *testing.T, file string) {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsExist(err) {
			t.Errorf("log file %s not found", file)
		}
		t.Error("log file error : ", err)
		return
	}

	defer func() {
		err = os.Remove(file)
		if err != nil {
			fmt.Println("file remove error : ", err)
		}
	}()

	dat, err := ioutil.ReadFile(file)
	if err != nil {
		t.Error("read log file error : ", err)
		return
	}

	now := time.Now().Format("2006/01/02 15:04:05")
	expect := fmt.Sprintf("[INFO] %s %s", now, "a b c")

	if strings.Trim(string(dat), "\r\n") != expect {
		t.Errorf("log content error, expect %s, get %s", expect, string(dat))
	}
}