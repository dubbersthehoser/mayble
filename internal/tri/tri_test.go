package tri

import (
	"testing"
	"slices"
	"cmp"
)

func TestPrefixSearch(t *testing.T) {
	type input struct {
		text string
		handle int64
	}
	prefixes := []input{
		input{
			text: "apple",
			handle: 1,
		},
		input{
			text: "apache",
			handle: 2,
		},
		input{
			text: "april",
			handle: 3,
		},
		input{
			text: "zephyr",
			handle: 4,
		},
		input{
			text: "zeppelin",
			handle: 5,
		},
		input{
			text: "zen",
			handle: 6,
		},
		input{
			text: "zero",
			handle: 0,
		},
	}

	tri := Tri{}
	for _, prefix := range prefixes {
		tri.InsertPrefix(prefix.text, prefix.handle)
	}

	type testCase struct {
		search string
		handles []int64
	}

	cases := []testCase{
		testCase{
			search: "ap",
			handles: []int64{1, 2, 3},
		},
		testCase{
			search: "app",
			handles: []int64{1},
		},
		testCase{
			search: "z",
			handles: []int64{4, 5, 6, 0},
		},
		testCase{
			search: "ze",
			handles: []int64{4, 5, 6, 0},
		},
		testCase{
			search: "zep",
			handles: []int64{4, 5},
		},
		testCase{
			search: "zepp",
			handles: []int64{5},
		},
		testCase{
			search: "zen",
			handles: []int64{6},
		},
		testCase{
			search: "zero",
			handles: []int64{0},
		},
		testCase{
			search: "",
			handles: []int64{},
		},
	}


	for i, _case := range cases {
		ok, handles := tri.SearchPrefix(_case.search)
		if !ok {
			t.Fatalf("case %d, search found nothing", i)
		}
		if len(handles) != len(_case.handles) {
			t.Fatalf("case %d, missing handels expect %d, got %d", i, len(handles), len(_case.handles))
		}
		slices.SortFunc(_case.handles, cmp.Compare)
		slices.SortFunc(handles, cmp.Compare)
		for i := range handles {
			expect := _case.handles[i]
			actual := handles[i]
			if expect != actual {
				t.Errorf("case %d, missmatch handles", i)
			}
		}
	}
}
