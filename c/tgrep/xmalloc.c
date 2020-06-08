#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "xmalloc.h"

/*  1 if calloc is known to be compatible with GNU calloc.  This
    matters if we are not also using the calloc module, which defines
    HAVE_CALLOC and supports the GNU API even on non-GNU platforms.  */
#if defined HAVE_CALLOC || defined __GLIBC__
    enum { HAVE_GNU_CALLOC = 1 };
#else
    enum { HAVE_GNU_CALLOC = 0 };
#endif



void xalloc_die(void)
{
    printf("memory exhausted\n");
    abort();
}

/*  Allocate N bytes of memory dynamically, with error checking. */
void *
xmalloc (size_t n)
{
    void *p = malloc(n);
    if (!p && n != 0)
        xalloc_die();
    return p;
}

/*  Change the size of an allocated block of memory P to N bytes,
    with error checking. */
void *
xrealloc (void *p, size_t n)
{
    p = realloc(p, n);
    if (!p && n != 0)
        xalloc_die();
    return p;
}

/*  Allocate S bytes of zeroed memory dynnamically, with error checking.
    There's no need for xnzalloc (N, S), since it would be equivalent
    to xcalloc (N, S). */
void *
xzalloc (size_t s)
{
    return memset(xmalloc(s), 0, s);
}

/*  Allocate zeroed memory for N elements of S bytes, with error
    checking. S must be nonzero. */
void *
xcalloc (size_t n, size_t s)
{
    void *p;
    /*  Test for overflow, since some calloc implementations don't have
        proper overflow checks.  But omit overflow and size-zero tests if
        HAVE_GNU_CALLOC, since GNU calloc catches overflow and never
        returns NULL if successful.  */
    if ((! HAVE_GNU_CALLOC && xalloc_oversized(n, s))
            || (! (p = calloc(n, s)) && (HAVE_GNU_CALLOC || n != 0)))
        xalloc_die();
    return p;
}

/*  Clone an object P of size S, with error checking. There's no need
    for xnmemdup (P, N, S), since xmemdup(P, N * S) works without any
    need for arithmetic overflow check. */
void *
xmemdup (void const *p, size_t s)
{
    return memcpy(xmalloc(s), p, s);
}

/*  Clone STRING. */
char *
xstrdup (char const *string)
{
    return xmemdup(string, strlen(string) + 1);
}
