package memstorage

import (
	"testing"

	"microcontrollers/internal/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHomeRamRepository_CreateHome(t *testing.T) {
	r := NewQueries()

	type args struct {
		homeId   string
		clientId string
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    bool
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {},
			input: args{
				homeId:   "asd",
				clientId: "0xArt",
			},
			wantErr: false,
		},
		{
			name: "Exist",
			mock: func() {},
			input: args{
				homeId:   "asd",
				clientId: "0xArt",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			_, err := r.CreateHome(tt.input.homeId, tt.input.clientId)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestHomeRamRepository_GetHome(t *testing.T) {
	r := NewQueries()

	type output struct {
		h *entity.Home
		b bool
	}

	tests := []struct {
		name    string
		mock    func()
		input   string
		want    output
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				r.CreateHome("newHome", "0xArt")
			},
			input: "newHome",
			want:  output{h: &entity.Home{ID: "newHome", ClientId: "0xArt"}, b: true},
		},
		{
			name: "Exist",
			mock: func() {
				r.CreateHome("newHome", "0xArt")
			},
			input: "asd",
			want:  output{h: &entity.Home{}, b: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			h, err := r.GetHome(tt.input)
			assert.Equal(t, err, tt.want.b)

			if tt.want.b {
				assert.Equal(t, h.ID, tt.want.h.ID)
			}
		})
	}
}

func TestHomeRamRepository_UpdateHome(t *testing.T) {
	r := NewQueries()

	type args struct {
		id string
		in entity.UpdateHomeInput
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    bool
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				r.CreateHome("newHome", "0xArt")
			},
			input: args{
				id: "newHome",
				in: entity.UpdateHomeInput{
					Temperature: stringPointer("35.5"),
					IsRobbery:   boolPointer(true),
				},
			},
			want: true,
		},
		{
			name: "Ok_WithoutIsRobberyAndIsLedTurned",
			mock: func() {
				r.CreateHome("newHome", "0xArt")
			},
			input: args{
				id: "newHome",
				in: entity.UpdateHomeInput{
					Temperature: stringPointer("39.5"),
				},
			},
			want: true,
		},
		{
			name: "Ok_WithoutFields",
			mock: func() {
				r.CreateHome("newHome", "0xArt")
			},
			input: args{
				id: "newHome",
				in: entity.UpdateHomeInput{},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			_, err := r.UpdateHome(tt.input.id, tt.input.in)
			require.NoError(t, err)
			//assert.Equal(t, ok, tt.want)
		})
	}
}

func TestHomeRamRepository_UpdateHomeSecurity(t *testing.T) {
	r := NewQueries()

	type args struct {
		id     string
		homeId string
		in     entity.UpdateHomeCommandInput
	}

	type output struct {
		cid string
		res bool
		sm  bool
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		want    output
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				r.CreateHome("newHome", "0xArt")
			},
			input: args{
				id: "0xArt",
				in: entity.UpdateHomeCommandInput{
					SecureMode:  boolPointer(true),
					NewClientId: stringPointer("0xMe"),
				},
				homeId: "newHome",
			},
			want: output{
				sm:  true,
				cid: "0xMe",
				res: true,
			},
		},
		{
			name: "Ok_WithoutIsRobberyAndIsLedTurned",
			mock: func() {
				r.CreateHome("new", "0xArt")
			},
			input: args{
				id: "0xArt",
				in: entity.UpdateHomeCommandInput{
					SecureMode: boolPointer(false),
				},
				homeId: "new",
			},
			want: output{
				sm:  false,
				res: true,
				cid: "0xArt",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			home, err := r.UpdateHomeInfo(tt.input.id, tt.input.in)

			require.NoError(t, err)

			assert.Equal(t, home.SecureMode, tt.want.sm)
			assert.Equal(t, home.ClientId, tt.want.cid)
		})
	}
}

func stringPointer(s string) *string {
	return &s
}

func boolPointer(b bool) *bool {
	return &b
}
