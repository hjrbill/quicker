package test

import (
	"fmt"
	reverseindex "github.com/hjrbill/quicker/internal/reverse_index"
	"github.com/huandu/skiplist"
	"strings"
	"testing"
)

var l1, l2, l3 *skiplist.SkipList

func init() {
	l1 = skiplist.New(skiplist.Uint64)
	l1.Set(uint64(5), 0)
	l1.Set(uint64(1), 0)
	l1.Set(uint64(4), 0)
	l1.Set(uint64(9), 0)
	l1.Set(uint64(11), 0)
	l1.Set(uint64(7), 0)
	//skip list 内部会自动做排序，排完序之后为 1 4 5 7 9 11

	l2 = skiplist.New(skiplist.Uint64)
	l2.Set(uint64(4), 0)
	l2.Set(uint64(5), 0)
	l2.Set(uint64(9), 0)
	l2.Set(uint64(8), 0)
	l2.Set(uint64(2), 0)
	//skip list 内部会自动做排序，排完序之后为 2 4 5 8 9

	l3 = skiplist.New(skiplist.Uint64)
	l3.Set(uint64(3), 0)
	l3.Set(uint64(5), 0)
	l3.Set(uint64(7), 0)
	l3.Set(uint64(9), 0)
	//skip list 内部会自动做排序，排完序之后为 3 5 7 9
}

func TestIntersectionOfSkipList(t *testing.T) {
	intersection := reverseindex.IntersectionOfSkipList()
	if intersection != nil {
		node := intersection.Front()
		for node != nil {
			fmt.Printf("%d ", node.Key().(uint64))
			node = node.Next()
		}
	}
	fmt.Println("\n" + strings.Repeat("-", 50))

	intersection = reverseindex.IntersectionOfSkipList(l1)
	if intersection != nil {
		node := intersection.Front()
		for node != nil {
			fmt.Printf("%d ", node.Key().(uint64))
			node = node.Next()
		}
	}
	fmt.Println("\n" + strings.Repeat("-", 50))

	intersection = reverseindex.IntersectionOfSkipList(l1, l2)
	if intersection != nil {
		node := intersection.Front()
		for node != nil {
			fmt.Printf("%d ", node.Key().(uint64))
			node = node.Next()
		}
	}
	fmt.Println("\n" + strings.Repeat("-", 50))

	intersection = reverseindex.IntersectionOfSkipList(l1, l2, l3)
	if intersection != nil {
		node := intersection.Front()
		for node != nil {
			fmt.Printf("%d ", node.Key().(uint64))
			node = node.Next()
		}
	}
	fmt.Println("\n" + strings.Repeat("-", 50))
}

func TestUnionOfSetSkipList(t *testing.T) {
	union := reverseindex.UnionOfSkipList()
	if union != nil {
		node := union.Front()
		for node != nil {
			fmt.Printf("%d ", node.Key().(uint64))
			node = node.Next()
		}
	}
	fmt.Println("\n" + strings.Repeat("-", 50))

	union = reverseindex.UnionOfSkipList(l1)
	if union != nil {
		node := union.Front()
		for node != nil {
			fmt.Printf("%d ", node.Key().(uint64))
			node = node.Next()
		}
	}
	fmt.Println("\n" + strings.Repeat("-", 50))

	union = reverseindex.UnionOfSkipList(l1, l2)
	if union != nil {
		node := union.Front()
		for node != nil {
			fmt.Printf("%d ", node.Key().(uint64))
			node = node.Next()
		}
	}
	fmt.Println("\n" + strings.Repeat("-", 50))

	union = reverseindex.UnionOfSkipList(l1, l2, l3)
	if union != nil {
		node := union.Front()
		for node != nil {
			fmt.Printf("%d ", node.Key().(uint64))
			node = node.Next()
		}
	}
	fmt.Println("\n" + strings.Repeat("-", 50))
}
