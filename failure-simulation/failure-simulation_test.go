package failuresimulation

import (
	configextractor "bartering/config-extractor"
	"fmt"
	"reflect"
	"testing"
)

func TestExtractFailureModel(t *testing.T) {

	config1 := configextractor.Config{FailureModel: "weibull"}
	config2 := configextractor.Config{FailureModel: "lognormal"}

	config3 := configextractor.Config{FailureModel: "random"}

	probabilityLaw1, _ := ExtractFailureModel(config1)
	probabilityLaw2, _ := ExtractFailureModel(config2)
	probabilityLaw3, _ := ExtractFailureModel(config3)

	if reflect.ValueOf(probabilityLaw1).Pointer() != reflect.ValueOf(DrawNumberWeibull).Pointer() || reflect.ValueOf(probabilityLaw2).Pointer() != reflect.ValueOf(DrawNumberLognormal).Pointer() || reflect.ValueOf(probabilityLaw3).Pointer() != reflect.ValueOf(func(float64, float64) float64 { return 0.0 }).Pointer() {
		t.Errorf("function not returning right probability law drawing function")
	}
}

func TestDrawNumberWeibull(t *testing.T) {
	// not sure what to do here yet ...
	sessionLength := DrawNumberWeibull(1.5, 100)
	fmt.Println(sessionLength)

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

	if cf1 != 0.5 || cf2 != 0.3 || cf3 != 0.7 || cf4 != 0.0 {
		t.Errorf("Connectivity not extracted well - should get 0.5 for peer, 0.3 for peeper and 0.7 for benefactor")
	}

}

func TestFailure(t *testing.T) {
	// config := configextractor.Config{NodeProfile: "benefactor", FailureModel: "weibull"}
	// Failure(config, 1.0, 100)
}
