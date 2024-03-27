package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/oracle/oci-go-sdk/v49/common"
	"github.com/oracle/oci-go-sdk/v49/objectstorage"
)

func Upload(filebytes []byte, fileName *string) (*string, error) {

	cfg, _ := common.ConfigurationProviderFromFile("/home/dojeto/.oci/config", "DEFAULT")

	client, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(cfg)
	if err != nil {
		fmt.Println("Error creating client:", err)
		return nil, err
	}
	// /home/ubuntu/.oci

	namespace := os.Getenv("NAMESPACE")

	bucketName := os.Getenv("BUCKETNAME")
	err = os.WriteFile("temp.pdf", filebytes, 0644)
	if err != nil {
		fmt.Println("Error writing PDF file:", err)
		return nil, err
	}
	wd, _ := os.Getwd()
	filePath := wd + "/temp.pdf"

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	putObjectRequest := objectstorage.PutObjectRequest{
		NamespaceName: common.String(namespace),
		BucketName:    common.String(bucketName),
		ObjectName:    common.String(*fileName),
		ContentType:   common.String("application/pdf"),
		PutObjectBody: file,
	}

	// Upload the file
	_, err = client.PutObject(context.Background(), putObjectRequest)
	if err != nil {
		fmt.Println("Error uploading file:", err)
		return nil, err
	}

	fmt.Println("File uploaded successfully.")
	expirationTime := time.Now().Add(24 * time.Hour)
	parRequest := objectstorage.CreatePreauthenticatedRequestRequest{
		CreatePreauthenticatedRequestDetails: objectstorage.CreatePreauthenticatedRequestDetails{
			Name:        common.String("hahathisismyfirst"),
			ObjectName:  fileName,
			AccessType:  objectstorage.CreatePreauthenticatedRequestDetailsAccessTypeObjectread,
			TimeExpires: &common.SDKTime{Time: expirationTime}, // Set expiration time (24 hours from now)
		},
		NamespaceName: &namespace,
		BucketName:    &bucketName,
	}

	// Send the PAR request
	parResponse, err := client.CreatePreauthenticatedRequest(context.Background(), parRequest)
	if err != nil {
		fmt.Println("Error creating pre-authenticated request:", err)
		return nil, err
	}

	// Print the PAR URL
	fmt.Println("Direct access link for the PDF:", *parResponse.PreauthenticatedRequest.AccessUri)

	os.Remove("temp.pdf")

	return parResponse.PreauthenticatedRequest.AccessUri, nil

}
