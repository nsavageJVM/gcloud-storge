package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"
	// "github.com/aws/aws-sdk-go/service"
)

const (
	// This can be changed to any valid object name.
	objectName = "test-file"
	// This scope allows the application full control over resources in Google Cloud Storage
	scope = storage.DevstorageFullControlScope
	bucketName = "cloud-demo-eddy.appspot.com"
	projectID  = "cloud-demo-eddy"
)

var (

	fileName   = flag.String("file", "", "The file to upload.")

	command   = flag.String("cmd", "", "The command to run")
)




func main() {

	// Authentication is provided by the gcloud tool when running locally
	client, err := google.DefaultClient(context.Background(), scope)
	if err != nil {
		log.Fatalf("Unable to get default client: %v", err)
	}
	service, err := storage.New(client)
	if err != nil {
		log.Fatalf("Unable to create storage service: %v", err)
	}


	flag.PrintDefaults()
	flag.Parse()


		if *command == "" {
			fmt.Println("please enter a command")
			var newcmmd string

			if _, err := fmt.Scanf("%s", &newcmmd); err != nil {
				fmt.Printf("%s\n", err)
				return
			}

			switch newcmmd {

			case "list":
				listBuckets(service)

			case "load":
				if *fileName == "" {
					fmt.Println("please enter a file")
					var newfile string
					if _, err := fmt.Scanf("%s", &newfile); err != nil {
						fmt.Printf("%s\n", err)
						return
					}
					uploadFile(service, newfile)
				} else {
					uploadFile(service, "")
				}


			case "listfiles":
				listFiles(service)

			case "geturl":
				getRemoteUrl(service)

			}

		} else {
			switch *command {

			case "list":
				listBuckets(service)

			case "load":
				if *fileName == "" {
					log.Fatalf("File argument is required. See --help.")
				}
				uploadFile(service, "")

			case "listfiles":
				listFiles(service)

			case "geturl":
				getRemoteUrl(service)

			}
		}
 //}




	//// If the bucket already exists and the user has access, warn the user
	//if _, err := service.Buckets.Get(bucketName).Do(); err == nil {
	//	fmt.Printf("Bucket %s already exists - skipping buckets.insert call.", bucketName)
	//} else {
	//	fatalf(service, "Failed creating bucket %s: %v", bucketName, err)
	//}
	//





}

//<editor-fold defaultstate="collapsed"  desc="==  utility error function ==" >
func fatalf(service *storage.Service, errorMessage string, args ...interface{}) {
	// restoreOriginalState(service)
	log.Fatalf("Dying with error:\n"+errorMessage, args...)
}
//</editor-fold>

//<editor-fold defaultstate="collapsed"  desc="==  utility function  List all buckets  ==" >
func listBuckets(service *storage.Service) {

	if res, err := service.Buckets.List(projectID).Do(); err == nil {
		fmt.Println("Buckets:")
		for _, item := range res.Items {
			fmt.Println(item.Id)
		}
		fmt.Println()
	} else {
		fatalf(service, "Buckets.List failed: %v", err)
	}

}
//</editor-fold>


//<editor-fold defaultstate="collapsed"  desc="==  utility function  upload files  ==" >
func uploadFile(service *storage.Service, newFile string) {

	// Insert an object into a bucket.
	object := &storage.Object{Name: objectName}
	if *fileName == "" {
		file, err := os.Open(newFile)
		if err != nil {
			fatalf(service, "Error opening %q: %v", newFile, err)
		}
		if res, err := service.Objects.Insert(bucketName, object).Media(file).Do(); err == nil {
			fmt.Printf("Created object %v at location %v\n\n", res.Name, res.SelfLink)
		} else {
			fatalf(service, "Objects.Insert failed: %v", err)
		}

	} else {
		file, err := os.Open(*fileName)
		if err != nil {
			fatalf(service, "Error opening %q: %v", *fileName, err)
		}
		if res, err := service.Objects.Insert(bucketName, object).Media(file).Do(); err == nil {
			fmt.Printf("Created object %v at location %v\n\n", res.Name, res.SelfLink)
		} else {
			fatalf(service, "Objects.Insert failed: %v", err)
		}
	}
}
//</editor-fold>


//<editor-fold defaultstate="collapsed"  desc="==  utility function  list files  ==" >
func listFiles(service *storage.Service) {

	// List all objects in a bucket using pagination
	var objects []string
	pageToken := ""
	for {
		call := service.Objects.List(bucketName)
		if pageToken != "" {
			call = call.PageToken(pageToken)
		}
		res, err := call.Do()
		if err != nil {
			fatalf(service, "Objects.List failed: %v", err)
		}
		for _, object := range res.Items {
			objects = append(objects, object.Name)
		}
		if pageToken = res.NextPageToken; pageToken == "" {
			break
		}
	}

	fmt.Printf("Objects in bucket %v:\n", bucketName)
	for _, object := range objects {
		fmt.Println(object)
	}
	fmt.Println()

}
//</editor-fold>


//<editor-fold defaultstate="collapsed"  desc="==  utility function get files from a bucket.  ==" >
func getRemoteUrl(service *storage.Service) {
	// Get an object from a bucket.
	if res, err := service.Objects.Get(bucketName, objectName).Do(); err == nil {
		fmt.Printf("The media download link for %v/%v is %v.\n\n", bucketName, res.Name, res.MediaLink)
	} else {
		fatalf(service, "Failed to get %s/%s: %s.", bucketName, objectName, err)
	}
}
//</editor-fold>