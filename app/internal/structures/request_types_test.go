package structures

import (
	"hometask/pkg/json_utils"
	"testing"
	"time"
)

func TestPutBirthdayRequest_Validate(t *testing.T) {
	type fields struct {
		DateOfBirth json_utils.JSONTime
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "birthday is tomorrow",
			fields:  fields{DateOfBirth: json_utils.JSONTime{Time: time.Now().AddDate(0, 0, 1)}},
			wantErr: true,
		},
		{
			name:    "birthday is literally today",
			fields:  fields{DateOfBirth: json_utils.JSONTime{Time: time.Now()}},
			wantErr: true,
		},
		{
			name:    "birthday is today, but in the past",
			fields:  fields{DateOfBirth: json_utils.JSONTime{Time: time.Now().AddDate(-10, 0, 0)}},
			wantErr: false,
		},
		{
			name:    "birthday is in the past",
			fields:  fields{DateOfBirth: json_utils.JSONTime{Time: time.Now().AddDate(-10, 0, 10)}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &PutBirthdayRequest{
				DateOfBirth: tt.fields.DateOfBirth,
			}
			if err := r.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
