package lib

import (
	"bytes"
	"encoding/json"
	"os"
	"strconv"
)

type Output struct{}

type Result struct {
	Id *string
}

func (o Output) OutputFile(tf_resource_name string, result []Result) error {
	f, err := os.OpenFile("output.tf", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
		replaced, err := convertJsonTf([]byte(block))
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

func convertJsonTf(block []byte) (replaced []byte, err error) {
	block = []byte(block)
	replaced = bytes.ReplaceAll(block, []byte("\"id\":"), []byte("  id = "))
	replaced = bytes.ReplaceAll(replaced, []byte("\"to\":"), []byte("  to = "))
	replaced = bytes.ReplaceAll(replaced, []byte(","), []byte("\n"))
	replaced = bytes.ReplaceAll(replaced, []byte("{"), []byte("import {\n"))
	replaced = bytes.ReplaceAll(replaced, []byte("}"), []byte("\n}\n\n"))
	return replaced, nil
}

type ImportBlock struct {
	Id *string `json:"id"`
	To string  `json:"to"`
}
