package sync

import (
	"regexp"
)

func (filter *Filter) compile() {
	if len(filter.excludeRegexp) == 0 {
		filter.excludeRegexp = make([]*regexp.Regexp, len(filter.Exclude))
		for i, filterVal := range filter.Exclude {
			filter.excludeRegexp[i] = regexp.MustCompile(filterVal)
		}
	}

	if len(filter.includeRegexp) == 0 {
		filter.includeRegexp = make([]*regexp.Regexp, len(filter.Include))
		for i, filterVal := range filter.Include {
			filter.includeRegexp[i] = regexp.MustCompile(filterVal)
		}
	}
}

func (filter *Filter) CalcExcludes(lines []string) []string {
	filter.compile()
	return filter.calculateMatching(filter.excludeRegexp, lines)
}

func (filter *Filter) CalcIncludes(lines []string) []string {
	filter.compile()
	return filter.calculateMatching(filter.includeRegexp, lines)
}

func (filter *Filter) calculateMatching(regexpList []*regexp.Regexp, lines []string) []string {
	var ret []string

	for _, filterRegexp := range regexpList {
		for _, value := range lines {
			if filterRegexp.MatchString(value) == true {
				ret = append(ret, value)
			}
		}
	}

	return ret
}
