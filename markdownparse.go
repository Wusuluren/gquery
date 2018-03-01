package gquery

import (
	"fmt"
	"strconv"
	"strings"
)

type GqueryMarkdown struct {
	treeRoot *MarkdownNode
}

var _ = fmt.Println

const (
	MdNone = iota
	MdAll
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

func (gq *GqueryMarkdown) parse(markdown string) *MarkdownNode {
	lines := strings.Split(markdown, "\n")
	nodeList := make([]*MarkdownNode, 0)
	curTag := MdNone
	for i := 0; i < len(lines); i++ {
		getValue := func() string {
			value := make([]string, 0)
			for k := i + 1; k < len(lines); k++ {
				if ok, idx := isParagraph(lines[k]); ok {
					value = append(value, lines[k][idx:])
					i++
				} else {
					break
				}
			}
			return strings.Join(value, "\n")
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
						_type: MdTitle,
						attr: map[string]string{
							"level": level,
						},
						text:     line[tabNum+idx:],
						html:     line[tabNum:],
						value:    getValue(),
						tabNum:   tabNum,
						children: make([]*MarkdownNode, 0),
					})
					break
				} else if ok, idx := isOrderList(line[j:]); ok {
					nodeList = append(nodeList, &MarkdownNode{
						_type:    MdOrderList,
						text:     line[tabNum+idx:],
						html:     line[tabNum:],
						value:    getValue(),
						tabNum:   tabNum,
						children: make([]*MarkdownNode, 0),
					})
					break
				} else if ok, idx := isUnorderList(line[j:]); ok {
					nodeList = append(nodeList, &MarkdownNode{
						_type:    MdUnorderList,
						text:     line[tabNum+idx:],
						html:     line[tabNum:],
						value:    getValue(),
						tabNum:   tabNum,
						children: make([]*MarkdownNode, 0),
					})
					break
				} else if ok, idx, num := isQuote(line[j:]); ok {
					nodeList = append(nodeList, &MarkdownNode{
						_type: MdUnorderList,
						attr: map[string]string{
							"num": num,
						},
						text:     line[tabNum+idx:],
						html:     line[tabNum:],
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

func (gq *GqueryMarkdown) Gquery(Type int) []*MarkdownNode {
	return gq.treeRoot.Gquery(Type)
}

func NewMarkdown(str string) *GqueryMarkdown {
	md := &GqueryMarkdown{}
	md.treeRoot = md.parse(str)
	return md
}
