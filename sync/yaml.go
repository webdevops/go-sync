package sync

import "strings"

type YamlStringArray struct {
	Multi  []string
	Single  string
}

func (ysa *YamlStringArray) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(&multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		ysa.Single = single
	} else {
		ysa.Multi = multi
	}
	return nil
}

func (ysa *YamlStringArray) String() string {
	return ysa.ToString(";")
}

func (ysa *YamlStringArray) ToString(sep string) string {
	if len(ysa.Multi) >= 1 {
		return strings.Join(ysa.Multi, sep)
	} else {
		return ysa.Single
	}
}

func (ysa *YamlStringArray) Array() []string {
	if len(ysa.Multi) >= 1 {
		return ysa.Multi
	} else if ysa.Single != "" {
		return []string{ysa.Single}
	} else {
		return []string{}
	}
}
