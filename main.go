package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/qcozof/asteroid/global"
	"github.com/qcozof/asteroid/model"
	"github.com/qcozof/asteroid/service"
	"github.com/qcozof/asteroid/utils"

	myNotify "github.com/qcozof/my-notify/notify"
)

type ActionType string

const (
	Init      ActionType = "init"
	Watch     ActionType = "watch"
	Uninstall ActionType = "uninstall"
	Reset     ActionType = "reset"
)

const noticeTitleLen = 50

//go:embed build/description.txt
var projectDescription string

//go:embed build/file_list.sql
var initSql string

var commandUtils utils.CommandUtils

func main() {
	var act ActionType
	var site ActionType
	var actSupported = fmt.Sprintf("--act [%s, %s, %s] --site [all, siteName1, siteName2, siteName3]", Init, Watch, Uninstall)

	// ./asteroid --act init --site all [or site1|site2...]
	// ./asteroid --act watch --site all [or site1|site2...]
	// ./asteroid --act uninstall --site all [or site1|site2...]
	flag.StringVar((*string)(&act), "act", "", actSupported)
	flag.StringVar((*string)(&site), "site", "", actSupported)
	flag.Parse()

	utils.Teal(projectDescription+"%s", "")

	if err := global.InitProjDir(); err != nil {
		utils.Fata("global.InitProjDir:%v", err)
		commandUtils.PressAnyKeyToContinue()
	}

	if err := global.InitLog(true); err != nil {
		utils.Fata("global.InitLog:%v", err)
		commandUtils.PressAnyKeyToContinue()
	}

	log.Println("program started !")

	if len(act) == 0 {
		utils.Fata("missing param --act, usage:%s", actSupported)
		commandUtils.PressAnyKeyToContinue()
	}

	if act != Reset && len(site) == 0 {
		utils.Fata("missing param --site usage:%s", actSupported)
		commandUtils.PressAnyKeyToContinue()
	}

	if err := global.InitConfig(); err != nil {
		utils.Fata("global.InitConfig:", err)
		commandUtils.PressAnyKeyToContinue()
	}

	if err := global.InitDB(); err != nil {
		utils.Fata("global.InitDB:%v", err)
		commandUtils.PressAnyKeyToContinue()
	}

	db, err := global.GormDB.DB()
	if err != nil {
		utils.Fata(err.Error())
		commandUtils.PressAnyKeyToContinue()
	}
	defer db.Close()

	myNotify.InitConfig(global.ConfigFile)
	go global.InitLog(false)

	if act == Reset {
		if err = service.Reset(initSql); err != nil {
			utils.Fata(err.Error())
		} else {
			utils.OK("Reset successfully !")
		}

		commandUtils.PressAnyKeyToContinue()
	}

	siteList, err := getMatchSites(string(site))
	if err != nil {
		utils.Fata(err.Error())
		commandUtils.PressAnyKeyToContinue()
	}

	siteCount := len(siteList)

	for _, siteModel := range siteList {
		siteDir := strings.TrimSpace(siteModel.SiteDir)
		if siteDir == "" {

			utils.Fata("siteDir cannot be empty.")
			commandUtils.PressAnyKeyToContinue()
		}

		for _, ex := range siteModel.ExcludeDir {
			if !strings.Contains(ex, siteDir) {
				utils.Warn("Please verify the configuration file to ensure that the 'exclude-dir' under the [%s] "+
					"is a subdirectory of the 'site-dir'.", siteModel.SiteName)
			}
		}

		if !utils.ExistsDir(siteDir) {
			utils.Fata("siteDir is not a dir.")
			commandUtils.PressAnyKeyToContinue()
		}

		siteDirName, err := utils.GetLastDirName(siteDir, true)
		if err != nil {
			utils.Fata("GetLastDirName err:%v", err.Error())
			commandUtils.PressAnyKeyToContinue()
		}

		go grt(act, siteModel, actSupported, siteDir, siteDirName)
	}

	actDoCount := 0
	for {
		select {
		case info, ok := <-global.BroadcastInfoList:
			if ok {
				log.Println(info)
				continue
			}
			utils.Warn("info chan NOT OK")
		case info, ok := <-global.BroadcastHighlightInfoList:
			if ok {
				if info == string(act) {
					actDoCount++
				} else {
					utils.Warn(info)
					log.Println(info)
				}

				if actDoCount == siteCount {
					utils.OK("All sites have been %sed successfully !", act)
					commandUtils.PressAnyKeyToContinue()
					break
				}

				continue
			}
			utils.Warn("highlight info chan NOT OK")

		case errMsg, ok := <-global.BroadcastErrorList:
			if ok {
				utils.Fata(errMsg)
				continue
			}
			utils.Warn("errMsg chan NOT OK")

		case notice, ok := <-global.BroadcastNoticeList:
			if ok {
				utils.Warn(notice)
				log.Println(notice)

				title := notice
				if len(title) > noticeTitleLen {
					title = title[0:noticeTitleLen]
				}
				go myNotify.NotifyAll("File modified ! "+title, notice, "")
				continue
			}
			utils.Warn("notice chan NOT OK")

		default:
			fmt.Print(".")
		}

		time.Sleep(time.Second)
	}
}

func grt(act ActionType, siteModel model.SiteModel, actSupported, siteDir, siteDirName string) {
_____monitor:
	var err error
	var tips string

	switch act {
	case Init:
		err = service.InitService(siteModel, siteDirName)
		tips = "InitService err:"

	case Watch:
		if !utils.ExistsDir(global.RepositoryDir) {
			err = errors.New("Please run init first. ")
			break
		}

		err = service.MonitorService(siteModel, siteDir, siteDirName)
		tips = "MonitorService err:"

	case Uninstall:
		err = service.Uninstall(siteDirName)
		tips = "Uninstall err:"

	default:
		err = errors.New(fmt.Sprintf("act:%s not match ! \nPlease use %s \n", act, actSupported))
	}

	if err != nil {
		global.ErrorToChan(tips, err)

		commandUtils.PressAnyKeyToContinue()
		return
	}

	if act == Watch {
		global.InfoToChan(fmt.Sprintf(`"%s" under watching ...`, siteDir))

		seconds := global.Config.WatchInterval
		//global.InfoToChan(fmt.Sprintf("sleep %d second(s) ...\n", seconds))
		go utils.Countdown(seconds)
		time.Sleep(time.Duration(seconds) * time.Second)

		goto _____monitor
	}

	global.InfoHighlightToChan(fmt.Sprintf("%s %v", siteModel.SiteDir, act))
	global.InfoHighlightToChan(act)
}

func getMatchSites(siteNameStr string) ([]model.SiteModel, error) {
	var siteList []model.SiteModel
	tmpSiteName := strings.TrimSpace(siteNameStr)
	if tmpSiteName == "" {
		return nil, errors.New("siteName cannot be empty")
	}

	if strings.ToLower(tmpSiteName) == "all" {
		return global.Config.SiteList, nil
	}

	sliSite := strings.Split(tmpSiteName, "|")
	for _, m := range global.Config.SiteList {
		for _, sn := range sliSite {
			if strings.TrimSpace(sn) == strings.TrimSpace(m.SiteName) {
				siteList = append(siteList, m)
			}
		}
	}

	if len(siteList) == 0 {
		return nil, errors.New("no siteNameStr match")
	}

	return siteList, nil
}
