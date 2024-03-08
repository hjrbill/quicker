package job

import (
	"encoding/csv"
	"github.com/gogo/protobuf/proto"
	model "github.com/hjrbill/quicker/demo/dal"
	"github.com/hjrbill/quicker/index_service"
	"github.com/hjrbill/quicker/pb"
	farmhash "github.com/leemcloughlin/gofarmhash"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// BuildIndexFromFile
// @param csvFile
// @param indexer 使用的 indexer 服务
// @param total 分布式 worker 总数，单体服务填 0
// @param workerIndex 本 worker 序号，单体服务填 0
func BuildIndexFromFile(csvFile string, indexer index_service.IIndexer, total, workerIndex int) {
	file, err := os.Open(csvFile)
	if err != nil {
		log.Printf("open csv file failed: %v", err)
		return
	}

	timeZone, err := time.LoadLocation("Asia/Shanghai") // 设置时区为东 8 区
	if err != nil {
		log.Printf("load time location failed: %v", err)
		return
	}

	reader := csv.NewReader(file)
	cnt := 0
	for {
		record, err := reader.Read() // record[0~9] 分别对应 video 的一个字段
		if err != nil {
			if err != io.EOF { // 如果不是因为是文件末尾，打日志说明读取失败
				log.Printf("read csv file failed: %v", err)
			}
			break
		}

		if len(record) < 10 { // 读取到的数据长度小于 10，为不完整数据，不做处理
			continue
		}

		if total > 0 && int(farmhash.Hash32([]byte(record[0])))%total != workerIndex { // 避免同一机器部署多个 worker 时对同文件同一数据的多次读取
			continue
		}

		BV := strings.TrimLeft(record[0], "https://www.bilibili.com/video/")
		video := &model.Video{
			ID:     BV,
			Title:  record[1],
			Author: record[3],
		}
		if len(record[2]) > 4 {
			postTime, err := time.ParseInLocation("2006/1/2 15:4", record[2], timeZone)
			if err == nil {
				video.PostTime = postTime.Unix()
			} else {
				log.Printf("parse time failed: %v", err)
			}
		}

		view, _ := strconv.ParseInt(record[4], 10, 32)
		video.View = int32(view)

		like, _ := strconv.ParseInt(record[5], 10, 32)
		video.Like = int32(like)

		coin, _ := strconv.ParseInt(record[6], 10, 32)
		video.Coin = int32(coin)

		favorite, _ := strconv.ParseInt(record[7], 10, 32)
		video.Favorite = int32(favorite)

		share, _ := strconv.ParseInt(record[8], 10, 32)
		video.Share = int32(share)

		keywords := strings.Split(record[9], ",")
		for _, kw := range keywords {
			kw = strings.TrimSpace(kw)
			if len(kw) > 0 {
				video.KeyWords = append(video.KeyWords, strings.ToLower(kw))
			}
		}

		err = AddVideoToIndex(video, indexer) // 将视频写入索引
		if err != nil {
			log.Printf("add video to index failed: %v", err)
			return
		}

		cnt++ // 计数
	}
	log.Printf("add %d documents to index successfully", cnt)
}

// AddVideoToIndex 将一个视频写入索引，如果存在就更新，如果不存在就创建，可用于实时更新
func AddVideoToIndex(video *model.Video, indexer index_service.IIndexer) error {
	doc := pb.Document{
		Id:          video.ID,
		BitsFeature: 0,
		Keywords:    nil,
	}

	doc.BitsFeature = GetClassBits(video.KeyWords) // 构建快速过滤使用的 类别 bit

	keywords := make([]*pb.Keyword, 0, len(video.KeyWords)+1)
	keywords = append(doc.Keywords, &pb.Keyword{
		Field: "author",                      // 来源：作者
		Word:  strings.ToLower(video.Author), // 关键词
	})
	for _, kw := range video.KeyWords {
		keywords = append(doc.Keywords, &pb.Keyword{
			Field: "content",           // 来源：视频内容
			Word:  strings.ToLower(kw), // 关键词
		})
	}
	doc.Keywords = keywords

	bytes, err := proto.Marshal(video)
	if err != nil {
		return err
	}
	doc.Bytes = bytes

	_, err = indexer.AddDoc(doc)
	if err != nil {
		return err
	}
	return nil
}
