package main

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/johansundell/atatikiFmsBot/fmsadmin"
)

func init() {
	key := commandFunc{"!open ([\\d]+)", "Open a file, ex !open 34", "", categoryAdmin}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(ctx *context.Context, command string) (string, error) {
		if found, _ := regexp.MatchString(key.command, command); found {
			re, _ := regexp.Compile(key.command)
			strs := re.FindStringSubmatch(command)
			if len(strs) == 2 {
				id, err := strconv.Atoi(strs[1])
				if err != nil {
					return "", err
				}

				s := fmsadmin.NewServer(settings.Url, settings.User, settings.Pass)
				if err := s.Login(); err != nil {
					return "", err
				}
				err = s.OpenFile(id)
				s.Logout()
				if err != nil {
					return err.Error(), nil
				}
				return fmt.Sprintf("File with id %d is open", id), nil
			}
		}
		return "", nil
	}
}
