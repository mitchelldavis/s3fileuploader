package main

import (
    "log"
    "flag"
    "os"
	"path"
    "path/filepath"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
    "github.com/mitchelldavis/s3fileupload/pkg/cognitoprovider"
    "github.com/mitchelldavis/s3fileupload/pkg/progressreader"
)

var regionArg, usernameArg, clientIdArg, userPoolIdArg, identityPoolIdArg string
var filenameArg, bucketArg string

func init() {
    flag.StringVar(&usernameArg, "username", "", "The Username to sign in with.")
    flag.StringVar(&regionArg, "region", "", "The AWS region the user and identity pools exist in.")
    flag.StringVar(&clientIdArg, "clientId", "", "The Client Id to pass to the identity provider.")
    flag.StringVar(&userPoolIdArg, "userPoolId", "", "The User Pool Id that has the user information we're authenticating against.")
    flag.StringVar(&identityPoolIdArg, "identityPoolId", "", "The Identity Pool Id that we're authenticating against.")
    flag.StringVar(&filenameArg, "filename", "", "The File to upload.")
    flag.StringVar(&bucketArg, "bucket", "", "The Bucket to upload the file to.")
    flag.Parse()
}

type Configuration struct {
    Region string `json:"region"`
    ClientId string `json:"clientid"`
    UserPoolId string `json:"userpoolid"`
    IdentityPoolId string `json:"identitypoolid"`
    Username string `json:"username"`
    Bucket string `json:"bucket"`
}

func validateArgs() {
    if filenameArg == "" {
        log.Fatal("You must supply a filename to upload.")
    } 

    if usernameArg == "" {
        log.Fatal("You must supply a username.")
    }

    if regionArg == "" {
        log.Fatal("You must supply a region.")
    }

    if clientIdArg == "" {
        log.Fatal("You must supply an id for the client you're using.")
    }

    if userPoolIdArg == "" {
        log.Fatal("You must supply a User Pool Id.")
    }

    if identityPoolIdArg == "" {
        log.Fatal("You must supply an Identity Pool Id.")
    }

    if bucketArg == "" {
        log.Fatal("You must supply a bucket to upload the file too.")
    }
}

func main() {

    validateArgs()

    cp := cognitoprovider.New(&regionArg, &usernameArg, &clientIdArg, &userPoolIdArg, &identityPoolIdArg)

    sess, err := session.NewSession(&aws.Config{
        Region: aws.String(regionArg),
        Credentials: credentials.NewCredentials(cp),
    })

    uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024
	})

    f, err  := os.Open(filenameArg)
    if err != nil {
        log.Fatalf("failed to open file %q, %v", filenameArg, err)
    }
    defer f.Close()

    fileInfo, err := f.Stat()
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

    pr := &progressreader.ProgressReader{
        File: f,
        Size: fileInfo.Size(),
        Progress: func(size, read int64) {
			// I have no idea why the read length need to be div 2,
			// maybe the request read once when Sign and actually send call ReadAt again
			// It works for me
            log.Printf("total read:%d    progress:%d%%\n", read/2, int(float32(read*100/2)/float32(size)))
        },
    }

    // Upload the file to S3.
    result, err := uploader.Upload(&s3manager.UploadInput{
        Bucket: aws.String(bucketArg),
        Key:    aws.String(path.Join(usernameArg, filepath.Base(filenameArg))),
        //Key:    aws.String(filepath.Base(filenameArg)),
        Body:   pr,
    })
    if err != nil {
        log.Fatalf("failed to upload file, %v", err)
    }
    log.Println("file uploaded to ", result.Location)

}
