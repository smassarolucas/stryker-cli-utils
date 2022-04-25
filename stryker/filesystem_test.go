package stryker

import (
	"reflect"
	"testing"

	iowrap "github.com/spf13/afero"
)

func init() {
	FS = iowrap.NewMemMapFs()
	FSUtil = &iowrap.Afero{Fs: FS}
}

func TestGetStrykerConfigFileNames(t *testing.T) {
	t.Run("should get all the fileNames with correct format", func(t *testing.T) {
		FS.MkdirAll(".", 0755)
		correctFiles := []string{"fp1-stryker-config.json", "fp2-stryker-config.json"}
		wrongFiles := []string{"wrong-config.json", "appSettings.json"}

		allFileNames := append(correctFiles, wrongFiles...)
		for _, fileName := range allFileNames {
			FS.Create(fileName)
		}

		got, err := getStrykerConfigFileNames()
		want := correctFiles

		if err != nil {
			t.Fatalf("Wasn't expecting an error, but got %v", err.Error())
		}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Wanted %v but got %v", want, got)
		}
	})

	t.Run("should return error when couldn't open current directory", func(t *testing.T) {
		FS.RemoveAll(".")

		got, err := getStrykerConfigFileNames()

		if err == nil {
			t.Fatalf("Was expecting an error, but got none")
		}

		if got != nil {
			t.Fatalf("Wanted no value, but got %v", got)
		}
	})

	t.Run("should return error when there are no config files", func(t *testing.T) {
		FS.MkdirAll(".", 0755)
		FS.Chmod(".", 0755)

		got, err := getStrykerConfigFileNames()
		want := "there are no Stryker config files"
		if err == nil {
			t.Fatalf("Was expecting an error, but got none")
		}

		if err.Error() != want {
			t.Fatalf("Wanted error %v, but got error %v", want, err)
		}

		if got != nil {
			t.Fatalf("Wanted no value, but got %v", got)
		}
	})
}
