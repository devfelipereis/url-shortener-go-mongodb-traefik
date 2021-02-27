up:
	docker-compose up --build --remove-orphans

down:
	docker-compose down

stop:
	docker-compose stop

scale:
	docker-compose up --scale api=3 -d