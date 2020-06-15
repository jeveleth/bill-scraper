
down:
	docker-compose down

build:
	docker-compose up --build -d 

logs-all:
	docker-compose logs -f

run-app:
	docker run --rm -e OPENSTATES_API_KEY=${OPENSTATES_API_KEY} -v $(pwd):/tmp -ti bill-scraper_app:latest ./app -state $(state) -session $(session) -phrase $(phrase) -num-bills $(number) 
	
start:
	-${MAKE} down
	-${MAKE} build

test:
	echo $(foo)
