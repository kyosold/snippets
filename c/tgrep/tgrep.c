#define _GNU_SOURCE
#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <unistd.h>
#include <getopt.h>

#include "utils.h"
#include "xmalloc.h"

struct matcher_t
{
    int start_hour;
    int end_hour;
    int ignore_case;
    int verbose;
    char pattern[2048];
    char file[2048];
    int show_filename;
};

struct line_time_t
{
    int month;
    int day;
    int hour;
    int minute;
    int second;
};

struct ret_pos_t
{
    long pos;
    struct line_time_t linetime;
};

const float step_c = 0.2;

void copy_linetime(struct line_time_t *dst, struct line_time_t *src);

unsigned char eolbyte = '\n';

char *buffer;                   // Base of buffer
size_t bufalloc;                // Allocated buffer size, counting slop.
#define INITIAL_BUFSIZE 32768   // Initial buffer size, not counting slop.
// char *bufbeg;       // Beginning of user-visiable stuff.
// char *buflim;       // Limit of user-visiable stuff.
size_t pagesize;    // alignment of memory pages

/*  Return VAL aligned to the next multiple of ALIGNMENT.  VAL can be
    an integer or a pointer.  Both args must be free of side effects.  */
#define ALIGN_TO(val, alignment) \
    ((size_t) (val) % (alignment) == 0 \
    ? (val) \
    : (val) + ((alignment) - (size_t) (val) % (alignment)))

int reset(int fd, char const *file)
{
    if (!pagesize) {
        pagesize = getpagesize();
        if (pagesize == 0 || 2 * pagesize + 1 <= pagesize) {
            abort();
        }
        bufalloc = ALIGN_TO (INITIAL_BUFSIZE, pagesize) + pagesize + 1;
        buffer = (char *) xmalloc(bufalloc);
        // printf("buffer alloc: %d\n", bufalloc);
    }

    // bufbeg = buflim = ALIGN_TO (buffer + 1, pagesize);
    // printf("bufbeg = buflim = %d\n", bufbeg);
    return 1;
}


int get_line_time(char *line, struct line_time_t *linetime)
{
    char month[128] = {0};
    char day[128] = {0};
    char time[128] = {0};
    int i = 0, j = 0;
    int x = 0;
    for (i = 0; i < strlen(line); i++) {
        if (line[i] == ' ') {
            if (line[i+1] == ' ') {
                continue;
            } else {
                x++;
                j = 0;
            }
            continue;
        }

        if (x == 0) {
            month[j++] = line[i]; 
            month[j] = 0;
        } else if (x == 1) {
            day[j++] = line[i];
            day[j] = 0;
        } else if (x == 2) {
            time[j++] = line[i];
            time[j] = 0;
        } else {
            break;
        }
    }

    int mm = conv_month_to_int(month);
    int dd = atoi(day);
    int hou = 0;
    int min = 0;
    int sec = 0;
    i = 0;
    j = 0;
    x = 0;
    char tmp[128] = {0};
    while (time[i] != 0) {
        if (time[i] == ':') {
            if (x == 0) {
                hou = atoi(tmp);
            } else if (x == 1) {
                min = atoi(tmp);
            } else {
                sec = atoi(tmp);
            }
            x++;
            i++;
            j = 0;
            continue;
        }
        tmp[j++] = time[i++];
        tmp[j] = 0;
    }
    sec = atoi(tmp);

    // printf("month(%d) day(%d) hour(%d) min(%d) sec(%d)\n",
        // mm, dd, hou, min, sec);
    
    linetime->month = mm;
    linetime->day = dd;
    linetime->hour = hou;
    linetime->minute = min;
    linetime->second = sec;

    return 0;
}

int get_pos(char *file, int match_hour, int v, int flag, struct ret_pos_t *retpos)
{
    int is_a = 0, is_b = 0;
    struct line_time_t linetime;
    struct line_time_t linetime_last;
    char sline[4096] = {0};

    struct stat sb;
    if (stat(file, &sb) == -1) {
        printf("Error: %s\n", strerror(errno));
        return 1;
    }

    long spos = 0;
    long epos = sb.st_size;
    long cpos = epos / 2;

    FILE *fp = fopen(file, "r");
    if (fp == NULL) {
        printf("Error: %s\n", strerror(errno));
        return 1;
    }

    if (match_hour == 0 && flag == 0) {
        // 从0点开始
        if (fgets(sline, sizeof(sline), fp) == NULL) {
            printf("Error: %s\n", strerror(errno));
            goto FAIL_EXIT;
        }
        get_line_time(sline, &linetime);
        if (v > 0) {
            printf("  Z: MatchHour[%d] LineHour[%d] spos:0\n", match_hour, linetime.hour);
        }
        retpos->pos = 0;
        copy_linetime(&retpos->linetime, &linetime);
        // retpos->linetime.month = linetime.month;
        // retpos->linetime.day = linetime.day;
        // retpos->linetime.hour = linetime.hour;
        // retpos->linetime.minute = linetime.minute;
        // retpos->linetime.second = linetime.second;

        goto SUCC_EXIT;
    }

    fseek(fp, cpos, SEEK_SET);
    fgets(sline, sizeof(sline), fp);

    while (fgets(sline, sizeof(sline), fp) != NULL) {
        get_line_time(sline, &linetime);
        if (match_hour <= linetime.hour) { // 在上面
            is_a = 1;
            if (v > 0) {
                printf("  A: MatchHour[%d] LineHour[%d] spos:%ld cpos:%ld epos:%ld\n", match_hour, linetime.hour, spos, cpos, epos);
            }

            if (match_hour == 0 && flag == 1 && linetime.hour == 0) {
                // 修正: -e 为0的情况
                copy_linetime(&linetime, &linetime_last);
                // linetime.month = linetime_last.month;
                // linetime.day = linetime_last.day;
                // linetime.hour = linetime_last.hour;
                // linetime.minute = linetime_last.minute;
                // linetime.second = linetime_last.second;
                break;
            }

            if (is_b) {
                epos = cpos;
                break;
            }
            if (cpos == epos) {
                break;
            }

            epos = cpos;
            cpos = epos / 2;
            fseek(fp, cpos, SEEK_SET);
            fgets(sline, sizeof(sline), fp);

        } else if (match_hour > linetime.hour) {    // 在下面
            is_b = 1;
            if (v > 0) {
                printf("  B: MatchHour[%d] LineHour[%d] spos:%ld cpos:%ld epos:%ld\n", match_hour, linetime.hour, spos, cpos, epos);
            }

            if (is_a) {
                spos = cpos;
                break;
            }
            if (spos == cpos) {
                break;
            }

            spos = cpos;
            cpos = (epos - cpos) / 2 + cpos;
            fseek(fp, cpos, SEEK_SET);
            fgets(sline, sizeof(sline), fp);
        }

        copy_linetime(&linetime_last, &linetime);
        // linetime_last.month = linetime.month;
        // linetime_last.day = linetime.day;
        // linetime_last.hour = linetime.hour;
        // linetime_last.minute = linetime.minute;
        // linetime_last.second = linetime.second;
    }
    if (feof(fp)) {
        printf("---EOF---\n");
        goto SUCC_EXIT;
    }

    copy_linetime(&linetime_last, &linetime);
    // linetime_last.month = linetime.month;
    // linetime_last.day = linetime.day;
    // linetime_last.hour = linetime.hour;
    // linetime_last.minute = linetime.minute;
    // linetime_last.second = linetime.second;

    long step = (float)(epos - spos) * step_c;
    if (v > 0) {
        printf("  Start Step: MatchHour[%d] LineHour[%d] Step:%ld spos:%ld --- epos:%ld\n", match_hour, linetime.hour, step, spos, epos);
    }

    if (step > 0) {
        cpos = spos + step;
        fseek(fp, cpos, SEEK_SET);
        fgets(sline, sizeof(sline), fp);
        while (fgets(sline, sizeof(sline), fp) != NULL) {
            get_line_time(sline, &linetime);
            if (match_hour > linetime.hour) {
                // 继续向下找
                if (v > 0) {
                    printf("  S: [STEP] MatchHour[%d] LineHour[%d]: cpos:%ld\n", match_hour, linetime.hour, cpos);
                }
                if (cpos > epos) {
                    printf("  S: [STEP] E cpos(%ld) > epos(%ld), exit.\n", cpos, epos);
                    break;
                }

                copy_linetime(&linetime_last, &linetime);
                cpos = cpos + step;
                fseek(fp, cpos, SEEK_SET);
                fgets(sline, sizeof(sline), fp);
            } else {
                if (flag == 0) {
                    if (v > 0) {
                        printf("  S: [STEP] E MatchHour[%d] LineHour[%d] cpos:%ld\n", match_hour, linetime.hour, cpos);
                    }
                    cpos = cpos - step;
                    copy_linetime(&linetime, &linetime_last);
                    break;
                } else {
                    if (match_hour < linetime.hour) {
                        if (v > 0) {
                            printf("  S: [STEP] E MatchHour[%d] LineHour[%d] cpos:%ld\n", match_hour, linetime.hour, cpos);
                        }
                        break;
                    }
                    if (v > 0) {
                        printf("  S: [STEP] MatchHour[%d] LineHour[%d] cpos:%ld\n", match_hour, linetime.hour, cpos);
                    }
                    cpos = cpos + step;
                    fseek(fp, cpos, SEEK_SET);
                    fgets(sline, sizeof(sline), fp);
                }
            }
        }

    } else {
        if (flag == 0) {
            cpos = spos;
        } else {
            cpos = epos;
        }
    }

    retpos->pos = cpos;
    copy_linetime(&retpos->linetime, &linetime);
    // retpos->linetime.month = linetime.month;
    // retpos->linetime.day = linetime.day;
    // retpos->linetime.hour = linetime.hour;
    // retpos->linetime.minute = linetime.minute;
    // retpos->linetime.second = linetime.second;

    goto SUCC_EXIT;

FAIL_EXIT:
    fclose(fp);
    return 1;

SUCC_EXIT:
    fclose(fp);
    return 0;
}

void grep_str(struct matcher_t *matcher, char *buf, char *substr)
{
    char *pos = NULL;
    if (matcher->ignore_case) {
        pos = strcasestr(buf, substr);
    } else {
        pos = strstr(buf, substr);
    }

    if (pos != NULL) {
        // found
        if (matcher->show_filename) {
            printf("%s: %s%c", matcher->file, buf, eolbyte);
        } else {
            printf("%s%c", buf, eolbyte);
        }
    }
}


int do_search(struct matcher_t *matcher)
{
    int hour_matched = 0;
    int ret = 0;

    // get start pos
    struct ret_pos_t ret_pos_s;
    ret = get_pos(matcher->file, matcher->start_hour, matcher->verbose, 0, &ret_pos_s);
    if (ret != 0) {
        return 1;
    }
    // printf("Start pos:%ld LineTime:%d:%d\n", ret_pos_s.pos, ret_pos_s.linetime.hour, ret_pos_s.linetime.minute);

    // get end pos
    struct ret_pos_t ret_pos_e;
    if (matcher->end_hour >= 0) {
        if (matcher->verbose > 0) {
            printf("\n");
        }
        ret = get_pos(matcher->file, matcher->end_hour, matcher->verbose, 1, &ret_pos_e);
        if (ret != 0) {
            return 1;
        }
    } else {
        struct stat sb;
        if (stat(matcher->file, &sb) == -1) {
            printf("Error: %s\n", strerror(errno));
            return 1;
        }
        ret_pos_e.pos = sb.st_size;
        ret_pos_e.linetime.hour = -1;
        ret_pos_e.linetime.minute = -1;
    }

    char scal_str[128] = {0};
    kscal(ret_pos_e.pos - ret_pos_s.pos, scal_str, sizeof(scal_str));

    printf("--------------------------\n");
    printf("Final: Start:[%ld][%d:%d] --- End:[%ld][%d:%d] Need Scan[%s]\n", 
        ret_pos_s.pos, ret_pos_s.linetime.hour, ret_pos_s.linetime.minute, 
        ret_pos_e.pos, ret_pos_e.linetime.hour, ret_pos_e.linetime.minute,
        scal_str);
    printf("--------------------------\n");
    printf("\n");

    int fd = open(matcher->file, O_RDONLY);
    if (fd == -1) {
        printf("Error: open file(%s) fail: %s\n", matcher->file, strerror(errno));
        return 1;
    }

    if (!reset(fd, matcher->file)) {
        goto ERR_EXIT;
    }

    lseek(fd, ret_pos_s.pos, SEEK_SET);

    char *line_buf = NULL;
    char *line_end = NULL;
    int len = 0;
    long total_bytes = ret_pos_e.pos - ret_pos_s.pos;
    while (total_bytes > 0) {
        if (total_bytes < bufalloc)
            bufalloc = total_bytes;

        ssize_t n = read (fd, buffer + len, bufalloc);
        if (n < 0) {
            printf("Error: read file(%s) fail: %s\n", matcher->file, strerror(errno));
            close(fd);
            return 1;
        } else if (n == 0) {
            printf("Error: read eof\n");
            close(fd);
            return 0;
        }

        total_bytes -= n;

        if (len == 0) {
            // bufbeg 是第一行开始的指针
            line_buf = memchr(buffer, eolbyte, n);
            if (line_buf) {
                line_buf++;
                int m = line_buf - buffer;
                n -= m;
            } else {
                printf("Error: buffer alloca is small, is not include \\n at buffer, (you should inc buffer alloca size)\n");
                goto ERR_EXIT;
            }
        } else {
            line_buf = buffer;
            n += len;
        }

        if (line_buf) {
            for (;;) {
                line_end = memchr(line_buf, eolbyte, n);
                if (line_end) {
                    *line_end = 0;
                } else {
                    // 保留 line_buf，不是完整的一行
                    memmove(buffer, line_buf, n);
                    len = n;
                    break;
                }

                grep_str(matcher, line_buf, matcher->pattern);
                n -= (line_end - line_buf + 1);
                line_buf = line_end + 1;
            }

        }

    }

ERR_EXIT:
    close(fd);
    return 1;

OK_EXIT:
    close(fd);
    return 0;
}

void usage(char *prog)
{
    printf("----------------------------------------------\n");
    printf("Usage:\n");
    printf("  %s -v -s15 -e17 [pattern] [file]\n", prog);
    printf("\n");
    printf("Options:\n");
    printf("  -s: 从指定的小时开始查找\n");
    printf("  -e: 到指定的小时结束查找, 不写默认到文件结尾\n");
    printf("  -i: 不区分大小写\n");
    printf("  -v: 显示详细信息\n");
    printf("  -S: 显示文件名\n");
    printf("\n");
    printf("Others:\n");
    printf("  1. 如果文件是gzip，先解压再查询，如:\n");
    printf("    a. gzip -dvc abc.0.gz > abc.0\n");
    printf("    b. tgrep -st=8 -et=10 'pattern' abc.0\n");
    printf("----------------------------------------------\n");
    printf("\n");
}

int main(int argc, char **argv)
{
    if (argc < 2) {
        usage(argv[0]);
        return 1;
    }

    struct matcher_t matcher;
    matcher.start_hour = -1;
    matcher.end_hour = -1;
    matcher.ignore_case = 0;
    matcher.verbose = 0;
    matcher.show_filename = 0;

    int i;
    int ch;
    const char *args = "s:e:Sivh";
    while ((ch = getopt(argc, argv, args)) != -1) {
        switch (ch)
        {
        case 's':
            matcher.start_hour = atoi(optarg);
            break;
        case 'e':
            matcher.end_hour = atoi(optarg);
            break;
        case 'i':
            matcher.ignore_case = 1;
            break;
        case 'v':
            matcher.verbose = 1;
            break;
        case 'S':
            matcher.show_filename = 1;
            break;
        
        case 'h':
        default:
            usage(argv[0]);
            exit(1);
        }
    }
    if ((argc - optind) != 2) {
        usage(argv[0]);
        return 1;
    }

    snprintf(matcher.pattern, sizeof(matcher.pattern), "%s", argv[optind++]);
    snprintf(matcher.file, sizeof(matcher.file), "%s", argv[optind++]);

    printf("**************************\n");
    printf("Start Hour: %d\n", matcher.start_hour);
    printf("End Hour: %d\n", matcher.end_hour);
    printf("Ignore Case: %d\n", matcher.ignore_case);
    printf("Verbose: %d\n", matcher.verbose);
    printf("Pattern: %s\n", matcher.pattern);
    printf("File: %s\n", matcher.file);
    printf("**************************\n");

    if (matcher.start_hour > matcher.end_hour) {
        printf("Error: -s can't more than -e value\n");
        return 1;
    }

    char *ext = strrchr(matcher.file, '.');
    if (ext && (
            strcasecmp(ext, ".gz") == 0 ||
            strcasecmp(ext, ".gzip") == 0)
            ) {
        printf("Error: file(%s) is gzip file\n", matcher.file);
        printf("  Run: 'gzip -dvc %s > tmp.log'\n", matcher.file);
        return 1;
    }

    do_search(&matcher);

    // printf("month(%d) day(%d) hour(%d) min(%d) sec(%d)\n", linetime.month, linetime.day, linetime.hour, linetime.minute, linetime.second);

    return 0;
}

void copy_linetime(struct line_time_t *dst, struct line_time_t *src) 
{
    dst->month = src->month;
    dst->day = src->day;
    dst->hour = src->hour;
    dst->minute = src->minute;
    dst->second = src->second;
}

