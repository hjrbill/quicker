package pb

import (
	"strings"
)

// NewTermQuery 创建叶节点 TermQuery（作为叶节点，无逻辑关系，只用于存储 keyword）
func NewTermQuery(field, word string) *TermQuery {
	if len(word) == 0 {
		return nil
	}
	keyword := &Keyword{
		Field: field,
		Word:  word,
	}
	return &TermQuery{Keyword: keyword}
}

func (t *TermQuery) Empty() bool {
	return t.Keyword == nil && len(t.Must) == 0 && len(t.Should) == 0
}

func (t *TermQuery) And(terms ...*TermQuery) *TermQuery {
	if len(terms) == 0 {
		return t
	}
	must := make([]*TermQuery, 0, len(terms)+1)
	if !t.Empty() {
		must = append(must, t)
	}
	for _, term := range terms {
		if !term.Empty() {
			must = append(must, term)
		}
	}
	return &TermQuery{Must: must}
}

func (t *TermQuery) Or(terms ...*TermQuery) *TermQuery {
	if len(terms) == 0 {
		return t
	}
	should := make([]*TermQuery, 0, len(terms)+1)
	if !t.Empty() {
		should = append(should, t)
	}
	for _, term := range terms {
		if !term.Empty() {
			should = append(should, term)
		}
	}
	return &TermQuery{Should: should}
}

func (t *TermQuery) ToString() string {
	if t.Keyword != nil {
		return t.Keyword.ToString()
	} else if len(t.Must) > 0 {
		if len(t.Must) == 1 {
			return t.Must[0].ToString()
		}
		sb := strings.Builder{}
		// 避免存在父节点，预先加上括号
		sb.WriteByte('(')
		// 递归遍历子树，并拼接逻辑符号
		for _, termQuery := range t.Must {
			s := termQuery.ToString()
			if len(s) > 0 {
				sb.WriteString(termQuery.ToString())
				sb.WriteString("&&")
			}
		}
		str := sb.String()
		// 去除末尾多拼接的逻辑符
		str = str[:len(str)-2]
		str = str + ")"
		return str
	} else if len(t.Should) > 0 {
		if len(t.Should) == 1 {
			return t.Should[0].ToString()
		}
		sb := strings.Builder{}
		sb.WriteByte('(')
		for _, termQuery := range t.Should {
			s := termQuery.ToString()
			if len(s) > 0 {
				sb.WriteString(termQuery.ToString())
				sb.WriteString("||")
			}
		}
		str := sb.String()
		str = str[:len(str)-2] + ")"
		return str
	}
	return ""
}
