package validate

import "testing"

func TestValidateStorageShareDirectoryName(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "",
			Expected: false,
		},
		{
			Input:    "abc123",
			Expected: true,
		},
		{
			Input:    "123abc",
			Expected: true,
		},
		{
			Input:    "123-abc",
			Expected: true,
		},
		{
			Input:    "-123-abc",
			Expected: false,
		},
		{
			Input:    "123-abc-",
			Expected: false,
		},
		{
			Input:    "123--abc",
			Expected: false,
		},
	}

	for _, v := range testCases {
		t.Logf("[DEBUG] Test Input %q", v.Input)

		warnings, errors := StorageShareDirectoryName(v.Input, "name")
		if len(warnings) != 0 {
			t.Fatalf("Expected no warnings but got %d", len(warnings))
		}

		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
