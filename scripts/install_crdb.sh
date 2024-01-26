#!/bin/bash

LINUX_VERSION=v23.1.14.linux-amd64
MAC_VERSION=v23.1.14.darwin-10.9-amd64

download_crdb() {
	mkdir -p ./crdbTemp

	if [ "$(expr substr $(uname -s) 1 6)" == "Darwin" ]; then
		curl https://binaries.cockroachdb.com/cockroach-$MAC_VERSION.tgz -o ./crdbTemp/crdb.tgz
		curl https://binaries.cockroachdb.com/cockroach-sql-$MAC_VERSION.tgz -o ./crdbTemp/crdbshell.tgz
	elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
		curl https://binaries.cockroachdb.com/cockroach-$LINUX_VERSION.tgz -o ./crdbTemp/crdb.tgz
		curl https://binaries.cockroachdb.com/cockroach-sql-$LINUX_VERSION.tgz -o ./crdbTemp/crdbshell.tgz
	else
		echo "Unknown architecture $(uname -s)"
		rm -rf ./crdbTemp
		return 1
	fi
}

uncompress_crdb() {
	tar -xvzf ./crdbTemp/crdb.tgz -C ./crdbTemp/
	tar -xvzf ./crdbTemp/crdbshell.tgz -C ./crdbTemp/
}

install_crdb() {
	sudo mkdir -p /usr/local/lib/cockroach

	sudo cp -i ./crdbTemp/*/lib/libgeos* /usr/local/lib/cockroach/
	sudo cp -i ./crdbTemp/*/cockroach /usr/local/bin/
	sudo cp -i ./crdbTemp/*/cockroach-sql /usr/local/bin/

	rm -rf ./crdbTemp
}

download_crdb
uncompress_crdb
install_crdb