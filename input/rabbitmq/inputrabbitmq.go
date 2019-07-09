const ModuleName = "rabbitmq"

type InputConfig struct {
	config.InputConfig

	Host string `json:"host"`
	Port int `json:"port"`
	Queue string `json:"queue"`
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
		}
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
}