#include <stdio.h>
#include <stdlib.h>
#include <readline/readline.h>
#include <readline/history.h>

char* _go_readline_complete() {
  extern char*_goReadlineComplete();
  return (char*)_goReadlineComplete();
}

int _go_history_len() {
  HIST_ENTRY** e = history_list();
  if (e == NULL) return 0;
  int i = 0;
  while (e[i] != NULL) i++;
  return i;
}
