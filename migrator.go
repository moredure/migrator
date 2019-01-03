package main

import (
	"github.com/go-redis/redis"
	"net"
	"strings"
)

type Migrator struct {
	from *redis.Client
	to *redis.Client
	maxOutBuff int64
	maxOutBuffCommands int64
}

func NewMigrator(fromRedisClient FromRedisClient, toRedisClient ToRedisClient) *Migrator {
	return &Migrator{fromRedisClient, toRedisClient, -1, -1}
}

func (m *Migrator) Migrate() {
	defer m.to.Close()
	defer m.from.Close()
	m.PrepareTo()
	m.WaitForUp()
	m.WaitForComplete()
	m.OnComplete()
}

func (m *Migrator) PrepareTo() {
	Must(m.to.ConfigSet("slave-read-only", "yes").Result())
	Must(m.to.ConfigSet("masterauth", m.from.Options().Password).Result())
	host, port, _ := net.SplitHostPort(m.from.Options().Addr)
	Must(m.to.SlaveOf(host, port).Result())
}

func (m *Migrator) WaitForUp() {
	for {
		info := Must(m.to.Info().Result())
		infoParsed := ParseMetrics(info)
		if infoParsed["master_link_status"].(string) == "up" {
			break
		}
	}
}

func (m *Migrator) WaitForComplete() {
	for {
		clientList := Must(m.from.ClientList().Result())
		for _, client := range strings.Split(clientList, "\r\n") {
			c := ParseClient(client)
			if c["flags"] == "S" {
				a := c["omem"].(int64)
				if m.maxOutBuff == 0 || m.maxOutBuff < a {
					m.maxOutBuff = a
				}
				t := convert(c["obl"].(int64)) + c["oll"].(int64)
				if m.maxOutBuffCommands == 0 || m.maxOutBuffCommands < t {
					m.maxOutBuffCommands = t
				}
			}
		}
		if m.maxOutBuff == 0 && m.maxOutBuffCommands == 0 {
			break
		}
	}
}

func convert(x int64) int64 {
	if x > 0 {
		return 1
	}
	return 0
}

func (m *Migrator) OnComplete() {
	if _, err := m.to.SlaveOf("no", "one").Result(); err != nil {
		panic(err)
	}
	if _, err := m.to.ConfigSet("slave-read-only", "no").Result(); err != nil {
		panic(err)
	}
}

func Must(result string, err error) string {
	if err != nil {
		panic(err)
	}
	return result
}