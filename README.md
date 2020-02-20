# Migrator


Script for realtime migration from one redis source to another in Golang. Inspired by RedisLabs https://github.com/RedisLabs/redis-migrate

Google Wire https://github.com/google/wire used to simplify dependency management

Just specify ```$REDIS_SOURCE``` and ```$REDIS_DESTINATION``` and the deal will be done.

Requires replicaof or slaveof commands be available in your redis installation

Licensed under SSPL license. Please do not use this script to provide service to anybody except for your personal needs.

Used by [ScaleChamp](https://scalechamp.com) and [ScalableSpace](https://www.scalablespace.net)
