module Prereqs {

use GMP;
use Math;
use Time;

const E:real = 2.718281828459045235;

var minusOne:mpz_t;
var zero:mpz_t;
var one:mpz_t;
var two:mpz_t;
var three:mpz_t;
var four:mpz_t;


/* returns approx the time since epoch */
proc secondsElapsed():real {

	var y,m,d:int;

	(y, m, d) = Time.getCurrentDate();
	var s:real = Time.getCurrentTime();

	return s + d*60*60*24 + 60*60*24*30*m + (y - 1970)*60*60*24*30*12;
}


/* PORT: exactly as in go */
proc isPrimeBruteForceSmallInt(n:int(64)):bool {

	if (n == 2) {
		return true;
	} else if (n < 2 || n % 2 == 0) {
		return false;
	}

	var max:int(64) = ((sqrt(n:real)) + 1):int(64);

	for i in 3..max {
		if (n % i == 0) {
			return false;
		}
	}

	return true;
}


proc initPrimes() {
	mpz_init(minusOne); mpz_set_si(minusOne, -1);
	mpz_init(zero); mpz_set_ui(zero, 0);
	mpz_init(one); mpz_set_ui(one, 1);
	mpz_init(two); mpz_set_ui(two, 2);
	mpz_init(three); mpz_set_ui(three, 3);
	mpz_init(four); mpz_set_ui(four, 4);
}

}
