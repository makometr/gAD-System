package domain

type Calculator struct {
	repo Repository
}

func NewCalculator(repo Repository) *Calculator {
	return &Calculator{repo: repo}
}

func (c Calculator) Calculate(expr string) (string, error) {
	return c.repo.doCalculate(expr)
}
