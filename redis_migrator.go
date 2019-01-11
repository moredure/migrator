package main

import (
	"log"
	"net"
)

type redisMigrator struct {
	RedisClients
	maxOutBuff         int64
	maxOutBuffCommands int64
}

func (m *redisMigrator) Migrate() {
	defer m.From.Close()
	defer m.To.Close()
	log.Println("Started...")
	m.prepareTo()
	log.Println("Waiting...")
	m.waitForUp()
	log.Println("Waiting...")
	m.waitForComplete()
	log.Println("Finish...")
	m.onComplete()
}

func (m *redisMigrator) prepareTo() {
	if err := m.To.ConfigSet("slave-read-only", "yes").Err(); err != nil {
		panic(err)
	}
	if err := m.To.ConfigSet("masterauth", m.From.Options().Password).Err(); err != nil {
		panic(err)
	}
	if host, port, err := net.SplitHostPort(m.From.Options().Addr); err != nil {
		panic(err)
	} else if err := m.To.SlaveOf(host, port).Err(); err != nil {
		panic(err)
	}
}

func (m *redisMigrator) waitForUp() {
	info := new(Info)
	for {
		if err := m.To.Info().Scan(info); err != nil {
			panic(err)
		}
		if info.MasterLinkStatus == "up" {
			return
		}
	}
}

func (m *redisMigrator) waitForComplete() {
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

func (m *redisMigrator) onComplete() {
	if err := m.To.SlaveOf("no", "one").Err(); err != nil {
		panic(err)
	}
	if err := m.To.ConfigSet("slave-read-only", "no").Err(); err != nil {
		panic(err)
	}
}

func NewMigrator(clients RedisClients) Migrator {
	return &redisMigrator{clients, 0, 0}
}
