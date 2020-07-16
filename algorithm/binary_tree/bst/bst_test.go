package bst

import (
	"fmt"
	"testing"
)

//二叉查找树
//某个节点的左子树的所有节点的值都比这个节点的值小
//某个节点的右子树的所有节点的值都比这个节点值大

type BiSearchTree struct {
	Data   int
	Lchild *BiSearchTree
	Rchild *BiSearchTree
}

//创建一颗新树
func NewBiSearchTree(data int) *BiSearchTree {
	return &BiSearchTree{Data: data}
}

//添加元素
func (bst *BiSearchTree) Add(data int) *BiSearchTree {
	if bst == nil {
		return NewBiSearchTree(data)
	}

	if data < bst.Data {
		bst.Lchild = bst.Lchild.Add(data)
	} else {
		bst.Rchild = bst.Rchild.Add(data)
	}

	return bst
}

//是否包含某一个元素
func (bst *BiSearchTree) Contains(data int) bool {
	if bst == nil {
		return false
	}

	v := bst.compareTo(data)

	if v < 0 {
		return bst.Lchild.Contains(data)
	} else if v > 0 {
		return bst.Rchild.Contains(data)
	} else {
		return true
	}
}

func (bst *BiSearchTree) compareTo(data int) int {
	return data - bst.Data
}

//移除元素,有3种情况
//情况 1：如果删除的节点没有右孩子，那么就选择它的左孩子来代替原来的节点。二叉查找树的性质保证了被删除节点的左子树必然符合二叉查找树的性质。
//因此左子树的值要么都大于，要么都小于被删除节点的父节点的值，这取决于被删除节点是左孩子还是右孩子。
//因此用被删除节点的左子树来替代被删除节点，是完全符合二叉搜索树的性质的。
//情况 2：如果被删除节点的右孩子没有左孩子，那么这个右孩子被用来替换被删除节点。
//因为被删除节点的右孩子都大于被删除节点左子树的所有节点，同时也大于或小于被删除节点的父节点，
//这同样取决于被删除节点是左孩子还是右孩子。因此，用右孩子来替换被删除节点，符合二叉查找树的性质。
//情况 3：如果被删除节点的右孩子有左孩子，就需要用被删除节点右孩子的左子树中的最下面的节点来替换它，
//就是说，我们用被删除节点的右子树中最小值的节点来替换。
func (bst *BiSearchTree) Remove(data int) *BiSearchTree {
	if bst == nil {
		return bst
	}
	compareResult := bst.compareTo(data)
	if compareResult < 0 {
		bst.Lchild = bst.Lchild.Remove(data)
	} else if compareResult > 0 {
		bst.Rchild = bst.Rchild.Remove(data)
	} else if bst.Lchild != nil && bst.Rchild != nil { //第三种情况
		bst.Data = bst.Rchild.FindMin()
		bst.Rchild = bst.Rchild.Remove(bst.Data)
	} else if bst.Lchild != nil { //第一种情况
		bst = bst.Lchild
	} else { //第二种情况
		bst = bst.Rchild
	}
	return bst
}

//查找最小值
func (bst *BiSearchTree) FindMin() int {
	if bst == nil {
		fmt.Println("tree is empty")
		return -1
	}

	if bst.Lchild == nil {
		return bst.Data
	} else {
		return bst.Lchild.FindMin()
	}
}

//获取树所有的元素值，按从小到大排序
func (bst *BiSearchTree) GetAll() []int {
	values := []int{}
	return appendValue(values, bst)
}

func appendValue(values []int, bst *BiSearchTree) []int {
	if bst != nil {
		values = appendValue(values, bst.Lchild)
		values = append(values, bst.Data)
		values = appendValue(values, bst.Rchild)
	}
	return values
}

func TestBST(t *testing.T) {
	binaryTree := NewBiSearchTree(50)
	binaryTree.Add(20)
	binaryTree.Add(10)
	binaryTree.Add(100)
	binaryTree.Add(60)
	binaryTree.Add(70)
	binaryTree.Add(5)
	binaryTree.Add(35)
	binaryTree.Add(40)
	fmt.Println(binaryTree.GetAll())

	fmt.Println(binaryTree.Contains(30))
	fmt.Println(binaryTree.Contains(20))

	fmt.Println(binaryTree.FindMin())

	binaryTree.Remove(20)
	fmt.Println(binaryTree.GetAll())
}
