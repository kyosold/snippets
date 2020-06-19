#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <stdlib.h>
#include "entities.h"

int main(int argc, char **argv)
{
    char *text = argv[1];
    char *dec = (char *)malloc(strlen(text) + 1);
    if (dec) {
        int n = decode_html_entities_utf8(dec, text);
        printf("decode(%d): %s\n", n, dec);
        free(dec);
    }

    return 0;
}

