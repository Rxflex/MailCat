# Email Classification with Gemini

This project represents a system for the automatic classification of incoming emails using Google's Gemini model. It connects to an IMAP server, processes emails, classifies them, and moves them to appropriate folders based on their content.

## Donations

If you find this project useful, you can support it through Revolut: https://revolut.me/chuk1.

## Dependencies

For proper functioning of the project, the following dependencies need to be installed:

- [go-imap](https://github.com/BrianLeishman/go-imap) — a library for working with IMAP servers.
- [generative-ai-go](https://github.com/google/generative-ai-go) — a library for working with Google Gemini.
- [google-api-go-client](https://github.com/googleapis/google-api-go-client) — a client for Google APIs.

### Installing Dependencies

To install the dependencies, use [Go Modules](https://blog.golang.org/using-go-modules).

1. Initialize Go Modules if not already done:

   ```bash
   go mod init email-classification
   ```

2. Install the required packages:

   ```bash
   go get github.com/BrianLeishman/go-imap
   go get github.com/google/generative-ai-go/genai
   go get google.golang.org/api/option
   ```

## Setup

1. **Optional**: You can create a `credentials.json` file in the root directory of the project, which will contain your login and password if you wish to save credentials for future use.

   File format:

   ```json
   {
       "login": "your-email",
       "password": "your-password"
   }
   ```

2. If you don't want to use the `credentials.json` file, the program will prompt you to enter login and password every time you run it.

3. Fill in the `apiKey` variable in the code with your Google Gemini API key.

4. Modify the IMAP server settings if necessary.

## Usage

1. Run the project with the following command:

   ```bash
   go run main.go
   ```

2. The program will prompt you to enter your login and password, and then it will begin processing emails.

3. The program will classify emails and move them to appropriate folders.

   Example program output:

   ```
   Connecting to IMAP...
   Login successful. Starting email processing...
   Getting folders...
   Folder: INBOX
   Processing email: Hello, World!
   Classifying email... Subject: Hello, World! Senders: sender@example.com
   Email classified as: INBOX
   Email not moved to folder: INBOX
   ```

## Screenshots

**Screenshot 1**: Screen for saving credentials  
![remember.png](assets%2Fremember.png)

**Screenshot 2**: Login screen  
![login.png](assets%2Flogin.png)

**Screenshot 3**: Email processing and classification process  
![result.png](assets%2Fresult.png)
