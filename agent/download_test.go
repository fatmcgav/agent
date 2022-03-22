//go:build !windows
// +build !windows

package agent

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetTargetPath(t *testing.T) {
	// If the destination ends in /,
	// then no part of destination is stripped
	// before joining with path
	assert.Equal(t, "a/a", getTargetPath("a", "a/"))
	assert.Equal(t, "c/a/a/b", getTargetPath("a/b", "c/a/"))
	// if the last part of destination
	// does not match the first part of path,
	// then just join path to destination
	assert.Equal(t, "a", getTargetPath("a", "."))
	assert.Equal(t, "a/b/c", getTargetPath("a/b/c", "."))
	assert.Equal(t, "a/b/a/b", getTargetPath("a/b", "a/b"))
	// If the last part of destination
	// matches the first part of path,
	// then remove the last part of destination
	// before joining with path
	assert.Equal(t, "a", getTargetPath("a", "a"))
	assert.Equal(t, "c/a/b", getTargetPath("a/b", "c/a"))
	assert.Equal(t, "lambda.zip", getTargetPath("lambda.zip", "lambda.zip"))
	// when the path[0], path[100] matches
	// destination[-1], destination[-2]
	// then the last 2 characters of destination
	// are removed before joining path
	// This is a tricky one to use.
	assert.Equal(t,
		"a/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/b",
		getTargetPath(
			"a/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/x/b",
			"b/a"),
	)
	// Gotcha: this is not what you want.
	assert.Equal(t, "a/lambda.zip/a/lambda.zip", getTargetPath("a/lambda.zip", "a/lambda.zip"))

	// Test absolute paths
	assert.Equal(t, "/var/logs/app/a.log", getTargetPath("app/a.log", "/var/logs"))          // no match, no trailing
	assert.Equal(t, "/var/logs/app/a.log", getTargetPath("app/a.log", "/var/logs/app"))      // match, no trailing
	assert.Equal(t, "/var/logs/app/app/a.log", getTargetPath("app/a.log", "/var/logs/app/")) // match, trailing

	// artifact_download documentation examples
	assert.Equal(t, "foo/app/app/logs/a.log", getTargetPath("app/logs/a.log", "foo/app/"))
	assert.Equal(t, "foo/app/logs/a.log", getTargetPath("app/logs/a.log", "foo/app"))
	assert.Equal(t, "app/logs/a.log", getTargetPath("app/logs/a.log", "."))
}
