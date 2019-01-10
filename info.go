package main

import "github.com/microredis/tools/encoding/info"

type Info struct {
	MasterLinkStatus string `info:"master_link_status"`
}

func (i *Info) UnmarshalBinary(data []byte) error {
	return info.Unmarshal(data, i)
}
