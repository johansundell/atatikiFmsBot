package main

import (
	"context"
	"fmt"
	//	"github.com/johansundell/atatikiFmsBot/fmsadmin"
	//	"net/http"
)

func init() {
	key := commandFunc{"!files", "List server files", "", categoryAdmin}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(ctx *context.Context, command string) (string, error) {
		if command == key.command {

			files, err := server.GetFiles()

			if err != nil {
				return "", err
			}
			msg := "Open files\n"
			closed := "\nClosed files\n"
			for _, r := range files.Files.Files {
				//msg += r.Filename + "\n"
				if r.Status == "NORMAL" {
					msg += fmt.Sprintf("Id: %s, Name: %s, Clients: %d\n", r.ID, r.Filename, r.Clients)
				} else {
					closed += "Id: " + r.ID + ", " + r.Filename + "\n"
				}
			}
			msg += closed + "\n"

			msg += fmt.Sprintf("FileMaker Pro: %d\n", files.FmproCount)
			msg += fmt.Sprintf("FileMaker Go: %d\n", files.FmgoCount)
			msg += fmt.Sprintf("FileMaker WebDirect: %d\n", files.FmwebdCount)
			msg += fmt.Sprintf("FileMaker Data API: %d\n", files.FmdapiCount)
			msg += fmt.Sprintf("Connected but not counted: %d\n", files.FmmiscCount)
			return msg, nil
		}
		return "", nil
	}
}
