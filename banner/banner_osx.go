// +build darwin

package banner

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa

#import <Foundation/Foundation.h>
#import <objc/runtime.h>
#include <AppKit/AppKit.h>
#include <errno.h>

@implementation NSBundle(noti)
- (NSString *)notiIdentifier {
	return @"com.apple.terminal";
}
@end

@interface NotiDelegate : NSObject<NSUserNotificationCenterDelegate>
@property (nonatomic, assign) BOOL delivered;
@end

@implementation NotiDelegate
- (void) userNotificationCenter:(NSUserNotificationCenter *)center didActivateNotification:(NSUserNotification *)notification {
	self.delivered = YES;
}
- (void) userNotificationCenter:(NSUserNotificationCenter *)center didDeliverNotification:(NSUserNotification *)notification {
	[NSApp activateIgnoringOtherApps:YES];
	self.delivered = YES;
}
@end

void BannerNotify(const char* title, const char* message, const char* sound) {
	errno = 0;
	@autoreleasepool {
		Class nsBundle = objc_getClass("NSBundle");
		method_exchangeImplementations(
			class_getInstanceMethod(nsBundle, @selector(bundleIdentifier)),
			class_getInstanceMethod(nsBundle, @selector(notiIdentifier))
		);

		NotiDelegate *notiDel = [[NotiDelegate alloc]init];
		notiDel.delivered = NO;

		NSUserNotificationCenter *nc = [NSUserNotificationCenter defaultUserNotificationCenter];
		nc.delegate = notiDel;

		NSUserNotification *nt = [[NSUserNotification alloc] init];
		nt.title = [NSString stringWithUTF8String:title];
		nt.informativeText = [NSString stringWithUTF8String:message];

		if ([[NSString stringWithUTF8String:sound] isEqualToString: @"_mute"] == NO) {
			nt.soundName = [NSString stringWithUTF8String:sound];
		}

		[nc deliverNotification:nt];

		while (notiDel.delivered == NO) {
			[[NSRunLoop currentRunLoop] runUntilDate:[NSDate dateWithTimeIntervalSinceNow:0.1]];
		}
	}
}
*/
import "C"
import (
	"unsafe"

	"github.com/variadico/noti"
)

const (
	soundEnv     = "NOTI_SOUND"
	soundFailEnv = "NOTI_SOUND_FAIL"
)

// Notify displays a NSUserNotification.
func Notify(n noti.Params) error {
	var sound string
	if n.Failure {
		sound = n.Config.Get(soundFailEnv)
		if sound == "" {
			sound = "Basso"
		}
	} else {
		sound = n.Config.Get(soundEnv)
		if sound == "" {
			sound = "Ping"
		}
	}

	t := C.CString(n.Title)
	m := C.CString(n.Message)
	s := C.CString(sound)
	defer C.free(unsafe.Pointer(t))
	defer C.free(unsafe.Pointer(m))
	defer C.free(unsafe.Pointer(s))

	C.BannerNotify(t, m, s)

	return nil
}
