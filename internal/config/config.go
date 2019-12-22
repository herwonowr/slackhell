package config

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config ...
type Config struct {
	Account struct {
		ID       string
		RealName string
	}
	Slack struct {
		Token string
	}
	Database struct {
		Path string
	}
	Log struct {
		Debug bool
	}
}

// NewConfig ...
func NewConfig(cmd *cobra.Command) (*Config, error) {
	cfg := new(Config)
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return nil, err
	}

	viper.SetConfigType("toml")
	viper.SetEnvPrefix("SLACKHELL")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if configFile, _ := cmd.Flags().GetString("config"); configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("slackhell")
		viper.AddConfigPath("./config")
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := checkConfig(viper.GetViper()); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func checkConfig(v *viper.Viper) error {
	if !v.IsSet("account.id") {
		return errors.New("missing slack account id configuration, for initialization")
	}
	if !v.IsSet("account.realname") {
		return errors.New("missing slack account real name configuration, for initialization")
	}
	if !v.IsSet("slack.token") {
		return errors.New("missing slack token configuration")
	}
	if !v.IsSet("database.path") {
		return errors.New("missing database file path configuration")
	}
	if !v.IsSet("log.debug") {
		return errors.New("missing debug configuration")
	}
	return nil
}
