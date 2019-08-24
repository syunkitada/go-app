package authproxy

import (
	"os"
	"path/filepath"

	"github.com/golang/glog"
	"github.com/spf13/cobra"

	"github.com/syunkitada/goapp/pkg/authproxy/authproxy_api"
	"github.com/syunkitada/goapp/pkg/authproxy/config"
	"github.com/syunkitada/goapp/pkg/authproxy/ctl"
	"github.com/syunkitada/goapp/pkg/base/base_config"
	"github.com/syunkitada/goapp/pkg/lib/logger"
)

var baseConf base_config.Config
var appConf config.Config

var rootCmd = &cobra.Command{}

func Main() {
	if err := rootCmd.Execute(); err != nil {
		glog.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initMain)
	base_config.InitFlags(rootCmd, &baseConf)

	rootCmd.AddCommand(authproxy_api.RootCmd)
	rootCmd.AddCommand(ctl.RootCmd)
}

func initMain() {
	os.Setenv("LANG", "en_US.UTF-8")
	baseConf.BaseDir = filepath.Join(os.Getenv("HOME"), ".goapp")
	baseConf.LogTimeFormat = "2006-01-02T15:04:05Z09:00"
	base_config.InitConfig(&baseConf, &appConf)
	logger.InitLogger(&baseConf)
}
