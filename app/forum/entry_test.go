package forum

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestProcessEntryWithMaliciousContent(t *testing.T) {
	entry := &Entry{
		Content: "<b>Test Content</b><script>malicious code</script>",
		Author:  "<i>Test Author</i><script>malicious code</script>",
	}

	processedEntry := entry.Process()

	assert.Empty(t, processedEntry.Id)
	assert.NotEmpty(t, processedEntry.Created)
	assert.Equal(t, "<b>Test Content</b>", processedEntry.Content)
	assert.Equal(t, "<i>Test Author</i>", processedEntry.Author)
}

func TestProcessEntryWithEmptyContent(t *testing.T) {
	entry := &Entry{
		Content: "",
		Author:  "Test Author",
	}

	processedEntry := entry.Process()

	assert.Empty(t, processedEntry.Id)
	assert.NotEmpty(t, processedEntry.Created)
	assert.Equal(t, "", processedEntry.Content)
	assert.Equal(t, "Test Author", processedEntry.Author)
}

func TestProcessEntryWithPaddedContent(t *testing.T) {
	entry := &Entry{
		Content: "   Test Content   ",
		Author:  "   Test Author   ",
	}

	processedEntry := entry.Process()

	assert.Equal(t, "Test Content", processedEntry.Content)
	assert.Equal(t, "Test Author", processedEntry.Author)
}

func TestValidateReturnsErrorForEmptyContent(t *testing.T) {
	entry := &Entry{
		Content:         "",
		Author:          "Test Author",
		CaptchaResponse: "valid-captcha",
	}

	err := entry.Validate()

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "Content: zero value")
}

func TestValidateReturnsErrorForEmptyAuthor(t *testing.T) {
	entry := &Entry{
		Content:         "Test Content",
		Author:          "",
		CaptchaResponse: "valid-captcha",
	}

	err := entry.Validate()

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "Author: zero value")
}

func TestValidateReturnsErrorForContentExceedingMaxLength(t *testing.T) {
	entry := &Entry{
		Content:         strings.Repeat("a", 801),
		Author:          "Test Author",
		CaptchaResponse: "valid-captcha",
	}

	err := entry.Validate()

	assert.NotNil(t, err)
}
