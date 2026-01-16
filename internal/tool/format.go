package tool

import (
	"go/format"
)

func FormatGoCode(code string) string {
	formatted, err := format.Source([]byte(code))
	if err != nil {
		return code
	}
	return string(formatted)
}
