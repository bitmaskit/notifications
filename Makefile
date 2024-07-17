# Services
SERVICES = backend frontend router slack

# Targets
.PHONY: all build test run

all: build test run
build:
	@for service in $(SERVICES); do \
		docker build -t $$service ./$$service; \
	done

test:
	 @for service in $(SERVICES); do \
	  if [ -d $$service ]; then \
		cd $$service && go test -v ./...; \
	  else \
		echo "Directory $$service does not exist"; \
	  fi \
	 done

run:
	 docker-compose up

stop:
	 docker-compose down

clean:
	 @for service in $(SERVICES); do \
	  docker rmi $$service; \
	 done
	 docker-compose down --rmi all --volumes --remove-orphans