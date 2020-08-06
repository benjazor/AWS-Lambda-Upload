package main

import "encoding/json"

type FunctionList struct {
	Functions []Function
}
type Function struct {
	FunctionName string
}

func NewFunctionList() (FunctionList, error) {
	// Get the list of functions
	data, err := run("aws", "lambda", "list-functions")
	if err != nil {
		return FunctionList{}, err
	}

	// Parse JSON
	var result FunctionList
	err = json.Unmarshal(data, &result)
	if err != nil {
		return FunctionList{}, err
	}

	return result, err
}

func (fl FunctionList) HasFunction(functionName string) bool {
	for _, value := range fl.Functions {
		if value.FunctionName == functionName {
			return true
		}
	}
	return false
}