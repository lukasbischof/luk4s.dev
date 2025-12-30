package forum

import (
	"github.com/microcosm-cc/bluemonday"
	"gopkg.in/validator.v2"
	"strings"
	"sync"
	"time"
)

var (
	once   sync.Once
	policy *bluemonday.Policy
)

type Entry struct {
	Id              int
	Content         string    `json:"content" form:"content" validate:"nonzero,max=800"`
	Author          string    `json:"author" form:"author" validate:"nonzero,max=100"`
	CaptchaResponse string    `form:"h-captcha-response" json:"-" validate:"nonzero"`
	Created         time.Time `json:"created"`
}

func (entry *Entry) Process() *Entry {
	once.Do(func() {
		policy = bluemonday.StrictPolicy()
		policy.AllowElements("b")
		policy.AllowElements("i")
	})

	entry.Content = policy.Sanitize(
		strings.TrimSpace(entry.Content),
	)

	// Strip all HTML tags from Author field
	strictPolicy := bluemonday.StrictPolicy()
	entry.Author = strictPolicy.Sanitize(
		strings.TrimSpace(entry.Author),
	)

	entry.Created = time.Now()

	return entry
}

func (entry *Entry) Validate() error {
	return validator.Validate(entry)
}
