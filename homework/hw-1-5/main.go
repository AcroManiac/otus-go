package main

import (
	"fmt"
	"github.com/AcroManiac/otus-go/homework/hw-1-5/doublylinkedlist"
)

func main() {
	list := doublylinkedlist.List{}
	list.PushBack(2)
	fmt.Printf("%d", list.First().Value())
}
