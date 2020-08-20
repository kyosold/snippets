//
//  AppDelegate.m
//  STools
//
//  Created by Grayson on 2020/8/3.
//  Copyright © 2020 Grayson. All rights reserved.
//

#import "AppDelegate.h"

@interface AppDelegate ()

@end

@implementation AppDelegate

@synthesize smtpWC;
@synthesize codecWC;
@synthesize cryptoWC;

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification {
    // Insert code here to initialize your application
    
    // get a reference to the storyboard
    NSStoryboard *storyBoard = [NSStoryboard storyboardWithName:@"Main" bundle:nil];
    
    // instantiate your window controller
    smtpWC = [storyBoard instantiateControllerWithIdentifier:@"SmtpWindowController"];
    smtpWC.window.delegate = self;
    [smtpWC showWindow:self];
    self.smtpMenuItem.state = NSControlStateValueOn;
    
    codecWC = [storyBoard instantiateControllerWithIdentifier:@"CodecWindowController"];
    codecWC.window.delegate = self;
    self.codecMenuItem.state = NSControlStateValueOff;
    
    cryptoWC = [storyBoard instantiateControllerWithIdentifier:@"CryptoWindowController"];
    cryptoWC.window.delegate = self;
    self.cryptoMenuItem.state = NSControlStateValueOff;
    

}


- (void)applicationWillTerminate:(NSNotification *)aNotification {
    // Insert code here to tear down your application
}

- (BOOL)applicationShouldTerminateAfterLastWindowClosed:(NSApplication *)sender
{
    return YES;
}


- (BOOL)windowShouldClose:(NSWindow *)sender
{
    if ([sender isEqual:smtpWC.window]) {
        NSLog(@"Close SMTP Window");
        self.smtpMenuItem.state = NSControlStateValueOff;
    } else if ([sender isEqual:codecWC.window]) {
        NSLog(@"Close Codec Window");
        self.codecMenuItem.state = NSControlStateValueOff;
    } else if ([sender isEqual:cryptoWC.window]) {
        NSLog(@"Close Crypto Window");
        self.cryptoMenuItem.state = NSControlStateValueOff;
    } else {
        NSLog(@"Close Window");
    }
    NSLog(@"Close Window");
    return YES;
}



// 打开发送邮件窗口
- (IBAction)openSendEmailWindow:(id)sender {
    NSLog(@"Open SendEmail Window");

    self.smtpMenuItem.state = NSControlStateValueOn;
    [smtpWC showWindow:self];
}
// 打开编解码窗口
- (IBAction)openCodecWindow:(id)sender {
    self.codecMenuItem.state = NSControlStateValueOn;
    [codecWC showWindow:self];
}
// 打开加密窗口
- (IBAction)openCryptoWindow:(id)sender {
    self.cryptoMenuItem.state = NSControlStateValueOn;
    [cryptoWC showWindow:self];
}



@end
