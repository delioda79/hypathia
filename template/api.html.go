// Code generated by hero.
// source: /Users/angelosvalsamis/newgo/hypatia/template/api.html
// DO NOT EDIT!
package template

import (
	"bytes"

	"github.com/shiyanhui/hero"
	"github.com/taxibeat/hypatia/scrape"
)

func ApiRender(doc scrape.DocDef, buffer *bytes.Buffer) {
	buffer.WriteString(`<!DOCTYPE html>
<!doctype html> <!-- Important: must specify -->
<html>
<head>
    <meta charset="utf-8"> <!-- Important: rapi-doc uses utf8 charecters -->
    <script src="https://unpkg.com/rapidoc/dist/rapidoc-min.js"></script>
    <link href="https://fonts.googleapis.com/css?family=Varela+Round" rel="stylesheet">
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
    <script src="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js" integrity="sha384-JjSmVgyd0p3pXB1rRibZUAYoIIy6OrQ6VrjIEaFf/nJGzIxFDsf4x0xIM+B07jRM" crossorigin="anonymous"></script>
    <script src="https://code.jquery.com/jquery-3.3.1.slim.min.js" integrity="sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo" crossorigin="anonymous"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js" integrity="sha384-UO2eT0CpHqdSJQ6hJty5KVphtPhzWj9WO1clHTMGa3JDZwrnQq4sF86dIHNDz0W1" crossorigin="anonymous"></script>
</head>

<body>
<header style="background-color: #23D2AA"><a href="/"><img src="https://cdn.onelogin.com/images/brands/logos/login/af659f099a3585125f9cb0ec562868f7e70b1406.png?1532445519"/></a></header>
`)
	buffer.WriteString(`
        <rapi-doc spec-url="/spec/`)
	hero.EscapeHTML(doc.RepoName, buffer)
	buffer.WriteString(`/`)
	hero.FormatInt(int64(doc.Type), buffer)
	buffer.WriteString(`" allow-spec-url-load="false"
                  allow-spec-file-load="false" header-color="#FFFFFF" allow-search="true"
                  regular-font="'Varela Round', 'Arial Rounded MT Bold', 'Helvetica Rounded' ">
        </rapi-doc>
    `)

	buffer.WriteString(`
</body>
</html>`)

}
