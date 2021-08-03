#!/bin/bash

docker run --rm \
	--name preed-db \
	-e POSTGRES_USER=preed \
	-e POSTGRES_PASSWORD=preed \
	-e POSTGRES_DB=preed \
	-p 5432:5432 \
 	-v $HOME/.preed/db:/var/lib/postgresql/data \
	-d \
	postgres:10
