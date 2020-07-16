package bt

import (
	"fmt"
	"testing"
)

/**
二叉树遍历
*/

//二叉树
type Node struct {
	data  string
	left  *Node
	right *Node
}

func init() {

}

func TestR(t *testing.T) {
	//生成二叉树
	nodeG := Node{data: "g", left: nil, right: nil}
	nodeF := Node{data: "f", left: &nodeG, right: nil}
	nodeE := Node{data: "e", left: nil, right: nil}
	nodeD := Node{data: "d", left: &nodeE, right: nil}
	nodeC := Node{data: "c", left: nil, right: nil}
	nodeB := Node{data: "b", left: &nodeD, right: &nodeF}
	nodeA := Node{data: "a", left: &nodeB, right: &nodeC}
	preOrderRecursive(nodeA)
	//preOrder(&nodeA)
	//minOrder(&nodeA)
}

func preOrderRecursive(node Node) {

	if node.left != nil {
		preOrderRecursive(*node.left)
	}
	//fmt.Println(node.data)
	// 在这里输出就是中序
	if node.right != nil {
		preOrderRecursive(*node.right)
	}
	// 在这里输出是后序
	fmt.Println(node.data)
}

type seqStack struct {
	data []*Node
	top  int //数组下标
}

//前序
func preOrder(node *Node) {
	var s seqStack
	s.top = -1 //空
	s.data = make([]*Node, 0)
	if node == nil {
		panic("tree is empty")
	} else {
		for node != nil || s.top != -1 {
			if node != nil {
				s.top++
				fmt.Println(node.data)
				s.data = append(s.data, node)
				node = node.left
			} else {
				node = s.data[s.top]
				s.data = s.data[:s.top]
				s.top--
				node = node.right
			}
		}
	}
	return
}

//中序
func minOrder(node *Node) {
	var s seqStack
	s.data = make([]*Node, 0)
	s.top = -1

	if node == nil {
		panic("tree is empty")
	} else {
		for node != nil || s.top != -1 {
			if node != nil {
				s.top++
				s.data = append(s.data, node)
				node = node.left
			} else {
				node = s.data[s.top]
				s.data = s.data[:s.top]
				fmt.Println(node.data)
				s.top--
				node = node.right
			}
		}
	}
}

//后序，暂时写不出来
func postOrder(node *Node) {
	var s seqStack
	s.data = make([]*Node, 0)
	s.top = -1

	if node == nil {
		panic("tree is empty")
	} else {
		for node != nil || s.top != -1 {

		}
	}
}
