include .env
export $(shell sed 's/=.*//' .env)

run:
	@go run -ldflags "-s -w" cmd/movie_to_gif_bot/main.go

docker-run:
	@sudo docker run -p 443:443 -p 80:80 --rm --env-file=.env impu1se/movie_to_gif_bot

docker-db:
	@sudo docker run -d --name db --rm -v ~/go-projects/data/postgres:/var/lib/postgresql/data -e POSTGRES_HOST_AUTH_METHOD=trust -p 5432:5432 postgres