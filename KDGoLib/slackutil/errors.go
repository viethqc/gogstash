package slackutil

import "github.com/viethqc/gogstash/KDGoLib/errutil"

// errors
var (
	ErrSlackAuthTestFailed = errutil.NewFactory("slack client auth test failed")
	ErrSlackFormatFailed1  = errutil.NewFactory("format slack message failed %q")
)
