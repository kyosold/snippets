//
//  AppDelegate.h
//  STools
//
//  Created by Grayson on 2020/8/3.
//  Copyright © 2020 Grayson. All rights reserved.
//

#import <Cocoa/Cocoa.h>


@interface AppDelegate : NSObject <NSApplicationDelegate, NSWindowDelegate>
{
    NSStatusItem *statusItem;
}

@property (weak) IBOutlet NSMenuItem *smtpMenuItem;
@property (weak) IBOutlet NSMenuItem *codecMenuItem;
@property (weak) IBOutlet NSMenuItem *cryptoMenuItem;

@property NSWindowController *smtpWC;
@property NSWindowController *codecWC;
@property NSWindowController *cryptoWC;


@property (weak) IBOutlet NSMenu *statusMenu;
@property (weak) IBOutlet NSMenuItem *sendMailStatusMenuItem;
@property (weak) IBOutlet NSMenuItem *codecStatusMenuItem;
@property (weak) IBOutlet NSMenuItem *cryptoStatusMenuItem;



- (IBAction)openSendMailAction:(id)sender;
- (IBAction)openCodecAction:(id)sender;
- (IBAction)openCryptoAction:(id)sender;
- (IBAction)exitApplicationAction:(id)sender;


@end

