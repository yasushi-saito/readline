char* _go_readline_complete() {
  extern char*_goReadlineComplete();
  return (char*)_goReadlineComplete();
}
