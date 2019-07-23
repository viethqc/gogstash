package inputbeats

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"github.com/elastic/go-lumber/server"
	reuse "github.com/libp2p/go-reuseport"
	codecjson "github.com/viethqc/gogstash/codec/json"
	"github.com/viethqc/gogstash/config"
	"github.com/viethqc/gogstash/config/goglog"
	"github.com/viethqc/gogstash/config/logevent"
)

// ModuleName is the name used in config file
const ModuleName = "beats"

// InputConfig holds the configuration json fields and internal objects
type InputConfig struct {
	config.InputConfig

	// The IP address to listen on, defaults to "0.0.0.0"
	Host string `json:"host"`
	// The port to listen on.
	Port int `json:"port"`
	// Here we enable SO_REUSEPORT, see more information:
	// https://github.com/libp2p/go-reuseport
	ReusePort bool `json:"reuseport"`

	// Enable ssl transport, defaults to false
	SSL bool `json:"ssl"`
	// SSL certificate to use.
	SSLCertificate string `json:"ssl_certificate"`
	// SSL key to use.
	SSLKey string `json:"ssl_key"`
	// SSL Verify, defaults to false
	SSLVerify bool `json:"ssl_verify"`

	tlsConfig *tls.Config
}

// DefaultInputConfig returns an InputConfig struct with default values
func DefaultInputConfig() InputConfig {
	return InputConfig{
		InputConfig: config.InputConfig{
			CommonConfig: config.CommonConfig{
				Type: ModuleName,
			},
		},
		Host: "0.0.0.0",
	}
}

// InitHandler initialize the input plugin
func InitHandler(ctx context.Context, raw *config.ConfigRaw) (config.TypeInputConfig, error) {
	conf := DefaultInputConfig()
	err := config.ReflectConfig(raw, &conf)
	if err != nil {
		return nil, err
	}

	if !conf.SSL {
		if conf.SSLCertificate != "" {
			goglog.Logger.Warn("beats input: SSL Certificate will not be used")
		}
		if conf.SSLKey != "" {
			goglog.Logger.Warn("beats input: SSL Key will not be used")
		}
	} else {
		// SSL enabled
		cer, err := tls.LoadX509KeyPair(conf.SSLCertificate, conf.SSLKey)
		if err != nil {
			return nil, err
		}
		conf.tlsConfig = &tls.Config{Certificates: []tls.Certificate{cer}}
		if !conf.SSLVerify {
			conf.tlsConfig.InsecureSkipVerify = true
		}
	}

	conf.Codec, err = config.GetCodecDefault(ctx, *raw, codecjson.ModuleName)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

// Start wraps the actual function starting the plugin
func (t *InputConfig) Start(ctx context.Context, msgChan chan<- logevent.LogEvent) (err error) {
	addr := fmt.Sprintf("%s:%d", t.Host, t.Port)
	s, err := server.ListenAndServeWith(func(network, addr string) (l net.Listener, err error) {
		if t.ReusePort {
			l, err = reuse.Listen(network, addr)
		} else {
			l, err = net.Listen(network, addr)
		}
		if err != nil {
			return nil, err
		}
		if t.SSL {
			l = tls.NewListener(l, t.tlsConfig)
		}
		return l, err
	}, addr, server.JSONDecoder(t.Codec.DecodeEvent))
	if err != nil {
		return err
	}
	defer s.Close()
	goglog.Logger.Infof("beats input: start listening on %s", addr)

	for {
		select {
		case <-ctx.Done():
			goglog.Logger.Info("input beats stopped")
			return nil
		case data := <-s.ReceiveChan():
			for _, e := range data.Events {
				msgChan <- e.(logevent.LogEvent)
			}
			data.ACK()
		}
	}
	return nil
}
