package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
)

func run(prog string, args  ...string) ([]byte, error) {
	cmd := exec.Command(prog, args...)
	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return []byte{}, err
	}

	errPipe, err := cmd.StderrPipe()
	if err != nil {
		return []byte{}, err
	}

	err = cmd.Start()
	if err != nil {
		return []byte{}, err
	}
	var outBuf bytes.Buffer
	var errBuf bytes.Buffer

	ch := make (chan bool)
	go func() {
		io.Copy(&outBuf, outPipe)
		ch <- true
	}()
	io.Copy(&errBuf, errPipe)

	_ = <- ch

	err = cmd.Wait()
	if err != nil {
		return outBuf.Bytes(), err
	}

	if len(errBuf.Bytes()) != 0 {
		return outBuf.Bytes(), fmt.Errorf("%s", errBuf.Bytes())
	}

	return outBuf.Bytes(), nil
}

type RoleList struct {
	Roles []Role
}

type Role struct {
	RoleName string
	Arn		 string
}

func RoleMap()(map[string]string, error) {
	res := make(map[string]string)
	data, err := run("aws", "iam", "list-roles")
	if err != nil {
		return res, err
	}

	var roles RoleList
	err = json.Unmarshal(data, &roles)
	if err != nil {
		return res, err
	}

	for _, value := range roles.Roles {
		res[value.RoleName] = value.Arn
	}
	return res, nil
}