package gquery

import (
	"strings"
)

type HtmlNode struct {
	label    string
	id       string
	class    []string
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
	var labelName, idName, className string
	var attrName, attrValue string
	classNames := make([]string, 0)
	attr := make(map[string]string)
	if hn.label == "" {
		return false
	}
	if selector == "*" {
		return true
	}
	type Idx struct {
		idx int
		c   byte
	}
	idxs := make([]Idx, 0)
	for i := 0; i < len(selector); i++ {
		c := selector[i]
		if c == '.' || c == '#' || c == '[' || c == ']' || c == '=' {
			idxs = append(idxs, Idx{
				idx: i,
				c:   c,
			})
		}
	}
	idxsSize := len(idxs)
	if idxsSize > 0 {
		labelName = selector[0:idxs[0].idx]
		for i, idx := range idxs {
			switch idx.c {
			case '.':
				if i+1 < idxsSize {
					className = selector[idx.idx+1 : idxs[i+1].idx]
				} else {
					className = selector[idx.idx+1:]
				}
				classNames = append(classNames, className)
			case '#':
				if idName == "" {
					if i+1 < idxsSize {
						idName = selector[idx.idx+1 : idxs[i+1].idx]
					} else {
						idName = selector[idx.idx+1:]
					}
				}
			case '[':
				if i+1 < idxsSize {
					attrName = selector[idx.idx+1 : idxs[i+1].idx]
				}
				attrValue = ""
			case '=':
				if i+1 < idxsSize {
					attrValue = selector[idx.idx+1 : idxs[i+1].idx]
					attrValue = strings.TrimLeft(attrValue, "'")
					attrValue = strings.TrimRight(attrValue, "'")
				}
			case ']':
				attr[attrName] = attrValue
			}
		}
	} else {
		labelName = selector
	}

	if labelName != "" && hn.label != labelName {
		return false
	}
	if idName != "" && hn.id != idName {
		return false
	}
	if len(classNames) > 0 {
		for _, className := range classNames {
			isFind := false
			for _, class := range hn.class {
				if class == className {
					isFind = true
					break
				}
			}
			if !isFind {
				return false
			}
		}
	}
	for attrName, attrValue := range attr {
		if value, ok := hn.attr[attrName]; ok {
			if attrValue != "" && value != attrValue {
				return false
			}
		} else {
			return false
		}
	}
	return true
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

func (hn *HtmlNode) Gquery(selector string) HtmlNodes {
	children := make([]*HtmlNode, 0)
	if hn.isFitSelector(selector) {
		children = append(children, hn)
	}
	for _, child := range hn.children {
		children = append(children, child.Gquery(selector)...)
	}
	return children
}

func (hn *HtmlNode) Parent() *HtmlNode {
	return hn.parent
}

func (hn *HtmlNode) Parents() HtmlNodes {
	parents := make([]*HtmlNode, 0)
	for parent := hn.parent; parent != nil; parent = parent.parent {
		parents = append(parents, parent)
	}
	return parents
}

func (hn *HtmlNode) ParentsUntil(selector string) HtmlNodes {
	parents := make([]*HtmlNode, 0)
	for parent := hn.parent; parent != nil; parent = parent.parent {
		if !parent.isFitSelector(selector) {
			break
		}
		parents = append(parents, parent)
	}
	return parents
}

func (hn *HtmlNode) Children(selector string) HtmlNodes {
	children := make([]*HtmlNode, 0)
	if selector == "*" {
		children = append(children, hn.children...)
	} else {
		for _, node := range hn.children {
			if node.isFitSelector(selector) {
				children = append(children, node)
			}
		}
	}
	return children
}

func (hn *HtmlNode) Find(selector string) HtmlNodes {
	children := make([]*HtmlNode, 0)
	if selector == "*" {
		children = append(children, hn.children...)
	} else {
		for _, node := range hn.children {
			if node.isFitSelector(selector) {
				children = append(children, node)
			}
			children = append(children, node.Find(selector)...)
		}
	}
	return children
}

func (hn *HtmlNode) Siblings(selector string) HtmlNodes {
	siblings := make([]*HtmlNode, 0)
	if selector == "*" {
		for _, node := range hn.parent.children {
			if node != hn {
				siblings = append(siblings, node)
			}
		}
	} else {
		for _, node := range hn.children {
			if node.isFitSelector(selector) && node != hn {
				siblings = append(siblings, node)
			}
		}
	}
	return siblings
}

func (hn *HtmlNode) Next() *HtmlNode {
	sibling := &HtmlNode{}
	findSelf := false
	for _, node := range hn.parent.children {
		if findSelf && node.isFitNode(hn) {
			sibling = node
			break
		}
		if node == hn {
			findSelf = true
		}
	}
	return sibling
}

func (hn *HtmlNode) NextAll() HtmlNodes {
	siblings := make([]*HtmlNode, 0)
	findSelf := false
	for _, node := range hn.parent.children {
		if findSelf && node.isFitNode(hn) {
			siblings = append(siblings, node)
		}
		if node == hn {
			findSelf = true
		}
	}
	return siblings
}

func (hn *HtmlNode) NextUntil(selector string) HtmlNodes {
	siblings := make([]*HtmlNode, 0)
	findSelf := false
	for _, node := range hn.parent.children {
		if findSelf && hn.isFitSelector(selector) {
			break
		}
		if findSelf && node.isFitNode(hn) {
			siblings = append(siblings, node)
		}
		if node == hn {
			findSelf = true
		}
	}
	return siblings
}

func (hn *HtmlNode) Prev() *HtmlNode {
	sibling := &HtmlNode{}
	findSelf := false
	for i := len(hn.parent.children) - 1; i >= 0; i-- {
		node := hn.parent.children[i]
		if findSelf && node.isFitNode(hn) {
			sibling = node
			break
		}
		if node == hn {
			findSelf = true
		}
	}
	return sibling
}

func (hn *HtmlNode) PrevAll() HtmlNodes {
	siblings := make([]*HtmlNode, 0)
	findSelf := false
	for i := len(hn.parent.children) - 1; i >= 0; i-- {
		node := hn.parent.children[i]
		if findSelf && node.isFitNode(hn) {
			siblings = append(siblings, node)
		}
		if node == hn {
			findSelf = true
		}
	}
	return siblings
}

func (hn *HtmlNode) PrevUntil(selector string) HtmlNodes {
	siblings := make([]*HtmlNode, 0)
	findSelf := false
	for i := len(hn.parent.children) - 1; i >= 0; i-- {
		node := hn.parent.children[i]
		if findSelf && node.isFitSelector(selector) {
			break
		}
		if findSelf && node.isFitNode(hn) {
			siblings = append(siblings, node)
		}
		if node == hn {
			findSelf = true
		}
	}
	return siblings
}

func (hn *HtmlNode) First(selector string) *HtmlNode {
	child := &HtmlNode{}
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
	child := &HtmlNode{}
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

func (hn HtmlNodes) Eq(idx int) *HtmlNode {
	child := &HtmlNode{}
	for i, node := range hn {
		if i >= idx {
			child = node
			break
		}
	}
	return child
}

func (hn HtmlNodes) Filter(selector string) HtmlNodes {
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

func (hn HtmlNodes) Not(selector string) HtmlNodes {
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

func (hn *HtmlNode) SetText(str string) {
	hn.text = str
}

func (hn *HtmlNode) SetHtml(str string) {
	hn.html = str
}

func (hn *HtmlNode) SetValue(str string) {
	hn.value = str
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
