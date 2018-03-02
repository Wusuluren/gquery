package gquery

type MarkdownNode struct {
	_type    int
	attr     map[string]string
	text     string
	html     string
	value    string
	parent   *MarkdownNode
	children []*MarkdownNode
	tabNum   int
}

type MarkdownNodes []*MarkdownNode

func (md *MarkdownNode) isFitType(Type int) bool {
	return md._type == Type
}

func (md *MarkdownNode) isFitNode(node *MarkdownNode) bool {
	return md._type == node._type
}

func (md *MarkdownNode) Gquery(Type int) []*MarkdownNode {
	return md.Find(Type)
}

func (md *MarkdownNode) Parent() *MarkdownNode {
	return md.parent
}

func (md *MarkdownNode) Parents() []*MarkdownNode {
	parents := make([]*MarkdownNode, 0)
	for parent := md.parent; parent != nil; parent = parent.parent {
		parents = append(parents, parent)
	}
	return parents
}

func (md *MarkdownNode) ParentsUntil(Type int) []*MarkdownNode {
	parents := make([]*MarkdownNode, 0)
	for parent := md.parent; parent != nil; parent = parent.parent {
		if Type != MdAll && parent._type != Type {
			break
		}
		parents = append(parents, parent)
	}
	return parents
}

func (md *MarkdownNode) Children(Type int) []*MarkdownNode {
	children := make([]*MarkdownNode, 0)
	if Type == MdAll {
		children = append(children, md.children...)
	} else {
		for _, node := range md.children {
			if node.isFitType(Type) {
				children = append(children, node)
			}
		}
	}
	return children
}

func (md *MarkdownNode) Find(Type int) []*MarkdownNode {
	children := make([]*MarkdownNode, 0)
	if Type == MdAll {
		children = append(children, md.children...)
	} else {
		for _, node := range md.children {
			if node.isFitType(Type) {
				children = append(children, node)
			}
			children = append(children, node.Find(Type)...)
		}
	}
	return children
}

func (md *MarkdownNode) Siblings(Type int) []*MarkdownNode {
	siblings := make([]*MarkdownNode, 0)
	if Type == MdAll {
		for _, node := range md.parent.children {
			if node != md {
				siblings = append(siblings, node)
			}
		}
	} else {
		for _, node := range md.children {
			if node.isFitType(Type) && node != md {
				siblings = append(siblings, node)
			}
		}
	}
	return siblings
}

func (md *MarkdownNode) Next() *MarkdownNode {
	var sibling *MarkdownNode
	findSelf := false
	for _, node := range md.parent.children {
		if findSelf && md.isFitNode(node) {
			sibling = node
			break
		}
		if node == md {
			findSelf = true
		}
	}
	return sibling
}

func (md *MarkdownNode) NextAll() []*MarkdownNode {
	siblings := make([]*MarkdownNode, 0)
	findSelf := false
	for _, node := range md.parent.children {
		if findSelf && md.isFitNode(node) {
			siblings = append(siblings, node)
		}
		if node == md {
			findSelf = true
		}
	}
	return siblings
}

func (md *MarkdownNode) NextUntil(Type int) []*MarkdownNode {
	siblings := make([]*MarkdownNode, 0)
	findSelf := false
	for _, node := range md.parent.children {
		if findSelf && node.isFitType(Type) {
			break
		}
		if findSelf && md.isFitNode(node) {
			siblings = append(siblings, node)
		}
		if node == md {
			findSelf = true
		}
	}
	return siblings
}

func (md *MarkdownNode) Prev() *MarkdownNode {
	var sibling *MarkdownNode
	findSelf := false
	for i := len(md.parent.children) - 1; i >= 0; i-- {
		node := md.parent.children[i]
		if findSelf && md.isFitNode(node) {
			sibling = node
			break
		}
		if node == md {
			findSelf = true
		}
	}
	return sibling
}

func (md *MarkdownNode) PrevAll() []*MarkdownNode {
	siblings := make([]*MarkdownNode, 0)
	findSelf := false
	for i := len(md.parent.children) - 1; i >= 0; i-- {
		node := md.parent.children[i]
		if findSelf && md.isFitNode(node) {
			siblings = append(siblings, node)
		}
		if node == md {
			findSelf = true
		}
	}
	return siblings
}

func (md *MarkdownNode) PrevUntil(Type int) []*MarkdownNode {
	siblings := make([]*MarkdownNode, 0)
	findSelf := false
	for i := len(md.parent.children) - 1; i >= 0; i-- {
		node := md.parent.children[i]
		if findSelf && node.isFitType(Type) {
			break
		}
		if findSelf && md.isFitNode(node) {
			siblings = append(siblings, node)
		}
		if node == md {
			findSelf = true
		}
	}
	return siblings
}

func (md *MarkdownNode) First(Type int) *MarkdownNode {
	var child *MarkdownNode
	if Type == MdAll {
		if len(md.children) > 0 {
			child = md.children[0]
		}
	} else {
		for _, node := range md.children {
			if node.isFitType(Type) {
				child = node
				break
			}
		}
	}
	return child
}

func (md *MarkdownNode) Last(Type int) *MarkdownNode {
	var child *MarkdownNode
	childrenNum := len(md.children)
	if Type == MdAll {
		if childrenNum > 0 {
			child = md.children[childrenNum-1]
		}
	}
	for i := childrenNum - 1; i >= 0; i-- {
		node := md.children[i]
		if node.isFitType(Type) {
			child = node
			break
		}
	}
	return child
}

func (md *MarkdownNode) Eq(idx int) *MarkdownNode {
	var child *MarkdownNode
	ctr := 0
	for _, node := range md.children {
		ctr++
		if ctr >= idx {
			child = node
			break
		}
	}
	return child
}

func (md MarkdownNodes) Filter(Type int) []*MarkdownNode {
	children := make([]*MarkdownNode, 0)
	if Type == MdAll {
		children = append(children, md...)
	} else {
		for _, node := range md {
			if node.isFitType(Type) {
				children = append(children, node)
			}
		}
	}
	return children
}

func (md MarkdownNodes) Not(Type int) []*MarkdownNode {
	children := make([]*MarkdownNode, 0)
	if Type != MdAll {
		for _, node := range md {
			if node._type != Type {
				children = append(children, node)
			}
		}
	}
	return children
}

func (md *MarkdownNode) Text() string {
	return md.text
}

func (md *MarkdownNode) Html() string {
	return md.html
}

func (md *MarkdownNode) Value() string {
	return md.value
}

func (md *MarkdownNode) Append(node *MarkdownNode) {
	if md.children == nil {
		md.children = make([]*MarkdownNode, 0)
	}
	node.parent = md
	md.children = append(md.children, node)
}

func (md *MarkdownNode) Prepend(node *MarkdownNode) {
	if md.children == nil {
		md.children = make([]*MarkdownNode, 0)
	}
	node.parent = md
	children2 := make([]*MarkdownNode, 0)
	children2 = append(children2, node)
	children2 = append(children2, md.children...)
	md.children = children2
}

func (md *MarkdownNode) After(node *MarkdownNode) {
	idx := 0
	for i, node := range md.parent.children {
		if node == md {
			idx = i
			break
		}
	}
	children2 := md.parent.children[0 : idx+1]
	children2 = append(children2, node)
	children2 = append(children2, md.parent.children[idx+1:]...)
	md.parent.children = children2
}

func (md *MarkdownNode) Before(node *MarkdownNode) {
	idx := 0
	for i, node := range md.parent.children {
		if node == md {
			idx = i
			break
		}
	}
	children2 := md.parent.children[0:idx]
	children2 = append(children2, node)
	children2 = append(children2, md.parent.children[idx:]...)
	md.parent.children = children2
}

func (md *MarkdownNode) Remove() {
	idx := 0
	for i, node := range md.parent.children {
		if node == md {
			idx = i
			break
		}
	}
	children2 := md.parent.children[idx+1:]
	md.parent.children = md.parent.children[0:idx]
	md.parent.children = append(md.parent.children, children2...)
}

func (md *MarkdownNode) Empty() {
	md.children = md.children[0:0]
}
