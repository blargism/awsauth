package aws

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
)

var newLine string

type AWSConfig struct {
	ProfileName string
	Values      map[string]string
}

func init() {
	os := runtime.GOOS
	switch os {
	case "windows":
		newLine = "\r\n"
	case "darwin":
		newLine = "\r"
	case "linux":
		newLine = "\r"
	default:
		newLine = "\r"
	}
}

func writeConfig(fileName string, config *AWSConfig) {
	home, err := os.UserHomeDir()

	if err != nil {
		log.Fatalln("could not access home directory")
	}

	if _, e := os.Stat(home + "/.aws"); e != nil {
		if os.IsNotExist(e) {
			if err := os.Mkdir(home+"/.aws", 0755); err != nil {
				log.Fatalln("could not create .aws directory")
			}
		} else {
			log.Fatalln("could not find or create .aws folder")
		}
	}

	filePath := home + "/.aws/" + fileName

	contents := "[" + config.ProfileName + "]" + newLine

	for k, v := range config.Values {
		contents += k + " = " + v + newLine
	}

	buf := []byte(contents)
	err = ioutil.WriteFile(filePath, buf, 644)

	if err != nil {
		log.Fatalf("Could not write to file: %s\n", filePath)
	}

	fmt.Println("new credentials written to ~/.aws/credentials")
}
