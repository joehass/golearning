package minKoa

type handleFunc func()

type chainFunc []handleFunc

type Engine struct {
	Mid   chainFunc
	index int
}

func New() *Engine {
	return &Engine{
		Mid:   make(chainFunc, 0),
		index: -1,
	}
}

func (e *Engine) Use(f handleFunc) {
	e.Mid = append(e.Mid, f)
}

func (e *Engine) Next() {
	e.index++
	for e.index < len(e.Mid) {
		e.Mid[e.index]()
		e.index++
	}
}

func (e *Engine) reset() {
	e.index = -1
	e.Mid = e.Mid[0:0]
}

func (e *Engine) Run() {
	e.Next()
}
