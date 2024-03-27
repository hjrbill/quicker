package test

import (
	"bytes"
	"fmt"
	"github.com/bytedance/sonic"
	model "github.com/hjrbill/quicker/demo/gen/video"
	"github.com/hjrbill/quicker/demo/param"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestSearch(t *testing.T) {
	client := http.Client{
		Timeout: 100 * time.Millisecond,
	}
	request := param.SearchRequest{
		Keywords: []string{"go", "gin"},
		Classes:  []string{"科技", "编程"},
		ViewFrom: 1000, //播放量大于 1000
		ViewTo:   0,    //播放量不设上限
	}
	bs, _ := sonic.Marshal(request)
	resp, err := client.Post("http://127.0.0.1:5678/search", "application/json", bytes.NewReader(bs))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	defer resp.Body.Close()
	content, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		var datas []model.Video
		sonic.Unmarshal(content, &datas)
		for _, data := range datas {
			fmt.Printf("%s %d %s %s\n", data.ID, data.View, data.Title, strings.Join(data.KeyWords, "|"))
		}
	} else {
		fmt.Println(resp.Status)
		t.Fail()
	}
}
