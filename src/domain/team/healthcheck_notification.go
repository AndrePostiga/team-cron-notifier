package team

type HealthCheckNotification struct {
	Feature
	urls           []string
	authentication struct {
		authType     string
		authToken    string
		authUser     string
		authPassword string
	}
}

func NewHealthCheckNotification(enable bool, urls []string) HealthCheckNotification {
	return HealthCheckNotification{
		Feature: Feature{enable: enable},
		urls:    urls,
	}
}
