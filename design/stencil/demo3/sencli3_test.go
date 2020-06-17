package demo3

import (
	"fmt"
	"testing"
)

//模版模式

type docSuper struct {
	getContent func() string
}

func (d docSuper) doOperate() {
	fmt.Println("对这个文档做了一些处理，文档是：", d.getContent())
}

type localDoc struct {
	docSuper
}

func newLocalDoc() *localDoc {
	c := new(localDoc)
	c.docSuper.getContent = c.getContent
	return c
}

func (c *localDoc) getContent() string {
	return "this is a localDoc"
}

type netDoc struct {
	docSuper
}

func newNetDoc() *netDoc {
	c := new(netDoc)
	c.docSuper.getContent = c.getContent
	return c
}

func (c *netDoc) getContent() string {
	return "this is a netDoc"
}

func TestDoc(t *testing.T) {
	netDoc := newNetDoc()
	lcDoc := newLocalDoc()

	netDoc.doOperate()
	lcDoc.doOperate()
}
