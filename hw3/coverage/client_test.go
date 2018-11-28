package main

// тут писать код тестов

import (
	"testing"
	//"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"time"
)

type TestCase struct {
	tClient   *SearchClient
	tRequest  *SearchRequest
	tResponse *SearchResponse
	tError    string
}

const (
	fileName = "dataset.xml"
)

func TestUrl(t *testing.T) {
	ts := httptest.NewServer(SearchServer(fileName))
	defer ts.Close()

	tCase := TestCase{
			tClient: &SearchClient{
				AccessToken: "UseOrDie",
				URL:         "",
			},
			tRequest: &SearchRequest{
				Limit:      10,
				Offset:     1,
				Query:      "Snow",
				OrderField: "Id",
				OrderBy:    0,
			},
			tError: "unknown error Get ?limit=11&offset=1&order_by=0&order_field=Id&query=Snow: unsupported protocol scheme \"\"",
	}

	_, err := tCase.tClient.FindUsers(*tCase.tRequest)

	if err == nil {
		t.Errorf("Expected error: %v", err.Error())
		t.FailNow()
	}

	if err.Error() != tCase.tError{
		t.Errorf("Unexpected error: %v", err.Error())
		t.FailNow()
	}
}

func TimeoutError(w http.ResponseWriter, r *http.Request){
	time.Sleep( time.Second )
}

func TestTimeoutError(t *testing.T){
	ts := httptest.NewServer(http.HandlerFunc(TimeoutError))
	defer ts.Close()

	tCase := TestCase{
			tClient: &SearchClient{
				AccessToken: "UseOrDie",
				URL:         ts.URL,
			},
			tRequest: &SearchRequest{
				Limit:      2,
				Offset:     1,
				Query:      "Snow",
				OrderField: "Id",
				OrderBy:    0,
			},
			tError: "timeout for limit=3&offset=1&order_by=0&order_field=Id&query=Snow",
	}

	_, err := tCase.tClient.FindUsers(*tCase.tRequest)

	if err == nil {
		t.Errorf("Expected error: %v", err.Error())
		t.FailNow()
	}

	if err.Error() != tCase.tError{
		t.Errorf("Unexpected error: %v", err.Error())
		t.FailNow()
	}
}

func UnpackErrorJson(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}

func TestUnpackErrorJson(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(UnpackErrorJson))
	defer ts.Close()

	tCase := TestCase{
		tClient: &SearchClient{
			AccessToken: "UseOrDie",
			URL:         ts.URL,
		},
		tRequest: &SearchRequest{
			Limit:      1,
			Offset:     0,
			Query:      "Boyd Wolf",
			OrderField: "Id",
			OrderBy:    -1,
		},
		tError: "cant unpack error json: unexpected end of JSON input",
	}

	_, err := tCase.tClient.FindUsers(*tCase.tRequest)

	if err == nil {
		t.Errorf("Expected error: %v", err.Error())
		t.FailNow()
	}

	if err.Error() != tCase.tError {
		t.Errorf("Unexpected error: %v", err.Error())
		t.FailNow()
	}
}

func UnpackResultJson(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Mons The Best"))
}

func TestUnpackResultJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(UnpackResultJson))
	defer ts.Close()

	tCase := TestCase{
		tClient: &SearchClient{
			AccessToken: "UseOrDie",
			URL:         ts.URL,
		},
		tRequest: &SearchRequest{
			Limit:      1,
			Offset:     0,
			Query:      "Boyd Wolf",
			OrderField: "Id",
			OrderBy:    -1,
		},
		tError: "cant unpack result json: invalid character 'M' looking for beginning of value",
	}

	_, err := tCase.tClient.FindUsers(*tCase.tRequest)

	if err == nil {
		t.Errorf("Expected error: %v", err.Error())
		t.FailNow()
	}

	if err.Error() != tCase.tError {
		t.Errorf("Unexpected error: %v", err.Error())
		t.FailNow()
	}
}

func TestOpenXmlFile(t *testing.T) {
	ts := httptest.NewServer(SearchServer("error.xml"))
	defer ts.Close()
	tCase := TestCase{
		tClient: &SearchClient{
			AccessToken: "UseOrDie",
			URL:         ts.URL,
		},
		tRequest: &SearchRequest{
			Limit:      1,
			Offset:     0,
			Query:      "Boyd Wolf",
			OrderField: "Id",
			OrderBy:    -1,
		},
		tError: "SearchServer fatal error",
	}

	_, err := tCase.tClient.FindUsers(*tCase.tRequest)

	if err == nil {
		t.Errorf("Expected error: %v", err.Error())
		t.FailNow()
	}

	if err.Error() != tCase.tError {
		t.Errorf("Unexpected error: %v", err.Error())
		t.FailNow()
	}
}

func TestBadRequest(t *testing.T) {
	ts := httptest.NewServer(SearchServer(fileName))
	defer ts.Close()
	tCase := TestCase{
		tClient: &SearchClient{
			AccessToken: "UseOrDie",
			URL:         ts.URL,
		},
		tRequest: &SearchRequest{
			Limit:      1,
			Offset:     0,
			Query:      "Boyd Wolf",
			OrderField: "Name",
			OrderBy:    -10,
		},
		tError: "unknown bad request error: BadOrderBy",
	}

	_, err := tCase.tClient.FindUsers(*tCase.tRequest)

	if err == nil {
		t.Errorf("Expected error: %v", err.Error())
		t.FailNow()
	}

	if err.Error() != tCase.tError {
		t.Errorf("Unexpected error: %v", err.Error())
		t.FailNow()
	}

}
func TestOrderField(t *testing.T) {
	ts := httptest.NewServer(SearchServer(fileName))
	defer ts.Close()

	tCase := TestCase{
		tClient: &SearchClient{
			AccessToken: "UseOrDie",
			URL:         ts.URL,
		},
		tRequest: &SearchRequest{
			Limit:      1,
			Offset:     0,
			Query:      "Boyd Wolf",
			OrderField: "Gender",
			OrderBy:    -1,
		},
		tError: "OrderField Gender invalid",
	}

	_, err := tCase.tClient.FindUsers(*tCase.tRequest)

	if err == nil {
		t.Errorf("Expected error: %v", err.Error())
		t.FailNow()
	}

	if err.Error() != tCase.tError {
		t.Errorf("Unexpected error: %v", err.Error())
		t.FailNow()
	}
}

func TestAccessToken(t *testing.T) {
	ts := httptest.NewServer(SearchServer(fileName))
	defer ts.Close()

	tCase := TestCase{
		tClient: &SearchClient{
			AccessToken: "DontUse",
			URL:         ts.URL,
		},
		tRequest: &SearchRequest{
			Limit:      1,
			Offset:     0,
			Query:      "Boyd Wolf",
			OrderField: "",
			OrderBy:    -1,
		},
		tError: "Bad AccessToken",
	}

	_, err := tCase.tClient.FindUsers(*tCase.tRequest)

	if err == nil {
		t.Errorf("Expected error: %v", err.Error())
		t.FailNow()
	}

	if err.Error() != tCase.tError {
		t.Errorf("Unexpected error: %v", err.Error())
		t.FailNow()
	}

}

func TestOffset(t *testing.T) {
	ts := httptest.NewServer(SearchServer(fileName))
	defer ts.Close()

	tCase := TestCase{
		tRequest: &SearchRequest{
			Limit:      1,
			Offset:     -1,
			Query:      "Boyd Wolf",
			OrderField: "",
			OrderBy:    -1,
		},
		tError: "offset must be > 0",
	}
	_, err := tCase.tClient.FindUsers(*tCase.tRequest)

	if err == nil {
		t.Errorf("Expected error: %v", err.Error())
		t.FailNow()
	}

	if err.Error() != tCase.tError {
		t.Errorf("Unexpected error: %v", err.Error())
		t.FailNow()
	}
}

func TestLimit(t *testing.T) {
	ts := httptest.NewServer(SearchServer(fileName))
	defer ts.Close()

	tCase := TestCase{
		tRequest: &SearchRequest{
			Limit:      -1,
			Offset:     0,
			Query:      "Boyd Wolf",
			OrderField: "",
			OrderBy:    -1,
		},
		tError: "limit must be > 0",
	}
	_, err := tCase.tClient.FindUsers(*tCase.tRequest)

	if err == nil {
		t.Errorf("Expected error: %v", err.Error())
		t.FailNow()
	}

	if err.Error() != tCase.tError {
		t.Errorf("Unexpected error: %v", err.Error())
		t.FailNow()
	}

}

func TestRequests(t *testing.T) {
	ts := httptest.NewServer(SearchServer(fileName))
	defer ts.Close()

	cases := []TestCase{
		TestCase{
			tClient: &SearchClient{
				AccessToken: "UseOrDie",
				URL:         ts.URL,
			},
			tRequest: &SearchRequest{
				Limit:      2,
				Offset:     0,
				Query:      "",
				OrderField: "Id",
				OrderBy:    -1,
			},
			tResponse: &SearchResponse{
				Users: []User{
					User{
						Id:     34,
						Name:   "Kane Sharp",
						Age:    34,
						About:  "Lorem proident sint minim anim commodo cillum. Eiusmod velit culpa commodo anim consectetur consectetur sint sint labore. Mollit consequat consectetur magna nulla veniam commodo eu ut et. Ut adipisicing qui ex consectetur officia sint ut fugiat ex velit cupidatat fugiat nisi non. Dolor minim mollit aliquip veniam nostrud. Magna eu aliqua Lorem aliquip.\n",
						Gender: "male",
					},
					User{
						Id:     33,
						Name:   "Twila Snow",
						Age:    36,
						About:  "Sint non sunt adipisicing sit laborum cillum magna nisi exercitation. Dolore officia esse dolore officia ea adipisicing amet ea nostrud elit cupidatat laboris. Proident culpa ullamco aute incididunt aute. Laboris et nulla incididunt consequat pariatur enim dolor incididunt adipisicing enim fugiat tempor ullamco. Amet est ullamco officia consectetur cupidatat non sunt laborum nisi in ex. Quis labore quis ipsum est nisi ex officia reprehenderit ad adipisicing fugiat. Labore fugiat ea dolore exercitation sint duis aliqua.\n",
						Gender: "female",
					},
				},
				NextPage: true,
			},
		},
		TestCase{
			tClient: &SearchClient{
				AccessToken: "UseOrDie",
				URL:         ts.URL,
			},
			tRequest: &SearchRequest{
				Limit:      2,
				Offset:     0,
				Query:      "",
				OrderField: "Id",
				OrderBy:    1,
			},
			tResponse: &SearchResponse{
				Users: []User{
					User{
						Id:     0,
						Name:   "Boyd Wolf",
						Age:    22,
						About:  "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n",
						Gender: "male",
					},
					User{
						Id:     1,
						Name:   "Hilda Mayer",
						Age:    21,
						About:  "Sit commodo consectetur minim amet ex. Elit aute mollit fugiat labore sint ipsum dolor cupidatat qui reprehenderit. Eu nisi in exercitation culpa sint aliqua nulla nulla proident eu. Nisi reprehenderit anim cupidatat dolor incididunt laboris mollit magna commodo ex. Cupidatat sit id aliqua amet nisi et voluptate voluptate commodo ex eiusmod et nulla velit.\n",
						Gender: "female",
					},
				},
				NextPage: true,
			},
		},
		TestCase{
			tClient: &SearchClient{
				AccessToken: "UseOrDie",
				URL:         ts.URL,
			},
			tRequest: &SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "Wolf",
				OrderField: "",
				OrderBy:    0,
			},
			tResponse: &SearchResponse{
				Users: []User{
					User{
						Id:     0,
						Name:   "Boyd Wolf",
						Age:    22,
						About:  "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n",
						Gender: "male",
					},
				},
				NextPage: false,
			},
		},

		TestCase{
			tClient: &SearchClient{
				AccessToken: "UseOrDie",
				URL:         ts.URL,
			},
			tRequest: &SearchRequest{
				Limit:      2018,
				Offset:     0,
				Query:      "Boyd",
				OrderField: "Id",
				OrderBy:    -1,
			},
			tResponse: &SearchResponse{
				Users: []User{
					User{
						Id:     0,
						Name:   "Boyd Wolf",
						Age:    22,
						About:  "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n",
						Gender: "male",
					},
				},
				NextPage: false,
			},
		},
		TestCase{
			tClient: &SearchClient{
				AccessToken: "UseOrDie",
				URL:         ts.URL,
			},
			tRequest: &SearchRequest{
				Limit:      2,
				Offset:     0,
				Query:      "kek",
				OrderField: "Name",
				OrderBy:    1,
			},
			tResponse: &SearchResponse{
				Users: []User{
					User{
						Id:     15,
						Name:   "Allison Valdez",
						Age:    21,
						About:  "Labore excepteur voluptate velit occaecat est nisi minim. Laborum ea et irure nostrud enim sit incididunt reprehenderit id est nostrud eu. Ullamco sint nisi voluptate cillum nostrud aliquip et minim. Enim duis esse do aute qui officia ipsum ut occaecat deserunt. Pariatur pariatur nisi do ad dolore reprehenderit et et enim esse dolor qui. Excepteur ullamco adipisicing qui adipisicing tempor minim aliquip.\n",
						Gender: "male",
					},

					User{
						Id:     16,
						Name:   "Annie Osborn",
						Age:    35,
						About:  "Consequat fugiat veniam commodo nisi nostrud culpa pariatur. Aliquip velit adipisicing dolor et nostrud. Eu nostrud officia velit eiusmod ullamco duis eiusmod ad non do quis.\n",
						Gender: "female",
					},
				},
				NextPage: true,
			},
		},
		TestCase{
			tClient: &SearchClient{
				AccessToken: "UseOrDie",
				URL:         ts.URL,
			},
			tRequest: &SearchRequest{
				Limit:      2,
				Offset:     0,
				Query:      "",
				OrderField: "Name",
				OrderBy:    -1,
			},
			tResponse: &SearchResponse{
				Users: []User{
					User{
						Id:     13,
						Name:   "Whitley Davidson",
						Age:    40,
						About:  "Consectetur dolore anim veniam aliqua deserunt officia eu. Et ullamco commodo ad officia duis ex incididunt proident consequat nostrud proident quis tempor. Sunt magna ad excepteur eu sint aliqua eiusmod deserunt proident. Do labore est dolore voluptate ullamco est dolore excepteur magna duis quis. Quis laborum deserunt ipsum velit occaecat est laborum enim aute. Officia dolore sit voluptate quis mollit veniam. Laborum nisi ullamco nisi sit nulla cillum et id nisi.\n",
						Gender: "male",
					},
					User{
						Id:     33,
						Name:   "Twila Snow",
						Age:    36,
						About:  "Sint non sunt adipisicing sit laborum cillum magna nisi exercitation. Dolore officia esse dolore officia ea adipisicing amet ea nostrud elit cupidatat laboris. Proident culpa ullamco aute incididunt aute. Laboris et nulla incididunt consequat pariatur enim dolor incididunt adipisicing enim fugiat tempor ullamco. Amet est ullamco officia consectetur cupidatat non sunt laborum nisi in ex. Quis labore quis ipsum est nisi ex officia reprehenderit ad adipisicing fugiat. Labore fugiat ea dolore exercitation sint duis aliqua.\n",
						Gender: "female",
					},
				},
				NextPage: true,
			},
		},
		TestCase{
			tClient: &SearchClient{
				AccessToken: "UseOrDie",
				URL:         ts.URL,
			},
			tRequest: &SearchRequest{
				Limit:      2,
				Offset:     0,
				Query:      "kek",
				OrderField: "Age",
				OrderBy:    1,
			},
			tResponse: &SearchResponse{
				Users: []User{
					User{
						Id:     1,
						Name:   "Hilda Mayer",
						Age:    21,
						About:  "Sit commodo consectetur minim amet ex. Elit aute mollit fugiat labore sint ipsum dolor cupidatat qui reprehenderit. Eu nisi in exercitation culpa sint aliqua nulla nulla proident eu. Nisi reprehenderit anim cupidatat dolor incididunt laboris mollit magna commodo ex. Cupidatat sit id aliqua amet nisi et voluptate voluptate commodo ex eiusmod et nulla velit.\n",
						Gender: "female",
					},
					User{
						Id:     15,
						Name:   "Allison Valdez",
						Age:    21,
						About:  "Labore excepteur voluptate velit occaecat est nisi minim. Laborum ea et irure nostrud enim sit incididunt reprehenderit id est nostrud eu. Ullamco sint nisi voluptate cillum nostrud aliquip et minim. Enim duis esse do aute qui officia ipsum ut occaecat deserunt. Pariatur pariatur nisi do ad dolore reprehenderit et et enim esse dolor qui. Excepteur ullamco adipisicing qui adipisicing tempor minim aliquip.\n",
						Gender: "male",
					},
				},
				NextPage: true,
			},
		},
		TestCase{
			tClient: &SearchClient{
				AccessToken: "UseOrDie",
				URL:         ts.URL,
			},
			tRequest: &SearchRequest{
				Limit:      2,
				Offset:     0,
				Query:      "",
				OrderField: "Age",
				OrderBy:    -1,
			},
			tResponse: &SearchResponse{
				Users: []User{
					User{
						Id:     32,
						Name:   "Christy Knapp",
						Age:    40,
						About:  "Incididunt culpa dolore laborum cupidatat consequat. Aliquip cupidatat pariatur sit consectetur laboris labore anim labore. Est sint ut ipsum dolor ipsum nisi tempor in tempor aliqua. Aliquip labore cillum est consequat anim officia non reprehenderit ex duis elit. Amet aliqua eu ad velit incididunt ad ut magna. Culpa dolore qui anim consequat commodo aute.\n",
						Gender: "female",
					},
					User{
						Id:     13,
						Name:   "Whitley Davidson",
						Age:    40,
						About:  "Consectetur dolore anim veniam aliqua deserunt officia eu. Et ullamco commodo ad officia duis ex incididunt proident consequat nostrud proident quis tempor. Sunt magna ad excepteur eu sint aliqua eiusmod deserunt proident. Do labore est dolore voluptate ullamco est dolore excepteur magna duis quis. Quis laborum deserunt ipsum velit occaecat est laborum enim aute. Officia dolore sit voluptate quis mollit veniam. Laborum nisi ullamco nisi sit nulla cillum et id nisi.\n",
						Gender: "male",
					},
				},
				NextPage: true,
			},
		},
	}

	for caseNum, item := range cases {
		resp, err := item.tClient.FindUsers(*item.tRequest)

		if err == nil && item.tError != "" {
			t.Errorf("[%d] expect error: %#v", caseNum, item.tError)
			t.FailNow()
		}
		if err != nil && err.Error() != item.tError {
			t.Errorf("[%d] unexpected error: %v", caseNum, err)
			t.FailNow()
		}
		if !reflect.DeepEqual(item.tResponse, resp) {
			t.Errorf("[%d] wrong result, expected %#v, got %#v", caseNum, item.tResponse, resp)
			t.FailNow()
		}

	}

}

/**/
