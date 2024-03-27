package param

type SearchRequest struct {
	Author   string   // 作者
	Classes  []string // 搜索类别，任一命中
	Keywords []string // 搜索关键词，全部命中
	ViewFrom int      //视频播放量下限
	ViewTo   int      //视频播放量上限
}
