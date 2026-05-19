package ipa

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type AnalysisResult struct {
	OK                bool               `json:"ok"`
	Command           string             `json:"command"`
	Target            string             `json:"target"`
	Platform          string             `json:"platform"`
	Encryption        EncryptionInfo     `json:"encryption"`
	Bundle            BundleInfo         `json:"bundle"`
	Permissions       PermissionsInfo    `json:"permissions"`
	URLSchemes        []string           `json:"urlSchemes"`
	QuerySchemes      []string           `json:"querySchemes"`
	BackgroundModes   []string           `json:"backgroundModes"`
	Frameworks        FrameworksInfo     `json:"frameworks"`
	Resources         ResourcesInfo      `json:"resources"`
	TechStackInference TechStackInference `json:"techStackInference"`
	LLMContext        LLMContext         `json:"llmContext"`
}

type EncryptionInfo struct {
	LikelyEncrypted bool   `json:"likelyEncrypted"`
	Reason          string `json:"reason"`
}

type BundleInfo struct {
	Name              string   `json:"name"`
	BundleID          string   `json:"bundleId"`
	Version           string   `json:"version"`
	Build             string   `json:"build"`
	MinimumOSVersion  string   `json:"minimumOSVersion"`
	DeviceFamilies    []string `json:"deviceFamilies"`
}

type PermissionDetail struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	Risk        string `json:"risk"`
}

type PermissionsInfo struct {
	PhotoLibrary      bool              `json:"photoLibrary"`
	Camera            bool              `json:"camera"`
	Microphone        bool              `json:"microphone"`
	Location          bool              `json:"location"`
	Contacts          bool              `json:"contacts"`
	Calendars         bool              `json:"calendars"`
	Tracking          bool              `json:"tracking"`
	SpeechRecognition bool              `json:"speechRecognition"`
	FaceID            bool              `json:"faceID"`
	Bluetooth         bool              `json:"bluetooth"`
	Motion            bool              `json:"motion"`
	Details           []PermissionDetail `json:"details"`
}

type FrameworksInfo struct {
	System         []string `json:"system"`
	ThirdPartyHints []string `json:"thirdPartyHints"`
}

type ResourcesInfo struct {
	AssetCatalogs int `json:"assetCatalogs"`
	Storyboards   int `json:"storyboards"`
	Nibs          int `json:"nibs"`
	StringsFiles  int `json:"stringsFiles"`
	JSONFiles     int `json:"jsonFiles"`
	MLModels      int `json:"mlModels"`
	Fonts         int `json:"fonts"`
	Images        int `json:"images"`
	AudioFiles    int `json:"audioFiles"`
	AppExtensions int `json:"appExtensions"`
}

type TechStackInference struct {
	PossibleLanguages []string `json:"possibleLanguages"`
	PossibleFrameworks []string `json:"possibleFrameworks"`
	PossibleSDKs      []string `json:"possibleSDKs"`
	Capabilities      []string `json:"capabilities"`
}

type LLMContext struct {
	Summary              string   `json:"summary"`
	RecommendedQuestions []string `json:"recommendedQuestions"`
}

var permissionKeys = map[string]struct {
	field string
	risk  string
}{
	"NSCameraUsageDescription":                  {"Camera", "high"},
	"NSMicrophoneUsageDescription":              {"Microphone", "high"},
	"NSPhotoLibraryUsageDescription":            {"PhotoLibrary", "medium"},
	"NSPhotoLibraryAddUsageDescription":         {"PhotoLibrary", "medium"},
	"NSLocationWhenInUseUsageDescription":       {"Location", "high"},
	"NSLocationAlwaysAndWhenInUseUsageDescription": {"Location", "high"},
	"NSContactsUsageDescription":                {"Contacts", "high"},
	"NSCalendarsUsageDescription":               {"Calendars", "medium"},
	"NSUserTrackingUsageDescription":            {"Tracking", "high"},
	"NSSpeechRecognitionUsageDescription":       {"SpeechRecognition", "medium"},
	"NSFaceIDUsageDescription":                  {"FaceID", "low"},
	"NSBluetoothAlwaysUsageDescription":         {"Bluetooth", "medium"},
	"NSMotionUsageDescription":                  {"Motion", "medium"},
}

var thirdPartyFrameworkHints = map[string]string{
	"Flutter":             "Flutter",
	"App":                 "Flutter",
	"Hermes":              "React Native",
	"React":               "React Native",
	"UnityFramework":      "Unity",
	"Capacitor":           "Capacitor/Hybrid",
	"Cordova":             "Cordova/Hybrid",
	"FirebaseCore":        "Firebase",
	"FirebaseAnalytics":   "Firebase",
	"FirebaseCrashlytics": "Firebase",
	"FirebaseMessaging":   "Firebase",
	"FirebaseRemoteConfig": "Firebase",
	"Sentry":              "Sentry",
	"RevenueCat":          "RevenueCat",
}

var systemFrameworkCapabilities = map[string]string{
	"Vision":      "可能使用 Apple Vision 框架",
	"CoreImage":   "可能做图像处理",
	"CoreML":      "可能使用本地 AI 模型",
	"Metal":       "可能使用 GPU 加速",
	"AVFoundation": "可能处理相机、视频、音频",
	"StoreKit":    "可能有内购",
	"ARKit":       "可能使用 AR 功能",
	"MapKit":      "可能使用地图",
	"WebKit":      "可能内嵌网页",
	"HealthKit":   "可能访问健康数据",
	"HomeKit":     "可能控制智能家居",
	"CloudKit":    "可能使用 iCloud 同步",
	"CoreLocation": "可能使用定位服务",
	"CoreBluetooth": "可能使用蓝牙",
	"Photos":      "可能访问相册",
	"Contacts":    "可能访问通讯录",
}

func Analyze(ipaPath string) (*AnalysisResult, error) {
	if _, err := os.Stat(ipaPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("IPA file not found: %s", ipaPath)
	}

	if !strings.HasSuffix(strings.ToLower(ipaPath), ".ipa") {
		return nil, fmt.Errorf("invalid file extension, expected .ipa: %s", ipaPath)
	}

	tmpDir, err := os.MkdirTemp("", "appinsight-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	cmd := exec.Command("unzip", "-q", "-o", ipaPath, "-d", tmpDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("failed to unzip IPA: %w\n%s", err, string(output))
	}

	payloadDir := filepath.Join(tmpDir, "Payload")
	entries, err := os.ReadDir(payloadDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read Payload directory: %w", err)
	}

	var appPath string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), ".app") {
			appPath = filepath.Join(payloadDir, entry.Name())
			break
		}
	}
	if appPath == "" {
		return nil, fmt.Errorf("no .app directory found in Payload")
	}

	plistData, err := readInfoPlist(appPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read Info.plist: %w", err)
	}

	bundle := extractBundleInfo(plistData)
	permissions := extractPermissions(plistData)
	urlSchemes := extractURLSchemes(plistData)
	querySchemes := extractQuerySchemes(plistData)
	backgroundModes := extractBackgroundModes(plistData)
	frameworks := scanFrameworks(appPath)
	resources := scanResources(appPath)
	techStack := inferTechStack(frameworks, resources, permissions)
	llmCtx := buildLLMContext(bundle, permissions, frameworks, techStack)

	return &AnalysisResult{
		OK:       true,
		Command:  "analyze-ipa",
		Target:   ipaPath,
		Platform: "ios",
		Encryption: EncryptionInfo{
			LikelyEncrypted: true,
			Reason:          "App Store IPA is usually encrypted. Binary-level analysis is limited.",
		},
		Bundle:            bundle,
		Permissions:       permissions,
		URLSchemes:        urlSchemes,
		QuerySchemes:      querySchemes,
		BackgroundModes:   backgroundModes,
		Frameworks:        frameworks,
		Resources:         resources,
		TechStackInference: techStack,
		LLMContext:        llmCtx,
	}, nil
}

func readInfoPlist(appPath string) (map[string]interface{}, error) {
	plistPath := filepath.Join(appPath, "Info.plist")
	if _, err := os.Stat(plistPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Info.plist not found at %s", plistPath)
	}

	cmd := exec.Command("plutil", "-convert", "json", "-o", "-", plistPath)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("plutil conversion failed: %w", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(output, &data); err != nil {
		return nil, fmt.Errorf("failed to parse Info.plist JSON: %w", err)
	}

	return data, nil
}

func getStringFromPlist(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}

func extractBundleInfo(plist map[string]interface{}) BundleInfo {
	families := []string{}
	if raw, ok := plist["UIDeviceFamily"]; ok {
		if arr, ok := raw.([]interface{}); ok {
			for _, v := range arr {
				switch fmt.Sprintf("%v", v) {
				case "1":
					families = append(families, "iPhone")
				case "2":
					families = append(families, "iPad")
				default:
					families = append(families, fmt.Sprintf("Unknown(%v)", v))
				}
			}
		}
	}

	return BundleInfo{
		Name:             getStringFromPlist(plist, "CFBundleDisplayName"),
		BundleID:         getStringFromPlist(plist, "CFBundleIdentifier"),
		Version:          getStringFromPlist(plist, "CFBundleShortVersionString"),
		Build:            getStringFromPlist(plist, "CFBundleVersion"),
		MinimumOSVersion: getStringFromPlist(plist, "MinimumOSVersion"),
		DeviceFamilies:   families,
	}
}

func extractPermissions(plist map[string]interface{}) PermissionsInfo {
	info := PermissionsInfo{
		Details: []PermissionDetail{},
	}

	for key, meta := range permissionKeys {
		desc := getStringFromPlist(plist, key)
		present := desc != ""
		switch meta.field {
		case "PhotoLibrary":
			if present {
				info.PhotoLibrary = true
			}
		case "Camera":
			info.Camera = present
		case "Microphone":
			info.Microphone = present
		case "Location":
			if present {
				info.Location = true
			}
		case "Contacts":
			info.Contacts = present
		case "Calendars":
			info.Calendars = present
		case "Tracking":
			info.Tracking = present
		case "SpeechRecognition":
			info.SpeechRecognition = present
		case "FaceID":
			info.FaceID = present
		case "Bluetooth":
			info.Bluetooth = present
		case "Motion":
			info.Motion = present
		}
		if present {
			info.Details = append(info.Details, PermissionDetail{
				Key:         key,
				Description: desc,
				Risk:        meta.risk,
			})
		}
	}

	return info
}

func extractURLSchemes(plist map[string]interface{}) []string {
	schemes := []string{}
	if raw, ok := plist["CFBundleURLTypes"]; ok {
		if arr, ok := raw.([]interface{}); ok {
			for _, item := range arr {
				if m, ok := item.(map[string]interface{}); ok {
					if urlSchemes, ok := m["CFBundleURLSchemes"]; ok {
						if sArr, ok := urlSchemes.([]interface{}); ok {
							for _, s := range sArr {
								if str, ok := s.(string); ok {
									schemes = append(schemes, str)
								}
							}
						}
					}
				}
			}
		}
	}
	return schemes
}

func extractQuerySchemes(plist map[string]interface{}) []string {
	schemes := []string{}
	if raw, ok := plist["LSApplicationQueriesSchemes"]; ok {
		if arr, ok := raw.([]interface{}); ok {
			for _, v := range arr {
				if s, ok := v.(string); ok {
					schemes = append(schemes, s)
				}
			}
		}
	}
	return schemes
}

func extractBackgroundModes(plist map[string]interface{}) []string {
	modes := []string{}
	if raw, ok := plist["UIBackgroundModes"]; ok {
		if arr, ok := raw.([]interface{}); ok {
			for _, v := range arr {
				if s, ok := v.(string); ok {
					modes = append(modes, s)
				}
			}
		}
	}
	return modes
}

func scanFrameworks(appPath string) FrameworksInfo {
	info := FrameworksInfo{
		System:         []string{},
		ThirdPartyHints: []string{},
	}

	fwPath := filepath.Join(appPath, "Frameworks")
	entries, err := os.ReadDir(fwPath)
	if err != nil {
		return info
	}

	seenThirdParty := map[string]bool{}
	for _, entry := range entries {
		name := entry.Name()
		baseName := strings.TrimSuffix(name, filepath.Ext(name))

		if isAppleFramework(baseName) {
			info.System = append(info.System, baseName)
		} else {
			if hint, ok := thirdPartyFrameworkHints[baseName]; ok {
				if !seenThirdParty[hint] {
					seenThirdParty[hint] = true
					info.ThirdPartyHints = append(info.ThirdPartyHints, hint)
				}
			} else {
				info.ThirdPartyHints = append(info.ThirdPartyHints, baseName)
			}
		}
	}

	return info
}

func isAppleFramework(name string) bool {
	applePrefixes := []string{"UI", "NS", "Core", "AV", "CF", "GL", "Metal", "Vision", "Map", "Store", "AR", "Web", "Health", "Home", "Cloud", "Photo", "Contact", "Event", "Local", "Audio", "Video", "Game", "Scene", "Swift", "Foundation", "Network", "Security", "Quartz", "Sprite", "Watch", "Pass", "Social", "Media", "Ad", "Car", "Class", "Intents", "Siri", "User", "Device", "Background", "Push", "Natural"}
	for _, prefix := range applePrefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}

func scanResources(appPath string) ResourcesInfo {
	info := ResourcesInfo{}

	filepath.Walk(appPath, func(path string, fi os.FileInfo, err error) error {
		if err != nil || fi.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		name := strings.ToLower(fi.Name())

		switch ext {
		case ".car":
			info.AssetCatalogs++
		case ".storyboardc":
			info.Storyboards++
		case ".nib":
			info.Nibs++
		case ".strings":
			info.StringsFiles++
		case ".json":
			info.JSONFiles++
		case ".mlmodelc":
			info.MLModels++
		case ".appex":
			info.AppExtensions++
		case ".ttf", ".otf", ".ttc":
			info.Fonts++
		case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".pdf", ".svg":
			info.Images++
		case ".mp3", ".wav", ".m4a", ".aac", ".caf", ".ogg":
			info.AudioFiles++
		}

		if strings.HasSuffix(name, ".appex") {
			info.AppExtensions++
		}

		return nil
	})

	return info
}

func inferTechStack(fw FrameworksInfo, res ResourcesInfo, perm PermissionsInfo) TechStackInference {
	langs := []string{}
	frameworks := []string{}
	sdks := []string{}
	capabilities := []string{}

	hasFlutter := false
	hasReactNative := false
	hasUnity := false
	hasHybrid := false

	for _, tp := range fw.ThirdPartyHints {
		switch tp {
		case "Flutter":
			hasFlutter = true
			frameworks = append(frameworks, "Flutter")
		case "React Native":
			hasReactNative = true
			frameworks = append(frameworks, "React Native")
		case "Unity":
			hasUnity = true
			frameworks = append(frameworks, "Unity")
		case "Capacitor/Hybrid", "Cordova/Hybrid":
			hasHybrid = true
			frameworks = append(frameworks, tp)
		case "Firebase":
			sdks = append(sdks, "Firebase")
		case "Sentry":
			sdks = append(sdks, "Sentry")
		case "RevenueCat":
			sdks = append(sdks, "RevenueCat")
		}
	}

	if hasFlutter {
		langs = append(langs, "Dart")
	}
	if hasReactNative {
		langs = append(langs, "JavaScript/TypeScript")
	}
	if hasUnity {
		langs = append(langs, "C#")
	}
	if hasHybrid {
		langs = append(langs, "HTML/CSS/JavaScript")
	}

	if !hasFlutter && !hasReactNative && !hasUnity && !hasHybrid {
		langs = append(langs, "Swift/Objective-C (推测)")
	}

	for _, sys := range fw.System {
		if cap, ok := systemFrameworkCapabilities[sys]; ok {
			capabilities = append(capabilities, cap)
		}
	}

	if res.MLModels > 0 {
		capabilities = append(capabilities, "可能使用 Core ML 本地模型")
	}

	if perm.Tracking {
		capabilities = append(capabilities, "可能做用户追踪/广告归因")
	}
	if perm.Camera && perm.Microphone {
		capabilities = append(capabilities, "可能做视频录制/通话")
	}
	if perm.Location {
		capabilities = append(capabilities, "可能使用定位服务")
	}

	if len(langs) == 0 {
		langs = []string{"Unknown"}
	}
	if len(frameworks) == 0 {
		frameworks = []string{"Native iOS (推测)"}
	}

	return TechStackInference{
		PossibleLanguages:  langs,
		PossibleFrameworks: frameworks,
		PossibleSDKs:       sdks,
		Capabilities:       capabilities,
	}
}

func buildLLMContext(bundle BundleInfo, perm PermissionsInfo, fw FrameworksInfo, tech TechStackInference) LLMContext {
	summary := fmt.Sprintf(
		"%s (Bundle ID: %s, Version: %s) 是一个 iOS 应用。技术栈推测: %v。使用了 %d 个系统 Framework 和 %d 个第三方 SDK。声明了 %d 项权限。",
		bundle.Name, bundle.BundleID, bundle.Version,
		tech.PossibleFrameworks,
		len(fw.System), len(fw.ThirdPartyHints),
		len(perm.Details),
	)

	questions := []string{}
	questions = append(questions, fmt.Sprintf("%s 使用了哪些第三方 SDK？各自的作用是什么？", bundle.Name))
	if len(tech.PossibleFrameworks) > 0 {
		questions = append(questions, fmt.Sprintf("%s 是否真的使用了 %s？如何确认？", bundle.Name, strings.Join(tech.PossibleFrameworks, "、")))
	}
	if perm.Tracking {
		questions = append(questions, fmt.Sprintf("%s 的用户追踪实现方式是什么？", bundle.Name))
	}
	if perm.Camera || perm.Microphone {
		questions = append(questions, fmt.Sprintf("%s 的相机/麦克风权限用于什么功能？", bundle.Name))
	}
	questions = append(questions, "结合 App Store 截图和用户评论，可以推断出哪些核心功能？")

	return LLMContext{
		Summary:              summary,
		RecommendedQuestions: questions,
	}
}

func ExtractStrings(binaryPath string, maxLines int) ([]string, error) {
	cmd := exec.Command("strings", binaryPath)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("strings command failed: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) > maxLines {
		lines = lines[:maxLines]
	}

	return lines, nil
}
