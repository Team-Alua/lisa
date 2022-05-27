package command

import (
	"github.com/Team-Alua/lisa/internal/cecie"
	"github.com/Team-Alua/lisa/internal/client"
	"time"
	"fmt"
	"os"
	"strings"
	"archive/zip"
)

func executeCreate(cc *cecie.Connection, c * client.Content, rc * zip.ReadCloser) (string, error) {

	if err := cc.SendCommand("ReserveSaveContainer"); err != nil {
		return "", err
	}

	if err := cc.SendTarget(c.Target); err != nil {
		return "", err
	}

	if err := cc.CheckOkay(); err != nil {
		return "", err
	}

	if err := cc.SendCommand("DeleteSaveContainer"); err != nil {
		return "", err
	}

	if err := cc.CheckOkay(); err != nil {
		return "", err
	}

	if err := cc.SendCommand("MountSaveContainer"); err != nil {
		return "", err
	}

	if err := cc.CheckOkay(); err != nil {
		return "", err
	}

	if err := cc.SendCommand("ModifySaveContainer"); err != nil {
		return "", err
	}

	if err := cc.CheckOkay(); err != nil {
		return "", err
	}


	for _, file := range rc.File {
		if !file.Mode().IsRegular() {
			continue
		}

		if strings.HasPrefix(file.Name, "files/") {
			relativePath := file.Name[6:]

			if err := cc.SendCommand("UploadFile"); err != nil {
				return "", err
			}
			if err := cc.SendFileHeader(relativePath, file.UncompressedSize64); err != nil {
				return "", err
			}
			fr, err := file.Open()
			if err != nil {
				return "", err
			}

			if err := cc.SendZipFile(fr); err != nil {
				fr.Close()
				return "", err
			}
			fr.Close()

			if err := cc.CheckOkay(); err != nil {
				return "", err
			}
		}
	}

	if err := cc.SendCommand("Finish"); err != nil {
		return "", err
	}

	if err := cc.CheckOkay(); err != nil {
		return "", err
	}

	if err := cc.SendCommand("UnmountSaveContainer"); err != nil {
		return "", err
	}

	if err := cc.CheckOkay(); err != nil {
		return "", err
	}

	if err := cc.SendCommand("DownloadSaveContainer"); err != nil {
		return "", err
	}

	if err := cc.CheckOkay(); err != nil {
		return "", err
	}

	zipName :=  fmt.Sprintf("%s-%s-%d.zip", c.Target.TitleId, c.Target.DirectoryName, time.Now().Unix())
	outZip, err := os.Create(zipName)
	if err != nil {
		return "", err
	}
	defer outZip.Close()
	w := zip.NewWriter(outZip)
	defer w.Close()
	outRaw := fmt.Sprintf("PS4/SAVEDATA/%08x/%s/%s", c.AccountId, c.Target.TitleId, c.Target.DirectoryName)

	exportFiles := []string{
		outRaw + ".bin",
		outRaw,
	}

	exportFileSizes := []int64{
		96,
		32768 * int64(c.Target.Blocks),
	}

	for index, fileName := range exportFiles {
		wtr, err := w.Create(fileName)
		if err != nil {
			return zipName, err
		}

		if err := cc.ReceiveZipFile(wtr, exportFileSizes[index]); err != nil {
			return zipName, err
		}


		if err := cc.CheckOkay(); err != nil {
			return zipName, err
		}
	}

// 	if err := cc.SendCommand("DeleteSaveContainer"); err != nil {
// 		return zipName, err
// 	}
// 
// 	if err := cc.CheckOkay(); err != nil {
// 		return zipName, err
// 	}

	return zipName, nil
}
