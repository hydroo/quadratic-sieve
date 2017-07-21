#include "prereqs.h"


#include <math.h>


const double E = 2.718281828459045235;


double secondsElapsed() {
    struct timeval timer;
    gettimeofday(&timer, 0);
    return timer.tv_sec+(timer.tv_usec/1000000.0);
}


/* PORT: exactly as in go */
int isPrimeBruteForceSmallInt(long long int n) {
    if (n == 2) {
        return 1;
    } else if (n < 2 || n % 2 == 0) {
        return 0;
    }

    long long int max = (long long int)(sqrt((double)n)) + 1;

    for (long long int i = 3; i <= max; i += 2) {
        if (n % i == 0) {
            return 0;
        }
    }
        
    return 1;
}


void initPrimes() {
    mpz_init(minusOne); mpz_set_si(minusOne, -1);
    mpz_init(zero); mpz_set_si(zero, 0);
    mpz_init(one); mpz_set_si(one, 1);
    mpz_init(two); mpz_set_si(two, 2);
    mpz_init(three); mpz_set_si(three, 3);
    mpz_init(four); mpz_set_si(four, 4);
}


#ifndef OPENMP
int omp_get_max_threads() {
    return 1;
}


int omp_get_thread_num() {
    return 0;
}
#endif

