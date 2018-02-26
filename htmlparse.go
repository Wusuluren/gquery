package gquery

import (
	"fmt"
	"strings"
)

type HtmlNode struct {
	Label     string
	Attribute map[string]string
	Text      string
	parent    *HtmlNode
	child     []*HtmlNode
	isPaired  bool
}

const (
	Label = iota
	AttributeName
	AttributeValue
	Value
	OpenTag
	CloseTag
	Text
	Comment
)

func (hn *HtmlNode) Parent() *HtmlNode {
	return hn.parent
}

func (hn *HtmlNode) Children(str string) []*HtmlNode {
	result := make([]*HtmlNode, 0)
	if str != "*" && str != "" {
		tmps := strings.Split(str, ".")
		labelName := tmps[0]
		className := ""
		idName := ""
		if len(tmps) > 1 {
			if tmps[1][0] == '#' {
				idName = tmps[1][1:]
			} else {
				className = tmps[1]
			}
		}
		for _, node := range hn.child {
			if node.Label == labelName {
				if className != "" {
					if value, ok := node.Attribute["class"]; ok {
						//fmt.Println(className, value)
						//if className == value {
						if ReStrCmp(value, className) {
							//fmt.Println(node)
							result = append(result, node)
						}
					}
				} else if idName != "" {
					if value, ok := node.Attribute["id"]; ok {
						if idName == value {
							result = append(result, node)
						}
					}
				} else {
					result = append(result, node)
				}
			}
		}
	} else {
		result = append(result, hn.child...)
	}
	return result
}

func (hn *HtmlNode) Find(str string) *HtmlNode {
	result := &HtmlNode{}
	if str != "*" && str != "" {
		tmps := strings.Split(str, ".")
		labelName := tmps[0]
		className := ""
		idName := ""
		if len(tmps) > 1 {
			if tmps[1][0] == '#' {
				idName = tmps[1][1:]
			} else {
				className = tmps[1]
			}
		}
		for _, node := range hn.child {
			if node.Label == labelName {
				if className != "" {
					if value, ok := node.Attribute["class"]; ok {
						//fmt.Println(value, className)
						//if className == value {
						if ReStrCmp(value, className) {
							result = node
							break
						}
					}
				} else if idName != "" {
					if value, ok := node.Attribute["id"]; ok {
						if idName == value {
							result = node
							break
						}
					}
				} else {
					result = node
					break
				}
			}
		}
	} else {
		if len(hn.child) > 0 {
			result = hn.child[0]
		}
	}
	//fmt.Println(str, result)
	return result
}

func (hn *HtmlNode) Next() *HtmlNode {
	result := &HtmlNode{}
	return result
}

func (hn *HtmlNode) First(str string) *HtmlNode {
	result := &HtmlNode{}
	if str != "" {
		tmps := strings.Split(str, ".")
		idName := tmps[0]
		className := ""
		if len(tmps) > 1 {
			className = tmps[1]
		}
		for _, node := range hn.child {
			if node.Label == idName {
				if className != "" {
					if value, ok := node.Attribute["class"]; ok {
						if ReStrCmp(value, className) {
							result = node
							break
						}
					}
				} else {
					result = node
					break
				}
			}
		}
	} else {
		if len(hn.child) > 0 {
			result = hn.child[0]
		}
	}
	return result
}

func (hn *HtmlNode) Last(str string) *HtmlNode {
	result := &HtmlNode{}
	childNum := len(hn.child)
	if str != "" && childNum > 0 {
		tmps := strings.Split(str, ".")
		idName := tmps[0]
		className := ""
		if len(tmps) > 1 {
			className = tmps[1]
		}
		for i := childNum - 1; i >= 0; i-- {
			node := hn.child[i]
			if node.Label == idName {
				if className != "" {
					if value, ok := node.Attribute["class"]; ok {
						if ReStrCmp(value, className) {
							result = node
							break
						}
					}
				} else {
					result = node
					break
				}
			}
		}
	} else {
		if childNum > 0 {
			result = hn.child[childNum-1]
		}
	}
	return result
}

func (hn *HtmlNode) Eq(str string, idx int) *HtmlNode {
	result := &HtmlNode{}
	ctr := 0
	if str != "*" && str != "" {
		tmps := strings.Split(str, ".")
		idName := tmps[0]
		className := ""
		if len(tmps) > 1 {
			className = tmps[1]
		}
		for _, node := range hn.child {
			if node.Label == idName {
				if className != "" {
					if value, ok := node.Attribute["class"]; ok {
						if ReStrCmp(value, className) {
							if ctr == idx {
								result = node
								break
							} else {
								ctr++
							}
						}
					}
				} else {
					if ctr == idx {
						result = node
						break
					} else {
						ctr++
					}
				}
			}
		}
	} else {
		for i, node := range hn.child {
			if i == idx {
				result = node
				break
			}
		}
	}
	return result
}

func printNodeTree(node *HtmlNode, tabNum int) {
	tabArry := make([]string, tabNum)
	for i := 0; i < tabNum; i++ {
		tabArry[i] = "\t"
	}
	tabStr := strings.Join(tabArry, "")
	if node.Label != "" {
		fmt.Println(tabStr, "label: ", node.Label)
	}
	if node.Text != "" {
		fmt.Println(tabStr, "text: ", node.Text)
	}
	if len(node.Attribute) > 0 {
		fmt.Println(tabStr, "attr: ", node.Attribute)
	}
	if len(node.child) > 0 {
		fmt.Println(tabStr, "child: ", len(node.child))
		for _, child := range node.child {
			printNodeTree(child, tabNum+1)
		}
	}
}

func printNodeList(nodeList []*HtmlNode) {
	for _, node := range nodeList {
		printNodeTree(node, 0)
	}
}

func ParseHtml(html string) []*HtmlNode {
	nodes := make([]string, 0, 64)
	inQuote := false
	beginIdx := 0
	endIdx := 0
	except := OpenTag
	for i := 0; i < len(html); i++ {
		cur := html[i]
		if cur == '"' {
			inQuote = !inQuote
		} else if !inQuote {
			if cur == '<' {
				if html[i+1] == '!' && html[i+2] == '-' && html[i+3] == '-' {
					beginIdx = i
					except = Comment
				}
				if except != Comment {
					if except == Text {
						endIdx = i
						nodes = append(nodes, html[beginIdx:endIdx])
					}
					beginIdx = i
					except = CloseTag
				}
			} else if cur == '>' {
				if except == Comment && html[i-1] == '-' && html[i-2] == '-' {
					endIdx = i + 1
					nodes = append(nodes, html[beginIdx:endIdx])
					beginIdx = i + 1
					except = OpenTag
				} else if except != Comment {
					endIdx = i + 1
					nodes = append(nodes, html[beginIdx:endIdx])
					beginIdx = i + 1
					except = OpenTag
				}
			} else if except != Text {
				if cur != ' ' && cur != '\n' && cur != '\t' && cur != '\r' {
					if except == OpenTag {
						except = Text
					}
				}
			}
		}
	}
	//for _,node := range nodes {
	//	fmt.Println(node)
	//}

	HtmlNodeList := make([]*HtmlNode, 0)
	for i := 0; i < len(nodes); i++ {
		node := nodes[i]
		except = OpenTag
		inQuote = false
		beginIdx = 0
		endIdx = 0
		label := ""
		attr := make(map[string]string)
		text := ""
		attrName := ""
		attrValue := ""
		for j := 0; j < len(node); j++ {
			cur := node[j]
			if cur == '"' {
				inQuote = !inQuote
			} else if !inQuote {
				if cur == '<' {
					if node[j+1] == '!' {
						beginIdx = j + 4
						except = Comment
					}
					if except == OpenTag {
						beginIdx = j + 1
						except = Label
					}
				} else if cur == '>' {
					if node[j-1] == '/' {
						endIdx = j - 1
					} else {
						endIdx = j
					}
					if except == AttributeName {
						attrName = strings.Trim(node[beginIdx:endIdx], "\t\n\r ")
						attr[attrName] = ""
						beginIdx = endIdx + 1
						except = OpenTag
					} else if except == AttributeValue {
						attrValue = node[beginIdx:endIdx]
						attr[attrName] = attrValue
						beginIdx = endIdx + 1
						except = OpenTag
					} else if except == Label {
						label = node[beginIdx:endIdx]
						beginIdx = endIdx + 1
						except = OpenTag
					} else if except == Comment {
						if node[j-1] == '-' {
							endIdx = j - 2
							text = node[beginIdx:endIdx]
							beginIdx = j + 1
							//fmt.Println(text)
							except = OpenTag
							break
						}
					}
				} else if cur == ' ' {
					endIdx = j
					if except == Label {
						label = node[beginIdx:endIdx]
						beginIdx = endIdx + 1
						except = AttributeName
					} else if except == AttributeValue {
						attrValue = node[beginIdx:endIdx]
						attr[attrName] = attrValue
						beginIdx = endIdx + 1
						except = AttributeName
					}
				} else if cur == '=' {
					endIdx = j
					attrName = strings.Trim(node[beginIdx:endIdx], "\t\n\r ")
					beginIdx = endIdx + 1
					except = AttributeValue
				} else {
					if except == OpenTag {
						text = node
						break
					}
				}
			}

		}
		HtmlNodeList = append(HtmlNodeList, &HtmlNode{
			Label:     label,
			Text:      text,
			Attribute: attr,
		})
	}
	//fmt.Println(len(HtmlNodeList))
	//printNodeList(HtmlNodeList)

	htmlNodeTree := make([]*HtmlNode, 0)
	for i := 0; i < len(HtmlNodeList); i++ {
		cur := HtmlNodeList[i]
		if len(cur.Label) > 0 {
			if cur.Label[0] == '/' { //reduce
				label2 := cur.Label[1:len(cur.Label)]
				for j := len(htmlNodeTree) - 1; j >= 0; j-- { // find last pair
					label := htmlNodeTree[j].Label
					if len(label) > 0 && !htmlNodeTree[j].isPaired {
						if label == label2 {
							if j == len(htmlNodeTree)-1 {
								htmlNodeTree[j].isPaired = true
								break
							}
							childList := make([]*HtmlNode, 0)
							for _, node := range htmlNodeTree[j+1:] {
								childList = append(childList, node)
							}
							htmlNodeTree[j].child = childList
							htmlNodeTree[j].isPaired = true

							htmlNodeTree2 := make([]*HtmlNode, 0)
							for _, node := range htmlNodeTree[0 : j+1] {
								htmlNodeTree2 = append(htmlNodeTree2, node)
							}
							htmlNodeTree = htmlNodeTree2
							//printNodeList(htmlNodeTree)
							//fmt.Println("+++++++++++++++++++++++")
							break
						}
					}
				}
			} else {
				htmlNodeTree = append(htmlNodeTree, cur)
			}
		} else {
			htmlNodeTree = append(htmlNodeTree, cur)
		}
	}
	//fmt.Println(len(htmlNodeTree), len(HtmlNodeList))
	//for _, node := range htmlNodeTree {
	//	printNodeTree(node, 0)
	//}

	return htmlNodeTree
}
