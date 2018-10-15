package fmsadmin

import (
	"fmt"
)

func (s *Server) OpenFile(id int) error {
	url := s.url + fmt.Sprintf("/fmi/admin/api/v1/databases/%d/open", id)
	b, err := s.makeCall(url, "PUT", nil)
	_ = b
	return err
}
