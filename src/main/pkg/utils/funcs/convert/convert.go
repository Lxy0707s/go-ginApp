package convert

import "strings"

var bonreeProv2BsyProv = make(map[string]string)
var bonreeISP2BsyISP = make(map[string]string)
var ipLibISP2BsyISP = make(map[string]string)
var region2Chiness = make(map[string]string)
var iplibISPChina2Ename = make(map[string]string)
var iplibCountryChinaToBsyEname = make(map[string]string)

func init() {
	bonreeProv2BsyProv["河南"] = "henan"
	bonreeProv2BsyProv["湖北"] = "hubei"
	bonreeProv2BsyProv["湖南"] = "hunan"
	bonreeProv2BsyProv["安徽"] = "anhui"
	bonreeProv2BsyProv["北京"] = "beijing"
	bonreeProv2BsyProv["福建"] = "fujian"
	bonreeProv2BsyProv["甘肃"] = "gansu"
	bonreeProv2BsyProv["广东"] = "guangdong"
	bonreeProv2BsyProv["山东"] = "shandong"
	bonreeProv2BsyProv["青海"] = "qinghai"
	bonreeProv2BsyProv["宁夏"] = "ningxia"
	bonreeProv2BsyProv["内蒙古"] = "neimenggu"
	bonreeProv2BsyProv["辽宁"] = "liaoning"
	bonreeProv2BsyProv["江西"] = "jiangxi"
	bonreeProv2BsyProv["江苏"] = "jiangsu"
	bonreeProv2BsyProv["吉林"] = "jilin"
	bonreeProv2BsyProv["黑龙江"] = "heilongjiang"
	bonreeProv2BsyProv["河北"] = "hebei"
	bonreeProv2BsyProv["海南"] = "hainan"
	bonreeProv2BsyProv["贵州"] = "guizhou"
	bonreeProv2BsyProv["广西"] = "guangxi"
	bonreeProv2BsyProv["重庆"] = "chongqing"
	bonreeProv2BsyProv["浙江"] = "zhejiang"
	bonreeProv2BsyProv["新疆"] = "xinjiang"
	bonreeProv2BsyProv["香港"] = "xianggang"
	bonreeProv2BsyProv["西藏"] = "xizang"
	bonreeProv2BsyProv["天津"] = "tianjin"
	bonreeProv2BsyProv["台湾"] = "taiwan"
	bonreeProv2BsyProv["四川"] = "sichuan"
	bonreeProv2BsyProv["上海"] = "shanghai"
	bonreeProv2BsyProv["陕西"] = "shan3xi"
	bonreeProv2BsyProv["山西"] = "shanxi"
	bonreeProv2BsyProv["云南"] = "yunnan"

	bonreeISP2BsyISP["中国移动"] = "yd"
	bonreeISP2BsyISP["北京市歌华宽带"] = "gh"
	bonreeISP2BsyISP["广电宽带"] = "gd"
	bonreeISP2BsyISP["华数宽带"] = "hs"
	bonreeISP2BsyISP["中国电信"] = "dx"
	bonreeISP2BsyISP["中国教育网"] = "jyw"
	bonreeISP2BsyISP["中国联通"] = "lt"
	bonreeISP2BsyISP["中国铁通"] = "tt"

	ipLibISP2BsyISP["chinatelecom"] = "dx"
	ipLibISP2BsyISP["chinaunicom"] = "lt"
	ipLibISP2BsyISP["chinamobile"] = "yd"
	ipLibISP2BsyISP["chinarailcom"] = "tt"
	ipLibISP2BsyISP["chinaedu"] = "jyw"
	ipLibISP2BsyISP["wasu"] = "hs"
	ipLibISP2BsyISP["catv"] = "gd"

	region2Chiness["huadong"] = "华东"
	region2Chiness["huabei"] = "华北"
	region2Chiness["huazhong"] = "华中"
	region2Chiness["huanan"] = "华南"
	region2Chiness["xinan"] = "西南"
	region2Chiness["xibei"] = "西北"
	region2Chiness["dongbei"] = "东北"

	iplibISPChina2Ename["电信"] = "dx"
	iplibISPChina2Ename["联通"] = "lt"
	iplibISPChina2Ename["移动"] = "yd"
	iplibISPChina2Ename["铁通"] = "tt"

	iplibCountryChinaToBsyEname["中国"] = "China"
	iplibCountryChinaToBsyEname["俄罗斯"] = "Russia"
	iplibCountryChinaToBsyEname["印度"] = "India"
}

// RegionFromBsyToChineseFunc convert region name from bsy to chinese
func RegionFromBsyToChineseFunc(region string) string {

	if len(region2Chiness[region]) > 0 {
		return region2Chiness[region]
	}
	return region
}

//ProvFromBonreeToBsy convert prov
func ProvFromBonreeToBsy(cProv string) string {
	if len(bonreeProv2BsyProv[cProv]) > 0 {
		return bonreeProv2BsyProv[cProv]
	}
	// return "unknow"
	return cProv
}

//ProvFromIPlibChinaToBsy convert prov
func ProvFromIPlibChinaToBsy(cProv string) string {
	return ProvFromBonreeToBsy(cProv)
}

//ISPFromBonreeToBsy convert isp from bonree to bsy
func ISPFromBonreeToBsy(cISP string) string {
	if len(bonreeISP2BsyISP[cISP]) > 0 {
		return bonreeISP2BsyISP[cISP]
	}
	return "unknow"
}

//ISPFromIPLibChinaToBsy convert isp from bonree to bsy
func ISPFromIPLibChinaToBsy(cISP string) string {
	if len(iplibISPChina2Ename[cISP]) > 0 {
		return iplibISPChina2Ename[cISP]
	}
	// return "unknow"
	return cISP
}

func CountryFromIPLibChinaToBsy(cCountry string) string {
	if len(iplibCountryChinaToBsyEname[cCountry]) > 0 {
		return iplibCountryChinaToBsyEname[cCountry]
	}
	return cCountry
}

// ISPFromIPLibToBonree convert isp
func ISPFromIPLibToBonree(isp string) string {
	isp = ISPFromIPLibToBsy(isp)
	for bonreeISP, bsyISP := range bonreeISP2BsyISP {
		if bsyISP == isp {
			return bonreeISP
		}
	}
	return "unknow"
}

//ISPFromIPLibToBsy convert isp
func ISPFromIPLibToBsy(ipLibISP string) string {
	ipLibISP = strings.ToLower(ipLibISP)
	if len(ipLibISP2BsyISP[ipLibISP]) > 0 {
		return ipLibISP2BsyISP[ipLibISP]
	}
	return "unknow"
}

//ProvFromIPLibToBsy convert province from iplib to bsy
func ProvFromIPLibToBsy(prov string) string {
	prov = strings.ToLower(prov)
	if prov == "shaanxi" {
		return "shan3xi"
	}
	if prov == "nei mongol" {
		return "neimenggu"
	}
	return prov
}

//ProvFromIPLibToBonree convert province from iplib to bonree
func ProvFromIPLibToBonree(prov string) string {
	prov = ProvFromIPLibToBsy(prov)
	for bonreeProv, bsyProv := range bonreeProv2BsyProv {
		if bsyProv == prov {
			return bonreeProv
		}
	}
	return "unknow"
}
