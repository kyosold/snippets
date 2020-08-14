//
//  DumpTableViewData.h
//  STools
//
//  Created by Grayson on 2020/8/7.
//  Copyright © 2020 Grayson. All rights reserved.
//

#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>

NS_ASSUME_NONNULL_BEGIN

@interface DumpTableViewData : NSObject <NSTableViewDelegate, NSTableViewDataSource>

@property (weak) IBOutlet NSTableView *dumpTableView;

@property NSMutableArray *rowData;

- (void)appendMsgToDumpTableView:(NSDictionary *)dict;

@end

NS_ASSUME_NONNULL_END
