module Main {

use Factorize;
use GMP;
use Prereqs;

config var min:string = "1";
config var step:string = "1";
config var count:string = "-1";

proc main() {

	initPrimes();

	var minInt:mpz_t; mpz_init(minInt);
	var stepInt:mpz_t; mpz_init(stepInt);
	var countInt:mpz_t; mpz_init(countInt);

	var errMin:int = gmp_sscanf (min, "%Zi", minInt);
	var errStep:int = gmp_sscanf (step, "%Zi", stepInt);
	var errCount:int = gmp_sscanf (count, "%Zi", countInt);

	if (errMin == 0 || errStep == 0 || errCount == 0) {
		writeln("not a valid number: ", min , ", ", step,", or ",count);
		exit(-1);
	}


	var i:mpz_t; mpz_init_set(i, minInt);
	mpz_clear(minInt);

	while mpz_cmp(countInt, zero) != 0 {
		var x:mpz_t; mpz_init(x);
		var y:mpz_t; mpz_init(y);

		factorize(i, x, y);

		if (mpz_cmp(x, zero) != 0 && mpz_cmp(y, zero) != 0) {
			//gmp_printf("%Zi + %Zi %Zi\n", i, x, y);
		} else {
			//gmp_printf("%Zi -\n", i);
		}

		mpz_sub(countInt, countInt, one);

		mpz_clear(x);
		mpz_clear(y);

		mpz_add(i, i, stepInt);
	}

	mpz_clear(i);
	mpz_clear(countInt);
	mpz_clear(stepInt);
}

}
