package main

import "github.com/microredis/tools/encoding/client/list"

type ClientList []Client

func (c *ClientList) UnmarshalBinary(data []byte) error {
	return list.Unmarshal(data, c)
}


