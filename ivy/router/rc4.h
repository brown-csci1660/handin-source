#include <stdint.h>

typedef struct rc4_cipher *rc4_cipher;

rc4_cipher rc4_new_cipher(char *key, int len);
void rc4_free(rc4_cipher c);

void rc4_xor_key_stream(rc4_cipher c, char *dst, char *src, int len);