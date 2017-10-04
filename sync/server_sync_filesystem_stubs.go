package sync

import (
	"errors"
	"bufio"
	"strings"
	"path"
	"path/filepath"
	"github.com/webdevops/go-shell"
	"github.com/webdevops/go-stubfilegenerator"
)

// General sync
func (filesystem *Filesystem) SyncStubs() {
	switch filesystem.Connection.GetType() {
	case "local":
		fallthrough
	case "ssh":
		filesystem.generateStubs()
	case "docker":
		errors.New("Docker not supported")
	}
}

// Sync filesystem using rsync
func (filesystem *Filesystem) generateStubs() {
	cmd := shell.Cmd(filesystem.Connection.CommandBuilder("find", filesystem.Path, "-type", "f")...)
	output := cmd.Run().Stdout.String()

	removeBasePath := filesystem.Path
	localBasePath := filesystem.localPath()

	// build list and filter it by user filter list
	fileList := []string{}
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		fileList = append(fileList, strings.TrimPrefix(scanner.Text(), removeBasePath))
	}
	fileList = filesystem.Filter.ApplyFilter(fileList)

	// generate stubs
	stubGen := stubfilegenerator.StubGenerator()
	for _, filePath := range fileList {
		localPath := path.Join(localBasePath, filePath)
		localAbsPath, _ := filepath.Abs(localPath)

		stubGen.TemplateVariables["PATH"] = localPath
		stubGen.GenerateStub(localAbsPath)
	}
}
