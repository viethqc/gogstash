package errutil

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type invalidFormatter struct {
}

func (t *invalidFormatter) Format(errin error) (errtext string, err error) {
	return "", New("invalid formatter")
}

func (t *invalidFormatter) FormatSkip(errin error, skip int) (errtext string, err error) {
	return "", New("invalid formatter")
}

func Test_Trace(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	logger := defaultLogger
	defer func() {
		defaultLogger = logger
	}()

	buffer := &bytes.Buffer{}
	SetDefaultLogger(loggerImpl{LoggerOptions{
		ErrorOutput:    buffer,
		TraceFormatter: NewJSONFormatter(),
	}})

	errchild1 := New("test error 1")
	errchild2 := New("test error 2")
	errtest := New("test error", errchild1, errchild2)

	Trace(errtest)

	errtext := buffer.String()
	require.Contains(errtext, `"test error 1"`)
	require.Contains(errtext, `"test error 2"`)
	require.Contains(errtext, `trace_test.go:42`)
	require.True(strings.HasSuffix(errtext, "\n"))
}

func Test_TraceWrap(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	logger := defaultLogger
	defer func() {
		defaultLogger = logger
	}()

	buffer := &bytes.Buffer{}
	SetDefaultLogger(loggerImpl{LoggerOptions{
		ErrorOutput:    buffer,
		TraceFormatter: NewJSONFormatter(),
	}})

	errchild1 := New("test error 1")
	errchild2 := New("test error 2")
	errtest := New("test error", errchild1, errchild2)

	TraceWrap(errtest, New("test error wrapper"))

	errtext := buffer.String()
	require.Contains(errtext, `"test error wrapper"`)
	require.Contains(errtext, `"test error"`)
	require.Contains(errtext, `"test error 1"`)
	require.Contains(errtext, `"test error 2"`)
	require.True(strings.HasSuffix(errtext, "\n"))
}

func Test_TraceWrapNil(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	logger := defaultLogger
	defer func() {
		defaultLogger = logger
	}()

	buffer := &bytes.Buffer{}
	SetDefaultLogger(loggerImpl{LoggerOptions{
		ErrorOutput:    buffer,
		TraceFormatter: NewJSONFormatter(),
	}})

	TraceWrap(nil, New("test error wrapper"))

	errtext := buffer.String()
	require.Zero(errtext)
}

func Test_TracePanic(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	logger := defaultLogger
	defer func() {
		defaultLogger = logger
	}()

	buffer := &bytes.Buffer{}
	SetDefaultLogger(loggerImpl{LoggerOptions{
		ErrorOutput:    buffer,
		TraceFormatter: &invalidFormatter{},
	}})

	errtest := NewErrors(
		New("test error 1"),
		New("test error 2"),
	)

	require.Panics(func() {
		Trace(errtest)
	})
}
