package configextractor

import (
	"bartering/utils"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Port                                      int     `yaml:"Port"`
	TotalStorage                              int     `yaml:"TotalStorage"`
	BarteringInitialScore                     float64 `yaml:"BarteringInitialScore"`
	BarteringFactorAcceptableRatio            float64 `yaml:"BarteringFactorAcceptableRatio"`
	BarteringRatioIncreaseRate                float64 `yaml:"BarteringRatioIncreaseRate"`
	StoragerequestsScoreDecreaseRefusedStoReq float64 `yaml:"StoragerequestsScoreDecreaseRefusedStoReq"`
	StoragetestingTimerTimeoutSec             float64 `yaml:"StoragetestingTimerTimeoutSec"`
	StoragetestingTestingPeriod               float64 `yaml:"StoragetestingTestingPeriod"`
	StoragetestingFailedTestTimeoutDecrease   float64 `yaml:"StoragetestingFailedTestTimeoutDecrease"`
	StoragetestingFailedTestWrongAnsDecrease  float64 `yaml:"StoragetestingFailedTestWrongAnsDecrease"`
	StoragetestingPassedTestIncrease          float64 `yaml:"StoragetestingPassedTestIncrease"`
}

func ConfigExtractor(path string) Config {
	file, err := ioutil.ReadFile(path)
	utils.ErrorHandler(err)

	config := Config{}

	err = yaml.Unmarshal(file, &config)
	utils.ErrorHandler(err)

	return config
}
