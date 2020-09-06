package main

import (
	"context"
	"encoding/json"
	"github.com/google/go-github/v32/github"
	"github.com/ldez/ghactions"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type JestResult struct {
	NumFailedTests int
	Success        bool
	TestResults    []*TestResult
}

type TestResult struct {
	Message string
	Name    string
	Status  string
}

type Slack struct {
	Text        string       `json:"text"`
	Username    string       `json:"username"`
	IconEmoji   string       `json:"icon_emoji"`
	Mrkdwn      bool         `json:"mrkdwn"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	AuthorName string   `json:"author_name"`
	AuthorIcon string   `json:"author_icon"`
	Fallback   string   `json:"fallback"`
	Color      string   `json:"color"`
	MrkdwnIn   []string `json:"mrkdwn_in"`
	Fields     []Field  `json:"fields"`
}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

func main() {
	webhookURL := os.Getenv("INCOMING_WEBHOOK_URL")
	filePath := os.Getenv("JEST_FILE_PATH")

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	var jestResultData JestResult
	err = json.Unmarshal(b, &jestResultData)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	action := ghactions.NewAction(ctx)

	err = action.OnPush(func(client *github.Client, event *github.PushEvent) error {
		return handler(ctx, client, event, jestResultData, string(webhookURL))
	}).Run()
	if err != nil {
		log.Fatal(err)
	}
}

func handler(ctx context.Context, client *github.Client, event *github.PushEvent, result JestResult, webhookURL string) error {
	if result.Success {
		return nil
	}

	var attachments []Attachment

	attachments = append(attachments, Attachment{
		*event.Sender.Login,
		*event.Sender.AvatarURL,
		"failed",
		"danger",
		[]string{"fields"},
		[]Field{
			Field{
				Title: "Compare",
				Value: *event.Compare,
				Short: false,
			},
		},
	})

	for _, t := range result.TestResults {
		if t.Status != "failed" {
			continue
		}

		attachments = append(attachments, Attachment{
			*event.Sender.Login,
			*event.Sender.AvatarURL,
			"failed",
			"danger",
			[]string{"fields"},
			[]Field{
				Field{
					Title: t.Name + "",
					Value: "```" + string(t.Message) + "```",
					Short: false,
				},
			},
		})
	}

	slack := Slack{
		"Failed Test",
		"Negima",
		":sob:",
		true,
		attachments,
	}

	p, err := json.Marshal(slack)
	if err != nil {
		return err
	}
	resp, err := http.PostForm(webhookURL, url.Values{"payload": {string(p)}})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
