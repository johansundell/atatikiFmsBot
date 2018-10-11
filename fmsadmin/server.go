package fmsadmin

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	url       string
	user      string
	pass      string
	token     string
	tokenTime time.Time
}

func NewServer(url, user, pass string) Server {
	return Server{url: url, user: user, pass: pass}
}

func (s *Server) Login() error {
	url := s.url + "/fmi/admin/api/v1/user/login"
	loginInfo := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{s.user, s.pass}
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(loginInfo)
	if err != nil {
		return err
	}
	// TODO: Remove this and use real cert instead
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Post(url, "application/json", b)
	if err != nil {
		return err
	}
	/*req := http.NewRequest("POST", url, b)
		req.Header.Set("Content-Type", "application/json", b)
		tr := &http.Transport{
	        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	    }
		client := &http.Client{Transport: tr}*/

	defer resp.Body.Close()
	loginResult := struct {
		Result int    `json:"result"`
		Token  string `json:"token"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&loginResult); err != nil {
		return err
	}
	//fmt.Printf("%+v\n", loginResult)
	//fmt.Println(loginResult)
	if loginResult.Result == 0 {
		s.tokenTime = time.Now()
		s.token = loginResult.Token
	}

	return nil
}

func (s *Server) Logout() error {
	url := s.url + "/fmi/admin/api/v1/user/logout"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+s.token)
	req.Header.Add("Content-Type", "application/json")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	logoutInfo := struct {
		Result int `json:"result"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&logoutInfo); err != nil {
		return err
	}
	if logoutInfo.Result != 0 {
		return errors.New("Failed to logout")
	}
	fmt.Println("We could log out")
	return nil
}

func (s *Server) GetFiles() (f Files, err error) {
	url := s.url + "/fmi/admin/api/v1/databases"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return f, err
	}
	req.Header.Add("Authorization", "Bearer "+s.token)
	req.Header.Add("Content-Type", "application/json")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
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
