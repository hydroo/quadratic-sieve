#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "factorize.h"


int main(int argc, char** args) {

    initPrimes();

    mpz_t min; mpz_init(min);
    mpz_t step; mpz_init(step);
    mpz_t count; mpz_init(count);


    for (int i = 1; i < argc; i += 1) {

        if (strcmp(args[i],"-h") == 0 || strcmp(args[i],"--help") == 0) {

            printf("factor <min> <step> <count>\n");
            printf("\n");
            printf("    breaks n down into two factors\n");
            printf("\n");
            printf("    default is 1 1 +inf\n");
            exit(-1);

        } else {

            if (args[i][0] != '-') {

                int err = 0;

                if (mpz_cmp_si(min, 0) == 0) {
                    err = gmp_sscanf (args[i], "%Zi", min);
                } else if (mpz_cmp_si(step, 0) == 0) {
                    err = gmp_sscanf (args[i], "%Zi", step);
                } else if (mpz_cmp_si(count, 0) == 0) {
                    err = gmp_sscanf (args[i], "%Zi", count);
                } else {
                    fprintf(stderr, "too many arguments");
                    exit(-1);
                }

                if (err == 0) {
                    fprintf(stderr, "not a valid number: %s", args[i]);
                    exit(-1);
                }

            } else {
                fprintf(stderr, "unknown argument: %s", args[i]);
                exit(-1);
            }
        }
    }

    if (mpz_cmp_si(min, 0) == 0) {
        mpz_set_si(min, 1);
    }

    if (mpz_cmp_si(step, 0) == 0) {
        mpz_set_si(step, 1);
    }

    if (mpz_cmp_si(count, 0) == 0) {
        mpz_set_si(count, -1);
    }

    mpz_t i; mpz_init(i);
    mpz_set(i, min);
    mpz_clear(min);

    for (; mpz_cmp(count, zero) != 0; mpz_add(i, i, step)) {

        mpz_t x; mpz_init(x);
        mpz_t y; mpz_init(y);

        factorize(i, &x, &y);

        if (mpz_cmp(x, zero) != 0 && mpz_cmp(y, zero) != 0) {
            //gmp_printf("%Zi + %Zi %Zi\n", i, x, y);
        } else {
            //gmp_printf("%Zi -\n", i);
        }

        mpz_sub(count, count, one);

        mpz_clear(x);
        mpz_clear(y);
    }

    mpz_clear(i);
    mpz_clear(count);
    mpz_clear(step);
}

