package filterratelimit

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/viethqc/gogstash/config"
	"github.com/viethqc/gogstash/config/goglog"
	"github.com/viethqc/gogstash/config/logevent"
)

func init() {
	goglog.Logger.SetLevel(logrus.DebugLevel)
	config.RegistFilterHandler(ModuleName, InitHandler)
}

func Test_filter_ratelimit_module(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	ctx := context.Background()
	conf, err := config.LoadFromYAML([]byte(strings.TrimSpace(`
debugch: true
filter:
  - type: rate_limit
    rate: 10
    burst: 1
	`)))
	require.NoError(err)
	require.NoError(conf.Start(ctx))

	start := time.Now()

	conf.TestInputEvent(logevent.LogEvent{})
	conf.TestInputEvent(logevent.LogEvent{})
	conf.TestInputEvent(logevent.LogEvent{})
	conf.TestInputEvent(logevent.LogEvent{})
	conf.TestInputEvent(logevent.LogEvent{})

	_, err = conf.TestGetOutputEvent(100 * time.Millisecond)
	require.NoError(err)
	_, err = conf.TestGetOutputEvent(100 * time.Millisecond)
	require.NoError(err)
	_, err = conf.TestGetOutputEvent(100 * time.Millisecond)
	require.NoError(err)
	_, err = conf.TestGetOutputEvent(100 * time.Millisecond)
	require.NoError(err)
	_, err = conf.TestGetOutputEvent(100 * time.Millisecond)
	require.NoError(err)

	require.WithinDuration(start.Add(400*time.Millisecond), time.Now(), 150*time.Millisecond)
}

func Test_filter_ratelimit_module_burst(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	ctx := context.Background()
	conf, err := config.LoadFromYAML([]byte(strings.TrimSpace(`
debugch: true
filter:
  - type: rate_limit
    rate: 10
    burst: 4
	`)))
	require.NoError(err)
	require.NoError(conf.Start(ctx))

	time.Sleep(600 * time.Millisecond)

	start := time.Now()

	conf.TestInputEvent(logevent.LogEvent{})
	conf.TestInputEvent(logevent.LogEvent{})
	conf.TestInputEvent(logevent.LogEvent{})
	conf.TestInputEvent(logevent.LogEvent{})
	conf.TestInputEvent(logevent.LogEvent{})
	conf.TestInputEvent(logevent.LogEvent{})

	_, err = conf.TestGetOutputEvent(100 * time.Millisecond)
	require.NoError(err)
	_, err = conf.TestGetOutputEvent(100 * time.Millisecond)
	require.NoError(err)
	_, err = conf.TestGetOutputEvent(100 * time.Millisecond)
	require.NoError(err)
	_, err = conf.TestGetOutputEvent(100 * time.Millisecond)
	require.NoError(err)
	_, err = conf.TestGetOutputEvent(100 * time.Millisecond)
	require.NoError(err)
	_, err = conf.TestGetOutputEvent(100 * time.Millisecond)
	require.NoError(err)

	require.WithinDuration(start.Add(150*time.Millisecond), time.Now(), 100*time.Millisecond)
}

func Test_filter_ratelimit_module_delay(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	ctx := context.Background()
	conf, err := config.LoadFromYAML([]byte(strings.TrimSpace(`
debugch: true
filter:
  - type: rate_limit
    rate: 10
    burst: 1
	`)))
	require.NoError(err)
	require.NoError(conf.Start(ctx))

	time.Sleep(500 * time.Millisecond)

	start := time.Now()

	conf.TestInputEvent(logevent.LogEvent{})
	conf.TestInputEvent(logevent.LogEvent{})
	conf.TestInputEvent(logevent.LogEvent{})
	conf.TestInputEvent(logevent.LogEvent{})
	conf.TestInputEvent(logevent.LogEvent{})

	_, err = conf.TestGetOutputEvent(100 * time.Millisecond)
	require.NoError(err)
	_, err = conf.TestGetOutputEvent(100 * time.Millisecond)
	require.NoError(err)
	_, err = conf.TestGetOutputEvent(100 * time.Millisecond)
	require.NoError(err)
	_, err = conf.TestGetOutputEvent(100 * time.Millisecond)
	require.NoError(err)
	_, err = conf.TestGetOutputEvent(100 * time.Millisecond)
	require.NoError(err)

	require.WithinDuration(start.Add(350*time.Millisecond), time.Now(), 150*time.Millisecond)
}
