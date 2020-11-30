package responses

type MathCaptchaResponse struct {
	ID     string `json:"id"`
	ImgB64 string `json:"img_b64"`
}

type SlideCapchaXYResponse struct {
	x     int64              `json:"x"`
	y     int64              `json:"y"`
	ID    string             `json:"id"`
	MapXY []map[string]int64 `json:"xy"`
}

func (resp *SlideCapchaXYResponse) GetX() int64 {
	return resp.x
}

func (resp *SlideCapchaXYResponse) GetY() int64 {
	return resp.y
}

func (resp *SlideCapchaXYResponse) SetX(x int64) *SlideCapchaXYResponse {
	resp.x = x
	return resp
}

func (resp *SlideCapchaXYResponse) SetY(y int64) *SlideCapchaXYResponse {
	resp.y = y
	return resp
}
