package gophia

/*
#cgo LDFLAGS: -lsophia
#include <sophia.h>

int sp_ctl_dir(void *p, uint32_t access, char *dir) {
	return sp_ctl(p, SPDIR, access, dir);
}

int go_sp_comparator(void *a, size_t asz, void *b, size_t bsz, void *arg);

int sp_ctl_cmp(void *p, void *cmp) {
	return sp_ctl(p, SPCMP, &go_sp_comparator, cmp);
}

int sp_ctl_page(void *p, uint32_t count) {
	return sp_ctl(p, SPPAGE, count);
}

int sp_ctl_gc(void *p, int active) {
	return sp_ctl(p, SPGC, active);
}

int sp_ctl_gcf(void *p, double factor) {
	return sp_ctl(p, SPGCF, factor);
}

int sp_ctl_grow(void *p, uint32_t newsize, double factor) {
	return sp_ctl(p, SPGROW, newsize, factor);
}

int sp_ctl_merge(void *p, int merge) {
	return sp_ctl(p, SPMERGE, merge);
}

int sp_ctl_mergewm(void *p, uint32_t watermark) {
	return sp_ctl(p, SPMERGEWM, watermark);
}

*/
import "C"
