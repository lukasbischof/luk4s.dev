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
	assert.Equal(t, "Test Author", processedEntry.Author)
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

func TestValidateReturnsErrorForAuthorExceedingMaxLength(t *testing.T) {
	entry := &Entry{
		Content:         "Test Content",
		Author:          strings.Repeat("a", 101),
		CaptchaResponse: "valid-captcha",
	}

	err := entry.Validate()

	assert.NotNil(t, err)
}

func TestValidateSucceedsForAuthorAtMaxLength(t *testing.T) {
	entry := &Entry{
		Content:         "Test Content",
		Author:          strings.Repeat("a", 100),
		CaptchaResponse: "valid-captcha",
	}

	err := entry.Validate()

	assert.Nil(t, err)
}

func TestProcessEntryStripsAllHtmlTagsFromAuthor(t *testing.T) {
	entry := &Entry{
		Content: "<b>Test Content</b>",
		Author:  "<b>Bold</b><i>Italic</i><div>Div</div><span>Span</span>Plain",
	}

	processedEntry := entry.Process()

	assert.Equal(t, "<b>Test Content</b>", processedEntry.Content)
	assert.Equal(t, "BoldItalicDivSpanPlain", processedEntry.Author)
}

func TestProcessEntryAllowsHtmlTagsInContent(t *testing.T) {
	entry := &Entry{
		Content: "<b>Bold</b> and <i>italic</i> content",
		Author:  "Test Author",
	}

	processedEntry := entry.Process()

	assert.Equal(t, "<b>Bold</b> and <i>italic</i> content", processedEntry.Content)
	assert.Equal(t, "Test Author", processedEntry.Author)
}

func TestProcessEntryStripsComplexHtmlFromAuthor(t *testing.T) {
	entry := &Entry{
		Content: "Test Content",
		Author:  "<a href='evil.com'>Link</a><p>Paragraph</p><h1>Header</h1>",
	}

	processedEntry := entry.Process()

	assert.Equal(t, "Test Content", processedEntry.Content)
	assert.Equal(t, "LinkParagraphHeader", processedEntry.Author)
}

func TestProcessEntryStripsNestedHtmlFromAuthor(t *testing.T) {
	entry := &Entry{
		Content: "Test Content",
		Author:  "<div><span><b>Nested</b></span></div>",
	}

	processedEntry := entry.Process()

	assert.Equal(t, "Test Content", processedEntry.Content)
	assert.Equal(t, "Nested", processedEntry.Author)
}

func TestValidateSucceedsForValidEntry(t *testing.T) {
	entry := &Entry{
		Content:         "Valid content",
		Author:          "Valid Author",
		CaptchaResponse: "valid-captcha",
	}

	err := entry.Validate()

	assert.Nil(t, err)
}

func TestProcessEntryWithAuthorAtBoundary(t *testing.T) {
	entry := &Entry{
		Content: "Test Content",
		Author:  strings.Repeat("a", 99),
	}

	processedEntry := entry.Process()

	assert.Equal(t, strings.Repeat("a", 99), processedEntry.Author)
	assert.Equal(t, "Test Content", processedEntry.Content)
}

func TestProcessEntryPreservesPlainTextAuthor(t *testing.T) {
	entry := &Entry{
		Content: "Test Content",
		Author:  "John Doe",
	}

	processedEntry := entry.Process()

	assert.Equal(t, "John Doe", processedEntry.Author)
}

func TestProcessEntryWithSpecialCharactersInAuthor(t *testing.T) {
	entry := &Entry{
		Content: "Test Content",
		Author:  "Author & Co. <author@example.com>",
	}

	processedEntry := entry.Process()

	// Bluemonday escapes special HTML characters
	assert.Equal(t, "Author &amp; Co. ", processedEntry.Author)
}
