package pagination

type Pagination struct {
	Limit     int `form:"limit"`
	Page      int `form:"page"`
	Offset    int
}

const MaximumLimit = 100

func (p *Pagination) Paginate() {
	p.ValidatePagination()
	p.Offset = p.Limit * (p.Page - 1)
}

func (p *Pagination) ValidatePagination() {
	if p.Page < 1 || p.Limit < 1 {
		p.Page, p.Limit = 1, 10
	}
	if p.Limit > MaximumLimit {
		p.Limit = MaximumLimit
	}
}