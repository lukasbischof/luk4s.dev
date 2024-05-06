package forum

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcessEntryWithMaliciousContent(t *testing.T) {
	entry := &Entry{
		Content: "<b>Test Content</b><script>malicious code</script>",
		Author:  "<i>Test Author</i><script>malicious code</script>",
	}

	processedEntry := entry.Process()

	assert.NotEmpty(t, processedEntry.Id)
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

	assert.NotEmpty(t, processedEntry.Id)
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

func TestToJsonReturnsValidJsonForValidEntry(t *testing.T) {
	entry := &Entry{
		Content: "Test Content",
		Author:  "Test Author",
	}

	json, err := entry.ToJson()

	assert.Nil(t, err)
	assert.NotEmpty(t, json)
	assert.Contains(t, json, "Test Content")
	assert.Contains(t, json, "Test Author")
}
