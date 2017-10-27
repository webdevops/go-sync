package sync

import (
	"github.com/webdevops/go-shell/commandbuilder"
)

type YamlCommandBuilderArgument struct {
	commandbuilder.Argument
}

func (yarg *YamlCommandBuilderArgument) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var argument commandbuilder.Argument
	err := unmarshal(&argument)
	if err == nil {
		// valid argument
		yarg.Argument = argument
	} else {
		// try to parse as string
		var config string
		err := unmarshal(&config)
		if err != nil {
			return err
		}

		err = yarg.Set(config)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ysa *YamlCommandBuilderArgument) String() string {
	return ""
}
