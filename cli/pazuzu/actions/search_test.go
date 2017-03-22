package actions

import (
	"reflect"
	"testing"

	"github.com/zalando-incubator/pazuzu/shared"
	"github.com/zalando-incubator/pazuzu/mock"
)

func TestSearchHandler(t *testing.T) {
	type args struct {
		feature string
		storage *mock.TestStorage
	}

	featureMeta := mock.GetTestFeatureMeta()

	tests := []struct {
		name    string
		args    args
		want    []shared.FeatureMeta
		wantErr bool
	}{
		{"Plain argument", args{"python", &mock.TestStorage{}}, []shared.FeatureMeta{featureMeta}, false},
		{"Regex argument", args{"pyth*", &mock.TestStorage{}}, []shared.FeatureMeta{featureMeta}, false},
		{"Incorrect regex", args{"pyth)", &mock.TestStorage{}}, nil, true},
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
