package myfile

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//WriteJSON write input data as json format with filename
func WriteJSON(inputData interface{}, filename string) error {
	data, err := json.MarshalIndent(inputData, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, data, 0666)
	if err != nil {
		return err
	}
	return nil
}

//ReadJSON read json from filename
func ReadJSON(v interface{}, filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func ReadAll(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

func WriteFile(inputData []byte, filename string) error {
	err := ioutil.WriteFile(filename, inputData, 0666)
	if err != nil {
		return err
	}
	return nil
}
