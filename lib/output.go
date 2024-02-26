package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
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
	outputs := [][]byte{}
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
		outputs = append(outputs, replaced)
		if _, err := f.Write(outputs[i]); err != nil {
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

func convertJsonToTfFile(block []byte) (replaced []byte, err error) {
	replaced = bytes.ReplaceAll(block, []byte("CidrBlock"), []byte("cidr_block"))
	return replaced, nil
}

func OutputTfFile[T any](x T) error {
	json, err := json.Marshal(x)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", json)
	f, err := os.OpenFile("output.tf", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	replaced, err := convertJsonToTfFile([]byte(json))
	if err != nil {
		return err
	}
	if _, err := f.Write(replaced); err != nil {
		return err
	}
	return nil
}

type ImportBlock struct {
	Id *string `json:"id"`
	To string  `json:"to"`
}
