package stryker

import (
	"errors"
	"io/fs"
	"log"
	"path/filepath"
	"regexp"

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
