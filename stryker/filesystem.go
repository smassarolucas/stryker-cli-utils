package stryker

import (
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

func getStrykerConfigFileNames() []string {
	dir, err := FS.Open(currentDirName)
	if err != nil {
		log.Fatalf("Couldn't open directory %v because of error %v", currentDirName, err.Error())
	}
	defer dir.Close()

	files, err := dir.Readdirnames(-1)
	if err != nil {
		log.Fatalf("Couldn't read names of directory %v because of error %v", currentDirName, err.Error())
	}

	regex, err := regexp.Compile(isStrykerConfigRegex)
	if err != nil {
		log.Fatalf("Couldn't compile the regex %v because of error %v", isStrykerConfigRegex, err.Error())
	}

	fileNames := []string{}
	for _, fileName := range files {
		if match := regex.MatchString(fileName); match {
			fileNames = append(fileNames, fileName)
		}
	}

	return fileNames
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
