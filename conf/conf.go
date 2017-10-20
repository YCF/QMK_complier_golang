package conf

import (
	"fmt"

	"github.com/go-ini/ini"
)

// GetOption n
func GetOption(section string, key string) (conf string) {
	cfg, err := ini.InsensitiveLoad("conf/conf.ini")
	conf = cfg.Section(section).Key(key).String()
	if err != nil {
		fmt.Println(err)
	}

	return conf
}
