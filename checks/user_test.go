package checks

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

func TestChecksCheckUser(t *testing.T) {

	check := &CheckUser{}

	cr, err := check.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	results, ok := cr.([]*resultUser)
	if !ok {
		t.Fatal("False type")
	}

	js, _ := json.Marshal(results)
	fmt.Println(string(js))
}
