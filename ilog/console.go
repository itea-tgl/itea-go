package ilog

import (
	"log"
	"time"
)

type Console struct {

}

func (*Console) Init() {

}

func (*Console) Done() bool{
	return true
}

func (*Console) Debug(v ...interface{}){
	log.Println("[DEBUG]", time.Now().Format("2006/01/02 15:04:05"),  v)
}

func (*Console) Info(v ...interface{}){
	log.Println("[INFO]", time.Now().Format("2006/01/02 15:04:05"), v)
}

func (*Console) Error(v ...interface{}){
	log.Println("[ERROR]", time.Now().Format("2006/01/02 15:04:05"), v)
}

func (*Console) Fatal(v ...interface{}){
	log.Println("[FATAL]", time.Now().Format("2006/01/02 15:04:05"), v)
}

func LogConsole() ILog {
	return &Console{}
}