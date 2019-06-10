package api2html

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewApiDef(t *testing.T) {
	apiDef := NewApiDef("filename", "fileAPI")
	assert.Equal(t, "filename", apiDef.fileName)
	assert.Equal(t, "fileAPI", apiDef.fileAPI)
}

func TestTransformer_ApplySuccess(t *testing.T) {
	apiDefs := []ApiDef{
		{
			fileAPI: `{
			  "openapi": "3.0.0",
			  "info": {
				"version": "1.0.0",
				"title": "Swagger Petstore"
			  },
			  "servers": [
				{
				  "url": "http://petstore.swagger.io/api"
				}
			  ],
			  "paths": {
				"/pets": {
				  "delete": {
					"operationId": "deletePet",
					"responses": {
					  "204": {
						"description": "pet deleted"
					  }
					}
				  }
				}
			  }
			}
			`,
			fileName: "test",
		},
	}

	actual := (&Transformer{}).Apply(apiDefs)

	assert.NotNil(t, actual[apiDefs[0].fileName])

	//check if files are deleted successfully
	if _, err := os.Stat(apiFilePref + apiDefs[0].fileName); !os.IsNotExist(err) {
		t.Errorf("Error file = %s is not removed.", apiFilePref+apiDefs[0].fileName)
	}
	if _, err := os.Stat(htmlFilePref + apiDefs[0].fileName); !os.IsNotExist(err) {
		t.Errorf("Error file = %s is not removed.", htmlFilePref+apiDefs[0].fileName)
	}
}

func TestTransformer_ApplyInvalidAPIAndFail(t *testing.T) {
	apiDefs := []ApiDef{
		{
			fileAPI:  `{}`,
			fileName: "test",
		},
	}

	actual := (&Transformer{}).Apply(apiDefs)

	assert.Equal(t, []byte{}, actual[apiDefs[0].fileName])

	//check if files are deleted successfully
	if _, err := os.Stat(apiFilePref + apiDefs[0].fileName); !os.IsNotExist(err) {
		t.Errorf("Error file = %s is not removed.", apiFilePref+apiDefs[0].fileName)
	}
	if _, err := os.Stat(htmlFilePref + apiDefs[0].fileName); !os.IsNotExist(err) {
		t.Errorf("Error file = %s is not removed.", htmlFilePref+apiDefs[0].fileName)
	}
}
