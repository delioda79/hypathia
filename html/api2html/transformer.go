package api2html

import (
	"github.com/taxibeat/hypatia/html"
	"io/ioutil"
	"os"
	"os/exec"
)

type Transformer struct {
	html.Transformer
}

type ApiDef struct {
	fileName string
	fileAPI  string
}

func NewApiDef(fileName string, fileAPI string) ApiDef {
	return ApiDef{
		fileName: fileName,
		fileAPI:  fileAPI,
	}
}

func (at *Transformer) Apply(asyncDefFiles []ApiDef) map[string][]byte {
	asyncDs := make(map[string][]byte)
	for _, asyncFile := range asyncDefFiles {
		rawFile := constructHTML(asyncFile)
		asyncDs[asyncFile.fileName] = rawFile
	}
	return asyncDs
}

const (
	htmlFilePref = "html_"
	apiFilePref  = "api_"
)

func constructHTML(asyncDefFile ApiDef) []byte {
	err := ioutil.WriteFile(htmlFilePref+asyncDefFile.fileName, []byte(asyncDefFile.fileAPI), 0644)
	if err != nil {
		println(err.Error())
		return []byte{}
	}
	defer os.Remove(htmlFilePref + asyncDefFile.fileName)

	app := "api2html"

	arg0 := "-o"
	arg1 := apiFilePref + asyncDefFile.fileName
	arg2 := "-c"
	arg3 := "./static/beat.png"
	arg4 := "-t"
	arg5 := "Atelier Cave Light"
	arg6 := "-u"
	arg7 := "/"
	arg8 := htmlFilePref + asyncDefFile.fileName

	cmd := exec.Command(app, arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8)
	cmd.Stdout = nil
	err = cmd.Run()
	if err != nil {
		println(err.Error())
		return []byte{}
	}

	defer os.Remove(apiFilePref + asyncDefFile.fileName)

	htmlRaw, err := ioutil.ReadFile(apiFilePref + asyncDefFile.fileName)
	if err != nil {
		println(err.Error())
		return []byte{}
	}
	return htmlRaw
}
