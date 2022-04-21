package helpers

import (
	"log"
	"os"
	"regexp"
)

const (
	strykerConfigNameSufix = "-stryker-config.json"
	strykerConfigRegex     = ".*\\" + strykerConfigNameSufix + "$"
)

func GetStrykerConfigFileNames() []string {
	dirName := "."

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
