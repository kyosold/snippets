//
//  CryptoViewController.h
//  STools
//
//  Created by Grayson on 2020/8/14.
//  Copyright © 2020 Grayson. All rights reserved.
//

#import <Cocoa/Cocoa.h>

NS_ASSUME_NONNULL_BEGIN

@interface CryptoViewController : NSViewController

@property (unsafe_unretained) IBOutlet NSTextView *srcCodeTextView;
@property (weak) IBOutlet NSPopUpButton *cryptoPopUpButton;
@property (unsafe_unretained) IBOutlet NSTextView *resultTextView;

- (IBAction)cryptoClickAction:(id)sender;

@end

NS_ASSUME_NONNULL_END
