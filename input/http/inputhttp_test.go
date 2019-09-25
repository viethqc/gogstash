package inputhttp

import (
	"context"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/viethqc/gogstash/config"
	"github.com/viethqc/gogstash/config/goglog"
)

func init() {
	goglog.Logger.SetLevel(logrus.DebugLevel)
	config.RegistInputHandler(ModuleName, InitHandler)
	config.RegistCodecHandler(config.DefaultCodecName, config.DefaultCodecInitHandler)
}

func TestMain(m *testing.M) {
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("foo"))
	})

	go func() {
		if err := http.ListenAndServe("127.0.0.1:8090", nil); err != nil {
			goglog.Logger.Fatal(err)
		}
	}()

	ret := m.Run()

	os.Exit(ret)
}

func Test_input_http_module(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	ctx := context.Background()
	conf, err := config.LoadFromYAML([]byte(strings.TrimSpace(`
debugch: true
input:
  - type: http
    method: GET
    url: "http://127.0.0.1:8090/"
    interval: 3
    codec:
      type: "default"
	`)))
	require.NoError(err)
	require.NoError(conf.Start(ctx))

	time.Sleep(500 * time.Millisecond)
	if event, err := conf.TestGetOutputEvent(100 * time.Millisecond); assert.NoError(err) {
		assert.Equal("foo", event.Message)
	}
}
