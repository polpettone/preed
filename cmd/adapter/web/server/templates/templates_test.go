package templates

import (
	"testing"
)

func Test_NewTestTemplateCache_happy_path(t *testing.T) {

	cache, err := NewTemplateCache("test-template-data/sub1")

	if err != nil {
		t.Errorf("wanted not error got %s", err)
	}

	if len(cache) != 1 {
		t.Errorf("wanted len %d got %d", 1, len(cache))
	}

}

func Test_NewTemplateCache_error_if_no_partial(t *testing.T) {
	_, err := NewTemplateCache("test-template-data/sub2")
	if err == nil {
		t.Errorf("wanted  error got none")
	}
}

func Test_NewTemplateCache_error_if_no_layout(t *testing.T) {
	_, err := NewTemplateCache("test-template-data/sub3")
	if err == nil {
		t.Errorf("wanted  error got none")
	}
}
