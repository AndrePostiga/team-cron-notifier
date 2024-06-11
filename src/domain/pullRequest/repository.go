package pullRequest

type Repository struct {
	name string
}

func NewRepository(name string) Repository {
	return Repository{
		name: name,
	}
}

func (r Repository) Name() string {
	return r.name
}
