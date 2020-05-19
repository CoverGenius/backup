package helpers

import (
	"testing"
	//	"os"
)

func TestReadFile(t *testing.T) {
	testfile := "../main.go"
	result := ReadFile(&testfile)
	if len(*result) < 1 {
		t.Errorf("Error reading file: %s", testfile)
	}
}

func TestIsFileExists(t *testing.T) {
	fake := "/sdfsadf/asdf12"
	existing := "/etc/crontab"

	if IsFileExists(&fake) == true {
		t.Errorf("File: %s should not exist", fake)
	}

	if IsFileExists(&existing) == false {
		t.Errorf("File: %s should exists on the system", existing)
	}
}

func TestIsCommandExists(t *testing.T) {
	fake := "asdf12"
	existing := "ls"

	if IsCommandExists(&fake) == true {
		t.Errorf("Command: %s should not exist", fake)
	}

	if IsCommandExists(&existing) == false {
		t.Errorf("Command: %s should exists on the system", existing)
	}
}

// TODO: Intercept log.Fatalf
//func TestRemoveDir(t *testing.T) {
//	test1 := "/tmp/test123"
//	test2 := "/var/../tmp/test123"
//	test3 := "/tmp"
//	test4 := "/etc/../tmp"
//	test5 := "/etc/../tmp/"
//	os.Mkdir(test1, 0700)
//	RemoveDir(&test1)
//	os.Mkdir(test2, 0700)
//	RemoveDir(&test2)
//	RemoveDir(&test3)
//	RemoveDir(&test4)
//	RemoveDir(&test5)
//}
