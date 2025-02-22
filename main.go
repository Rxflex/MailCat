package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/BrianLeishman/go-imap"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
)

const (
	apiKey      = "API_KEY"
	imapServer  = "MAIL_SERVER"
	imapPort    = 993
	inboxFolder = "INBOX"
	shortAlias  = "SHORT_ALIAS"
)

const requestDelay = 3 * time.Second

var categories = []string{"INBOX", "INBOX.Spam"}
var categoriesString = strings.Join(categories, ", ")

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func getCredentials() (string, string) {
	var creds Credentials
	credsFile := "credentials.json"
	reader := bufio.NewReader(os.Stdin)

	if _, err := os.Stat(credsFile); err == nil {
		data, err := ioutil.ReadFile(credsFile)
		if err == nil {
			err := json.Unmarshal(data, &creds)
			if err == nil {
				fmt.Printf("Use saved account: %s? (y/n): ", creds.Login)
				choice, _ := reader.ReadString('\n')
				choice = strings.TrimSpace(choice)
				if choice == "y" {
					return creds.Login, creds.Password
				}
			}
		}
	}

	fmt.Print("Enter login (without " + shortAlias + ", e.g., 'dev'): ")
	login, _ := reader.ReadString('\n')
	login = strings.TrimSpace(login)

	if !strings.Contains(login, "@") {
		login += shortAlias
	}

	fmt.Print("Enter password: ")
	bytePassword, _ := reader.ReadString('\n')
	password := strings.TrimSpace(bytePassword)

	creds.Login = login
	creds.Password = password
	data, _ := json.Marshal(creds)
	err := ioutil.WriteFile(credsFile, data, 0644)
	check(err)

	return login, password
}

func classifyEmail(emailSubject string, emailSender string, emailContent string) (string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx, genai.Text("Here is the subject of the letter, your task is to answer me in one word the type of this letter. Available categories are: "+categoriesString+". Mail sender is: "+emailSender+". The email subject is: "+emailSubject+". The email content is: "+emailContent))

	if err != nil {
		return "", err
	}

	printResponse(resp)

	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		return fmt.Sprintf("%+v", resp.Candidates[0].Content.Parts[0]), nil
	}

	return "Other", nil
}

func printResponse(resp *genai.GenerateContentResponse) {
	fmt.Printf("Response from Gemini: %+v\n", resp)

	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				printStruct(part)
			}
		}
	}
	fmt.Println("---")
}

func printStruct(s interface{}) {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	if v.Kind() == reflect.String {
		fmt.Printf("String: %s\n", v.String())
		return
	}

	fmt.Println("Type:", t)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := t.Field(i).Name
		fmt.Printf("%s: %v\n", fieldName, field.Interface())
	}
}

func mapCategoryToFolder(category string) string {
	category = strings.TrimSpace(strings.ToLower(category))

	for _, cat := range categories {
		if strings.ToLower(strings.TrimSpace(cat)) == category {
			return cat
		}
	}

	return "INBOX"
}

func main() {
	imap.Verbose = false
	imap.RetryCount = 3

	login, password := getCredentials()

	log.Println("Connecting to IMAP...")

	im, err := imap.New(login, password, imapServer, imapPort)
	check(err)
	defer im.Close()

	log.Println("Login successful. Starting email processing...")

	log.Println("Getting folders...")

	folders, err := im.GetFolders()
	check(err)

	for _, folder := range folders {
		log.Println("Folder: ", folder)
	}

	err = im.SelectFolder(inboxFolder)
	check(err)

	uids, err := im.GetUIDs("ALL")
	check(err)

	if len(uids) == 0 {
		log.Println("No new emails.")
		return
	}

	const batchSize = 50
	totalMessages := len(uids)
	for start := 1; start <= totalMessages; start += batchSize {
		end := start + batchSize - 1
		if end > totalMessages {
			end = totalMessages
		}

		emails, err := im.GetEmails(uids...)
		check(err)

		for _, email := range emails {
			time.Sleep(requestDelay)
			log.Println("Processing email:", email.Subject)
			var senders []string
			for _, addr := range email.From {
				senders = append(senders, addr)
			}
			emailSenders := strings.Join(senders, ",")
			log.Println("Classifying email... Subject:", email.Subject, "Senders:", emailSenders)
			category, err := classifyEmail(email.Subject, emailSenders, email.HTML)
			check(err)
			log.Println("Email classified as:", category)

			folder := mapCategoryToFolder(category)

			if folder != "INBOX" {
				err = im.MoveEmail(email.UID, folder)
				check(err)
				log.Println("Email moved to folder:", folder)
			} else {
				log.Println("Email not moved to folder:", folder)
			}
		}
	}

	log.Println("All emails processed.")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
