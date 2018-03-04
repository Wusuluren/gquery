package gquery_test

import (
	"github.com/wusuluren/gquery"
	"testing"
)

func printHtmlNodeList(t *testing.T, nodeList []*gquery.HtmlNode) {
	for _, node := range nodeList {
		t.Log(node)
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
    <div id="test" class="test1 test2" attr="attr">
	</div>
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
	t.Log(gq.Gquery("head")[0].First("meta").Attr("charset"))
	t.Log(gq.Gquery("meta").Eq(1).Attr("name"))
	t.Log(gq.Gquery("meta").Eq(1) == gq.Gquery("[name]")[0])

	t.Log("test insert")
	node := gquery.NewHtmlNode(map[string]interface{}{
		"label": "p",
		"value": "test",
		"text":  "test",
		"html":  "test",
	})
	t.Log(node)
	gq.Gquery("body")[0].Append(node)
	t.Log(gq.Gquery("body")[0].Last("p") == node)
}
