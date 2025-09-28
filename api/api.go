package api

import (
	"fmt"
	"io"

	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httputil"

	"os"
)

type loggingTransport struct{}

type CliConfig struct {
	// shared parameters
	LogLevel  string `mapstructure:"log-level"`
	Cookies   string `mapstructure:"cookies"`
	FirstName string `mapstructure:"first-name"`
}

type HouseHoldMember struct {
	AvatarUri  string `json:"avatarUri"`
	Role       string `json:"role"`
	FirstName  string `json:"firstName"`
	DirectedId string `json:"directedId"`
}

type HouseHold struct {
	HouseholdId string            `json:"householdId"`
	Members     []HouseHoldMember `json:"members"`
}

type CurfewConfig struct {
	End     string `json:"end"`
	Type    any    `json:"type"`
	Start   string `json:"start"`
	Enabled bool   `json:"enabled"`
}

type TimeLimits struct {
	ContentTimeLimitsEnabled bool           `json:"contentTimeLimitsEnabled"`
	ContentTimeLimits        map[string]int `json:"contentTimeLimits"`
}

type ContentGoals struct {
	Category_BOOK    int `json:"category_BOOK"`
	Category_VIDEO   int `json:"category_VIDEO"`
	Category_APP     int `json:"category_APP"`
	Category_AUDIBLE int `json:"category_AUDIBLE"`
	Category_WEB     int `json:"category_WEB"`
}

type GoalsConfig struct {
	ContentGoals      ContentGoals `json:"contentGoals"`
	LearnFirstEnabled bool         `json:"learnFirstEnabled"`
}

type PeriodConfig struct {
	Type             string         `json:"type"`
	Name             string         `json:"name"`
	Enabled          bool           `json:"enabled"`
	CurfewConfigList []CurfewConfig `json:"curfewConfigList"`
	Time             uint64         `json:"time"`
	TimeLimits       TimeLimits     `json:"timeLimits"`
	GoalsConfig      GoalsConfig    `json:"goalsConfig"`
}

type PeriodConfigs struct {
	PeriodConfigurations []PeriodConfig `json:"periodConfigurations"`
}

type HTTPParameter struct {
	Key   string
	Value string
}

func (s *loggingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	bytes, _ := httputil.DumpRequestOut(r, true)

	resp, err := http.DefaultTransport.RoundTrip(r)
	// err is returned after dumping the response

	respBytes, _ := httputil.DumpResponse(resp, true)
	bytes = append(bytes, respBytes...)

	fmt.Printf("%s\n", bytes)

	return resp, err
}

func ConfigureLogging(level string) error {
	var logLevel slog.Level

	switch level {
	case "error":
		logLevel = slog.LevelError
	case "warn":
		logLevel = slog.LevelWarn
	case "info":
		logLevel = slog.LevelInfo
	case "debug":
		logLevel = slog.LevelDebug
	default:
		return fmt.Errorf("unknown log-level: %s", level)
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)

	slog.SetDefault(logger)

	return nil
}

func (c CliConfig) httpRequest(url string, parms []HTTPParameter, data any) error {
	cookies, err := http.ParseCookie(c.Cookies)
	if err != nil {
		return err
	}

	res, err := http.Get(url)
	if err != nil {
		return err
	}

	// Create a new GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("accept-language", "de-DE,de;q=0.9")
	req.Header.Add("priority", "u=1, i")
	req.Header.Add("referer", "https://eltern.amazon.de/settings/timelimits")
	req.Header.Add("sec-ch-ua", "\"Chromium\";v=\"140\", \"Not=A?Brand\";v=\"24\", \"Google Chrome\";v=\"140\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36")

	for _, cookie := range cookies {
		slog.Debug("getHousehold()", "cookie", fmt.Sprintf("%#v", cookie))
		req.AddCookie(cookie)
	}

	q := req.URL.Query()
	for _, parm := range parms {
		q.Add(parm.Key, parm.Value)
	}
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}

	if c.LogLevel == "debug" {
		client.Transport = &loggingTransport{}
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if resp.StatusCode > 299 {
		return fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		return err
	}
	slog.Debug("getHousehold()", "body", body)

	if err := json.Unmarshal(body, data); err != nil {
		return fmt.Errorf("json.Unmarshal() failed: %s\n\nTried to parse this content:\n%s", err, string(body))
	}

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))

	return nil
}

func (c CliConfig) getHousehold() (HouseHold, error) {
	houseHold := HouseHold{}

	err := c.httpRequest("https://eltern.amazon.de/ajax/get-household", []HTTPParameter{}, &houseHold)
	if err != nil {
		return houseHold, err
	}

	for _, member := range houseHold.Members {
		slog.Debug("getHousehold()", "member", fmt.Sprintf("%#v", member))
	}

	return houseHold, nil
}

func (c CliConfig) getDirectedId() (string, error) {
	foundMembers := []string{}
	houseHold, err := c.getHousehold()
	if err != nil {
		return "", err
	}

	for _, member := range houseHold.Members {
		if member.FirstName == c.FirstName {
			return member.DirectedId, nil
		}

		foundMembers = append(foundMembers, member.FirstName)
	}

	return "", fmt.Errorf("Did not find first-name: %q in members: %#v", c.FirstName, foundMembers)
}

func (c CliConfig) GetTimes() error {
	if err := ConfigureLogging(c.LogLevel); err != nil {
		return err
	}

	slog.Info("GetTimes() called")

	directedId, err := c.getDirectedId()
	if err != nil {
		return err
	}

	periodConfigs := PeriodConfigs{}

	err = c.httpRequest("https://eltern.amazon.de/ajax/get-time-limit-v2", []HTTPParameter{{"childDirectedId", directedId}}, &periodConfigs)
	if err != nil {
		return err
	}

	for _, periodConfig := range periodConfigs.PeriodConfigurations {
		slog.Info("GetTimes()", "periodConfig", fmt.Sprintf("%#v", periodConfig))
	}

	return nil
}

func (c CliConfig) SetTimes() error {
	if err := ConfigureLogging(c.LogLevel); err != nil {
		return err
	}

	slog.Info("SetTimes() called")

	return nil
}
