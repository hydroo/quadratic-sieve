CC = gcc
INCLUDE =
CFLAGS = -DOPENMP --openmp -Wall -g -std=c99 -DDEBUG $(INCLUDE)
#CFLAGS = -Wall -g -std=c99 -DDEBUG $(INCLUDE)
LN = gcc
LIBS = -lgmp -lm -lgomp
LFLAGS = $(LIBS)

BINARY = factorize

OBJECTS = main.o factorize.o prereqs.o
INCLUDES = prereqs.h factorize.h
SOURCES = main.c factorize.c prereqs.c


all: dep $(BINARY)

$(BINARY): $(OBJECTS)
	$(LN) -o $(BINARY) $^ $(LFLAGS)
	@echo ""

%.o: %.c
	$(CC) $(CFLAGS) -c $<

dep: $(SOURCES)
	$(CC) $(CFLAGS) -MM $^ > .depend

-include .depend

.PHONY: clean
clean:
	-rm -f $(BINARY) $(OBJECTS) *.h.gch

