package disc_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/urfave/cli/v2/disc"
)

func TestSlug(t *testing.T) {
	slug := &disc.Slug{}
	func(ctx context.Context) {
		_, ok := ctx.(*disc.Slug)
		if !ok {
			t.Error("I don't know how type checking works")
			return
		}
	}(slug)
	fmt.Println("stu")
}
