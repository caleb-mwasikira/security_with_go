package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/joho/godotenv"
)

var (
	secret_key    string
	password_file string = "test_data/shadow"
)

func generateSalt() (string, error) {
	random_bytes := make([]byte, 32)
	_, err := rand.Read(random_bytes)
	if err != nil {
		return "", err
	}
	salt := hex.EncodeToString(random_bytes)
	return salt, nil
}

func hashPassword(plain_text, salt string) string {
	hash := hmac.New(sha256.New, []byte(secret_key))
	io.WriteString(hash, plain_text+salt)
	hashed_value := hash.Sum(nil)
	return hex.EncodeToString(hashed_value)
}

func inputStr(input_type, prompt string, min_len int) string {
	var value string

	for {
		fmt.Print(prompt)
		fmt.Scanf("%s", &value)

		if len(value) > min_len {
			break
		}
		fmt.Printf("%v cannot be less than %v characters: please try again\n\n", input_type, min_len)
	}
	return value
}

/*
asks the user to enter their username and password
*/
func getLoginCredentials() (string, string) {
	min_uname_len := 3
	min_pass_len := 8
	username := inputStr("username", "Enter your username: ", min_uname_len)
	password := inputStr("password", "Enter your password: ", min_pass_len)

	return username, password
}

/*
saves user login credentials to password file in the
format:

	username:$id$salt$hashed
*/
func savePassword(username, salt, hashed_password string) error {
	file, err := os.OpenFile(password_file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// $5$ is sha-256 hashing algorithm
	// https://www.cyberciti.biz/faq/understanding-etcshadow-file/
	hashing_alg := 5
	line := fmt.Sprintf("%v:$%v$%v$%v\n", username, hashing_alg, salt, hashed_password)
	_, err = file.WriteString(line)
	return err
}

/*
search for username in password file
*/
func findUserAccount(username string) (string, error) {
	file, err := os.OpenFile(password_file, os.O_RDONLY, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// read file line-by-line until we find the user or not
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, username) {
			return line, nil
		}
	}
	return "", fmt.Errorf("user account with username '%v' not found", username)
}

func main() {
	// load secret key from env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file; %v", err)
	}
	secret_key = os.Getenv("SECRET_KEY")

	// select user action
	user_action := 0
	valid_user_actions := []int{1, 2, 3}

	for {
		fmt.Println("Select action:\n   1. SignUp\n   2. Login\n   3. Exit")
		fmt.Print("> ")
		fmt.Scanf("%d", &user_action)

		if slices.Contains(valid_user_actions, user_action) {
			break
		}
		fmt.Printf("incorrect choice %v\n", user_action)
	}
	println()

	switch user_action {
	case 1: // signup
		fmt.Println("creating a new account")
		username, password := getLoginCredentials()

		// first check if username already exists
		user_account, _ := findUserAccount(username)
		if len(user_account) > 0 {
			log.Fatalf("account with username %v already exists", username)
		}

		salt, err := generateSalt()
		if err != nil {
			log.Fatalf("error generating salt; %v", err)
		}

		hashed_password := hashPassword(password, salt)
		err = savePassword(username, salt, hashed_password)
		if err != nil {
			log.Fatalf("error saving password to file; %v", err)
		}
		fmt.Println("account created successfuly")

	case 2: // login
		fmt.Println("logging in to your existing account")
		username, password := getLoginCredentials()
		user_account, err := findUserAccount(username)
		if err != nil {
			log.Fatalf("error logging in: %v", err)
		}
		id_salt_hash := strings.TrimPrefix(user_account, fmt.Sprintf("%v:", username))
		id_salt_hash = strings.TrimPrefix(id_salt_hash, "$")
		fields := strings.Split(id_salt_hash, "$")
		salt := fields[1]
		db_password_hash := fields[2]

		// verify password match
		hashed_password := hashPassword(password, salt)
		if hashed_password == db_password_hash {
			fmt.Println("Login success")
		} else {
			fmt.Println("Login failed")
		}
	case 3:
		fmt.Println("exiting program... Bye :)")
		os.Exit(1)
	}
}
