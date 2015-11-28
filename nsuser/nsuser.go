// +build darwin

// Package nsuser can be used to display an NSUserNotification on OS X.
package nsuser

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

#include <errno.h>
#import <Foundation/Foundation.h>
#import <objc/runtime.h>
#include <AppKit/AppKit.h>

@implementation NSBundle(noti)
- (NSString *)notiIdentifier {
    return @"com.apple.terminal";
}
@end

@interface NotiDelegate : NSObject<NSUserNotificationCenterDelegate>
@property (nonatomic, assign) BOOL delivered;
@end

@implementation NotiDelegate
- (void)userNotificationCenter:(NSUserNotificationCenter *)center didDeliverNotification:(NSUserNotification *)notification {
    self.delivered = YES;
}
@end

void notify(const char* title, const char* subtitle, const char* text, const char* sound, const char* img) {
	errno = 0;
    @autoreleasepool {
        Class nsBundle = objc_getClass("NSBundle");
        method_exchangeImplementations(
            class_getInstanceMethod(nsBundle, @selector(bundleIdentifier)),
            class_getInstanceMethod(nsBundle, @selector(notiIdentifier))
        );

        NotiDelegate *delegate = [[NotiDelegate alloc]init];
        delegate.delivered = NO;

        NSUserNotificationCenter *nc = [NSUserNotificationCenter defaultUserNotificationCenter];
        nc.delegate = delegate;

		NSString *nSound = [NSString stringWithUTF8String:sound];

        NSUserNotification *nt = [[NSUserNotification alloc] init];
        nt.title = [NSString stringWithUTF8String:title];
        nt.informativeText = [NSString stringWithUTF8String:text];
        nt.soundName = NSUserNotificationDefaultSoundName;

        if ([[NSString stringWithUTF8String:sound] length] != 0) {
			nt.soundName = [NSString stringWithUTF8String:sound];
        }
		if ([[NSString stringWithUTF8String:subtitle] length] != 0) {
			nt.subtitle = [NSString stringWithUTF8String:subtitle];
		}
		if ([[NSString stringWithUTF8String:img] length] != 0) {
			nt.contentImage = [[NSImage alloc] initWithContentsOfFile:[NSString stringWithUTF8String:img]];
		}

        [nc deliverNotification:nt];

        while (delegate.delivered == NO) {
            [[NSRunLoop currentRunLoop] runUntilDate:[NSDate dateWithTimeIntervalSinceNow:0.1]];
        }
    }
}
*/
import "C"

import "unsafe"

type Notification struct {
	Title           string
	Subtitle        string
	InformativeText string
	SoundName       string
	ContentImage    string
}

func (n *Notification) GetTitle() string {
	return n.Title
}

func (n *Notification) SetTitle(t string) {
	n.Title = t
}

func (n *Notification) GetMessage() string {
	return n.InformativeText
}

func (n *Notification) SetMessage(m string) {
	n.InformativeText = m
}

func (n *Notification) Notify() error {
	t := C.CString(n.Title)
	sb := C.CString(n.Subtitle)
	i := C.CString(n.InformativeText)
	sn := C.CString(n.SoundName)
	c := C.CString(n.ContentImage)
	defer C.free(unsafe.Pointer(t))
	defer C.free(unsafe.Pointer(sb))
	defer C.free(unsafe.Pointer(i))
	defer C.free(unsafe.Pointer(sn))
	defer C.free(unsafe.Pointer(c))

	C.notify(t, sb, i, sn, c)

	return nil
}
