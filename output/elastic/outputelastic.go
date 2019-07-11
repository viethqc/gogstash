package outputelastic

import (
	"context"
	"crypto/sha1"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"github.com/tsaikd/KDGoLib/errutil"
	"github.com/tsaikd/gogstash/config"
	"github.com/tsaikd/gogstash/config/goglog"
	"github.com/tsaikd/gogstash/config/logevent"
	"gopkg.in/olivere/elastic.v2"
)

// ModuleName is the name used in config file
const ModuleName = "elastic"

// OutputConfig holds the configuration json fields and internal objects
type OutputConfig struct {
	config.OutputConfig
	URL                  []string `json:"url"` // elastic API entrypoints
	resolvedURLs         []string // URLs after resolving environment vars
	Index                string   `json:"index"`             // index name to log
	DocumentType         string   `json:"document_type"`     // type name to log
	DocumentID           string   `json:"document_id"`       // id to log, used if you want to control id format
	RetryOnConflict      int      `json:"retry_on_conflict"` // the number of times Elasticsearch should internally retry an update/upserted document
	Action               string   `json:"action"`
	RetryInitialInterval int      `json:"retry_initial_interval"`
	RetryMaxInterval     int      `json:"retry_max_interval"`
	RetriableCode        []int    `json:"retriable_code"`

	Sniff bool `json:"sniff"` // find all nodes of your cluster, https://github.com/olivere/elastic/wiki/Sniffing

	BulkActions int `json:"bulk_actions,omitempty"`

	// BulkSize specifies when to flush based on the size (in bytes) of the actions
	// currently added. Defaults to 5 MB and can be set to -1 to be disabled.
	BulkSize int `json:"bulk_size,omitempty"`

	// BulkFlushInterval specifies when to flush at the end of the given interval.
	// Defaults to 30 seconds. If you want the bulk processor to
	// operate completely asynchronously, set both BulkActions and BulkSize to
	// -1 and set the FlushInterval to a meaningful interval.
	BulkFlushInterval time.Duration `json:"bulk_flush_interval"`

	// ExponentialBackoffInitialTimeout used to set the first/minimal interval in elastic.ExponentialBackoff
	// Defaults to 10s
	ExponentialBackoffInitialTimeout string `json:"exponential_backoff_initial_timeout,omitempty"`
	exponentialBackoffInitialTimeout time.Duration

	// ExponentialBackoffMaxTimeout used to set the maximum wait interval in elastic.ExponentialBackoff
	// Defaults to 5m
	ExponentialBackoffMaxTimeout string `json:"exponential_backoff_max_timeout,omitempty"`
	exponentialBackoffMaxTimeout time.Duration

	// SSLCertValidation Option to validate the server's certificate. Disabling this severely compromises security.
	// For more information on disabling certificate verification please read https://www.cs.utexas.edu/~shmat/shmat_ccs12.pdf
	SSLCertValidation bool `json:"ssl_certificate_validation,omitempty"`

	client        *elastic.Client        // elastic client instance
	processor     *elastic.BulkProcessor // elastic bulk processor
	retryCount    map[string]int
	retryableCode map[int]bool
}

// DefaultOutputConfig returns an OutputConfig struct with default values
func DefaultOutputConfig() OutputConfig {
	return OutputConfig{
		OutputConfig: config.OutputConfig{
			CommonConfig: config.CommonConfig{
				Type: ModuleName,
			},
		},
		RetryOnConflict:                  1,
		BulkActions:                      1000,    // 1000 actions
		BulkSize:                         5 << 20, // 5 MB
		BulkFlushInterval:                30 * time.Second,
		ExponentialBackoffInitialTimeout: "10s",
		ExponentialBackoffMaxTimeout:     "5m",
		SSLCertValidation:                true,
	}
}

// errors
var (
	ErrorCreateClientFailed1 = errutil.NewFactory("create elastic client failed: %q")
)

type errorLogger struct {
	logger logrus.FieldLogger
}

// Printf log format string to error level
func (l *errorLogger) Printf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

type jsonDecoder struct{}

// Decode decodes with jsoniter.Unmarshal
func (u *jsonDecoder) Decode(data []byte, v interface{}) error {
	return jsoniter.Unmarshal(data, v)
}

// InitHandler initialize the output plugin
func InitHandler(ctx context.Context, raw *config.ConfigRaw) (config.TypeOutputConfig, error) {
	conf := DefaultOutputConfig()
	err := config.ReflectConfig(raw, &conf)
	if err != nil {
		return nil, err
	}

	// map Printf to error level
	logger := &errorLogger{logger: goglog.Logger}

	// replace env var names with values on URL config
	for _, url := range conf.URL {
		newURL := logevent.FormatWithEnv(url)
		conf.resolvedURLs = append(conf.resolvedURLs, newURL)
	}

	options := []elastic.ClientOptionFunc{
		elastic.SetURL(conf.resolvedURLs...),
		elastic.SetSniff(conf.Sniff),
		elastic.SetErrorLog(logger),
		elastic.SetDecoder(&jsonDecoder{}),
	}

	// set httpclient explicitly if we need to avoid https cert checks
	if !conf.SSLCertValidation {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		options = append(options, elastic.SetHttpClient(client))
	}

	if conf.client, err = elastic.NewClient(options...); err != nil {
		return nil, ErrorCreateClientFailed1.New(err, conf.URL)
	}

	conf.exponentialBackoffInitialTimeout, err = time.ParseDuration(conf.ExponentialBackoffInitialTimeout)
	if err != nil {
		return nil, err
	}

	conf.exponentialBackoffMaxTimeout, err = time.ParseDuration(conf.ExponentialBackoffMaxTimeout)
	if err != nil {
		return nil, err
	}

	conf.retryCount = make(map[string]int)
	conf.retryableCode = make(map[int]bool)
	for _, retryCode := range conf.RetriableCode {
		conf.retryableCode[retryCode] = true
	}

	conf.processor, err = conf.client.BulkProcessor().
		Name("gogstash-output-elastic").
		BulkActions(conf.BulkActions).
		BulkSize(conf.BulkSize).
		FlushInterval(conf.BulkFlushInterval * time.Second).
		After(conf.BulkAfter).
		Do()
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

// BulkAfter execute after a commit to Elasticsearch
func (t *OutputConfig) BulkAfter(executionID int64, requests []elastic.BulkableRequest, response *elastic.BulkResponse, err error) {
	if err == nil && response.Errors {
		// find failed requests, and log it
		for i, item := range response.Items {
			for _, v := range item {
				if v.Error != "" {
					data, err := requests[i].Source()
					if err != nil {
						goglog.Logger.Error(err)
						continue
					}

					goglog.Logger.Errorf("%s: bulk processor request %s failed: %s", ModuleName, requests[i].String(), v.Error)
					if _, ok := t.retryableCode[v.Status]; ok {
						t.retry(data[1], v)
					}
				}
			}
		}
	}
}

func (t *OutputConfig) retry(data string, errorInfo *elastic.BulkResponseItem) error {
	status := errorInfo.Status

	h := sha1.New()
	h.Write([]byte(data))
	bs := h.Sum(nil)
	sha1_data := fmt.Sprintf("%x", bs)

	goglog.Logger.Errorf("data: %s", data)
	goglog.Logger.Errorf("status code: %d", status)
	goglog.Logger.Errorf("sha1: %s", sha1_data)

	if _, ok := t.retryCount[sha1_data]; ok {
		t.retryCount[sha1_data] += 1
	} else {
		t.retryCount[sha1_data] = 1
	}

	goglog.Logger.Infof("Retry number: %d", t.retryCount[sha1_data])

	if t.retryCount[sha1_data] > t.RetryInitialInterval {
		goglog.Logger.Infof("Retry over quata")
		delete(t.retryCount, sha1_data)
		return nil
	}

	event := logevent.LogEvent{
		Timestamp: time.Now(),
		Extra:     nil,
	}

	kk := make(map[string]interface{})
	if err := jsoniter.Unmarshal([]byte(data), &kk); err != nil {
		goglog.Logger.Error(err.Error())
	}

	event.Extra = kk["doc"].(map[string]interface{})

	ticker := time.NewTicker(time.Duration(t.RetryMaxInterval) * time.Second)
	go func(ticker *time.Ticker) {
		<-ticker.C
		t.Output(nil, event)
	}(ticker)

	return nil
}

// Output event
func (t *OutputConfig) Output(ctx context.Context, event logevent.LogEvent) (err error) {
	index := event.Format(t.Index)
	// elastic index name should be lowercase
	index = strings.ToLower(index)
	doctype := event.Format(t.DocumentType)
	id := event.Format(t.DocumentID)
	action := event.Format(t.Action)

	switch action {
	case "index":
		indexRequest := elastic.NewBulkIndexRequest().
			Index(index).
			Type(doctype).
			Id(id).
			Doc(event.Extra)
		t.processor.Add(indexRequest)
		break
	case "update":
		updateRequest := elastic.NewBulkUpdateRequest().
			Index(index).
			Type(doctype).
			Id(id).
			Doc(event.Extra)

		t.processor.Add(updateRequest)
		break
	}

	return
}
