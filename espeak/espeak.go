// +build linux freebsd

// Package espeak can be used to synthesize speech using eSpeak on Linux and
// FreeBSD.
//
// To compile this package, you'll need to install the eSpeak library. On
// Ubuntu, you can install it with this command.
//    sudo apt-get install libespeak-dev
package espeak

/*
#cgo LDFLAGS: -lespeak

#include <stdlib.h>
#include <stdio.h>
#include <errno.h>
#include <espeak/speak_lib.h>
#include <string.h>

int notify2(const char* message, const char* voice) {
	errno = 0;
	espeak_Initialize(AUDIO_OUTPUT_PLAYBACK, 500, NULL, 0);

	unsigned int sz = strlen(message)+1;
	espeak_POSITION_TYPE pos_type;
	espeak_Synth(message, sz, 0, pos_type, 0, espeakCHARS_UTF8, NULL, NULL);

	espeak_Synchronize();
}

const char* availableVoices() {
	const espeak_VOICE **voices = espeak_ListVoices(NULL);
	const espeak_VOICE *v;
	int i;
	const char *p;

	for (i = 0; (v = voices[i]) != NULL; i++) {
		printf("%s\n", v->name);
		printf("%s\n", v->languages);
	}

	int sz = sizeof(voices)/sizeof(voices[0]);

	printf("size: %d\n", sz);
	printf("DONE\n");
	fflush(stdout);
	return NULL;
}
*/
import "C"

import "unsafe"

// Notification is an eSpeak notification.
type Notification struct {
	Message string
	Voice   string
}

// GetMessage returns a notification's message.
func (n *Notification) GetMessage() string {
	return n.Message
}

// SetMessage sets a notification's message.
func (n *Notification) SetMessage(m string) {
	n.Message = m
}

// Notify speaks a notification's message. If the voice field is set, then Notify will
// attempt to use that voice. Otherwise, it'll use the default voice.
func (n *Notification) Notify() error {
	m := C.CString(n.Message)
	v := C.CString("")
	defer C.free(unsafe.Pointer(m))
	defer C.free(unsafe.Pointer(v))

	C.availableVoices()
	C.notify2(m, v)
	return nil
}
