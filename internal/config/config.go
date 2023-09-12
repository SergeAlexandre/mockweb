package config

import (
	"github.com/go-logr/logr"
	"github.com/spf13/pflag"
	"mockweb/pkg/misc"
	"mockweb/pkg/skserver"
)

type Config struct {
	Log    misc.LogConfig  `yaml:"log"`
	Server skserver.Config `yaml:"server"`
}

var Conf Config
var Log logr.Logger

func Setup() {
	pflag.StringVar(&Conf.Server.Interface, "interface", "0.0.0.0", "Listening interface")
	pflag.IntVar(&Conf.Server.Port, "port", 0, "Listening port (Default 80 or 443")
	pflag.StringVar(&Conf.Server.CertDir, "cerDir", "", "Directory containing server key and cert. Clear text if not defined")
	pflag.StringVar(&Conf.Server.CertName, "certName", "tls.crt", "Server certificate file name")
	pflag.StringVar(&Conf.Server.KeyName, "keyName", "tls.key", "Server key file name")
	pflag.StringVar(&Conf.Log.Level, "logLevel", "INFO", "Log level (PANIC|FATAL|ERROR|WARN|INFO|DEBUG|TRACE)")
	pflag.StringVar(&Conf.Log.Mode, "logMode", "json", "Log mode: 'dev' or 'json'")

	pflag.Parse()

	Conf.Server.Ssl = Conf.Server.CertDir != ""

	if Conf.Server.Port == 0 {
		if Conf.Server.Ssl {
			Conf.Server.Port = 443
		} else {
			Conf.Server.Port = 80
		}
	}

	var err error
	Log, err = misc.HandleLog(&Conf.Log)
	if err != nil {
		panic(err)
	}

}
