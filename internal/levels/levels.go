package levels

type Level struct {
	Min   int
	Title string
}

var All = []Level{
	{0, "Novice"},
	{32, "Apprentice"},
	{64, "Student"},
	{96, "Scholar"},
	{128, "Adept"},
	{160, "Expert"},
	{192, "Master"},
	{224, "Grandmaster"},
	{256, "Math Architect"},
}

func Compute(mastered int) string {
	title := All[0].Title
	for _, l := range All {
		if mastered >= l.Min {
			title = l.Title
		}
	}
	return title
}
