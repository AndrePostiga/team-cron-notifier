package pullRequest

type Author struct {
	name      string
	avatarUrl string
}

func NewAuthor(name string, avatarUrl string) Author {
	return Author{
		name:      name,
		avatarUrl: avatarUrl,
	}
}

func (a *Author) Name() string {
	return a.name
}

func (a *Author) AvatarUrl() string {
	return a.avatarUrl
}
