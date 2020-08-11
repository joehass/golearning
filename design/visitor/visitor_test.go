package visitor

import "testing"

func TestV1(t *testing.T) {
	e := new(Element)
	e.Accept(new(ProductionVisitor))
	e.Accept(new(TestingVisitor))

	m := new(EnvExample)
	m.Print(new(ProductionVisitor))
	m.Print(new(TestingVisitor))
}
