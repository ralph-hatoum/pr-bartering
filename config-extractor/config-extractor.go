package configextractor

import (
	"bartering/utils"
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

/*
	Package to read and parse the config file into a Config struct
*/

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
	FailureModel                              string  `yaml:"FailureModel"`
	NodeProfile                               string  `yaml:"NodeProfile"`
}

func ConfigExtractor(path string) Config {

	/*
		To read the config yaml file and extract the config
		Input : path to the config.yaml file
		Output : Config object
	*/
	file, err := ioutil.ReadFile(path)
	utils.ErrorHandler(err)

	config := Config{}

	err = yaml.Unmarshal(file, &config)
	utils.ErrorHandler(err)

	return config
}

func ConfigPrinter(conf Config) {

	/*
		To print a config object in a neat format
		Input : config object
	*/

	toPrint := fmt.Sprintf(`
	Read config -- launching node with the following parameters :
	Port : %d
	Node total storage : %d
	Initial scores attributed to peers : %f
	Maximum acceptable ratio factor : %f
	Ratio increase rate : %f
	Score decrease upon refused storage request : %f
	Proof of storage timeout in seconds : %f
	Storage testing period in seconds : %f
	Score decrease upon failed test in case of timeout : %f
	Score decrease upon failed test in case of wrong answer : %f
	Score increase upon succesful test : %f
	Failure mode : %s
	Node profile : %s
	`, conf.Port,
		conf.TotalStorage,
		conf.BarteringInitialScore,
		conf.BarteringFactorAcceptableRatio,
		conf.BarteringRatioIncreaseRate,
		conf.StoragerequestsScoreDecreaseRefusedStoReq,
		conf.StoragetestingTimerTimeoutSec,
		conf.StoragetestingTestingPeriod,
		conf.StoragetestingFailedTestTimeoutDecrease,
		conf.StoragetestingFailedTestWrongAnsDecrease,
		conf.StoragetestingPassedTestIncrease,
		conf.FailureModel,
		conf.NodeProfile)

	fmt.Println(toPrint)
}
