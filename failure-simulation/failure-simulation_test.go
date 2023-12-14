package failuresimulation

import (
	configextractor "bartering/config-extractor"
	"fmt"
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

	if reflect.ValueOf(probabilityLaw1).Pointer() != reflect.ValueOf(DrawNumberWeibull).Pointer() || reflect.ValueOf(probabilityLaw2).Pointer() != reflect.ValueOf(DrawNumberLognormal).Pointer() || reflect.ValueOf(probabilityLaw3).Pointer() != reflect.ValueOf(func(float64, float64) float64 { return 0.0 }).Pointer() {
		t.Errorf("function not returning right probability law drawing function")
	}
}

func TestDrawNumberWeibull(t *testing.T) {
	// not sure what to do here yet ...

}

func TestDrawNumberLognormal(t *testing.T) {
	// not sure what to do here yet ...
}

func TestExtractConnectivityFactor(t *testing.T) {
	config1 := configextractor.Config{NodeProfile: "peer"}
	config2 := configextractor.Config{NodeProfile: "peeper"}
	config3 := configextractor.Config{NodeProfile: "benefactor"}

	config4 := configextractor.Config{NodeProfile: "random"}

	cf1, _ := ExtractConnectivityFactor(config1)
	cf2, _ := ExtractConnectivityFactor(config2)
	cf3, _ := ExtractConnectivityFactor(config3)
	cf4, _ := ExtractConnectivityFactor(config4)

	fmt.Println(cf1, cf2, cf3, cf2, cf4)

}
