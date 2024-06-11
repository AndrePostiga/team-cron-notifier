package pullRequest

import (
	"time"
)

type Labels []string

type Priority string

const (
	High           Priority = "High"
	Medium         Priority = "Medium"
	Low            Priority = "Low"
	ReadyForDeploy Priority = "ReadyForDeploy"
)

type State string

const (
	open   State = "OPEN"
	closed State = "CLOSED"
)

type PullRequest struct {
	title                  string
	number                 uint
	url                    string
	author                 Author
	createdAt              time.Time
	openedDays             uint
	state                  State
	reviews                []Review
	numberOfApproves       uint
	numberOfRequestChanges uint
	priority               Priority
	labels                 Labels
	repository             Repository
}

func NewPullRequest(
	title string,
	number uint,
	url string,
	author Author,
	createdAt time.Time,
	state State,
	labels []string,
	reviews []Review,
	repository Repository) PullRequest {

	pr := PullRequest{
		title:      title,
		number:     number,
		url:        url,
		author:     author,
		createdAt:  createdAt,
		state:      state,
		labels:     labels,
		reviews:    reviews,
		repository: repository,
	}

	pr.calculateOpenedDays()
	pr.calculateNumberOfApprovesAndRequestChanges()
	pr.setPriority()
	return pr
}

func (pr *PullRequest) IsReadyForDeploy() bool {
	return pr.numberOfApproves >= 2
}

func (pr *PullRequest) calculateNumberOfApprovesAndRequestChanges() {

	for _, review := range pr.reviews {
		if review.state == approved {
			pr.numberOfApproves++
		} else if review.state == changesRequested {
			pr.numberOfRequestChanges++
		}
	}
}

func (pr *PullRequest) setPriority() {

	if pr.IsReadyForDeploy() {
		pr.priority = ReadyForDeploy
		return
	}

	precedence := map[Priority]int{
		High:   3,
		Medium: 2,
		Low:    1,
	}

	actualHighest := Low

	for _, label := range pr.labels {
		if rank, found := precedence[Priority(label)]; found {
			if rank > precedence[actualHighest] {
				actualHighest = Priority(label)
			}
		}
	}

	pr.priority = actualHighest
}

func (pr *PullRequest) calculateOpenedDays() {
	pr.openedDays = uint(time.Since(pr.createdAt).Hours() / 24)
}

// Getters

func (pr *PullRequest) URL() string {
	return pr.url
}

func (pr *PullRequest) Title() string {
	return pr.title
}

func (pr *PullRequest) NumberOfApproves() uint {
	return pr.numberOfApproves
}

func (pr *PullRequest) NumberOfRequestChanges() uint {
	return pr.numberOfRequestChanges
}

func (pr *PullRequest) Author() Author {
	return pr.author
}

func (pr *PullRequest) GetOpenedDays() uint {
	return pr.openedDays
}

func (pr *PullRequest) Priority() Priority {
	return pr.priority
}

func (pr Priority) ToString() string {
	return string(pr)
}

func (pr *PullRequest) Repository() Repository {
	return pr.repository
}
