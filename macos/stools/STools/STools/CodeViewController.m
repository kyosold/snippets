//
//  CodeViewController.m
//  STools
//
//  Created by Grayson on 2020/8/13.
//  Copyright © 2020 Grayson. All rights reserved.
//

#import "CodeViewController.h"

@interface CodeViewController ()

@end

@implementation CodeViewController

- (void)viewDidLoad {
    [super viewDidLoad];
    // Do view setup here.
    
    self.srcCodeTextView.textColor = [NSColor whiteColor];
    self.resultTextView.textColor = [NSColor whiteColor];
}

- (IBAction)decodeClickAction:(id)sender {
    NSString *srcString = self.srcCodeTextView.string;
    
    if ([self.codecPopUpButton.titleOfSelectedItem isEqualToString:@"Base64"]) {
        
        if ([self.charsetPopUpButton.titleOfSelectedItem isEqualToString:@"UTF-8"]) {
            NSData *data = [[NSData alloc] initWithBase64EncodedString:srcString options:0];
            NSString *result = [[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding];
            if (result == nil) {
                result = @"";
                [self sendAlertWithString:@"Decode Fail" subTitle:@"Try Change Charset Again"];
            }
            self.resultTextView.string = result;
            
        } else if ([self.charsetPopUpButton.titleOfSelectedItem isEqualToString:@"GBK"]) {
            NSStringEncoding gbkEncoding = CFStringConvertEncodingToNSStringEncoding(kCFStringEncodingGB_18030_2000);
            NSData *data = [[NSData alloc] initWithBase64EncodedString:srcString options:0];
            NSString *result = [[NSString alloc] initWithData:data encoding:gbkEncoding];
            if (result == nil) {
                result = @"";
                [self sendAlertWithString:@"Decode Fail" subTitle:@"Try Change Charset Again"];
            }
            self.resultTextView.string = result;
        }
        
    } else if ([self.codecPopUpButton.titleOfSelectedItem isEqualToString:@"URL"]) {
        
        NSString *result;
        
        if ([self.charsetPopUpButton.titleOfSelectedItem isEqualToString:@"UTF-8"]) {
//            self.resultTextView.string = [srcString stringByReplacingPercentEscapesUsingEncoding:NSUTF8StringEncoding];
            result = [srcString stringByRemovingPercentEncoding];
        
        } else if ([self.charsetPopUpButton.titleOfSelectedItem isEqualToString:@"GBK"]) {
            NSStringEncoding gbkEncoding = CFStringConvertEncodingToNSStringEncoding(kCFStringEncodingGB_18030_2000);
            result = [srcString stringByReplacingPercentEscapesUsingEncoding:gbkEncoding];
            
        }
        
        if (result == nil) {
            result = @"";
            [self sendAlertWithString:@"Decode Fail" subTitle:@"Try Change Charset Again"];
        }
        self.resultTextView.string = result;
    
    }
}

- (IBAction)exchangeClickAction:(id)sender {
    NSString *srcString = [NSString stringWithString:self.srcCodeTextView.string];
    self.srcCodeTextView.string = self.resultTextView.string;
    self.resultTextView.string = srcString;
}

- (IBAction)encodeClickAction:(id)sender {
    NSString *srcString = self.srcCodeTextView.string;
    
    if (srcString.length <= 0)
        return;
    
    if ([self.codecPopUpButton.titleOfSelectedItem isEqualToString:@"Base64"]) {
        
        if ([self.charsetPopUpButton.titleOfSelectedItem isEqualToString:@"UTF-8"]) {
            NSData *data = [srcString dataUsingEncoding:NSUTF8StringEncoding];
            NSString *result = [data base64EncodedStringWithOptions:NSDataBase64Encoding64CharacterLineLength];
            self.resultTextView.string = result;
            
        } else if ([self.charsetPopUpButton.titleOfSelectedItem isEqualToString:@"GBK"]) {
            NSStringEncoding gbkEncoding = CFStringConvertEncodingToNSStringEncoding(kCFStringEncodingGB_18030_2000);
            NSData *data = [srcString dataUsingEncoding:gbkEncoding];
            NSString *result = [data base64EncodedStringWithOptions:0];
            self.resultTextView.string = result;
        }
        
    } else if ([self.codecPopUpButton.titleOfSelectedItem isEqualToString:@"URL"]) {
        
        
        if ([self.charsetPopUpButton.titleOfSelectedItem isEqualToString:@"UTF-8"]) {
            NSData *data = [srcString dataUsingEncoding:NSUTF8StringEncoding];
            NSString *urlString = [[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding];
            self.resultTextView.string = [urlString stringByAddingPercentEncodingWithAllowedCharacters:NSCharacterSet.URLQueryAllowedCharacterSet];
            
        } else if ([self.charsetPopUpButton.titleOfSelectedItem isEqualToString:@"GBK"]) {
            NSStringEncoding gbkEncoding = CFStringConvertEncodingToNSStringEncoding(kCFStringEncodingGB_18030_2000);
            NSData *data = [srcString dataUsingEncoding:gbkEncoding];
            NSString *urlString = [[NSString alloc] initWithData:data encoding:gbkEncoding];
            self.resultTextView.string = [urlString stringByAddingPercentEncodingWithAllowedCharacters:NSCharacterSet.URLQueryAllowedCharacterSet];
        }
//        self.resultTextView.string = [srcString stringByAddingPercentEncodingWithAllowedCharacters:NSCharacterSet.URLQueryAllowedCharacterSet];
    }
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


- (void)sendAlertWithString:(NSString *)title subTitle:(nonnull NSString *)subTitle
{
    NSUserNotification *noti = [[NSUserNotification alloc] init];
    noti.title = title;
    noti.subtitle = subTitle;
    noti.hasActionButton = YES;
    noti.actionButtonTitle = @"OK";
    noti.otherButtonTitle = @"Cancel";
    
    [[NSUserNotificationCenter defaultUserNotificationCenter] scheduleNotification:noti];
    [[NSUserNotificationCenter defaultUserNotificationCenter] setDelegate:self];
    [NSTimer scheduledTimerWithTimeInterval:5.0 repeats:NO block:^(NSTimer * _Nonnull timer) {
        [[NSUserNotificationCenter defaultUserNotificationCenter] removeDeliveredNotification:noti];
    }];
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

@end
