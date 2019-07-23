package logrusutil

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/viethqc/gogstash/KDGoLib/errutil"
	"github.com/viethqc/gogstash/KDGoLib/runtimecaller"
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

func (t *ConsoleLogFormatter) getFunctionName(functionName string) string {
	arrData := strings.Split(functionName, ".")
	if len(arrData) != 0 {
		return arrData[len(arrData)-1]
	}

	return ""
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

	if t.Flag&(Lshortfile|Llongfile) != 0 {
		var filelinetext string
		filters := append([]runtimecaller.Filter{filterLogrusRuntimeCaller}, t.RuntimeCallerFilters...)
		if callinfo, ok := errutil.RuntimeCaller(1+t.CallerOffset, filters...); ok {
			if t.Flag&Lshortfile != 0 {
				filelinetext = fmt.Sprintf("[%s:%d][%s]", callinfo.FileName(), callinfo.Line(), t.getFunctionName(callinfo.PCFunc().Name()))
			} else {
				filelinetext = fmt.Sprintf("[%s/%s:%d][%s]", callinfo.PackageName(), callinfo.FileName(), callinfo.Line())
			}

			filelinetext, addspaceflag = addspace(filelinetext, addspaceflag)
		}

		if _, err = buffer.WriteString(filelinetext); err != nil {
			err = errutil.New("write fileline to buffer failed", err)
			return
		}
	}

	if t.Flag&Llevel != 0 {
		leveltext := fmt.Sprintf("[%s]: ", strings.ToUpper(entry.Level.String()))
		leveltext, addspaceflag = addspace(leveltext, addspaceflag)
		if _, err = buffer.WriteString(leveltext); err != nil {
			err = errutil.New("write level to buffer failed", err)
			return
		}
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
