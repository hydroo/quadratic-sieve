CHPL = chpl
INCLUDE =
#CFLAGS = -DOPENMP --openmp -Wall -g -std=c99 -DDEBUG $(INCLUDE)
LIBS = -lgmp -lm
CFLAGS = $(LIBS)

BINARY = factorize-binary

SOURCES = main.chpl prereqs.chpl factorize.chpl


all: $(BINARY)

$(BINARY): $(SOURCES)
	$(CHPL) /usr/include/gmp.h $(CFLAGS) -o $(BINARY) $^

.PHONY: clean
clean:
	-rm -f $(BINARY)

