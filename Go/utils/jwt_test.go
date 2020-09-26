package utils

import (
	"github.com/dgrijalva/jwt-go"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Generate"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(CustomClaims{
				UID:            "",
				Username:       "",
				Name:           "",
				Phone:          "",
				Flag:           "",
				Role:           0,
				StandardClaims: jwt.StandardClaims{},
			})
			if token == "" {
				t.Errorf("err:%s", err.Error())
			}
		})
	}
}

//func TestJWT_parseToken(t *testing.T) {
//	type fields struct {
//		SigningKey []byte
//	}
//	type args struct {
//		tokenString string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    *CustomClaims
//		wantErr bool
//	}{
//		{"parse", fields{SigningKey: []byte(SignKey)}, args{tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJ0ZXN0IiwibmFtZSI6InRlc3QiLCJwaG9uZSI6InRlc3QiLCJyb2xlIjowLCJleHAiOjE1OTk2NDk0NjUsImlzcyI6IkpXVFRva2VuTWFuYWdlciIsIm5iZiI6MTU5OTY0NDg2NX0.RPuSnfvTdgytXTJuhjcrC7Qs_kb3XSLZ6kxuGBu_WSk"}, &CustomClaims{
//			UID:   "test",
//			Name:  "test",
//			Phone: "test",
//			Role:  0,
//			StandardClaims: jwt.StandardClaims{
//				NotBefore: time.Now().Unix() - 1000, // 签名生效时间
//				ExpiresAt: time.Now().Unix() + 3600, // 过期时间 一小时
//				Issuer:    Issuer,                   //签名的发行者
//			},
//		}, false},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			j := &JWT{
//				SigningKey: tt.fields.SigningKey,
//			}
//			got, err := j.parseToken(tt.args.tokenString)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("parseToken() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("parseToken() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestJWT_refreshToken(t *testing.T) {
//	type fields struct {
//		SigningKey []byte
//	}
//	type args struct {
//		tokenString string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    string
//		wantErr bool
//	}{
//		{"parse", fields{SigningKey: []byte(SignKey)}, args{tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJ0ZXN0IiwibmFtZSI6InRlc3QiLCJwaG9uZSI6InRlc3QiLCJyb2xlIjowLCJleHAiOjE1OTk2NDk0NjUsImlzcyI6IkpXVFRva2VuTWFuYWdlciIsIm5iZiI6MTU5OTY0NDg2NX0.RPuSnfvTdgytXTJuhjcrC7Qs_kb3XSLZ6kxuGBu_WSk"},
//			"",
//			false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			j := &JWT{
//				SigningKey: tt.fields.SigningKey,
//			}
//			got, err := j.refreshToken(tt.args.tokenString)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("refreshToken() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if got != tt.want {
//				t.Errorf("refreshToken() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
