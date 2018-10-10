package main

type Config struct {
	Token string `json:"token"`
	// TODO: Store these in a database and load them from there
	Url  string `json:"url"`
	User string `json:"user"`
	Pass string `json:"pass"`
}
