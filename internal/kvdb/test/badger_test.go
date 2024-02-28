package test

import (
	"github.com/hjrbill/quicker/internal/kvdb"
	"github.com/hjrbill/quicker/pkg/util"
	"testing"
)

func TestBadgerTest(t *testing.T) {
	setup = func() {
		var err error
		db, err = kvdb.NewKVDB(kvdb.BADGER, util.RootPath+"temp/test/badger_db")
		if err != nil {
			panic(err)
		}
	}

	t.Run("test_badger_db", testPipeline)
}
