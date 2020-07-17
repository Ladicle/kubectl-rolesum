package cmd

import (
	"testing"

	"github.com/Ladicle/kubectl-rolesum/pkg/util/subject"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		desc      string
		args      []string
		option    Option
		wantError bool
	}{
		{
			desc: "Use Service Account",
			args: []string{"ci-bot"},
			option: Option{
				SubjectKind: subject.KindSA,
			},
		},
		{
			desc: "Use User",
			args: []string{"alice"},
			option: Option{
				SubjectKind: subject.KindUser,
			},
		},
		{
			desc: "Use Group",
			args: []string{"developer"},
			option: Option{
				SubjectKind: subject.KindGroup,
			},
		},
		{
			desc:      "No arguments",
			wantError: true,
		},
		{
			desc: "Unknown subject kind",
			option: Option{
				SubjectKind: "unknown",
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			err := tt.option.Validate(nil, tt.args)
			if err != nil {
				if !tt.wantError {
					t.Fatalf("o.Validate(): %v", err)
				}
				return
			}
			if tt.wantError {
				t.Fatal("o.Validate() should fail")
			}
		})
	}
}
