//
//  CryptoViewController.m
//  STools
//
//  Created by Grayson on 2020/8/14.
//  Copyright © 2020 Grayson. All rights reserved.
//

#import "CryptoViewController.h"
#import<CommonCrypto/CommonDigest.h>

@interface CryptoViewController ()

@end

@implementation CryptoViewController

- (void)viewDidLoad {
    [super viewDidLoad];
    // Do view setup here.
    
    self.srcCodeTextView.textColor = [NSColor whiteColor];
    self.resultTextView.textColor = [NSColor whiteColor];
}

- (IBAction)cryptoClickAction:(id)sender {
    NSString *srcString = self.srcCodeTextView.string;
    
    if ([self.cryptoPopUpButton.titleOfSelectedItem isEqualToString:@"MD5"]) {
        NSData *data = [srcString dataUsingEncoding:NSUTF8StringEncoding];
        uint8_t digest[CC_MD5_DIGEST_LENGTH];
        CC_MD5(data.bytes, (unsigned int)data.length, digest);
        NSMutableString *output = [NSMutableString stringWithCapacity:CC_MD5_DIGEST_LENGTH * 2];
        for (int i=0; i<CC_MD5_DIGEST_LENGTH; i++) {
            [output appendFormat:@"%02x", digest[i]];
        }
        self.resultTextView.string = output;
        
    } else if ([self.cryptoPopUpButton.titleOfSelectedItem isEqualToString:@"SHA1"]) {
        NSData *data = [srcString dataUsingEncoding:NSUTF8StringEncoding];
        uint8_t digest[CC_SHA1_DIGEST_LENGTH];
        CC_SHA1(data.bytes, (unsigned int)data.length, digest);
        NSMutableString *output = [NSMutableString stringWithCapacity:CC_SHA1_DIGEST_LENGTH * 2];
        for (int i=0; i<CC_SHA1_DIGEST_LENGTH; i++) {
            [output appendFormat:@"%02x", digest[i]];
        }
        self.resultTextView.string = output;
        
    } else if ([self.cryptoPopUpButton.titleOfSelectedItem isEqualToString:@"SHA256"]) {
        NSData *data = [srcString dataUsingEncoding:NSUTF8StringEncoding];
        uint8_t digest[CC_SHA256_DIGEST_LENGTH];
        CC_SHA256(data.bytes, (unsigned int)data.length, digest);
        NSMutableString *output = [NSMutableString stringWithCapacity:CC_SHA256_DIGEST_LENGTH * 2];
        for (int i=0; i<CC_SHA256_DIGEST_LENGTH; i++) {
            [output appendFormat:@"%02x", digest[i]];
        }
        self.resultTextView.string = output;
        
    }
}
@end
