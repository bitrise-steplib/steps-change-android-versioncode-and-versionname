package main

import (
	"io"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func Test_typeConv(t *testing.T) {
	if string(rune(323)) != "Ń" {
		t.Fatal("invalid type conversion")
	}
	if strconv.Itoa(323) != "323" {
		t.Fatal("invalid type conversion")
	}
}

func Test_regexPatterns(t *testing.T) {
	for _, tt := range []struct {
		sampleContent string
		want          string
		regexPattern  string
	}{
		// versionCode regex
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

		// versionName regex
		{`versionName "1.0"`, `"1.0"`, versionNameRegexPattern},
		{`versionName "1.0"//close comment`, `"1.0"`, versionNameRegexPattern},
		{`versionName "1.0" // far comment`, `"1.0"`, versionNameRegexPattern},
		{`versionName '1.0'`, `'1.0'`, versionNameRegexPattern},
		{`versionName '1.0'//close comment`, `'1.0'`, versionNameRegexPattern},
		{`versionName '1.0' // far comment`, `'1.0'`, versionNameRegexPattern},
		{`versionName = '1.0' // far comment`, `'1.0'`, versionNameRegexPattern},

		{`versionName myWar`, "myWar", versionNameRegexPattern},
		{`versionName myWar//close comment`, "myWar", versionNameRegexPattern},
		{`versionName myWar // far comment`, "myWar", versionNameRegexPattern},
		{`versionName = myWar // far comment`, "myWar", versionNameRegexPattern},

		{`versionName "1.0"` + "\n", `"1.0"`, versionNameRegexPattern},
		{`versionName "1.0"//close comment` + "\n", `"1.0"`, versionNameRegexPattern},
		{`versionName="1.0" // far comment` + "\n", `"1.0"`, versionNameRegexPattern},
		{`versionName '1.0'` + "\n", `'1.0'`, versionNameRegexPattern},
		{`versionName '1.0'//close comment` + "\n", `'1.0'`, versionNameRegexPattern},
		{`versionName '1.0' // far comment` + "\n", `'1.0'`, versionNameRegexPattern},
		{`versionName = '1.0' // far comment` + "\n", `'1.0'`, versionNameRegexPattern},

		{`versionName myWar` + "\n", "myWar", versionNameRegexPattern},
		{`versionName myWar//close comment` + "\n", "myWar", versionNameRegexPattern},
		{`versionName myWar // far comment` + "\n", "myWar", versionNameRegexPattern},
		{`versionName = myWar // far comment` + "\n", "myWar", versionNameRegexPattern},

		{`versionName myWar` + "\n", "myWar", versionNameRegexPattern},
		{`versionName myWar//close comment` + "\n", "myWar", versionNameRegexPattern},
		{`versionName myWar // far comment` + "\n", "myWar", versionNameRegexPattern},
		{`versionName=myWar // far comment` + "\n", "myWar", versionNameRegexPattern},
	} {
		t.Run(tt.sampleContent, func(t *testing.T) {
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

func TestBuildGradleVersionUpdater_UpdateVersion(t *testing.T) {
	tests := []struct {
		name              string
		buildGradleReader io.Reader
		newVersionCode    int
		versionCodeOffset int
		newVersionName    string

		want    UpdateResult
		wantErr bool
	}{
		// versionCode update
		{
			name:              "Updates versionCode value",
			buildGradleReader: strings.NewReader("versionCode 1"),
			newVersionCode:    555,
			want:              UpdateResult{NewContent: "versionCode 555", FinalVersionCode: "555", UpdatedVersionCodes: 1},
		},
		{
			name:              "Updates versionCode variable",
			buildGradleReader: strings.NewReader("versionCode rootProject.ext.versionCode"),
			newVersionCode:    555,
			want:              UpdateResult{NewContent: "versionCode 555", FinalVersionCode: "555", UpdatedVersionCodes: 1},
		},
		{
			name:              "versionCode needs to be a positive integer",
			buildGradleReader: strings.NewReader("versionCode 1"),
			newVersionCode:    0,
			want:              UpdateResult{NewContent: "versionCode 1", FinalVersionCode: "1"},
		},
		{
			name:              "Does not touch ABI version code mapping",
			buildGradleReader: strings.NewReader(`def versionCodes = ["armeabi-v7a": 1, "x86": 2, "arm64-v8a": 3, "x86_64": 4]`),
			newVersionCode:    555,
			want:              UpdateResult{NewContent: `def versionCodes = ["armeabi-v7a": 1, "x86": 2, "arm64-v8a": 3, "x86_64": 4]`},
		},
		{
			name:              "Does not touch per ABI version selector",
			buildGradleReader: strings.NewReader(`versionCodes.get(abi) * 1000 + defaultConfig.versionCode`),
			newVersionCode:    555,
			want:              UpdateResult{NewContent: `versionCodes.get(abi) * 1000 + defaultConfig.versionCode`},
		},
		// versionName update
		{
			name:              "Updates versionName value with single quote",
			buildGradleReader: strings.NewReader(`versionName "0.9.0"`),
			newVersionName:    `"1.1.0"`,
			want:              UpdateResult{NewContent: `versionName "1.1.0"`, FinalVersionName: `"1.1.0"`, UpdatedVersionNames: 1},
		},
		{
			name:              "Updates versionName value with double quote",
			buildGradleReader: strings.NewReader(`versionName '0.9.0'`),
			newVersionName:    `"1.1.0"`,
			want:              UpdateResult{NewContent: `versionName "1.1.0"`, FinalVersionName: `"1.1.0"`, UpdatedVersionNames: 1},
		},
		{
			name:              "Updates versionName variable",
			buildGradleReader: strings.NewReader("versionName rootProject.ext.versionName"),
			newVersionName:    `"1.1.0"`,
			want:              UpdateResult{NewContent: `versionName "1.1.0"`, FinalVersionName: `"1.1.0"`, UpdatedVersionNames: 1},
		},
		{
			name:              "versionName needs to be a not empty string",
			buildGradleReader: strings.NewReader(`versionName "1.0.0"`),
			newVersionName:    "",
			want:              UpdateResult{NewContent: `versionName "1.0.0"`, FinalVersionName: `"1.0.0"`},
		},
		{
			name:              "Adds quotation mark to newVersionName if missing",
			buildGradleReader: strings.NewReader("versionName rootProject.ext.versionName"),
			newVersionName:    `1.1.0`,
			want:              UpdateResult{NewContent: `versionName "1.1.0"`, FinalVersionName: `"1.1.0"`, UpdatedVersionNames: 1},
		},
		{
			name:              "Adds quotation mark to newVersionName if leading is missing",
			buildGradleReader: strings.NewReader("versionName rootProject.ext.versionName"),
			newVersionName:    `1.1.0"`,
			want:              UpdateResult{NewContent: `versionName "1.1.0"`, FinalVersionName: `"1.1.0"`, UpdatedVersionNames: 1},
		},
		{
			name:              "Adds quotation mark to newVersionName if traling is missing",
			buildGradleReader: strings.NewReader("versionName rootProject.ext.versionName"),
			newVersionName:    `"1.1.0`,
			want:              UpdateResult{NewContent: `versionName "1.1.0"`, FinalVersionName: `"1.1.0"`, UpdatedVersionNames: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewBuildGradleVersionUpdater(tt.buildGradleReader)
			got, err := u.UpdateVersion(tt.newVersionCode, tt.versionCodeOffset, tt.newVersionName)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildGradleVersionUpdater.UpdateVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildGradleVersionUpdater.UpdateVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
