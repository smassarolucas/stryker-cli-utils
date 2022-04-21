package helpers

import (
	"fmt"
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

	for index, fileName := range files {
		if match := regex.MatchString(fileName); match {
			fmt.Printf("File name of file %d is %v\n", index, fileName)
		}
	}

	return make([]string, 0)
}
