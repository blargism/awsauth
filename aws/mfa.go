package aws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type MFAResponse struct {
	Credentials struct {
		SecretAccessKey string `json:"SecretAccessKey"`
		SessionToken    string `json:"SessionToken"`
		Expiration      string `json:"Expiration"`
		AccessKeyId     string `json:"AccessKeyId"`
	} `json:"Credentials"`
}

func LoginWithMFA(accessID string, accessKey string, mfaDevice string, code string) {
	fmt.Printf("running session auth with mfa code with device %s\n", mfaDevice)
	cmd := exec.Command("aws")
	cmd.Args = []string{"", "--profile", "default", "sts", "get-session-token", "--output=json", "--serial-number", mfaDevice, "--token-code", code}
	fmt.Printf("running: %s\n", cmd.String())
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "AWS_PROFILE=default")
	cmd.Env = append(cmd.Env, "AWS_ACCESS_ID="+accessID)
	cmd.Env = append(cmd.Env, "AWS_ACCESS_KEY="+accessKey)

	writeConfig("credentials", &AWSConfig{
		ProfileName: "default",
		Values: map[string]string{
			"aws_access_key_id":     accessID,
			"aws_secret_access_key": accessKey,
		},
	})

	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut

	err := cmd.Run()

	if err != nil {
		fmt.Println(errOut.String())
		log.Fatalf("Failed to run aws mfa auth: \n%s\n", err)
	}

	res := &MFAResponse{}
	json.Unmarshal(out.Bytes(), res)

	fmt.Printf("created new session token expiring %s\n", res.Credentials.Expiration)

	values := map[string]string{
		"aws_access_key_id":     res.Credentials.AccessKeyId,
		"aws_secret_access_key": res.Credentials.SecretAccessKey,
		"aws_session_token":     res.Credentials.SessionToken,
	}
	config := &AWSConfig{
		ProfileName: "default",
		Values:      values,
	}

	writeConfig("credentials", config)
}
