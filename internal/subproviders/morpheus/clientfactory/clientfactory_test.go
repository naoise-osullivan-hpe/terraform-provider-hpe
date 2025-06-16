// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package clientfactory_test

import (
	"context"
	"crypto/x509"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/clientfactory"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/model"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/testhelpers"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestSecureTLS(t *testing.T) {
	defer testhelpers.RecordResult(t)
	server := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			// Simulate a simple 200 OK response
			w.WriteHeader(http.StatusOK)
		}))
	defer server.Close()

	m := model.SubModel{
		URL:         types.StringValue(server.URL),
		Username:    types.StringValue("user"),
		Password:    types.StringValue("secret"),
		AccessToken: types.StringValue("token"),
		Insecure:    types.BoolValue(false),
	}
	cf := clientfactory.New(m)
	c, err := cf.NewClient(context.Background())
	if err != nil {
		t.Fatal("Failed to create client", err)
	}
	u := c.UsersAPI.GetUser(context.Background(), 1)
	_, _, err = u.Execute()
	if err == nil {
		t.Fatal("Failed to raise error", err)
	}
	var certErr x509.UnknownAuthorityError
	if !errors.As(err, &certErr) {
		t.Fatalf("Expected UnknownAuthorityError, got: %v", err)
	}
}

func TestInsecureTLS(t *testing.T) {
	defer testhelpers.RecordResult(t)
	server := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			// Simulate a simple 200 OK response
			w.WriteHeader(http.StatusOK)
		}))
	defer server.Close()

	m := model.SubModel{
		URL:         types.StringValue(server.URL),
		Username:    types.StringValue("user"),
		Password:    types.StringValue("secret"),
		AccessToken: types.StringValue("token"),
		Insecure:    types.BoolValue(true),
	}
	cf := clientfactory.New(m)
	c, err := cf.NewClient(context.Background())
	if err != nil {
		t.Fatal("Failed to create client", err)
	}
	u := c.UsersAPI.GetUser(context.Background(), 1)
	_, _, err = u.Execute()
	if err != nil {
		t.Fatal("Unexpected error", err)
	}
}
