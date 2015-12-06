// +build darwin

// Package nsspeechsynthesizer speaks a notification using NSSpeechSynthesizer.
package nsspeechsynthesizer

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import <Foundation/Foundation.h>
#import <objc/runtime.h>
#include <AppKit/AppKit.h>
#include <errno.h>

int notify2(const char* message, const char* voice) {
    @autoreleasepool {
        NSSpeechSynthesizer *synth = [[NSSpeechSynthesizer alloc] init];

        NSString *v = [NSString stringWithUTF8String:voice];
        if ([v length] != 0) {
            if ([synth setVoice:v] == NO) {
                return 1;
            }
        }

        if ([synth startSpeakingString:[NSString stringWithUTF8String:message]] == NO) {
            return 2;
        }

        while ([NSSpeechSynthesizer isAnyApplicationSpeaking] == YES) {
            [[NSRunLoop currentRunLoop] runUntilDate:[NSDate dateWithTimeIntervalSinceNow:0.1]];
        }
    }
    return 0;
}

const char* availableVoices() {
    NSArray *vs = [NSSpeechSynthesizer availableVoices];
    return [[vs componentsJoinedByString:@"|"] UTF8String];
}
*/
import "C"

import (
	"errors"
	"fmt"
	"strings"
	"unsafe"
)

const (
	voicePrefix = "com.apple.speech.synthesis.voice."
)

var (
	ErrInvalidVoice = errors.New("Voice does not exist")
	ErrFailedSpeech = errors.New("Failed to start speaking")
)

// Notification is a NSSpeechSynthesizer notification.
type Notification struct {
	Message string
	Voice   string

	voices []string
}

// GetMessage returns a notification's message.
func (n *Notification) GetMessage() string {
	return n.Message
}

// SetMessage sets a notification's message.
func (n *Notification) SetMessage(m string) {
	n.Message = m
}

// Notify speaks a notification's message. If the voice field is set, then it'll
// use that voice. It'll return an error if there was a problem executing say.
func (n *Notification) Notify() error {
	m := C.CString(n.Message)
	v := C.CString("")
	defer C.free(unsafe.Pointer(m))
	defer C.free(unsafe.Pointer(v))

	if n.Voice != "" {
		if !n.VoiceExists(n.Voice) {
			return ErrInvalidVoice
		}
		v = C.CString(n.prepare(n.Voice))
	}

	rt, _ := C.notify2(m, v)
	switch rt {
	case 1:
		return fmt.Errorf("Failed to set voice %q", n.Voice)
	case 2:
		return ErrFailedSpeech
	}

	return nil
}

// GetVoices returns a slice of available voices.
func (n *Notification) GetVoices() []string {
	if n.voices != nil && len(n.voices) != 0 {
		// It's already initialized. No need to dive again.
		return n.voices
	}

	vs, _ := C.availableVoices()
	defer C.free(unsafe.Pointer(vs))

	full := strings.Split(C.GoString(vs), "|")
	n.voices = make([]string, 0, len(full))

	for _, v := range full {
		n.voices = append(n.voices, strings.TrimPrefix(v, voicePrefix))
	}

	return n.voices
}

// VoiceExists checks if a given voice exists. It ignores case.
func (n *Notification) VoiceExists(v string) bool {
	for _, voice := range n.GetVoices() {
		if strings.EqualFold(v, voice) {
			return true
		}
	}

	return false
}

// prepare returns a voice with the correct capitalization and voice prefix. If
// v doesn't exist, then an empty string is returned. This is necessary because
// startSpeakingString cares about both of those things.
func (n *Notification) prepare(v string) string {
	for _, voice := range n.GetVoices() {
		if strings.EqualFold(v, voice) {
			return voicePrefix + voice
		}
	}

	return ""
}
