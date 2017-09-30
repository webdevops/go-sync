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
func (filesystem *Filesystem) SyncStubs(server *Server) {
	switch server.Connection.GetType() {
	case "ssh":
		filesystem.generateStubs(server)
	case "docker":
		errors.New("Docker not supported")
	}
}

// Sync filesystem using rsync
func (filesystem *Filesystem) generateStubs(server *Server) {
	cmd := shell.Cmd(server.Connection.CommandBuilder("find", filesystem.Path, "-type f")...)
	output := cmd.Run().Stdout.String()

	removeBasePath := filesystem.Path
	localBasePath := filesystem.localPath(server)

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

		stubGen.TemplateVariables["PATH"] = localPath
		localAbsPath, _ := filepath.Abs(localPath)
		stubGen.GenerateStub(localAbsPath)
	}
}
