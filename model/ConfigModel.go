package model

type ConfigModel struct {
	AsteroidDataDir string      `json:"asteroidDataDir" yaml:"asteroid-data-dir"`
	WatchInterval   int64       `json:"monitorInterval" yaml:"watch-interval"`
	SiteList        []SiteModel `json:"siteList" yaml:"site-list"`
}

type SiteModel struct {
	SiteName   string   `json:"siteName" yaml:"site-name"`
	SiteDir    string   `json:"siteDir" yaml:"site-dir"`
	IncludeExt []string `json:"includeExt" yaml:"include-ext"`
	ExcludeDir []string `json:"excludeDir" yaml:"exclude-dir"`
}
