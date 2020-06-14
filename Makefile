
down:
	docker-compose down

build:
	docker-compose up --build -d 

logs-all:
	docker-compose logs -f

run-app:
	docker run -v $(pwd):/tmp -ti bill-scraper_app:latest ./app -state $(state) -session $(session) -phrase $(phrase) -num-bills $(number) 
	
start:
	-${MAKE} down

test:
	echo $(foo)
