package main

import (
	"fmt"
	"testing"
	"strings"
)

func  TestReWithSpace(t *testing.T)  {
	str:=strings.Map(reWithSpace,`	testing
	`)

	for _,c:= range str{
		if c=='\t'{
			t.Errorf("tabs should have been replaced")
		}
		if c=='\n'{
			t.Errorf("newline should have been replaced")
		}
	}
	fmt.Println(str)
}