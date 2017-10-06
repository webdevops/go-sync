package main

import (
	"fmt"
	"context"
	"log"
	"runtime"
	"strings"
	"github.com/google/go-github/github"
	"github.com/inconshreveable/go-update"
	"net/http"
)

type SelfUpdateCommand struct {
	CurrentVersion      string
	GithubOrganization  string
	GithubRepository    string
	GithubAssetTemplate string
}

func (conf *SelfUpdateCommand) Execute(args []string) error {
	fmt.Println("Starting self update")

	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), conf.GithubOrganization, conf.GithubRepository)

	if _, ok := err.(*github.RateLimitError); ok {
		log.Println("GitHub rate limit, please try again later")
	}

	fmt.Println(fmt.Sprintf(" - latest version is %s", release.GetName()))

	// check if latest version is current version
	if release.GetName() == conf.CurrentVersion {
		fmt.Println(" - already using the latest version")
		return nil
	}

	// translate OS names
	os := runtime.GOOS
	switch (runtime.GOOS) {
	case "darwin":
		os = "osx"
	}

	// translate arch names
	arch := runtime.GOARCH
	switch (arch) {
	case "amd64":
		arch = "x64"
	case "386":
		arch = "x32"
	}

	// build asset name
	assetName := conf.GithubAssetTemplate
	assetName = strings.Replace(assetName, "%OS%", os, -1)
	assetName = strings.Replace(assetName, "%ARCH%", arch, -1)

	// search assets in release for the desired filename
	fmt.Println(fmt.Sprintf(" - searching for asset \"%s\"", assetName))
	for _, asset := range release.Assets {
		if asset.GetName() == assetName {
			downloadUrl := asset.GetBrowserDownloadURL()
			fmt.Println(fmt.Sprintf(" - found new update url \"%s\"", downloadUrl))
			conf.runUpdate(downloadUrl)
			fmt.Println(fmt.Sprintf(" - finished update to version %s", release.GetName()))
			return nil
		}
	}

	fmt.Println(" - unable to find asset, please contact maintainer")
	return nil
}

func (conf *SelfUpdateCommand) runUpdate(url string) error {
	fmt.Println(" - downloading update")
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println(" - applying update")
	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		// error handling
		fmt.Println(fmt.Sprintf(" - updating application failed: %s", err))
	}
	return err
}
