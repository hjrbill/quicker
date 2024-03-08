package test

import (
	"fmt"
	"github.com/hjrbill/quicker/demo/job"
	"testing"
)

func TestBits(t *testing.T) {
	fmt.Printf("%064b\n", job.GetClassBits([]string{"五月天", "北京", "资讯", "热点"}))
}
