package failuresimulation

import (
	configextractor "bartering/config-extractor"
	"reflect"
	"testing"
)

func TestExtractFailureModelNodeProfile(t *testing.T) {

	config1 := configextractor.Config{FailureModel: "weibull"}
	config2 := configextractor.Config{FailureModel: "lognormal"}

	config3 := configextractor.Config{FailureModel: "random"}

	probabilityLaw1, _ := ExtractFailureModelNodeProfile(config1)
	probabilityLaw2, _ := ExtractFailureModelNodeProfile(config2)
	probabilityLaw3, _ := ExtractFailureModelNodeProfile(config3)

	if reflect.ValueOf(probabilityLaw1).Pointer() != reflect.ValueOf(DrawNumberWeibull).Pointer() && reflect.ValueOf(probabilityLaw2).Pointer() != reflect.ValueOf(DrawNumberLognormal).Pointer() && reflect.ValueOf(probabilityLaw3).Pointer() != reflect.ValueOf(func(float64, float64) float64 { return 0.0 }).Pointer() {
		t.Errorf("function not returning right probability law drawing function")
	}
}

func TestDrawNumberWeibull(t *testing.T) {

}

func TestDrawNumberLognormal(t *testing.T) {

}
