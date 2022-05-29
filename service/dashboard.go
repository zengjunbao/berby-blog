package service

import "github.com/cjyzwg/forestblog/config"

func SetThemeColor(index int)  {
	config.Cfg.ThemeColor = config.Cfg.ThemeOption[index]
}
