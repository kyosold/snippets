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
    self.sendMailStatusMenuItem.state = NSControlStateValueOn;
    
    codecWC = [storyBoard instantiateControllerWithIdentifier:@"CodecWindowController"];
    codecWC.window.delegate = self;
    self.codecMenuItem.state = NSControlStateValueOff;
    self.codecStatusMenuItem.state = NSControlStateValueOff;
    
    cryptoWC = [storyBoard instantiateControllerWithIdentifier:@"CryptoWindowController"];
    cryptoWC.window.delegate = self;
    self.cryptoMenuItem.state = NSControlStateValueOff;
    self.cryptoStatusMenuItem.state = NSControlStateValueOff;
    

    // 添加状态栏菜单
    statusItem = [[NSStatusBar systemStatusBar] statusItemWithLength:NSVariableStatusItemLength];
    statusItem.button.image = [NSImage imageNamed:@"icon_24x24"];
    statusItem.button.toolTip = @"STools";
    statusItem.button.cell.highlighted = NO;
    statusItem.menu = _statusMenu;
}


- (void)applicationWillTerminate:(NSNotification *)aNotification {
    // Insert code here to tear down your application
}

- (BOOL)applicationShouldTerminateAfterLastWindowClosed:(NSApplication *)sender
{
    return YES;
}

- (void)windowWillClose:(NSNotification *)notification
{
    NSWindow *win = notification.object;
    if (win == smtpWC.window) {
        self.sendMailStatusMenuItem.state = NSControlStateValueOff;
        self.smtpMenuItem.state = NSControlStateValueOff;
    } else if (win == codecWC.window) {
        self.codecStatusMenuItem.state = NSControlStateValueOff;
        self.codecMenuItem.state = NSControlStateValueOff;
    } else if (win == cryptoWC.window) {
        self.cryptoStatusMenuItem.state = NSControlStateValueOff;
        self.cryptoMenuItem.state = NSControlStateValueOff;
    }
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


- (void)openSendMailWin
{
    NSLog(@"Open SendEmail Window");

    self.smtpMenuItem.state = NSControlStateValueOn;
    self.sendMailStatusMenuItem.state = NSControlStateValueOn;
    [smtpWC showWindow:self];
}

- (void)openCodecWin
{
    self.codecMenuItem.state = NSControlStateValueOn;
    self.codecStatusMenuItem.state = NSControlStateValueOn;
    [codecWC showWindow:self];
}

- (void)openCryptoWin
{
    self.cryptoMenuItem.state = NSControlStateValueOn;
    self.cryptoStatusMenuItem.state = NSControlStateValueOn;
    [cryptoWC showWindow:self];
}


// 打开发送邮件窗口
- (IBAction)openSendEmailWindow:(id)sender {
    [self openSendMailWin];
}
// 打开编解码窗口
- (IBAction)openCodecWindow:(id)sender {
    [self openCodecWin];
}
// 打开加密窗口
- (IBAction)openCryptoWindow:(id)sender {
    [self openCryptoWin];
}



- (IBAction)exitApplicationAction:(id)sender {
    [NSApplication.sharedApplication terminate:self];
}

- (IBAction)openCryptoAction:(id)sender {
    [self openCryptoWin];
}

- (IBAction)openCodecAction:(id)sender {
    [self openCodecWin];
}

- (IBAction)openSendMailAction:(id)sender {
    [self openSendMailWin];
}
@end
