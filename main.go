package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
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
)

//go:embed misc/description.txt
var projectDescription string

const noticeTitleLen = 50

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

	args := flag.Args()
	if len(args) > 0 && args[0] == "version" {
		fmt.Println(utils.Pur("v0.1"))
		pressAnyKeyToContinue()
	}

	fmt.Println(utils.Pur(projectDescription))

	if len(act) == 0 {
		log.Println(utils.Fata("missing param --act, usage:"), actSupported)
		pressAnyKeyToContinue()
	}
	if len(site) == 0 {
		log.Println(utils.Fata("missing param --site usage:"), actSupported)
		pressAnyKeyToContinue()
	}

	if err := global.InitProjDir(); err != nil {
		log.Println("global.InitProjDir:", utils.Fata(err))
		pressAnyKeyToContinue()
	}

	if err := global.InitLog(); err != nil {
		log.Println("global.InitLog:", utils.Fata(err))
		pressAnyKeyToContinue()
	}

	if err := global.InitConfig(); err != nil {
		log.Println("global.InitConfig:", utils.Fata(err))
		pressAnyKeyToContinue()
	}

	if err := global.InitDB(); err != nil {
		log.Println("global.InitDB:", utils.Fata(err))
		pressAnyKeyToContinue()
	}

	db, err := global.GormDB.DB()
	if err != nil {
		log.Println(utils.Fata(err))
		pressAnyKeyToContinue()
	}
	defer db.Close()

	myNotify.InitConfig(global.ConfigFile)

	siteList, err := getMatchSites(string(site))
	if err != nil {
		log.Println(utils.Fata(err))
		pressAnyKeyToContinue()
	}

	for _, siteModel := range siteList {
		siteDir := siteModel.SiteDir
		if strings.TrimSpace(siteDir) == "" {
			log.Println(utils.Fata("siteDir cannot be empty."))
			pressAnyKeyToContinue()
		}

		if !utils.ExistsDir(siteDir) {
			log.Println(utils.Fata("siteDir is not a dir."))
			pressAnyKeyToContinue()
		}

		siteDirName, err := utils.GetLastDirName(siteDir, true)
		if err != nil {
			log.Println(utils.Fata("GetLastDirName err:", err.Error()))
			pressAnyKeyToContinue()
		}

		go grt(act, siteModel, actSupported, siteDir, siteDirName)
	}

	for {
		select {
		case info, ok := <-global.BroadcastInfoList:
			if ok {
				log.Println(info)
				continue
			}
			fmt.Println("info chan NOT OK")
		case info, ok := <-global.BroadcastHighlightInfoList:
			if ok {
				log.Println(utils.Tea(info))
				continue
			}
			fmt.Println("highlight info chan NOT OK")

		case err, ok := <-global.BroadcastErrorList:
			if ok {
				log.Println(utils.Fata(err))
				continue
			}
			fmt.Println("err chan NOT OK")

		case notice, ok := <-global.BroadcastNoticeList:
			if ok {
				log.Println(utils.Warn(notice))

				title := notice
				if len(title) > noticeTitleLen {
					title = title[0:noticeTitleLen]
				}
				go myNotify.NotifyAll("File modified ! "+title, notice, "")
				continue
			}
			fmt.Println("notice chan NOT OK")

		default:
			fmt.Print(".")
		}

		/*		if act != Watch {
				time.Sleep(time.Second * 5)
				break
			}*/

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

		pressAnyKeyToContinue()
		return
	}

	if act == Watch {
		global.InfoToChan(fmt.Sprintf("[ %s ] under watching ...", siteDir))

		seconds := global.Config.MonitorInterval
		global.InfoToChan(fmt.Sprintf("sleep %d second(s) ...\n", seconds))
		go utils.Countdown(seconds)
		time.Sleep(time.Duration(seconds) * time.Second)

		global.LogDate = time.Now()
		goto _____monitor
	}

	global.InfoHighlightToChan(siteModel.SiteDir + " OK.")
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

func pressAnyKeyToContinue() {
	fmt.Println("Press any key to exit.")
	var input string
	fmt.Scanln(&input)
	os.Exit(0)
}
