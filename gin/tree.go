package gin

import (
	"bytes"
	"golearning/gin/internal/bytesconv"
	"net/url"
)

var (
	strColon = []byte(":")
	strStar  = []byte("*")
)

type nodeType uint8

const (
	static   nodeType = iota //普通节点
	root                     //根节点
	param                    //参数路由，比如 /user/:id
	catchAll                 //匹配所有内容的路由，比如 /article/*key
)

//其中 path和indices是使用了前缀书的逻辑
type node struct {
	path      string        //当前节点相对路径（与祖先节点的path拼接可得到完整路径）
	indices   string        //孩子节点的path[0]组成的字符串
	priority  uint32        //当前节点及子孙节点的实际路由数量
	children  []*node       //孩子节点
	handlers  HandlersChain //当前节点的处理函数，包括中间件
	nType     nodeType      //节点类型
	wildChild bool          //孩子节点是否有通配符
	fullPath  string        //全路径
}

type nodeValue struct {
	handlers HandlersChain
	params   *Params
	tsr      bool
	fullPath string
}

//返回给定key注册过的handle
func (n *node) getValue(path string, params *Params, unescape bool) (value nodeValue) {
walk: // Outer loop for walking the tree
	for {
		prefix := n.path
		if len(path) > len(prefix) {
			if path[:len(prefix)] == prefix {
				path = path[len(prefix):]
				// If this node does not have a wildcard (param or catchAll)
				// child, we can just look up the next child node and continue
				// to walk down the tree
				if !n.wildChild {
					idxc := path[0]
					for i, c := range []byte(n.indices) {
						if c == idxc {
							n = n.children[i]
							continue walk
						}
					}

					// Nothing found.
					// We can recommend to redirect to the same URL without a
					// trailing slash if a leaf exists for that path.
					value.tsr = (path == "/" && n.handlers != nil)
					return
				}

				// Handle wildcard child
				n = n.children[0]
				switch n.nType {
				case param:
					// Find param end (either '/' or path end)
					end := 0
					for end < len(path) && path[end] != '/' {
						end++
					}

					// Save param value
					if params != nil {
						if value.params == nil {
							value.params = params
						}
						// Expand slice within preallocated capacity
						i := len(*value.params)
						*value.params = (*value.params)[:i+1]
						val := path[:end]
						if unescape {
							if v, err := url.QueryUnescape(val); err == nil {
								val = v
							}
						}
						(*value.params)[i] = Param{
							Key:   n.path[1:],
							Value: val,
						}
					}

					// we need to go deeper!
					if end < len(path) {
						if len(n.children) > 0 {
							path = path[end:]
							n = n.children[0]
							continue walk
						}

						// ... but we can't
						value.tsr = (len(path) == end+1)
						return
					}

					if value.handlers = n.handlers; value.handlers != nil {
						value.fullPath = n.fullPath
						return
					}
					if len(n.children) == 1 {
						// No handle found. Check if a handle for this path + a
						// trailing slash exists for TSR recommendation
						n = n.children[0]
						value.tsr = (n.path == "/" && n.handlers != nil)
					}
					return

				case catchAll:
					// Save param value
					if params != nil {
						if value.params == nil {
							value.params = params
						}
						// Expand slice within preallocated capacity
						i := len(*value.params)
						*value.params = (*value.params)[:i+1]
						val := path
						if unescape {
							if v, err := url.QueryUnescape(path); err == nil {
								val = v
							}
						}
						(*value.params)[i] = Param{
							Key:   n.path[2:],
							Value: val,
						}
					}

					value.handlers = n.handlers
					value.fullPath = n.fullPath
					return

				default:
					panic("invalid node type")
				}
			}
		}

		if path == prefix {
			// We should have reached the node containing the handle.
			// Check if this node has a handle registered.
			if value.handlers = n.handlers; value.handlers != nil {
				value.fullPath = n.fullPath
				return
			}

			// If there is no handle for this route, but this route has a
			// wildcard child, there must be a handle for this path with an
			// additional trailing slash
			if path == "/" && n.wildChild && n.nType != root {
				value.tsr = true
				return
			}

			// No handle found. Check if a handle for this path + a
			// trailing slash exists for trailing slash recommendation
			for i, c := range []byte(n.indices) {
				if c == '/' {
					n = n.children[i]
					value.tsr = (len(n.path) == 1 && n.handlers != nil) ||
						(n.nType == catchAll && n.children[0].handlers != nil)
					return
				}
			}

			return
		}

		// Nothing found. We can recommend to redirect to the same URL with an
		// extra trailing slash if a leaf exists for that path
		value.tsr = (path == "/") ||
			(len(prefix) == len(path)+1 && prefix[len(path)] == '/' &&
				path == prefix[:len(prefix)-1] && n.handlers != nil)
		return
	}

}

//前缀树
func (n *node) addRoute(path string, handlers HandlersChain) {
	fullPath := path
	n.priority++

	//空树
	if len(n.path) == 0 && len(n.children) == 0 {
		n.insertChild(path, fullPath, handlers)
		n.nType = root
		return
	}

	parentFullPathIndex := 0

walk:
	for {
		//找到最长公共前缀
		//这意味着公共前缀不包含:和*
		//因为现有的key不能包含这些字符
		i := longestCommonPrefix(path, n.path)

		//将路径短的作为父节点，路径长的作为子节点
		if i < len(n.path) {
			child := node{
				path:      n.path[i:],
				indices:   n.indices,
				children:  n.children,
				handlers:  n.handlers,
				wildChild: n.wildChild,
				fullPath:  n.fullPath,
				priority:  n.priority - 1,
			}

			n.children = []*node{&child}
			n.indices = bytesconv.BytesToString([]byte{n.path[i]})
			n.path = path[:i]
			n.handlers = nil
			n.wildChild = false
			n.fullPath = fullPath[:parentFullPathIndex+i]
		}

		//创建一个新的node children
		if i < len(path) {
			path = path[i:]

			//if n.wildChild {
			//	parentFullPathIndex += len(n.path)
			//	n = n.children[0]
			//	n.priority++
			//
			//	if len(path) >= len(n.path) && n.path
			//}

			c := path[0]

			//TODO：插入到param后

			//TODO：检查是否存入下一个孩子节点
			for i, max := 0, len(n.indices); i < max; i++ {
				//前缀相等
				if c == n.indices[i] {
					parentFullPathIndex += len(n.path)
					i = n.incrementChildPrio(i)
					n = n.children[i]
					continue walk
				}
			}

			if c != ':' && c != '*' {
				n.indices += bytesconv.BytesToString([]byte{c})
				child := &node{
					fullPath: fullPath,
				}
				n.children = append(n.children, child)
				n.incrementChildPrio(len(n.indices) - 1)
				n = child
			}
			n.insertChild(path, fullPath, handlers)
			return
		}
		n.handlers = handlers
		n.fullPath = fullPath
		return
	}

}

//增加指定孩子的优先级，并在必要时排序
func (n *node) incrementChildPrio(pos int) int {
	cs := n.children
	cs[pos].priority++
	prio := cs[pos].priority

	//调整顺序（移动到前面）,按优先级从大到小排序
	newPos := pos
	for ; newPos > 0 && cs[newPos-1].priority < prio; newPos-- {
		cs[newPos-1], cs[newPos] = cs[newPos], cs[newPos-1]
	}

	if newPos != pos {
		n.indices = n.indices[:newPos] +
			n.indices[pos:pos+1] +
			n.indices[newPos:pos] + n.indices[pos+1:]
	}

	return newPos

}

//增加通配符的孩子节点
func (n *node) insertChild(path string, fullPath string, handlers HandlersChain) {
	for {
		wildcard, i, valid := findWildcard(path)
		if i < 0 { //没找到通配符
			break
		}

		//通配符不能包含:和*
		if !valid {
			panic("only one wildcard per path segment is allowed, has: '" +
				wildcard + "' in path '" + fullPath + "'")
		}

		//通配符必须匹配一个变量
		if len(wildcard) < 2 {
			panic("wildcards must be named with a non-empty name in path '" + fullPath + "'")
		}

		if len(n.children) > 0 {
			panic("wildcard segment '" + wildcard +
				"' conflicts with existing children in path '" + fullPath + "'")
		}

		//这里开始判断参数路由，即:路由
		if wildcard[0] == ':' {
			if i > 0 {
				//插入通配符前缀
				n.path = path[:i]
				path = path[i:]
			}

			n.wildChild = true
			child := &node{
				nType:    param,
				path:     wildcard,
				fullPath: fullPath,
			}
			n.children = []*node{child}
			n = child
			n.priority++

			//TODO: 判断第二个通配符
			if len(wildcard) < len(path) {

			}

			n.handlers = handlers
			return
		}

		//这里开始判断匹配所有的路由，即*路由
		if i+len(wildcard) != len(path) {
			panic("catch-all conflicts with existing handle for the path segment root in path '" + fullPath + "'")
		}

		i--
		if path[i] != '/' {
			panic("no / before catch-all in path '" + fullPath + "'")
		}

		n.path = path[:i]

		child := &node{
			wildChild: true,
			nType:     catchAll,
			fullPath:  fullPath,
		}

		n.children = []*node{child}
		n.indices = string('/')
		n = child
		n.priority++

		//第二个节点，保持变量
		child = &node{
			path:     path[i:],
			nType:    catchAll,
			handlers: handlers,
			priority: 1,
			fullPath: fullPath,
		}
		n.children = []*node{child}
		return
	}

	//没有找到通配符
	n.path = path
	n.handlers = handlers
	n.fullPath = fullPath
}

func countParams(path string) uint16 {
	var n uint16
	s := bytesconv.StringToBytes(path)
	n += uint16(bytes.Count(s, strColon))
	n += uint16(bytes.Count(s, strStar))
	return n
}

func longestCommonPrefix(a, b string) int {
	i := 0
	max := min(len(a), len(b))
	for i < max && a[i] == b[i] {
		i++
	}

	return i
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

//检索通配符，并检查是否有无效字符，如果没有通配符，则返回-1
//例如/foo:id/ 则返回:id，:和*通配符只能在/ /之间出现一次
func findWildcard(path string) (wildcard string, i int, valid bool) {
	for start, c := range []byte(path) {
		if c != ':' && c != '*' {
			continue
		}

		valid = true
		for end, c := range []byte(path[start+1:]) {
			switch c {
			case '/':
				return path[start : start+1+end], start, valid
			case ':', '*':
				valid = false
			}
			return path[start:], start, valid
		}
	}
	return "", -1, false
}

type methodTree struct {
	method string
	root   *node
}

type methodTrees []methodTree

func (trees methodTrees) get(method string) *node {
	for _, tree := range trees {
		if tree.method == method {
			return tree.root
		}
	}
	return nil
}

//单个url参数，包含key和value
type Param struct {
	Key   string
	Value string
}

//路由器返回的参数切片
type Params []Param

func (ps Params) Get(name string) (string, bool) {
	for _, entry := range ps {
		if entry.Key == name {
			return entry.Value, true
		}
	}

	return "", false
}

//根据给定的name在url中进行匹配，返回第一个匹配到的数据
//如果没找到，则返回空字符串
func (ps Params) ByName(name string) (va string) {
	va, _ = ps.Get(name)
	return
}
