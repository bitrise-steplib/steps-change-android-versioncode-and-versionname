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
	return nil
}

type updateFn func(line string, lineNum int, matches []string) string

func findAndUpdate(reader io.Reader, update map[string]updateFn) (string, error) {
	scanner := bufio.NewScanner(reader)
	updatedLines := []string{}

	reByPattern := make(map[string]*regexp.Regexp, len(update))
	for pattern := range update {
		reByPattern[pattern] = regexp.MustCompile(pattern)
	}

	for lineNum := 0; scanner.Scan(); lineNum++ {
		line := scanner.Text()

		updated := false
		for pattern, fn := range update {
			re := reByPattern[pattern]
			if match := re.FindStringSubmatch(strings.TrimSpace(line)); len(match) == 2 {
				updatedLine := fn(line, lineNum, match)
				updatedLines = append(updatedLines, updatedLine)
				updated = true
				break
			}
		}
		if !updated {
			updatedLines = append(updatedLines, line)
		}
	}

	return strings.Join(updatedLines, "\n"), scanner.Err()
}

func exportOutput(outputs map[string]string) error {
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

	finalVersionCode := ""
	finalVersionName := ""

	updatedVersionCodeNum := 0
	updatedVersionNameNum := 0

	updatedBuildGradleContent, err := findAndUpdate(f, map[string]updateFn{
		`^versionCode (?P<version_code>.*)`: func(line string, lineNum int, match []string) string {
			oldVersionCode := match[1]
			finalVersionCode = oldVersionCode
			updatedLine := ""

			if configs.NewVersionCode != "" {
				finalVersionCode = strconv.Itoa(newVersionCode)
				updatedLine = strings.Replace(line, oldVersionCode, finalVersionCode, -1)
				updatedVersionCodeNum++
				log.Printf("updating line (%d): %s -> %s", lineNum, line, updatedLine)
			}

			return updatedLine
		},
		`^versionName "(?P<version_code>.*)"`: func(line string, lineNum int, match []string) string {
			oldVersionName := match[1]
			finalVersionName = oldVersionName
			updatedLine := ""

			if configs.NewVersionName != "" {
				finalVersionName = configs.NewVersionName
				updatedLine = strings.Replace(line, oldVersionName, finalVersionName, -1)
				updatedVersionNameNum++
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
	if err := exportOutput(map[string]string{
		"ANDROID_VERSION_NAME": finalVersionName,
		"ANDROID_VERSION_CODE": finalVersionCode,
	}); err != nil {
		logFail("Failed to export outputs, error: %s", err)
	}

	if err := fileutil.WriteStringToFile(configs.BuildGradlePth, updatedBuildGradleContent); err != nil {
		logFail("Failed to write build.gradle file, error: %s", err)
	}

	fmt.Println()
	log.Donef("%d versionCode updated", updatedVersionCodeNum)
	log.Donef("%d versionName updated", updatedVersionNameNum)
}
