package configurator

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func InitConfigs(configurationServiceName string, path []string, config interface{}) interface{} {

	log.Println("Loading configurations...")

	viper.SetConfigName(configurationServiceName)

	for _, r := range path {
		viper.AddConfigPath(r)
	}
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	for _, k := range viper.AllKeys() {
		value := viper.GetString(k)
		if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
			viper.Set(k, getEnvOrDefault(strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}")))
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}
	return config
}

func getEnvOrDefault(env string) string {
	var res string
	if idx := strings.IndexByte(env, ':'); idx >= 0 {
		res = os.Getenv(env[:strings.IndexByte(env, ':')])
		if len(res) == 0 {
			res = env[(strings.IndexByte(env, ':') + 1):]
		}
	} else {
		res = os.Getenv(env)
	}
	return res
}
