package utils

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
)

// terraform provider constants
const (
	// utility constants

	LoggerID = "terraform-provider-vra7"
)

// UnmarshalJSON  decodes json
func UnmarshalJSON(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}

// MarshalToJSON the object to JSON and convert to *bytes.Buffer
func MarshalToJSON(v interface{}) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	err := json.NewEncoder(buffer).Encode(v)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

// ConvertInterfaceToString cpnverts interface to string
func ConvertInterfaceToString(interfaceData interface{}) string {
	var stringData string
	if reflect.ValueOf(interfaceData).Kind() == reflect.Float64 {
		stringData =
			strconv.FormatFloat(interfaceData.(float64), 'f', 0, 64)
	} else if reflect.ValueOf(interfaceData).Kind() == reflect.Float32 {
		stringData =
			strconv.FormatFloat(interfaceData.(float64), 'f', 0, 32)
	} else if reflect.ValueOf(interfaceData).Kind() == reflect.Int {
		stringData = strconv.Itoa(interfaceData.(int))
	} else if reflect.ValueOf(interfaceData).Kind() == reflect.String {
		stringData = interfaceData.(string)
	} else if reflect.ValueOf(interfaceData).Kind() == reflect.Bool {
		stringData = strconv.FormatBool(interfaceData.(bool))
	}
	return stringData
}

// ReplaceValueInRequestTemplate replaces the value for a given key in a catalog
// request template.
func ReplaceValueInRequestTemplate(templateInterface map[string]interface{}, field string, value interface{}) (map[string]interface{}, bool) {
	var replaced bool
	//Iterate over the map to get field provided as an argument
	for key := range templateInterface {
		//If value type is map then set recursive call which will fiend field in one level down of map interface
		if reflect.ValueOf(templateInterface[key]).Kind() == reflect.Map {
			template, _ := templateInterface[key].(map[string]interface{})
			templateInterface[key], replaced = ReplaceValueInRequestTemplate(template, field, value)
			if replaced == true {
				return templateInterface, true
			}
		} else if key == field {
			//If value type is not map then compare field name with provided field name
			//If both matches then update field value with provided value
			templateInterface[key] = value
			return templateInterface, true
		}
	}
	//Return updated map interface type
	return templateInterface, replaced
}

// AddValueToRequestTemplate modeled after replaceValueInRequestTemplate
// for values being added to template vs updating existing ones
func AddValueToRequestTemplate(templateInterface map[string]interface{}, field string, value interface{}) map[string]interface{} {
	//simplest case is adding a simple value. Leaving as a func in case there's a need to do more complicated additions later
	//	templateInterface[data]
	for k, v := range templateInterface {
		if reflect.ValueOf(v).Kind() == reflect.Map && k == "data" {
			template, _ := v.(map[string]interface{})
			v = AddValueToRequestTemplate(template, field, value)
		} else { //if i == "data" {
			templateInterface[field] = value
		}
	}
	//Return updated map interface type
	return templateInterface
}

// ResourceMapper returns the mapping of resource attributes from ResourceView APIs
// to Catalog Item Request Template APIs
func ResourceMapper() map[string]string {
	m := make(map[string]string)
	m["MachineName"] = "name"
	m["MachineDescription"] = "description"
	m["MachineMemory"] = "memory"
	m["MachineStorage"] = "storage"
	m["MachineCPU"] = "cpu"
	m["MachineStatus"] = "status"
	m["MachineType"] = "type"
	return m
}
