package model

import "time"

type Secrets struct {
	GoogleToken           string        `yaml:"googleToken"`
	SlackToken            string        `yaml:"slackToken"`
	SlackChannel          string        `yaml:"slackChannel"`
	SlackHeartbeatChannel string        `yaml:"slackHeartbeatChannel"`
	Interval              time.Duration `yaml:"interval"`
	CitiesToSkip          []string      `yaml:"citiesToSkip"`
}

type RefreshToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
}

type JWT struct {
	Exp int64 `json:"exp"`
}

type SlackReq struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

type SlackResp struct {
	Ok bool `json:"ok"`
}

type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
	State  string `json:"state"`
	Zip    string `json:"zip"`
}

type Opening struct {
	Grades      []string `json:"grades"`
	DisplayDate string   `json:"displayDate"`
	IsMultiDay  bool     `json:"isMultiDay"`
	School      struct {
		ID       string  `json:"id"`
		Name     string  `json:"name"`
		Address  Address `json:"address"`
		TimeZone string  `json:"timeZone"`
	} `json:"school"`
	Payment struct {
		IsHourly        bool      `json:"isHourly"`
		PayValue        float64   `json:"payValue"`
		PaymentDate     time.Time `json:"paymentDate"`
		FullDayPayValue float64   `json:"fullDayPayValue"`
		HalfDayPayValue float64   `json:"halfDayPayValue"`
		TotalPayValue   float64   `json:"totalPayValue"`
		HourlyPayValue  float64   `json:"hourlyPayValue"`
	} `json:"payment"`
	DidSubSubmitAvailability bool      `json:"didSubSubmitAvailability"`
	StartDate                time.Time `json:"startDate"`
	CreatedAt                time.Time `json:"createdAt"`
	DisplayTime              string    `json:"displayTime"`
	EndDate                  time.Time `json:"endDate"`
	IsMultiWeek              bool      `json:"isMultiWeek"`
	Status                   string    `json:"status"`
	ID                       string    `json:"id"`
	Subjects                 []string  `json:"subjects"`
	HasFeedback              bool      `json:"hasFeedback"`
	Intervals                []struct {
		EndDate     time.Time `json:"endDate"`
		EndOffset   int       `json:"endOffset"`
		StartDate   time.Time `json:"startDate"`
		StartOffset int       `json:"startOffset"`
	} `json:"intervals"`
	Memo string `json:"memo"`
}

type Openings struct {
	Data []Opening `json:"data"`
}
