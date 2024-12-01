package plan

import "goraj/limited-network-driver/pkg/billing"

type Mb int
type Percentage int

type Plan struct {
	monthlyData Mb
	softLimit   Percentage
}

func NewPlan(data Mb, softLimit Percentage) *Plan {
	return &Plan{
		monthlyData: data,
		softLimit:   softLimit,
	}
}

func (plan *Plan) getCurrentThreshold(billing billing.Billing) int {
	monthlyData := int(plan.monthlyData)

	dailyData := monthlyData / billing.GetDaysInCurrentBillingPeriod()
	threshold := dailyData * billing.GetBillingPeriodCurrentDay()

	return threshold
}

func (plan *Plan) getSoftLimit(billing billing.Billing) int {
	percentageSetting := int(plan.softLimit)
	threshold := plan.getCurrentThreshold(billing)
	return threshold * percentageSetting / 100
}
