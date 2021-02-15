// +build integration

package rgw

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"testing"
)

// Run with S3_ACCESS_KEY="…" S3_SECRET_KEY="…" CEPH_CONF="…ceph.conf"
func TestSimpleSetupAndTouchFile(t *testing.T) {
	ret, rgw := Create()
	if ret != 0 {
		t.Fatalf("RGW Create failed: %v", ret)
	}

	ret, rgwfs := Mount(rgw,
		"test",
		os.Getenv("S3_ACCESS_KEY"),
		os.Getenv("S3_SECRET_KEY"),
		0)
	if ret == 0 {
		fmt.Printf("RGW Mounted: %v\n", rgwfs)
	} else {
		t.Fatalf("Failed to mount: %v\n", ret)
	}

	ret, statvfs := StatFs(rgwfs, rgwfs.root_fh, 0)
	if ret == 0 {
		fmt.Printf("Statfs: %+v\n", statvfs)
	} else {
		t.Fatalf("Statfs failed: %v", ret)
	}

	stat := NewStat(0, 0, 0755)
	createMask := SetAttrUID | SetAttrGID | SetAttrMode
	newDirName := uuid.NewString()
	ret, dirFh := Mkdir(rgwfs, rgwfs.root_fh, newDirName, stat,
		createMask, 0)
	if ret == 0 {
		fmt.Printf("Created new directory: %v  %+v %v %v\n", newDirName, stat, dirFh, ret)
	} else {
		t.Fatalf("Failed to create %v: %v\n", newDirName, ret)
	}

	stat = NewStat(0, 0, 0644)

	newFileName := uuid.NewString()
	ret, fh := CreateFile(rgwfs, dirFh, newFileName, stat,
		createMask, 0, 0)
	if ret == 0 {
		fmt.Printf("Created new file %v in %v: %+v %v %v\n", newFileName, newDirName, stat, fh, ret)
	} else {
		t.Fatalf("Failed to create %v: %v", newFileName, ret)
	}

	// dirCh := make(chan *ReaddirResult)
	// Readdir(rgwfs, rgwfs.root_fh, "", 0, dirCh)

	err, eof := ReadDir(rgwfs, rgwfs.root_fh, "", ReaddirFlagDotDot, func(name string, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Printf("Readdir: %v\n", name)
		return nil
	})
	fmt.Println(err, eof)
	fmt.Println("Read first entry:")

	err, eof = ReadDir(rgwfs, rgwfs.root_fh, "", ReaddirFlagDotDot, func(name string, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Printf("Readdir: %v\n", name)
		return ErrStopReadDir
	})

	ret = Umount(rgwfs, 0)
	fmt.Println(ret)
}
