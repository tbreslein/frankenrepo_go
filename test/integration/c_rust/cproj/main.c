#include <stdio.h>

extern int add();

int main(void) {
  printf("2 + 3 = %d\n", add(2, 3));
  return 0;
}
