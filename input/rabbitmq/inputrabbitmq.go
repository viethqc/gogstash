package inputrabbitmq

import (
	"context"
	"fmt"
	"time"

	//"github.com/viethqc/gogstash/KDGoLib/errutil"
	"github.com/streadway/amqp"
	codecjson "github.com/viethqc/gogstash/codec/json"
	"github.com/viethqc/gogstash/config"
	"github.com/viethqc/gogstash/config/goglog"
	"github.com/viethqc/gogstash/config/logevent"
)

const ModuleName = "rabbitmq"

type InputConfig struct {
	config.InputConfig

	Host          string `json:"host"`
	Port          int    `json:"port"`
	Queue         string `json:"queue"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	PrefetchCount int    `json:"prefetch_count"`
}

func DefaultInputConfig() InputConfig {
	return InputConfig{
		InputConfig: config.InputConfig{
			CommonConfig: config.CommonConfig{
				Type: ModuleName,
			},
		},
		PrefetchCount: 1,
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

	err = ch.Qos(
		t.PrefetchCount, // prefetch count
		0,               // prefetch size
		false,           // global
	)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		goglog.Logger.Infof("Failed to register a consumer: %s", err.Error())
		return err
	}

	for msg := range msgs {
		goglog.Logger.Info(string(msg.Body))
		if ok, err := t.Codec.Decode(ctx, string(msg.Body), nil, msgChan); ok == true && err == nil {
			msg.Ack(false)
		} else {
			msg.Nack(false, true)
			time.Sleep(5 * time.Second)
		}
	}

	return nil
}
