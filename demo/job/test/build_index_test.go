package test

import (
	"github.com/hjrbill/quicker/demo/job"
	"github.com/hjrbill/quicker/index_service"
	"github.com/hjrbill/quicker/internal/kvdb"
	"github.com/hjrbill/quicker/pkg/path"
	"os"
	"testing"
)

var (
	dbType  = kvdb.BOLT
	dbPath  = path.RootPath + "temp/test/bolt_db"
	indexer *index_service.Indexer
)

func Init() {
	os.Remove(dbPath) //先删除原有的索引文件
	indexer = new(index_service.Indexer)
	if err := indexer.Init(50000, dbType, dbPath); err != nil {
		panic(err)
	}
}

func TestBuildIndexFromFile(t *testing.T) {
	Init()
	defer indexer.Close()
	csvFile := path.RootPath + "demo/deployments/video.csv"
	job.BuildIndexFromFile(csvFile, indexer, 0, 0)
}
