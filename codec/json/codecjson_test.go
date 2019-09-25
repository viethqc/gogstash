package codecjson

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/viethqc/gogstash/config/goglog"
	"github.com/viethqc/gogstash/config/logevent"
)

func init() {
	goglog.Logger.SetLevel(logrus.DebugLevel)
}

func TestDecode(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(assert)
	require := require.New(t)
	require.NotNil(require)

	ctx := context.Background()
	codec, err := InitHandler(ctx, nil)
	require.NoError(err)

	msgChan := make(chan logevent.LogEvent, 1)

	ok, err := codec.Decode(ctx, []byte(`{"foo":"bar"}`), nil, msgChan)
	require.NoError(err)
	assert.True(ok)
	require.Len(msgChan, 1)
	event := <-msgChan
	assert.Equal(map[string]interface{}{"foo": "bar"}, event.Extra)
	assert.Equal("", event.Message)

	// string should be ok
	ok, err = codec.Decode(ctx, `{"foo":"bar"}`, nil, msgChan)
	require.NoError(err)
	assert.True(ok)
	require.Len(msgChan, 1)
	event = <-msgChan
	assert.Equal(map[string]interface{}{"foo": "bar"}, event.Extra)
	assert.Equal("", event.Message)

	// map[string]interface{} should be ok
	ok, err = codec.Decode(ctx, map[string]interface{}{"foo": "bar"}, nil, msgChan)
	require.NoError(err)
	assert.True(ok)
	require.Len(msgChan, 1)
	event = <-msgChan
	assert.Equal(map[string]interface{}{"foo": "bar"}, event.Extra)
	assert.Equal("", event.Message)

	// ok will be true, as message sent
	ok, err = codec.Decode(ctx, []byte(`{"foo":"bar"dr.who}`), nil, msgChan)
	require.Error(err) // fail to decode
	assert.True(ok)
	require.Len(msgChan, 1)
	event = <-msgChan
	assert.Equal([]string{ErrorTag}, event.Tags)
	assert.Equal(`{"foo":"bar"dr.who}`, event.Message)

	// timestamp will be parsed
	ts := time.Date(2019, time.January, 4, 0, 55, 36, 0, time.UTC)
	ok, err = codec.Decode(ctx, []byte(`{"@timestamp":"`+ts.Format(time.RFC3339)+`"}`), nil, msgChan)
	require.NoError(err)
	assert.True(ok)
	require.Len(msgChan, 1)
	event = <-msgChan
	assert.Equal(ts, event.Timestamp)

	// message will map to event.Message
	ok, err = codec.Decode(ctx, []byte(`{"message":"hello"}`), nil, msgChan)
	require.NoError(err)
	assert.True(ok)
	require.Len(msgChan, 1)
	event = <-msgChan
	assert.Equal("hello", event.Message)

	// tags will map to event.Tags
	ok, err = codec.Decode(ctx, []byte(`{"tags":["foo","bar"]}`), nil, msgChan)
	require.NoError(err)
	assert.True(ok)
	require.Len(msgChan, 1)
	event = <-msgChan
	assert.Equal([]string{"foo", "bar"}, event.Tags)

	// merge & override extra
	ok, err = codec.Decode(ctx, []byte(`{"foo":"bar2"}`), map[string]interface{}{
		"foo": "bar",
		"one": "more thing",
	}, msgChan)
	require.NoError(err)
	assert.True(ok)
	require.Len(msgChan, 1)
	event = <-msgChan
	assert.Equal(map[string]interface{}{
		"foo": "bar2",
		"one": "more thing",
	}, event.Extra)
}
