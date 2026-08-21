package unicode

import (
	"fmt"
	"regexp"
)

var (
	re = regexp.MustCompile(`^(?:U\+[0-9A-F]{4}|U\+[1-9A-F][0-9A-F]{4,5})$`)
)

type (
	CodePoint string
)

func (cp CodePoint) String() string {
	s := string(cp)

	if !re.Match([]byte(s)) {
		return s
	}

	var r rune
	_, err := fmt.Sscanf(s, "%U", &r)
	if err != nil {
		return s
	}
	return string(r)
}
