package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Faagerholm/clockify-cli/domain"
	"github.com/spf13/viper"
)

// GetAllEntries returns all entries from the Clockify API.
// From the first day, until today
func GetAllEntries() (
	*domain.ResultUser,
	error,
) {
	var first_day_str, last_day_str string

	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)

	last_day_str = now.Format("2006-01-02T15:04:05Z")
	first_day_str = fmt.Sprintf("%d-01-01T00:00:00Z", now.Year())

	users, err := getEntries(first_day_str, last_day_str)
	if err != nil {
		return nil, err
	}
	if len(users) != 1 {
		return nil, errors.New(fmt.Sprintf("Invalid entries count found, expected 1, got %d", len(users)))
	}
	user := users[0]
	year := now.Year() - 1

	for true {
		first_day_str = fmt.Sprintf("%d-01-01T00:00:00Z", year)
		last_day_str = fmt.Sprintf("%d-12-31T23:59:59Z", year)

		users, err := getEntries(first_day_str, last_day_str)
		if err != nil {
			return nil, err
		}
		if len(users) != 1 {
			break
		}
		user.Duration += users[0].Duration
		user.Entries = append(user.Entries, users[0].Entries...)
		year--
	}

	if err != nil {
		log.Fatal(err)
	}
	return &user, err
}

func GetUserEntries(start_date, end_date time.Time) (*domain.ResultUser, error) {
	res, err := getEntries(start_date.Format("2006-01-02T15:04:05Z"), end_date.Format("2006-01-02T15:04:05Z"))
	if err != nil {
		return nil, err
	}
	if len(res) != 1 {
		return nil, errors.New(fmt.Sprintf("Invalid user count found, expected 1, got %d", len(res)))
	}

	return &res[0], nil
}

func getEntries(start_date, end_date string) ([]domain.ResultUser, error) {
	user := viper.GetString("USER-ID")

	if len(user) == 0 {
		// TODO: Notify user to run setup (run it for them)
		return nil, errors.New("HANDLE ME")
	}

	reqBody, err := marshalRequestBody(user, start_date, end_date)
	if err != nil {
		log.Fatal(err)
	}

	res, err := requestReport(reqBody)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return res.Entries, nil
}

func requestReport(body []byte) (*domain.Result, error) {
	key := viper.GetString("API-KEY")
	workspace := viper.GetString("WORKSPACE")

	if len(key) == 0 {
		return nil, errors.New("API KEY NOT SET! Please run clockify app-key")
	}

	client := &http.Client{}
	req, _ := http.NewRequest("POST", fmt.Sprintf("https://reports.api.clockify.me/v1/workspaces/%s/reports/summary", workspace), bytes.NewBuffer(body))
	req.Header.Set("X-API-KEY", key)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	res := new(domain.Result)

	json.NewDecoder(resp.Body).Decode(&res)
	return res, err
}

func marshalRequestBody(user, start_date, end_date string) ([]byte, error) {

	return json.Marshal(domain.Report{
		Start: start_date,
		End:   end_date,
		SummaryFilter: &domain.Report_filter{
			Groups: []string{"user", "date"},
		},
		SortOrder: "Ascending",
		Users: &domain.Report_user{
			Ids:      []string{user},
			Contains: "CONTAINS",
			Status:   "All",
		},
	})
}
