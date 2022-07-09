#include "libfarmfa.h"

int main() {
  fm_dealer_t h = fm_dealer_init();

  fm_dealer_add_player(h, "key1", "age1g24fwwnw9y37as3h7d67pjl2eepesznydeka2kppqh43kewt0qxq70qxf8");
  fm_dealer_add_player(h, "key2", "age1u683v4u6ef32keav4j5wvd0xn3zmhr59cnkqwvgwl22lddwc6y2scuya6j");
  fm_dealer_add_player(h, "key3", "age18mvvtfm9sr49yh3fdhvgrs9pjg74s3m2m2tpqdmh27pc3gwn65cs8vqsk8");

  fm_dealer_set_note(h, "testing FFI integration");
  fm_dealer_set_secret(h, "HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ");

  fm_dealer_create_tocs(h);

  fm_dealer_free(h);
}
