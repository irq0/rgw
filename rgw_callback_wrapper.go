package rgw

/*
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>
#include <sys/stat.h>

bool ReaddirCallback(char *name, void *arg, uint64_t offset,
		struct stat *st, uint32_t st_mask, uint32_t flags);

bool ReaddirCallbackCgo(const char *name, void *arg, uint64_t offset,
		struct stat *st, uint32_t st_mask, uint32_t flags)
{
	return ReaddirCallback((char*)name, arg, offset, st, st_mask, flags);
}
*/
import "C"
