package main

import (
	"context"
	"fmt"
	"sort"
)

func init() {
	key := commandFunc{"!help", "List the help", "", catgoryHelp}
	lockMap.Lock()
	defer lockMap.Unlock()
	botFuncs[key] = func(ctx *context.Context, command string) (string, error) {
		if command == key.command {
			var keys = make([]commandFunc, 0, len(botFuncs))
			for k := range botFuncs {
				keys = append(keys, k)
			}
			sort.Slice(keys, func(i, j int) bool {
				if keys[i].category == keys[j].category {
					return keys[i].command < keys[j].command
				}
				return keys[i].category < keys[j].category
			})
			msg := "**BOT COMMANDS**\n"
			var c category
			for _, v := range keys {
				if v.category != categoryHidden {
					if c != v.category {
						msg += fmt.Sprintf("\n%s\n", v.category)
						c = v.category
					}
					msg += fmt.Sprintf("%s - %s\n", v.command, v.helpText)
				}
			}
			return msg, nil
		}
		return "", nil
	}
}
