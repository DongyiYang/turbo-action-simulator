package rest


var RoutesPaths []Route

func init() {
	RoutesPaths = []Route{
		{
			"/api/actions",
			[]string{"POST"},
		},
		{
			"/api/actions/{id}",
			[]string{"GET", "DELETE"},
		},
		{
			"/api/discovery",
			[]string{"POST"},
		},
	}
}

type Route struct{
	Path string
	Method []string
}