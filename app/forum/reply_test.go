package forum

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestProcessReplyWithMaliciousContent(t *testing.T) {
	reply := &Reply{
		ForumEntryId: 1,
		Content:      "<b>Test Reply Content</b><script>malicious code</script>",
	}

	processedReply := reply.Process()

	assert.Empty(t, processedReply.Id)
	assert.NotEmpty(t, processedReply.Created)
	assert.Equal(t, 1, processedReply.ForumEntryId)
	assert.Equal(t, "<b>Test Reply Content</b>", processedReply.Content)
}

func TestProcessReplyWithEmptyContent(t *testing.T) {
	reply := &Reply{
		ForumEntryId: 1,
		Content:      "",
	}

	processedReply := reply.Process()

	assert.Empty(t, processedReply.Id)
	assert.NotEmpty(t, processedReply.Created)
	assert.Equal(t, 1, processedReply.ForumEntryId)
	assert.Equal(t, "", processedReply.Content)
}

func TestProcessReplyWithPaddedContent(t *testing.T) {
	reply := &Reply{
		ForumEntryId: 1,
		Content:      "   Test Reply Content   ",
	}

	processedReply := reply.Process()

	assert.Equal(t, 1, processedReply.ForumEntryId)
	assert.Equal(t, "Test Reply Content", processedReply.Content)
}

func TestProcessReplyWithAllowedHtmlTags(t *testing.T) {
	reply := &Reply{
		ForumEntryId: 1,
		Content:      "<b>Bold content</b> and <i>italic content</i>",
	}

	processedReply := reply.Process()

	assert.Equal(t, "<b>Bold content</b> and <i>italic content</i>", processedReply.Content)
}

func TestProcessReplyWithDisallowedHtmlTags(t *testing.T) {
	reply := &Reply{
		ForumEntryId: 1,
		Content:      "<b>Bold</b><div>Not allowed</div><i>Italic</i>",
	}

	processedReply := reply.Process()

	assert.Equal(t, "<b>Bold</b>Not allowed<i>Italic</i>", processedReply.Content)
}

func TestValidateReplyReturnsErrorForEmptyContent(t *testing.T) {
	reply := &Reply{
		ForumEntryId: 1,
		Content:      "",
	}

	err := reply.Validate()

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "Content: zero value")
}

func TestValidateReturnsErrorForZeroForumEntryId(t *testing.T) {
	reply := &Reply{
		ForumEntryId: 0,
		Content:      "Test Reply Content",
	}

	err := reply.Validate()

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "ForumEntryId: zero value")
}

func TestValidateReplyReturnsErrorForContentExceedingMaxLength(t *testing.T) {
	reply := &Reply{
		ForumEntryId: 1,
		Content:      strings.Repeat("a", 801),
	}

	err := reply.Validate()

	assert.NotNil(t, err)
}

func TestValidateSucceedsForValidReply(t *testing.T) {
	reply := &Reply{
		ForumEntryId: 1,
		Content:      "Valid reply content",
	}

	err := reply.Validate()

	assert.Nil(t, err)
}

func TestValidateSucceedsForContentAtMaxLength(t *testing.T) {
	reply := &Reply{
		ForumEntryId: 1,
		Content:      strings.Repeat("a", 800),
	}

	err := reply.Validate()

	assert.Nil(t, err)
}