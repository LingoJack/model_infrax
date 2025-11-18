package tool

import (
	"fmt"
	"testing"
	"time"
)

func TestDeepCopyToSameType(t *testing.T) {
	type UserEntity struct {
		ID       *int64
		Name     *string
		Age      *int
		Birthday *time.Time
	}

	type UserVo struct {
		ID   int64
		Name string
	}

	user := UserEntity{
		ID:       Int64Ptr(1),
		Age:      IntPtr(18),
		Birthday: TimePtr(time.Now()),
	}

	userVo := UserVo{
		Name: "john",
	}

	user, err := DeepCopyToByJson[UserEntity](userVo)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(JsonifyIndent(user))
}
