package inputrabbitmq

import (
	"context"
	"fmt"

	//"github.com/tsaikd/KDGoLib/errutil"
	"github.com/streadway/amqp"
	codecjson "github.com/tsaikd/gogstash/codec/json"
	"github.com/tsaikd/gogstash/config"
	"github.com/tsaikd/gogstash/config/goglog"
	"github.com/tsaikd/gogstash/config/logevent"
)

const ModuleName = "rabbitmq"

type InputConfig struct {
	config.InputConfig

	Host     string `json:"host"`
	Port     int    `json:"port"`
	Queue    string `json:"queue"`
	Username string `json:"username"`
	// SSL Verify, defaults to false
	Password string `json:"password"`
}

func DefaultInputConfig() InputConfig {
	return InputConfig{
		InputConfig: config.InputConfig{
			CommonConfig: config.CommonConfig{
				Type: ModuleName,
			},
		},
	}
}

// InitHandler initialize the input plugin
func InitHandler(ctx context.Context, raw *config.ConfigRaw) (config.TypeInputConfig, error) {
	conf := DefaultInputConfig()
	err := config.ReflectConfig(raw, &conf)
	if err != nil {
		return nil, err
	}

	conf.Codec, err = config.GetCodecDefault(ctx, *raw, codecjson.ModuleName)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

// Start wraps the actual function starting the plugin
func (t *InputConfig) Start(ctx context.Context, msgChan chan<- logevent.LogEvent) (err error) {
	goglog.Logger.Infof("start rabbitmq")

	authen := fmt.Sprintf("amqp://%s:%s@%s:%d/", t.Username, t.Password, t.Host, t.Port)
	conn, err := amqp.Dial(authen)
	if err != nil {
		goglog.Logger.Infof("Failed to connect to RabbitMQ: %s", err.Error())
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		goglog.Logger.Infof("Failed to open a channel: %s", err.Error())
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		t.Queue, // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		goglog.Logger.Infof("Failed to declare a queue: %s", err.Error())
		return err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		goglog.Logger.Infof("Failed to register a consumer: %s", err.Error())
		return err
	}

	for {
		msg := <-msgs
		goglog.Logger.Info(string(msg.Body))
		t.Codec.Decode(ctx, string(msg.Body), nil, msgChan)
	}

	return nil
}
