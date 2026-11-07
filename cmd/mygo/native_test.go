package main

import (
	"strings"
	"testing"
)

func TestMergeModulePlists(t *testing.T) {
	a, err := parsePlist([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0"><dict>
	<key>NSSpeechRecognitionUsageDescription</key><string>Recognize speech.</string>
	<key>com.apple.security.application-groups</key><array><string>group.a</string></array>
	<key>Count</key><integer>3</integer>
	<key>On</key><true/>
</dict></plist>`))
	if err != nil {
		t.Fatal(err)
	}
	b, err := parsePlist([]byte(`<plist version="1.0"><dict>
	<key>com.apple.security.application-groups</key><array><string>group.a</string><string>group.b</string></array>
	<key>NSSpeechRecognitionUsageDescription</key><string>Recognize speech.</string>
</dict></plist>`))
	if err != nil {
		t.Fatal(err)
	}
	merged := map[string]any{}
	owners := map[string]string{}
	if err := mergePlist(merged, a, "modules/a", owners, ""); err != nil {
		t.Fatal(err)
	}
	if err := mergePlist(merged, b, "modules/b", owners, ""); err != nil {
		t.Fatal(err)
	}
	if groups := merged["com.apple.security.application-groups"].([]any); len(groups) != 2 {
		t.Errorf("arrays merge by their values: %v", groups)
	}
	if merged["Count"] != 3.0 || merged["On"] != true {
		t.Errorf("merged = %v", merged)
	}
	c := map[string]any{"NSSpeechRecognitionUsageDescription": "Other words."}
	err = mergePlist(merged, c, "modules/c", owners, "")
	if err == nil || !strings.Contains(err.Error(), "modules/a and modules/c") {
		t.Errorf("conflicting modules: %v", err)
	}
}
