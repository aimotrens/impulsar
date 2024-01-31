package model

import "strings"

type FlagArray []string

func (i *FlagArray) String() string {
	return strings.Join(*i, ", ")
}

func (i *FlagArray) Set(value string) error {
	*i = append(*i, value)
	return nil
}
