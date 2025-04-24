package initialize

import (
	"license/jetbrains/util"
	"license/logger"
)

// InitJetbrains initialize JetBrains components
func InitJetbrains() {
	logger.Info("init fake cert")
	fakeCert := util.GetFake()
	fakeCert.LoadOrGenerate()
	err := fakeCert.LoadRootCert()
	if err != nil {
		logger.Error("load root ca err %e", err)
	}
	err = fakeCert.GenerateRootCert()
	if err != nil {
		logger.Error("generate jet ca err %e", err)
	}

	err = fakeCert.LoadCert()
	if err != nil {
		logger.Error("load my ca err %e", err)
	}
	logger.Info("init fake cert done")
}
