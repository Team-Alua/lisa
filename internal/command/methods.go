package command

import (
	"strings"
	"archive/zip"
	"github.com/Team-Alua/lisa/internal/client"
	"github.com/Team-Alua/lisa/internal/cecie"
)



type commandFunction func(*cecie.Connection, *client.Content, *zip.ReadCloser) error

var executors = map[string]commandFunction {
	"CREATE": executeCreate,
}

func executeCreate(cc *cecie.Connection, c * client.Content, rc * zip.ReadCloser) error {
	cc.SendCommand("ReserveSaveContainer")
	cc.SendTarget(c.Target)

	if err := cc.CheckOkay(); err != nil {
		return err
	}

	cc.SendCommand("DeleteSaveContainer")
	if err := cc.CheckOkay(); err != nil {
		return err
	}

	cc.SendCommand("MountSaveContainer")
	if err := cc.CheckOkay(); err != nil {
		return err
	}

	cc.SendCommand("ModifySaveContainer")
	if err := cc.CheckOkay(); err != nil {
		return err
	}

	for _, file := range rc.File {
		if strings.HasPrefix(file.Name, "files/") {
			relativePath := file.Name[6:]
			cc.SendCommand("UploadFile")
			cc.SendFileHeader(relativePath, file.UncompressedSize64)
			fr, err := file.Open()
			if err != nil {
				return err
			}
			
			if err := cc.SendZipFile(fr); err != nil {
				return err
			}
		}
	}

	cc.SendCommand("Finish")
	if err := cc.CheckOkay(); err != nil {
		return err
	}

	cc.SendCommand("UnmountSaveContainer")
	if err := cc.CheckOkay(); err != nil {
		return err
	}

	return nil
}

func Execute(cc * cecie.Connection, c * client.Content, rc * zip.ReadCloser) error {
	err := executors[c.Type](cc, c, rc)
	return err
}

