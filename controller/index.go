package controller

import (
	"bufio"
	"fmt"
	"github.com/cjyzwg/forestblog/config"
	"github.com/cjyzwg/forestblog/helper"
	"github.com/cjyzwg/forestblog/models"
	"github.com/cjyzwg/forestblog/service"
	"net/http"
	"os"
	"strconv"
)

func Index(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		helper.WriteErrorHtml(w, err.Error())
		return
	}

	page, err := strconv.Atoi(r.Form.Get("page"))
	if err != nil {
		page = 1
	}

	template, err := helper.HtmlTemplate("index")
	if err != nil {
		helper.WriteErrorHtml(w, err.Error())
		return
	}

	searchKey := r.Form.Get("search")

	markdownPagination, err := service.GetArticleList(page, "/", searchKey)
	if err != nil {
		helper.WriteErrorHtml(w, err.Error())
		return
	}

	err = template.Execute(w, map[string]interface{}{
		"Title":  "首页",
		"Data":   markdownPagination,
		"Config": config.Cfg,
	})

	if err != nil {
		helper.WriteErrorHtml(w, err.Error())
		return
	}

}

func Categories(w http.ResponseWriter, r *http.Request) {

	template, err := helper.HtmlTemplate("categories")

	if err != nil {
		helper.WriteErrorHtml(w, err.Error())
		return
	}

	categories, err := service.GetCategories()
	if err != nil {
		helper.WriteErrorHtml(w, err.Error())
		return
	}
	err = template.Execute(w, map[string]interface{}{
		"Title":  "分类",
		"Data":   categories,
		"Config": config.Cfg,
	})

	if err != nil {
		helper.WriteErrorHtml(w, err.Error())
		return
	}
}

func About(w http.ResponseWriter, r *http.Request) {

	template, err := helper.HtmlTemplate("about")
	if err != nil {
		helper.WriteErrorHtml(w, err.Error())
		return
	}

	markdown, err := models.ReadMarkdownBody("/About.md")

	if err != nil {
		helper.WriteErrorHtml(w, err.Error())
		return
	}

	err = template.Execute(w, map[string]interface{}{
		"Title":  "关于",
		"Data":   markdown,
		"Config": config.Cfg,
	})
	if err != nil {
		helper.WriteErrorHtml(w, err.Error())
		return
	}
}

func HandleActivity(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // 解析参数，默认是不会解析的
	if r.Method == "GET" {
		template, err := helper.HtmlTemplate("my_test")
		if err != nil {
			helper.WriteErrorHtml(w, err.Error())
			return
		}

		err = template.Execute(w, map[string]interface{}{
			"Title":  "关于",
			"Data":   nil,
			"Config": config.Cfg,
		})
		if err != nil {
			helper.WriteErrorHtml(w, err.Error())
			return
		}
		return
	} else if r.Method == "POST" {
		filePath := "activity.txt"
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("文件打开失败", err)
		}
		// 及时关闭file句柄
		defer file.Close()
		// 写入文件时，使用带缓存的 *Writer
		write := bufio.NewWriter(file)

		name := r.Form.Get("name")
		address := r.Form.Get("address")
		phone := r.Form.Get("phone")
		time := r.Form.Get("time")

		str := fmt.Sprintf("姓名：%v \n地址：%v \n电话：%v \n预约时间：%v \n",name,address,phone,time)
		write.WriteString(str)

		// Flush将缓存的文件真正写入到文件中
		write.Flush()
		helper.SedResponse(w, "预约成功")
		return
	}

}
