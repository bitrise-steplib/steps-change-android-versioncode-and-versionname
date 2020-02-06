package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
)

const (
	// versionCode — A positive integer [...] -> https://developer.android.com/studio/publish/versioning
	versionCodeRegexPattern = `^versionCode(?:\s|=)+(.*?)(?:\s|\/\/|$)`
	// versionName — A string used as the version number shown to users [...] -> https://developer.android.com/studio/publish/versioning
	versionNameRegexPattern = `^versionName(?:=|'|"|\s)+(.*?)(?:'|"|\s|\/\/|$)`
)

type config struct {
	BuildGradlePth    string  `env:"build_gradle_path,file"`
	NewVersionName    *string `env:"new_version_name"`
	NewVersionCode    *int    `env:"new_version_code"`
	VersionCodeOffset int     `env:"version_code_offset"`
}

type updateFn func(line string, lineNum int, matches []string) string

func findAndUpdate(reader io.Reader, update map[*regexp.Regexp]updateFn) (string, error) {
	scanner := bufio.NewScanner(reader)
	var updatedLines []string

	for lineNum := 0; scanner.Scan(); lineNum++ {
		line := scanner.Text()

		updated := false
		for re, fn := range update {
			if match := re.FindStringSubmatch(strings.TrimSpace(line)); len(match) == 2 {
				if updatedLine := fn(line, lineNum, match); updatedLine != "" {
					updatedLines = append(updatedLines, updatedLine)
					updated = true
					break
				}
			}
		}
		if !updated {
			updatedLines = append(updatedLines, line)
		}
	}

	return strings.Join(updatedLines, "\n"), scanner.Err()
}

func exportOutputs(outputs map[string]string) error {
	for envKey, envValue := range outputs {
		cmd := command.New("envman", "add", "--key", envKey, "--value", envValue)
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}

func main() {
	var cfg config
	if err := stepconf.Parse(&cfg); err != nil {
		failf("Issue with input: %s", err)
	}
	stepconf.Print(cfg)
	fmt.Println()

	if cfg.NewVersionName == nil && cfg.NewVersionCode == nil {
		failf("Neither NewVersionCode nor NewVersionName are provided, however one of them is required.")
	}

	//
	// find versionName & versionCode with regexp
	fmt.Println()
	log.Infof("Updating versionName and versionCode in: %s", cfg.BuildGradlePth)

	f, err := os.Open(cfg.BuildGradlePth)
	if err != nil {
		failf("Failed to read build.gradle file, error: %s", err)
	}

	var finalVersionCode, finalVersionName string
	var updatedVersionCodes, updatedVersionNames int

	updatedBuildGradleContent, err := findAndUpdate(f, map[*regexp.Regexp]updateFn{
		regexp.MustCompile(versionCodeRegexPattern): func(line string, lineNum int, match []string) string {
			oldVersionCode := match[1]
			finalVersionCode = oldVersionCode
			updatedLine := ""

			if cfg.NewVersionCode != nil {
				finalVersionCode = strconv.Itoa(*cfg.NewVersionCode + cfg.VersionCodeOffset)
				updatedLine = strings.Replace(line, oldVersionCode, finalVersionCode, -1)
				updatedVersionCodes++
				log.Printf("updating line (%d): %s -> %s", lineNum, line, updatedLine)
			}

			return updatedLine
		},
		regexp.MustCompile(versionNameRegexPattern): func(line string, lineNum int, match []string) string {
			oldVersionName := match[1]
			finalVersionName = oldVersionName
			updatedLine := ""

			if cfg.NewVersionName != nil {
				finalVersionName = *cfg.NewVersionName
				updatedLine = strings.Replace(line, oldVersionName, finalVersionName, -1)
				updatedVersionNames++
				log.Printf("updating line (%d): %s -> %s", lineNum, line, updatedLine)
			}

			return updatedLine
		},
	})
	if err != nil {
		failf("Failed to scann build.gradle file, error: %s", err)
	}

	//
	// export outputs
	if err := exportOutputs(map[string]string{
		"ANDROID_VERSION_NAME": finalVersionName,
		"ANDROID_VERSION_CODE": finalVersionCode,
	}); err != nil {
		failf("Failed to export outputs, error: %s", err)
	}

	if err := fileutil.WriteStringToFile(cfg.BuildGradlePth, updatedBuildGradleContent); err != nil {
		failf("Failed to write build.gradle file, error: %s", err)
	}

	fmt.Println()
	log.Donef("%d versionCode updated", updatedVersionCodes)
	log.Donef("%d versionName updated", updatedVersionNames)
}
