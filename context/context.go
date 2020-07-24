package context

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"
)

//context 实现

type Context interface {
	DeadLine() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}

var Canceled = errors.New("context canceled")

var DeadlineExceeded error = deadlineExceededError{}

type deadlineExceededError struct{}

func (deadlineExceededError) Error() string {
	return "context deadline exceeded"
}

type emptyCtx int

func (*emptyCtx) DeadLine() (deadline time.Time, ok bool) {
	return
}

func (*emptyCtx) Done() <-chan struct{} {
	return nil
}

func (*emptyCtx) Err() error {
	return nil
}

func (*emptyCtx) Value(key interface{}) interface{} {
	return nil
}

func (e *emptyCtx) String() string {
	switch e {
	case background:
		return "context.Background"
	case todo:
		return "context.TODO"
	}
	return "unknown empty Context"
}

var (
	background = new(emptyCtx)
	todo       = new(emptyCtx)
)

func Background() Context {
	return background
}

func TODO() Context {
	return todo
}

var closedchan = make(chan struct{})

func init() {
	close(closedchan)
}

type canceler interface {
	cancel(removeFromParent bool, err error)
	Done() <-chan struct{}
}

type cancelCtx struct {
	Context

	mu       sync.Mutex
	done     chan struct{}
	children map[canceler]struct{}
	err      error
}

func (c *cancelCtx) Value(key interface{}) interface{} {
	if key == &cancelCtxKey {
		return c
	}

	return c.Context.Value(key)
}

func (c *cancelCtx) Done() <-chan struct{} {
	c.mu.Lock()
	if c.done == nil {
		c.done = make(chan struct{})
	}
	d := c.done
	c.mu.Unlock()
	return d
}

func (c *cancelCtx) Err() error {
	c.mu.Lock()
	err := c.err
	c.mu.Unlock()
	return err
}

type CancelFunc func()

func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {
	c := newCancelCtx(parent)
	propagateCancel(parent, &c)
	return &c, func() { c.cancel(true, Canceled) }
}

func newCancelCtx(parent Context) cancelCtx {
	return cancelCtx{Context: parent}
}

func propagateCancel(parent Context, child canceler) {
	if parent.Done() == nil {
		return
	}

	if p, ok := parentCancelCtx(parent); ok {
		p.mu.Lock()
		if p.err != nil { //context已经停止，需要停止子context
			child.cancel(false, p.err)
		} else {
			if p.children == nil {
				p.children = make(map[canceler]struct{})
			}
			p.children[child] = struct{}{}
		}
		p.mu.Unlock()
	}

	return
}

//关闭cancel Context 及其子onctxt
func (c *cancelCtx) cancel(removeFromParent bool, err error) {
	if err == nil {
		panic("context:internal error:missing cancel error")
	}
	c.mu.Lock()
	if c.err != nil {
		c.mu.Unlock()
		return
	}

	c.err = err
	if c.done == nil {
		c.done = closedchan
	} else {
		close(c.done)
	}

	for child := range c.children {
		child.cancel(false, err)
	}

	c.children = nil
	c.mu.Unlock()

	if removeFromParent {
		removeChild(c.Context, c)
	}

	return
}

var cancelCtxKey int

//判断context是否是cancelCtx，并判断是否关闭
func parentCancelCtx(parent Context) (*cancelCtx, bool) {
	done := parent.Done()
	if done == closedchan || done == nil {
		return nil, false
	}
	p, ok := parent.Value(&cancelCtxKey).(*cancelCtx)
	if !ok {
		return nil, false
	}
	p.mu.Lock()
	ok = p.done == done
	p.mu.Unlock()
	if !ok {
		return nil, false
	}
	return p, true
}

//把父context的子context删除
func removeChild(parent Context, child canceler) {
	p, ok := parentCancelCtx(parent)
	if !ok {
		return
	}
	fmt.Println(p)
	return
}

func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
	if cur, ok := parent.DeadLine(); ok && cur.Before(d) {
		return WithCancel(parent)
	}

	c := &timerCtx{
		cancelCtx: newCancelCtx(parent),
		deadline:  d,
	}

	propagateCancel(parent, c)
	dur := time.Until(d)
	if dur <= 0 {
		c.cancel(true, DeadlineExceeded)
		return c, func() {
			c.cancel(false, Canceled)
		}
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err == nil {
		c.timer = time.AfterFunc(dur, func() {
			c.cancel(true, DeadlineExceeded)
		})
	}

	return c, func() {
		c.cancel(true, Canceled)
	}
}

type stringer interface {
	String() string
}

func contextName(c Context) string {
	if s, ok := c.(stringer); ok {
		return s.String()
	}

	return reflect.TypeOf(c).String()
}

func (c *cancelCtx) String() string {
	return contextName(c.Context) + ".WithCancel"
}

type timerCtx struct {
	cancelCtx
	timer *time.Timer

	deadline time.Time
}

func (c *timerCtx) Deadline() (deadline time.Time, ok bool) {
	return c.deadline, true
}

func (c *timerCtx) String() string {
	return contextName(c.cancelCtx.Context) + ".WithDeadline(" +
		c.deadline.String() + " [" +
		time.Until(c.deadline).String() + "])"
}

func (c *timerCtx) cancel(removeFromParent bool, err error) {
	c.cancelCtx.cancel(false, err)
	if removeFromParent {
		removeChild(c.cancelCtx.Context, c)
	}

	c.mu.Lock()
	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}

	c.mu.Unlock()
}

func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	return WithDeadline(parent, time.Now().Add(timeout))
}

func WithValue(parent Context, key, val interface{}) Context {
	if key == nil {
		panic("nil key")
	}

	if !reflect.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}

	return &valueCtx{parent, key, val}
}

type valueCtx struct {
	Context
	key, val interface{}
}

func stringify(v interface{}) string {
	switch s := v.(type) {
	case stringer:
		return s.String()
	case string:
		return s
	}
	return "<not Stringers>"
}

func (c *valueCtx) String() string {
	return contextName(c.Context) + ".WithValue(type " +
		reflect.TypeOf(c.key).String() +
		", val" + stringify(c.val) + ")"
}

func (c *valueCtx) Value(key interface{}) interface{} {
	if c.key == key {
		return c.val
	}

	return c.Context.Value(key)
}
