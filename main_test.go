package main

import (
	"regexp"
	"testing"
)

const (
	sampleCodeWithoutComments = `versionCode 1`
	sampleNameWithoutComments = `versionName "1.0"`

	sampleCodeWithComments = `versionCode 1//close comment`
	sampleNameWithComments = `versionName "1.0" // far comment`
)

func Test_regexPatterns(t *testing.T) {
	tests := []struct {
		name          string
		sampleContent string
		regexPattern  string
		want          string
	}{
		{"versionCode check without comments", sampleCodeWithoutComments, versionCodeRegexPattern, "1"},
		{"versionName check without comments", sampleNameWithoutComments, versionNameRegexPattern, "1.0"},
		{"versionCode check with comments", sampleCodeWithComments, versionCodeRegexPattern, "1"},
		{"versionName check with comments", sampleNameWithComments, versionNameRegexPattern, "1.0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := regexp.MustCompile(tt.regexPattern).FindStringSubmatch(tt.sampleContent)
			if len(got) == 0 {
				t.Errorf("regex(%s) didn't match for content: %s\n\n got: %s", tt.regexPattern, tt.sampleContent, got)
				return
			}
			if got[1] != tt.want {
				t.Errorf("got: (%v), want: (%v)", got[1], tt.want)
			}
		})
	}
}
