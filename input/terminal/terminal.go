//go:build !nodefsystem
// +build !nodefsystem

package terminal

/*
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>

#if defined(_WIN32) || defined(_WIN64)
#include <windows.h>
static HANDLE hStdin, hStdout;
static DWORD  hStdinOldMode;
static DWORD  hStdoutOldMode;
#else
#include <unistd.h>
#include <termios.h>
static struct termios term;
#endif

static bool is_raw_enabled = false;

void enableRawMode(void)
{
	if (is_raw_enabled) return;

#if defined(_WIN32) || defined(_WIN64)
	SetConsoleMode(hStdin, ENABLE_VIRTUAL_TERMINAL_INPUT | ENABLE_PROCESSED_INPUT);
	SetConsoleMode(hStdout, ENABLE_VIRTUAL_TERMINAL_PROCESSING | ENABLE_PROCESSED_OUTPUT);
#else
	tcgetattr(STDIN_FILENO, &term);
	term.c_lflag &= ~(ECHO | ICANON);
	tcsetattr(STDIN_FILENO, TCSAFLUSH, &term);
#endif

	is_raw_enabled = true;
}

void disableRawMode(void)
{
	if (!is_raw_enabled) return;

#if defined(_WIN32) || defined(_WIN64)
	SetConsoleMode(hStdin, hStdinOldMode);
	SetConsoleMode(hStdout, hStdoutOldMode);
#else
	term.c_lflag |= ECHO | ICANON;
  	tcsetattr(STDIN_FILENO, TCSAFLUSH, &term);
#endif

	is_raw_enabled = false;
}

void registerAutoDisable(void)
{
#if defined(_WIN32) || defined(_WIN64)
	hStdin = GetStdHandle(STD_INPUT_HANDLE);
	hStdout = GetStdHandle(STD_OUTPUT_HANDLE);
	GetConsoleMode(hStdin, &hStdinOldMode);
	GetConsoleMode(hStdout, &hStdoutOldMode);
#endif
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
