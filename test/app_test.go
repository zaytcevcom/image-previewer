package test

import (
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

func TestIntegration(t *testing.T) {

	t.Run("success JPG", func(t *testing.T) {

		statusCode, result, err := sendRequest("fill/200/200/nginx/images/1.jpg")

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, statusCode)
		require.Greater(t, len(result), 0)
	})

	t.Run("success PNG", func(t *testing.T) {

		statusCode, result, err := sendRequest("fill/200/200/nginx/images/2.png")

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, statusCode)
		require.Greater(t, len(result), 0)
	})

	t.Run("method not found", func(t *testing.T) {

		statusCode, result, err := sendRequest("not/found/path")

		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, statusCode)
		require.Equal(t, "404 Not Found\n", string(result))
	})

	t.Run("invalid url (zero values)", func(t *testing.T) {

		statusCode, result, err := sendRequest("fill/0/0/nginx/images/1.jpg")

		require.NoError(t, err)
		require.Equal(t, http.StatusBadGateway, statusCode)
		require.Equal(t, "502 Bad Gateway\n", string(result))
	})

	t.Run("invalid url (negative values)", func(t *testing.T) {

		statusCode, result, err := sendRequest("fill/-20/-30/nginx/images/1.jpg")

		require.NoError(t, err)
		require.Equal(t, http.StatusBadGateway, statusCode)
		require.Equal(t, "502 Bad Gateway\n", string(result))
	})

	t.Run("invalid url (huge values)", func(t *testing.T) {

		statusCode, result, err := sendRequest("fill/80000/80000/nginx/images/1.jpg")

		require.NoError(t, err)
		require.Equal(t, http.StatusBadGateway, statusCode)
		require.Equal(t, "502 Bad Gateway\n", string(result))
	})

	t.Run("host not found", func(t *testing.T) {

		statusCode, result, err := sendRequest("fill/200/200/nginx-not-found/images/2.png")

		require.NoError(t, err)
		require.Equal(t, http.StatusBadGateway, statusCode)
		require.Equal(t, "502 Bad Gateway\n", string(result))
	})

	t.Run("image not found", func(t *testing.T) {

		statusCode, result, err := sendRequest("fill/200/200/nginx/images/not-found.png")

		require.NoError(t, err)
		require.Equal(t, http.StatusBadGateway, statusCode)
		require.Equal(t, "502 Bad Gateway\n", string(result))
	})

	t.Run("incorrect file", func(t *testing.T) {

		statusCode, result, err := sendRequest("fill/200/200/nginx/images/test.txt")

		require.NoError(t, err)
		require.Equal(t, http.StatusBadGateway, statusCode)
		require.Equal(t, "502 Bad Gateway\n", string(result))
	})

}

func sendRequest(endpoint string) (int, []byte, error) {

	host := "http://previewer:80/" // ci/cd
	//host := "http://localhost:8880/" // local

	response, err := http.Get(host + endpoint)
	if err != nil {
		return 0, nil, err
	}

	result, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, nil, err
	}

	return response.StatusCode, result, nil
}
