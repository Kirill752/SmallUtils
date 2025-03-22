package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type TestCase struct {
	SearchRequest
	Response SearchResponse
	Err      string
}

// http://localhost:8080/?query=Boyd&order_field=Id&orderBy=1&limit=10&offset=0
func TestFindUsers(t *testing.T) {
	TestCases := []TestCase{
		{
			SearchRequest: SearchRequest{
				Query:      "Boyd",
				OrderField: "Id",
				OrderBy:    1,
				Limit:      10,
				Offset:     0,
			},
			Response: SearchResponse{
				Users: []User{
					{
						Id:     0,
						Name:   "Boyd Wolf",
						Age:    22,
						About:  "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n",
						Gender: "male",
					},
				},
			},
		},
		{
			SearchRequest: SearchRequest{
				Query:      "Boyd",
				OrderField: "Id",
				OrderBy:    1,
				Limit:      -1,
				Offset:     0,
			},
			Response: SearchResponse{
				Users: []User{
					{
						Id:     0,
						Name:   "Boyd Wolf",
						Age:    22,
						About:  "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n",
						Gender: "male",
					},
				},
			},
			Err: "limit must be > 0",
		},
		{
			SearchRequest: SearchRequest{
				Query:      "Boyd",
				OrderField: "Id",
				OrderBy:    1,
				Limit:      30,
				Offset:     0,
			},
			Response: SearchResponse{
				Users: []User{
					{
						Id:     0,
						Name:   "Boyd Wolf",
						Age:    22,
						About:  "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n",
						Gender: "male",
					},
				},
			},
		},
		{
			SearchRequest: SearchRequest{
				Query:      "Boyd",
				OrderField: "Id",
				OrderBy:    1,
				Limit:      10,
				Offset:     -5,
			},
			Response: SearchResponse{
				Users: []User{
					{
						Id:     0,
						Name:   "Boyd Wolf",
						Age:    22,
						About:  "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n",
						Gender: "male",
					},
				},
			},
			Err: "offset must be > 0",
		},
		{
			SearchRequest: SearchRequest{
				Query:      "Boyd",
				OrderField: "Id",
				OrderBy:    1,
				Limit:      0,
				Offset:     0,
			},
			Response: SearchResponse{
				Users: []User{
					{
						Id:     0,
						Name:   "Boyd Wolf",
						Age:    22,
						About:  "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n",
						Gender: "male",
					},
				},
			},
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(SearchServerHandler))
	for caseNum, tc := range TestCases {
		// url = url + "query=" + tc.Query + "&order_field=" + tc.OrderField + "&orderBy=" + strconv.Itoa(tc.OrderBy) +
		// 	"&limit=" + strconv.Itoa(tc.Limit) + "&offset=0" + strconv.Itoa(tc.Offset)
		// req := httptest.NewRequest("GET", url, nil)
		// w := httptest.NewRecorder()
		srv := new(SearchClient)
		srv.URL = ts.URL
		resp, err := srv.FindUsers(tc.SearchRequest)
		if err != nil {
			if err.Error() != tc.Err {
				t.Errorf("[%d] Error FindUsers: %v", caseNum, err)
			}
		} else {
			for i, r := range resp.Users {
				if r.Name != tc.Response.Users[i].Name {
					t.Errorf("[%d] Wrong name of user:\n Expected: %v\n Got: %v", caseNum, tc.Response.Users[i].Name, r.Name)
				}
				if r.Id != tc.Response.Users[i].Id {
					t.Errorf("[%d] Wrong ID of user:\n Expected: %v\n Got: %v", caseNum, tc.Response.Users[i].Id, r.Id)
				}
				if r.Age != tc.Response.Users[i].Age {
					t.Errorf("[%d] Wrong Age of user:\n Expected: %v\n Got: %v", caseNum, tc.Response.Users[i].Age, r.Age)
				}
				if r.About != tc.Response.Users[i].About {
					t.Errorf("[%d] Wrong About of user:\n Expected: %v\n Got: %v", caseNum, tc.Response.Users[i].About, r.About)
				}
				if r.Gender != tc.Response.Users[i].Gender {
					t.Errorf("[%d] Wrong Gender of user:\n Expected: %v\n Got: %v", caseNum, tc.Response.Users[i].Gender, r.Gender)
				}
			}
		}
	}
}

func TestErrors(t *testing.T) {
	TestCases := []TestCase{
		{
			SearchRequest: SearchRequest{
				Query:      "StatusUnauthorized",
				OrderField: "",
				OrderBy:    0,
				Limit:      0,
				Offset:     0,
			},
			Err: "Bad AccessToken",
		},
		{
			SearchRequest: SearchRequest{
				Query:      "StatusInternalServerError",
				OrderField: "",
				OrderBy:    0,
				Limit:      0,
				Offset:     0,
			},
			Err: "SearchServer fatal error",
		},
		{
			SearchRequest: SearchRequest{
				Query:      "StatusBadRequest",
				OrderField: "",
				OrderBy:    0,
				Limit:      0,
				Offset:     0,
			},
			Err: "cant unpack error json: unexpected end of JSON input",
		},
		{
			SearchRequest: SearchRequest{
				Query:      "StatusBadRequest_ErrorOrder",
				OrderField: "abs",
				OrderBy:    0,
				Limit:      0,
				Offset:     0,
			},
			Err: "OrderFeld abs invalid",
		},
		{
			SearchRequest: SearchRequest{
				Query:      "StatusBadRequest_Unknown",
				OrderField: "abs",
				OrderBy:    0,
				Limit:      0,
				Offset:     0,
			},
			Err: "unknown bad request error: Unknown",
		},
		{
			SearchRequest: SearchRequest{
				Query:      "StatusOk_ErrorJSON",
				OrderField: "",
				OrderBy:    0,
				Limit:      0,
				Offset:     0,
			},
			Err: "cant unpack result json: unexpected end of JSON input",
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(SearchServerErrors))
	for caseNum, tc := range TestCases {
		srv := new(SearchClient)
		srv.URL = ts.URL
		_, err := srv.FindUsers(tc.SearchRequest)
		if err != nil {
			if err.Error() != tc.Err {
				t.Errorf("[%d] expected %s, found %s", caseNum, tc.Err, err.Error())
			}
		}
	}
	ts.Close()
}

func TestTimeOut(t *testing.T) {
	TestCases := []TestCase{
		{
			SearchRequest: SearchRequest{
				Query:      "Timeout",
				OrderField: "",
				OrderBy:    0,
				Limit:      0,
				Offset:     0,
			},
			Err: "timeout for limit=1&offset=0&order_by=0&order_field=&query=Timeout",
		},
	}
	ts := httptest.NewServer(http.TimeoutHandler(http.HandlerFunc(SearchServerErrors), 10*time.Second, "timeout"))
	for caseNum, tc := range TestCases {
		srv := new(SearchClient)
		srv.URL = ts.URL
		_, err := srv.FindUsers(tc.SearchRequest)
		if err != nil {
			if err.Error() != tc.Err {
				t.Errorf("[%d] expected %s, found %s", caseNum, tc.Err, err.Error())
			}
		}
	}
	ts.Close()
}
func TestUnknownError(t *testing.T) {
	TestCases := []TestCase{
		{
			SearchRequest: SearchRequest{
				Query:      "Unknown",
				OrderField: "",
				OrderBy:    0,
				Limit:      0,
				Offset:     0,
			},
			Err: `unknown error Get "htp:\\rr?limit=1&offset=0&order_by=0&order_field=&query=Unknown": unsupported protocol scheme "htp"`,
		},
	}
	for caseNum, tc := range TestCases {
		srv := new(SearchClient)
		srv.URL = "htp:\\rr"
		_, err := srv.FindUsers(tc.SearchRequest)
		if err != nil {
			if err.Error() != tc.Err {
				t.Errorf("[%d] expected %s, found %s", caseNum, tc.Err, err.Error())
			}
		}
	}
}
