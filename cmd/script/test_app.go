package main

import (
	"bytes"
	"fmt"
	"github.com/goiiot/libmqtt/cmd/utils"
	"github.com/goiiot/libmqtt/common"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

// CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o test_app test_app.go
// nohup sh test.sh >> log.out 2>&1 &
func main() {
	args := os.Args[1:]
	for i, arg := range args {
		fmt.Printf("参数 %d:%s \n", i, arg)
	}
	if len(args) == 3 {
		fmt.Printf("start app")
		file := args[0]
		url := args[1]
		token := args[2]
		files := []utils.UploadFile{
			{Name: "file", Filepath: file},
		}
		headers := make(map[string]string)
		reqs := make(map[string]string)
		headers["Authorization"] = token
		resp := utils.PostFile(url, reqs, files, headers)
		fmt.Printf("sucesss:%s", resp)
	} else {
		fmt.Printf("test success")
	}
}
func postFile(filename string, targetUrl string, token string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("file", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	//打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}
	defer func(fh *os.File) {
		_ = fh.Close()
		_ = bodyWriter.Close()
	}(fh)

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}

	contentType := bodyWriter.FormDataContentType()

	req, err := http.NewRequest("POST", targetUrl, bodyBuf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", token)
	resp, err := common.Client.Do(req)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	fmt.Println(string(respBody))
	return nil
}
