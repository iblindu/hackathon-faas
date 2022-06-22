package main

import (
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
	uploader := s3manager.NewUploader(Sess)
	bucket := os.Getenv("BUCKET_NAME")
	file, err := os.Open("testfile")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer file.Close()
	//upload to the s3 bucket
	up, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		ACL:    aws.String("public-read"),
		Key:    aws.String("testfile"),
		Body:   file,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write([]byte(up.Location))
}
func actions(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "actions\n")

}

func main() {
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/actions", actions)

	http.ListenAndServe(":8090", nil)
}
