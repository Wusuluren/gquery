package gquery

import (
	"fmt"
	"strings"
)

var _ = fmt.Println

type HtmlNode struct {
	label    string
	attr     map[string]string
	text     string
	value    string
	html     string
	parent   *HtmlNode
	children []*HtmlNode
	isPaired bool
}

type HtmlNodes []*HtmlNode

func (hn *HtmlNode) isFitSelector(selector string) bool {
	var idName, className, labelName string
	var tmp []string
	if selector == "*" {
		return true
	}
	//match attr
	if selector[0] == '[' && selector[len(selector)-1] == ']' {
		selector = selector[1 : len(selector)-1]
		tmp = strings.Split(selector, "=")
		if len(tmp) > 1 {
			attrName := tmp[0]
			attrValue := tmp[1]
			return reStrCmp(StrStrMap(hn.attr).Get(attrName), attrValue)
		}
		return StrStrMap(hn.attr).Get(selector) != ""
	}
	//match class
	tmp = strings.Split(selector, ".")
	if len(tmp) > 1 {
		labelName = tmp[0]
		className = tmp[1]
		return hn.label == labelName && reStrCmp(StrStrMap(hn.attr).Get("class"), className)
	}
	//match id
	tmp = strings.Split(selector, "#")
	if len(tmp) > 1 {
		labelName = tmp[0]
		idName = tmp[1]
		return hn.label == labelName && reStrCmp(StrStrMap(hn.attr).Get("id"), idName)
	}
	return hn.label == selector
}

func (hn *HtmlNode) isFitNode(node *HtmlNode) bool {
	fit := hn.label == node.label
	if !fit {
		return false
	}
	//fit = StrStrMap(hn.attr).Get("id") ==  StrStrMap(node.attr).Get("id")
	//if !fit {
	//	return false
	//}
	//fit = StrStrMap(hn.attr).Get("class") ==  StrStrMap(node.attr).Get("class")
	//if !fit {
	//	return false
	//}
	return true
}

func (hn *HtmlNode) Gquery(selector string) []*HtmlNode {
	return hn.Find(selector)
}

func (hn *HtmlNode) Parent() *HtmlNode {
	return hn.parent
}

func (hn *HtmlNode) Parents() []*HtmlNode {
	parents := make([]*HtmlNode, 0)
	for parent := hn.parent; parent != nil; parent = parent.parent {
		parents = append(parents, parent)
	}
	return parents
}

func (hn *HtmlNode) ParentsUntil(selector string) []*HtmlNode {
	parents := make([]*HtmlNode, 0)
	for parent := hn.parent; parent != nil; parent = parent.parent {
		if !hn.isFitSelector(selector) {
			break
		}
		parents = append(parents, parent)
	}
	return parents
}

func (hn *HtmlNode) Children(selector string) []*HtmlNode {
	children := make([]*HtmlNode, 0)
	if selector == "*" {
		children = append(children, hn.children...)
	} else {
		for _, node := range hn.children {
			if hn.isFitSelector(selector) {
				children = append(children, node)
			}
		}
	}
	return children
}

func (hn *HtmlNode) Find(selector string) []*HtmlNode {
	children := make([]*HtmlNode, 0)
	if selector == "*" {
		children = append(children, hn.children...)
	} else {
		for _, node := range hn.children {
			if hn.isFitSelector(selector) {
				children = append(children, node)
			}
			children = append(children, node.Find(selector)...)
		}
	}
	return children
}

func (hn *HtmlNode) Siblings(selector string) []*HtmlNode {
	siblings := make([]*HtmlNode, 0)
	if selector == "*" {
		for _, node := range hn.parent.children {
			if node != hn {
				siblings = append(siblings, node)
			}
		}
	} else {
		for _, node := range hn.children {
			if hn.isFitSelector(selector) && node != hn {
				siblings = append(siblings, node)
			}
		}
	}
	return siblings
}

func (hn *HtmlNode) Next() *HtmlNode {
	var sibling *HtmlNode
	findSelf := false
	for _, node := range hn.parent.children {
		if findSelf && hn.isFitNode(node) {
			sibling = node
			break
		}
		if node == hn {
			findSelf = true
		}
	}
	return sibling
}

func (hn *HtmlNode) NextAll() []*HtmlNode {
	siblings := make([]*HtmlNode, 0)
	findSelf := false
	for _, node := range hn.parent.children {
		if findSelf && hn.isFitNode(node) {
			siblings = append(siblings, node)
		}
		if node == hn {
			findSelf = true
		}
	}
	return siblings
}

func (hn *HtmlNode) NextUntil(selector string) []*HtmlNode {
	siblings := make([]*HtmlNode, 0)
	findSelf := false
	for _, node := range hn.parent.children {
		if findSelf && hn.isFitSelector(selector) {
			break
		}
		if findSelf && hn.isFitNode(node) {
			siblings = append(siblings, node)
		}
		if node == hn {
			findSelf = true
		}
	}
	return siblings
}

func (hn *HtmlNode) Prev() *HtmlNode {
	var sibling *HtmlNode
	findSelf := false
	for i := len(hn.parent.children) - 1; i >= 0; i-- {
		node := hn.parent.children[i]
		if findSelf && hn.isFitNode(node) {
			sibling = node
			break
		}
		if node == hn {
			findSelf = true
		}
	}
	return sibling
}

func (hn *HtmlNode) PrevAll() []*HtmlNode {
	siblings := make([]*HtmlNode, 0)
	findSelf := false
	for i := len(hn.parent.children) - 1; i >= 0; i-- {
		node := hn.parent.children[i]
		if findSelf && hn.isFitNode(node) {
			siblings = append(siblings, node)
		}
		if node == hn {
			findSelf = true
		}
	}
	return siblings
}

func (hn *HtmlNode) PrevUntil(selector string) []*HtmlNode {
	siblings := make([]*HtmlNode, 0)
	findSelf := false
	for i := len(hn.parent.children) - 1; i >= 0; i-- {
		node := hn.parent.children[i]
		if findSelf && node.isFitSelector(selector) {
			break
		}
		if findSelf && hn.isFitNode(node) {
			siblings = append(siblings, node)
		}
		if node == hn {
			findSelf = true
		}
	}
	return siblings
}

func (hn *HtmlNode) First(selector string) *HtmlNode {
	var child *HtmlNode
	if selector == "*" {
		if len(hn.children) > 0 {
			child = hn.children[0]
		}
	} else {
		for _, node := range hn.children {
			if node.isFitSelector(selector) {
				child = node
				break
			}
		}
	}
	return child
}

func (hn *HtmlNode) Last(selector string) *HtmlNode {
	var child *HtmlNode
	childrenNum := len(hn.children)
	if selector == "*" {
		if childrenNum > 0 {
			child = hn.children[childrenNum-1]
		}
	}
	for i := childrenNum - 1; i >= 0; i-- {
		node := hn.children[i]
		if node.isFitSelector(selector) {
			child = node
			break
		}
	}
	return child
}

func (hn *HtmlNode) Eq(idx int) *HtmlNode {
	var child *HtmlNode
	ctr := 0
	for _, node := range hn.children {
		ctr++
		if ctr >= idx {
			child = node
			break
		}
	}
	return child
}

func (hn HtmlNodes) Filter(selector string) []*HtmlNode {
	children := make([]*HtmlNode, 0)
	if selector == "*" {
		children = append(children, hn...)
	} else {
		for _, node := range hn {
			if node.isFitSelector(selector) {
				children = append(children, node)
			}
		}
	}
	return children
}

func (hn HtmlNodes) Not(selector string) []*HtmlNode {
	children := make([]*HtmlNode, 0)
	if selector != "*" {
		for _, node := range hn {
			if !node.isFitSelector(selector) {
				children = append(children, node)
			}
		}
	}
	return children
}

func (hn *HtmlNode) Text() string {
	return hn.text
}

func (hn *HtmlNode) Html() string {
	return hn.html
}

func (hn *HtmlNode) Value() string {
	return hn.value
}

func (hn *HtmlNode) Attr(selector string) string {
	return StrStrMap(hn.attr).Get(selector)
}

func (hn *HtmlNode) Append(node *HtmlNode) {
	if hn.children == nil {
		hn.children = make([]*HtmlNode, 0)
	}
	node.parent = hn
	hn.children = append(hn.children, node)
}

func (hn *HtmlNode) Prepend(node *HtmlNode) {
	if hn.children == nil {
		hn.children = make([]*HtmlNode, 0)
	}
	node.parent = hn
	children2 := make([]*HtmlNode, 0)
	children2 = append(children2, node)
	children2 = append(children2, hn.children...)
	hn.children = children2
}

func (hn *HtmlNode) After(node *HtmlNode) {
	idx := 0
	for i, node := range hn.parent.children {
		if node == hn {
			idx = i
			break
		}
	}
	children2 := hn.parent.children[0 : idx+1]
	children2 = append(children2, node)
	children2 = append(children2, hn.parent.children[idx+1:]...)
	hn.parent.children = children2
}

func (hn *HtmlNode) Before(node *HtmlNode) {
	idx := 0
	for i, node := range hn.parent.children {
		if node == hn {
			idx = i
			break
		}
	}
	children2 := hn.parent.children[0:idx]
	children2 = append(children2, node)
	children2 = append(children2, hn.parent.children[idx:]...)
	hn.parent.children = children2
}

func (hn *HtmlNode) Remove() {
	idx := 0
	for i, node := range hn.parent.children {
		if node == hn {
			idx = i
			break
		}
	}
	children2 := hn.parent.children[idx+1:]
	hn.parent.children = hn.parent.children[0:idx]
	hn.parent.children = append(hn.parent.children, children2...)
}

func (hn *HtmlNode) Empty() {
	hn.children = hn.children[0:0]
}
