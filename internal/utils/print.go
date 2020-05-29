package utils

import "fmt"

type VSlice []string

func (s VSlice) String() string {
	var str string
	for _, i := range s {
		str += fmt.Sprintf("%s\n", i)
	}
	return str
}
