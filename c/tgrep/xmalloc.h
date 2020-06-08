#ifndef __XMALLOC_H__
#define __XMALLOC_H__

#include <stddef.h>

/*  This function is always triggered when memory is exhausted.
    It must be defined by the application, either explicitly
    or by using gnulib's xalloc-die module.  This is the
    function to call when one wants the program to die because of a
    memory allocation failure.  */
void xalloc_die (void) ;

void *xmalloc (size_t s) ;
void *xzalloc (size_t s) ;
void *xcalloc (size_t n, size_t s) ;
void *xrealloc (void *p, size_t s);
void *x2realloc (void *p, size_t *pn);
void *xmemdup (void const *p, size_t s) ;
char *xstrdup (char const *str) ;


/*  Return 1 if an array of N objects, each of size S, cannot exist due
    to size arithmetic overflow.  S must be positive and N must be
    nonnegative.  This is a macro, not an inline function, so that it
    works correctly even when SIZE_MAX < N.

    By gnulib convention, SIZE_MAX represents overflow in size
    calculations, so the conservative dividend to use here is
    SIZE_MAX - 1, since SIZE_MAX might represent an overflowed value.
    However, malloc (SIZE_MAX) fails on all known hosts where
    sizeof (ptrdiff_t) <= sizeof (size_t), so do not bother to test for
    exactly-SIZE_MAX allocations on such hosts; this avoids a test and
    branch when S is known to be 1. */
#define xalloc_oversized(n, s) \
    ((size_t) (sizeof(ptrdiff_t) <= sizeof(size_t) ? -1 : -2) / (s) < (n))

#endif
