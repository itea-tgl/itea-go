package system

import (
	"fmt"
	"github.com/itea-tgl/itea-go/ilog"
	"strings"
)

const LOG_KEY = "log"

type Log struct {
	Type 		string
	Logfile 	string
	Rotate 		bool
	Divide		bool
}

func InitLog() {
	logtype, logfile, rotate, divide := "", "", false, false
	if c := Conf.GetStruct(fmt.Sprintf("%s.%s", Conf.FileName, LOG_KEY), Log{}); c != nil {
		logConf := c.(*Log)
		if !strings.EqualFold(logConf.Type, "") {
			logtype = logConf.Type
		}
		if !strings.EqualFold(logConf.Logfile, "") {
			logfile = logConf.Logfile
		}
		rotate = logConf.Rotate
		divide = logConf.Divide
	}
	if strings.EqualFold(logtype, "file") {
		var opts []ilog.IOption
		if !strings.EqualFold(logfile, "") {
			opts = append(opts, ilog.WithFile(logfile))
		}
		if rotate {
			opts = append(opts, ilog.EnableRotate())
		}
		if divide {
			opts = append(opts, ilog.EnableDivide())
		}
		ilog.Init(ilog.LogFile, opts...)
	} else {
		ilog.Init(nil)
	}
}
