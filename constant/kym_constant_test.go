package constant

import "testing"

func TestIsValidKymStatus(t *testing.T) {
	type args struct {
		status string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "KymStatusNotAllowed",
			args: args{
				status: "12345",
			},
			want: false,
		},
		{
			name: "KymStatusAllowed",
			args: args{
				status: "rejected",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidKymStatus(tt.args.status); got != tt.want {
				t.Errorf("isValidKymStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
