package main

type Client struct {
	Flags string `client_list:"flags"`
	Omem  int64  `client_list:"omem"`
	Obl   int64  `client_list:"obl"`
	Oll   int64  `client_list:"oll"`
}
