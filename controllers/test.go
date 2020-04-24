package controllers

func Test(ctx *gin.WebContext) {
	type Movie struct {
		Month string `json:"month"`
		Title string `json:"title"`
	}
	var movie Movie
	movie.Month = "may"
	movie.Title = "can not sleep"
	str, _ := json.Marshal(movie)
	ctx.Json(200, movie)
	return
}
