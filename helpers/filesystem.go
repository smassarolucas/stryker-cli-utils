package helpers

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

const (
	strykerConfigNameSufix = "-stryker-config.json"
	strykerConfigRegex     = ".*\\" + strykerConfigNameSufix + "$"
	dirName                = "."
	strykerReportsDirName  = "./StrykerOutput"
	mutationReportRegex    = ".*\\-report.html$"
)

func GetStrykerConfigFileNames() []string {
	dir, err := os.Open(dirName)
	if err != nil {
		log.Fatalf("Couldn't open directory %v because of error %v", dirName, err.Error())
	}
	defer dir.Close()

	files, err := dir.Readdirnames(-1)
	if err != nil {
		log.Fatalf("Couldn't read names of directory %v because of error %v", dirName, err.Error())
	}

	regex, err := regexp.Compile(strykerConfigRegex)
	if err != nil {
		log.Fatalf("Couldn't compile the regex %v because of error %v", strykerConfigRegex, err.Error())
	}

	fileNames := []string{}

	for _, fileName := range files {
		if match := regex.MatchString(fileName); match {
			fileNames = append(fileNames, fileName)
		}
	}

	return fileNames
}

func GetMutationReportsFilePaths() []string {
	regex, err := regexp.Compile(mutationReportRegex)
	if err != nil {
		log.Fatalf("Couldn't compile the regex %v because of error %v", strykerConfigRegex, err.Error())
	}

	filePaths := []string{}
	err = filepath.Walk(strykerReportsDirName, func(path string, info fs.FileInfo, err error) error {
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
