# Codeforces Submission Tracker Bot

This project is a Telegram bot that tracks Codeforces submissions and sends daily reports of solved problems by users. The bot fetches submission data from the Codeforces API, analyzes the results, and sends a summary report via Telegram. It supports commands to trigger an immediate report and can be scheduled to send daily updates at a specific time.

## Features

- **Fetch Submission Data**: The bot fetches user submission data from Codeforces using the API.
- **Track Solved Problems**: It tracks solved problems for each user and reports them daily.
- **Telegram Integration**: Sends daily reports and updates directly to a Telegram chat.
- **Automatic Updates**: The bot can be scheduled to send reports at a specific time daily.
- **Customizable User List**: Users are read from a file, and their submissions are fetched and analyzed.
- **Markdown Formatting**: The report is formatted with Markdown for better readability in Telegram.

## Requirements

- **Go 1.18+**
- **Telegram Bot Token**: Obtain a bot token by creating a new bot on Telegram via the BotFather.
- **Codeforces Usernames**: A text file with a list of Codeforces usernames (e.g., `user_list.txt`).

## Setup

1. **Install Dependencies**:  
   - Run `go get` to install necessary dependencies:
    ```bash
    go get github.com/go-co-op/gocron
2. **Configuration**:  
   - Open the `main.go` file.
   - Set the `bot_token` with your Telegram bot token.
   - Set the `chat_id` with the desired Telegram chat ID where the reports should be sent.
   - Prepare the `user_list.txt` file containing the list of Codeforces usernames.
3. **Schedule Daily Reports**
   - The bot will send a daily report at 08:24 AM (Asia/Dhaka time). You can change the time by modifying the following line in the main() function:
     ```bash
     scheduler.Every(1).Day().At("08:24").Do(func() {...}) 

4. **Running the Bot**
   - To run the bot, use the following command:
     ```bash
     go run main.go

## Output Example

Here is an example of the daily report that the bot sends to the Telegram chat:

<img src="https://i.imgur.com/8tEH1Uh.jpeg" alt="Telegram Bot Report" width="400">

### Sample Telegram Bot Report

The following is an example of a formatted report sent by the bot, showing the solved problems for the day:

