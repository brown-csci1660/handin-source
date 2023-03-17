#include <stdlib.h>
#include "rc4.h"

struct rc4_cipher {
	uint32_t s[256];
	uint8_t i, j;
};

rc4_cipher rc4_new_cipher(char *key, int len) {
	rc4_cipher c = malloc(sizeof(*c));
	c->i = 0;
	c->j = 0;
	for (int i = 0; i < 256; i++) {
		c->s[i] = (uint32_t)i;
	}
	uint8_t j = 0;
	for (int i = 0; i < 256; i++) {
		j += ((uint8_t)c->s[i]) + key[i % len];
		uint32_t tmp = c->s[i];
		c->s[i] = c->s[j];
		c->s[j] = tmp;
	}
	return c;
}

void rc4_free(rc4_cipher c) {
	free(c);
}

void rc4_xor_key_stream(rc4_cipher c, char *dst, char *src, int len) {
	uint8_t i = c->i, j = c->j;
	for (int idx = 0; idx < len; idx++) {
		i += 1;
		j += (uint8_t)(c->s[i]);
		uint32_t tmp = c->s[i];
		c->s[i] = c->s[j];
		c->s[j] = tmp;
		dst[idx] = src[idx] ^ (uint8_t)(c->s[(uint8_t)(c->s[i]+c->s[j])]);
	}
	c->i = i;
	c->j = j;
}