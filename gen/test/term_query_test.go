package test

import (
	"fmt"
	pb "github.com/hjrbill/quicker/gen"
	"testing"
)

const FIELD = ""

func TestTermQuery(t *testing.T) {
	A := pb.NewTermQuery(FIELD, "") //基础 Expression
	B := pb.NewTermQuery(FIELD, "B")
	C := pb.NewTermQuery(FIELD, "C")
	D := pb.NewTermQuery(FIELD, "D")
	E := &pb.TermQuery{} //空 Term
	F := pb.NewTermQuery(FIELD, "F")
	G := pb.NewTermQuery(FIELD, "G")
	H := pb.NewTermQuery(FIELD, "H")

	var q *pb.TermQuery

	fmt.Println(1)
	q = A
	//手动调用 ToString()，避免 println 函数会自动调用 String() 方法
	fmt.Println(q.ToString())

	fmt.Println(2)
	q = B.Or(C)
	fmt.Println(q.ToString())

	fmt.Println(3)
	// (A||B||C)&& D)||E&& ((F||G)&& H
	q = A.Or(B).Or(C).And(D).Or(E).And(F.Or(G)).And(H)
	fmt.Println(q.ToString())
}
