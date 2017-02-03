package config

import (
	"flag"
	"log"
	"os"
	"github.com/go-gcfg/gcfg.git"
)

type ApplicationConfig struct {
	General struct {
		Listen string
	}
	Cluster struct {
		Server    string
		Arguments string
	}
	Storage struct {
		Address  string
		Password string
		Database int64
	}
}

type ApplicationFlags struct {
	ConfigFile string
}

var application_flags *ApplicationFlags
var application_config *ApplicationConfig

func GetConfig() *ApplicationConfig {
	if application_config == nil {
		application_config = &ApplicationConfig{}

		application_config.General.Listen = "0.0.0.0:8080"
		application_config.Cluster.Server = "./ioq3ded.x86_64"
		application_config.Cluster.Arguments = "+set net_ip $address +set net_port $port"

		if err := gcfg.ReadFileInto(application_config, GetFlags().ConfigFile); err != nil {
			log.Print(err)
			os.Exit(1)
		}
	}
	return application_config
}

func GetFlags() *ApplicationFlags {
	if application_flags == nil {
		application_flags = &ApplicationFlags{}
		flag.StringVar(&application_flags.ConfigFile, "c", "config.ini", "config file")
		flag.Parse()
	}
	return application_flags
}
