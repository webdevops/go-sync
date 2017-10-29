package sync

import (
	"regexp"
)

// Compile filter regexp (and cache them)
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

// Apply filter (exclude/include) and get filtered list
func (filter *Filter) ApplyFilter(lines []string) []string {
	filter.compile()

	if len(filter.Include) > 0 {
		lines = filter.calculateMatching(filter.includeRegexp, lines)
	}

	if len(filter.Exclude) > 0 {
		excludes := filter.calculateMatching(filter.excludeRegexp, lines)

		tmp := []string{}
		for _, line := range lines {
			if ! stringInSlice(line, excludes) {
				tmp = append(tmp, line)
			}
		}

		lines = tmp
	}

	return lines
}

// Apply exclude filter only and get filtered excludes
func (filter *Filter) CalcExcludes(lines []string) []string {
	filter.compile()
	return filter.calculateMatching(filter.excludeRegexp, lines)
}

// Apply includes filter only and get filtered includes
func (filter *Filter) CalcIncludes(lines []string) []string {
	filter.compile()
	return filter.calculateMatching(filter.includeRegexp, lines)
}

// Calculate matches using regexp array
func (filter *Filter) calculateMatching(regexpList []*regexp.Regexp, lines []string) (matches []string) {
	for _, filterRegexp := range regexpList {
		for _, value := range lines {
			if filterRegexp.MatchString(value) == true {
				matches = append(matches, value)
			}
		}
	}

	return
}

// check if string exists in slice
func stringInSlice(a string, list []string) (status bool) {
	status = false

	for _, b := range list {
		if b == a {
			status = true
			return
		}
	}
	return
}
