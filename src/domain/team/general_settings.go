package team

import "errors"

type Config interface {
	GetConfig() any
}

type GeneralSettings struct {
	configs []Config
}

func NewGeneralSettings(configs ...Config) GeneralSettings {
	return GeneralSettings{
		configs: configs,
	}
}

func (gs GeneralSettings) GetSlackConfig() (SlackConfig, error) {
	for _, c := range gs.configs {
		if sc, ok := c.(SlackConfig); ok {
			return sc, nil
		}
	}

	return SlackConfig{}, errors.New("Slack config not found")
}
