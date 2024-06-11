package pullRequest

type ReviewState string

const (
	approved         ReviewState = "APPROVED"
	changesRequested             = "CHANGES_REQUESTED"
	comment                      = "COMMENT"
)

type Review struct {
	state         ReviewState
	approvesCount uint
}

func NewReview(state ReviewState) Review {
	return Review{
		state: state,
	}
}
