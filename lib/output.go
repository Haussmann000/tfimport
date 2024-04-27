package lib

import (
	"bytes"
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

var (
	funcMap = map[string]interface{}{
		"parsetags": func(tags []types.Tag) (s string, err error) {
			result, err := parseTag(tags)
			if err != nil {
				return "", err
			}

			return string(result), err
		},
	}
)

func parseTag(tags []types.Tag) (conv []byte, err error) {
	json, err := json.Marshal(tags)
	if err != nil {
		return nil, err
	}
	conv = []byte("[]")
	if len(tags) != 0 {
		conv = convertJsonToTagString([]byte(json))
	}
	return conv, err
}

func convertJsonToTagString(block []byte) (replaced []byte) {
	replaced = bytes.ReplaceAll(block, []byte("\"Key\":"), []byte(""))
	replaced = bytes.ReplaceAll(replaced, []byte(",\"Value\":"), []byte(" = "))
	replaced = bytes.ReplaceAll(replaced, []byte(","), []byte(",\n\t\t"))
	replaced = bytes.ReplaceAll(replaced, []byte("["), []byte("[\n\t\t"))
	replaced = bytes.ReplaceAll(replaced, []byte("]"), []byte("\n\t]"))
	return replaced
}

func parseTemplate[T TfOutput](file io.Writer, output T, tmpStr string) (err error) {
	tmpl, _ := template.New("funcMap").Funcs(funcMap).Parse(tmpStr)
	err = tmpl.Execute(file, output)
	if err != nil {
		return err
	}
	return err
}
