service-up:
	sudo docker-compose up --remove-orphans --build

service-run:
	sudo docker-compose up

service-down:
	sudo docker-compose down

.PHONY: service-up service-run service-down