package test

import (
	"fmt"
	"testing"
	"xll-job/utils"
)

func TestArray(t *testing.T) {
	list := utils.NewArrayList[utils.String]()
	fmt.Println(list)
}
