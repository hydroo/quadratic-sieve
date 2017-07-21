#include "factorize.h"

#include <math.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>


static void squareRootCeil(mpz_t n, mpz_t *r);
static void factorBase(mpz_t n, mpz_t **factorbase, int *factorbaselen );
static void sieveInterval(mpz_t n, mpz_t *min, mpz_t *max);
static void sieve(mpz_t n, mpz_t* factorbase, int factorbaselen, mpz_t cMin,
        mpz_t cMax, mpz_t **retCis, mpz_t **retDis, int ***retExponents, int *sievelen);
static void combineRecursively(int start, int *currentExponents, mpz_t cMul, mpz_t *dis, mpz_t *cis, int **exponents,
        int sievelen, mpz_t *factorBase, int factorbaselen, mpz_t n, mpz_t *x, mpz_t *y);
static void combine(mpz_t *dis, mpz_t *cis, int **exponents, int sievelen,
        mpz_t *factorbase, int factorbaselen, mpz_t n, mpz_t *x, mpz_t *y);

/* PORT: exactly as in go */
void squareRootCeil(mpz_t n, mpz_t *r) {

    ASSERT(mpz_cmp(n, one) >= 0);

    mpz_t upperLimit; mpz_init(upperLimit); mpz_set_si(upperLimit, 2);
    mpz_t lowerLimit; mpz_init(lowerLimit); mpz_set_si(lowerLimit, 2);

    int upperLimitExp = (int)(ceil(((double)mpz_sizeinbase(n, 2)) / 2.0));
    int lowerLimitExp = (int)(floor(((double)mpz_sizeinbase(n, 2) - 1) / 2.0));

    mpz_pow_ui(upperLimit, upperLimit, upperLimitExp);
    mpz_pow_ui(lowerLimit, lowerLimit, lowerLimitExp);

    mpz_t middle; mpz_init(middle);
    mpz_t middleSquared; mpz_init(middleSquared);

    /* binary search */
    for (; mpz_cmp(upperLimit, lowerLimit) != 0;) {

        ASSERT(mpz_cmp(upperLimit, lowerLimit) != -1);

        mpz_add(middle,upperLimit, lowerLimit);
        mpz_fdiv_q(middle, middle, two);

        mpz_pow_ui(middleSquared, middle, 2);

        if (mpz_cmp(middleSquared, n) == -1) {

            if (mpz_cmp(lowerLimit, middle) == 0) {
                break;
            }

            mpz_set(lowerLimit, middle);
        } else {

            mpz_set(upperLimit, middle);
        }
    }

    mpz_set(*r, upperLimit);

    mpz_clear(upperLimit);
    mpz_clear(lowerLimit);
    mpz_clear(middle);
    mpz_clear(middleSquared);

    return;
}

/* PORT: exactly as in go */
void factorBase(mpz_t n, mpz_t **factorbase, int *factorbaselen )  {

    /* calculate 'S' upper bound for the primes to collect */
    double lnn = ((double) mpz_sizeinbase(n, 2)) * log(2.0);

    if (lnn < 1.0) {
      /* if this is not done. if n is 1 everything explodes sqrt(<0) = NaN */
      lnn = 1.0;
    }

    double lnlnn = log(lnn);
    double exp = sqrt(lnn * lnlnn) * 0.5;
    int S = (int)ceil(pow(E, exp)); // magic parameter (wikipedia)

    mpz_t *primes = (mpz_t*) malloc(sizeof(mpz_t) * S);

    int count = 0;
    mpz_t P; mpz_init(P);
    mpz_t nModP; mpz_init(nModP);
    for (int p = 2; p <= S; p += 1) {

        if (isPrimeBruteForceSmallInt(p) == 1) {

            mpz_set_si(P, p);

            mpz_mod(nModP, n, P);

            ASSERT(mpz_sizeinbase(nModP, 2) <= 63)

            /* iterate through the whole ring of Z_p and test wether some i^2 equals n -> it is square mod p */
            int isSquare = 0;
            mpz_t iSquaredModP;
            mpz_init(iSquaredModP);
            mpz_t i;
            mpz_init(i);

            for (;mpz_cmp(i, P) == -1; mpz_add(i, i, one)) {
                mpz_powm(iSquaredModP, i, two, P);
                if (mpz_cmp(iSquaredModP, nModP) == 0) {
                    isSquare = 1;
                    break;
                }
            }

            if (isSquare == 1) {
                mpz_init(primes[count]);
                mpz_set(primes[count], P);
                count += 1;
            }

            mpz_clear(i);;
            mpz_clear(iSquaredModP);;
        }
    }

    /* copy results into a new array to save space */
    *factorbase = (mpz_t*) malloc(sizeof(mpz_t) * (count + 1));
    mpz_init((*factorbase)[0]);
    mpz_set_si((*factorbase[0]), -1);
    for (int i = 0; i < count; i += 1) {
        mpz_init((*factorbase)[i+1]); mpz_set((*factorbase)[i+1], primes[i]);
        mpz_clear(primes[i]);
    }

    free(primes);
    mpz_clear(P);
    mpz_clear(nModP);

    *factorbaselen = count + 1;
}


/* PORT: exactly as in go */
void sieveInterval(mpz_t n, mpz_t *min, mpz_t *max) {

    double lnn = ((double) mpz_sizeinbase(n, 2)) * log(2.0);

    if (lnn < 1.0) {
      /* if this is not done. if n is 1 everything explodes sqrt(<0) = NaN */
      lnn = 1.0;
    }

    double lnlnn = log(lnn);
    int exp = (int) ceil(sqrt(lnn * lnlnn) / log(2));

    mpz_t L; mpz_init(L);
    mpz_pow_ui(L, two, exp); // magic parameter (wikipedia)

    mpz_t sqrtN; mpz_init(sqrtN);
    squareRootCeil(n, &sqrtN);

    mpz_sub(*min, sqrtN, L);
    mpz_add(*max, sqrtN, L);

    mpz_clear(L);
    mpz_clear(sqrtN);
}


/* PORT: exactly as in go */
void sieve(mpz_t n, mpz_t* factorbase, int factorbaselen, mpz_t cMin,
        mpz_t cMax, mpz_t **retCis, mpz_t **retDis, int ***retExponents, int *sievelen) {

    mpz_t intervalBig; mpz_init(intervalBig);
    mpz_sub(intervalBig, cMax, cMin);
    mpz_add(intervalBig, intervalBig, one);
    ASSERT(mpz_sizeinbase(intervalBig,2) <= 31);
    int interval = (int) mpz_get_si(intervalBig);
    mpz_clear(intervalBig);

    mpz_t *retDisTmp = (mpz_t*) malloc(sizeof(mpz_t) * interval);
    mpz_t *retCisTmp = (mpz_t*) malloc(sizeof(mpz_t) * interval);
    int **retExponentsTmp = (int**) malloc(sizeof(int*) * interval);

    int threads = omp_get_max_threads();

    #pragma omp parallel firstprivate(n, factorbase, factorbaselen, cMin, cMax, interval, retDisTmp, retCisTmp, retExponentsTmp, threads) num_threads(omp_get_max_threads())
    {
        int tasksRest = interval % threads;
        int tasksStep = interval / threads;

        int me = omp_get_thread_num();

        int min;
        int width;

        if (me < tasksRest) {
            min = (tasksStep + 1) * me;
            width = tasksStep + 1;
        } else {
            min = (tasksStep + 1) * (tasksRest);
            min += (me - tasksRest) * tasksStep;
            width = tasksStep;
        }

        mpz_t ci; mpz_init(ci);
        mpz_set(ci, cMin);
        mpz_t d; mpz_init(d);
        mpz_t rest; mpz_init(rest);
        mpz_t tmpRest; mpz_init(tmpRest);
        mpz_t tmpQuotient; mpz_init(tmpQuotient);
        int *exponents = malloc(sizeof(int) * factorbaselen);

        /* foreach c(i) */
        for (int h = min; h < min + width; h += 1) {

            mpz_add_ui(ci, cMin, h);

            retExponentsTmp[h] = NULL; /* used to determine wether this entry (h) is present/valid */

            /* calculate c(i)^1 - n */
            mpz_mul(d, ci, ci);
            mpz_sub(d, d, n);

            /* copy the result because we need to modify it and save c(i)^2 - n as well */
            mpz_set(rest, d);

            for (int i = 0; i < factorbaselen; i += 1) {

                exponents[i] = 0;

                int repeat = 1;

                /* repeat as long as rest % p == 0 -> add 1 to the exponent for each division by p */
                for (;repeat == 1;) {

                    if (i == 0) {
                        /* needs special handling: p = -1 */
                        if (mpz_sgn(rest) == -1) {
                            exponents[0] = 1;
                            mpz_mul(rest, rest, minusOne);
                        } else {
                            exponents[0] = 0;
                        }
                        repeat = 0;
                        continue;
                    }

                    mpz_fdiv_qr(tmpQuotient, tmpRest, rest, factorbase[i]);

                    if (mpz_cmp(tmpRest, zero) == 0) {
                        exponents[i] += 1;
                        mpz_set(rest, tmpQuotient);
                    } else {
                        repeat = 0;
                    }

                    if (mpz_cmp(tmpQuotient, zero) == 0) {
                        repeat = 0;
                    }

                }
            }

            /* if rest is 1 c(i)^2 - n has been successfully broken down into a number that can be represented through
            the factor base -> save c(i)^2 and the exponents for the prime factors */
            if (mpz_cmp(rest, one) == 0) {
                mpz_init_set(retDisTmp[h], d);
                mpz_init_set(retCisTmp[h], ci);
                retExponentsTmp[h] = (int*) malloc(sizeof(int) * factorbaselen);
                memcpy(retExponentsTmp[h], exponents, sizeof(int) * factorbaselen);
            }

        }

        mpz_clear(ci);
        mpz_clear(d);
        mpz_clear(rest);
        mpz_clear(tmpRest);
        mpz_clear(tmpQuotient);
        free(exponents);
    }



    int retIterator = 0;
    for (int i = 0; i < interval; i += 1) {
        if (retExponentsTmp[i] != NULL) {
            retIterator += 1;
        }
    }


    /* copy the result into new, shorter arrays */
    *retCis = (mpz_t*) malloc(sizeof(mpz_t) * retIterator);
    *retDis = (mpz_t*) malloc(sizeof(mpz_t) * retIterator);
    *retExponents = (int**) malloc(sizeof(int*) * retIterator);

    *sievelen = retIterator;

    for (int i = 0, j = 0; i < interval; i += 1) {
        if (retExponentsTmp[i] != NULL) {
            mpz_init_set((*retCis)[j], retCisTmp[i]); // c(i)
            mpz_clear(retCisTmp[i]);
            mpz_init_set((*retDis)[j], retDisTmp[i]); // c(i)^2 - n
            mpz_clear(retDisTmp[i]);
            (*retExponents)[j] = retExponentsTmp[i];

            j += 1;
        }
    }
    free(retCisTmp);
    free(retDisTmp);
    free(retExponentsTmp);
}


/* PORT: exactly as in go */
void combineRecursively(int start, int *currentExponents, mpz_t cMul, mpz_t *dis, mpz_t *cis, int **exponents,
        int sievelen, mpz_t *factorbase, int factorbaselen, mpz_t n, mpz_t *x, mpz_t *y) {

    for (int i = start; i < sievelen; i += 1) {

        int *newCurrentExponents = (int*) malloc(sizeof(int) * factorbaselen);
        for (int j = 0; j < factorbaselen; j += 1) {
            newCurrentExponents[j] = exponents[i][j] + currentExponents[j];
        }

        mpz_t newCMul; mpz_init(newCMul);
        mpz_mul(newCMul, cMul, cis[i]);
        mpz_mod(newCMul, newCMul, n);

        int even = 1;
        for (int j = 0; j < factorbaselen; j += 1) {
            if (newCurrentExponents[j] % 2 == 1) {
                even = 0;
                break;
            }
        }

        if (even == 1) {

            mpz_t b; mpz_init_set(b, one);
            for (int j = 0; j < factorbaselen; j += 1) {

                if (j == 0) {
                    continue; /* minus one makes screw ups */
                }

                for (int k = 0; k < newCurrentExponents[j]/2; k += 1) {
                    mpz_mul(b, b, factorbase[j]);
                }
            }

            mpz_t a; mpz_init_set(a, newCMul);

            mpz_t xTmp; mpz_init(xTmp);
            mpz_add(xTmp, a, b);
            mpz_mod(xTmp, xTmp, n);

            mpz_t yTmp; mpz_init(yTmp);
            mpz_sub(yTmp, a, b);
            mpz_mod(yTmp, yTmp, n);



            if (mpz_cmp(xTmp, zero) == 0 || mpz_cmp(xTmp, one) == 0 || mpz_cmp(yTmp, zero) == 0 || mpz_cmp(yTmp, one) == 0) {
                /* no trivial divisors */
                mpz_clear(b);
                mpz_clear(a);
                mpz_clear(xTmp);
                mpz_clear(yTmp);
                free(newCurrentExponents);
                mpz_clear(newCMul);
                continue;
            }

            mpz_t xTimesY; mpz_init(xTimesY);
            mpz_t test; mpz_init(test);
            mpz_t testMod; mpz_init(testMod);

            mpz_mul(xTimesY, xTmp, yTmp);
            mpz_fdiv_qr(test, testMod, xTimesY, n);

            if (mpz_cmp(testMod,zero) != 0) {
                mpz_clear(b);
                mpz_clear(a);
                mpz_clear(xTmp);
                mpz_clear(yTmp);
                mpz_clear(xTimesY);
                mpz_clear(test);
                mpz_clear(testMod);
                free(newCurrentExponents);
                mpz_clear(newCMul);
                continue;
            }

            mpz_t gcd; mpz_init(gcd);

            for (;mpz_cmp(test,one) == 1;) {

                mpz_gcd(gcd, xTmp, test);

                if (mpz_cmp(gcd, one) == 1) {
                    mpz_fdiv_q(xTmp, xTmp, gcd);
                    mpz_fdiv_q(test, test, gcd);
                }

                mpz_gcd(gcd, yTmp, test);

                if (mpz_cmp(gcd, one) == 1) {
                    mpz_fdiv_q(yTmp, yTmp, gcd);
                    mpz_fdiv_q(test, test, gcd);
                }

            }

            mpz_mul(xTimesY, xTmp, yTmp);

            if (mpz_cmp(xTimesY, n) == 0) {
                mpz_set(*x, xTmp);
                mpz_set(*y, yTmp);

                mpz_clear(b);
                mpz_clear(a);
                mpz_clear(xTmp);
                mpz_clear(yTmp);
                mpz_clear(xTimesY);
                mpz_clear(test);
                mpz_clear(testMod);
                mpz_clear(gcd);
                free(newCurrentExponents);
                mpz_clear(newCMul);
                return;
            }

            mpz_clear(b);
            mpz_clear(a);
            mpz_clear(xTmp);
            mpz_clear(yTmp);
            mpz_clear(xTimesY);
            mpz_clear(test);
            mpz_clear(testMod);
            mpz_clear(gcd);
        }

        combineRecursively(i+1, newCurrentExponents, newCMul, dis, cis, exponents,
                sievelen, factorbase, factorbaselen, n, x, y);

        free(newCurrentExponents);
        mpz_clear(newCMul);

        if (mpz_cmp(*x, zero) != 0 && mpz_cmp(*y, zero) != 0) {
            return;
        }
    }
}


/* usually this step should be performed by a solver for a system of linear equations in Z_2
but this is too much work for now, so it just tests all (all subsets of the powerset of all exponent vectors)
the linear combinations of exponent vectors for evenness

returns x, y = 0, 0 if nothing is found */
/* PORT: exactly as in go */
void combine(mpz_t *dis, mpz_t *cis, int **exponents, int sievelen,
        mpz_t *factorbase, int factorbaselen, mpz_t n, mpz_t *x, mpz_t *y) {

    mpz_set(*x, zero);
    mpz_set(*y, zero);

    if (sievelen == 0) {
        return;
    }

    int *currentExponents = (int*) malloc(sizeof(int) * factorbaselen);
    for (int i = 0; i < factorbaselen; i += 1) {
        currentExponents[i] = 0;
    }

    mpz_t cMul; mpz_init_set_si(cMul, 1);

    combineRecursively(0, currentExponents, cMul, dis, cis, exponents, sievelen, factorbase, factorbaselen, n, x, y);

    free(currentExponents);
    mpz_clear(cMul);
}


/* returns x, y = 0, 0 upon failure */
void factorize(mpz_t n, mpz_t *x, mpz_t *y) {

    double t1 = secondsElapsed();

    int factorbaselen;
    mpz_t *factorbase;
    factorBase(n, &factorbase, &factorbaselen);

    double t2 = secondsElapsed();

    mpz_t min; mpz_init(min);
    mpz_t max; mpz_init(max);
    sieveInterval(n, &min, &max);

    double t3 = secondsElapsed();

    mpz_t *cis;
    mpz_t *dis;
    int **exponents;
    int sievelen;

    sieve(n, factorbase, factorbaselen, min, max, &cis, &dis, &exponents, &sievelen);

    double t4 = secondsElapsed();

    const int combinationsLog2 = 20;

    if (sievelen > combinationsLog2) {
        /* avoid building too large powersets */
        sievelen = 20;
    }


    /* usually you want to solve this through an LGS ... TODO? */
    combine(dis, cis, exponents, sievelen, factorbase, factorbaselen, n, x, y);

    double t5 = secondsElapsed();

    char s[1000] = "";
    sprintf(s,"%f%f%f%f%f", t1,t2,t3,t4,t5); // Q_UNUSED(...), (void) ...

    if (mpz_cmp(*x, *y) == 1) {
        mpz_t tmp; mpz_init_set(tmp, *x);
        mpz_set(*x, *y);
        mpz_set(*y, tmp);
        mpz_clear(tmp);
    }

    if (mpz_cmp(*x, zero) != 0 && mpz_cmp(*y, zero) != 0) {
        gmp_printf("%Zi + %Zi %Zi wall %.6f sieve %.6f combining %.6f\n", n, *x, *y, t5-t1, t4-t3, t5-t4);
    } else {
        gmp_printf("%Zi - - - wall %.6f sieve %.6f combining %.6f\n", n, t5-t1, t4-t3, t5-t4);
    }

    /* free everything */
    for (int i = 0;i < factorbaselen; i += 1) {
        mpz_clear(factorbase[i]);
    }
    free(factorbase);
    mpz_clear(min);
    mpz_clear(max);

    for (int i = 0; i < sievelen; i += 1) {
        mpz_clear(cis[i]);
        mpz_clear(dis[i]);
        free(exponents[i]);
    }
    free(cis);
    free(dis);
    free(exponents);
}

