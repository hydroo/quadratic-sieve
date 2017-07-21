module Factorize {

use GMP;
use Math;
use Prereqs;

config var numThreads:int = 1;


proc squareRootCeil(n:mpz_t):mpz_t {

	if (mpz_cmp(n, one) == -1) {
		halt("n cannot be less than one");
	}

	var upperLimit:mpz_t; mpz_init_set_si(upperLimit, 2);
	var lowerLimit:mpz_t; mpz_init_set_si(lowerLimit, 2);

	var upperLimitExp = ceil((mpz_sizeinbase(n, 2):real) / 2.0):uint;
	var lowerLimitExp = floor((mpz_sizeinbase(n, 2):real - 1) / 2.0):uint;

	mpz_pow_ui(upperLimit, upperLimit, upperLimitExp);
	mpz_pow_ui(lowerLimit, lowerLimit, lowerLimitExp);

	var middle:mpz_t; mpz_init(middle);
	var middleSquared:mpz_t; mpz_init(middleSquared);

	/* binary search */
	while mpz_cmp(upperLimit, lowerLimit) != 0 {

		if (mpz_cmp(upperLimit, lowerLimit) == -1) {
			halt("shouldnt happen -> bug");
		}

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

	mpz_clear(lowerLimit);
	mpz_clear(middle);
	mpz_clear(middleSquared);

	return upperLimit;
}

proc factorBase(n:mpz_t) {

	var a:[0..9]int;

	/* calculate 'S' upper bound for the primes to collect */
	var lnn = mpz_sizeinbase(n, 2) * log(2.0);

	if (lnn < 1.0) {
	  /* if this is not done. if n is 1 everything explodes sqrt(<0) = NaN */
	  lnn = 1.0;
	}


	var lnlnn = log(lnn);
	var expo = sqrt(lnn * lnlnn) * 0.5;
	var S = ceil(exp(expo)):int; // magic parameter (wikipedia)

	var primes:[0..S-1]mpz_t;

	var count = 0;
	var P:mpz_t; mpz_init(P);
	var nModP:mpz_t; mpz_init(nModP);
	for p in [2..S] {

		if (isPrimeBruteForceSmallInt(p) == 1) {

			mpz_set_si(P, p);

			mpz_mod(nModP, n, P);

			if(mpz_sizeinbase(nModP, 2) > 63) {
				halt("");
			}

			/* iterate through the whole ring of Z_p and test wether some i^2 equals n -> it is square mod p */
			var isSquare = 0;
			var iSquaredModP:mpz_t; mpz_init(iSquaredModP);
			var i:mpz_t; mpz_init(i);

			while mpz_cmp(i, P) == -1  {
				mpz_powm(iSquaredModP, i, two, P);
				if (mpz_cmp(iSquaredModP, nModP) == 0) {
					isSquare = 1;
					break;
				}
			
				  mpz_add(i, i, one);
			  }

			if (isSquare == 1) {
				mpz_init(primes[count]);
				mpz_set(primes[count], P);
				count += 1;
			}

			mpz_clear(i);
			mpz_clear(iSquaredModP);
		  }
	}


	/* copy results into a new array to save space and prepend -1 */
	var factorbase:[0..count]mpz_t;
	mpz_init(factorbase(0));
	mpz_set_si(factorbase(0), -1);
	for i in [0..count-1] {
		mpz_init(factorbase(i+1)); mpz_set(factorbase(i+1), primes(i));
		mpz_clear(primes(i));
	}

	mpz_clear(P);
	mpz_clear(nModP);

	return factorbase;
}


proc sieveInterval(n:mpz_t):(mpz_t, mpz_t) {

	var lnn = mpz_sizeinbase(n, 2):real * log(2.0);

	if (lnn < 1.0) {
		/* if this is not done. if n is 1 everything explodes sqrt(<0) = NaN */
		lnn = 1.0;
	}

	var lnlnn = log(lnn);
	var exp = ceil(sqrt(lnn * lnlnn) / log(2)):uint;

	var L:mpz_t; mpz_init(L);
	mpz_pow_ui(L, two, exp); // magic parameter (wikipedia)

	var sqrtN = squareRootCeil(n);

	var min:mpz_t; mpz_init(min);
	var max:mpz_t; mpz_init(max);

	mpz_sub(min, sqrtN, L);
	mpz_add(max, sqrtN, L);

	mpz_clear(L);
	mpz_clear(sqrtN);

	return (min, max);
}


proc sieve(n:mpz_t, factorbase:[?factorbasedomain]mpz_t, cMin:mpz_t, cMax:mpz_t) { 

	var intervalBig:mpz_t; mpz_init(intervalBig);
	mpz_sub(intervalBig, cMax, cMin);
	mpz_add(intervalBig, intervalBig, one);
	if (mpz_sizeinbase(intervalBig,2) > 31) {
		halt("interval too large, use big ints (code newly)\n");
	}

	var interval = mpz_get_si(intervalBig):int;
	mpz_clear(intervalBig);

	var retCisTmp:[0..interval-1]mpz_t;
	var retDisTmp:[0..interval-1]mpz_t;
	var retExponentsTmp:[0..interval-1, 0..factorbasedomain.high]int;
	var retValid:[0..interval-1]bool;

	coforall me in [0..numThreads-1] {

		var threads = numThreads;
		var tasksRest = interval % threads;
		var tasksStep = interval / threads;

		var min:int;
		var width:int;

		if (me < tasksRest) {
			min = (tasksStep + 1) * me;
			width = tasksStep + 1;
		} else {
			min = (tasksStep + 1) * (tasksRest);
			min += (me - tasksRest) * tasksStep;
			width = tasksStep;
		}

		var ci:mpz_t; mpz_init_set(ci, cMin);
		var d:mpz_t; mpz_init(d);
		var rest:mpz_t; mpz_init(rest);
		var tmpRest:mpz_t; mpz_init(tmpRest);
		var tmpQuotient:mpz_t; mpz_init(tmpQuotient);

		/* foreach c(i) */
		for h in [min..min+width-1] {

			mpz_add_ui(ci, cMin, h:uint(64));

			retValid(h) = false;

			/* calculate c(i)^1 - n */
			mpz_mul(d, ci, ci);
			mpz_sub(d, d, n);

			/* copy the result because we need to modify it and save c(i)^2 - n as well */
			mpz_set(rest, d);

			for i in factorbasedomain {

				retExponentsTmp(h,i) = 0;

				var repeat = true;

				/* repeat as long as rest % p == 0 -> add 1 to the exponent for each division by p */
				while repeat == true {

					if (i == 0) {
						/* needs special handling: p = -1 */
						if (mpz_sgn(rest) == -1) {
							retExponentsTmp(h,0) = 1;
							mpz_mul(rest, rest, minusOne);
						} else {
							retExponentsTmp(h,0) = 0;
						}
						repeat = false;
						continue;
					}

					mpz_fdiv_qr(tmpQuotient, tmpRest, rest, factorbase[i]);

					if (mpz_cmp(tmpRest, zero) == 0) {
						retExponentsTmp(h,i) += 1;
						mpz_set(rest, tmpQuotient);
					} else {
						  repeat = false;
					}

					if (mpz_cmp(tmpQuotient, zero) == 0) {
						repeat = false;
					}

				}
			}

			/* if rest is 1 c(i)^2 - n has been successfully broken down into a number that can be represented through
			the factor base -> save c(i)^2 and the exponents for the prime factors */
			if (mpz_cmp(rest, one) == 0) {
				mpz_init_set(retDisTmp(h), d);
				mpz_init_set(retCisTmp(h), ci);
				retValid(h) = true;
			}

		}

		mpz_clear(ci);
		mpz_clear(d);
		mpz_clear(rest);
		mpz_clear(tmpRest);
		mpz_clear(tmpQuotient);
	}


	var retIterator = 0;
	for i in [0..interval-1] {
		if (retValid(i) == true) {
			retIterator += 1;
		}
	}

	/* copy the result into new, shorter arrays */
	var retCis:[0..retIterator-1]mpz_t;
	var retDis:[0..retIterator-1]mpz_t;
	var retExponents:[0..retIterator-1, 0..factorbasedomain.high]int;

	var j = 0;
	for i in [0..interval-1] {
		if (retValid(i) == true) {
			mpz_init_set(retCis(j), retCisTmp(i)); // c(i)
			mpz_clear(retCisTmp(i));
			mpz_init_set(retDis(j), retDisTmp(i)); // c(i)^2 - n
			mpz_clear(retDisTmp(i));
			retExponents(j,..) = retExponentsTmp(i,..);
			j += 1;
		}
	}

	return (retCis, retDis, retExponents);
}


proc combineRecursively(start:int, currentExponents:[?zz]int, cMul:mpz_t, dis:[?sievedomain]mpz_t, cis:[?zzz]mpz_t, exponents:[?zzzz]int, factorbase:[?factorbasedomain]mpz_t, n:mpz_t, inout x:mpz_t, inout y:mpz_t) {

	for i in [start..sievedomain.high] {

		var newCurrentExponents:[factorbasedomain]int;
		newCurrentExponents(..) = exponents(i,..) + currentExponents(..);

		var newCMul:mpz_t; mpz_init(newCMul);
		mpz_mul(newCMul, cMul, cis(i));
		mpz_mod(newCMul, newCMul, n);

		var even = true;
		for exp in newCurrentExponents {
			if (exp % 2 == 1) {
				even = false;
				break;
			}
		}

		if (even == true) {

			var b:mpz_t; mpz_init_set(b, one);
			for j in factorbasedomain {

				if (j == 0) {
					continue; /* minus one makes screw ups */
				}

				for k in [0..(newCurrentExponents(j)/2):int-1] {
					mpz_mul(b, b, factorbase(j));
				}
			}

			var a:mpz_t; mpz_init_set(a, newCMul);

			var xTmp:mpz_t; mpz_init(xTmp);
			mpz_add(xTmp, a, b);
			mpz_mod(xTmp, xTmp, n);

			var yTmp:mpz_t; mpz_init(yTmp);
			mpz_sub(yTmp, a, b);
			mpz_mod(yTmp, yTmp, n);

			if (mpz_cmp(xTmp, zero) == 0 || mpz_cmp(xTmp, one) == 0 || mpz_cmp(yTmp, zero) == 0 || mpz_cmp(yTmp, one) == 0) {
				/* no trivial divisors */
				mpz_clear(b);
				mpz_clear(a);
				mpz_clear(xTmp);
				mpz_clear(yTmp);
				mpz_clear(newCMul);
				continue;
			}

			var xTimesY:mpz_t; mpz_init(xTimesY);
			var test:mpz_t; mpz_init(test);
			var testMod:mpz_t; mpz_init(testMod);

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
				mpz_clear(newCMul);
				continue;
			}

			var gcd:mpz_t; mpz_init(gcd);

			while mpz_cmp(test,one) == 1 {

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
				mpz_set(x, xTmp);
				mpz_set(y, yTmp);
				mpz_clear(b);
				mpz_clear(a);
				mpz_clear(xTmp);
				mpz_clear(yTmp);
				mpz_clear(xTimesY);
				mpz_clear(test);
				mpz_clear(testMod);
				mpz_clear(gcd);
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

		combineRecursively(i+1, newCurrentExponents, newCMul, dis, cis, exponents, factorbase, n, x, y);

		mpz_clear(newCMul);

		if (mpz_cmp(x, zero) != 0 && mpz_cmp(y, zero) != 0) {
			return;
		}
	}
}


/* usually this step should be performed by a solver for a system of linear equations in Z_2
but this is too much work for now, so it just tests all (all subsets of the powerset of all exponent vectors)
the linear combinations of exponent vectors for evenness

returns x, y = 0, 0 if nothing is found */
proc combine(dis:[?sievedomain]mpz_t, cis:[?zz]mpz_t, exponents:[?zzz]int,
		factorbase:[?factorbasedomain]mpz_t, n:mpz_t, inout x:mpz_t, inout y:mpz_t) {

	mpz_init_set_si(x, 0);
	mpz_init_set_si(y, 0);
	
	if (sievedomain.high == -1) {
		return;
	}

	var currentExponents:[factorbasedomain]int;
	for i in factorbasedomain {
		currentExponents(i) = 0;
	}

	var cMul:mpz_t; mpz_init_set_si(cMul, 1);


	combineRecursively(0, currentExponents, cMul, dis, cis, exponents, factorbase, n, x, y);

	mpz_clear(cMul);
}


/* returns x, y = 0, 0 upon failure */
proc factorize(n:mpz_t, inout x:mpz_t, inout y:mpz_t) {

	var t1 = secondsElapsed();

	var factorbase = factorBase(n);

	var t2 = secondsElapsed();

	var (min, max) = sieveInterval(n);

	var t3 = secondsElapsed();

	var (cis, dis, exponents) = sieve(n, factorbase, min, max);

	var t4 = secondsElapsed();

	const combinationsLog2:int = 20;

	var cd:subdomain(cis.domain);
	var dd:subdomain(dis.domain);
	var ed:subdomain(exponents.domain);

	if (cis.domain.high >= combinationsLog2-1) {
		/* avoid building too large powersets */
		cd = [0..19];
		dd = [0..19];
		ed = [0..19,exponents.domain.low(2)..exponents.domain.high(2)];
	} else {
		cd = cis.domain;
		dd = dis.domain;
		ed = exponents.domain;
	}

	/* usually you want to solve this through an LGS ... TODO? */
	combine(dis(dd), cis(cd), exponents(ed), factorbase, n, x, y);

	var t5 = secondsElapsed();

	if (mpz_cmp(x, y) == 1) {
		var tmp:mpz_t; mpz_init_set(tmp, x);
		mpz_set(x, y);
		mpz_set(y, tmp);
		mpz_clear(tmp);
	}

	if (mpz_cmp(x, zero) != 0 && mpz_cmp(y, zero) != 0) {
		gmp_printf("%Zi + %Zi %Zi wall %.6f sieve %.6f combining %.6f\n", n, x, y, t5-t1, t4-t3, t5-t4);
	} else {
		gmp_printf("%Zi - - - wall %.6f sieve %.6f combining %.6f\n", n, t5-t1, t4-t3, t5-t4);
	}

	/* free everything */
	for p in factorbase {
		  mpz_clear(p);
	}
	mpz_clear(min);
	mpz_clear(max);

	for i in cis.domain {
		mpz_clear(cis(i));
		mpz_clear(dis(i));
	}
}


}

