// Based on https://github.com/norio-nomura/usernotification (WTFPL, 2013)

#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>
#import <objc/runtime.h>

@implementation NSBundle (swizle)

// Overriding bundleIdentifier works, but overriding NSUserNotificationAlertStyle does not work.
- (NSString *)__bundleIdentifier
{
    if (self == [NSBundle mainBundle]) {
        return @"com.apple.terminal";
    }

    return [self __bundleIdentifier];
}

@end

BOOL installNSBundleHook()
{
    Class c = objc_getClass("NSBundle");
    if (c) {
        method_exchangeImplementations(class_getInstanceMethod(c, @selector(bundleIdentifier)),
                                       class_getInstanceMethod(c, @selector(__bundleIdentifier)));
        return YES;
    }

    return NO;
}

@interface NotificationCenterDelegate : NSObject <NSUserNotificationCenterDelegate>

@property (nonatomic, assign) BOOL keepRunning;

@end

@implementation NotificationCenterDelegate

- (void)userNotificationCenter:(NSUserNotificationCenter *)center didDeliverNotification:(NSUserNotification *)notification
{
    self.keepRunning = NO;
}

- (BOOL)userNotificationCenter:(NSUserNotificationCenter *)center shouldPresentNotification:(NSUserNotification *)notification
{
    return YES;
}

@end

void Send(const char *title, const char *subtitle, const char *informativeText, const char *contentImage, const char *soundName)
{
    @autoreleasepool {
        if (!installNSBundleHook()) {
            return;
        }

        NSUserNotificationCenter *nc = [NSUserNotificationCenter defaultUserNotificationCenter];
        NotificationCenterDelegate *ncDelegate = [[NotificationCenterDelegate alloc] init];
        ncDelegate.keepRunning = YES;
        nc.delegate = ncDelegate;

        NSUserNotification *note = [[NSUserNotification alloc] init];
        note.title = [NSString stringWithUTF8String:title];
        note.subtitle = [NSString stringWithUTF8String:subtitle];
        note.informativeText = [NSString stringWithUTF8String:informativeText];
        note.soundName = [NSString stringWithUTF8String:soundName];
        // note.contentImage = [[NSImage alloc] initWithContentsOfFile:[NSString stringWithUTF8String:contentImage]];
		[note setValue:[[NSImage alloc] initWithContentsOfFile:[NSString stringWithUTF8String:contentImage]] forKey:@"_identityImage"];


        [nc deliverNotification:note];

		int i = 0;
        while (ncDelegate.keepRunning) {
            [[NSRunLoop currentRunLoop] runUntilDate:[NSDate dateWithTimeIntervalSinceNow:0.1]];
			i++;
			if (i > 1000) {
				break;
			}
        }
    }
}
