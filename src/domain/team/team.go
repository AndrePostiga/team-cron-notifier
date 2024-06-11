package team

type Feature struct {
	enable bool
}

type Team struct {
	name                    string
	prNotification          PrNotification
	healthCheckNotification HealthCheckNotification
	notificationSettings    GeneralSettings
	githubApiToken          string
}

func NewTeam(name string, prs []string, healthchecks []string, notificationConfigs GeneralSettings, githubApiToken string) Team {
	return Team{
		name:                    name,
		prNotification:          NewPrNotification(true, prs),
		healthCheckNotification: NewHealthCheckNotification(true, healthchecks),
		notificationSettings:    notificationConfigs,
		githubApiToken:          githubApiToken,
	}
}

func (t Team) IsPrNotificationEnabled() bool {
	return t.prNotification.Feature.enable
}

func (t Team) Name() string {
	return t.name
}

func (t Team) GithubApiToken() string {
	return t.githubApiToken
}

func (t Team) PrNotification() PrNotification {
	return t.prNotification
}

func (t Team) NotificationSettings() GeneralSettings {
	return t.notificationSettings
}
