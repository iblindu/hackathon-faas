package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

var Sess = ConnectAws()

func ConnectAws() *session.Session {
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	endpoint := os.Getenv("ENDPOINT")
	sess, err := session.NewSession(
		&aws.Config{
			Endpoint: aws.String(endpoint),
			Region:   aws.String("eu-central-1"),
			Credentials: credentials.NewStaticCredentials(
				accessKeyID,
				secretAccessKey,
				"", // a token will be created when the session it's used.
			),
		})
	if err != nil {
		panic(err)
	}
	return sess
}

func upload(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "upload \n")
	err := req.ParseMultipartForm(32 << 20) // maxMemory 32MB
	if err != nil {
		fmt.Printf("parser %v", err)
		return
	}
	file, handler, err := req.FormFile("myfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Fprintf(w, "Handler %v", handler.Header)
	uploader := s3manager.NewUploader(Sess)
	bucket := os.Getenv("BUCKET_NAME")
	//upload to the s3 bucket
	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		ACL:    aws.String("public-read"),
		Key:    aws.String(handler.Filename),
		Body:   file,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write([]byte(up.Location))
}
func actions(w http.ResponseWriter, req *http.Request) {

	videos := make(map[string]string)
	videos["name"] = "action"
	jData, err := json.Marshal(videos)
	if err != nil {
		// handle error
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)

}

func main() {
	http.Handle("/", http.FileServer(http.Dir("../")))
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/actions", actions)

	http.ListenAndServe(":8090", nil)
}
