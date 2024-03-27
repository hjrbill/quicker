package pb

func (k *Keyword) ToString() string {
	if len(k.Word) > 0 {
		return k.Field + "\001" + k.Word // 使用分隔符分割字段并拼接返回
	} else {
		return ""
	}
}
