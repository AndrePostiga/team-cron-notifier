package team

type PrNotification struct {
	Feature
	repositories      []string
	userExclusionList []string
}

func NewPrNotification(enable bool, repositories []string, excludePrFromUsers []string) PrNotification {
	return PrNotification{
		Feature:           Feature{enable: enable},
		repositories:      repositories,
		userExclusionList: excludePrFromUsers,
	}
}

func (p PrNotification) Repositories() []string {
	return p.repositories
}

func (p PrNotification) UserExclusionList() []string {
	return p.userExclusionList
}
