package jobs

import (
	"github.com/mt1976/frantic-core/commonConfig"
	logger "github.com/mt1976/frantic-core/logHandler"
	trnsl8r "github.com/mt1976/trnsl8r_connect"
)

var cfg *commonConfig.Settings
var appName string
var translationServiceRequest trnsl8r.Request

func init() {
	cfg = commonConfig.Get()
	appName = cfg.GetApplication_Name()
	trnsServerProtocol := cfg.GetTranslationServer_Protocol()
	trnsServerHost := cfg.GetTranslationServer_Host()
	trnsServerPort := cfg.GetTranslationServer_Port()
	trnsLocale := cfg.GetTranslation_Locale()
	err := error(nil)
	translationServiceRequest, err = trnsl8r.NewRequest().FromOrigin(appName).WithHost(trnsServerHost).WithPort(trnsServerPort).WithProtocol(trnsServerProtocol).WithLogger(logger.TranslationLogger).WithFilter(trnsl8r.LOCALE, trnsLocale)
	if err != nil {
		panic(err)
	}
}
