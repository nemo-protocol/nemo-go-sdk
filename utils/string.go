package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func FindFunctionInBytecode(bytecode string, functionName string) string {
	lines := strings.Split(bytecode, "\n")
	funcPattern := regexp.MustCompile(fmt.Sprintf(`public %s[(<][^{]*{`, functionName))
	for _, line := range lines {
		if funcPattern.MatchString(line) {
			return line
		}
	}
	return ""
}

func ExtractFunctionArgs(funcDef string) string {
	argPattern := regexp.MustCompile(`\(Arg0[^)]*\)`)
	matches := argPattern.FindString(funcDef)
	return matches
}

func Contains(elems []string, elem string) bool {
	for _, e := range elems {
		if elem == e {
			return true
		}
	}
	return false
}
