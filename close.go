package main

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
)

func init() {
	key := commandFunc{"!close id ([\\d]+) message:(.*)", "Close file, ex !close id 4 message: I am closing this file", "", categoryAdmin}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(ctx *context.Context, command string) (string, error) {
		if found, _ := regexp.MatchString(key.command, command); found {
			re, _ := regexp.Compile(key.command)
			strs := re.FindStringSubmatch(command)
			if len(strs) == 3 {
				id, err := strconv.Atoi(strs[1])
				if err != nil {
					return "", err
				}
				message := strs[2]

				err = server.Close(id, message)
				if err != nil {
					return err.Error(), nil
				}
				return fmt.Sprintf("File with id %d is closed", id), nil
			}
		}
		return "", nil
	}
}
