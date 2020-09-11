package crypto

import (
	"crypto/cipher"
	"reflect"
	"testing"
)

func TestChaCha20Poly1305_DecryptedFromBase64(t *testing.T) {
	type fields struct {
		aead cipher.AEAD
	}
	type args struct {
		text string
	}
	instance := ChaCha20Poly1305{}
	instance.Init()
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{"Decryption 1", fields{aead: instance.aead}, args{text: "h7NYDdgHCpzR6ReKQ7TvslnrFDhYZL2VHzfdYhSD4wU="}, []byte("test"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ChaCha20Poly1305{
				aead: tt.fields.aead,
			}
			got, err := c.DecryptedFromBase64(tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecryptedFromBase64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecryptedFromBase64() got = %v, want %v", got, tt.want)
			}
		})
	}
}
