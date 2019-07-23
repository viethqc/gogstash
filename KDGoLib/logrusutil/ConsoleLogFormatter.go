package logrusutil

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/viethqc/KDGoLib/errutil"
	"github.com/viethqc/KDGoLib/runtimecaller"
)

// flags
const (
	Llongfile = 1 << iota
	Lshortfile
	Ltime
	Llevel
	Lhostname
	Lloggername
	LstdFlags = Ltime | Lshortfile | Llevel | Lhostname | Lloggername
)

// ConsoleLogFormatter suitable formatter for console
type ConsoleLogFormatter struct {
	TimestampFormat      string
	Flag                 int
	CallerOffset         int
	LoggerName           string
	HostName             string
	RuntimeCallerFilters []runtimecaller.Filter
}

func addspace(text string, addspaceflag bool) (string, bool) {
	if addspaceflag {
		return " " + text, true
	}
	return text, true
}

func filterLogrusRuntimeCaller(callinfo runtimecaller.CallInfo) (valid bool, stop bool) {
	return !strings.Contains(callinfo.PackageName(), "github.com/sirupsen/logrus"), false
}

// Format output logrus entry
func (t *ConsoleLogFormatter) Format(entry *logrus.Entry) (data []byte, err error) {
	buffer := bytes.Buffer{}
	addspaceflag := false

	if t.Flag == 0 {
		t.Flag = LstdFlags
	}

	if t.TimestampFormat == "" {
		t.TimestampFormat = time.RFC3339
	}

	if t.Flag&Ltime != 0 {
		timetext := entry.Time.Format(t.TimestampFormat)
		timetext, addspaceflag = addspace(timetext, addspaceflag)
		if _, err = buffer.WriteString(timetext); err != nil {
			err = errutil.New("write timestamp to buffer failed", err)
			return
		}
	}

	if t.Flag&(Lshortfile|Llongfile) != 0 {
		var filelinetext string
		filters := append([]runtimecaller.Filter{filterLogrusRuntimeCaller}, t.RuntimeCallerFilters...)
		if callinfo, ok := errutil.RuntimeCaller(1+t.CallerOffset, filters...); ok {
			if t.Flag&Lshortfile != 0 {
				filelinetext = fmt.Sprintf("[%s:%d]", callinfo.FileName(), callinfo.Line())
			} else {
				filelinetext = fmt.Sprintf("[%s/%s:%d]", callinfo.PackageName(), callinfo.FileName(), callinfo.Line())
			}

			filelinetext, addspaceflag = addspace(filelinetext, addspaceflag)
		}

		if _, err = buffer.WriteString(filelinetext); err != nil {
			err = errutil.New("write fileline to buffer failed", err)
			return
		}
	}

	if t.Flag&Lloggername != 0 {
		loggerNameText := fmt.Sprintf("[%s]", t.LoggerName)
		loggerNameText, addspaceflag = addspace(loggerNameText, addspaceflag)
		if _, err = buffer.WriteString(loggerNameText); err != nil {
			err = errutil.New("write level to buffer failed", err)
			return
		}
	}

	if t.Flag&Lhostname != 0 {
		hostnameText := fmt.Sprintf("[%s]", t.HostName)
		hostnameText, addspaceflag = addspace(hostnameText, addspaceflag)
		if _, err = buffer.WriteString(hostnameText); err != nil {
			err = errutil.New("write level to buffer failed", err)
			return
		}
	}

	if t.Flag&Llevel != 0 {
		leveltext := fmt.Sprintf("[%s]", entry.Level.String())
		leveltext, addspaceflag = addspace(leveltext, addspaceflag)
		if _, err = buffer.WriteString(leveltext); err != nil {
			err = errutil.New("write level to buffer failed", err)
			return
		}
	}

	f := ""
	l := 0
	fn := ""
	fnName := ""
	if pc, file, line, ok := runtime.Caller(2); ok {
		f = file
		l = line
		fun := runtime.FuncForPC(pc)
		fn = fun.Name()
		fnName = filepath.Ext(fn)

		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				f = file[i+1:]
				break
			}
		}

		if len(fnName) > 0 && fnName[0] == '.' {
			fnName = fnName[1:]
		}
	}

	fileText := fmt.Sprintf("[%s:%d]", f, l)
	fileText, addspaceflag = addspace(fileText, addspaceflag)
	if _, err = buffer.WriteString(fileText); err != nil {
		err = errutil.New("write level to buffer failed", err)
		return
	}

	message := entry.Message
	message, _ = addspace(message, addspaceflag)
	if _, err = buffer.WriteString(message); err != nil {
		err = errutil.New("write message to buffer failed", err)
		return
	}

	if err = buffer.WriteByte('\n'); err != nil {
		err = errutil.New("write newline to buffer failed", err)
		return
	}

	data = buffer.Bytes()
	return
}
