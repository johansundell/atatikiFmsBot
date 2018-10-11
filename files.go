package main

import (
	"context"

	"github.com/johansundell/atatikiFmsBot/fmsadmin"
	//	"net/http"
)

func init() {
	key := commandFunc{"!files", "List server files", "", categoryAdmin}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(ctx *context.Context, command string) (string, error) {
		if command == key.command {
			s := fmsadmin.NewServer(settings.Url, settings.User, settings.Pass)
			if err := s.Login(); err != nil {
				return "", err
			}
			defer s.Logout()
			files, err := s.GetFiles()
			if err != nil {
				return "", err
			}
			msg := ""
			for _, r := range files.Files.Files {
				msg += r.Filename + "\n"
			}
			return msg, nil
		}
		return "", nil
	}
}
