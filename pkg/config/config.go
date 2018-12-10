package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

var (
	Conf              Config
	configDir         string
	configFile        string
	enableDebug       bool
	enableDevelop     bool
	enableDatabaseLog bool
)

var (
	glogV               int
	glogLogtostderr     bool
	glogStderrthreshold int
	glogAlsologtostderr bool
	glogVmodule         string
	glogLogDir          string
	glogLogBacktraceAt  string
)

func InitFlags(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().StringVar(&configDir, "config-dir", "", "config directory (default is $PWD/ci/etc)")
	rootCmd.PersistentFlags().StringVar(&configFile, "config-file", "", "config file (default is config.toml)")
	rootCmd.PersistentFlags().BoolVar(&enableDebug, "debug", false, "enable debug mode")
	rootCmd.PersistentFlags().BoolVar(&enableDevelop, "develop", false, "enable develop mode")
	rootCmd.PersistentFlags().BoolVar(&enableDatabaseLog, "database-log", false, "enable database logging")

	// glog flags
	rootCmd.PersistentFlags().IntVar(&glogV, "glog-v", 0, "log level for V logs")
	rootCmd.PersistentFlags().BoolVar(&glogLogtostderr, "glog-logtostderr", true, "log to standard error instead of files")
	rootCmd.PersistentFlags().IntVar(&glogStderrthreshold, "glog-stderrthreshold", 0, "logs at or above this threshold go to stderr")
	rootCmd.PersistentFlags().BoolVar(&glogAlsologtostderr, "glog-alsologtostderr", false, "log to standard error as well as files")
	rootCmd.PersistentFlags().StringVar(&glogVmodule, "glog-vmodule", "", "comma-separated list of pattern=N settings for file-filtered logging")
	rootCmd.PersistentFlags().StringVar(&glogLogDir, "glog-log-dir", "", "If non-empty, write log files in this directory")
	rootCmd.PersistentFlags().StringVar(&glogLogBacktraceAt, "glog-log-backtrace-at", ":0", "when logging hits line file:N, emit a stack trace")
}

func InitConfig() {
	_ = flag.CommandLine.Parse([]string{})
	flagShim(map[string]string{
		"v":                fmt.Sprint(glogV),
		"logtostderr":      fmt.Sprint(glogLogtostderr),
		"stderrthreshold":  fmt.Sprint(glogStderrthreshold),
		"alsologtostderr":  fmt.Sprint(glogAlsologtostderr),
		"vmodule":          glogVmodule,
		"log_dir":          glogLogDir,
		"log_backtrace_at": glogLogBacktraceAt,
	})

	if configDir == "" {
		configDir = os.Getenv("CONFIG_DIR")
		if configDir == "" {
			pwd := os.Getenv("PWD")
			configDir = filepath.Join(pwd, "ci", "etc")
		}
	}

	if configFile == "" {
		configFile = os.Getenv("CONFIG_FILE")
		if configFile == "" {
			configFile = "config.toml"
		}
	}

	if enableDebug {
		enableDatabaseLog = true
	}

	var err error

	hostname, err := os.Hostname()
	if err != nil {
		glog.Fatal(err)
	}

	defaultConfig := DefaultConfig{
		Host:              hostname,
		ConfigDir:         configDir,
		ConfigFile:        filepath.Join(configDir, configFile),
		EnableDebug:       enableDebug,
		EnableDevelop:     enableDevelop,
		EnableDatabaseLog: enableDatabaseLog,
	}

	err = loadConfig(&defaultConfig)
	if err != nil {
		glog.Fatal(err)
	}
}

func flagShim(fakeVals map[string]string) {
	flag.VisitAll(func(fl *flag.Flag) {
		if val, ok := fakeVals[fl.Name]; ok {
			fl.Value.Set(val)
		}
	})
}

func loadConfig(defaultConfig *DefaultConfig) error {
	newConfig := newConfig(defaultConfig)
	_, err := toml.DecodeFile(defaultConfig.ConfigFile, newConfig)
	if err != nil {
		return err
	}
	Conf = *newConfig

	return nil
}
