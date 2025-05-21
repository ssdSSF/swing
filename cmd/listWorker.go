/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/ssdSSF/swing/pkg/model"
)

// listWorkerCmd represents the listWorker command
var listWorkerCmd = &cobra.Command{
	Use:   "list-worker",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := slog.Default()

		accessToken, err := refreshToken()
		if err != nil {
			logger.Error("err refreshToken", "err", err)
		}
		notified := map[string]bool{}
		for {

			hour, _, _ := time.Now().Local().Clock()
			// between 00:00 to 4:59, this worker also goes to sleep
			if hour >= 0 && hour < 5 {
				time.Sleep(1 * time.Minute)
			}

			if needRefresh(accessToken) {
				accessToken, err = refreshToken()
				if err != nil {
					logger.Error("err refreshToken", "err", err)
				}
			}

			op, err := openings(accessToken)
			if err != nil {
				logger.Error("err openings", "err", err)
			}

			var dataIds []string
			for _, data := range op.Data {
				// skip the cities that you don't want to go
				if slices.Contains(cmdSecrets.CitiesToSkip, strings.ToLower(data.School.Address.City)) {
					continue
				}

				dataIds = append(dataIds, data.ID)

				_, sent := notified[data.ID]

				if sent {
					continue
				} else {
					notified[data.ID] = true
				}

				err := slack(toMessage(data))
				if err != nil {
					logger.Error("err slack(data.ID)", "err", err)
				}
			}

			var toRemove []string
			for id := range notified {
				if !slices.Contains(dataIds, id) {
					toRemove = append(toRemove, id)
				}
			}

			for _, id := range toRemove {
				delete(notified, id)
			}

			time.Sleep(time.Second * cmdSecrets.Interval)
		}
	},
}

func refreshToken() (string, error) {
	tokenUrl := "https://securetoken.googleapis.com/v1/token?key=AIzaSyDKBOlXVXPr4RY_Px3tsWbsIa3ozigoH18"
	resp, err := http.PostForm(tokenUrl, url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {cmdSecrets.GoogleToken},
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("resp.StatusCode != http.StatusOK")
	}

	refreshToken := model.RefreshToken{}
	json.NewDecoder(resp.Body).Decode(&refreshToken)

	return refreshToken.AccessToken, nil
}

func needRefresh(authToken string) bool {

	parts := strings.Split(authToken, ".")
	if len(parts) != 3 { // "aa.bb.cc"
		return true
	}

	jwt := model.JWT{}

	err := json.NewDecoder(strings.NewReader(parts[1])).Decode(&jwt)
	if err != nil {
		return true
	}

	now := time.Now().Unix()

	// if expiring in less than 300 seconds, refresh
	return jwt.Exp-now < 300
}

func openings(accessToken string) (openings model.Openings, err error) {
	client := http.Client{}
	openingUrl := "https://prod-api.aws.swingeducation.com/api/sub/request/openings"
	openingReq, err := http.NewRequest(http.MethodGet, openingUrl, nil)
	if err != nil {
		return
	}
	openingReq.Header.Add("authorization", fmt.Sprintf("Token %s", accessToken))

	openingResp, err := client.Do(openingReq)
	if err != nil {
		return
	}

	err = json.NewDecoder(openingResp.Body).Decode(&openings)

	return
}

func slack(message string) (err error) {
	client := http.Client{}
	slackUrl := "https://slack.com/api/chat.postMessage"

	slackReq := model.SlackReq{
		Channel: cmdSecrets.SlackChannel,
		Text:    message,
	}
	body, err := json.Marshal(slackReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, slackUrl, bytes.NewBuffer(body))
	if err != nil {
		return
	}

	req.Header.Add("Content-type", "application/json; charset=utf-8")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cmdSecrets.SlackToken))

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	slackResp := model.SlackResp{}
	err = json.NewDecoder(resp.Body).Decode(&slackResp)

	if !slackResp.Ok {
		err = fmt.Errorf("slackResp is not ok")
	}
	return
}

func init() {
	rootCmd.AddCommand(listWorkerCmd)
}
