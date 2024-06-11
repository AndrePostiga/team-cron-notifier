package team

type SlackConfig struct {
	webhookSecretEnvName string
}

func NewSlackConfig(webhookSecretEnvName string) SlackConfig {
	return SlackConfig{
		webhookSecretEnvName: webhookSecretEnvName,
	}
}

func (sc SlackConfig) GetConfig() any {
	return sc
}

func (sc SlackConfig) GetWebhookSecretEnvName() string {
	return sc.webhookSecretEnvName
}
