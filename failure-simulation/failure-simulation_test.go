package failuresimulation

import (
	configextractor "bartering/config-extractor"
	"errors"

	"gonum.org/v1/gonum/stat/distuv"
)

func extractFailureModelNodeProfile(config configextractor.Config) (func(float64, float64) float64, error) {
	if config.FailureModel == "weibull" {
		return drawNumberWeibull, nil
	} else if config.FailureModel == "lognormal" {
		return drawNumberLognormal, nil
	} else {
		return func(float64, float64) float64 { return 0.0 }, errors.New("failure model not recognized")
	}
}

func drawNumberWeibull(shape float64, scale float64) float64 {
	return 0.0
}

func drawNumberLognormal(shape float64, scale float64) float64 {
	return 0.0
}

func Failure(config configextractor.Config) {

	if config.FailureModel == "weibull" {

	}

	shape := 2.0
	scale := 1.5

	weibullDist := distuv.Weibull{
		K:      shape,
		Lambda: scale,
	}

}
