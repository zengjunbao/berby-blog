## 快速开始

### 接口部分：
- 【文章列表】博客地址+/apis/articlelist，参数：page || search || 无
- 【文章内容】博客地址+/apis/articlecontent，参数：path(文件路径名，如：/DNS/a.md)
- 【文章分类】博客地址+/apis/category，参数：无
- 【文章分类的内容】博客地址+/apis/categorycontent，参数：name(分类名)

#### 地址：[cjblog](https://github.com/cjyzwg/cjblog)  

#### 运行（熟悉golang）
1. git clone http://github.com/cjyzwg/cjblog 
2. cd cjblog
3. go mod tidy(前提已经开启export GO111MODULE=on）
4. go run main.go打开浏览器，访问http://localhost:8081 即可
5. 删除resources/blog_docs/content中所有内容，分类名为文件夹名如：DNS，文件为名为：a.md 
6. 所有配置项在config.json中，均可以修改
```
{
  "port": 8081,
  "pageSize": 7,
  "descriptionLen": 200,
  "author": "Chensir",
  "webHookSecret": "cj",
  "timeLayout": "2006.01.02 15:04",
  "siteName": "Chensir's Personal Blog",
  "documentPath": "resources/blog_docs",
  "htmlKeywords": "forest blog,Golang,前端",
  "htmlDescription": "Chensir's Personal blog",
  "categoryListFileNumber": 6,
  "themeColor": "#9c27b0",
  "dashboardEntrance": "/admin",
  "themeOption": ["#673ab7","#f44336","#9c27b0","#2196f3","#607d8b","#795548"]
}
```
**注：请确保go版本在1.11以上**

#### 注意事项：
1. 所有的笔记都是markdown文件,.md结尾
2. 服务端代码中不涉及mysql部分，本着简洁的目的，通过生成的cache文件来访问，若要添加，可自行添加或者联系我
3. 可将此代码部署到服务器，或者部署到本机，服务器代理转发到本机即可


## 功能

- [x] 文章展示
- [x] 评论展示
- [x] 搜索文章功能
- [x] 文章评论功能
- [x] 评论功能内容识别
- [x] 友链展示
- [x] 点赞功能（云函数）
- [x] 文章浏览统计功能（云函数）
- [ ] 用户回复评论追评功能
- [ ] 生成海报
- [ ] 博主查看评论功能
- [ ] 博主回复评论功能
- [ ] 支付宝小程序抽屉功能修复...

## 感谢

UniBlog的诞生离不开下面这些项目：

- **[WeHalo](https://github.com/aquanlerou/WeHalo)： 一个简约清爽的开源博客微信小程序。**
- **[ColorUI](https://github.com/weilanwl/ColorUI)：鲜亮的高饱和色彩，专注视觉的小程序组件库**
- **[ForestBlog](https://github.com/xusenlin/forest-blog)：golang简约版的博客应用**
