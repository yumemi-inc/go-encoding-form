package form

import (
	"strings"
)

type Tag struct {
	Key       string
	OmitEmpty bool
}

func ParseTag(s string) Tag {
	parts := strings.Split(s, ",")
	tag := Tag{
		Key: parts[0],
	}

	for _, v := range parts[1:] {
		if v == "omitempty" {
			tag.OmitEmpty = true
		}
	}

	return tag
}
