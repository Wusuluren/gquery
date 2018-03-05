package gquery

import (
	"fmt"
	"strings"
)

var _ = fmt.Println

type GqueryHtml struct {
	treeRoot *HtmlNode
}

const (
	cLabel = iota
	cAttrName
	cAttrValue
	cValue
	cOpenTag
	cCloseTag
	cText
	cComment
)

func printNodeTree(node *HtmlNode, tabNum int) {
	tabArry := make([]string, tabNum)
	for i := 0; i < tabNum; i++ {
		tabArry[i] = "\t"
	}
	tabStr := strings.Join(tabArry, "")
	if node.label != "" {
		fmt.Println(tabStr, "label: ", node.label)
	}
	if node.text != "" {
		fmt.Println(tabStr, "text: ", node.text)
	}
	if len(node.attr) > 0 {
		fmt.Println(tabStr, "attr: ", node.attr)
	}
	if len(node.children) > 0 {
		fmt.Println(tabStr, "child: ", len(node.children))
		for _, child := range node.children {
			printNodeTree(child, tabNum+1)
		}
	}
}

func printNodeList(nodeList []*HtmlNode) {
	for _, node := range nodeList {
		printNodeTree(node, 0)
	}
}

func (gq *GqueryHtml) parse(html string) *HtmlNode {
	stats := make([]string, 0, 64)
	inQuote := false
	beginIdx := 0
	endIdx := 0
	except := cOpenTag
	for i := 0; i < len(html); i++ {
		cur := html[i]
		if cur == '"' {
			inQuote = !inQuote
		} else if !inQuote {
			if cur == '<' {
				if html[i+1] == '!' && html[i+2] == '-' && html[i+3] == '-' {
					beginIdx = i
					except = cComment
				}
				if except != cComment {
					if except == cText {
						endIdx = i
						stats = append(stats, html[beginIdx:endIdx])
					}
					beginIdx = i
					except = cCloseTag
				}
			} else if cur == '>' {
				if except == cComment && html[i-1] == '-' && html[i-2] == '-' {
					endIdx = i + 1
					stats = append(stats, html[beginIdx:endIdx])
					beginIdx = i + 1
					except = cOpenTag
				} else if except != cComment {
					endIdx = i + 1
					stats = append(stats, html[beginIdx:endIdx])
					beginIdx = i + 1
					except = cOpenTag
				}
			} else if except != cText {
				if cur != ' ' && cur != '\n' && cur != '\t' && cur != '\r' {
					if except == cOpenTag {
						except = cText
					}
				}
			}
		}
	}
	//for _,stat := range stats {
	//	fmt.Println(stat)
	//}

	nodeList := make([]*HtmlNode, 0)
	for i := 0; i < len(stats); i++ {
		node := stats[i]
		except = cOpenTag
		inQuote = false
		beginIdx = 0
		endIdx = 0
		label := ""
		attr := make(map[string]string)
		text := ""
		attrName := ""
		attrValue := ""
		id := ""
		classes := make([]string, 0)
		for j := 0; j < len(node); j++ {
			cur := node[j]
			if cur == '"' {
				inQuote = !inQuote
			} else if !inQuote {
				if cur == '<' {
					if node[j+1] == '!' {
						beginIdx = j + 4
						except = cComment
					}
					if except == cOpenTag {
						beginIdx = j + 1
						except = cLabel
					}
				} else if cur == '>' {
					if node[j-1] == '/' {
						endIdx = j - 1
					} else {
						endIdx = j
					}
					if except == cAttrName {
						attrName = strings.Trim(node[beginIdx:endIdx], "\t\n\r ")
						attr[attrName] = ""
						beginIdx = endIdx + 1
						except = cOpenTag
					} else if except == cAttrValue {
						attrValue = node[beginIdx:endIdx]
						attrValue = strings.TrimLeft(attrValue, "\"")
						attrValue = strings.TrimRight(attrValue, "\"")
						if attrName == "id" {
							id = attrValue
						} else if attrName == "class" {
							classes = strings.Split(attrValue, " ")
						} else {
							attr[attrName] = attrValue
						}
						beginIdx = endIdx + 1
						except = cOpenTag
					} else if except == cLabel {
						label = node[beginIdx:endIdx]
						beginIdx = endIdx + 1
						except = cOpenTag
					} else if except == cComment {
						if node[j-1] == '-' {
							endIdx = j - 2
							text = node[beginIdx:endIdx]
							beginIdx = j + 1
							except = cOpenTag
							break
						}
					}
				} else if cur == ' ' {
					endIdx = j
					if except == cLabel {
						label = node[beginIdx:endIdx]
						beginIdx = endIdx + 1
						except = cAttrName
					} else if except == cAttrValue {
						attrValue = node[beginIdx:endIdx]
						attrValue = strings.TrimLeft(attrValue, "\"")
						attrValue = strings.TrimRight(attrValue, "\"")
						if attrName == "id" {
							id = attrValue
						} else if attrName == "class" {
							classes = strings.Split(attrValue, " ")
						} else {
							attr[attrName] = attrValue
						}
						beginIdx = endIdx + 1
						except = cAttrName
					}
				} else if cur == '=' {
					endIdx = j
					attrName = strings.Trim(node[beginIdx:endIdx], "\t\n\r ")
					beginIdx = endIdx + 1
					except = cAttrValue
				} else {
					if except == cOpenTag {
						text = node
						break
					}
				}
			}

		}
		nodeList = append(nodeList, &HtmlNode{
			label: label,
			id:    id,
			class: classes,
			text:  text,
			html:  text, //FIXME!
			attr:  attr,
		})
	}
	//fmt.Println(len(nodeList))
	//printNodeList(nodeList)

	nodeTree := make([]*HtmlNode, 0)
	for i := 0; i < len(nodeList); i++ {
		cur := nodeList[i]
		if len(cur.label) > 0 {
			if cur.label[0] == '/' { //reduce
				label2 := cur.label[1:len(cur.label)]
				for j := len(nodeTree) - 1; j >= 0; j-- { // find last pair
					label := nodeTree[j].label
					if len(label) > 0 && !nodeTree[j].isPaired {
						if label == label2 {
							if j == len(nodeTree)-1 {
								nodeTree[j].isPaired = true
								break
							}
							children := make([]*HtmlNode, 0)
							children = append(children, nodeTree[j+1:]...)
							nodeTree[j].children = children
							nodeTree[j].isPaired = true

							nodeTree2 := make([]*HtmlNode, 0)
							nodeTree2 = append(nodeTree2, nodeTree[0:j+1]...)
							nodeTree = nodeTree2
							//printNodeList(nodeTree)
							//fmt.Println("+++++++++++++++++++++++")
							break
						}
					}
				}
			} else {
				nodeTree = append(nodeTree, cur)
			}
		} else {
			nodeTree = append(nodeTree, cur)
		}
	}
	//fmt.Println(len(nodeTree), len(nodeList))
	//for _, node := range nodeTree {
	//	printNodeTree(node, 0)
	//}

	nodeTreeRoot := &HtmlNode{
		children: nodeTree,
	}
	for _, node := range nodeTreeRoot.children {
		node.parent = nodeTreeRoot
	}
	return nodeTreeRoot
}

func (gq *GqueryHtml) Gquery(selector string) HtmlNodes {
	return gq.treeRoot.Gquery(selector)
}

func (gq *GqueryHtml) TreeRoot() *HtmlNode {
	return gq.treeRoot
}

func NewHtmlNode(conf map[string]interface{}) *HtmlNode {
	node := &HtmlNode{}
	for name, value := range conf {
		switch name {
		case "label":
			if label, ok := value.(string); ok {
				node.label = label
			}
		case "id":
			if id, ok := value.(string); ok {
				node.id = id
			}
		case "class":
			if class, ok := value.([]string); ok {
				node.class = class
			}
		case "text":
			if text, ok := value.(string); ok {
				node.text = text
			}
		case "html":
			if html, ok := value.(string); ok {
				node.html = html
			}
		case "value":
			if value, ok := value.(string); ok {
				node.value = value
			}
		case "attr":
			if attr, ok := value.(map[string]string); ok {
				node.attr = attr
			}
		}
	}
	return node
}

func NewHtml(html string) *GqueryHtml {
	gq := &GqueryHtml{}
	gq.treeRoot = gq.parse(html)
	return gq
}
