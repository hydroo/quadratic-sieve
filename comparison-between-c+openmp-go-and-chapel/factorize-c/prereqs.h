#ifndef DEBUG_HPP
#define DEBUG_HPP


#include <assert.h>
#include <gmp.h>
#include <stdio.h>
#include <time.h>
#include <sys/time.h>


#ifdef DEBUG
#define ASSERT(expression) \
    if ((expression) == 0) { \
        printf("\e[0;31m\033[1mASSERT\033[0m\e[0;30m"); \
        printf(" in %s:%d", __FILE__,  __LINE__); \
        printf(": \"" #expression "\"\n"); \
        assert(0); \
    }
#else
#define ASSERT(expression)
#endif


mpz_t minusOne;
mpz_t zero;
mpz_t one;
mpz_t two;
mpz_t three;
mpz_t four;


const double E;

int isPrimeBruteForceSmallInt(long long int n);

double secondsElapsed();

void initPrimes();


#ifndef OPENMP
int omp_get_max_threads();
int omp_get_thread_num();
#else
#include <omp.h>
#endif

#endif /* DEBUG_HPP */

