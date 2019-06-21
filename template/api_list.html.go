// Code generated by hero.
// source: /Users/deliodanna/goprojects/hypatia/template/api_list.html
// DO NOT EDIT!
package template

import (
	"bytes"

	"github.com/taxibeat/hypatia/scrape"

	"github.com/shiyanhui/hero"
)

func ApiList(docs []scrape.DocDef, buffer *bytes.Buffer) {
	buffer.WriteString(`<!DOCTYPE html>
<!doctype html> <!-- Important: must specify -->
<html>
<head>
    <meta charset="utf-8"> <!-- Important: rapi-doc uses utf8 charecters -->
    <script src="/static/js/rapidoc-min.js"></script>
    <link href="/static/css/font1.css" rel="stylesheet">
    <link rel="stylesheet" href="/static/css/bootstrap.min.css" crossorigin="anonymous">
    <script src="/static/js/bootstrap.min.js" crossorigin="anonymous"></script>
    <script src="/static/js/jquery-3.3.1.slim.min.js" crossorigin="anonymous"></script>
    <script src="/static/js/popper.min.js" crossorigin="anonymous"></script>
</head>

<body>
<header style="background-color: #23D2AA" class="row">
    <a href="/" class="col-1"><img src="/static/img/beat-sm.png"/></a>
    <span class="col-10 text-center" style="font-size: xx-large">HYPATIA</span>
    <img src="/static/img/hypatia.png" height="100" class="col-1 float-right">
</header>
`)
	buffer.WriteString(`
<div>
    <form action="/" method="POST">
        <input type="text" name="query"/>
        <input type="submit"/>
    </form>
</div>
<ul class="list-group">
`)
	for _, doc := range docs {
		buffer.WriteString(`
    <li class="list-group-item"><a  href='/doc/`)
		hero.EscapeHTML(doc.ID, buffer)
		buffer.WriteString(`'>Api Documentation for `)
		hero.EscapeHTML(doc.RepoName, buffer)
		buffer.WriteString(` ( `)
		hero.EscapeHTML(doc.Type.String(), buffer)
		buffer.WriteString(` )</a></li>
`)
	}
	buffer.WriteString(`
</ul>
`)

	buffer.WriteString(`
</body>
</html>`)

}
