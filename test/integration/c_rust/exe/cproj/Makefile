CC = gcc
LFLAGS = -L../../lib/rustproj/target/release -lrustproj

.PHONY = all clean runtest

all: main.out test.out

%.out: %.o
	${CC} $< -o $@ ${LFLAGS}

%.o: %.c
	${CC} -c $<

runtest: test.out
	./test.out

clean:
	rm -fr *.out
