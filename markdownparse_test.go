package gquery_test

import (
	"github.com/wusuluren/gquery"
	"testing"
)

func printNodeTree(t *testing.T, nodeTree *gquery.MarkdownNode) {
	t.Log(nodeTree)
	for _, node := range nodeTree.Children(gquery.MdNone) {
		printNodeTree(t, node)
	}
}

func TestParseMarkdown(t *testing.T) {
	testData := `
# Title
This is title
- baidu
http://www.baidu.com
	- baidu
	http://www.baidu.com
- google
http://www.google.com
`
	treeRoot := gquery.ParseMarkdown(testData)
	t.Log(len(treeRoot.Children(gquery.MdNone)))
	printNodeTree(t, treeRoot)

	t.Log("search")
	t.Log(treeRoot.Children(gquery.MdTitle)[0])
	t.Log(treeRoot.First(gquery.MdUnorderList))
}
