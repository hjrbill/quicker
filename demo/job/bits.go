package job

import "slices"

// 视频类别枚举
const (
	ZI_XUN    = 1 << iota // 1 << 0
	SHE_HUI               // 1 << 1
	RE_DIAN               // 1 << 2
	SHENG_HUO             // 1 << 3
	ZHI_SHI
	HUAN_QIU
	YOU_XI
	ZONG_HE
	RI_CHANG
	YING_SHI
	DONG_HUA
	KE_JI
	YU_LE
	BIAN_CHENG
)

// GetClassBits 从 Keywords 中提取类型，构建快速过滤使用的 类别 bit
func GetClassBits(keywords []string) uint64 {
	var bits uint64
	if slices.Contains(keywords, "资讯") {
		bits |= ZI_XUN //属于哪个类别，将对应的 bit 置为 1
	}
	if slices.Contains(keywords, "社会") {
		bits |= SHE_HUI
	}
	if slices.Contains(keywords, "热点") {
		bits |= RE_DIAN
	}
	if slices.Contains(keywords, "生活") {
		bits |= SHENG_HUO
	}
	if slices.Contains(keywords, "知识") {
		bits |= ZHI_SHI
	}
	if slices.Contains(keywords, "环球") {
		bits |= HUAN_QIU
	}
	if slices.Contains(keywords, "游戏") {
		bits |= YOU_XI
	}
	if slices.Contains(keywords, "综合") {
		bits |= ZONG_HE
	}
	if slices.Contains(keywords, "日常") {
		bits |= RI_CHANG
	}
	if slices.Contains(keywords, "影视") {
		bits |= YING_SHI
	}
	if slices.Contains(keywords, "动画") {
		bits |= DONG_HUA
	}
	if slices.Contains(keywords, "科技") {
		bits |= KE_JI
	}
	if slices.Contains(keywords, "娱乐") {
		bits |= YU_LE
	}
	if slices.Contains(keywords, "编程") {
		bits |= BIAN_CHENG
	}
	return bits
}
