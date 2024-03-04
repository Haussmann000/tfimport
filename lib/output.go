package lib

import (
	"bytes"
	"embed"
	"encoding/json"
	"io"
	"os"
	"regexp"
	"strconv"
	"text/template"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type Result struct {
	Id *string
}

func OutputFile(tf_resource_name string, result []Result) error {
	f, err := os.OpenFile("import.tf", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	for i, resource := range result {
		template := ImportBlock{}
		template.Id = resource.Id
		template.To = tf_resource_name + ".my_resource" + strconv.Itoa(i)
		block, err := json.Marshal(template)
		if err != nil {
			return err
		}
		replaced, err := convertJsonToImportBlock([]byte(block))
		if err != nil {
			return err
		}
		if _, err := f.Write(replaced); err != nil {
			return err
		}
	}
	return nil
}

func convertJsonToImportBlock(block []byte) (replaced []byte, err error) {
	replaced = bytes.ReplaceAll(block, []byte("\"id\":"), []byte("  id = "))
	replaced = bytes.ReplaceAll(replaced, []byte("\"to\":"), []byte("  to = "))
	replaced = bytes.ReplaceAll(replaced, []byte(","), []byte("\n"))
	replaced = bytes.ReplaceAll(replaced, []byte("{"), []byte("import {\n"))
	replaced = bytes.ReplaceAll(replaced, []byte("}"), []byte("\n}\n\n"))
	r := regexp.MustCompile(`(to = )(")(.*)(")`)
	replaced = r.ReplaceAll(replaced, []byte("$1$3"))
	return replaced, nil
}

func OutputTfFile[T TfOutput](x T, resource_name string) error {
	f, err := os.OpenFile("output.tf", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	parseTemplate(f, x, resource_name)
	defer f.Close()
	return nil
}

//go:embed template/*
var f embed.FS

var (
	funcMap = map[string]interface{}{
		"parsetags": func(tags []types.Tag) []types.Tag {
			// result, err := json.Marshal(tags)
			// if err != nil {
			// 	panic(err)
			// }
			return tags
		},
	}
)

func parseTemplate[T TfOutput](file io.Writer, output T, temp_path string) (err error) {
	tmpl, _ := template.ParseFS(f, "template/"+temp_path+".tmpl")
	err = tmpl.Funcs(funcMap).Execute(file, output)
	if err != nil {
		return err
	}
	return err
}
