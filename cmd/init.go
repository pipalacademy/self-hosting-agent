package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/valyala/fasthttp"
)

// initLogger initializes logger.
func initLogger(ko *koanf.Koanf) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	})
	if ko.String("app.log") == "debug" {
		logger.SetLevel(logrus.DebugLevel)
	}
	return logger
}

// initConfig loads config to `ko`
// object.
func initConfig(cfgDefault string, envPrefix string) (*koanf.Koanf, error) {
	var (
		ko = koanf.New(".")
		f  = flag.NewFlagSet("front", flag.ContinueOnError)
	)

	// Configure Flags.
	f.Usage = func() {
		fmt.Println(f.FlagUsages())
		os.Exit(0)
	}

	// Register `--config` flag.
	cfgPath := f.String("config", cfgDefault, "Path to a config file to load.")

	// Parse and Load Flags.
	err := f.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}

	// Load the config files from the path provided.
	err = ko.Load(file.Provider(*cfgPath), toml.Parser())
	if err != nil {
		return nil, err
	}

	// Load environment variables if the key is given
	// and merge into the loaded config.
	if envPrefix != "" {
		err = ko.Load(env.Provider(envPrefix, ".", func(s string) string {
			return strings.Replace(strings.ToLower(
				strings.TrimPrefix(s, envPrefix)), "__", ".", -1)
		}), nil)
		if err != nil {
			return nil, err
		}
	}

	return ko, nil
}

// initHTTPServer initialises HTTP server.
func initHTTPServer(ko *koanf.Koanf) *fasthttp.Server {
	return &fasthttp.Server{
		Name:               ko.String("server.name"),
		ReadTimeout:        ko.Duration("server.read_timeout"),
		WriteTimeout:       ko.Duration("server.write_timeout"),
		IdleTimeout:        ko.Duration("server.keepalive_timeout"),
		MaxRequestBodySize: ko.Int("server.max_body_size"),
		ReadBufferSize:     ko.Int("server.read_buffer_size"),
	}
}

func initOpts(ko *koanf.Koanf) Opts {
	return Opts{
		ServerAddr: ko.String("server.address"),
	}
}
