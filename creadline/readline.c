#include <readline/history.h>
#include <readline/readline.h>
#include <setjmp.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>

static sigjmp_buf _creadline_jmpbuf;

char* _creadline_complete() {
  extern char* _creadlineComplete();
  return (char*)_creadlineComplete();
}

void _creadline_signal_handler(int signo) {
  if (signo == SIGINT) {
    siglongjmp(_creadline_jmpbuf, 1);
  }
}

// Setup a signal handler, then call readline(). *sig is 0 on normal exit, 1 on
// signal.
char* _creadline_readline(const char* prompt, int* sig) {
  *sig = 0;
  signal(SIGINT, _creadline_signal_handler);
  if (sigsetjmp(_creadline_jmpbuf, 1) != 0) {
    rl_free_line_state();
    rl_cleanup_after_signal();
    *sig = 1;
    return 0;
  }
  char* input = readline(prompt);
  rl_clear_signals();
  return input;
}
