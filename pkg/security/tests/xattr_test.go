// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build linux && functionaltests

// Package tests holds tests related files
package tests

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"testing"
	"unsafe"

	"github.com/iceber/iouring-go"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sys/unix"

	"github.com/DataDog/datadog-agent/pkg/security/secl/model"
	"github.com/DataDog/datadog-agent/pkg/security/secl/rules"
)

func TestSetXAttr(t *testing.T) {
	SkipIfNotAvailable(t)

	rule := &rules.RuleDefinition{
		ID:         "test_rule_xattr",
		Expression: `((setxattr.file.path == "{{.Root}}/test-setxattr" && setxattr.file.uid == 98 && setxattr.file.gid == 99) || setxattr.file.path == "{{.Root}}/test-setxattr-link") && setxattr.file.destination.namespace == "user" && setxattr.file.destination.name == "user.test_xattr"`,
	}

	testDrive, err := newTestDrive(t, "xfs", nil, "")
	if err != nil {
		t.Fatal(err)
	}
	defer testDrive.Close()

	test, err := newTestModule(t, nil, []*rules.RuleDefinition{rule}, withDynamicOpts(dynamicTestOpts{testDir: testDrive.Root()}))
	if err != nil {
		t.Fatal(err)
	}
	defer test.Close()

	xattrName, err := syscall.BytePtrFromString("user.test_xattr")
	if err != nil {
		t.Fatal(err)
	}
	xattrNamePtr := unsafe.Pointer(xattrName)
	xattrValuePtr := unsafe.Pointer(&[]byte{})

	fileMode := 0o777
	expectedMode := uint16(applyUmask(fileMode))

	t.Run("setxattr", func(t *testing.T) {
		testFile, testFilePtr, err := test.CreateWithOptions("test-setxattr", 98, 99, fileMode)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(testFile)

		test.WaitSignal(t, func() error {
			_, _, errno := syscall.Syscall6(syscall.SYS_SETXATTR, uintptr(testFilePtr), uintptr(xattrNamePtr), uintptr(xattrValuePtr), 0, unix.XATTR_CREATE, 0)
			if errno != 0 {
				return error(errno)
			}
			return nil
		}, func(event *model.Event, _ *rules.Rule) {
			assert.Equal(t, "setxattr", event.GetType(), "wrong event type")
			assert.Equal(t, "user.test_xattr", event.SetXAttr.Name)
			assert.Equal(t, "user", event.SetXAttr.Namespace)
			assert.Equal(t, getInode(t, testFile), event.SetXAttr.File.Inode, "wrong inode")
			assertRights(t, event.SetXAttr.File.Mode, expectedMode)
			assertNearTime(t, event.SetXAttr.File.MTime)
			assertNearTime(t, event.SetXAttr.File.CTime)

			value, _ := event.GetFieldValue("event.async")
			assert.Equal(t, value.(bool), false)
		})
	})

	t.Run("lsetxattr", func(t *testing.T) {
		testFile, testFilePtr, err := test.Path("test-setxattr-link")
		if err != nil {
			t.Fatal(err)
		}

		testOldFile, _, err := test.CreateWithOptions("test-setxattr-old", 98, 99, fileMode)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(testOldFile)

		if err := os.Symlink(testOldFile, testFile); err != nil {
			t.Fatal(err)
		}
		defer os.Remove(testFile)

		test.WaitSignal(t, func() error {
			_, _, errno := syscall.Syscall6(syscall.SYS_LSETXATTR, uintptr(testFilePtr), uintptr(xattrNamePtr), uintptr(xattrValuePtr), 0, unix.XATTR_CREATE, 0)
			// Linux and Android don't support xattrs on symlinks according
			// to xattr(7), so just test that we get the proper error.
			// We should get the event though
			if errno != syscall.EACCES && errno != syscall.EPERM {
				return error(errno)
			}
			return nil
		}, func(event *model.Event, _ *rules.Rule) {
			assert.Equal(t, "setxattr", event.GetType(), "wrong event type")
			assert.Equal(t, "user.test_xattr", event.SetXAttr.Name)
			assert.Equal(t, "user", event.SetXAttr.Namespace)
			assert.Equal(t, getInode(t, testFile), event.SetXAttr.File.Inode, "wrong inode")
			assertRights(t, event.SetXAttr.File.Mode, 0777)
			assertNearTime(t, event.SetXAttr.File.MTime)
			assertNearTime(t, event.SetXAttr.File.CTime)

			value, _ := event.GetFieldValue("event.async")
			assert.Equal(t, value.(bool), false)
		})
	})

	t.Run("fsetxattr", func(t *testing.T) {
		testFile, _, err := test.CreateWithOptions("test-setxattr", 98, 99, fileMode)
		if err != nil {
			t.Fatal(err)
		}

		f, err := os.Open(testFile)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		defer os.Remove(testFile)

		test.WaitSignal(t, func() error {
			_, _, errno := syscall.Syscall6(syscall.SYS_FSETXATTR, f.Fd(), uintptr(xattrNamePtr), uintptr(xattrValuePtr), 0, unix.XATTR_CREATE, 0)
			if errno != 0 {
				return error(errno)
			}
			return nil
		}, func(event *model.Event, _ *rules.Rule) {
			assert.Equal(t, "setxattr", event.GetType(), "wrong event type")
			assert.Equal(t, "user.test_xattr", event.SetXAttr.Name)
			assert.Equal(t, "user", event.SetXAttr.Namespace)
			assert.Equal(t, getInode(t, testFile), event.SetXAttr.File.Inode, "wrong inode")
			assertRights(t, event.SetXAttr.File.Mode, expectedMode)
			assertNearTime(t, event.SetXAttr.File.MTime)
			assertNearTime(t, event.SetXAttr.File.CTime)

			value, _ := event.GetFieldValue("event.async")
			assert.Equal(t, value.(bool), false)
		})
	})

	t.Run("io_uring-fsetxattr", func(t *testing.T) {
		SkipIfNotAvailable(t)

		testFile, _, err := test.CreateWithOptions("test-setxattr", 98, 99, fileMode)
		if err != nil {
			t.Fatal(err)
		}
		defer syscall.Rmdir(testFile)

		iour, err := iouring.New(1)
		if err != nil {
			if errors.Is(err, unix.ENOTSUP) {
				t.Fatal(err)
			}
			t.Skip("io_uring not supported")
		}
		defer iour.Close()

		fd, err := unix.Open(testFile, unix.O_CREAT, uint32(fileMode))
		if err != nil {
			t.Fatal(err)
		}
		defer unix.Close(fd)

		prepRequest, err := iouring.Fsetxattr(int32(fd), "user.test_xattr", "foo", 0)
		if err != nil {
			t.Fatal(err)
		}

		ch := make(chan iouring.Result, 1)

		test.WaitSignal(t, func() error {
			if _, err = iour.SubmitRequest(prepRequest, ch); err != nil {
				return err
			}

			result := <-ch
			t.Logf("Got result from io-uring !")
			ret, err := result.ReturnInt()
			if err != nil {
				if err == syscall.EBADF || err == syscall.EINVAL {
					return ErrSkipTest{"fsetxattr not supported by io_uring"}
				}
				return err
			}

			if ret < 0 {
				return fmt.Errorf("failed to set extended attribute with io_uring: %d", ret)
			}

			t.Logf("Successfully set extended attribute ! %d", ret)
			return nil
		}, func(event *model.Event, rule *rules.Rule) {
			assertTriggeredRule(t, rule, "test_rule_xattr")
			assert.Equal(t, getInode(t, testFile), event.SetXAttr.File.Inode, "wrong inode")
			assertRights(t, uint16(event.SetXAttr.File.Mode), expectedMode)
			assertNearTime(t, event.SetXAttr.File.MTime)
			assertNearTime(t, event.SetXAttr.File.CTime)

			value, _ := event.GetFieldValue("event.async")
			assert.Equal(t, value.(bool), true)
		})
	})

	t.Run("io_uring-setxattr", func(t *testing.T) {
		SkipIfNotAvailable(t)

		testFile, _, err := test.CreateWithOptions("test-setxattr", 98, 99, fileMode)
		if err != nil {
			t.Fatal(err)
		}
		defer syscall.Rmdir(testFile)

		iour, err := iouring.New(1)
		if err != nil {
			if errors.Is(err, unix.ENOTSUP) {
				t.Fatal(err)
			}
			t.Skip("io_uring not supported")
		}
		defer iour.Close()

		prepRequest, err := iouring.Setxattr(testFile, "user.test_xattr", "foo", 0)
		if err != nil {
			t.Fatal(err)
		}

		ch := make(chan iouring.Result, 1)

		test.WaitSignal(t, func() error {
			if _, err = iour.SubmitRequest(prepRequest, ch); err != nil {
				return err
			}

			result := <-ch
			t.Logf("Got result from io-uring !")
			ret, err := result.ReturnInt()
			if err != nil {
				if err == syscall.EBADF || err == syscall.EINVAL {
					return ErrSkipTest{"fsetxattr not supported by io_uring"}
				}
				return err
			}

			if ret < 0 {
				return fmt.Errorf("failed to set extended attribute with io_uring: %d", ret)
			}

			t.Logf("Successfully set extended attribute ! %d", ret)
			return nil
		}, func(event *model.Event, rule *rules.Rule) {
			assertTriggeredRule(t, rule, "test_rule_xattr")
			assert.Equal(t, getInode(t, testFile), event.SetXAttr.File.Inode, "wrong inode")
			assertRights(t, uint16(event.SetXAttr.File.Mode), expectedMode)
			assertNearTime(t, event.SetXAttr.File.MTime)
			assertNearTime(t, event.SetXAttr.File.CTime)

			value, _ := event.GetFieldValue("event.async")
			assert.Equal(t, value.(bool), true)
		})
	})
}

func TestRemoveXAttr(t *testing.T) {
	SkipIfNotAvailable(t)

	ruleDefs := []*rules.RuleDefinition{
		{
			ID:         "test_rule_xattr",
			Expression: `((removexattr.file.path == "{{.Root}}/test-removexattr" && removexattr.file.uid == 98 && removexattr.file.gid == 99) || removexattr.file.path == "{{.Root}}/test-removexattr-link") && removexattr.file.destination.namespace == "user" && removexattr.file.destination.name == "user.test_xattr" `,
		},
	}

	testDrive, err := newTestDrive(t, "xfs", nil, "")
	if err != nil {
		t.Fatal(err)
	}
	defer testDrive.Close()

	test, err := newTestModule(t, nil, ruleDefs, withDynamicOpts(dynamicTestOpts{testDir: testDrive.Root()}))
	if err != nil {
		t.Fatal(err)
	}
	defer test.Close()

	xattrName, err := syscall.BytePtrFromString("user.test_xattr")
	if err != nil {
		t.Fatal(err)
	}
	xattrNamePtr := unsafe.Pointer(xattrName)

	fileMode := 0o777
	expectedMode := applyUmask(fileMode)

	t.Run("removexattr", func(t *testing.T) {
		testFile, testFilePtr, err := test.CreateWithOptions("test-removexattr", 98, 99, fileMode)
		if err != nil {
			t.Fatal(err)
		}

		f, err := os.Open(testFile)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		defer os.Remove(testFile)

		// set xattr
		_, _, errno := syscall.Syscall6(syscall.SYS_FSETXATTR, f.Fd(), uintptr(xattrNamePtr), 0, 0, 1, 0)
		if errno != 0 {
			t.Fatal(error(errno))
		}

		test.WaitSignal(t, func() error {
			_, _, errno = syscall.Syscall(syscall.SYS_REMOVEXATTR, uintptr(testFilePtr), uintptr(xattrNamePtr), 0)
			if errno != 0 {
				return error(errno)
			}
			return nil
		}, func(event *model.Event, _ *rules.Rule) {
			assert.Equal(t, "removexattr", event.GetType(), "wrong event type")
			assert.Equal(t, "user.test_xattr", event.RemoveXAttr.Name)
			assert.Equal(t, getInode(t, testFile), event.RemoveXAttr.File.Inode, "wrong inode")
			assertRights(t, event.RemoveXAttr.File.Mode, uint16(expectedMode))
			assertNearTime(t, event.RemoveXAttr.File.MTime)
			assertNearTime(t, event.RemoveXAttr.File.CTime)

			value, _ := event.GetFieldValue("event.async")
			assert.Equal(t, value.(bool), false)
		})
	})

	t.Run("lremovexattr", func(t *testing.T) {
		testFile, testFilePtr, err := test.Path("test-removexattr-link")
		if err != nil {
			t.Fatal(err)
		}

		testOldFile, _, err := test.CreateWithOptions("test-setxattr-old", 98, 99, fileMode)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(testOldFile)

		if err := os.Symlink(testOldFile, testFile); err != nil {
			t.Fatal(err)
		}
		defer os.Remove(testFile)

		// set xattr
		_, _, errno := syscall.Syscall6(syscall.SYS_LSETXATTR, uintptr(testFilePtr), uintptr(xattrNamePtr), 0, 0, 1, 0)
		// Linux and Android don't support xattrs on symlinks according
		// to xattr(7), so just test that we get the proper error.
		if errno != syscall.EACCES && errno != syscall.EPERM {
			t.Fatal(error(errno))
		}

		test.WaitSignal(t, func() error {
			_, _, errno = syscall.Syscall(syscall.SYS_LREMOVEXATTR, uintptr(testFilePtr), uintptr(xattrNamePtr), 0)
			// Linux and Android don't support xattrs on symlinks according
			// to xattr(7), so just test that we get the proper error.
			if errno != syscall.EACCES && errno != syscall.EPERM {
				return error(errno)
			}
			return nil
		}, func(event *model.Event, _ *rules.Rule) {
			assert.Equal(t, "removexattr", event.GetType(), "wrong event type")
			assert.Equal(t, "user.test_xattr", event.RemoveXAttr.Name)
			assert.Equal(t, getInode(t, testFile), event.RemoveXAttr.File.Inode, "wrong inode")
			assertRights(t, event.RemoveXAttr.File.Mode, 0777)
			assertNearTime(t, event.RemoveXAttr.File.MTime)
			assertNearTime(t, event.RemoveXAttr.File.CTime)

			value, _ := event.GetFieldValue("event.async")
			assert.Equal(t, value.(bool), false)
		})
	})

	t.Run("fremovexattr", func(t *testing.T) {
		testFile, _, err := test.CreateWithOptions("test-removexattr", 98, 99, fileMode)
		if err != nil {
			t.Fatal(err)
		}

		f, err := os.Open(testFile)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		defer os.Remove(testFile)

		// set xattr
		_, _, errno := syscall.Syscall6(syscall.SYS_FSETXATTR, f.Fd(), uintptr(xattrNamePtr), 0, 0, 1, 0)
		if errno != 0 {
			t.Fatal(error(errno))
		}

		test.WaitSignal(t, func() error {
			_, _, errno = syscall.Syscall(syscall.SYS_FREMOVEXATTR, f.Fd(), uintptr(xattrNamePtr), 0)
			if errno != 0 {
				return error(errno)
			}
			return nil
		}, func(event *model.Event, _ *rules.Rule) {
			if event.GetType() != "removexattr" {
				t.Errorf("expected removexattr event, got %s", event.GetType())
			}

			if event.RemoveXAttr.Name != "user.test_xattr" || event.RemoveXAttr.Namespace != "user" {
				t.Errorf("expected removexattr name user.test_xattr, got %s", event.RemoveXAttr.Name)
			}

			if inode := getInode(t, testFile); inode != event.RemoveXAttr.File.Inode {
				t.Errorf("expected inode %d, got %d", event.RemoveXAttr.File.Inode, inode)
			}

			if int(event.RemoveXAttr.File.Mode)&expectedMode != expectedMode {
				t.Errorf("expected initial mode %d, got %d", expectedMode, int(event.RemoveXAttr.File.Mode)&expectedMode)
			}

			assertNearTime(t, event.RemoveXAttr.File.MTime)
			assertNearTime(t, event.RemoveXAttr.File.CTime)
		})
	})
}
