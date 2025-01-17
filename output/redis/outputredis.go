package outputredis

import (
	"context"
	"time"

	"github.com/viethqc/gogstash/KDGoLib/errutil"
	"github.com/viethqc/gogstash/KDGoLib/timeutil"
	"github.com/viethqc/gogstash/config"
	"github.com/viethqc/gogstash/config/goglog"
	"github.com/viethqc/gogstash/config/logevent"
	"gopkg.in/redis.v5"
)

// ModuleName is the name used in config file
const ModuleName = "redis"

// ErrorTag tag added to event when process module failed
const ErrorTag = "gogstash_output_redis_error"

// OutputConfig holds the configuration json fields and internal objects
type OutputConfig struct {
	config.OutputConfig
	Host              []string `json:"host"`
	Key               string   `json:"key"`
	DataType          string   `json:"data_type,omitempty"` // one of ["list", "channel"]
	Timeout           int      `json:"timeout,omitempty"`
	ReconnectInterval int      `json:"reconnect_interval,omitempty"`
	Password          string   `json:"password,omitempty"`
	Db                int      `json:"db,omitempty"`
	Ttl               int      `json:"ttl,omitempty"`
	Connections       int      `json:"connections"` // maximum number of socket connections, default: 10

	client *redis.Client
}

// DefaultOutputConfig returns an OutputConfig struct with default values
func DefaultOutputConfig() OutputConfig {
	return OutputConfig{
		OutputConfig: config.OutputConfig{
			CommonConfig: config.CommonConfig{
				Type: ModuleName,
			},
		},
		Host:              []string{"localhost:6379"},
		Key:               "gogstash",
		DataType:          "list",
		Timeout:           5,
		ReconnectInterval: 1,
		Connections:       10,
	}
}

// errors
var (
	ErrorPingFailed           = errutil.NewFactory("ping redis server failed")
	ErrorEventMarshalFailed1  = errutil.NewFactory("event Marshal failed: %v")
	ErrorUnsupportedDataType1 = errutil.NewFactory("unsupported data type: %q")
)

// InitHandler initialize the output plugin
func InitHandler(ctx context.Context, raw *config.ConfigRaw) (config.TypeOutputConfig, error) {
	conf := DefaultOutputConfig()
	err := config.ReflectConfig(raw, &conf)
	if err != nil {
		return nil, err
	}

	if len(conf.Host) > 1 {
		goglog.Logger.Warn("deprecated: host number should be only 1")
	}

	conf.client = redis.NewClient(&redis.Options{
		Addr:     conf.Host[0],
		Password: conf.Password, // no password set
		DB:       conf.Db,       // use default DB
		PoolSize: conf.Connections,
	})

	if _, err = conf.client.Ping().Result(); err != nil {
		return nil, ErrorPingFailed.New(err)
	}

	return &conf, nil
}

// Output event
func (t *OutputConfig) Output(ctx context.Context, event logevent.LogEvent) (err error) {
	raw, err := event.MarshalJSON()
	if err != nil {
		return ErrorEventMarshalFailed1.New(err, event)
	}

	key := event.Format(t.Key)

	// try to log forever
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		switch t.DataType {
		case "list":
			if _, err = t.client.RPush(key, raw).Result(); err == nil {
				return
			}
		case "channel":
			if _, err = t.client.Publish(key, string(raw)).Result(); err == nil {
				return
			}
		case "key-value":
			if _, err = t.client.Set(key, string(raw), time.Duration(t.Ttl)*time.Second).Result(); err == nil {
				return
			}
		default:
			return ErrorUnsupportedDataType1.New(nil, t.DataType)
		}

		timeout := time.Duration(t.ReconnectInterval) * time.Second
		timeutil.ContextSleep(ctx, timeout)
	}
}

func (t *OutputConfig) IsRunning() (bool, error) {
	return true, nil
}
