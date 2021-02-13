package rgw

// #cgo LDFLAGS: -lrgw
// #include <stdio.h>
// #include <stdlib.h>
// #include <sys/stat.h>
// #include <rados/librgw.h>
// #include <rados/rgw_file.h>
import "C"
import "unsafe"

// Mask bits for new file / directory functions
const (
	SetAttrUID  uint32 = C.RGW_SETATTR_UID
	SetAttrGID  uint32 = C.RGW_SETATTR_GID
	SetAttrMode uint32 = C.RGW_SETATTR_MODE
)

// NewStat returns a C struct stat initialized with uid, gid, mode
func NewStat(uid, gid, mode int) *C.struct_stat {
	stat := new(C.struct_stat)
	stat.st_uid = C.uid_t(uid)
	stat.st_gid = C.gid_t(gid)
	stat.st_mode = C.mode_t(mode)
	return stat
}

// Create new RGW session
func Create() (int, C.librgw_t) {
	cstr := C.CString("")
	defer C.free(unsafe.Pointer(cstr))
	var ptr C.librgw_t
	ret := C.librgw_create(&ptr, C.int(1), &cstr)
	return int(ret), ptr
}

// Shutdown RGW session
func Shutdown(rgw C.librgw_t) {
	C.librgw_shutdown(rgw)
}

// Mount attaches to an RGW namespace and returns a handle to perform operations on that namespace
func Mount(rgw C.librgw_t, uid, key, secret string, flags uint32) (int, *C.struct_rgw_fs) {
	var rgwfs *C.struct_rgw_fs
	cuid := C.CString(uid)
	defer C.free(unsafe.Pointer(cuid))
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	csecret := C.CString(secret)
	defer C.free(unsafe.Pointer(csecret))
	ret := C.rgw_mount(rgw, cuid, ckey, csecret, &rgwfs, C.uint(flags))
	return int(ret), rgwfs
}

// Umount detaches from an RGW namespace
func Umount(rgwFs *C.struct_rgw_fs, flags uint32) int {
	ret := C.rgw_umount(rgwFs, C.uint(flags))
	return int(ret)
}

// StatFs returns filesystem attributes
func StatFs(rgwFs *C.struct_rgw_fs, parentFh *C.struct_rgw_file_handle, flags uint32) (int, C.struct_rgw_statvfs) {
	var statvfs C.struct_rgw_statvfs
	ret := C.rgw_statfs(rgwFs, parentFh, &statvfs, C.uint(flags))
	return int(ret), statvfs
}

// CreateFile creates and returns a file handle to a new file
func CreateFile(rgwFs *C.struct_rgw_fs, parentFh *C.struct_rgw_file_handle, name string,
	stat *C.struct_stat, mask uint32, posixFlags uint32, flags uint32) (int, *C.struct_rgw_file_handle) {

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var fh *C.struct_rgw_file_handle

	ret := C.rgw_create(rgwFs, parentFh, cname, stat, C.uint(mask), &fh, C.uint(posixFlags), C.uint(flags))

	return int(ret), fh
}

// Mkdir creates and returns a handle to a new directory
func Mkdir(rgwFs *C.struct_rgw_fs, parentFh *C.struct_rgw_file_handle, name string,
	stat *C.struct_stat, mask uint32, flags uint32) (int, *C.struct_rgw_file_handle) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var fh *C.struct_rgw_file_handle
	ret := C.rgw_mkdir(rgwFs, parentFh, cname, stat, C.uint(mask), &fh, C.uint(flags))
	return int(ret), fh
}
