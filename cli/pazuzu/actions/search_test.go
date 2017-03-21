package actions

import (
	"reflect"
	"testing"

	"github.com/zalando-incubator/pazuzu"
	"github.com/zalando-incubator/pazuzu/shared"
)

func TestSearchHandler(t *testing.T) {
	type args struct {
		feature string
		storage *pazuzu.TestStorage
	}

	featureMeta := pazuzu.GetTestFeatureMeta()

	tests := []struct {
		name    string
		args    args
		want    []shared.FeatureMeta
		wantErr bool
	}{
		{"Plain argument", args{"python", &pazuzu.TestStorage{}}, []shared.FeatureMeta{featureMeta}, false},
		{"Regex argument", args{"pyth*", &pazuzu.TestStorage{}}, []shared.FeatureMeta{featureMeta}, false},
		{"Incorrect regex", args{"pyth)", &pazuzu.TestStorage{}}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SearchHandler(tt.args.feature, tt.args.storage)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
