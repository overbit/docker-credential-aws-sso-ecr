package main

import (
	"fmt"
	"os"
	"bufio"
	"log"
	s "strings"
	"os/exec"
)

func getAwsConfigFile() string {
	var awsConfigFile string

	if os.Getenv("AWS_CONFIG_FILE") != "" {
		awsConfigFile = os.Getenv("AWS_CONFIG_FILE")
	} else {
		awsConfigFile = os.Getenv("HOME") + "/.aws/config"
	}

	return awsConfigFile
}

func parseFile(file string) []string {
	f, err := os.Open(file)

	if err != nil {
		log.Fatal("")
		log.Fatal(err)
	}

	defer f.Close()

	var lines []string
	
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return lines
}

type awsProfile struct {
	account string
	region string
	name string
}

func getAwsSsoProfile(awsAccount string, awsRegion string) string{	

	// An aws config file for SSO enabled accounts looks like
	//    [profile myprofile]
	//    sso_start_url = https://myurl
	//    sso_region = myregion
	//    sso_account_id = 123123
	//
	//    [profile myotherprofile]
	//    sso_start_url = https://myurl
	//    sso_region = myotherregion
	//    sso_account_id=456456
	lines := parseFile(getAwsConfigFile())
	
	var profiles []awsProfile

	var currentProfile string
	currentProfileIndex := -1
	for i := 0; i < len(lines); i++ {
		if s.Contains(lines[i], "[profile ") || s.Contains(lines[i], "[default]") {
			currentProfile = s.Replace(lines[i], "[profile ", "", -1)
			currentProfile = s.Replace(currentProfile, "]", "", -1)
			profiles = append(profiles, awsProfile{name: currentProfile})
			currentProfileIndex ++
		} else if s.Contains(lines[i], "sso_region = ") {
			profiles[currentProfileIndex].region = s.Replace(lines[i], "sso_region = ", "", -1)
		} else if s.Contains(lines[i], "sso_account_id =") {
			profiles[currentProfileIndex].account = s.Replace(lines[i], "sso_account_id = ", "", -1)
		}
	}

	var awsProfileName string
	for _, v := range profiles {
		if v.account == awsAccount && v.region == awsRegion {
			awsProfileName = v.name			
			break
		} 
	}

	// fmt.Println("AWS profile name: %s\n", awsProfileName)
	return awsProfileName
}

func getCredentials(serverUrl string) {
	// Server url for AWS ECR looks like 123123.dkr.ecr.us-east-1.amazonaws.com
	var awsAccount = s.Split(serverUrl, ".")[0]
	// fmt.Printf("AWS account: %s\n", awsAccount)
	
	var awsRegion = s.Split(serverUrl, ".")[3]
	// fmt.Printf("AWS region: %s\n", awsRegion)

	profileName := getAwsSsoProfile(awsAccount, awsRegion)

	cmd := exec.Command("aws", "--profile", profileName, "ecr", "get-login-password", "--region", awsRegion)
	stdout, err := cmd.Output()
	
	if err != nil {
		exec.Command("aws", "--profile", profileName, "sso", "login")
	}

	fmt.Printf("{\"ServerURL\": \"%s\", \"Username\": \"AWS\", \"Secret\": \"%s\" }\n", serverUrl, s.Replace(string(stdout), "\n", "", -1))
}

func main() {

	args := os.Args[1:]

	if (args[0] == "get") {
		// fmt.Println("Not supported. Only get is supported")
		getCredentials(args[1])
	}
	
}
