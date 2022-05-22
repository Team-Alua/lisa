package validator

import (
	"archive/zip"
	"github.com/goccy/go-yaml"
	"github.com/Team-Alua/lisa/internal/client"
	"regexp"
	"errors"
	"strings"
)

func checkContentFileTarget(ct *client.ContentTarget) error {
	titleIdRegex := regexp.MustCompile(`[A-Z]{4}[0-9]{5}`)
	if !titleIdRegex.MatchString(ct.TitleId) {
		return errors.New("Invalid Title Id supplied")
	}
	if len(ct.DirectoryName) == 0 || len(ct.DirectoryName) >= 32 {
		return errors.New("Invalid Directory Name supplied")
	}

	if ct.Blocks < 96 || ct.Blocks > 32768 {
		return errors.New("Invalid Blocks supplied")
	}

	return nil
}

func checkContentFile(content * client.Content) error {

	validCommandTypes := map[string]bool {
		"CREATE": true,
		"MODIFY": true,
		"DUMP": true,
	}

	if !validCommandTypes[content.Type] {
		return errors.New("Invalid Type supplied")
	}

	accountIdRegex := regexp.MustCompile(`(0x)?[0-9a-f]{16}`)
	if !accountIdRegex.MatchString(content.AccountId) {
		return errors.New("Invalid AccountId supplied")
	}


	return checkContentFileTarget(content.Target)
}

func checkZipFolderStructure(rc * zip.ReadCloser, content * client.Content) error {
	checkContainerFolder := false
	checkFilesFolder := false

	if content.Type == "MODIFY" || content.Type == "DUMP" {
		checkContainerFolder = true
	}

	if content.Type == "CREATE" || content.Type == "MODIFY" {
		checkFilesFolder = true
	}

	if checkFilesFolder {
		numOfFiles := 0
		for _,f := range rc.File {
			if f.Mode().IsDir() {
				continue
			}

			if strings.HasPrefix(f.Name, "files/") {
				numOfFiles += 1
			}
		}
		if numOfFiles == 0 {
			return errors.New("Must have non empty files folder.")
		}
	}

	if checkContainerFolder {

		containerFiles := map[string]uint64 {}

		for _,f := range rc.File {
			if f.Mode().IsDir() {
				continue
			}

			if strings.HasPrefix(f.Name, "container/") {
				containerFiles[f.Name] = f.UncompressedSize64
			}
		}

		if len(containerFiles) != 2 {
			return errors.New("Must supply exactly 2 save container files.")
		}

		titleId := content.Target.TitleId
		dirName := content.Target.DirectoryName
		rawDataPath := "container/" + titleId + "/" + dirName
		pfsKeyPath := rawDataPath + ".bin"

		if containerFiles[pfsKeyPath] != 96 {
			return errors.New("Did not supply the correct .bin file.")
		}

		targetSize := uint64(content.Target.Blocks) * 32768


		if containerFiles[rawDataPath] != targetSize {
			return errors.New("Did not supply the correct raw data file.")
		} 
	}

	return nil
}

func CheckZip(zipPath string) error {
	rc, err :=  zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer rc.Close()

	cRc, err := rc.Open("content.yml")
	if err != nil {
		return err
	}
	var content client.Content
	dec := yaml.NewDecoder(cRc)

	if err := dec.Decode(&content); err != nil {
		return err
	}

	err = checkContentFile(&content)
	if err != nil {
		return err
	}

	return checkZipFolderStructure(rc, &content)
}

