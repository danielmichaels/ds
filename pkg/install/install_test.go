package install

import (
	"testing"
)

func TestGoCheck(t *testing.T) {
	tt := []struct {
		name string
		exe  string
		want string
	}{
		{
			name: "ls executable available",
			exe:  "ls",
			want: "",
		},
		{
			name: "go executable available",
			exe:  "go",
			want: "",
		},
		{
			name: "executable unavailable errors",
			exe:  "bad123",
			want: `exec: "bad123": executable file not found in $PATH`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := exeCheck(tc.exe)
			if err != nil {
				if err.Error() != tc.want {
					t.Fatalf("test %s failed.\ngot:  %#v\nwant: %#v", tc.name, err.Error(), tc.want)
				}
			}

		})
	}
}
