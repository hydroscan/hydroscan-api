package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"strings"
)

// pod name is something like "hydroscan-6bff5fc87d-l864j"
// this method will remove the tailing hash, then return "hydroscan"
func getPodName() string {
	parts := strings.Split(viper.GetString("pod.name"), "-")
	return strings.Join(parts[:len(parts)-2], "-")
}

func Load() {
	// All capitalized ENVs with prefix 'HYDROSCAN' will have the most priority.
	viper.SetEnvPrefix("HYDROSCAN")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// Try to load config in local config file. This is used for development.
	// .config.yml should never be committed in git.
	viper.SetConfigFile(".config.yml")
	err := viper.ReadInConfig()

	if err == nil {
		log.Printf("Coinfigs are loaded from local .config.yml")
		return
	}

	if viper.GetString("etcd_url") != "" {
		etcd_url := viper.GetString("etcd_url")

		etcd_config_file := fmt.Sprintf("/k8s-app-configs/%s/%s", viper.GetString("pod.namespace"), getPodName())

		log.Printf("Loading Coinfigs from etcd: %s, %s", etcd_url, etcd_config_file)

		err = viper.AddRemoteProvider("etcd", etcd_url, etcd_config_file)

		if err != nil {
			panic(err)
		}

		viper.SetConfigType("yaml")

		err = viper.ReadRemoteConfig()

		if err != nil {
			panic(err)
		}

		log.Printf("Coinfigs are loaded from etcd: %s, %s", etcd_url, etcd_config_file)
		return
	}
}
