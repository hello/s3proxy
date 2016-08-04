package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type S3ProxyHandler struct {
	s3Client   *s3.S3
	bucketName string
}

func (h *S3ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "%v", err)
		return
	}

	key := time.Now().Format("2006-01-02-150405")
	input := &s3.PutObjectInput{
		Bucket: aws.String(h.bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
	}

	_, err = h.s3Client.PutObject(input)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "%v", err)
		return
	}
	fmt.Fprintf(w, "%s\n", key)
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("missing bucketname and port")
		return
	}

	bucketname := os.Args[1]
	port := os.Args[2]

	config := &aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials("AKIAJNV427WXSZAYVRXA", "7vArVV9VP7bZA/UeI8dSkRTpJ80ZzyPMHDg7mRFG", ""),
	}

	handler := &S3ProxyHandler{
		s3Client:   s3.New(session.New(), config),
		bucketName: bucketname,
	}
	err := http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatal(err)
	}
}
