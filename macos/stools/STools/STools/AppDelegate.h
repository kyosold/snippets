//
//  AppDelegate.h
//  STools
//
//  Created by Grayson on 2020/8/3.
//  Copyright © 2020 Grayson. All rights reserved.
//

#import <Cocoa/Cocoa.h>


@interface AppDelegate : NSObject <NSApplicationDelegate, NSWindowDelegate>

@property (weak) IBOutlet NSMenuItem *smtpMenuItem;
@property (weak) IBOutlet NSMenuItem *codecMenuItem;

@property NSWindowController *smtpWC;
@property NSWindowController *codecWC;

@end

