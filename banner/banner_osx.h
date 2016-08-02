// Based on https://github.com/norio-nomura/usernotification (WTFPL, 2013)

#import <Foundation/Foundation.h>
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

void BannerNotify(const char *title, const char *message, const char *sound)
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
        note.informativeText = [NSString stringWithUTF8String:message];

        if ([[NSString stringWithUTF8String:sound] isEqualToString:@"_mute"] == NO) {
            note.soundName = [NSString stringWithUTF8String:sound];
        }

        [nc deliverNotification:note];

        while (ncDelegate.keepRunning) {
            [[NSRunLoop currentRunLoop] runUntilDate:[NSDate dateWithTimeIntervalSinceNow:0.1]];
        }
    }
}
