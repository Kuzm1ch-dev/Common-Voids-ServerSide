package common

import (
	"encoding/json"
	"log"
)

type Package struct {
	Code int32
	Data string
	UUID string
}

func (p Package) Marshal() []byte {
	json_data, err := json.Marshal(p)
	if err != nil {
		log.Println("Error:" + err.Error())
	}
	return json_data
}

func UnMarshal(packageByte []byte) Package {
	var packageStruct Package
	err := json.Unmarshal(packageByte, &packageStruct)
	if err != nil {
		log.Println("Error:" + err.Error())
	}
	return packageStruct
}
