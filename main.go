package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hashicorp/vault/api"
	"strings"
	"os"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/logical"
	"fmt"
)

func main() {
	if (len(os.Args) <= 1) {
		usage()
	}
	pkcs7, instanceProfile := getEc2Info()
	vaultClient, _ := api.NewClient(nil)
	vaultLogin(vaultClient, pkcs7, instanceProfile)
	secretKey := os.Args[1]
	vaultReadSecret(vaultClient, secretKey)
}

func vaultReadSecret(vaultClient *api.Client, secretKey string) {
	secret, err := vaultClient.Logical().Read("secret/" + secretKey)
	if err != nil {
		log.Fatal(err)
	}
	if secret == nil {
		log.Fatalf("Key %s not found\n", secretKey)
	}
	if val, ok := secret.Data["secret"]; ok {
		fmt.Print(val)
	}
}

func vaultLogin(vaultClient *api.Client, pkcs7 string, instanceProfile string) {
	req := vaultClient.NewRequest("POST", "/v1/auth/aws-ec2/login")
	req.SetJSONBody(map[string]string{
		"role": instanceProfile,
		"pkcs7": pkcs7,
		"nonce":"vault-client-nonce",
	})
	resp, err := vaultClient.RawRequest(req)

	var response logical.Response
	err = jsonutil.DecodeJSONFromReader(resp.Body, &response)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Fatalf("Login failed with instnace profile=%s\n", instanceProfile)
	}

	if err != nil {
		log.Println(err)
	}
	log.Debug("ClientToken: " + response.Auth.ClientToken)
	vaultClient.SetToken(response.Auth.ClientToken)
}

func getEc2Info() (string, string) {
	ec2MetadataClient := ec2metadata.New(session.New())
	pkcs7, _ := ec2MetadataClient.GetDynamicData("instance-identity/pkcs7")
	if pkcs7 == "" {
		log.Fatalln("Unable to get PKCS7, the instance is not on AWS")
	}
	iamInfo, _ := ec2MetadataClient.IAMInfo()

	// arn:aws:iam::$acct_number:instance-profile/$role_name
	log.Debug(iamInfo.InstanceProfileArn)
	instanceProfile := strings.SplitN(iamInfo.InstanceProfileArn, "/", 2)[1]
	log.Debug(instanceProfile)
	return pkcs7, instanceProfile
}

func usage() {
	fmt.Printf("Usage: %s secret_key\n", os.Args[0])
	os.Exit(1)
}
