package lib

import (
	"bytes"
	"encoding/json"
	"os"
	"strconv"

	service "github.com/Haussmann000/tfimport/services"
)

func OutputFile(result service.VpcResults) error {
	f, err := os.OpenFile("output.tf", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	outputs := [][]byte{}
	for i, vpc := range result {
		template := service.ImportBlock{}
		template.Id = *vpc.VpcId
		template.To = VPC_RESOUCE + ".my_resource" + strconv.Itoa(i)
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
