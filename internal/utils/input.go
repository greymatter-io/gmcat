package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Input() (output string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("-->")
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	text = strings.TrimSpace(text)
	output = output + " " + text

	return output
}

//determine if the val is in slice.
//present: (first_location, bool)
//not present: (-1,false)
func Verify(val string, slice []string) (location int, valid bool) {
	for i, item := range slice {
		val = strings.TrimSpace(val)
		item = strings.TrimSpace(item)
		if item == val {
			return i, true
		}
	}
	return
}

//takes a string slice and asks for input.  if matches will return <string>, true.
func InputFromList(list []string) (input string, valid bool) {
	valid = false
	if len(list) > 0 {
		for valid == false {
			input = strings.TrimSpace(Input())
			_, valid = Verify(input, list)
			if !valid {
				fmt.Printf("Invalid choose from %s\n", list)
			}
		}
	} else {
		fmt.Print("No list passed")
	}
	return
}
