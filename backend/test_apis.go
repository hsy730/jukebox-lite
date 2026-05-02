//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func testAPI(name, url string) {
	fmt.Printf("\n========== %s ==========\n", name)
	fmt.Printf("URL: %s\n", url)
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %d\n", resp.StatusCode)
	
	if resp.StatusCode == 200 {
		var prettyJSON map[string]interface{}
		if err := json.Unmarshal(body, &prettyJSON); err == nil {
			formatted, _ := json.MarshalIndent(prettyJSON, "", "  ")
			fmt.Println(string(formatted))
		} else {
			fmt.Println(string(body))
		}
	} else {
		fmt.Println(string(body))
	}
}

func main() {
	apis := []struct {
		name string
		url  string
	}{
		{"Lokua-Search", "https://api.lokua.cn/music/search?keyword=周杰伦&source=netease&page=1&limit=3"},
		{"Lokua-Toplists", "https://api.lokua.cn/music/toplists?source=netease"},
		{"Lokua-Toplist", "https://api.lokua.cn/music/toplist?id=3778678&source=netease&page=1&limit=3"},
		{"Lokua-Info", "https://api.lokua.cn/music/info?id=001BLpXF2DyJe2&source=netease"},
	}
	
	for _, api := range apis {
		testAPI(api.name, api.url)
	}
}
