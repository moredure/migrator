# Migrator
[![Maintainability](https://api.codeclimate.com/v1/badges/4cd409c1b35085e147af/maintainability)](https://codeclimate.com/github/microredis/migrator/maintainability)

Script for realtime migration from one redis source to another.
Just specify ```$REDIS_SOURCE``` and ```$REDIS_DESTINATION``` and the deal will be done.

To use this tool
Use docker:

```shell
#!/bin/bash

docker run -it --rm \
-e REDIS_SOURCE=redis://localhost:6379 \
-e REDIS_DESTINATION=redis://localhost:6380 \
microredis/migrator
```
