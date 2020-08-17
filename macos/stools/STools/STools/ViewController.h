//
//  ViewController.h
//  STools
//
//  Created by Grayson on 2020/8/3.
//  Copyright © 2020 Grayson. All rights reserved.
//

#import <Cocoa/Cocoa.h>
#import "GCDAsyncSocket.h"


@interface ViewController : NSViewController <GCDAsyncSocketDelegate, NSTableViewDataSource, NSTableViewDelegate, NSUserNotificationCenterDelegate>
@property (nonatomic, strong) GCDAsyncSocket *asyncSocket;


@property (weak) IBOutlet NSTextField *ipTextField;
@property (weak) IBOutlet NSTextField *portTextField;

@property (weak) IBOutlet NSButton *authCheckBox;
@property (weak) IBOutlet NSTextField *userTextField;
@property (weak) IBOutlet NSSecureTextField *passwordSecureTextField;

@property (weak) IBOutlet NSTextField *envfromTextField;
@property (weak) IBOutlet NSTextField *envtoTextField;
@property (weak) IBOutlet NSScrollView *bodyScrollView;
@property (unsafe_unretained) IBOutlet NSTextView *bodyTextView;

@property (weak) IBOutlet NSButton *replaceHeadFieldCheckBox;

@property (weak) IBOutlet NSPopUpButton *protocalPopUpButton;
@property (weak) IBOutlet NSPopUpButton *cryptoPopUpButton;
@property (weak) IBOutlet NSTextField *sslPeerNameTextField;

@property (weak) IBOutlet NSPopUpButton *replaceHeadFieldPopUpButton;
- (IBAction)clickReplaceHeadFieldAction:(id)sender;


- (IBAction)clickCryptoPopUpButton:(id)sender;

@property (weak) IBOutlet NSButton *sendButton;

@property NSMutableArray *rowData;
@property (weak) IBOutlet NSTableView *dumpTableView;

- (void)dieConnect;

- (NSString *)base64EncodedString:(NSString *)str;
- (NSString *)base64DecodedString:(NSString *)str;

// status:
// YES: 设置button为 Send
// NO:  设置button为 Disconnect
- (void)setSendButtonStatus:(BOOL)state;

// 保存数据
- (void)saveConfigINIWithDict:(NSDictionary *)dict;
// 读取数据
- (NSDictionary *)readConfigINI;

- (void)setSocketSSL:(BOOL)state withPeerName:(NSString *)peerName;



// type:
//  "ME", "OTHER"
- (void)dumpTableViewAppendString:(NSString *)str withType:(NSString *)type;

- (NSColor *)getColorFromRGB:(unsigned char)r green:(unsigned char)g blue:(unsigned char)b;


-(void)doubleClickForTableViewCell:(id)sender;

- (NSString *)generateTradeNO;

@end

