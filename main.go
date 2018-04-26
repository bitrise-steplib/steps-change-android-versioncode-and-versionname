package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/input"
)

// ConfigsModel ...
type ConfigsModel struct {
	BuildGradlePth    string
	NewVersionName    string
	NewVersionCode    string
	VersionCodeOffset string
}

func createConfigsModelFromEnvs() ConfigsModel {
	return ConfigsModel{
		BuildGradlePth:    os.Getenv("build_gradle_path"),
		NewVersionName:    os.Getenv("new_version_name"),
		NewVersionCode:    os.Getenv("new_version_code"),
		VersionCodeOffset: os.Getenv("version_code_offset"),
	}
}

func (configs ConfigsModel) print() {
	log.Infof("Configs:")
	log.Printf("- BuildGradlePth: %s", configs.BuildGradlePth)
	log.Printf("- NewVersionName: %s", configs.NewVersionName)
	log.Printf("- NewVersionCode: %s", configs.NewVersionCode)
	log.Printf("- VersionCodeOffset: %s", configs.VersionCodeOffset)
}

func (configs ConfigsModel) validate() error {
	if err := input.ValidateIfPathExists(configs.BuildGradlePth); err != nil {
		return errors.New("issue with input BuildGradlePth: " + err.Error())
	}

	if err := input.ValidateIfNotEmpty(configs.NewVersionCode); err != nil {
		if err := input.ValidateIfNotEmpty(configs.NewVersionName); err != nil {
			return errors.New("neither NewVersionCode nor NewVersionName are provided however, one of them is required")
		}
	}

	return nil
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

func logFail(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}

func main() {
	//
	// validate and prepare inputs
	configs := createConfigsModelFromEnvs()

	fmt.Println()
	configs.print()

	if err := configs.validate(); err != nil {
		logFail("Issue with input: %s", err)
	}

	var newVersionCode int
	if configs.NewVersionCode != "" {
		var err error
		newVersionCode, err = strconv.Atoi(configs.NewVersionCode)
		if err != nil {
			logFail("Failed to convert to string: %s, error: %s", configs.NewVersionCode, err)
		}
	}

	if configs.VersionCodeOffset != "" {
		offset, err := strconv.Atoi(configs.VersionCodeOffset)
		if err != nil {
			logFail("Failed to convert to string: %s, error: %s", configs.VersionCodeOffset, err)
		}
		newVersionCode += offset
	}

	//
	// find versionName & versionCode with regexp
	fmt.Println()
	log.Infof("Updating versionName and versionCode in: %s", configs.BuildGradlePth)

	f, err := os.Open(configs.BuildGradlePth)
	if err != nil {
		logFail("Failed to read build.gradle file, error: %s", err)
	}

	var finalVersionCode, finalVersionName string
	var updatedVersionCodes, updatedVersionNames int

	updatedBuildGradleContent, err := findAndUpdate(f, map[*regexp.Regexp]updateFn{
		regexp.MustCompile(`^versionCode (?P<version_code>.*)`): func(line string, lineNum int, match []string) string {
			oldVersionCode := match[1]
			finalVersionCode = oldVersionCode
			updatedLine := ""

			if configs.NewVersionCode != "" {
				finalVersionCode = strconv.Itoa(newVersionCode)
				updatedLine = strings.Replace(line, oldVersionCode, finalVersionCode, -1)
				updatedVersionCodes++
				log.Printf("updating line (%d): %s -> %s", lineNum, line, updatedLine)
			}

			return updatedLine
		},
		regexp.MustCompile(`^versionName "(?P<version_code>.*)"`): func(line string, lineNum int, match []string) string {
			oldVersionName := match[1]
			finalVersionName = oldVersionName
			updatedLine := ""

			if configs.NewVersionName != "" {
				finalVersionName = configs.NewVersionName
				updatedLine = strings.Replace(line, oldVersionName, finalVersionName, -1)
				updatedVersionNames++
				log.Printf("updating line (%d): %s -> %s", lineNum, line, updatedLine)
			}

			return updatedLine
		},
	})
	if err != nil {
		logFail("Failed to scann build.gradle file, error: %s", err)
	}

	//
	// export outputs
	if err := exportOutputs(map[string]string{
		"ANDROID_VERSION_NAME": finalVersionName,
		"ANDROID_VERSION_CODE": finalVersionCode,
	}); err != nil {
		logFail("Failed to export outputs, error: %s", err)
	}

	if err := fileutil.WriteStringToFile(configs.BuildGradlePth, updatedBuildGradleContent); err != nil {
		logFail("Failed to write build.gradle file, error: %s", err)
	}

	fmt.Println()
	log.Donef("%d versionCode updated", updatedVersionCodes)
	log.Donef("%d versionName updated", updatedVersionNames)
}
