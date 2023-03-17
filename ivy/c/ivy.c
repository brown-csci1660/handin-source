#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>
#include <error.h>
#include <time.h>
#include "rc4.h"

void encryptPacket(uint16_t iv, char *k, char *m, char *c, int len);
void nextIV();
void printHex(char *buf, int len);
void parseHex(char *str, char *buf, int buflen);

void encryptPacket(uint16_t iv, char *k, char *m, char *c, int len) {
	char *key = malloc(10);
	key[0] = iv >> 8;
	key[1] = iv;
	for (int i = 0; i < 8; i++) {
		key[i + 2] = k[i];
	}
	rc4_cipher cipher = rc4_new_cipher(key, 10);
	rc4_xor_key_stream(cipher, c, m, len);
	rc4_free(cipher);
	free(key);
}

uint16_t iv;

// generate a new random IV
void nextIV() {
	iv = (uint16_t)rand();
}

int main() {
	unsigned char *key = malloc(8);
	// KEY_HEX should be defined at compile time
	// (for example, using '-DKEY_HEX="0011223344556677"');
	// we then initialize the key by parsing KEY_HEX
	// as a hexadecimal string
	parseHex(KEY_HEX, key, 8);

	srand(time(NULL));
	nextIV();

	char *c = malloc(8);
	encryptPacket(iv, key, key, c, 8);
	printf("%04x ", iv);
	printHex(c, 8);
	printf("\n");

	char *m = malloc(8);
	while (1) {
		// we're lazy - only handle 8-byte payloads (plus newline)
		char *linebuf = malloc(17);
		int code = read(0, linebuf, 17);
		// we're really lazy - only handle getting everything
		// at once from the read syscall
		if (code != 17) {
			if (code == 0) {
				// Got EOF
				return 0;
			}
			if (code > 0) {
				printf("Error: got partial read (plaintext must be 16 hexadecimal characters)\n");
				return 1;
			}
			perror("Error reading:");
			return 1;
		}

		parseHex(linebuf, m, 8);
		nextIV();
		encryptPacket(iv, key, m, c, 8);
		printf("%04x ", iv);
		printHex(c, 8);
		printf("\n");
	}
}

void printHex(char *buf, int len) {
	unsigned char *c = (unsigned char*)buf;
	for (int i = 0; i < len; i++) {
		printf("%02x", c[i]);
	}
}

void parseHex(char *str, char *buf, int buflen) {
	unsigned char *b = (unsigned char*)buf;
		for (int i = 0; i < buflen; i++) {
		int j = i * 2;
		unsigned char val = 0;
		if (str[j] >= '0' && str[j] <= '9') {
			val = str[j] - '0';
		} else {
			val = str[j] - 'a' + 10;
		}
		val <<= 4;
		if (str[j+1] >= '0' && str[j+1] <= '9') {
			val += str[j+1] - '0';
		} else {
			val += str[j+1] - 'a' + 10;
		}
		b[i] = val;
	}
}
