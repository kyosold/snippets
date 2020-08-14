//
//  CodeViewController.h
//  STools
//
//  Created by Grayson on 2020/8/13.
//  Copyright © 2020 Grayson. All rights reserved.
//

#import <Cocoa/Cocoa.h>

NS_ASSUME_NONNULL_BEGIN

@interface CodeViewController : NSViewController <NSUserNotificationCenterDelegate>

@property (weak) IBOutlet NSButton *exchangeButton;
@property (unsafe_unretained) IBOutlet NSTextView *srcCodeTextView;
@property (weak) IBOutlet NSPopUpButton *charsetPopUpButton;
@property (weak) IBOutlet NSPopUpButton *codecPopUpButton;


@property (weak) IBOutlet NSButton *encodeButton;
@property (unsafe_unretained) IBOutlet NSTextView *resultTextView;

- (IBAction)encodeClickAction:(id)sender;
- (IBAction)exchangeClickAction:(id)sender;
- (IBAction)decodeClickAction:(id)sender;


- (NSString *)base64EncodedString:(NSString *)str;
- (NSString *)base64DecodedString:(NSString *)str;

- (void)sendAlertWithString:(NSString *)title subTitle:(NSString *)subTitle;

@end

NS_ASSUME_NONNULL_END
