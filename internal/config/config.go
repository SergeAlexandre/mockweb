package config

import (
	"github.com/go-logr/logr"
	"github.com/spf13/pflag"
	"mockweb/pkg/misc"
	"mockweb/pkg/skserver"
	"time"
)

type Config struct {
	Log             misc.LogConfig  `yaml:"log"`
	Server          skserver.Config `yaml:"server"`
	NoCache         bool            `yaml:"noCache"`
	Name            string          `yaml:"name"`
	IdleTimeout     time.Duration   `yaml:"idleTimeout"`
	SessionLifetime time.Duration   `yaml:"sessionLifetime"`
}

var Conf Config
var Log logr.Logger

func Setup() {
	var idleTimeout string
	var sessionLifetime string

	pflag.StringVar(&Conf.Server.Interface, "interface", "0.0.0.0", "Listening interface")
	pflag.IntVar(&Conf.Server.Port, "port", 0, "Listening port (Default 80 or 443")
	pflag.StringVar(&Conf.Server.CertDir, "certDir", "", "Directory containing server key and cert. Clear text if not defined")
	pflag.StringVar(&Conf.Server.CertName, "certName", "tls.crt", "Server certificate file name")
	pflag.StringVar(&Conf.Server.KeyName, "keyName", "tls.key", "Server key file name")
	pflag.StringVar(&Conf.Log.Level, "logLevel", "INFO", "Log level (PANIC|FATAL|ERROR|WARN|INFO|DEBUG|TRACE)")
	pflag.StringVar(&Conf.Log.Mode, "logMode", "json", "Log mode: 'dev' or 'json'")
	pflag.StringVar(&Conf.Name, "name", "MOCKWEB", "Server name")
	pflag.BoolVar(&Conf.NoCache, "noCache", false, "Add no caching headers")
	pflag.StringVar(&idleTimeout, "idleTimeout", "15m", "The maximum length of time a session can be inactive before being expired")
	pflag.StringVar(&sessionLifetime, "sessionLifetime", "6h", "The absolute maximum length of time that a session is valid.")

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

	// ----------------------- Session handling
	Conf.IdleTimeout, err = time.ParseDuration(idleTimeout)
	if err != nil {
		panic(err)
	}
	Conf.SessionLifetime, err = time.ParseDuration(sessionLifetime)
	if err != nil {
		panic(err)
	}

}
