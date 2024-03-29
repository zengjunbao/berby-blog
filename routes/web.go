package routes

import (
	"net/http"

	"github.com/cjyzwg/forestblog/config"
	"github.com/cjyzwg/forestblog/controller"
)

func initWebRoute() {

	http.HandleFunc("/", controller.Index)
	http.HandleFunc("/categories", controller.Categories)
	http.HandleFunc("/about", controller.About)
	http.HandleFunc("/savedocument", controller.HandleDocumentData)
	http.HandleFunc("/del", controller.HandleDelData)

	http.HandleFunc("/activity", controller.HandleActivity)

	// 二级页面
	http.HandleFunc("/article", controller.Article)
	http.HandleFunc("/category", controller.CategoryArticle)
	http.HandleFunc(config.Cfg.DashboardEntrance, controller.Dashboard)
	// 静态文件服务器
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("resources/public"))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(config.Cfg.DocumentPath+"/assets"))))

}
