package main

import (
	"github.com/microredis/tools/closer"
	"log"
	"net"
)

type redisMigrator struct {
	RedisClients
	maxOutBuff         int64
	maxOutBuffCommands int64
}

func (m *redisMigrator) Migrate() {
	defer closer.MustClose(m.To, m.From)
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
		log.Fatalf("failed to set slave-read-only to yes with error %v\n", err)
	}
	if err := m.To.ConfigSet("masterauth", m.From.Options().Password).Err(); err != nil {
		log.Fatalf("failed to set masterauth password with error %v\n", err)
	}
	if host, port, err := net.SplitHostPort(m.From.Options().Addr); err != nil {
		log.Fatalf("failed to split host port of source redis with error %v\n", err)
	} else if err := m.To.SlaveOf(host, port).Err(); err != nil {
		log.Fatalf("failed to make target slave of host %s port %s with error %v\n", host, port, err)
	}
}

func (m *redisMigrator) waitForUp() {
	info := new(Info)
	for {
		if err := m.To.Info().Scan(info); err != nil {
			log.Fatalf("failed to scan info with error %v\n", err)
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
			log.Fatalf("failed to list clients with error %v\n", err)
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
		log.Fatalf("failed to make target slave of no one with error %v\n", err)
	}
	if err := m.To.ConfigSet("slave-read-only", "no").Err(); err != nil {
		log.Fatalf("failed to make target not slave-read-only with error %v\n", err)
	}
}

func NewMigrator(clients RedisClients) Migrator {
	return &redisMigrator{clients, 0, 0}
}
