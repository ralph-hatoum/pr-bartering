package failuresimulation

import (
	configextractor "bartering/config-extractor"
	"reflect"
	"testing"
)

func TestExtractFailureModelNodeProfile(t *testing.T) {

	config1 := configextractor.Config{FailureModel: "weibull"}
	config2 := configextractor.Config{FailureModel: "lognormal"}

	probabilityLaw1, _ := ExtractFailureModelNodeProfile(config1)
	probabilityLaw2, _ := ExtractFailureModelNodeProfile(config2)

	if reflect.ValueOf(probabilityLaw1).Pointer() != reflect.ValueOf(DrawNumberWeibull).Pointer() && reflect.ValueOf(probabilityLaw2).Pointer() != reflect.ValueOf(DrawNumberLognormal).Pointer() {
		t.Errorf("function not returning right probability law drawing function")
	}
}

func TestDrawNumberWeibull(t *testing.T) {

}

func TestDrawNumberLognormal(t *testing.T) {

}
