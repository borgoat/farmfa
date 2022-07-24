#include <stdlib.h>
#include <stdio.h>
#include "libfarmfa.h"

int test_player() {
    int ret = 0;

    fm_keypair keypair;

    ret = fm_player_create_key(&keypair);
    if (ret != 0) return ret;

    printf("public key:\t%s\nprivate key:\t%s\n\n", keypair.public_key, keypair.private_key);
    fm_keypair_free(&keypair);

    return 0;
}

int test_dealer() {
    int ret = 0;
    fm_dealer_t h;
    fm_dealer_init(&h);

    fm_dealer_add_player(h, "recipient1@example.com", "age1g24fwwnw9y37as3h7d67pjl2eepesznydeka2kppqh43kewt0qxq70qxf8");
    fm_dealer_add_player(h, "recipient2@example.com", "age1u683v4u6ef32keav4j5wvd0xn3zmhr59cnkqwvgwl22lddwc6y2scuya6j");
    fm_dealer_add_player(h, "recipient3@example.com", "age18mvvtfm9sr49yh3fdhvgrs9pjg74s3m2m2tpqdmh27pc3gwn65cs8vqsk8");

    fm_dealer_set_note(h, "testing FFI integration");
    fm_dealer_set_secret(h, "HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ");

    fm_encrypted_tocs enc_tocs;

    fm_dealer_create_tocs(h, &enc_tocs);

    for (int i = 0; i < enc_tocs.length; i++) {
        fm_encrypted_toc t = enc_tocs.items[i];
        printf("Recipient: %s\n", t.recipient);
        printf("Toc:\n%s\n", t.encrypted_toc);
    }

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
