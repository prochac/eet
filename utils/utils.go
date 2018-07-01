package utils

import (
	"fmt"
	"time"
)

//Prepares PKP for EET message
func preparePkp(dic string, provozovna int, pokladna string, uctenka string, datum time.Time, trzba float64) string {
	return fmt.Sprintf("%s|%d|%s|%s|%s|%0.2f",
		dic,
		provozovna,
		pokladna,
		uctenka,
		datum,
		trzba)
}
