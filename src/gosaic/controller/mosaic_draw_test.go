package controller

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMosaicDraw(t *testing.T) {
	env, out, err := setupControllerTest()
	if err != nil {
		t.Fatalf("Error getting test environment: %s\n", err.Error())
	}
	defer env.Close()

	dir, err := ioutil.TempDir("", "gosaic_test_mosaic_draw")
	if err != nil {
		t.Fatal("Error getting temp dir for mosaic draw test: %s\n", err.Error())
	}
	defer os.RemoveAll(dir)

	err = Index(env, []string{"testdata", "../service/testdata"})
	if err != nil {
		t.Fatalf("Error indexing images: %s\n", err.Error())
	}

	cover, macro := MacroAspect(env, "testdata/jumping_bunny.jpg", 1000, 1000, 2, 3, 10, "", "")
	if cover == nil || macro == nil {
		t.Fatal("Failed to create cover or macro")
	}

	err = PartialAspect(env, macro.Id)
	if err != nil {
		t.Fatalf("Error building partial aspects: %s\n", err.Error())
	}

	err = Compare(env, macro.Id)
	if err != nil {
		t.Fatalf("Comparing images: %s\n", err.Error())
	}

	mosaic := MosaicBuild(env, "Jumping Bunny", "best", macro.Id, 0)
	if mosaic == nil {
		t.Fatal("Failed to build mosaic")
	}

	err = MosaicDraw(env, mosaic.Id, filepath.Join(dir, "jumping_bunny_mosaic.jpg"))
	if err != nil {
		t.Fatalf("Error drawing mosaic: %s\n", err.Error())
	}

	result := out.String()
	expect := []string{
		"Drawing 150 mosaic partials...",
	}

	for _, e := range expect {
		if !strings.Contains(result, e) {
			t.Fatalf("Expected result to contain '%s', but it did not\n", e)
		}
	}

	for _, ne := range []string{"fail", "error"} {
		if strings.Contains(strings.ToLower(result), ne) {
			t.Fatalf("Did not expect result to contain: %s, but it did\n", ne)
		}
	}
}
