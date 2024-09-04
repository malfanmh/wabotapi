package usecase

import (
	"context"
	"github.com/malfanmh/wabotapi"
	"github.com/malfanmh/wabotapi/model"
	"github.com/malfanmh/wabotapi/repository"
	"testing"
)

func Test_useCase_generateMessageBody(t *testing.T) {
	db := wabotapi.OpenMysqlDB("wabot_user:password@tcp(127.0.0.1:3306)/wabot_db?parseTime=true")
	repo := repository.NewMysql(db)

	type fields struct {
		repo Repository
		wa   WhatsAppRepository
	}
	type args struct {
		ctx       context.Context
		clientID  int64
		messageID int64
		access    model.Access
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantJsonMessage string
		wantErr         bool
	}{
		{
			name: "text",
			fields: fields{
				repo: repo,
				wa:   nil,
			},
			args: args{
				ctx:       context.TODO(),
				clientID:  1,
				messageID: 1,
				access:    0,
			},
			wantJsonMessage: "{\"body\":\"Assalamualaikum {{.Name}} Selamat datang di Layanan WA Muhammadiyah. \\n\\nSebelum melanjutkan, harap masukan terlebih dahulu No KTA dan Nama Anggota Anda dengan format: \\n\\nNo. KTA#Nama Anggota\\ncontoh:\\n123456789#Ahmad Sayuri\"}",
			wantErr:         false,
		},
		{
			name: "list",
			fields: fields{
				repo: repo,
				wa:   nil,
			},
			args: args{
				ctx:       context.TODO(),
				clientID:  1,
				messageID: 2,
				access:    0,
			},
			wantJsonMessage: "{\"body\":\"Assalamualaikum {{.Name}} Selamat datang di Layanan WA Muhammadiyah. \\n\\nSebelum melanjutkan, harap masukan terlebih dahulu No KTA dan Nama Anggota Anda dengan format: \\n\\nNo. KTA#Nama Anggota\\ncontoh:\\n123456789#Ahmad Sayuri\"}",
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &useCase{
				repo: tt.fields.repo,
				wa:   tt.fields.wa,
			}
			gotJsonMessage, err := uc.generateMessageBody(tt.args.ctx, tt.args.clientID, tt.args.messageID, tt.args.access)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateMessageBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotJsonMessage != tt.wantJsonMessage {
				t.Errorf("generateMessageBody() gotJsonMessage = %v, want %v", gotJsonMessage, tt.wantJsonMessage)
			}
		})
	}
}
