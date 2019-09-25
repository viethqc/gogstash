package inputdockerlog

import (
	"io/ioutil"
	"time"
	"unsafe"

	"github.com/viethqc/gogstash/KDGoLib/futil"
	"github.com/viethqc/gogstash/KDGoLib/mmfile"
)

func NewSinceFile(filepath string) (sincefile *SinceFile, err error) {
	sincefile = &SinceFile{}
	err = sincefile.Open(filepath)
	return
}

type SinceFile struct {
	mmfile mmfile.MMFile
	Since  *time.Time
}

func (t *SinceFile) Open(filepath string) (err error) {
	if err = t.Close(); err != nil {
		return
	}

	if !futil.IsExist(filepath) {
		if err = ioutil.WriteFile(filepath, make([]byte, 32), 0644); err != nil {
			return
		}
	}

	if t.mmfile, err = mmfile.Open(filepath); err != nil {
		return
	}

	t.Since = (*time.Time)(unsafe.Pointer(&t.mmfile.Data()[0]))

	return
}

func (t *SinceFile) Close() (err error) {
	if t.mmfile != nil {
		if err = t.mmfile.Close(); err != nil {
			return
		}
	}
	t.Since = nil

	return
}
