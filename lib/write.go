package lib

import (
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
		template.To = "aws_vpc.my_vpc" + strconv.Itoa(i)
		block, err := json.Marshal(template)
		if err != nil {
			return err
		}
		outputs = append(outputs, block)
		if _, err := f.Write(outputs[i]); err != nil {
			return err
		}
	}
	return nil
}

// func readTemplate() (template []byte, err error) {
// 	template, err = os.ReadFile(TEMPLATEDIR)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return template, nil
// }

// func convertFile(json []byte, vpcId []byte) (replaced []byte, err error) {
// 	bytes.Replace(json, []byte(fmt.Sprintf(`${id}`)), vpcId, -1)
// 	fmt.Printf("%s", replaced)
// 	return replaced, nil
// }
