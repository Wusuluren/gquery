package gquery

import (
	"fmt"
	"strconv"
	"strings"
)

var _ = fmt.Println

type MarkdownNode struct {
	Type      int
	Attribute map[string]string
	Text      string
	RawText   string
	Value     []string
	parent    *MarkdownNode
	children  []*MarkdownNode
	tabNum    int
}

const (
	MdNone = iota
	MdTitle
	MdMajorTitle
	MdSubTitle
	MdParagraph
	MdOrderList
	MdUnorderList
	MdAttributeName
	MdAttributeValue
	MdQuote
	MdCodeblock
	MdInlineCode
	MdHref
	MdImage
	MdStrong
	MdTable
	MdSeparateLine
	MdDeleteLine
)

func isDigital(c byte) bool {
	return c >= '0' && c <= '9'
}
func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}
func isTitle(line string) (bool, int, string) {
	if len(line) >= 2 && line[0:2] == "# " {
		return true, 2, "1"
	} else if len(line) >= 3 && line[0:3] == "## " {
		return true, 3, "2"
	} else if len(line) >= 4 && line[0:4] == "### " {
		return true, 4, "3"
	} else if len(line) >= 5 && line[0:5] == "#### " {
		return true, 5, "4"
	} else if len(line) >= 6 && line[0:6] == "##### " {
		return true, 6, "5"
	} else if len(line) >= 7 && line[0:7] == "###### " {
		return true, 7, "6"
	} else {
		return false, 0, "0"
	}
}
func isParagraph(line string) (bool, int) {
	i := 0
	for _, c := range line {
		if c != '\t' && c != ' ' {
			break
		}
		i++
	}
	if i >= len(line) {
		return false, 0
	}
	c := line[i]
	return isAlpha(c) || isDigital(c), i
}
func isOrderList(line string) (bool, int) {
	if len(line) < 3 {
		return false, 0
	}
	return isDigital(line[0]) && line[1] == '.' && line[2] == ' ', 3
}
func isUnorderList(line string) (bool, int) {
	if len(line) < 2 {
		return false, 0
	}
	isValidStart := line[0] == '-' || line[0] == '+' || line[0] == '*'
	return isValidStart && line[1] == ' ', 2
}
func isQuote(line string) (bool, int, string) {
	num := 0
	for _, c := range line {
		if c != '>' {
			break
		}
		num++
	}
	if num == 0 {
		return false, 0, "0"
	}
	if line[num] == ' ' {
		return true, num + 1, strconv.Itoa(num)
	}
	return false, 0, "0"
}

func ParseMarkdown(markdown string) *MarkdownNode {
	lines := strings.Split(markdown, "\n")
	nodeList := make([]*MarkdownNode, 0)
	curTag := MdNone
	for i := 0; i < len(lines); i++ {
		getValue := func() []string {
			value := make([]string, 0)
			for k := i + 1; k < len(lines); k++ {
				if ok, idx := isParagraph(lines[k]); ok {
					value = append(value, lines[k][idx:])
					i++
				} else {
					break
				}
			}
			return value
		}
		line := lines[i]
		tabNum := 0
		if curTag == MdNone {
			for j := 0; j < len(line); j++ {
				if line[j] == '\t' { //should ignore blank space?
					tabNum++
					continue
				}
				if ok, idx, level := isTitle(line[j:]); ok {
					nodeList = append(nodeList, &MarkdownNode{
						Type: MdTitle,
						Attribute: map[string]string{
							"level": level,
						},
						Text:     line[tabNum+idx:],
						RawText:  line[tabNum:],
						Value:    getValue(),
						tabNum:   tabNum,
						children: make([]*MarkdownNode, 0),
					})
					break
				} else if ok, idx := isOrderList(line[j:]); ok {
					nodeList = append(nodeList, &MarkdownNode{
						Type:     MdOrderList,
						Text:     line[tabNum+idx:],
						RawText:  line[tabNum:],
						Value:    getValue(),
						tabNum:   tabNum,
						children: make([]*MarkdownNode, 0),
					})
					break
				} else if ok, idx := isUnorderList(line[j:]); ok {
					nodeList = append(nodeList, &MarkdownNode{
						Type:     MdUnorderList,
						Text:     line[tabNum+idx:],
						RawText:  line[tabNum:],
						Value:    getValue(),
						tabNum:   tabNum,
						children: make([]*MarkdownNode, 0),
					})
					break
				} else if ok, idx, num := isQuote(line[j:]); ok {
					nodeList = append(nodeList, &MarkdownNode{
						Type: MdUnorderList,
						Attribute: map[string]string{
							"num": num,
						},
						Text:     line[tabNum+idx:],
						RawText:  line[tabNum:],
						tabNum:   tabNum,
						children: make([]*MarkdownNode, 0),
					})
					break
				}
			}
		} else {

		}
	}

	nodeTree := make([]*MarkdownNode, 0)
	lastTabNum := -1
	for i := 0; i < len(nodeList); i++ {
		cur := nodeList[i]
		if cur.tabNum < lastTabNum { //reduce
			for j := len(nodeTree) - 1; j >= 0; j-- {
				if nodeTree[j].tabNum <= cur.tabNum {
					nodeTree[j].children = append(nodeTree[j].children, nodeTree[j+1:]...)
					for _, node := range nodeTree[j+1:] {
						node.parent = nodeTree[j]
					}
					nodeTree = nodeTree[0 : j+1]
					break
				}
			}
			nodeTree = append(nodeTree, cur)
		} else {
			nodeTree = append(nodeTree, cur)
		}
		lastTabNum = cur.tabNum
	}
	nodeTreeRoot := &MarkdownNode{
		children: nodeTree,
	}
	for _, node := range nodeTreeRoot.children {
		node.parent = nodeTreeRoot
	}
	return nodeTreeRoot
}

func (md *MarkdownNode) Parent() *MarkdownNode {
	return md.parent
}

func (md *MarkdownNode) Children(Type int) []*MarkdownNode {
	results := make([]*MarkdownNode, 0)
	if Type == MdNone {
		results = append(results, md.children...)
	} else {
		for _, node := range md.children {
			if node.Type == Type {
				results = append(results, node)
			}
		}
	}
	return results
}

func (md *MarkdownNode) Find(Type int) *MarkdownNode {
	var result *MarkdownNode
	for _, node := range md.children {
		if node.Type == Type {
			result = node
			break
		}
	}
	return result
}

func (md *MarkdownNode) Next() *MarkdownNode {
	var result *MarkdownNode
	for i, node := range md.parent.children {
		if node == md {
			if i+1 < len(md.parent.children) {
				result = md.parent.children[i+1]
			}
			break
		}
	}
	return result
}

func (md *MarkdownNode) First(Type int) *MarkdownNode {
	var result *MarkdownNode
	for _, node := range md.children {
		if node.Type == Type {
			result = node
			break
		}
	}
	return result
}

func (md *MarkdownNode) Last(Type int) *MarkdownNode {
	var result *MarkdownNode
	for i := len(md.children) - 1; i >= 0; i-- {
		node := md.children[i]
		if node.Type == Type {
			result = node
			break
		}
	}
	return result
}

func (md *MarkdownNode) Eq(Type, idx int) *MarkdownNode {
	var result *MarkdownNode
	ctr := 0
	for _, node := range md.children {
		if node.Type == Type {
			ctr++
			if ctr >= idx {
				result = node
				break
			}
		}
	}
	return result
}

func (md *MarkdownNode) Append(node *MarkdownNode) {
	if md.children == nil {
		md.children = make([]*MarkdownNode, 0)
	}
	node.parent = md
	md.children = append(md.children, node)
}

func (md *MarkdownNode) Remove() {
	idx := 0
	if md.parent == nil {
		fmt.Println(md)
	}
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
