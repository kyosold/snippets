#ifndef __UTILS_H_
#define __UTILS_H_

int get_uuid(char *uuid, size_t uuid_size);

int kscal(long long size, char *str, size_t str_size);

int conv_month_to_string(int month, char *str, size_t str_size);

int conv_month_to_int(char *str);

#endif