package conf

import (
	"reflect"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/tsingshaner/go-pkg/log/console"
)

type Options struct {
	// FilePath the path to the configuration file (default: config.json)
	FilePath string
	// Silence silence the output of config loading messages
	Silence bool
}

type config[T any] struct {
	Value   *T
	options *Options
	viper   *viper.Viper
}

// Read a simple function to read the configuration from the configuration file,
// will exit if load error
//
// accepted cmd args:
// --config  the path to the configuration file
// --silence silence the output of config loading
func Read[T any]() *T {
	c := New(new(T), ParseArgs())
	if err := c.Load(); err != nil {
		console.Fatal("load config error: %+v", err)
	}

	return c.Value
}

func New[T any](conf *T, opts *Options) *config[T] {
	if reflect.ValueOf(conf).Elem().Kind() != reflect.Struct {
		console.Fatal("store must be a struct ptr")
	}

	return &config[T]{conf, opts, viper.New()}
}

// Load read & unmarshal config from configuration file
func (c *config[T]) Load() error {
	c.viper.SetConfigFile(c.options.FilePath)

	if err := c.viper.ReadInConfig(); err != nil {
		return err
	}

	return c.unmarshal()
}

func (c *config[T]) unmarshal() error {
	if err := c.viper.Unmarshal(c.Value); err != nil {
		return err
	}

	if !c.options.Silence {
		showConfig(c.Value)
	}

	return nil
}

type Event = fsnotify.Event

// Watch watch the configuration file changes
// if you want to reload the configuration, you should call the Load method in the listener
func (c *config[T]) Watch(listener func(Event)) {
	c.viper.OnConfigChange(listener)
	c.viper.WatchConfig()
}
