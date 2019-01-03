package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/microredis/redisutil"
	"log"
	"net"
)

type Migrator struct {
	From               *redis.Client
	to                 *redis.Client
	maxOutBuff         int64
	maxOutBuffCommands int64
}

func NewMigrator(fromRedisClient FromRedisClient, toRedisClient ToRedisClient) *Migrator {
	return &Migrator{fromRedisClient, toRedisClient, 0, 0}
}

func (m *Migrator) Migrate() {
	defer m.to.Close()
	defer m.From.Close()
	log.Println("Started...")
	m.PrepareTo()
	log.Println("Waiting...")
	m.WaitForUp()
	log.Println("Waiting for complete...")
	m.WaitForComplete()
	log.Println("Finish...")
	m.OnComplete()
}

func (m *Migrator) PrepareTo() {
	if _, err := m.to.ConfigSet("slave-read-only", "yes").Result(); err != nil {
		panic(err)
	}
	if _, err := m.to.ConfigSet("masterauth", m.From.Options().Password).Result(); err != nil {
		panic(err)
	}
	if host, port, err := net.SplitHostPort(m.From.Options().Addr); err != nil {
		panic(err)
	} else if _, err := m.to.SlaveOf(host, port).Result(); err != nil {
		panic(err)
	}
}

func (m *Migrator) WaitForUp() {
	for {
		info, err := m.to.Info().Result()
		if err != nil {
			panic(err)
		}
		infoParsed := redisutil.ParseInfo(info)
		if infoParsed["master_link_status"].(string) == "up" {
			return
		}
	}
}

func (m *Migrator) WaitForComplete() {
	for {
		clientList, err := m.From.ClientList().Result()
		if err != nil {
			panic(err)
		}
		for _, client := range redisutil.ParseClientList(clientList) {
			if client["flags"] == "S" {
				omem := client["omem"].(int64)
				obl := client["obl"].(int64)
				oll := client["oll"].(int64)
				fmt.Println(omem, obl, oll)
				oblOllSum := toBinary(obl) + oll
				if m.maxOutBuff == 0 || m.maxOutBuff < omem {
					m.maxOutBuff = omem
				}
				if m.maxOutBuffCommands == 0 || m.maxOutBuffCommands < oblOllSum {
					m.maxOutBuffCommands = oblOllSum
				}
			}
		}
		if m.maxOutBuff == 0 && m.maxOutBuffCommands == 0 {
			return
		}
	}
}

func toBinary(x int64) int64 {
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
