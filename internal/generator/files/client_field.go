package files

import (
	"bytes"
	"sort"
	"strconv"
	"strings"
)

func NewField(name string) *Field {
	return &Field{
		Name: name,
	}
}

type Field struct {
	Comment string
	Name    string
	Type    string
	Tags    Tags
}

func (f *Field) AddTag(tagKey, tagValue string, flags ...string) {
	if f.Tags == nil {
		f.Tags = Tags{}
	}
	f.Tags[tagKey] = NewTag(tagKey, tagValue).WithFlags(flags...)
}

func AsCommentLines(comment string) string {
	buf := &bytes.Buffer{}
	for _, line := range strings.Split(comment, "\n") {
		buf.WriteString(`// ` + line + `
`)
	}
	return buf.String()
}

func (f *Field) Anonymous() bool {
	return f.Type == ""
}

func (f *Field) String() string {
	if f.Anonymous() {
		return AsCommentLines(f.Comment) + f.Name + ` ` + f.Tags.String() + `
`
	}
	return AsCommentLines(f.Comment) + f.Name + ` ` + f.Type + ` ` + f.Tags.String() + `
`
}

func NewTag(key, value string) *Tag {
	return &Tag{
		Key:   key,
		Value: value,
	}
}

type Tag struct {
	Key   string
	Value string
	Flags map[string]bool
}

func (tag Tag) WithFlags(flags ...string) *Tag {
	for _, flag := range flags {
		if tag.Flags == nil {
			tag.Flags = map[string]bool{}
		}
		tag.Flags[flag] = true
	}
	return &tag
}

func (tag *Tag) String() string {
	flags := make([]string, 0)
	if len(tag.Flags) > 0 {
		for flag := range tag.Flags {
			flags = append(flags, flag)
		}
		sort.Strings(flags)
	}

	return tag.Key + ":" + strconv.Quote(strings.Join(append([]string{tag.Value}, flags...), ","))
}

type Tags map[string]*Tag

func (tags Tags) String() string {
	if len(tags) > 0 {
		tagList := make([]string, 0)
		for tag := range tags {
			tagList = append(tagList, tag)
		}
		sort.Strings(tagList)

		buf := &bytes.Buffer{}

		tagCount := len(tagList)

		for i, tag := range tagList {
			if i == 0 {
				buf.WriteString("`")
			} else {
				buf.WriteString(" ")
			}

			buf.WriteString(tags[tag].String())

			if i == tagCount-1 {
				buf.WriteString("`")
			}
		}
		return buf.String()
	}
	return ""
}
