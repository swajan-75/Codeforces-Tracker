package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"github.com/go-co-op/gocron"
)

var bot_token string = "7166228483:AAGD2P3z0o004YCT9jPMTz_EogX3zBcMEo8"
var chat_id string = "1274939394"
const codeforces_api = "https://codeforces.com/api/user.status"

type Problem struct {
	ContestID int      `json:"contestId"`
	Index     string   `json:"index"`
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	Points    float64  `json:"points"`
	Rating    int      `json:"rating"`
	Tags      []string `json:"tags"`
}

type Submission struct {
	ID             int    `json:"id"`
	ContestID      int    `json:"contestId"`
	Index          string `json:"index"`
	Problem       Problem `json:"problem"`
	Verdict        string `json:"verdict"`
	CreationTime   int64  `json:"creationTimeSeconds"`
	
}

type APIResponse struct {
	Status string       `json:"status"`
	Result []Submission `json:"result"`
}
func read_user_name(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var usernames []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		username := strings.TrimSpace(scanner.Text())
		if username != "" {
			usernames = append(usernames, username)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return usernames, nil
}
func get_data(handle string, from, count int) ([]Submission, error) {
	api_url := fmt.Sprintf("%s?handle=%s&from=%d&count=%d", codeforces_api, handle, from, count)
	resp, err := http.Get(api_url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch submissions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if apiResp.Status != "OK" {
		return nil, fmt.Errorf("API error: %s", apiResp.Status)
	}

	return apiResp.Result, nil
}
func listen_command() {
	api_url_ := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates", bot_token)
	offset := 0

	for {
		resp, err := http.Get(fmt.Sprintf("%s?offset=%d", api_url_, offset))
		if err != nil {
			fmt.Println("Error fetching updates:", err)
			time.Sleep(5 * time.Second)
			continue
		}
		defer resp.Body.Close()

		var updates struct {
			Result []struct {
				UpdateID int `json:"update_id"`
				Message  struct {
					Text string `json:"text"`
				} `json:"message"`
			} `json:"result"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&updates); err != nil {
			fmt.Println("Error decoding updates:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, update := range updates.Result {
			offset = update.UpdateID + 1

			if strings.ToLower(update.Message.Text) == "update" {
				report := generate_daily_report()
				if err := sendMessage(report); err != nil {
					fmt.Println("Error sending update report:", err)
				} else {
					fmt.Println("Update report sent successfully ")
				}
			}
		}

		time.Sleep(2 * time.Second)
	}
}

func generate_daily_report() string {
	usernames, err := read_user_name("user_list.txt")
	if err != nil {
		fmt.Println("Error reading usernames:", err)
		return ""
	}

	report_message := fmt.Sprintf("*📅 Date:* %s\n\n", escapeMarkdownV2(time.Now().Format("2006-01-02")))

	for _, username := range usernames {
		count := 100
		submissions, err := get_data(username, 1, count)
		if err != nil {
			fmt.Printf("Error fetching submissions for user %s: %v\n", username, err)
			continue
		}

		total_solve, solved_problem := solve_count(submissions)

		user_report := fmt.Sprintf("*👨 User:* %s\n*✅ Solved Today:* %d\n\n",
			escapeMarkdownV2(username),
			total_solve)

		if len(solved_problem) > 0 {
			user_report += "*📊 Solved Problems:*\n"
			for _, problem := range solved_problem {
				user_report += fmt.Sprintf("🔓 %s\n", problem)
			}
			user_report += "⏤⏤⏤⏤⏤⏤⏤⏤⏤⏤⏤⏤⏤⏤⏤\n"
			report_message += user_report
		}

		time.Sleep(1000 * time.Millisecond)
	}

	return report_message
}
func sendMessage(message string) error {
	api_url_ := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", bot_token)
	params := url.Values{}
	params.Add("chat_id", chat_id)
	params.Add("text", message)
	params.Add("parse_mode", "MarkdownV2")

	resp, err := http.Get(fmt.Sprintf("%s?%s", api_url_, params.Encode()))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func escapeMarkdownV2(text string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)
	return replacer.Replace(text)
}


func solve_count(submissions []Submission) (int, []string) {
	unique_solve := make(map[string]bool)
	today := time.Now().Format("2006-01-02")
	solved_problems := []string{}

	for _, sub := range submissions {
		if sub.Verdict == "OK" {
	loc, _ := time.LoadLocation("Asia/Dhaka")
			submission_date := time.Unix(sub.CreationTime, 0).In(loc).Format("2006-01-02")			
			if submission_date == today {
				problemID := fmt.Sprintf("%d-%s-%s-%s", sub.ContestID, sub.Problem.Index, sub.Problem.Name, submission_date)
				if !unique_solve[problemID] {
					problemURL := fmt.Sprintf("https://codeforces.com/contest/%d/problem/%s", sub.ContestID, sub.Problem.Index)
					escapedName := escapeMarkdownV2(sub.Problem.Name)
					escapedURL := escapeMarkdownV2(problemURL)
					problemEntry := fmt.Sprintf("[%s](%s)", escapedName, escapedURL)
					solved_problems = append(solved_problems, problemEntry)
				}
				unique_solve[problemID] = true
			}
		}

	}

	return len(unique_solve), solved_problems
}


func send_message(message string) error {
	api_url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", bot_token)
	params := url.Values{}
	params.Add("chat_id", chat_id)
	params.Add("text", message)
	params.Add("parse_mode", "MarkdownV2")
	//fmt.Println("Sending request to:", fmt.Sprintf("%s?%s", api_url, params.Encode()))
	
	resp, err := http.Get(fmt.Sprintf("%s?%s", api_url, params.Encode()))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}


func main() {
	
	loc, err := time.LoadLocation("Asia/Dhaka")
	if err != nil {
		fmt.Println("Error loading timezone:", err)
		return
	}

	scheduler := gocron.NewScheduler(loc)

	scheduler.Every(1).Day().At("08:24").Do(func() {
		report := generate_daily_report()
		if err := sendMessage(report); err != nil {
			fmt.Println("Error sending daily report:", err)
		} else {
			fmt.Println("Daily report sent successfully!")
		}
	})

	scheduler.StartAsync()
	go listen_command()

	select {}
}
