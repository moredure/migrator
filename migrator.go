package main

import (
	"github.com/go-redis/redis"
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
	m.prepareTo()
	log.Println("Waiting...")
	m.waitForUp()
	log.Println("Waiting for complete...")
	m.waitForComplete()
	log.Println("Finish...")
	m.onComplete()
}

func (m *Migrator) prepareTo() {
	if err := m.to.ConfigSet("slave-read-only", "yes").Err(); err != nil {
		panic(err)
	}
	if err := m.to.ConfigSet("masterauth", m.From.Options().Password).Err(); err != nil {
		panic(err)
	}
	if host, port, err := net.SplitHostPort(m.From.Options().Addr); err != nil {
		panic(err)
	} else if err := m.to.SlaveOf(host, port).Err(); err != nil {
		panic(err)
	}
}

func (m *Migrator) waitForUp() {
	info := new(Info)
	for {
		if err := m.to.Info().Scan(info); err != nil {
			panic(err)
		}
		if info.MasterLinkStatus == "up" {
			return
		}
	}
}

func (m *Migrator) waitForComplete() {
	var clientList ClientList
	for {
		if err := m.From.ClientList().Scan(&clientList); err != nil {
			panic(err)
		}
		for _, client := range clientList {
			if client.Flags == "S" {
				oblAndOll := toBinary(client.Obl) + client.Oll
				if m.maxOutBuff == 0 || m.maxOutBuff < client.Omem {
					m.maxOutBuff = client.Omem
				}
				if m.maxOutBuffCommands == 0 || m.maxOutBuffCommands < oblAndOll {
					m.maxOutBuffCommands = oblAndOll
				}
			}
		}
		if (m.maxOutBuff + m.maxOutBuffCommands) == 0 {
			return
		}
	}
}

func (m *Migrator) onComplete() {
	if err := m.to.SlaveOf("no", "one").Err(); err != nil {
		panic(err)
	}
	if err := m.to.ConfigSet("slave-read-only", "no").Err(); err != nil {
		panic(err)
	}
}
