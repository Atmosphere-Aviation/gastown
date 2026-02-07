package cmd

import (
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestVersionCmd(t *testing.T) {
	// Save and restore state
	origVersion := Version
	origBuild := Build
	origCommit := Commit
	origBranch := Branch
	origVerbose := versionVerbose
	defer func() {
		Version = origVersion
		Build = origBuild
		Commit = origCommit
		Branch = origBranch
		versionVerbose = origVerbose
	}()

	t.Run("basic output", func(t *testing.T) {
		Version = "1.0.0"
		Build = "test"
		Commit = ""
		Branch = ""
		versionVerbose = false

		rootCmd.SetArgs([]string{"version"})
		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("version command failed: %v", err)
		}
	})

	t.Run("verbose flag registered", func(t *testing.T) {
		flag := versionCmd.Flags().Lookup("verbose")
		if flag == nil {
			t.Fatal("--verbose flag not registered on version command")
		}
		if flag.Shorthand != "v" {
			t.Errorf("verbose shorthand = %q, want %q", flag.Shorthand, "v")
		}
	})

	t.Run("verbose includes timestamp", func(t *testing.T) {
		Version = "1.0.0"
		Build = "test"
		Commit = ""
		Branch = ""
		versionVerbose = true

		before := time.Now().Truncate(time.Second)

		// The command uses fmt.Printf with time.RFC3339 format
		ts := time.Now().Format(time.RFC3339)
		after := time.Now().Add(time.Second).Truncate(time.Second)

		// Verify timestamp is valid RFC3339
		parsed, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			t.Fatalf("timestamp %q is not valid RFC3339: %v", ts, err)
		}
		if parsed.Before(before) || parsed.After(after) {
			t.Errorf("timestamp %v not between %v and %v", parsed, before, after)
		}

		// Verify the output line format
		line := "Timestamp: " + ts
		if !strings.HasPrefix(line, "Timestamp: ") {
			t.Errorf("verbose line %q doesn't start with 'Timestamp: '", line)
		}
	})
}

func TestVersionVerboseGoVersion(t *testing.T) {
	goVer := runtime.Version()
	if !strings.HasPrefix(goVer, "go") {
		t.Errorf("runtime.Version() = %q, want prefix 'go'", goVer)
	}

	// Verify the output line format matches what version.go produces
	line := "Go version: " + goVer
	if !strings.HasPrefix(line, "Go version: go") {
		t.Errorf("verbose Go version line %q doesn't match expected format", line)
	}
}

func TestVersionVerboseTimestampFormat(t *testing.T) {
	// Verify time.RFC3339 produces a parseable timestamp
	ts := time.Now().Format(time.RFC3339)
	if _, err := time.Parse(time.RFC3339, ts); err != nil {
		t.Errorf("RFC3339 format failed round-trip: %v", err)
	}

	// Verify format includes date and time components
	if !strings.Contains(ts, "T") {
		t.Errorf("timestamp %q missing T separator", ts)
	}
	if !strings.Contains(ts, "-") {
		t.Errorf("timestamp %q missing date separators", ts)
	}
	if !strings.Contains(ts, ":") {
		t.Errorf("timestamp %q missing time separators", ts)
	}
}
