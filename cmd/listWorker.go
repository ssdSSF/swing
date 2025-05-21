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
	"os"
	"slices"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/ssdSSF/swing/pkg/model"
)

var errCount int64

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

		go hearBeat(logger)

		// need to get a valid google token
		accessToken, err := refreshToken()
		if err != nil {
			logger.Error("err refreshToken", "err", err)
			os.Exit(1)
		}

		// and a success openings
		_, err = openings(accessToken)
		if err != nil {
			logger.Error("err openings with accessToken", "err", err)
			os.Exit(1)
		}

		// before we can say the worker is started
		logger.Info("started worker")

		notified := map[string]bool{}
		for {

			hour, _, _ := time.Now().Local().Clock()
			// between 00:00 to 4:59, this worker also goes to sleep
			if hour >= 0 && hour < 5 {
				// slow down at midnight
				time.Sleep(10 * time.Minute)
			}

			if needRefresh(accessToken) {
				accessToken, err = refreshToken()
				if err != nil {
					errCount++
					logger.Error("err refreshToken", "err", err)
				}
			}

			op, err := openings(accessToken)
			if err != nil {
				errCount++
				logger.Error("err openings", "err", err)
			}

			var dataIds []string
			for _, data := range op.Data {

				dataIds = append(dataIds, data.ID)

				_, sent := notified[data.ID]

				// we only want to notify once
				if sent {
					continue
				} else {
					notified[data.ID] = true
				}

				err := slack(toMessage(data), cmdSecrets.SlackChannel)
				if err != nil {
					errCount++
					logger.Error("err slack(data.ID)", "err", err)
				}
			}

			// to handle an opening is re-opened
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

func hearBeat(logger *slog.Logger) {

	// heartbeat first
	err := slack(fmt.Sprintf("started hearbeat, error count %d", errCount), cmdSecrets.SlackHeartbeatChannel)
	if err != nil {
		logger.Error("err start slack heatbeat", "err", err)
		os.Exit(1)
	}
	logger.Info("started heartbeat")

	for {
		// then sleep
		time.Sleep(10 * time.Minute)

		err := slack(fmt.Sprintf("hearbeat, error count %d", errCount), cmdSecrets.SlackHeartbeatChannel)
		if err != nil {
			errCount++
			logger.Error("err slack hearbeat", "err", err)
		}
	}
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

	var all model.Openings

	err = json.NewDecoder(openingResp.Body).Decode(&all)

	for _, op := range all.Data {
		// if school's city belongs to cities to skip
		if slices.Contains(cmdSecrets.CitiesToSkip, strings.ToLower(op.School.Address.City)) {
			// then skip
			continue
		}
		openings.Data = append(openings.Data, op)
	}

	return
}

func slack(message, channel string) (err error) {
	client := http.Client{}
	slackUrl := "https://slack.com/api/chat.postMessage"

	slackReq := model.SlackReq{
		Channel: channel,
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
