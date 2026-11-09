package javaruntime

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestResolveUsesExplicitJavaWithoutProvisioning(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("teste especifico para execucao no Windows")
	}

	tempDir := t.TempDir()
	javaPath := writeFakeJavaBinary(t, tempDir, "java.cmd")

	resolver := Resolver{
		BaseURL: "http://127.0.0.1:0",
		HomeDir: filepath.Join(tempDir, "runner-home"),
		GOOS:    "windows",
		GOARCH:  "amd64",
	}

	resolved, err := resolver.Resolve(javaPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved != javaPath {
		t.Fatalf("expected explicit java path, got %s", resolved)
	}
}

func TestResolveDownloadsAndReusesManagedRuntimeOnWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("teste especifico para execucao no Windows")
	}

	archiveBytes := buildWindowsArchive(t)
	checksum := sha256Hex(archiveBytes)
	serverHits := 0

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		serverHits++
		switch {
		case strings.Contains(request.URL.Path, "/checksum/"):
			response.WriteHeader(http.StatusOK)
			_, _ = response.Write([]byte(checksum))
		case strings.Contains(request.URL.Path, "/binary/"):
			response.WriteHeader(http.StatusOK)
			_, _ = response.Write(archiveBytes)
		default:
			http.NotFound(response, request)
		}
	}))
	defer server.Close()

	tempDir := t.TempDir()
	resolver := Resolver{
		BaseURL: server.URL,
		HomeDir: filepath.Join(tempDir, "runner-home"),
		GOOS:    "windows",
		GOARCH:  "amd64",
		Now: func() time.Time {
			return time.Date(2026, time.June, 2, 12, 0, 0, 0, time.UTC)
		},
	}

	javaPath, err := resolver.Resolve("")
	if err != nil {
		t.Fatalf("unexpected error resolving managed java: %v", err)
	}
	if _, err := os.Stat(javaPath); err != nil {
		t.Fatalf("expected managed java to exist: %v", err)
	}

	metadataPath := filepath.Join(filepath.Dir(filepath.Dir(javaPath)), metadataFileName)
	if _, err := os.Stat(metadataPath); err != nil {
		t.Fatalf("expected runtime metadata to exist: %v", err)
	}

	firstHitCount := serverHits
	reusedPath, err := resolver.Resolve("")
	if err != nil {
		t.Fatalf("unexpected error reusing managed java: %v", err)
	}
	if reusedPath != javaPath {
		t.Fatalf("expected reused runtime path %s, got %s", javaPath, reusedPath)
	}
	if serverHits != firstHitCount {
		t.Fatalf("expected no extra downloads after reuse, got %d hits then %d", firstHitCount, serverHits)
	}
}

func TestResolveReturnsHelpfulErrorWhenChecksumFails(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("teste especifico para execucao no Windows")
	}

	archiveBytes := buildWindowsArchive(t)

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		switch {
		case strings.Contains(request.URL.Path, "/checksum/"):
			response.WriteHeader(http.StatusOK)
			_, _ = response.Write([]byte(strings.Repeat("0", 64)))
		case strings.Contains(request.URL.Path, "/binary/"):
			response.WriteHeader(http.StatusOK)
			_, _ = response.Write(archiveBytes)
		default:
			http.NotFound(response, request)
		}
	}))
	defer server.Close()

	resolver := Resolver{
		BaseURL: server.URL,
		HomeDir: filepath.Join(t.TempDir(), "runner-home"),
		GOOS:    "windows",
		GOARCH:  "amd64",
	}

	_, err := resolver.Resolve("")
	if err == nil {
		t.Fatalf("expected checksum validation error")
	}
	if !strings.Contains(err.Error(), "integridade") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResolveReturnsHelpfulErrorWhenArchiveHasNoJava(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("teste especifico para execucao no Windows")
	}

	archiveBytes := buildBrokenWindowsArchive(t)
	checksum := sha256Hex(archiveBytes)

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		switch {
		case strings.Contains(request.URL.Path, "/checksum/"):
			response.WriteHeader(http.StatusOK)
			_, _ = response.Write([]byte(checksum))
		case strings.Contains(request.URL.Path, "/binary/"):
			response.WriteHeader(http.StatusOK)
			_, _ = response.Write(archiveBytes)
		default:
			http.NotFound(response, request)
		}
	}))
	defer server.Close()

	resolver := Resolver{
		BaseURL: server.URL,
		HomeDir: filepath.Join(t.TempDir(), "runner-home"),
		GOOS:    "windows",
		GOARCH:  "amd64",
	}

	_, err := resolver.Resolve("")
	if err == nil {
		t.Fatalf("expected extraction error")
	}
	if !strings.Contains(err.Error(), "executavel Java extraido") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResolveSupportsTarGzExtractionOnLinux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("teste especifico para execucao no Linux")
	}

	archiveBytes := buildLinuxArchive(t)
	checksum := sha256Hex(archiveBytes)

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		switch {
		case strings.Contains(request.URL.Path, "/checksum/"):
			response.WriteHeader(http.StatusOK)
			_, _ = response.Write([]byte(checksum))
		case strings.Contains(request.URL.Path, "/binary/"):
			response.WriteHeader(http.StatusOK)
			_, _ = response.Write(archiveBytes)
		default:
			http.NotFound(response, request)
		}
	}))
	defer server.Close()

	resolver := Resolver{
		BaseURL: server.URL,
		HomeDir: filepath.Join(t.TempDir(), "runner-home"),
		GOOS:    "linux",
		GOARCH:  "amd64",
	}

	javaPath, err := resolver.Resolve("")
	if err != nil {
		t.Fatalf("unexpected error resolving linux managed java: %v", err)
	}
	if !strings.HasSuffix(javaPath, filepath.Join("bin", "java")) {
		t.Fatalf("unexpected linux java path: %s", javaPath)
	}
}

func buildWindowsArchive(t *testing.T) []byte {
	t.Helper()

	var buffer bytes.Buffer
	writer := zip.NewWriter(&buffer)
	fileWriter, err := writer.CreateHeader(&zip.FileHeader{
		Name:   "OpenJDK17U-jre_x64_windows_hotspot_17.0.18_8/bin/java.cmd",
		Method: zip.Deflate,
	})
	if err != nil {
		t.Fatalf("failed to create zip header: %v", err)
	}
	if _, err := fileWriter.Write([]byte("@echo off\r\nif \"%1\"==\"-version\" exit /b 0\r\n")); err != nil {
		t.Fatalf("failed to write fake java.cmd: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close zip writer: %v", err)
	}
	return buffer.Bytes()
}

func buildBrokenWindowsArchive(t *testing.T) []byte {
	t.Helper()

	var buffer bytes.Buffer
	writer := zip.NewWriter(&buffer)
	fileWriter, err := writer.Create("OpenJDK17U-jre_x64_windows_hotspot_17.0.18_8/README.txt")
	if err != nil {
		t.Fatalf("failed to create zip entry: %v", err)
	}
	if _, err := fileWriter.Write([]byte("missing java")); err != nil {
		t.Fatalf("failed to write broken archive: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close zip writer: %v", err)
	}
	return buffer.Bytes()
}

func buildLinuxArchive(t *testing.T) []byte {
	t.Helper()

	var buffer bytes.Buffer
	gzipWriter := gzip.NewWriter(&buffer)
	tarWriter := tar.NewWriter(gzipWriter)

	content := []byte("#!/bin/sh\nif [ \"$1\" = \"-version\" ]; then exit 0; fi\n")
	header := &tar.Header{
		Name: "OpenJDK17U-jre_x64_linux_hotspot_17.0.18_8/bin/java",
		Mode: 0o755,
		Size: int64(len(content)),
	}

	if err := tarWriter.WriteHeader(header); err != nil {
		t.Fatalf("failed to write tar header: %v", err)
	}
	if _, err := tarWriter.Write(content); err != nil {
		t.Fatalf("failed to write tar content: %v", err)
	}
	if err := tarWriter.Close(); err != nil {
		t.Fatalf("failed to close tar writer: %v", err)
	}
	if err := gzipWriter.Close(); err != nil {
		t.Fatalf("failed to close gzip writer: %v", err)
	}
	return buffer.Bytes()
}

func sha256Hex(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func writeFakeJavaBinary(t *testing.T, dir string, name string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	content := "@echo off\r\nexit /b 0\r\n"
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatalf("failed to create fake java: %v", err)
	}
	return path
}
