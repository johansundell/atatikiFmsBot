package fmsadmin

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
)

type Status struct {
	Result            int `json:"result"`
	CacheSize         int `json:"cacheSize"`
	MaxFiles          int `json:"maxFiles"`
	MaxProConnections int `json:"maxProConnections"`
	MaxPSOS           int `json:"maxPSOS"`
}

func (s *Server) GetStatus() (*Status, error) {
	url := s.url + "/fmi/admin/api/v1/server/config/general"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+s.token)
	req.Header.Add("Content-Type", "application/json")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	status := &Status{}
	if err := json.NewDecoder(resp.Body).Decode(status); err != nil {
		return nil, err
	}
	return status, nil
}
