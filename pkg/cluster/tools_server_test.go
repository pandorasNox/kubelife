package cluster

import (
	"reflect"
	"testing"
)

func Test_extractProviderMachineTemplateValue(t *testing.T) {
	// 	in := `---
	// static:
	//   toolsServer:
	// toolsServer{}

	// 	in := `---
	// static:
	//   toolsServer:
	//     providerMachineTemplate:
	//       hetznerCloud:
	//       digitalocean:
	// `
	// toolsServer{}

	tests := []struct {
		name    string
		tsCfg   toolsServer
		want    reflect.Value
		wantErr bool
	}{
		{
			name:    "not set or empty toolsServer config",
			tsCfg:   toolsServer{},
			want:    reflect.Value{},
			wantErr: true,
		},
		{
			name: "one empty template privided",
			tsCfg: toolsServer{
				ProviderMachineTemplate: providerMachineTemplate{
					HetznerCloud: hetznerCloudMachine{},
				},
			},
			want:    reflect.Value{},
			wantErr: true,
		},
		{
			name: "two empty templates privided",
			tsCfg: toolsServer{
				ProviderMachineTemplate: providerMachineTemplate{
					HetznerCloud: hetznerCloudMachine{},
					Digitalocean: digitaloceanMachine{},
				},
			},
			want:    reflect.Value{},
			wantErr: true,
		},
		{
			name: "two non-empty templates privided",
			tsCfg: toolsServer{
				ProviderMachineTemplate: providerMachineTemplate{
					HetznerCloud: hetznerCloudMachine{
						ServerType: "someType",
					},
					Digitalocean: digitaloceanMachine{
						ServerType: "anotherType",
					},
				},
			},
			want:    reflect.Value{},
			wantErr: true,
		},
		{
			name: "one non-empty template & one empty provided",
			tsCfg: toolsServer{
				ProviderMachineTemplate: providerMachineTemplate{
					HetznerCloud: hetznerCloudMachine{
						ServerType: "someType",
					},
					Digitalocean: digitaloceanMachine{},
				},
			},
			want: extractIgnoreError(reflect.ValueOf(providerMachineTemplate{
				HetznerCloud: hetznerCloudMachine{
					ServerType: "someType",
				},
				Digitalocean: digitaloceanMachine{},
			})),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractProviderMachineTemplateValue(tt.tsCfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractProviderMachineTemplateValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want == (reflect.Value{}) && got != (reflect.Value{}) {
				t.Errorf("extractProviderMachineTemplateValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != (reflect.Value{}) && got == (reflect.Value{}) {
				t.Errorf("extractProviderMachineTemplateValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !(tt.want == (reflect.Value{}) && got == (reflect.Value{})) &&
				(!reflect.DeepEqual(got.Interface(), tt.want.Interface())) {
				t.Errorf("extractProviderMachineTemplateValue() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func extractIgnoreError(v reflect.Value) reflect.Value {
	res, _ := extractFirstNonEmpty(v)
	return res
}
