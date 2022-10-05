# Simple cpp app

.PHONY: run
run:
	@g++ main.cpp
	./a.out
	rm a.out

.PHONY: postgresql
postgresql: download build clean

postgres = postgresql-14.5

.PHONY: download
download:
	@mkdir tmp
	@curl -o tmp/postgresql.tar.gz https://ftp.postgresql.org/pub/source/v14.5/$(postgres).tar.gz
	@tar --directory tmp -xvf tmp/postgresql.tar.gz

.PHONY: build
build:
	unset MAKELEVEL && unset MAKEFLAGS && unset MFLAGS && \
	cd tmp/$(postgres) && \
	./configure && \
	make world
	cp -r tmp/$(postgres)/src/bin/* bin

.PHONY: clean
clean:
	@rm -r tmp

.PHONY: delete
delete:
	@rm -rf bin/*