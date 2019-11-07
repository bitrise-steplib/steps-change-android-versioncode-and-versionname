package main

import (
	"regexp"
	"testing"
)

// const (
// 	sampleCodeWithoutComments =
// 	sampleNameWithoutComments =

// 	sampleCodeWithComments =
// 	sampleNameWithComments =

// 	sampleCodeWithCommentsAndVar = `versionCode myVar//close comment`
// 	sampleNameWithCommentsAndVar = `versionName myVar // far comment`
// )

func Test_regexPatterns(t *testing.T) {
	for _, tt := range []struct {
		sampleContent string
		want          string
		regexPattern  string
	}{
		{`versionCode 1`, "1", versionCodeRegexPattern},
		{`versionCode 1//close comment`, "1", versionCodeRegexPattern},
		{`versionCode 1 // far comment`, "1", versionCodeRegexPattern},
		{`versionCode myWar`, "myWar", versionCodeRegexPattern},
		{`versionCode myWar//close comment`, "myWar", versionCodeRegexPattern},
		{`versionCode myWar // far comment`, "myWar", versionCodeRegexPattern},

		{`versionCode = 1`, "1", versionCodeRegexPattern},
		{`versionCode =1//close comment`, "1", versionCodeRegexPattern},
		{`versionCode= 1 // far comment`, "1", versionCodeRegexPattern},
		{`versionCode = myWar`, "myWar", versionCodeRegexPattern},
		{`versionCode   =  myWar//close comment`, "myWar", versionCodeRegexPattern},
		{`versionCode  = myWar // far comment`, "myWar", versionCodeRegexPattern},

		{`versionCode 1` + "\n", "1", versionCodeRegexPattern},
		{`versionCode 1//close comment` + "\n", "1", versionCodeRegexPattern},
		{`versionCode 1 // far comment` + "\n", "1", versionCodeRegexPattern},
		{`versionCode myWar` + "\n", "myWar", versionCodeRegexPattern},
		{`versionCode myWar//close comment` + "\n", "myWar", versionCodeRegexPattern},
		{`versionCode myWar // far comment` + "\n", "myWar", versionCodeRegexPattern},

		{`versionCode = 1` + "\n", "1", versionCodeRegexPattern},
		{`versionCode =1//close comment` + "\n", "1", versionCodeRegexPattern},
		{`versionCode= 1 // far comment` + "\n", "1", versionCodeRegexPattern},
		{`versionCode = myWar` + "\n", "myWar", versionCodeRegexPattern},
		{`versionCode   =  myWar//close comment` + "\n", "myWar", versionCodeRegexPattern},
		{`versionCode  = myWar // far comment` + "\n", "myWar", versionCodeRegexPattern},

		{`versionName "1.0"`, "1.0", versionNameRegexPattern},
		{`versionName "1.0"//close comment`, "1.0", versionNameRegexPattern},
		{`versionName "1.0" // far comment`, "1.0", versionNameRegexPattern},
		{`versionName '1.0'`, "1.0", versionNameRegexPattern},
		{`versionName '1.0'//close comment`, "1.0", versionNameRegexPattern},
		{`versionName '1.0' // far comment`, "1.0", versionNameRegexPattern},
		{`versionName = '1.0' // far comment`, "1.0", versionNameRegexPattern},

		{`versionName myWar`, "myWar", versionNameRegexPattern},
		{`versionName myWar//close comment`, "myWar", versionNameRegexPattern},
		{`versionName myWar // far comment`, "myWar", versionNameRegexPattern},
		{`versionName = myWar // far comment`, "myWar", versionNameRegexPattern},

		{`versionName "1.0"` + "\n", "1.0", versionNameRegexPattern},
		{`versionName "1.0"//close comment` + "\n", "1.0", versionNameRegexPattern},
		{`versionName="1.0" // far comment` + "\n", "1.0", versionNameRegexPattern},
		{`versionName '1.0'` + "\n", "1.0", versionNameRegexPattern},
		{`versionName '1.0'//close comment` + "\n", "1.0", versionNameRegexPattern},
		{`versionName '1.0' // far comment` + "\n", "1.0", versionNameRegexPattern},
		{`versionName = '1.0' // far comment` + "\n", "1.0", versionNameRegexPattern},

		{`versionName myWar` + "\n", "myWar", versionNameRegexPattern},
		{`versionName myWar//close comment` + "\n", "myWar", versionNameRegexPattern},
		{`versionName myWar // far comment` + "\n", "myWar", versionNameRegexPattern},
		{`versionName = myWar // far comment` + "\n", "myWar", versionNameRegexPattern},

		{`versionName myWar` + "\n", "myWar", versionNameRegexPattern},
		{`versionName myWar//close comment` + "\n", "myWar", versionNameRegexPattern},
		{`versionName myWar // far comment` + "\n", "myWar", versionNameRegexPattern},
		{`versionName=myWar // far comment` + "\n", "myWar", versionNameRegexPattern},
	} {
		got := regexp.MustCompile(tt.regexPattern).FindStringSubmatch(tt.sampleContent)
		if len(got) == 0 {
			t.Errorf("regex(%s) didn't match for content: %s\n\n got: %s", tt.regexPattern, tt.sampleContent, got)
			return
		}
		if got[1] != tt.want {
			t.Errorf("got: (%v), want: (%v)", got[1], tt.want)
		}
	}
}
