package jira

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestLinkCreation(t *testing.T) {
	jiraIntegration := &Integration{}
	link := jiraIntegration.CreateLink("khulnasoft-lab", "https://khulnasoft-lab.io")
	assert.Equal(t, "[khulnasoft-lab|https://khulnasoft-lab.io]", link)
}

func TestHorizontalLineCreation(t *testing.T) {
	jiraIntegration := &Integration{}
	horizontalLine := jiraIntegration.CreateHorizontalLine()
	assert.True(t, strings.Contains(horizontalLine, "----"))
}

func TestTableCreation(t *testing.T) {
	jiraIntegration := &Integration{}

	table, err := jiraIntegration.CreateTable([]string{"key", "value"}, [][]string{
		{"a", "b"},
		{"c"},
		{"d", "e"},
	})

	assert.Nil(t, err)
	expected := `| key | value |
| a | b |
| c |  |
| d | e |
`
	assert.Equal(t, expected, table)
}
