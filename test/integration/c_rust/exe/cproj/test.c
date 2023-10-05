#include <assert.h>
#include <stdio.h>

extern int add();

int main(void) {
  printf("Starting C tests...\n");
  assert(add(2, 3) == 5);
  assert(add(0, 3) == 3);
  assert(add(0, 0) == 0);
  printf("Finished C tests successfully!\n");
  return 0;
}
