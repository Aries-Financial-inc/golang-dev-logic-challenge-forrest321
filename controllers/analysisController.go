package controllers

import (
	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-forrest321/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

var long = models.Long
var short = models.Short
var call = models.Call
var put = models.Put

// AnalysisResponse represents the data structure of the analysis result
type AnalysisResponse struct {
	XYValues        []XYValue `json:"xy_values"`
	MaxProfit       float64   `json:"max_profit"`
	MaxLoss         float64   `json:"max_loss"`
	BreakEvenPoints []float64 `json:"break_even_points"`
}

// XYValue represents a pair of X and Y values
type XYValue struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// AnalysisHandler handles the request and responses
func AnalysisHandler(c *gin.Context) {
	var contracts []models.OptionsContract
	//Check request
	if err := c.ShouldBindJSON(&contracts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Ensure min / max contracts
	if len(contracts) == 0 || len(contracts) > 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Must contain between 1 and 4 options contracts"})
		return
	}

	var errs []error
	for _, ct := range contracts {
		err := ct.Validate()
		if len(err) > 0 {
			errs = append(errs, err...)
		}
	}

	if len(errs) > 0 {
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	response := AnalysisResponse{
		XYValues:        calculateXYValues(contracts),
		MaxProfit:       calculateMaxProfit(contracts),
		MaxLoss:         calculateMaxLoss(contracts),
		BreakEvenPoints: calculateBreakEvenPoints(contracts),
	}

	c.JSON(http.StatusOK, response)
}

// calculateXYValues calculates graph data for the contracts
func calculateXYValues(contracts []models.OptionsContract) []XYValue {
	var xyValues []XYValue
	var x, y float64

	for _, contract := range contracts {
		x = contract.StrikePrice
		if contract.GetType() == call {
			y = contract.Bid
		} else if contract.GetType() == put {
			y = contract.Ask
		}
		xyValues = append(xyValues, XYValue{X: x, Y: y})
	}

	return xyValues
}

// calculateMaxProfit calculates the maximum profit for the contracts
func calculateMaxProfit(contracts []models.OptionsContract) float64 {
	var profit, maxProfit float64

	for _, contract := range contracts {
		if contract.GetPosition() == long {
			profit = calculateLongProfit(contract)
		} else if contract.GetPosition() == short {
			profit = calculateShortProfit(contract)
		}
		if profit > maxProfit {
			maxProfit = profit
		}
	}

	return maxProfit
}

// calculateMaxLoss finds the maximum loss for the contracts
func calculateMaxLoss(contracts []models.OptionsContract) float64 {
	var loss, maxLoss float64

	for _, contract := range contracts {
		if contract.GetPosition() == long {
			loss = calculateLongLoss(contract)
		} else if contract.GetPosition() == short {
			loss = calculateShortLoss(contract)
		}
		if loss > maxLoss {
			maxLoss = loss
		}
	}

	return maxLoss
}

// calculateBreakEvenPoints finds the break even points for the contracts
func calculateBreakEvenPoints(contracts []models.OptionsContract) []float64 {
	var breakEvenPoints []float64

	for _, contract := range contracts {
		if contract.GetType() == call {
			breakEvenPoints = append(breakEvenPoints, calculateCallBreakEvenPoint(contract))
		} else if contract.GetType() == put {
			breakEvenPoints = append(breakEvenPoints, calculatePutBreakEvenPoint(contract))
		}
	}

	return breakEvenPoints
}

func calculateLongProfit(contract models.OptionsContract) float64 {
	if contract.GetType() == call {
		return contract.StrikePrice - contract.Bid
	} else if contract.GetType() == put {
		return contract.Ask - contract.StrikePrice
	}
	return 0
}

func calculateShortProfit(contract models.OptionsContract) float64 {
	if contract.GetType() == call {
		return contract.Bid - contract.StrikePrice
	} else if contract.GetType() == put {
		return contract.StrikePrice - contract.Ask
	}
	return 0
}

func calculateLongLoss(contract models.OptionsContract) float64 {
	if contract.GetType() == call {
		return contract.Bid - contract.StrikePrice
	} else if contract.GetType() == put {
		return contract.Ask - contract.StrikePrice
	}
	return 0
}

func calculateShortLoss(contract models.OptionsContract) float64 {
	if contract.GetType() == call {
		return contract.StrikePrice - contract.Bid
	} else if contract.GetType() == put {
		return contract.StrikePrice - contract.Ask
	}
	return 0
}

func calculateCallBreakEvenPoint(contract models.OptionsContract) float64 {
	if contract.GetPosition() == long {
		return contract.StrikePrice + contract.Bid
	} else if contract.GetPosition() == short {
		return contract.StrikePrice - contract.Bid
	}
	return 0
}

func calculatePutBreakEvenPoint(contract models.OptionsContract) float64 {
	if contract.GetPosition() == long {
		return contract.StrikePrice - contract.Ask
	} else if contract.GetPosition() == short {
		return contract.StrikePrice + contract.Ask
	}
	return 0
}
