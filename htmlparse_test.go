package gquery_test

import (
	"github.com/wusuluren/gquery"
	"testing"
)

func printHtmlNodeList(t *testing.T, nodeList []*gquery.HtmlNode) {
	for _, node := range nodeList {
		t.Log(node)
		//printNodeList(t, node.Children(gquery.MdAll))
	}
}

func TestParseHtml(t *testing.T) {
	testData := `
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Document</title>
</head>
<body>
    
</body>
</html>
`
	gq := gquery.NewHtml(testData)
	children := gq.Gquery("*")
	t.Log(len(children))
	printHtmlNodeList(t, children)

	t.Log("test search")
	t.Log(gq.Gquery("html")[0])
	for _, metaNode := range gq.Gquery("html")[0].Gquery("head") {
		t.Log(metaNode)
	}
	t.Log(gq.Gquery("title")[0].Text() == gq.Gquery("title")[0].Html())
}
