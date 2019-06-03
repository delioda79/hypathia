// Code generated by hero.
// source: /Users/angelosvalsamis/go/src/hypatia/template/api_list.html
// DO NOT EDIT!
package template

import (
	"bytes"
	"hypatia/scrape"

	"github.com/shiyanhui/hero"
)

func ApiList(docs []scrape.DocDef, buffer *bytes.Buffer) {
	buffer.WriteString(`<!DOCTYPE html>
<!doctype html> <!-- Important: must specify -->
<html>
<head>
    <meta charset="utf-8"> <!-- Important: rapi-doc uses utf8 charecters -->
    <script src="https://unpkg.com/rapidoc/dist/rapidoc-min.js"></script>
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
    <script src="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js" integrity="sha384-JjSmVgyd0p3pXB1rRibZUAYoIIy6OrQ6VrjIEaFf/nJGzIxFDsf4x0xIM+B07jRM" crossorigin="anonymous"></script>
    <script src="https://code.jquery.com/jquery-3.3.1.slim.min.js" integrity="sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo" crossorigin="anonymous"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js" integrity="sha384-UO2eT0CpHqdSJQ6hJty5KVphtPhzWj9WO1clHTMGa3JDZwrnQq4sF86dIHNDz0W1" crossorigin="anonymous"></script>
</head>

<body>
<header style="background-color: #23D2AA"><img src="https://cdn.onelogin.com/images/brands/logos/login/af659f099a3585125f9cb0ec562868f7e70b1406.png?1532445519"/></header>
`)
	buffer.WriteString(`
<ul class="list-group">
`)
	for _, doc := range docs {
		buffer.WriteString(`
    <li class="list-group-item"><a  href='/doc/`)
		hero.EscapeHTML(doc.RepoName, buffer)
		buffer.WriteString(`/`)
		hero.FormatInt(int64(doc.Type), buffer)
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
