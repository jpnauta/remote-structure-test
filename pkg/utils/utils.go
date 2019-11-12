package utils

import (
	"regexp"

	"github.com/sirupsen/logrus"
)

func CompileAndRunRegex(regex string, base string, shouldMatch bool) bool {
	r, rErr := regexp.Compile(regex)
	if rErr != nil {
		logrus.Errorf("Error compiling regex %s : %s", regex, rErr.Error())
		return false
	}
	return shouldMatch == r.MatchString(base)
}
