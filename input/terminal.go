package input

/*
#include <stdio.h>
#include <stdlib.h>
#include <termios.h>
#include <unistd.h>
#include <stdbool.h>

static bool is_raw_enabled = false;

void enableRawMode(void)
{
	if (is_raw_enabled) return;

	struct termios t;
	tcgetattr(STDIN_FILENO, &t);
	t.c_lflag &= ~(ECHO | ICANON);
	tcsetattr(STDIN_FILENO, TCSAFLUSH, &t);
	is_raw_enabled = true;
}

void disableRawMode(void)
{
	if (!is_raw_enabled) return;

	struct termios t;
  	tcsetattr(STDIN_FILENO, TCSAFLUSH, &t);
	is_raw_enabled = false;
}

void registerAutoDisable(void)
{
	atexit(disableRawMode);
}
*/
import "C"

func init() {
	C.registerAutoDisable()
}

func EnableRawMode() {
	C.enableRawMode()
}

func DisableRawMode() {
	C.disableRawMode()
}
