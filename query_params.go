package vocode

type ListNumbersQueryParams struct {
	Page     int
	Size     int
	SortCol  string
	SortDesc bool
}

type GetNumberQueryParams struct {
	PhoneNumber string
}
