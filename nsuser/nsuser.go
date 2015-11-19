// +build darwin

package nsuser

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
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

void Notify(const char* title, const char* subtitle, const char* text, const char* sound, const char* img) {
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

        NSUserNotification *noti = [[NSUserNotification alloc] init];
        noti.title = [NSString stringWithUTF8String:title];
        noti.informativeText = [NSString stringWithUTF8String:text];
        noti.soundName = NSUserNotificationDefaultSoundName;

        if ([[NSString stringWithUTF8String:sound] length] != 0) {
			noti.soundName = [NSString stringWithUTF8String:sound];
        }
		if ([[NSString stringWithUTF8String:subtitle] length] != 0) {
			noti.subtitle = [NSString stringWithUTF8String:subtitle];
		}
		if ([[NSString stringWithUTF8String:img] length] != 0) {
			noti.contentImage = [[NSImage alloc] initWithContentsOfFile:[NSString stringWithUTF8String:img]];
		}

        [nc deliverNotification:noti];

        while (delegate.delivered == NO) {
            [[NSRunLoop currentRunLoop] runUntilDate:[NSDate dateWithTimeIntervalSinceNow:0.1]];
        }
    }
}
*/
import "C"

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
	C.Notify(
		C.CString(n.Title),
		C.CString(n.Subtitle),
		C.CString(n.InformativeText),
		C.CString(n.SoundName),
		C.CString(n.ContentImage),
	)

	return nil
}
