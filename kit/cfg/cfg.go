package cfg

import (
	"path/filepath"
	"github.com/spf13/viper"
	"fmt"
)

type Config struct {
	v * viper.Viper
}

var c *Config

func Init(path string) error {
	var err error
	if c, err = NewConfig(path); err != nil {
		return err
	}

	return nil
}

func NewConfig(path string) (*Config, error) {

	v := viper.New()

	name := filepath.Base(path)
	v.SetConfigName(name[:len(name)-len(filepath.Ext(name))])

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil,fmt.Errorf("[cfg.NewConfig] read config file error: %s \n", err)
	}

	v.AddConfigPath(filepath.Dir(absPath))
	err = v.ReadInConfig()
	if err != nil {
		return nil,fmt.Errorf("[cfg.NewConfig] read config file error: %s \n", err)
	}

	c := &Config{
		v : v,
	}

	return c,nil
}

func GetString(key string) string {
	return c.GetString(key)
}
func (c *Config) GetString(key string) string {
	return c.v.GetString(key)
}
