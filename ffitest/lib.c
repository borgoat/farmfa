#include <stdlib.h>
#include <stdio.h>
#include "libfarmfa.h"

int test_player() {
    int ret = 0;

    char *public_key_buffer = malloc(sizeof(char) * 128);
    char *private_key_buffer = malloc(sizeof(char) * 128);

    ret = fm_player_create_key(public_key_buffer, private_key_buffer);
    if (ret != 0) return ret;

    printf("public key:\t%s\nprivate key:\t%s\n\n", public_key_buffer, private_key_buffer);
    free(public_key_buffer);
    free(private_key_buffer);

    return 0;
}

int test_dealer() {
    int ret = 0;
    fm_dealer_t h = fm_dealer_init();

    fm_dealer_add_player(h, "key1", "age1g24fwwnw9y37as3h7d67pjl2eepesznydeka2kppqh43kewt0qxq70qxf8");
    fm_dealer_add_player(h, "key2", "age1u683v4u6ef32keav4j5wvd0xn3zmhr59cnkqwvgwl22lddwc6y2scuya6j");
    fm_dealer_add_player(h, "key3", "age18mvvtfm9sr49yh3fdhvgrs9pjg74s3m2m2tpqdmh27pc3gwn65cs8vqsk8");

    fm_dealer_set_note(h, "testing FFI integration");
    fm_dealer_set_secret(h, "HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ");

    fm_dealer_create_tocs(h);

    fm_dealer_free(h);
}

int main() {
    int ret = 0;

    ret = test_player();
    if (ret != 0) return ret;

    ret = test_dealer();
    if (ret != 0) return ret;

    return ret;
}
