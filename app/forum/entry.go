package forum

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2/utils"
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
	Id      string
	Content string    `json:"content" validate:"nonzero"`
	Author  string    `json:"author" validate:"nonzero"`
	Created time.Time `json:"created"`
}

func FromJson(data []byte) (*Entry, error) {
	var entry Entry
	err := json.Unmarshal(data, &entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
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

	entry.Author = policy.Sanitize(
		strings.TrimSpace(entry.Author),
	)

	entry.Id = utils.UUID()
	entry.Created = time.Now()

	return entry
}

func (entry *Entry) Validate() error {
	return validator.Validate(entry)
}

func (entry *Entry) ToJson() (string, error) {
	marshal, err := json.Marshal(entry)
	if err != nil {
		return "", err
	}

	return string(marshal), nil
}
