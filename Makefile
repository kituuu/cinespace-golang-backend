build-image:
	docker build -t cine-golang .
run-container: 
	docker run -itd --rm -p 10000:10000 --name cine-golang-backend cine-golang 
	
stop-container:
	docker stop cine-golang-backend

start:
	docker compose up -d
stop:
	docker compose down

