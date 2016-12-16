package storageconnector

import (
	"fmt"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"regexp"
	"testing"
)

func TestPostgresStorage_SearchMeta(t *testing.T) {
	t.Run("FeatureContainingJava", func(t *testing.T) {
		name, _ := regexp.Compile("java")
		expected := 3

		mock.ExpectBegin()
		rows := sqlmock.NewRows([]string{"index", "name", "description", "author", "UpdatedAt", "dependencies", "snippet", "testSnippet"}).
			AddRow(1, "A-java-lein", "", "", "", "", "", "").
			AddRow(2, "B-java-node", "", "", "", "", "", "").
			AddRow(3, "java", "", "", "", "", "", "")

		mock.ExpectQuery("select").WillReturnRows(rows)
		features, err := pgStorage.SearchMeta(name)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(features)
		if len(features) != expected {
			t.Fatalf("Feature count should be %d but was %d", expected, len(features))
		}

		if features[0].Name != "java" {
			t.Errorf("Name of feature 0 should be 'java' but was '%s'", features[0].Name)
		}
		if features[1].Name != "java-python2" {
			t.Errorf("Name of feature 1 should be 'java-python2' but was '%s'", features[1].Name)
		}
	})
}

func TestPostgresStorage_GetMeta(t *testing.T) {
	t.Run("ExistingFeature", func(t *testing.T) {
		mock.ExpectBegin()
		rows := sqlmock.NewRows([]string{"index", "name", "description", "author", "UpdatedAt", "dependencies", "snippet", "testSnippet"}).
			AddRow(1, "java", "", "", "", "", "", "")
		mock.ExpectQuery("SELECT").WillReturnRows(rows)
		meta, err := pgStorage.GetMeta("java")
		if err != nil {
			t.Error(err)
		}

		if meta.Name != "java" {
			t.Errorf("Feature name should be 'java' but was '%s'", meta.Name)
		}
	})

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"index", "name", "description", "author", "UpdatedAt", "dependencies", "snippet", "testSnippet"})
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	t.Run("NonExistingFeature", func(t *testing.T) {
		_, err := pgStorage.GetMeta("reallynotafeature")
		if err == nil {
			t.Error("Error expected when getting metadata for non existing feature")
		}
	})
}

func TestPostgresStorage_Get(t *testing.T) {
	t.Run("ExistingFeatureWithoutSnippet", func(t *testing.T) {

		mock.ExpectBegin()
		rows := sqlmock.NewRows([]string{"index", "name", "description", "author", "UpdatedAt", "dependencies", "snippet", "testSnippet"}).
			AddRow(1, "java-python2", "", "", "", "", "", "")
		mock.ExpectQuery("SELECT").WillReturnRows(rows)
		feature, err := pgStorage.GetFeature("java-python2")
		if err != nil {
			t.Error(err)
		}

		if feature.Meta.Name != "java-python2" {
			t.Errorf("Feature name should be 'java-python2' but was '%s'", feature.Meta.Name)
		}

		if feature.Snippet != "" {
			t.Errorf("Feature snippet should be empty but was '%s'", feature.Snippet)
		}
	})

	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"index", "name", "description", "author", "UpdatedAt", "dependencies", "snippet", "testSnippet"}).
		AddRow(1, "java", "", "", "", "", "install java!!!", "")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	t.Run("ExistingFeatureWithSnippet", func(t *testing.T) {
		feature, err := pgStorage.GetFeature("java")
		if err != nil {
			t.Error(err)
		}

		if feature.Meta.Name != "java" {
			t.Errorf("Feature name should be 'java' but was '%s'", feature.Meta.Name)
		}

		if feature.Snippet == "" {
			t.Error("Feature snippet should not be empty", feature.Snippet)
		}
	})
	mock.ExpectBegin()
	rows = sqlmock.NewRows([]string{"index", "name", "description", "author", "UpdatedAt", "dependencies", "snippet", "testSnippet"})
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	t.Run("NonExistingFeature", func(t *testing.T) {
		_, err := pgStorage.GetFeature("reallynotafeature")
		if err == nil {
			t.Error("Error expected when getting non existing feature")
		}
	})
}

func TestPostgresStorage_Resolve(t *testing.T) {
	t.Run("FeatureWithoutDependencies", func(t *testing.T) {
		expected := 1

		list, features, err := pgStorage.Resolve("java")
		if err != nil {
			t.Error(err)
		}
		if len(list) != expected {
			t.Errorf("Feature count should be %d but was %d", expected, len(list))
		}

		if len(features) != expected {
			t.Errorf("Feature count should be %d but was %d", expected, len(features))
		}
		if _, ok := features["java"]; !ok {
			t.Error("Missing feature 'java'")
		}
	})

	t.Run("FeatureWithTwoDependencies", func(t *testing.T) {
		expected := 3

		list, features, err := pgStorage.Resolve("java-python2")
		if err != nil {
			t.Error(err)
		}

		if len(list) != expected {
			t.Errorf("Feature count should be %d but was %d", expected, len(list))
		}

		if len(features) != expected {
			t.Errorf("Feature count should be %d but was %d", expected, len(features))
		}

		if _, ok := features["java-python2"]; !ok {
			t.Error("Missing feature 'A-java-lein")
		}
		if _, ok := features["java"]; !ok {
			t.Error("Missing feature 'java")
		}
		if _, ok := features["python2"]; !ok {
			t.Error("Missing feature 'python2")
		}
	})
}
