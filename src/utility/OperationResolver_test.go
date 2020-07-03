package utility_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x1n13y84issmd42/oasis/src/contract"
	"github.com/x1n13y84issmd42/oasis/src/log"
	"github.com/x1n13y84issmd42/oasis/src/utility"
)

func Test_OperationResolver(T *testing.T) {
	log := log.NewPlain(0)
	spec := utility.Load("../../spec/test/oas3_ok.yaml", log)
	resolver := utility.NewOperationResolver(spec, log)

	getOpsIDs := func(ops []contract.Operation) []string {
		res := []string{}
		for _, op := range ops {
			res = append(res, op.ID())
		}
		return res
	}

	T.Run("*", func(T *testing.T) {
		expected := []string{
			"deleteUser",
			"getPetById",
			"getUserByName",
			"updateUser",
		}

		masks := []string{
			"*",
			"get*",
			"*",
		}

		assert.Equal(T, expected, getOpsIDs(resolver.Resolve(masks)))
	})

	T.Run("*User", func(T *testing.T) {
		expected := []string{
			"deleteUser",
			"updateUser",
		}

		masks := []string{
			"*User",
		}

		assert.Equal(T, expected, getOpsIDs(resolver.Resolve(masks)))
	})

	T.Run("get*", func(T *testing.T) {
		expected := []string{
			"getPetById",
			"getUserByName",
		}

		masks := []string{
			"get*",
		}

		assert.Equal(T, expected, getOpsIDs(resolver.Resolve(masks)))
	})
}
