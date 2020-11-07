package server

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

type ServerLogLevel string
const (
	Debug   ServerLogLevel = "DEBUG"
	Info                   = "INFO"
	Warning                = "WARNING"
	Error                  = "ERROR"
)

type ServerConfig struct {
	LogLevel               ServerLogLevel `json:"logLevel"`
	LogPath                string         `json:"logPath"`
	ShotgunConnectionRetry int            `json:"shotgunRetry"`
	ShotgunServerAddress   string         `json:"shotgunAddress"`
	ShotgunScriptName      string         `json:"shotgunScriptName"`
	ShotgunScriptKey       string         `json:"shotgunScriptKey"`
	PluginPath             string         `json:"pluginPath"`
	SMTPServerAddress      string         `json:"smtpAddress"`
	SMTPServerPort         int            `json:"smtpPort"`
	SMTPUsername           string         `json:"smtpUsername"`
	SMTPPassword           string         `json:"smtpPassword"`
	EmailFrom string `json:"emailFrom"`
	EmailTo string `json:"emailTo"`
	EmailSubject string `json:"emailSubject"`
}

func NewServerConfigFromFile(configPath string) (*ServerConfig, error) {
	conf := ServerConfig{}
	configData, _ := ioutil.ReadFile(configPath)
	err := json.Unmarshal(configData, &conf)
	if err != nil {
		logrus.WithError(err).Error("failed to unmarshal server configuration file")
		return nil, err
	}
	return &conf, nil
}