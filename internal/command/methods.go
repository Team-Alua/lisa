package command

import (
	"fmt"
	"errors"
	"archive/zip"
	"github.com/Team-Alua/lisa/internal/client"
	"github.com/Team-Alua/lisa/internal/cecie"
)



type commandFunction func(*cecie.Connection, *client.Content, *zip.ReadCloser) (string, error)

var executors = map[string]commandFunction {
	"CREATE": executeCreate,
	"DUMP": executeDump,
}


func Execute(cc * cecie.Connection, c * client.Content, rc * zip.ReadCloser) (string, error) {
	executor, ok := executors[c.Type]
	if !ok {
		return "", errors.New(fmt.Sprintf("%s is not implemented", string(c.Type)))
	}
	outZip, err := executor(cc, c, rc)
	return outZip, err
}

