#include <stdio.h>
#include <stdlib.h>
#include <readline/readline.h>
#include <readline/history.h>

char* _go_readline_complete() {
  extern char*_goReadlineComplete();
  return (char*)_goReadlineComplete();
}
