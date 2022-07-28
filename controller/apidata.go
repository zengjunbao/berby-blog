package controller

import (
	"encoding/json"
	"fmt"
	"github.com/cjyzwg/forestblog/config"
	"github.com/cjyzwg/forestblog/models"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
)

func HandleDocumentData(w http.ResponseWriter, r *http.Request) {
	// r.ParseForm() //解析参数，默认是不会解析的
	fmt.Println("method is:" + r.Method)
	if r.Method == "POST" {
		mr, err := r.MultipartReader()
		if err != nil {
			fmt.Println("r.MultipartReader() err,", err)
			return
		}
		// r.Body.Close()
		form, _ := mr.ReadForm(128)
		m := make(map[string]string)
		for k, v := range form.Value {
			fmt.Println("value,k,v = ", k, ",", v)
			if k == "title" {
				m["title"] = v[0]
			}
			if k == "body" {
				m["body"] = v[0]
			}
			if k == "category" {
				m["category"] = v[0]
			}
		}
		blogPath := config.CurrentDir + "/" + config.Cfg.DocumentPath
		filename := blogPath + "/content/" + m["category"] + "/" + m["title"] + ".md"
		newfile, error := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0766)
		if error != nil {
			fmt.Println(error)
		}
		io.WriteString(newfile, m["body"])

		newfile.Close()

		_, err = exec.LookPath("git")

		if err != nil {
			fmt.Println("请先安装git并克隆博客文档到" + blogPath)

		}
		fmt.Println("blogpath is:" + blogPath)
		fmt.Println("filename is:" + filename)

		// cmd := exec.Command("git", "pull")
		// cmd.Dir = blogPath
		// _, err = cmd.CombinedOutput()
		// if err != nil {
		//	fmt.Println("cmd.Run() failed with", err)
		//	return
		// }

		cmd1 := exec.Command("git", "add", filename)
		cmd1.Dir = blogPath
		_, err = cmd1.CombinedOutput()
		if err != nil {
			fmt.Println("cmd.Run() failed with", err)
			return
		}
		cmd2 := exec.Command("git", "commit", "-m", "add")
		cmd2.Dir = blogPath
		_, err = cmd2.CombinedOutput()
		if err != nil {
			fmt.Println("cmd2.Run() failed with", err)
			return
		}
		cmd3 := exec.Command("git", "push", "origin", "master")
		cmd3.Dir = blogPath
		_, err = cmd3.CombinedOutput()
		if err != nil {
			fmt.Println("cmd3.Run() failed with", err)
			return
		}

		var categorylists = "1"
		jsoncategorylists, _ := json.Marshal(categorylists)
		fmt.Println(string(jsoncategorylists))
		// 返回的这个是给json用的，需要去掉
		// w.Header().Set("Content-Length", strconv.Itoa(len(jsoncategorylists)))
		w.Write(jsoncategorylists)
		return

	}
}

func HandleDelData(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // 解析参数，默认是不会解析的
	if r.Method == "GET" {
		queryForm, err := url.ParseQuery(r.URL.RawQuery)
		fmt.Println(queryForm)
		if err != nil && len(queryForm["category"]) == 0 && len(queryForm["name"]) == 0 {
			fmt.Println("query is wrong", err)
			return
		}
		category := queryForm["category"][0]
		fmt.Println(category)
		name := queryForm["name"][0]
		fmt.Println(name)
		blogPath := config.CurrentDir + "/" + config.Cfg.DocumentPath
		filename := blogPath + "/content/" + category + "/" + name
		fmt.Println(filename)
		cmd := exec.Command("rm", "-rf", filename)
		_, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Println("cmd.Run() failed with", err)
			return
		}

		cmd1 := exec.Command("git", "add", filename)
		cmd1.Dir = blogPath
		_, err = cmd1.CombinedOutput()
		if err != nil {
			fmt.Println("cmd.Run() failed with", err)
			return
		}
		cmd2 := exec.Command("git", "commit", "-m", "add")
		cmd2.Dir = blogPath
		_, err = cmd2.CombinedOutput()
		if err != nil {
			fmt.Println("cmd2.Run() failed with", err)
			return
		}
		cmd3 := exec.Command("git", "push", "origin", "master")
		cmd3.Dir = blogPath
		_, err = cmd3.CombinedOutput()
		if err != nil {
			fmt.Println("cmd3.Run() failed with", err)
			return
		}

		var categorylists = "1"
		jsoncategorylists, _ := json.Marshal(categorylists)
		fmt.Println(string(jsoncategorylists))
		// 返回的这个是给json用的，需要去掉
		w.Header().Set("Content-Length", strconv.Itoa(len(jsoncategorylists)))
		w.Write(jsoncategorylists)
		return
	}
}

type MarkdowndetailsResult struct {
	Code    int                    `json:"code"`
	Content models.MarkdownDetails `json:"content"`
}

type CategoryResult struct {
	Code    int               `json:"code"`
	Content models.Categories `json:"content"`
}

type ContentResult struct {
	Code    int                       `json:"code"`
	Content models.MarkdownPagination `json:"content"`
}

type Note struct {
	Title       string `json:"title"`
	Date        string `json:"date"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Content     string `json:"content"`
}
type Notes []Note

type CategoryList struct {
	Category string `json:"category"`
}
type CategoryLists []CategoryList
