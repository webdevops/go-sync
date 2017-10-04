package main

import (
	"fmt"
)

type VersionCommand struct {
	ShowOnlyVersion  bool `long:"dump"  description:"show only version number and exit"`
	Name string
	Version string
	Author string
}

// Show app version
func (conf *VersionCommand) Execute(args []string) error {
	if conf.ShowOnlyVersion {
		fmt.Println(conf.Version)
	} else {
		fmt.Println(fmt.Sprintf("%s version %s", conf.Name, conf.Version))
		fmt.Println(fmt.Sprintf("Copyright (C) 2017 %s", conf.Author))
	}

	return nil
}
