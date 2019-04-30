package db_mysql

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func loginfo(f interface{}, v ...interface{}) {
	if !logTofile {
		return
	}
	writeFile("log", formatLog(f, v...))
}

func errorinfo(f interface{}, v ...interface{}) {
	if !errTofile {
		return
	}
	writeFile("error", formatLog(f, v...))
}

func writeFile(logtype, msg string) {
	defer checkerr()
	var fi *os.File
	var err error
	createDir("dblogs")
	var _fileName string = "dblogs/" + logtype + time.Now().Format("20060102") + ".log"
	_exist, _size := checkFile(_fileName)
	if _exist && _size > 2048*1024 {
		errr := os.Rename(_fileName, "dblogs/"+logtype+time.Now().Format("20060102150405")+".log")
		if errr != nil {
			//fmt.Println(errr.Error())
		}
		_exist = false
	}
	if _exist {
		fi, err = os.OpenFile(_fileName, os.O_APPEND, 0666)
	} else {
		fi, err = os.Create(_fileName)
	}
	if err == nil {
		io.WriteString(fi, fmt.Sprintf("#-------------%s--------------#\r", time.Now()))
		_, err = io.WriteString(fi, msg+"\r\r")
		//		w := bufio.NewWriter(fi)
		//		_,err = w.WriteString(msg+"\r")
		//		w.Flush()
		fi.Close()
	}
}

func createDir(dirName string) {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		os.MkdirAll(dirName, 0777)
	}
}

func checkFile(fileName string) (exist bool, size int64) {
	exist = true
	size = 0
	var fileInfo os.FileInfo
	var err error
	if fileInfo, err = os.Stat(fileName); os.IsNotExist(err) {
		exist = false
	}
	if exist {
		size = fileInfo.Size()
	}
	return
}

func formatLog(f interface{}, v ...interface{}) string {
	var msg string
	switch f.(type) {
	case string:
		msg = f.(string)
		if len(v) == 0 {
			return msg
		}
		if strings.Contains(msg, "%") && !strings.Contains(msg, "%%") {
			//format string
		} else {
			//do not contain format char
			msg += strings.Repeat(" %v", len(v))
		}
	default:
		msg = fmt.Sprint(f)
		if len(v) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(v))
	}
	return fmt.Sprintf(msg, v...)
}

func argsToString(args ...interface{}) string {
	var _sb bytes.Buffer
	for i := range args {
		if i != 0 && i != len(args) {
			_sb.WriteString(",")
		}
		_sb.WriteString(fmt.Sprint(args[i]))
	}
	return _sb.String()
}
