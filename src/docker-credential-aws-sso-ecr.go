package main

import (
	"fmt"
	"os"
	"bufio"
	s "strings"
	"os/exec"
	"encoding/json"
)

func parseFile(file string) []string {
	f, err := os.Open(file)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing config file: %s", err)
	}

	defer f.Close()

	var lines []string
	
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning lines of config file: %s", err)
	}
	return lines
}

type awsProfile struct {
	account string
	region string
	name string
}

func getAwsConfigFile() string {
	var awsConfigFile string

	if os.Getenv("AWS_CONFIG_FILE") != "" {
		awsConfigFile = os.Getenv("AWS_CONFIG_FILE")
	} else {
		awsConfigFile = os.Getenv("HOME") + "/.aws/config"
	}

	return awsConfigFile
}

func getAwsSsoProfile(awsAccount string, awsRegion string) string{	
	// Parsing the config file to extract the profile name
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
	awsProfileName = ""
	for _, v := range profiles {
		if v.account == awsAccount && v.region == awsRegion {
			awsProfileName = v.name			
			break
		} 
	}

	if awsProfileName == "" {
		fmt.Fprintln(os.Stderr, "AWS SSO profile missing for region %v in account %v", awsRegion, awsAccount )
		fmt.Fprintln(os.Stderr, "Try to add it by running: aws configure sso")
		os.Exit(1)
	}

	return awsProfileName
}

func callAwsCli(profileName string, awsRegion string) string{
	cmd := exec.Command("aws", "--profile", profileName, "ecr", "get-login-password", "--region", awsRegion)
	stdout, err := cmd.Output()
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "AWS SSO login expired: %s\n", err)
		loginCmd := exec.Command("aws", "--profile", profileName, "sso", "login")
		loginErr := loginCmd.Run()
		if loginErr != nil {
			fmt.Fprintln(os.Stderr, "AWS SSO login expired: %s", loginErr)
			os.Exit(1)
		} else {
			return callAwsCli(profileName, awsRegion)
		}
	}
	return string(stdout)
}

func getCredentials(serverUrl string) {
	// Server url for AWS ECR looks like 123123.dkr.ecr.us-east-1.amazonaws.com
	var awsAccount = s.Split(serverUrl, ".")[0]	
	var awsRegion = s.Split(serverUrl, ".")[3]

	profileName := getAwsSsoProfile(awsAccount, awsRegion)

	fmt.Fprintln(os.Stderr, "Profile name used: ", profileName)
	
	secret := callAwsCli(profileName, awsRegion)

	response := map[string]string {"ServerURL": serverUrl, "Username": "AWS", "Secret": secret }
    responseJson, _ := json.Marshal(response)
	fmt.Fprintf(os.Stdout,string(responseJson))	
}

func main() {

	args := os.Args[1:]
	fmt.Fprintln(os.Stderr, "##########################################################################")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "       Logging to AWS SSO with docker-credential-aws-sso-ecr")
	fmt.Fprintln(os.Stderr, "       https://github.com/overbit/docker-credential-aws-sso-ecr")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "##########################################################################")

	var payload string 
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		payload = fmt.Sprintf(scanner.Text())
	}

	if (args[0] == "get") {
		// fmt.Println("Not supported. Only get is supported")
		getCredentials(payload)
	}

	// if (args[0] == "store") {
	// }

	// if (args[0] == "erase") {
	// }
	
}
