package failuresimulation

import (
	configextractor "bartering/config-extractor"
	"errors"

	"gonum.org/v1/gonum/stat/distuv"
)

func ExtractFailureModelNodeProfile(config configextractor.Config) (func(float64, float64) float64, error) {
	if config.FailureModel == "weibull" {
		return DrawNumberWeibull, nil
	} else if config.FailureModel == "lognormal" {
		return DrawNumberLognormal, nil
	} else {
		return func(float64, float64) float64 { return 0.0 }, errors.New("failure model not recognized")
	}
}

func DrawNumberWeibull(shape float64, scale float64) float64 {
	weibullDist := distuv.Weibull{
		K:      shape,
		Lambda: scale,
	}
	return weibullDist.Rand()
}

func DrawNumberLognormal(Mu float64, Sigma float64) float64 {
	logNormalDist := distuv.LogNormal{
		Mu:    Mu,
		Sigma: Sigma,
	}
	return logNormalDist.Rand()
}

func Failure(config configextractor.Config, shape float64, scale float64) {

	// sessionLengthDraw, err := extractFailureModelNodeProfile(config)

	// utils.ErrorHandler(err)

}
