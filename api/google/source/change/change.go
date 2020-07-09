package change

import "fmt"

var changeURL = "https://go-review.googlesource.com/c/go/+/%d"

type List []Change

type Change struct {
	Project           string  `json:"project"`
	Branch            string  `json:"branch"`
	ChangeID          string  `json:"change_id"`
	Subject           string  `json:"subject"`
	Status            string  `json:"status"`
	Created           string  `json:"created"`
	Updated           string  `json:"updated"`
	Submitted         string  `json:"submitted"`
	Insertions        int     `json:"insertions"`
	Deletions         int     `json:"deletions"`
	TotalCommentCount int     `json:"total_comment_count"`
	HasReviewStarted  bool    `json:"has_review_started"`
	Number            int     `json:"_number"`
	SubmissionID      string  `json:"submission_id"`
	Owner             Person  `json:"owner"`
	Submitter         *Person `json:"submitter"`
	Labels            Labels  `json:"labels"`
}

type Operate struct {
	Approved *Person `json:"approved"`
	Optional bool    `json:"optional"`
}

type Person struct {
	AccountID int    `json:"_account_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
}

func (c Change) ChangeURL() string {
	return fmt.Sprintf(changeURL, c.Number)
}

type Labels map[string]Operate

func (l Labels) CodeReview() *Person {
	return l.personForLabel("Code-Review")
}

func (l Labels) RunTryBot() *Person {
	return l.personForLabel("Run-TryBot")
}

func (l Labels) TryBotResult() *Person {
	return l.personForLabel("TryBot-Result")
}

func (l Labels) personForLabel(label string) *Person {
	operate, ok := l[label]
	if !ok {
		return nil
	}

	return operate.Approved
}
