package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"io/ioutil"
)

type LProject struct {
	Name 	string
	Bucket 	string
	Role 	string
	path 	string
}

// Setup new project
func NewLProject(configLocation string) (LProject, error) {
	// Read data from config
	data, err := ioutil.ReadFile(configLocation)
	if err != nil {
		return LProject{}, err
	}
	
	var result LProject
	// Parse config
	err = json.Unmarshal(data, &result)
	if err != nil {
		return LProject{}, err
	}

	result.path = path.Dir(configLocation)

	// Check if role is valid
	if strings.HasPrefix(result.Role,"arn:") {
		return result, nil
	}

	// Check if project role exists in the current IAM
	roleMap, err := RoleMap()
	if err != nil {
		return LProject{}, err
	}
	newRole, ok := roleMap[result.Role]
	result.Role = newRole
	if !ok {
		return result, errors.New("Role Not found: " + result.Role)
	}

	return result, nil
}

// Upload a function to an S3 Bucket
func (lp LProject) UploadLambda(name string) error {
	fpath := path.Join(lp.path, name)

	os.Setenv("GOOS", "linux")
	os.Setenv("GOARCH", "amd64")

	_, err := run("go", "build", "-o", fpath, fpath+".go")
	if err != nil {
		return err
	}

	fmt.Println("Zipping to " + fpath + ".zip")
	_, err = run("zip", "-j", fpath+".zip", fpath)
	if err != nil {
		return err
	}

	lambdaName := lp.Name + "_" + name

	upCmd := exec.Command("aws", "s3", "cp", fpath+".zip", "s3://"+lp.Bucket+"/"+lambdaName+".zip")

	upOut, err := upCmd.StdoutPipe()
	if err != nil {
		return err
	}

	fmt.Println("Starting Upload of " + lambdaName)
	err = upCmd.Start()
	if err != nil {
		return err
	}

	io.Copy(os.Stdout, upOut)
	err = upCmd.Wait()
	if err != nil {
		return err
	}

	fl, err := NewFunctionList()
	if err != nil {
		return err
	}

	if fl.HasFunction(lambdaName) {
		resp, err := run("aws", "lambda", "update-function-code", "--function-name", lambdaName, "--s3-bucket", lp.Bucket, "--s3-key", lambdaName+".zip")
		if err != nil {
			return err
		}
		fmt.Println(string(resp))
		return nil
	}

	fmt.Println("Creating Function")
	res, err := run("aws", "lambda", "create-function", "--function-name", lambdaName, "--runtime", "go1.x", "--role", lp.Role, "--handler", name, "--code", "S3Bucket="+lp.Bucket+",S3Key="+lambdaName+".zip")
	if err != nil {
		return err
	}

	fmt.Println(string(res))

	return nil
}

func main() {
	lambdaName := flag.String("n", "", "Name of Lambda")
	configLocation := flag.String("c", "project.json", "Location of Config file")
	flag.Parse()

	project, err := NewLProject(*configLocation)
	if err != nil {
		log.Fatal(err)
	}

	err = project.UploadLambda(*lambdaName)
	if err != nil {
		log.Fatal(err)
	}
}

