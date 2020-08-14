//
//  DumpTableViewData.m
//  STools
//
//  Created by Grayson on 2020/8/7.
//  Copyright © 2020 Grayson. All rights reserved.
//

#import "DumpTableViewData.h"

@implementation DumpTableViewData

- (id)init{
    self = [super init];
    if (self) {
        self.rowData = [[NSMutableArray alloc] init];
    }
    return self;
}

- (NSInteger)numberOfRowsInTableView:(NSTableView *)tableView
{
    return self.rowData.count;
}

- (NSView *)tableView:(NSTableView *)tableView viewForTableColumn:(NSTableColumn *)tableColumn row:(NSInteger)row
{
    NSString *identifier = [tableColumn identifier];
    NSDictionary *dict = [self.rowData objectAtIndex:row];
    NSString *value = [dict objectForKey:identifier];
    if (value) {
        NSTableCellView *cell = [tableView makeViewWithIdentifier:identifier owner:self];
        cell.textField.stringValue = value;
        return cell;
    }
    return nil;
}

- (void)appendMsgToDumpTableView:(NSDictionary *)dict
{
    [self.rowData addObject:dict];
    [self.dumpTableView reloadData];
}

@end
