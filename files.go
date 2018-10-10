package main

import (
	"context"
	//	"net/http"
)

func init() {
	key := commandFunc{"!files", "List server files", "", categoryAdmin}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(ctx *context.Context, command string) (string, error) {
		if command == key.command {
			// Get the FMS data here
			return "found", nil
		}
		return "", nil
	}
}

type FilesResp struct {
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
