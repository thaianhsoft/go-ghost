package ghostengine

import (
	"fmt"
	"strings"
)

type node struct {
	index int
	key string
	childNodes []*node
	fullKey string
	isWild bool
}

func newNode() *node {
	return &node{
		index: -1,
	}
}

func (n *node) LenChilds() int {
	return len(n.childNodes)
}
func (n *node) matchChild(key string) (*node, int) {
	if n.childNodes == nil {
		n.childNodes = make([]*node, 0)
		return nil, -1
	}
		for i, child := range n.childNodes {
			if child.key == key {
				return child, i
			}
		}
	return nil, -1
}

func (n *node) matchChildStartFrom(index int, key string) (*node, int) {
	if index < len(n.childNodes) {
		for i := index; i < len(n.childNodes); i++ {
			if n.childNodes[i].key == key || n.childNodes[i].isWild {
				return n.childNodes[i], i
			}
		}
	}
	return nil, -1
}

func (n *node) showChilds() []node {
	shows := make([]node, len(n.childNodes))
	for i := range n.childNodes {
		shows[i] = *n.childNodes[i]
	}
	return shows
}

func (n *node) GetChildIth(index int) *node {
	if index < len(n.childNodes) {
		return n.childNodes[index]
	}
	return nil
}


type prefixTrie struct {
	root *node
}

func split(key string, sep string) []string {
	pt := strings.Split(key, "/")
	if pt[0] == "" {
		return pt[1:]
	}
	return pt
}

func (p *prefixTrie) insertKey(fullKey string, keys []string) {
	n := p.root
 	for i := range keys {
		child, _ := n.matchChild(keys[i])
		if child == nil {
			child = &node{
				key: keys[i],
				isWild: string(keys[i][0]) == ":",
			}
			n.childNodes = append(n.childNodes, child)
		}
		fmt.Println("parent: ", n.key, "child: ", child)
		n = child
	}
	n.fullKey = fullKey
}

func (p *prefixTrie) showTree(n *node, v *string) string {
	for i, child := range n.childNodes {
		if i == 0 {
			*v += fmt.Sprintf("parent key: %v, childs: [", n.key)
		}
		if i > 0 {
			*v += ", "
		}
		*v += child.key
		*v += "]\n"
		p.showTree(child, v)
	}
	return *v
}

func (p *prefixTrie) search(url string) []string {
	st := &st{}
	dummyNode := &pair{
		node: p.root,
		indexNB: 0,
	}
	splitPtr := split(url, "/")
	indexPtr := 0
	for {
		if child, indexMatch := dummyNode.matchChildStartFrom(dummyNode.indexNB, splitPtr[indexPtr]); child != nil {
			dummyNode.indexNB = indexMatch
			fmt.Println("next child", child.key, dummyNode.showChilds())
			if indexPtr == len(splitPtr) - 1 {
				if child.fullKey != "" {
					return splitPtr
				}
			} else {
				st.PushBack(dummyNode)
				st.PushBack(dummyNode)
				dummyNode = &pair{
					node: child,
					indexNB: 0,
				}
				splitPtr[indexPtr] = child.key
				indexPtr++
			}
		}
		popNode := st.Pop().(*pair)
		if !st.Empty() && st.Peak().(*pair) == popNode {
			// backtrack from neighbor back to parent
			if  popNode.indexNB < popNode.LenChilds() {
				// if child is not final
				st.PushBack(popNode)
				st.PushBack(popNode)
				nb := popNode.GetChildIth(popNode.indexNB)
				dummyNode = &pair{
					node: nb,
					indexNB: 0,
				}
			} else {
				fmt.Println("have backtrack decrease one")
				indexPtr-- // backtrack decrease one
			}
			popNode.indexNB++
		}
	}
}


type internalRouter struct {
	handlers map[string]HandlerFunc
	tries map[string]*prefixTrie
}

func newRouter() *internalRouter {
	return &internalRouter{
		handlers: make(map[string]HandlerFunc),
		tries: make(map[string]*prefixTrie),
	}
}

func (i *internalRouter) insertRoute(methodTrie string, url string, handlerFunc HandlerFunc) {
	splitPtr := split(url, "/")
	pattern := methodTrie + "-" + url
	if _, ok := i.tries[methodTrie]; !ok {
		i.tries[methodTrie] = &prefixTrie{
			root: &node{key: ""},
		}
	}
	i.tries[methodTrie].insertKey(url, splitPtr)
	i.handlers[pattern] = handlerFunc
}

func (i *internalRouter) getRoute(methodTrie string, url string) (fn HandlerFunc, param map[string]interface{}) {
	if _, ok := i.tries[methodTrie]; ok {
		var arrayPattern []string = i.tries[methodTrie].search(url)
		fmt.Println("have array pattern", arrayPattern)
		if arrayPattern == nil {
			return nil, nil
		}
		var pattern string = methodTrie + "-" + "/" + strings.Join(arrayPattern, "/")
		fmt.Println(pattern, i.handlers)
		if _, ok := i.handlers[pattern]; ok {
			// have handler
			splitPattern := split(url, "/")
			for i := range arrayPattern {
				if param == nil {
					param = make(map[string]interface{})
				}
				if arrayPattern[i] != splitPattern[i] {
					param[arrayPattern[i][1:]] = splitPattern[i]
				}
			}
			return i.handlers[pattern], param
		}
	}
	return nil, nil
}


type pair struct {
	indexNB int
	*node
}

type e struct {
	data interface{}
	down *e
}

type st struct {
	top *e
}

func (s *st) Pop() interface{} {
	popTop := s.top
	s.top = s.top.down
	return popTop.data
}

func (s *st) PushBack(data interface{}) {
	n := e{data: data}
	if s.top == nil {
		s.top = &n
	} else {
		s.top, n.down = &n, s.top
	}
}

func (s *st) Empty() bool {
	return s.top == nil
}

func (s *st) Peak() interface{} {
	return s.top.data
}

func (s *st) PeakNextTop() interface{} {
	return s.top.down.data
}