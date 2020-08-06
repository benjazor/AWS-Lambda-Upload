# AWS-Lambda-Upload
This is just a little project that I made during my Udemy course "Hands-on Serverless Computing with Go".  
It's a tool which simplifies to the process of uploading lambda functions to an S3 Bucket.
## How to use
### Go build
```shell script
go build -o AWS-Lambda-Upload .
./AWS-Lambda-Upload -c /path/to/project.json -n name-of-lambda-function 
```

### Go run
```shell script
go run *.go -c /path/to/project.json -n name-of-lambda-function 
```

### Arguments
- -c : config file path 
- -n : lambda function name (without the .go and must be in the smae directory as config)


## Config file structure
project.json:
```json
{
  "Name": "nameOfProject",
  "Bucket": "nameOfS3Bucket",
  "Role": "IAMRole"
}
```