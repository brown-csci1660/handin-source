#!/bin/bash

if [ $# -ne 2 ]; then
	echo "Usage: $0 <key-hex> <output-location>" >&2
	exit 1
fi

if [ $(echo -n "$1" | wc -c) -ne 16 ]; then
	echo "key must be 16 hexadecimal characters (8 bytes)" >&2
	exit 1
fi

echo -n "$1" | while read -n 1 c; do
	echo -n "$c" | grep '[0-9a-f]' >/dev/null
	if [ $? -ne 0 ]; then
		echo "key must be 16 hexadecimal characters (8 bytes)" >&2
		exit 1
	fi
done || exit 1 # inner exit will cause while loop (not script) to exit 1

gcc -std=gnu99 *.c -DKEY_HEX='"'$1'"' -o "$2"