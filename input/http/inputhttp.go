package inputhttp

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/viethqc/gogstash/config"
	"github.com/viethqc/gogstash/config/logevent"
)

// ModuleName is the name used in config file
const ModuleName = "http"

// ErrorTag tag added to event when process module failed
const ErrorTag = "gogstash_input_http_error"

// InputConfig holds the configuration json fields and internal objects
type InputConfig struct {
	config.InputConfig
	Method   string `json:"method,omitempty"` // one of ["HEAD", "GET"]
	URL      string `json:"url"`
	Interval int    `json:"interval,omitempty"`

	hostname string
}

// DefaultInputConfig returns an InputConfig struct with default values
func DefaultInputConfig() InputConfig {
	return InputConfig{
		InputConfig: config.InputConfig{
			CommonConfig: config.CommonConfig{
				Type: ModuleName,
			},
		},
		Method:   "GET",
		Interval: 60,
	}
}

// InitHandler initialize the input plugin
func InitHandler(ctx context.Context, raw *config.ConfigRaw) (config.TypeInputConfig, error) {
	conf := DefaultInputConfig()
	err := config.ReflectConfig(raw, &conf)
	if err != nil {
		return nil, err
	}

	if conf.hostname, err = os.Hostname(); err != nil {
		return nil, err
	}

	conf.Codec, err = config.GetCodec(ctx, *raw)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

// Start wraps the actual function starting the plugin
func (t *InputConfig) Start(ctx context.Context, msgChan chan<- logevent.LogEvent) (err error) {
	startChan := make(chan bool, 1) // startup tick
	ticker := time.NewTicker(time.Duration(t.Interval) * time.Second)
	defer ticker.Stop()

	startChan <- true

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-startChan:
			t.Request(ctx, msgChan)
		case <-ticker.C:
			t.Request(ctx, msgChan)
		}
	}
}

func (t *InputConfig) Request(ctx context.Context, msgChan chan<- logevent.LogEvent) {
	data, err := t.SendRequest()
	extra := map[string]interface{}{
		"host": t.hostname,
		"url":  t.URL,
	}
	if err != nil {
		event := logevent.LogEvent{
			Timestamp: time.Now(),
			Extra:     extra,
		}
		event.AddTag(ErrorTag)
		msgChan <- event
		return
	}

	t.Codec.Decode(ctx, data, extra, msgChan)

	return
}

func (t *InputConfig) SendRequest() (data []byte, err error) {
	var (
		res *http.Response
		raw []byte
	)
	switch t.Method {
	case "HEAD":
		res, err = http.Head(t.URL)
	case "GET":
		res, err = http.Get(t.URL)
	default:
		err = errors.New("Unknown method")
	}

	if err != nil {
		return
	}

	defer res.Body.Close()
	if raw, err = ioutil.ReadAll(res.Body); err != nil {
		return
	}
	data = bytes.TrimSpace(raw)

	return
}
