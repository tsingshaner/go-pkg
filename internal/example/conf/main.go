package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/tsingshaner/go-pkg/conf"
	"github.com/tsingshaner/go-pkg/log/console"
)

type sub struct {
	Key string `json:"key" yaml:"key" toml:"key"`
}

type nestPtr struct {
	Key    string `json:"key" yaml:"key" toml:"key"`
	Sub    sub    `json:"sub" yaml:"sub" toml:"sub"`
	SubPtr *sub   `json:"subPtr" yaml:"subPtr" toml:"subPtr"`
}

type Config struct {
	Ptr   *int     `json:"ptr" yaml:"ptr" toml:"ptr"`
	Level []string `json:"level" yaml:"level" toml:"level"`
	Nest  struct {
		Key string `json:"key" yaml:"key" toml:"key"`
	} `json:"nest" yaml:"nest" toml:"nest"`
	NestPtr *nestPtr `json:"nestPtr" yaml:"nestPtr" toml:"nestPtr"`
}

func main() {
	config := conf.New(&Config{}, conf.ParseArgs())

	config.Load()
	config.Watch(func(e conf.Event) {
		console.Info("配置文件发生变化, 重新加载配置。")
		config.Load()
	})

	println(*config.Value.Ptr)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	console.Info("程序运行中, 按 Ctrl+C 退出。")

	<-sigChan

	console.Info("收到 Ctrl+C, 程序退出。")
}
