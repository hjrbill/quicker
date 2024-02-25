package test

import (
	"quicker/internal/kvdb"
	"quicker/pkg/util"
	"testing"
)

func TestBoltDB(t *testing.T) {
	setup = func() {
		var err error
		db, err = kvdb.NewKVDB(kvdb.BOLT, util.RootPath+"temp/test/bolt_db")
		if err != nil {
			panic(err)
		}
	}
	t.Run("test_bolt_db", testPipeline)
}
