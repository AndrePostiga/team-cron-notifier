package team

type PrNotification struct {
	Feature
	repositories []string
}

func NewPrNotification(enable bool, repositories []string) PrNotification {
	return PrNotification{
		Feature:      Feature{enable: enable},
		repositories: repositories,
	}
}

func (p PrNotification) Repositories() []string {
	return p.repositories
}
