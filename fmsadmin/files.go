package fmsadmin

import (
	"encoding/json"
	"net/http"
)

func (s *Server) GetFiles() (f Files, err error) {
	url := s.url + "/fmi/admin/api/v1/databases"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return f, err
	}
	req.Header.Add("Authorization", "Bearer "+s.token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := getClient().Do(req)
	if err != nil {
		return f, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&f); err != nil {
		return f, err
	}
	return f, nil
}

type Files struct {
	Clients struct {
		Clients []struct {
			AppLanguage     string `json:"appLanguage"`
			AppType         string `json:"appType"`
			AppVersion      string `json:"appVersion"`
			ComputerName    string `json:"computerName"`
			Concurrent      bool   `json:"concurrent"`
			ConnectDuration string `json:"connectDuration"`
			ConnectTime     string `json:"connectTime"`
			Extpriv         string `json:"extpriv"`
			GuestFiles      []struct {
				AccountName string `json:"accountName"`
				Filename    string `json:"filename"`
				ID          string `json:"id"`
				PrivsetName string `json:"privsetName"`
			} `json:"guestFiles"`
			ID              string `json:"id"`
			Ipaddress       string `json:"ipaddress"`
			Macaddress      string `json:"macaddress"`
			OperatingSystem string `json:"operatingSystem"`
			Status          string `json:"status"`
			TeamLicensed    bool   `json:"teamLicensed"`
			UserName        string `json:"userName"`
		} `json:"clients"`
		Result int `json:"result"`
	} `json:"clients"`
	Files struct {
		Files []struct {
			Clients              int      `json:"clients"`
			DecryptHint          string   `json:"decryptHint"`
			EnabledExtPrivileges []string `json:"enabledExtPrivileges"`
			Filename             string   `json:"filename"`
			Folder               string   `json:"folder"`
			HasSavedDecryptKey   bool     `json:"hasSavedDecryptKey"`
			ID                   string   `json:"id"`
			IsEncrypted          bool     `json:"isEncrypted"`
			Size                 int      `json:"size"`
			Status               string   `json:"status"`
		} `json:"files"`
		Result int `json:"result"`
	} `json:"files"`
	FmdapiCount  int    `json:"fmdapiCount"`
	FmgoCount    int    `json:"fmgoCount"`
	FmmiscCount  int    `json:"fmmiscCount"`
	FmproCount   int    `json:"fmproCount"`
	FmwebdCount  int    `json:"fmwebdCount"`
	OpenDBCount  int    `json:"openDBCount"`
	Result       int    `json:"result"`
	Time         string `json:"time"`
	TotalDBCount int    `json:"totalDBCount"`
}
