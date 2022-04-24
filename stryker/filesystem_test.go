package stryker

import (
	"reflect"
	"testing"

	iowrap "github.com/spf13/afero"
)

func init() {
	FS = iowrap.NewMemMapFs()
	FSUtil = &iowrap.Afero{Fs: FS}

	FS.MkdirAll(".", 0755)
}

func TestGetStrykerConfigFileNames(t *testing.T) {
	t.Run("should get all the fileNames with correct format", func(t *testing.T) {
		correctFiles := []string{"fp1-stryker-config.json", "fp2-stryker-config.json"}
		wrongFiles := []string{"wrong-config.json", "appSettings.json"}

		allFileNames := append(correctFiles, wrongFiles...)
		for _, fileName := range allFileNames {
			FS.Create(fileName)
		}

		got := getStrykerConfigFileNames()
		want := correctFiles

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Wanted %v but got %v", want, got)
		}
	})

}
