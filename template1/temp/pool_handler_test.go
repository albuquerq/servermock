package temp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStackHandler(t *testing.T) {

	t.Run("nil stack handler", func(t *testing.T) {
		var h *StackHandler

		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "", nil)
		require.NoError(t, err)

		h.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("nil base handler", func(t *testing.T) {
		h := NewStackHandler(nil)

		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "", nil)
		require.NoError(t, err)

		h.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("schedule 2 items", func(t *testing.T) {
		h := NewStackHandler(func(w http.ResponseWriter, _ *http.Request) {
			fmt.Fprint(w, "hello with 0")
		})

		h.Push(func(w http.ResponseWriter, _ *http.Request) {
			fmt.Fprint(w, "hello with 1")
		})

		h.Push(func(w http.ResponseWriter, _ *http.Request) {
			fmt.Fprint(w, "hello with 2")
		})

		for i := 0; i < 5; i++ {
			w := httptest.NewRecorder()
			r, err := http.NewRequest(http.MethodGet, "", nil)
			require.NoError(t, err)

			h.ServeHTTP(w, r)

			assert.Equal(t, http.StatusOK, w.Result().StatusCode)
			t.Log(w.Body.String())
		}

	})

}
