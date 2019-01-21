#include <poll.h>
#include <readline/history.h>
#include <readline/readline.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

char* _creadline_complete() {
  extern char* _creadlineComplete();
  return (char*)_creadlineComplete();
}

static int _creadline_have_input;
static char* _creadline_input;
static int _creadline_pipefds[2];

void _creadline_init() {
  if (pipe(_creadline_pipefds) != 0) {
    perror("readline: pipe");
    _creadline_pipefds[0] = -1;
    _creadline_pipefds[1] = -1;
  }
}

static void _creadline_input_handler(char* line) {
  _creadline_input = line;
  _creadline_have_input = 1;
  rl_callback_handler_remove();
}

void _creadline_signal_handler(int signo) {
  // Inform the poll() thread about the signal.
  int retval __attribute__((unused));
  retval = write(_creadline_pipefds[1], "x" /*data is irrelevant*/, 1);
}

// Run a readline session. On SIGINT, set *sig=1 and bail.
char* _creadline_readline(const char* prompt, int* sig) {
  void (*old_sigint_handler)(int);
  char* input = 0;
  *sig = 0;
  _creadline_input = 0;
  _creadline_have_input = 0;
  old_sigint_handler = signal(SIGINT, _creadline_signal_handler);
  rl_callback_handler_install(prompt, _creadline_input_handler);
  while (!_creadline_have_input) {
    struct pollfd pfd[2];
    pfd[0].fd = 0;  // stdin
    pfd[0].events = POLLIN;
    pfd[0].revents = 0;
    pfd[1].fd = _creadline_pipefds[0];
    pfd[1].events = POLLIN;
    pfd[1].revents = 0;

    poll(pfd, 2, -1);
    if ((pfd[1].revents & POLLIN) != 0) {  // Signal.
      char buf[128] __attribute__((unused));
      int retval __attribute__((unused));
      retval = read(pfd[1].fd, buf, sizeof(buf));
      *sig = 1;
      break;
    }
    if ((pfd[0].revents & POLLIN) != 0) {
      rl_callback_read_char();
    }
    if ((pfd[0].revents & POLLERR) != 0 || (pfd[1].revents & POLLERR) != 0) {
      *sig = 1;
      // TODO(saito) Report errors more properly.
      break;
    }
  }
  rl_callback_handler_remove();
  signal(SIGINT, old_sigint_handler);
  input = _creadline_input;
  _creadline_input = 0;
  return input;
}
