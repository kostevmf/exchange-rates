package exrate

type Options struct {
	CurID    int     `json:"cur_id"`
	Date     string  `json:"date"`
	CurCode  string  `json:"cur_code"`
	CurName  string  `json:"cur_name"`
	CurScale int     `json:"cur_scale"`
	CurRate  float64 `json:"cur_rate"`
}
