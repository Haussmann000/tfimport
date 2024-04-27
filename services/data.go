package service

type DataContainer struct {
	DataService DataService
}

type DataService interface {
	OutputTffile(string) (string, error)
	Outputfile(string) (string, error)
}

type Output struct {
	TfOutput TfOutput
}

type TfOutput struct {
	id string
}

// func NewDataContainer() *DataContainer {
// 	return &DataContainer{
// 		DataService: NewOutputFile(),
// 	}
// }

// func NewOutputFile() *Output {
// 	return
// }
