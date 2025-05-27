/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/ssdSSF/swing/pkg/model"
)

// previewCmd represents the preview command
var previewCmd = &cobra.Command{
	Use:   "preview",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		example := `
    {
      "grades": [
        "Unknown"
      ],
      "displayDate": "Wed, May 21",
      "isMultiDay": false,
      "school": {
        "id": "56cfb707-a02e-444b-92a6-00a3a47f2543",
        "name": "Lighthouse K-5",
        "address": {
          "street": "444 Hegenberger Road",
          "city": "Oakland",
          "state": "CA",
          "zip": "94621"
        },
        "timeZone": "US/Pacific"
      },
      "payment": {
        "isHourly": true,
        "payValue": 3436.9999999999995,
        "paymentDate": "2025-05-30T16:30:00.000Z",
        "fullDayPayValue": 30624.0,
        "halfDayPayValue": 15312.0,
        "totalPayValue": 13749.999912000001,
        "hourlyPayValue": 3436.9999999999995
      },
      "didSubSubmitAvailability": false,
      "startDate": "2025-05-21T16:30:00.000Z",
      "createdAt": "2025-05-20T20:03:54.917Z",
      "displayTime": "9:30AM - 1:30PM",
      "endDate": "2025-05-21T20:30:00.000Z",
      "isMultiWeek": false,
      "status": "STATUS_OPEN",
      "id": "682ce02a-6e59-4fb8-8c02-a658d747f21e",
      "subjects": [
        "Unknown"
      ],
      "hasFeedback": false,
      "intervals": [
        {
          "endDate": "2025-05-21T20:30:00.000Z",
          "endOffset": 420,
          "startDate": "2025-05-21T16:30:00.000Z",
          "startOffset": 420
        }
      ],
      "memo": "You will receive more details about your teaching assignment when you check in at the school. You may be in one classroom or moving to different classrooms throughout the day."
    }
`
		opening := model.Opening{}
		json.NewDecoder(strings.NewReader(example)).Decode(&opening)

		fmt.Println(toMessage(opening))
	},
}

func toMessage(opening model.Opening) string {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf(`%s
%s
Grades: %s
Start: %s
End: %s
CreatedAt: %s
TotalPay: $%.2f
HourlyRate: $%.2f`, opening.School.Name, toAddress(opening.School.Address), opening.Grades,
		opening.StartDate.In(loc).Format(time.RFC1123), opening.EndDate.In(loc).Format(time.RFC1123), opening.CreatedAt.In(loc).Format(time.RFC1123),
		opening.Payment.TotalPayValue/100, opening.Payment.HourlyPayValue/100)
}

func toAddress(address model.Address) string {

	return fmt.Sprintf("%s, %s, %s %s", address.Street, address.City, address.State, address.Zip)
}

func init() {
	rootCmd.AddCommand(previewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// previewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// previewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
