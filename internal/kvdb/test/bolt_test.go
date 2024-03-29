package test

import (
	"github.com/hjrbill/quicker/internal/kvdb"
	"github.com/hjrbill/quicker/pkg/path"
	"testing"
)

func TestBoltDB(t *testing.T) {
	setup = func() {
		var err error
		db, err = kvdb.NewKVDB(kvdb.BOLT, path.RootPath+"temp/test/bolt_db")
		if err != nil {
			panic(err)
		}
	}
	t.Run("test_bolt_db", testPipeline)
}
