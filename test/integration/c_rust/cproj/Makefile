CC = gcc
LFLAGS = -L../rustproj/target/release -lrustproj

.PHONY = all clean runtest

all: main.out test.out

%.out: %.o
	${CC} ${LFLAGS} $< -o $@ 

%.o: %.c
	${CC} -c $<

runtest: test.out
	./test.out

clean:
	rm -fr *.out
