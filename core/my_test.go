package core

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Element struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *Element

	// The list to which this element belongs.
	list *List

	// The value stored with this element.
	Value interface{}
}

type List struct {
	root Element // sentinel list element, only &root, root.prev, and root.next are used
	len  int     // current list length excluding (this) sentinel element
}

func Test(t *testing.T) {
	//	list := new(List)
	//	fmt.Println(list)
	//	fmt.Println(*(&list))
	//	fmt.Println(list.root.prev)
	//	fmt.Println(list.root)
	//	list.root.prev = &list.root
	//	list.root.next = &list.root
	//	fmt.Println(list.root)

	var ps ProcessTask
	ps.Path = "aaa"
	sts := make([]SinkTask, 0, 1)
	var st SinkTask
	st.DataFilePath = "aa"
	st.DestFilePath = "bb"
	sts = append(sts, st)
	ps.SinkTasks = sts

	bb, _ := json.Marshal(ps)
	fmt.Println(string(bb))
}
