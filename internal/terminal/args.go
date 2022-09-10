package terminal

import (
	"regexp"
	"strings"

	"github.com/reconmap/shared-lib/pkg/api"
)

func ReplaceArgs(command *api.Command, vars []string) string {
	var updatedArgs = command.ContainerArgs
	for _, v := range vars {
		var tokens = strings.Split(v, "=")
		var validID = regexp.MustCompile("{{{" + tokens[0] + ".*?}}}")
		updatedArgs = validID.ReplaceAllString(updatedArgs, tokens[1])
	}

	return updatedArgs
}
