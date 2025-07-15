package forum

import (
	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/validator.v2"
	"strings"
	"time"
)

type Reply struct {
	Id           int
	ForumEntryId int       `json:"forum_entry_id" form:"forum_entry_id" validate:"nonzero"`
	Content      string    `json:"content" form:"content" validate:"nonzero,max=800"`
	Created      time.Time `json:"created"`
}

func (reply *Reply) Process() *Reply {
	once.Do(func() {
		policy = bluemonday.StrictPolicy()
		policy.AllowElements("b")
		policy.AllowElements("i")
	})

	reply.Content = policy.Sanitize(
		strings.TrimSpace(reply.Content),
	)

	reply.Created = time.Now()

	return reply
}

func (reply *Reply) Validate() error {
	return validator.Validate(reply)
}