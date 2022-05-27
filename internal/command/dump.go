package command

import (
	"github.com/Team-Alua/lisa/internal/cecie"
	"github.com/Team-Alua/lisa/internal/client"
	"fmt"
	"os"
	"time"
	"archive/zip"
)

func executeDump(cc *cecie.Connection, c * client.Content, rc * zip.ReadCloser) (string, error) {
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

	if err := cc.SendCommand("UnmountSaveContainer"); err != nil {
		return "", err
	}

	if err := cc.CheckOkay(); err != nil {
		return "", err
	}

	if err := cc.SendCommand("UploadSaveContainer"); err != nil {
		return "", err
	}

	if err := cc.CheckOkay(); err != nil {
		return "", err
	}


	pfsFilePath := fmt.Sprintf("container/%s/%s.bin", c.Target.TitleId, c.Target.DirectoryName);

	fr, err := rc.Open(pfsFilePath)
	
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

	rawFilePath := fmt.Sprintf("container/%s/%s", c.Target.TitleId, c.Target.DirectoryName);

	fr, err = rc.Open(rawFilePath)
	
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

	if err := cc.SendCommand("MountSaveContainer"); err != nil {
		return "", err
	}

	if err := cc.CheckOkay(); err != nil {
		return "", err
	}

	if err := cc.SendCommand("DumpSaveContainer"); err != nil {
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
	zw := zip.NewWriter(outZip)
	defer zw.Close()

	if err := cc.ReceiveContainerDump(zw, "files/"); err != nil {
		return zipName, err
	}

	if err := cc.SendCommand("UnmountSaveContainer"); err != nil {
		return zipName, err
	}

	if err := cc.CheckOkay(); err != nil {
		return zipName, err
	}

	if err := cc.SendCommand("DeleteSaveContainer"); err != nil {
		return zipName, err
	}

	if err := cc.CheckOkay(); err != nil {
		return zipName, err
	}


	return zipName, nil
}

