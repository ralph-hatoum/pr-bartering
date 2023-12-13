package failuresimulation

import (
	configextractor "bartering/config-extractor"
	"testing"
)

func TestExtractFailureModelNodeProfile(t *testing.T) {

	config1 := configextractor.Config{FailureModel: "weibull"}

	probabilityLaw, err := ExtractFailureModelNodeProfile(config1)
}

func TestDrawNumberWeibull(t *testing.T) {

}

func TestDrawNumberLognormal(t *testing.T) {

}
