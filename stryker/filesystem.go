package stryker

import (
	"bufio"
	"errors"
	"io/fs"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	iowrap "github.com/spf13/afero"
)

const (
	strykerConfigNameSufix = "-stryker-config.json"
	isStrykerConfigRegex   = ".*\\" + strykerConfigNameSufix + "$"
	currentDirName         = "."
)

var (
	FS     iowrap.Fs
	FSUtil *iowrap.Afero
)

func init() {
	FS = iowrap.NewOsFs()
	FSUtil = &iowrap.Afero{Fs: FS}
}

func getStrykerConfigFileNames() ([]string, error) {
	dir, err := FS.Open(currentDirName)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	files, err := dir.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	regex, err := regexp.Compile(isStrykerConfigRegex)
	if err != nil {
		return nil, err
	}

	fileNames := []string{}
	for _, fileName := range files {
		if match := regex.MatchString(fileName); match {
			fileNames = append(fileNames, fileName)
		}
	}

	if len(fileNames) == 0 {
		return nil, errors.New("there are no Stryker config files")
	}

	return fileNames, nil
}

func getProjectsToMutate() ([]string, string) {
	testProjectPath := getTestProjectPath()
	mutableProjects := getMutableProjects(testProjectPath)
	return mutableProjects, testProjectPath
}

const (
	isTestProjectRegex = `.*Tests.csproj`
)

func getTestProjectPath() string {
	regex, err := regexp.Compile(isTestProjectRegex)
	if err != nil {
		log.Fatalf("Couldn't compile the regex %v because of error %v", isTestProjectRegex, err.Error())
	}

	var testProjectPath string
	err = FSUtil.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && (info.Name() == "obj" || info.Name() == "bin") {
			return filepath.SkipDir
		}

		if match := regex.MatchString(info.Name()); match {
			log.Println(testProjectPath)
			testProjectPath = path
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Got error %v", err.Error())
	}
	return testProjectPath
}

const (
	isTestableProject      = `.*.csproj\" />`
	projectReferencePrefix = `<ProjectReference Include="`
	projectReferenceSuffix = `" />`
)

func getMutableProjects(testProjectPath string) []string {
	file, err := FS.Open(testProjectPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	regex, err := regexp.Compile(isTestableProject)
	if err != nil {
		log.Fatalf("Couldn't compile the regex %v because of error %v", isTestProjectRegex, err.Error())
	}

	filePaths := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if match := regex.MatchString(scanner.Text()); match {
			line = strings.Trim(line, " ")
			line = strings.TrimPrefix(line, projectReferencePrefix)
			line = strings.TrimSuffix(line, projectReferenceSuffix)
			splittedLine := strings.Split(line, "\\")
			line = splittedLine[len(splittedLine)-1]
			filePaths = append(filePaths, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return filePaths
}

const (
	strykerReportsDirName = "./StrykerOutput"
	isMutationReportRegex = ".*\\-report.html$"
)

func getMutationReportsFilePaths() []string {
	regex, err := regexp.Compile(isMutationReportRegex)
	if err != nil {
		log.Fatalf("Couldn't compile the regex %v because of error %v", isStrykerConfigRegex, err.Error())
	}

	filePaths := []string{}
	err = FSUtil.Walk(strykerReportsDirName, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if match := regex.MatchString(info.Name()); match {
			filePaths = append(filePaths, path)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Couldn't read filePaths for reports because of error %v", err.Error())
	}

	return filePaths
}

func writeToFile(content, fileName string) string {
	f, err := FS.Create(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	f.WriteString(content)

	filePath, _ := filepath.Abs(f.Name())
	return filePath
}

const (
	strykerOutputFolder = "StrykerOutput"
)

func deleteStrykerOutputFolder() {
	err := FS.RemoveAll(strykerOutputFolder)
	if err != nil {
		log.Fatalf("Couldn't remove %v because of %v", strykerOutputFolder, err)
	}
}
