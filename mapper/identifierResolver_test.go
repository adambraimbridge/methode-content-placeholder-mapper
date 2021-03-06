package mapper

import (
	"net/http"
	"strings"
	"testing"

	"github.com/Financial-Times/methode-content-placeholder-mapper/model"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestResolve_Ok(t *testing.T) {
	mockClient := new(model.MockDocStoreClient)
	mockClient.On("ContentQuery", "http://api.ft.com/system/FT-LABS-WP-1-24", "http://ftalphaville.ft.com/?p=2193913", "tid_1").Return(http.StatusMovedPermanently, "http://api.ft.com/content/5414b08f-5ae1-3bd6-9901-a9dd1bf9db03", nil)

	resolver := NewHttpIResolver(mockClient, map[string]string{"ftalphaville.ft.com": "FT-LABS-WP-1-24"})
	uuid, err := resolver.ResolveIdentifier("http://ftalphaville.ft.com/?p=2193913", "2193913", "tid_1")

	assert.NoError(t, err, "Should resolve fine.")
	assert.Equal(t, "5414b08f-5ae1-3bd6-9901-a9dd1bf9db03", uuid)
}

func TestResolve_NotInMap(t *testing.T) {
	mockClient := new(model.MockDocStoreClient)
	mockClient.On("ContentQuery", "http://api.ft.com/system/FT-LABS-WP-1-24", "http://ftalphaville.ft.com/?p=2193913", "tid_1").Return(http.StatusMovedPermanently, "http://api.ft.com/content/5414b08f-5ae1-3bd6-9901-a9dd1bf9db03", nil)

	resolver := NewHttpIResolver(mockClient, map[string]string{})
	_, err := resolver.ResolveIdentifier("http://ftalphaville.ft.com/?p=2193913", "2193913", "tid_1")

	assert.True(t, strings.Contains(err.Error(), "couldn't find authority in mapping table"))
}

func TestResolve_InvalidUuid(t *testing.T) {
	mockClient := new(model.MockDocStoreClient)
	mockClient.On("ContentQuery", "http://api.ft.com/system/FT-LABS-WP-1-24", "http://ftalphaville.ft.com/?p=2193913", "tid_1").Return(http.StatusMovedPermanently, "http://api.ft.com/content/5414b08f-xxxxx", nil)

	resolver := NewHttpIResolver(mockClient, map[string]string{"ftalphaville.ft.com": "FT-LABS-WP-1-24"})
	_, err := resolver.ResolveIdentifier("http://ftalphaville.ft.com/?p=2193913", "2193913", "tid_1")

	assert.True(t, strings.Contains(err.Error(), "invalid uuid"))
}

func TestResolve_InvalidLocation(t *testing.T) {
	mockClient := new(model.MockDocStoreClient)
	mockClient.On("ContentQuery", "http://api.ft.com/system/FT-LABS-WP-1-24", "http://ftalphaville.ft.com/?p=2193913", "tid_1").Return(http.StatusMovedPermanently, "wrong", nil)

	resolver := NewHttpIResolver(mockClient, map[string]string{"ftalphaville.ft.com": "FT-LABS-WP-1-24"})
	_, err := resolver.ResolveIdentifier("http://ftalphaville.ft.com/?p=2193913", "2193913", "tid_1")

	assert.True(t, strings.Contains(err.Error(), "invalid FT URI"))
}

func TestResolve_NotFound(t *testing.T) {
	mockClient := new(model.MockDocStoreClient)
	mockClient.On("ContentQuery", "http://api.ft.com/system/FT-LABS-WP-1-24", "http://ftalphaville.ft.com/?p=2193913", "tid_1").Return(http.StatusNotFound, "", nil)

	resolver := NewHttpIResolver(mockClient, map[string]string{"ftalphaville.ft.com": "FT-LABS-WP-1-24"})
	_, err := resolver.ResolveIdentifier("http://ftalphaville.ft.com/?p=2193913", "2193913", "tid_1")

	assert.True(t, strings.Contains(err.Error(), "404"))
}

func TestResolve_NetFail(t *testing.T) {
	mockClient := new(model.MockDocStoreClient)
	mockClient.On("ContentQuery", "http://api.ft.com/system/FT-LABS-WP-1-24", "http://ftalphaville.ft.com/?p=2193913", "tid_1").Return(-1, "", errors.New("Couldn't make HTTP call"))

	resolver := NewHttpIResolver(mockClient, map[string]string{"ftalphaville.ft.com": "FT-LABS-WP-1-24"})
	_, err := resolver.ResolveIdentifier("http://ftalphaville.ft.com/?p=2193913", "2193913", "tid_1")

	assert.Equal(t, "Couldn't make HTTP call", err.Error())
}

func TestContentExists_OK(t *testing.T) {
	mockClient := new(model.MockDocStoreClient)
	mockClient.On("ContentExists", "111", "tid_1").Return(true, nil)

	resolver := NewHttpIResolver(mockClient, map[string]string{})
	found, err := resolver.ContentExists("111", "tid_1")
	assert.NoError(t, err)
	assert.True(t, found)
}

func TestContentExists_NOK(t *testing.T) {
	mockClient := new(model.MockDocStoreClient)
	mockClient.On("ContentExists", "111", "tid_1").Return(false, errors.New("any"))

	resolver := NewHttpIResolver(mockClient, map[string]string{})
	found, err := resolver.ContentExists("111", "tid_1")
	assert.Error(t, err)
	assert.False(t, found)
}
