package fmsadmin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (s *Server) Close(id int, message string) error {
	url := s.url + fmt.Sprintf("/fmi/admin/api/v1/databases/%d/close", id)
	body := struct {
		Message string `json:"message"`
	}{message}
	b, err := s.makeCall(url, "PUT", &body)
	_ = b
	return err
}

func (s *Server) makeCall(url, method string, inf interface{}) ([]byte, error) {

	var err error = nil
	req := &http.Request{}
	if inf != nil {
		b := new(bytes.Buffer)
		err := json.NewEncoder(b).Encode(inf)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, url, b)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+s.token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := getClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
