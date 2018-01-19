package main

import (
	"os"
	"os/exec"

	"strconv"

	"bufio"
	"strings"

	"regexp"

	"fmt"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/input"
)

func logFail(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}

func main() {
	// inputs
	buildGradlePth := os.Getenv("build_gradle_path")
	if err := input.ValidateIfPathExists(buildGradlePth); err != nil {
		logFail("Issue with input build_gradle_path - %s", err)
	}

	newVersionCode := os.Getenv("new_version_code")
	newVersionName := os.Getenv("new_version_name")

	versionCodeOffset, err := strconv.Atoi(os.Getenv("version_code_offset"))
	if err != nil {
		versionCodeOffset=0
	}

	log.Infof("Configs:")
	log.Printf("- build_gradle_path: %s", buildGradlePth)
	log.Printf("- new_version_code: %s", newVersionCode)
	log.Printf("- new_version_name: %s", newVersionName)
	// ---

	//
	// find versionName & versionCode with regexp
	buildGradleContent, err := fileutil.ReadStringFromFile(buildGradlePth)
	if err != nil {
		logFail("Failed to read build.gradle file, error: %s", err)
	}

	versionCodePattern := `^versionCode (?P<version_code>.*)`
	versionCodeRegexp := regexp.MustCompile(versionCodePattern)

	versionNamePattern := `^versionName "(?P<version_code>.*)"`
	versionNameRegexp := regexp.MustCompile(versionNamePattern)

	updatedLines := []string{}
	updatedVersionCodeNum := 0
	updatedVersionNameNum := 0

	reader := strings.NewReader(buildGradleContent)
	scanner := bufio.NewScanner(reader)
	lineNum := 0

	fmt.Println()
	log.Infof("Updating build.gradle file")

	for scanner.Scan() {
		lineNum++

		line := scanner.Text()

		if match := versionCodeRegexp.FindStringSubmatch(strings.TrimSpace(line)); len(match) == 2 {
			oldVersionCode := match[1]
			cmdLog1, err1 := exec.Command("bitrise", "envman", "add", "--key", "GRADLE_VERSION_CODE", "--value", oldVersionCode).CombinedOutput()
			if err1 != nil {
				logFail("Failed to expose output with envman, error: %#v | output: %s", err1, cmdLog1)
			}
			if newVersionCode != "" {
				iNewVersionCode, _ := strconv.Atoi(newVersionCode)
				newVersionCode=strconv.Itoa(iNewVersionCode+versionCodeOffset)
				updatedLine := strings.Replace(line, oldVersionCode, newVersionCode, -1)
				updatedVersionCodeNum++

				log.Printf("updating line (%d): %s -> %s", lineNum, line, updatedLine)

				updatedLines = append(updatedLines, updatedLine)
				cmdLog1, err1 := exec.Command("bitrise", "envman", "add", "--key", "GRADLE_VERSION_CODE", "--value", newVersionCode).CombinedOutput()
				if err1 != nil {
					logFail("Failed to expose output with envman, error: %#v | output: %s", err1, cmdLog1)
				}
				continue
			}
		}

		if match := versionNameRegexp.FindStringSubmatch(strings.TrimSpace(line)); len(match) == 2 {
			oldVersionName := match[1]
			cmdLog2, err2 := exec.Command("bitrise", "envman", "add", "--key", "GRADLE_VERSION_NAME", "--value", oldVersionName).CombinedOutput()
			if err2 != nil {
				logFail("Failed to expose output with envman, error: %#v | output: %s", err2, cmdLog2)
			}
			if newVersionName != "" {
				updatedLine := strings.Replace(line, oldVersionName, newVersionName, -1)
				updatedVersionNameNum++

				log.Printf("updating line (%d): %s -> %s", lineNum, line, updatedLine)

				updatedLines = append(updatedLines, updatedLine)
				cmdLog2, err2 := exec.Command("bitrise", "envman", "add", "--key", "GRADLE_VERSION_NAME", "--value", newVersionName).CombinedOutput()
				if err2 != nil {
					logFail("Failed to expose output with envman, error: %#v | output: %s", err2, cmdLog2)
				}
				continue
			}
		}

		updatedLines = append(updatedLines, line)
	}
	if err := scanner.Err(); err != nil {
		logFail("Failed to scann build.gradle file, error: %s", err)
	}
	// ---

	fmt.Println()
	log.Donef("%d versionCode updated", updatedVersionCodeNum)
	log.Donef("%d versionName updated", updatedVersionNameNum)

	updatedBuildGradleContent := strings.Join(updatedLines, "\n")

	cmdLog, err := exec.Command("bitrise", "envman", "add", "--key", "GRADLE_FILE_PATH", "--value", buildGradlePth).CombinedOutput()
	if err != nil {
		logFail("Failed to expose output with envman, error: %#v | output: %s", err, cmdLog)
	}

	if err := fileutil.WriteStringToFile(buildGradlePth, updatedBuildGradleContent); err != nil {
		logFail("Failed to write build.gradle file, error: %s", err)
	}
}