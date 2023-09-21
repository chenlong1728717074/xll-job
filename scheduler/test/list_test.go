package test

import (
	"fmt"
	"testing"
	"xll-job/scheduler/util"
)

func TestArray(t *testing.T) {
	list := util.NewArrayList[util.String]()
	fmt.Println(list)
}
