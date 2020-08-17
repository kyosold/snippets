//
//  ViewController.m
//  STools
//
//  Created by Grayson on 2020/8/3.
//  Copyright © 2020 Grayson. All rights reserved.
//

#import "ViewController.h"

//#import "TcpAppDelegate.h"

typedef enum {
    GREET = 100,
    HELO = 101,
    EHLO = 1011,
    LHLO = 1012,
    TLS = 1013,
    AUTH = 102,
    USER = 103,
    PASS = 104,
    FROM = 105,
    RCPT = 106,
    DATA = 107,
    EOM = 108,
    QUIT = 109,
    END = 110
} STATE;
STATE readState = GREET;

NSTimeInterval rwTimeout = 5.0;

int rcptIndex = 0;
NSUInteger rcptCount = 0;
NSArray *rcptList = nil;
BOOL isTLSSocket = NO;

@implementation ViewController

- (void)viewDidLoad {
    [super viewDidLoad];

    // Do any additional setup after loading the view.
    self.title = @"SMTP";
    
    self.rowData = [[NSMutableArray alloc] init];
    _dumpTableView.intercellSpacing = NSMakeSize(0, 0);
    _dumpTableView.selectionHighlightStyle = NSTableViewSelectionHighlightStyleRegular;
    
    _dumpTableView.doubleAction = @selector(doubleClickForTableViewCell:);
    
    _bodyTextView.textColor = [NSColor whiteColor];
    
    NSDictionary *iniDict = [self readConfigINI];
    NSLog(@"iniDict:\n%@", iniDict);
    NSString *ip = [iniDict objectForKey:@"ip"];
    NSString *port = [iniDict objectForKey:@"port"];
    NSString *isAuth = [iniDict objectForKey:@"isAuth"];
    NSString *user = [iniDict objectForKey:@"user"];
    NSString *password = [iniDict objectForKey:@"password"];
    NSString *envFrom = [iniDict objectForKey:@"envFrom"];
    NSString *envTo = [iniDict objectForKey:@"envTo"];
    NSString *protocol = [iniDict objectForKey:@"protocol"];
    NSString *crypto = [iniDict objectForKey:@"crypto"];
    NSString *sslPeerName = [iniDict objectForKey:@"sslPeerName"];
    NSString *isRepaceHF = [iniDict objectForKey:@"isReplaceHF"];
    NSString *replaceHF = [iniDict objectForKey:@"replaceHF"];
    NSString *body = [iniDict objectForKey:@"body"];
    
    _ipTextField.stringValue = ip ? : @"";
    _portTextField.stringValue = port ? : @"";
    if ([isAuth isEqualToString:@"yes"]) {
        _authCheckBox.state = NSControlStateValueOn;
        _userTextField.editable = YES;
        _passwordSecureTextField.editable = YES;
    } else {
        _authCheckBox.state = NSControlStateValueOff;
        _userTextField.editable = NO;
        _passwordSecureTextField.editable = NO;
    }
    _userTextField.stringValue = user ? : @"";
    _passwordSecureTextField.stringValue = password ? : @"";
    _envfromTextField.stringValue = envFrom ? : @"";
    _envtoTextField.stringValue = envTo ? : @"";
    _protocalPopUpButton.title = protocol ? : @"";
    _cryptoPopUpButton.title = crypto ? : @"";
    if ([crypto isEqualToString:@"SSL"] ||
        [crypto isEqualToString:@"TLS"]) {
        _sslPeerNameTextField.editable = YES;
        _sslPeerNameTextField.backgroundColor = [NSColor clearColor];
    } else {
        _sslPeerNameTextField.editable = NO;
        _sslPeerNameTextField.backgroundColor = [NSColor grayColor];
    }
    _sslPeerNameTextField.stringValue = sslPeerName ? : @"";
    
    if ([isRepaceHF isEqualToString:@"yes"]) {
        _replaceHeadFieldCheckBox.state = NSControlStateValueOn;
        _replaceHeadFieldPopUpButton.enabled = YES;
    } else {
        _replaceHeadFieldCheckBox.state = NSControlStateValueOff;
        _replaceHeadFieldPopUpButton.enabled = NO;
    }
    _replaceHeadFieldPopUpButton.title = replaceHF;
    _bodyTextView.string = body ? : @"";
}


- (void)setRepresentedObject:(id)representedObject {
    [super setRepresentedObject:representedObject];
    // Update the view, if already loaded.
}

// 勾选认证复选框
- (IBAction)clickAuthenticateAction:(id)sender {
    if (_authCheckBox.state) {
        _userTextField.editable = YES;
        _passwordSecureTextField.editable = YES;
        _userTextField.backgroundColor = [NSColor clearColor];
        _passwordSecureTextField.backgroundColor = [NSColor clearColor];
        [_userTextField becomeFirstResponder];
    } else {
        _userTextField.editable = NO;
        _passwordSecureTextField.editable = NO;
        _userTextField.backgroundColor = [NSColor grayColor];
        _passwordSecureTextField.backgroundColor = [NSColor grayColor];
    }
}


- (IBAction)cleanButtonAction:(id)sender {
//    _dumpTextView.string = @"";
    if (self.rowData.count <= 0)
        return;
    
    [self.rowData removeAllObjects];
    [_dumpTableView reloadData];
    // 定位光标到新添加的行
//    [self.dumpTableView editColumn:0 row:0 withEvent:nil select:YES];
}

- (IBAction)sendButtonAction:(id)sender {
    NSLog(@"%@", _sendButton.title);
    
    if ([_sendButton.title isEqual: @"Disconnect"]) {
        NSLog(@"Need call disconnect socket");
        [self.asyncSocket disconnect];
        self.asyncSocket = nil;
        return;
    }
    
    [self setSendButtonStatus:NO];
    
    NSString *ip = _ipTextField.stringValue;
    ip = [ip stringByTrimmingCharactersInSet:[NSCharacterSet whitespaceCharacterSet]];
    
    NSString *port = _portTextField.stringValue;
    port = [port stringByTrimmingCharactersInSet:[NSCharacterSet whitespaceCharacterSet]];
    
    NSString *user = @"";
    NSString *password = @"";
    BOOL isAuth = _authCheckBox.state;
    if (isAuth) {
        user = _userTextField.stringValue;
        user = [user stringByTrimmingCharactersInSet:[NSCharacterSet whitespaceCharacterSet]];
        password = _passwordSecureTextField.stringValue;
        password = [password stringByTrimmingCharactersInSet:[NSCharacterSet whitespaceCharacterSet]];
    }
    
    NSString *envFrom = _envfromTextField.stringValue;
    envFrom = [envFrom stringByTrimmingCharactersInSet:[NSCharacterSet whitespaceCharacterSet]];
    
    NSString *envto = _envtoTextField.stringValue;
    envto = [envto stringByTrimmingCharactersInSet:[NSCharacterSet whitespaceCharacterSet]];
    
    NSString *sslPeerName = _sslPeerNameTextField.stringValue;
    sslPeerName = [sslPeerName stringByTrimmingCharactersInSet:[NSCharacterSet whitespaceCharacterSet]];
    
    NSString *isReplaceHF = @"no";
    if (_replaceHeadFieldCheckBox.state == YES) {
        isReplaceHF = @"yes";
    }
    NSString *replaceHF = _replaceHeadFieldPopUpButton.titleOfSelectedItem;
    
    NSString *body = _bodyTextView.string;
    
    NSString *protocol = _protocalPopUpButton.titleOfSelectedItem;
    NSString *crypto = _cryptoPopUpButton.titleOfSelectedItem;
    
    if (ip.length <= 0 ||
        port.length <= 0 ||
        envFrom.length <= 0 ||
        envto.length <= 0 ||
        body.length <= 0) {
        NSLog(@"Argument Error");
        [self dumpTableViewAppendString:@"Argument Error" withType:@"ME"];
        [self setSendButtonStatus:YES];
        return;
    }
    if (isAuth) {
        if (user.length <= 0 ||
            password.length <= 0) {
            NSLog(@"SASL Argument Error");
            [self dumpTableViewAppendString:@"Argument Error" withType:@"ME"];
            [self setSendButtonStatus:YES];
            return;
        }
    }
    
    rcptList = [_envtoTextField.stringValue componentsSeparatedByString:@","];
    rcptCount = rcptList.count;
    rcptIndex = 0;
        
    NSMutableDictionary *args = [NSMutableDictionary dictionary];
    [args setValue:ip forKey:@"ip"];
    [args setValue:port forKey:@"port"];
    [args setValue:(isAuth ? @"yes" : @"no") forKey:@"isAuth"];
    [args setValue:user forKey:@"user"];
    [args setValue:password forKey:@"password"];
    [args setValue:envFrom forKey:@"envFrom"];
    [args setValue:envto forKey:@"envTo"];
    [args setValue:protocol forKey:@"protocol"];
    [args setValue:crypto forKey:@"crypto"];
    [args setValue:sslPeerName forKey:@"sslPeerName"];
    [args setValue:isReplaceHF forKey:@"isReplaceHF"];
    [args setValue:replaceHF forKey:@"replaceHF"];
    [args setValue:body forKey:@"body"];
    
    [self saveConfigINIWithDict:[NSDictionary dictionaryWithDictionary:args]];
    
    // Connect to Remote
    if (self.asyncSocket == nil) {
        self.asyncSocket = [[GCDAsyncSocket alloc] initWithDelegate:self delegateQueue:dispatch_get_main_queue()];

    }
    if (!self.asyncSocket.isConnected) {
        [self dumpTableViewAppendString:[NSString stringWithFormat:@"Connecting to %@:%d ...", ip, [port intValue]] withType:@"ME"];
        
        NSError *error;
        BOOL ret = [self.asyncSocket connectToHost:ip onPort:[port intValue] withTimeout:rwTimeout error:&error];
        if (!ret) {
            NSLog(@"Error Connect: %@", error);
            [self setSendButtonStatus:YES];
            return;
        }
    }
    
}


// status:
// YES: 设置button为 Send
// NO:  设置button为 Disconnect
- (void)setSendButtonStatus:(BOOL)state
{
    if (state) {
        _sendButton.title = @"Send";
//        _dumpTextView.editable = NO;
//        _dumpTextView.backgroundColor = [NSColor grayColor];
    } else {
        _sendButton.title = @"Disconnect";
//        _dumpTextView.editable = YES;
//        _dumpTextView.backgroundColor = [NSColor clearColor];
    }
}


#pragma Socket Delegate
// 连接断开
- (void)socketDidDisconnect:(GCDAsyncSocket *)sock withError:(NSError *)err
{
    NSLog(@"SocketDidDisconnect:%p withError:%@", sock, err);
    if (err != nil) {
        [self dumpTableViewAppendString:[NSString stringWithFormat:@"Disconnect Connection. \n\nError: (%@)", err] withType:@"ME"];
    } else {
        [self dumpTableViewAppendString:@"Disconnect Connection." withType:@"ME"];
    }

    
    [self setSendButtonStatus:YES];
    
    rcptCount = 0;
    rcptIndex = 0;
    rcptList = nil;
    isTLSSocket = NO;
}

// 连接成功
- (void)socket:(GCDAsyncSocket *)sock didConnectToHost:(NSString *)host port:(uint16_t)port
{
    NSLog(@"Socket:%p didConnectToHost:%@ port:%hu", sock, host, port);
    
    // 输出
    [self dumpTableViewAppendString:[NSString stringWithFormat:@"-------------------------------\nConnect Remote(%@:%hu) Success...", host, port] withType:@"ME"];
    
    NSString *crypto = _cryptoPopUpButton.titleOfSelectedItem;
    NSString *sslPeerName = [_sslPeerNameTextField.stringValue stringByTrimmingCharactersInSet:[NSCharacterSet whitespaceCharacterSet]];
    if (![crypto isEqualToString:@"No Crypto"]) {
        if (sslPeerName.length <= 0) {
            [self dumpTableViewAppendString:@"Use SSL/TLS Must be set 'SSLPeerName'" withType:@"ME"];

            [self dieConnect];
            return;
        }
        
        if ([crypto isEqualToString:@"SSL"]) {
            [self setSocketSSL:YES withPeerName:sslPeerName];
        }
        
    }
    
    // 连接成功或收到消息，必须开始 read, 否则将无法收到消息,
    // 不read的话，缓存区将会被关闭
    // -1 表示无限时长, 永久不失效
    readState = GREET;
    [self.asyncSocket readDataWithTimeout:rwTimeout tag:readState];
}

- (void)socket:(GCDAsyncSocket *)sock didReceiveTrust:(SecTrustRef)trust completionHandler:(void (^)(BOOL))completionHandler
{
    NSLog(@"Socket: shouldTrustPeer:");

//    dispatch_queue_t bgQueue = dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0);
//    dispatch_async(bgQueue, ^{
        CFErrorRef errorRef;
        bool ret = SecTrustEvaluateWithError(trust, &errorRef);
        if (ret) {
            NSLog(@"Certificate Match");
            [self dumpTableViewAppendString:@"+++[Certificate Match]+++" withType:@"ME"];
            completionHandler(YES);
        } else {
            NSLog(@"Certificate Not Match");
            [self dumpTableViewAppendString:@"+++[Certificate Not Match]+++" withType:@"ERROR"];
            completionHandler(NO);
        }
        
//    });
    
    // --- 获取证书并验证 START ----
    // Server Certificate
//    SecCertificateRef serverCertificate = SecTrustGetCertificateAtIndex(trust, 0);
//    CFDataRef serverCertificateData = SecCertificateCopyData(serverCertificate);
//
//    const UInt8 *serverData = CFDataGetBytePtr(serverCertificateData);
//    const CFIndex serverDataSize = CFDataGetLength(serverCertificateData);
//    NSData *cert1 = [NSData dataWithBytes:serverData length:(NSUInteger)serverDataSize];
//
//    // Local certificate
//    NSString *localCertFilePath = [[NSBundle mainBundle] pathForResource:@"LocalCertificate" ofType:@"cer"];
//    NSData *localCertData = [NSData dataWithContentsOfFile:localCertFilePath];
//    CFDataRef myCertData = (__bridge CFDataRef)localCertData;
//
//    const UInt8 *localData = CFDataGetBytePtr(myCertData);
//    const CFIndex localDataSize = CFDataGetLength(myCertData);
//    NSData *cert2 = [NSData dataWithBytes:localData length:(NSUInteger)localDataSize];
//
//    if (cert1 == nil || cert2 == nil) {
//        NSLog(@"Certificate NULL");
//        return;
//    }
//
//    const BOOL equal = [cert1 isEqualToData:cert2];
//    if (equal) {
//        NSLog(@"Certificate Match");
//        completionHandler(YES);
//    } else {
//        NSLog(@"Certificate Not Match");
//        completionHandler(NO);
//    }
    // --- 获取证书并验证 END ----
}

- (void)socketDidSecure:(GCDAsyncSocket *)sock
{
    NSLog(@"socketDidSecure:");
}




// 发送数据到服务器成功
- (void)socket:(GCDAsyncSocket *)sock didWriteDataWithTag:(long)tag
{
    NSLog(@"%ld Send Data Succ.", tag);
}

// 接收服务器返回数据
- (void)socket:(GCDAsyncSocket *)sock didReadData:(NSData *)data withTag:(long)tag
{
    NSLog(@"Read: tag(%ld) length(%ld)", tag, data.length);
    
    NSString *inStr = [[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding];
    NSLog(@"%@", inStr);
    
    [self dumpTableViewAppendString:inStr withType:@"OTHER"];
    
    
    // ---- Parse Response ----
    BOOL supportTLS = NO;
    BOOL supportAuth = NO;
    NSString *code = @"";
    
    if (readState == GREET ||
        readState == HELO ||
        readState == TLS ||
        readState == AUTH ||
        readState == USER ||
        readState == PASS ||
        readState == FROM ||
        readState == RCPT ||
        readState == DATA ||
        readState == QUIT ||
        readState == END) {
        code = [inStr substringToIndex:3];
        
    } else if (readState == EHLO ||
               readState == LHLO) {
        NSArray *codeLineList = [[inStr stringByTrimmingCharactersInSet:[NSCharacterSet whitespaceAndNewlineCharacterSet]]  componentsSeparatedByString:@"\n"];
        NSLog(@"codeLineList: %@", codeLineList);
        for (int i=0; i<codeLineList.count; i++) {
            if ([codeLineList[i] localizedCaseInsensitiveContainsString:@"-STARTTLS"]) {
                supportTLS = YES;
            } else if ([codeLineList[i] localizedCaseInsensitiveContainsString:@"-AUTH=LOGIN"]) {
                supportAuth = YES;
            }
            code = [inStr substringToIndex:3];
        }
        
    }
    
    
    if (readState == GREET ||
        readState == TLS) {
        if (![code isEqualToString:@"220"]) {
            [self dieConnect];
            return;
        }
    } else if (readState == END) {
        if (![code isEqualToString:@"221"]) {
            [self dieConnect];
            return;
        }
    } else if (readState == HELO ||
               readState == EHLO ||
               readState == LHLO ||
               readState == FROM ||
               readState == RCPT ||
               readState == QUIT) {
        if (![code isEqualToString:@"250"]) {
            [self dieConnect];
            return;
        }
    } else if (readState == AUTH ||
               readState == USER) {
        if (![code isEqualToString:@"334"]) {
            [self dieConnect];
            return;
        }
    } else if (readState == PASS) {
        if (![code isEqualToString:@"235"]) {
            [self dieConnect];
            return;
        }
    } else if (readState == DATA) {
        if (![code isEqualToString:@"354"]) {
            [self dieConnect];
            return;
        }
    }
    
    // ---- SEND ----
    NSString *outStr = nil;
    if (tag == GREET) {
        NSString *protocol = _protocalPopUpButton.titleOfSelectedItem;
        BOOL isAuth = _authCheckBox.state;
        if ([protocol isEqualToString:@"SMTP"]) {
            if (isAuth) {
                outStr = [NSString stringWithFormat:@"EHLO %@\r\n", [[NSHost currentHost] name]];
                readState = EHLO;
            } else {
                outStr = [NSString stringWithFormat:@"HELO %@\r\n", [[NSHost currentHost] name]];
                readState = HELO;
            }
            
        } else if ([protocol isEqualToString:@"LMTP"]) {
            outStr = [NSString stringWithFormat:@"LHLO %@\r\n", [[NSHost currentHost] name]];
            readState = LHLO;
        } else {
            [self dumpTableViewAppendString:@"Unknow Protocol" withType:@"ME"];
            [self.asyncSocket disconnect];
            self.asyncSocket = nil;
            return;
        }
        
    } else if (tag == HELO) {
        outStr = [NSString stringWithFormat:@"MAIL FROM:<%@>\r\n", [_envfromTextField.stringValue stringByTrimmingCharactersInSet:[NSCharacterSet whitespaceCharacterSet]]];
        readState = FROM;
        
    } else if (tag == EHLO) {
        if ((isTLSSocket == NO) && [_cryptoPopUpButton.titleOfSelectedItem isEqualToString:@"TLS"]) {
            outStr = @"STARTTLS\r\n";
            readState = TLS;
            
        } else {
            if (supportAuth) {
                outStr = @"AUTH LOGIN\r\n";
                readState = AUTH;
            } else {
                outStr = [NSString stringWithFormat:@"MAIL FROM:<%@>\r\n", [_envfromTextField.stringValue stringByTrimmingCharactersInSet:[NSCharacterSet whitespaceCharacterSet]]];
                readState = FROM;
            }
        }
        
    } else if (tag == TLS) {
        NSString *sslPeerName = [_sslPeerNameTextField.stringValue stringByTrimmingCharactersInSet:[NSCharacterSet whitespaceCharacterSet]];
        
        [self setSocketSSL:YES withPeerName:sslPeerName];
        
        isTLSSocket = YES;
        
        NSString *protocol = _protocalPopUpButton.titleOfSelectedItem;
        if ([protocol isEqualToString:@"SMTP"]) {
            outStr = [NSString stringWithFormat:@"EHLO %@\r\n", [[NSHost currentHost] name]];
            readState = EHLO;
        } else {
            outStr = [NSString stringWithFormat:@"LHLO %@\r\n", [[NSHost currentHost] name]];
            readState = LHLO;
        }
    
    } else if (tag == LHLO) {
        outStr = [NSString stringWithFormat:@"MAIL FROM:<%@>\r\n", [_envfromTextField.stringValue stringByTrimmingCharactersInSet:[NSCharacterSet whitespaceCharacterSet]]];
        readState = FROM;
        
    } else if (tag == AUTH) {
        NSString *user = [_userTextField.stringValue stringByTrimmingCharactersInSet:[NSCharacterSet whitespaceCharacterSet]];
        NSString *userB64 = [self base64EncodedString:user];
        outStr = [NSString stringWithFormat:@"%@\r\n", userB64];
        readState = USER;
    
    } else if (tag == USER) {
        NSString *pass = [_passwordSecureTextField.stringValue stringByTrimmingCharactersInSet:[NSCharacterSet whitespaceCharacterSet]];
        NSString *passB64 = [self base64EncodedString:pass];
        outStr = [NSString stringWithFormat:@"%@\r\n", passB64];
        readState = PASS;
    
    } else if (tag == PASS) {
        outStr = [NSString stringWithFormat:@"MAIL FROM:<%@>\r\n", [_envfromTextField.stringValue stringByTrimmingCharactersInSet:[NSCharacterSet whitespaceCharacterSet]]];
        readState = FROM;
        
    } else if (tag == FROM) {
        NSString *to = [rcptList[rcptIndex] stringByTrimmingCharactersInSet:[NSCharacterSet whitespaceCharacterSet]];
        outStr = [NSString stringWithFormat:@"RCPT TO:<%@>\r\n", to];
        
        if ((rcptIndex + 1) == rcptCount) {
            readState = RCPT;
        } else {
            readState = FROM;
        }
        rcptIndex++;
    
    } else if (tag == RCPT) {
        outStr = @"DATA\r\n";
        readState = DATA;
        
    } else if (tag == DATA) {
        BOOL isHead = YES;
        NSArray *bodyList = [_bodyTextView.string componentsSeparatedByString:@"\n"];
        NSMutableString *newBodyString = [[NSMutableString alloc] init];
        for (int i=0; i<bodyList.count; i++) {
            NSLog(@"%@", bodyList[i]);
            if ([bodyList[i] length] <= 0) {
                [newBodyString appendString:@"\n"];
                continue;
            }
            
            if (isHead) {
                if (_replaceHeadFieldCheckBox.state == NSControlStateValueOn) {
                    NSArray *headItem = [bodyList[i] componentsSeparatedByString:@":"];
                    if (headItem.count > 1) {
                        if ([headItem[0] caseInsensitiveCompare:_replaceHeadFieldPopUpButton.titleOfSelectedItem] == NSOrderedSame) {
                            if ([_replaceHeadFieldPopUpButton.titleOfSelectedItem caseInsensitiveCompare:@"message-id"] == NSOrderedSame) {
                                NSString *messageid = [self generateTradeNO];
                                [newBodyString appendFormat:@"Message-ID: %@\n", messageid];
                                isHead = NO;
                                continue;
                            } else if ([_replaceHeadFieldPopUpButton.titleOfSelectedItem caseInsensitiveCompare:@"Date"] == NSOrderedSame) {
                                NSDate *now = [NSDate date];
                                NSDateFormatter *dateFmt = [[NSDateFormatter alloc] init];
                                [dateFmt setDateFormat:@"EEE, d MMM yyyy HH:mm:ss ZZZ"];
                                NSString *nowDate = [dateFmt stringFromDate:now];
                                [newBodyString appendFormat:@"Date: %@\n", nowDate];
                                isHead = NO;
                                continue;
                            }
                        }
                    }
                }
            }
            
            unichar ch = [bodyList[i] characterAtIndex:0];
            NSLog(@"ch=(%c)", ch);
            if (ch == '.') {
                [newBodyString appendFormat:@".%@\n", bodyList[i]];
            } else {
                if (isHead && [bodyList[i] length] == 1 && ch == '\r') {
                    isHead = NO;
                }
                [newBodyString appendFormat:@"%@\n", bodyList[i]];
            }
        }
        outStr = [NSString stringWithFormat:@"%@\r\n.\r\n", newBodyString];
        readState = QUIT;
        
    } else if (tag == QUIT) {
        outStr = @"QUIT\r\n";
        readState = END;
        
    } else if (tag == END) {
        [self dieConnect];
        return;
    }
    
    [self.asyncSocket writeData:[outStr dataUsingEncoding:NSUTF8StringEncoding] withTimeout:5.0 tag:readState];
    
//    _dumpTextView.string = [NSString stringWithFormat:@"%@%@", _dumpTextView.string, outStr];
    [self dumpTableViewAppendString:outStr withType:@"ME"];
    
    
    [self.asyncSocket readDataWithTimeout:rwTimeout tag:readState];
}


- (void)dieConnect
{
    [self.asyncSocket disconnect];
    self.asyncSocket = nil;
}

- (NSString *)base64EncodedString:(NSString *)str
{
    NSData *data = [str dataUsingEncoding:NSUTF8StringEncoding];
    return [data base64EncodedStringWithOptions:0];
}

- (NSString *)base64DecodedString:(NSString *)str
{
    NSData *data = [[NSData alloc] initWithBase64EncodedString:str options:0];
    return [[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding];
}

- (void)saveConfigINIWithDict:(NSDictionary *)dict
{
    // 获取ini路径
    NSString *iniPath = NSSearchPathForDirectoriesInDomains(NSDocumentDirectory, NSUserDomainMask, YES)[0];
    
    // 拼接plist路径
    NSString *filePath = [iniPath stringByAppendingPathComponent:@"stools.plist"];
    
    // 存储字典
    [dict writeToFile:filePath atomically:YES];
}
- (NSDictionary *)readConfigINI
{
    NSString *iniPath = NSSearchPathForDirectoriesInDomains(NSDocumentDirectory, NSUserDomainMask, YES)[0];
    NSString *filePath = [iniPath stringByAppendingPathComponent:@"stools.plist"];
    NSDictionary *dict = [NSDictionary dictionaryWithContentsOfFile:filePath];
    NSLog(@"Read Ini %@:\n%@", filePath, dict);
    return dict;
}

- (IBAction)clickCryptoPopUpButton:(id)sender {
    NSLog(@"Crypto Selected value is: %@", [(NSPopUpButton *)sender titleOfSelectedItem]);
    NSString *val = [(NSPopUpButton *)sender titleOfSelectedItem];
    if ([val isEqualToString:@"No Crypto"]) {
        _sslPeerNameTextField.editable = NO;
        _sslPeerNameTextField.backgroundColor = [NSColor grayColor];
    } else {
        _sslPeerNameTextField.editable = YES;
        _sslPeerNameTextField.backgroundColor = [NSColor clearColor];
    }
}

- (IBAction)clickReplaceHeadFieldAction:(id)sender {
    if (_replaceHeadFieldCheckBox.state == NSControlStateValueOn) {
        _replaceHeadFieldPopUpButton.enabled = YES;
    } else {
        _replaceHeadFieldPopUpButton.enabled = NO;
    }
}

- (void)setSocketSSL:(BOOL)state withPeerName:(NSString *)peerName
{
    if (!state) {
        [self.asyncSocket startTLS:nil];
        return;
    }
    
    if (peerName.length <= 0) {
        [self dumpTableViewAppendString:@"se SSL/TLS Must be set 'SSLPeerName'" withType:@"ME"];
        [self dieConnect];
        return;
    }
    

    NSMutableDictionary *opts = [NSMutableDictionary dictionaryWithCapacity:3];
    [opts setObject:[NSNumber numberWithBool:YES] forKey:GCDAsyncSocketManuallyEvaluateTrust];
    [opts setObject:[NSString stringWithString:peerName] forKey:GCDAsyncSocketSSLPeerName];


    [self.asyncSocket startTLS:opts];
}



- (void)dumpTableViewAppendString:(NSString *)str withType:(NSString *)type
{
    NSArray *strList = [str componentsSeparatedByString:@"\n"];
    for (NSInteger i=0; i<strList.count; i++) {
        if ([strList[i] length] <= 0) continue;
        if ([strList[i] characterAtIndex:0] == '\r') continue;
        
        NSMutableDictionary *msg = [[NSMutableDictionary alloc] init];
        [msg setValue:strList[i] forKey:@"text"];
        [msg setValue:type forKey:@"type"];
        [self.rowData addObject:msg];
    }

    [self.dumpTableView reloadData];
    
    NSInteger numberOfRows = self.dumpTableView.numberOfRows;
    if (numberOfRows > 0)
        [self.dumpTableView scrollRowToVisible:(numberOfRows - 1)];
    
    // 定位光标到新添加的行
    [self.dumpTableView editColumn:0 row:(self.rowData.count - 1) withEvent:nil select:YES];
}

#pragma NSTableView Delegate

// 返回数据行数
- (NSInteger)numberOfRowsInTableView:(NSTableView *)tableView
{
    return self.rowData.count;
}

// 返回row对应的自定义视图
- (NSView *)tableView:(NSTableView *)tableView viewForTableColumn:(NSTableColumn *)tableColumn row:(NSInteger)row
{
    
    NSDictionary *dict = [self.rowData objectAtIndex:row];
    NSString *type = [dict objectForKey:@"type"];
    NSString *text = [dict objectForKey:@"text"];
    
    if (!text) {
        return nil;
    }
    
    NSTableCellView *cell = [tableView makeViewWithIdentifier:@"dump" owner:self];
    cell.textField.stringValue = text;
    cell.textField.drawsBackground = YES;
    
    if ([type isEqualToString:@"ME"]) {
        cell.textField.textColor = [self getColorFromRGB:93 green:155 blue:194];
    } else if ([type isEqualToString:@"ERROR"]) {
        cell.textField.textColor = [self getColorFromRGB:249 green:141 blue:57];
    } else {
        cell.textField.textColor = [self getColorFromRGB:133 green:90 blue:194];
    }
    return cell;
}

// 点击选中row事件
- (void)tableViewSelectionDidChange:(NSNotification *)notification
{
    NSLog(@"Selected Row:%ld", (long)_dumpTableView.selectedRow);
    [self.dumpTableView editColumn:0 row:self.dumpTableView.selectedRow withEvent:nil select:YES];
}


- (NSColor *)getColorFromRGB:(unsigned char)r green:(unsigned char)g blue:(unsigned char)b
{
    CGFloat rFloat = r/255.0;
    CGFloat gFloat = g/255.0;
    CGFloat bFloat = b/255.0;
    
    return [NSColor colorWithCalibratedRed:rFloat green:gFloat blue:bFloat alpha:1.0];
}


- (void)doubleClickForTableViewCell:(id)sender
{
    NSInteger rowIndex = _dumpTableView.clickedRow;
    NSLog(@"Double Click Row:%ld", rowIndex);
    if (rowIndex > self.rowData.count || rowIndex < 0)
        return;
    NSDictionary *dict = [self.rowData objectAtIndex:rowIndex];
    NSString *type = [dict objectForKey:@"type"];
    NSString *text = [dict objectForKey:@"text"];
    
    // NSPasteBoard 使用 ------
    // --- 放内容 ----
    NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
    [pasteboard clearContents];
    [pasteboard setString:text forType:NSPasteboardTypeString];
    
    // --- 取内容 ----
//    NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
//    NSArray *types = [pasteboard types];
//    if ([types containsObject:NSPasteboardTypeString]) {
//        NSString *s = [pasteboard stringForType:NSPasteboardTypeString];
//        do something...
//    }
    

    NSUserNotification *noti = [[NSUserNotification alloc] init];
    noti.title = @"Copy Finished";
//    noti.subtitle = @"小标题";
//    noti.informativeText = @"详细文字说明";
    noti.hasActionButton = YES;
    noti.actionButtonTitle = @"OK";
    noti.otherButtonTitle = @"Cancel";
    
    [[NSUserNotificationCenter defaultUserNotificationCenter] scheduleNotification:noti];
    [[NSUserNotificationCenter defaultUserNotificationCenter] setDelegate:self];
    [NSTimer scheduledTimerWithTimeInterval:2.0 repeats:NO block:^(NSTimer * _Nonnull timer) {
        [[NSUserNotificationCenter defaultUserNotificationCenter] removeDeliveredNotification:noti];
    }];
    
    
    NSLog(@"Double Click Row Index:%ld Type:%@ Text:%@", rowIndex, type, text);
}

- (void)userNotificationCenter:(NSUserNotificationCenter *)center didDeliverNotification:(NSUserNotification *)notification
{
    NSLog(@"通知已经递交!");
}
- (void)userNotificationCenter:(NSUserNotificationCenter *)center didActivateNotification:(NSUserNotification *)notification
{
    NSLog(@"用户点击了通知!");
//    [[NSUserNotificationCenter defaultUserNotificationCenter] removeDeliveredNotification:noti];
    [center removeDeliveredNotification:notification];
}
- (BOOL)userNotificationCenter:(NSUserNotificationCenter *)center shouldPresentNotification:(NSUserNotification *)notification
{
    return YES;
}


//生成随机数算法 ,随机字符串，不长于32位
//微信支付API接口协议中包含字段nonce_str，主要保证签名不可预测。
//我们推荐生成随机数算法如下：调用随机数函数生成，将得到的值转换为字符串。
- (NSString *)generateTradeNO
{
    static int kNumber = 15;
    NSString *sourceStr = @"0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ";
    NSMutableString *resultStr = [[NSMutableString alloc] init];
    
    //  srand函数是初始化随机数的种子，为接下来的rand函数调用做准备。
    //  time(0)函数返回某一特定时间的小数值。
    //  这条语句的意思就是初始化随机数种子，time函数是为了提高随机的质量（也就是减少重复）而使用的。
    
    //　srand(time(0)) 就是给这个算法一个启动种子，也就是算法的随机种子数，有这个数以后才可以产生随机数,用1970.1.1至今的秒数，初始化随机数种子。
    //　Srand是种下随机种子数，你每回种下的种子不一样，用Rand得到的随机数就不一样。为了每回种下一个不一样的种子，所以就选用Time(0)，Time(0)是得到当前时时间值（因为每时每刻时间是不一样的了）。
    srand((unsigned)time(0));
    for (int i=0; i < kNumber; i++) {
        unsigned index = rand() % [sourceStr length];
        NSString *oneStr = [sourceStr substringWithRange:NSMakeRange(index, 1)];
        [resultStr appendString:oneStr];
    }
    return resultStr;;
}


@end
