package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"hackathon/helpers"
	"log"
	"net/http"
	"os"
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

	jsonValue, err := json.Marshal(map[string]string{"url": up.Location})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	faasUrl := os.Getenv("URL")
	faasApiKey := os.Getenv("API_KEY")
	faasReq, err := http.NewRequest("POST", faasUrl, bytes.NewBuffer(jsonValue))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	faasReq.Header.Set("Content-Type", "application/json")
	faasReq.Header.Set("api-key", faasApiKey)

	client := &http.Client{}
	resp, err := client.Do(faasReq)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer resp.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	jData, err := json.Marshal(up.Location)
	if err != nil {
		// handle error
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}
func actions(w http.ResponseWriter, req *http.Request) {

	var video helpers.Videos
	var url string
	json.NewDecoder(req.Body).Decode(&url)
	fmt.Print(req.Body)
	helpers.DB.Where("name = ?", &url).First(&video)
	jData, err := json.Marshal(video.Action)
	if err != nil {
		// handle error
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func main() {
	dbDns := os.Getenv("DB_DNS")
	if dbDns == "" {
		log.Fatal(1)
	}
	helpers.SetupConnection(dbDns)

	http.Handle("/", http.FileServer(http.Dir("../")))
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/actions", actions)

	http.ListenAndServe(":8090", nil)
}
