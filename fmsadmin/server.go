package fmsadmin

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	url       string
	user      string
	pass      string
	token     string
	tokenTime time.Time // TODO: Are we using this
	sync.RWMutex
	quit chan struct{}
}

func NewServer(url, user, pass string) Server {
	return Server{url: url, user: user, pass: pass}
}

func getClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return client
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

	req, _ := http.NewRequest("POST", url, b)
	req.Header.Set("Content-Type", "application/json")
	resp, err := getClient().Do(req)

	defer resp.Body.Close()
	loginResult := struct {
		Result int    `json:"result"`
		Token  string `json:"token"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&loginResult); err != nil {
		return err
	}
	if loginResult.Result == 0 {
		s.tokenTime = time.Now()
		s.token = loginResult.Token

		// We have a token, keep it and renew when needed
		ticker := time.NewTicker(14 * time.Minute)
		s.quit = make(chan struct{})
		go func() {
			for {
				select {
				case <-ticker.C:
					s.Lock()
					s.Logout()
					s.Login()
					s.Unlock()

				case <-s.quit:
					ticker.Stop()
					return
				}
			}
		}()
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
	resp, err := getClient().Do(req)
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
	close(s.quit)
	return nil
}
