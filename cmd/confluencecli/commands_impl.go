package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/darthkoax/confluencecli/internal/api"
	"github.com/darthkoax/confluencecli/internal/client"
	"github.com/darthkoax/confluencecli/internal/config"
)

func getServices(configPath string) *api.Services {
	if configPath == "" {
		configPath = config.DefaultConfigPath()
	}
	cfg, err := loadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}
	c, err := createClient(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
		os.Exit(1)
	}
	return api.NewServices(c)
}

func loadConfig(configPath string) (*config.Config, error) {
	return config.Load(configPath)
}

func createClient(cfg *config.Config) (*client.Client, error) {
	return client.New(cfg)
}

func outputJSON(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}

func parseFlag(args []string, flag string) string {
	for i := 0; i < len(args)-1; i++ {
		if args[i] == flag {
			return args[i+1]
		}
	}
	return ""
}

func parseFlagInt(args []string, flag string, defaultVal int) int {
	val := parseFlag(args, flag)
	if val == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return i
}

func split(s, sep string) []string {
	var result []string
	for {
		idx := indexOf(s, sep)
		if idx < 0 {
			result = append(result, s)
			break
		}
		result = append(result, s[:idx])
		s = s[idx+len(sep):]
	}
	return result
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func handleContentGet(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: content ID required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	expand := parseFlag(args, "--expand")
	var expandList []string
	if expand != "" {
		expandList = split(expand, ",")
	}
	content, err := services.Content.Get(context.Background(), args[0], expandList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(content)
}

func handleContentList(args []string) {
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	contentType := parseFlag(args, "--type")
	spaceKey := parseFlag(args, "--space-key")
	title := parseFlag(args, "--title")
	start := parseFlagInt(args, "--start", 0)
	limit := parseFlagInt(args, "--limit", 25)
	expand := parseFlag(args, "--expand")
	var expandList []string
	if expand != "" {
		expandList = split(expand, ",")
	}
	result, err := services.Content.GetAll(context.Background(), contentType, spaceKey, title, start, limit, expandList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}

func handleContentCreate(args []string) {
	jsonPayload := parseFlag(args, "--json")
	if jsonPayload == "" {
		fmt.Fprintf(os.Stderr, "Error: --json flag required\n")
		os.Exit(1)
	}
	var content api.Content
	if err := json.Unmarshal([]byte(jsonPayload), &content); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	result, err := services.Content.Create(context.Background(), &content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}

func handleContentUpdate(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: content ID required\n")
		os.Exit(1)
	}
	jsonPayload := parseFlag(args, "--json")
	if jsonPayload == "" {
		fmt.Fprintf(os.Stderr, "Error: --json flag required\n")
		os.Exit(1)
	}
	var content api.Content
	if err := json.Unmarshal([]byte(jsonPayload), &content); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	result, err := services.Content.Update(context.Background(), args[0], &content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}

func handleContentDelete(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: content ID required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	err := services.Content.Delete(context.Background(), args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Content %s deleted successfully\n", args[0])
}

func handleContentHistory(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: content ID required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	expand := parseFlag(args, "--expand")
	var expandList []string
	if expand != "" {
		expandList = split(expand, ",")
	}
	history, err := services.Content.GetHistory(context.Background(), args[0], expandList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(history)
}

func handleContentChildren(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: content ID required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	expand := parseFlag(args, "--expand")
	var expandList []string
	if expand != "" {
		expandList = split(expand, ",")
	}
	children, err := services.Content.GetChildren(context.Background(), args[0], expandList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(children)
}

func handleContentDescendants(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: content ID required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	expand := parseFlag(args, "--expand")
	var expandList []string
	if expand != "" {
		expandList = split(expand, ",")
	}
	descendants, err := services.Content.GetDescendants(context.Background(), args[0], expandList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(descendants)
}

func handleContentVersions(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: content ID required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	start := parseFlagInt(args, "--start", 0)
	limit := parseFlagInt(args, "--limit", 25)
	result, err := services.Content.GetVersions(context.Background(), args[0], start, limit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}

func handleContentComments(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: content ID required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	result, err := services.Content.GetComments(context.Background(), args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}

func handleContentAttachments(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: content ID required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	result, err := services.Content.GetAttachments(context.Background(), args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}

func handleSpaceList(args []string) {
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	spaceType := parseFlag(args, "--type")
	status := parseFlag(args, "--status")
	start := parseFlagInt(args, "--start", 0)
	limit := parseFlagInt(args, "--limit", 25)
	expand := parseFlag(args, "--expand")
	var expandList []string
	if expand != "" {
		expandList = split(expand, ",")
	}
	result, err := services.Spaces.GetAll(context.Background(), spaceType, status, start, limit, expandList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}

func handleSpaceGet(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: space key required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	expand := parseFlag(args, "--expand")
	var expandList []string
	if expand != "" {
		expandList = split(expand, ",")
	}
	space, err := services.Spaces.Get(context.Background(), args[0], expandList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(space)
}

func handleSpaceCreate(args []string) {
	jsonPayload := parseFlag(args, "--json")
	if jsonPayload == "" {
		fmt.Fprintf(os.Stderr, "Error: --json flag required\n")
		os.Exit(1)
	}
	var space api.Space
	if err := json.Unmarshal([]byte(jsonPayload), &space); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	result, err := services.Spaces.Create(context.Background(), &space)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}

func handleSpaceDelete(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: space key required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	err := services.Spaces.Delete(context.Background(), args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Space %s deleted successfully\n", args[0])
}

func handleSpaceContent(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: space key required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	contentType := parseFlag(args, "--type")
	start := parseFlagInt(args, "--start", 0)
	limit := parseFlagInt(args, "--limit", 25)
	expand := parseFlag(args, "--expand")
	var expandList []string
	if expand != "" {
		expandList = split(expand, ",")
	}
	result, err := services.Spaces.GetContent(context.Background(), args[0], contentType, start, limit, expandList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}

func handleSearchCQL(args []string) {
	cql := parseFlag(args, "--cql")
	if cql == "" {
		fmt.Fprintf(os.Stderr, "Error: --cql required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	start := parseFlagInt(args, "--start", 0)
	limit := parseFlagInt(args, "--limit", 25)
	expand := parseFlag(args, "--expand")
	var expandList []string
	if expand != "" {
		expandList = split(expand, ",")
	}
	result, err := services.Search.Search(context.Background(), cql, start, limit, expandList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}

func handleSearchUser(args []string) {
	cql := parseFlag(args, "--cql")
	if cql == "" {
		fmt.Fprintf(os.Stderr, "Error: --cql required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	start := parseFlagInt(args, "--start", 0)
	limit := parseFlagInt(args, "--limit", 25)
	result, err := services.Search.SearchUser(context.Background(), cql, start, limit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}

func handleUserCurrent(args []string) {
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	expand := parseFlag(args, "--expand")
	var expandList []string
	if expand != "" {
		expandList = split(expand, ",")
	}
	user, err := services.Users.GetCurrent(context.Background(), expandList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(user)
}

func handleUserAnonymous(args []string) {
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	expand := parseFlag(args, "--expand")
	var expandList []string
	if expand != "" {
		expandList = split(expand, ",")
	}
	user, err := services.Users.GetAnonymous(context.Background(), expandList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(user)
}

func handleUserGet(args []string) {
	username := parseFlag(args, "--username")
	key := parseFlag(args, "--key")
	accountID := parseFlag(args, "--account-id")
	if username == "" && key == "" && accountID == "" {
		fmt.Fprintf(os.Stderr, "Error: --username, --key, or --account-id required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	expand := parseFlag(args, "--expand")
	var expandList []string
	if expand != "" {
		expandList = split(expand, ",")
	}
	user, err := services.Users.Get(context.Background(), username, key, accountID, expandList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(user)
}

func handleUserUnknown(args []string) {
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	expand := parseFlag(args, "--expand")
	var expandList []string
	if expand != "" {
		expandList = split(expand, ",")
	}
	user, err := services.Users.GetUnknown(context.Background(), expandList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(user)
}

func handleUserEmail(args []string) {
	email := parseFlag(args, "--email")
	if email == "" {
		fmt.Fprintf(os.Stderr, "Error: --email required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	expand := parseFlag(args, "--expand")
	var expandList []string
	if expand != "" {
		expandList = split(expand, ",")
	}
	user, err := services.Users.GetUserByEmail(context.Background(), email, expandList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(user)
}

func handleGroupGet(args []string) {
	groupName := parseFlag(args, "--groupname")
	if groupName == "" {
		if len(args) == 0 {
			fmt.Fprintf(os.Stderr, "Error: --groupname required\n")
			os.Exit(1)
		}
		groupName = args[0]
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	expand := parseFlag(args, "--expand")
	var expandList []string
	if expand != "" {
		expandList = split(expand, ",")
	}
	group, err := services.Groups.Get(context.Background(), groupName, expandList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(group)
}

func handleGroupMembers(args []string) {
	groupName := parseFlag(args, "--groupname")
	if groupName == "" {
		if len(args) == 0 {
			fmt.Fprintf(os.Stderr, "Error: --groupname required\n")
			os.Exit(1)
		}
		groupName = args[0]
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	start := parseFlagInt(args, "--start", 0)
	limit := parseFlagInt(args, "--limit", 25)
	members, err := services.Groups.GetMembers(context.Background(), groupName, start, limit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(members)
}

func handleGroupAddMember(args []string) {
	groupName := parseFlag(args, "--groupname")
	if groupName == "" {
		fmt.Fprintf(os.Stderr, "Error: --groupname required\n")
		os.Exit(1)
	}
	username := parseFlag(args, "--username")
	if username == "" {
		fmt.Fprintf(os.Stderr, "Error: --username required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	user := &api.User{Username: username}
	err := services.Groups.AddMember(context.Background(), groupName, user)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("User %s added to group %s successfully\n", username, groupName)
}

func handleGroupRemoveMember(args []string) {
	groupName := parseFlag(args, "--groupname")
	if groupName == "" {
		fmt.Fprintf(os.Stderr, "Error: --groupname required\n")
		os.Exit(1)
	}
	username := parseFlag(args, "--username")
	if username == "" {
		fmt.Fprintf(os.Stderr, "Error: --username required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	err := services.Groups.RemoveMember(context.Background(), groupName, username)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("User %s removed from group %s successfully\n", username, groupName)
}

func handleSettingsSystemInfo(args []string) {
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	info, err := services.Settings.GetSystemInfo(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(info)
}

func handleSettingsThemeGet(args []string) {
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	theme, err := services.Settings.GetTheme(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(theme)
}

func handleSettingsThemeSet(args []string) {
	themeKey := parseFlag(args, "--theme-key")
	if themeKey == "" {
		fmt.Fprintf(os.Stderr, "Error: --theme-key required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	err := services.Settings.SetTheme(context.Background(), themeKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Theme set to %s successfully\n", themeKey)
}

func handleSettingsLookAndFeelGet(args []string) {
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	spaceKey := parseFlag(args, "--space-key")
	laf, err := services.Settings.GetLookAndFeel(context.Background(), spaceKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(laf)
}

func handleSettingsLookAndFeelSet(args []string) {
	jsonPayload := parseFlag(args, "--json")
	if jsonPayload == "" {
		fmt.Fprintf(os.Stderr, "Error: --json flag required\n")
		os.Exit(1)
	}
	var laf api.LookAndFeel
	if err := json.Unmarshal([]byte(jsonPayload), &laf); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	err := services.Settings.SetLookAndFeel(context.Background(), &laf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Look and feel updated successfully")
}

func handleSettingsLookAndFeelReset(args []string) {
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	spaceKey := parseFlag(args, "--space-key")
	err := services.Settings.ResetLookAndFeel(context.Background(), spaceKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Look and feel reset successfully")
}

func handleAuditGet(args []string) {
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	startDate := parseFlag(args, "--start-date")
	endDate := parseFlag(args, "--end-date")
	searchString := parseFlag(args, "--search")
	start := parseFlagInt(args, "--start", 0)
	limit := parseFlagInt(args, "--limit", 25)
	result, err := services.Audit.Get(context.Background(), startDate, endDate, searchString, start, limit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}

func handleAuditSince(args []string) {
	number := parseFlag(args, "--number")
	if number == "" {
		fmt.Fprintf(os.Stderr, "Error: --number required (e.g., 1m, 7d, 24h)\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	searchString := parseFlag(args, "--search")
	start := parseFlagInt(args, "--start", 0)
	limit := parseFlagInt(args, "--limit", 25)
	result, err := services.Audit.GetSince(context.Background(), number, searchString, start, limit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}

func handleAuditRetention(args []string) {
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	if len(args) > 0 && args[0] == "set" {
		number := parseFlagInt(args, "--number", 0)
		units := parseFlag(args, "--units")
		if number == 0 || units == "" {
			fmt.Fprintf(os.Stderr, "Error: --number and --units required\n")
			os.Exit(1)
		}
		err := services.Audit.SetRetentionPeriod(context.Background(), number, units)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Audit retention period updated successfully")
		return
	}
	retention, err := services.Audit.GetRetentionPeriod(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(retention)
}

func handleAuditExport(args []string) {
	startDate := parseFlag(args, "--start-date")
	endDate := parseFlag(args, "--end-date")
	if startDate == "" || endDate == "" {
		fmt.Fprintf(os.Stderr, "Error: --start-date and --end-date required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	result, err := services.Audit.Export(context.Background(), startDate, endDate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}

func handleTemplateListContent(args []string) {
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	spaceKey := parseFlag(args, "--space-key")
	start := parseFlagInt(args, "--start", 0)
	limit := parseFlagInt(args, "--limit", 25)
	expand := parseFlag(args, "--expand")
	var expandList []string
	if expand != "" {
		expandList = split(expand, ",")
	}
	result, err := services.Templates.GetContentTemplates(context.Background(), spaceKey, start, limit, expandList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}

func handleTemplateGetContent(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: template ID required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	expand := parseFlag(args, "--expand")
	var expandList []string
	if expand != "" {
		expandList = split(expand, ",")
	}
	tmpl, err := services.Templates.GetContentTemplate(context.Background(), args[0], expandList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(tmpl)
}

func handleTemplateCreateContent(args []string) {
	jsonPayload := parseFlag(args, "--json")
	if jsonPayload == "" {
		fmt.Fprintf(os.Stderr, "Error: --json flag required\n")
		os.Exit(1)
	}
	var tmpl api.ContentTemplate
	if err := json.Unmarshal([]byte(jsonPayload), &tmpl); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	result, err := services.Templates.CreateContentTemplate(context.Background(), &tmpl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}

func handleTemplateUpdateContent(args []string) {
	jsonPayload := parseFlag(args, "--json")
	if jsonPayload == "" {
		fmt.Fprintf(os.Stderr, "Error: --json flag required\n")
		os.Exit(1)
	}
	var tmpl api.ContentTemplate
	if err := json.Unmarshal([]byte(jsonPayload), &tmpl); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	result, err := services.Templates.UpdateContentTemplate(context.Background(), &tmpl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}

func handleTemplateDeleteContent(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: template ID required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	err := services.Templates.RemoveContentTemplate(context.Background(), args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Content template %s deleted successfully\n", args[0])
}

func handleTemplateListBlueprint(args []string) {
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	spaceKey := parseFlag(args, "--space-key")
	start := parseFlagInt(args, "--start", 0)
	limit := parseFlagInt(args, "--limit", 25)
	expand := parseFlag(args, "--expand")
	var expandList []string
	if expand != "" {
		expandList = split(expand, ",")
	}
	result, err := services.Templates.GetBlueprintTemplates(context.Background(), spaceKey, start, limit, expandList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}

func handleTemplateGetBlueprint(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: blueprint ID required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	expand := parseFlag(args, "--expand")
	var expandList []string
	if expand != "" {
		expandList = split(expand, ",")
	}
	tmpl, err := services.Templates.GetBlueprintTemplate(context.Background(), args[0], expandList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(tmpl)
}

func handleContentStateList(args []string) {
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	start := parseFlagInt(args, "--start", 0)
	limit := parseFlagInt(args, "--limit", 25)
	result, err := services.ContentStates.GetAll(context.Background(), start, limit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}

func handleContentStateGet(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: state ID required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	state, err := services.ContentStates.Get(context.Background(), args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(state)
}

func handleContentStateCreate(args []string) {
	jsonPayload := parseFlag(args, "--json")
	if jsonPayload == "" {
		fmt.Fprintf(os.Stderr, "Error: --json flag required\n")
		os.Exit(1)
	}
	var state api.ContentState
	if err := json.Unmarshal([]byte(jsonPayload), &state); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	result, err := services.ContentStates.Create(context.Background(), &state)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}

func handleContentStateUpdate(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: state ID required\n")
		os.Exit(1)
	}
	jsonPayload := parseFlag(args, "--json")
	if jsonPayload == "" {
		fmt.Fprintf(os.Stderr, "Error: --json flag required\n")
		os.Exit(1)
	}
	var state api.ContentState
	if err := json.Unmarshal([]byte(jsonPayload), &state); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	result, err := services.ContentStates.Update(context.Background(), args[0], &state)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}

func handleContentStateDelete(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: state ID required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	err := services.ContentStates.Delete(context.Background(), args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Content state %s deleted successfully\n", args[0])
}

func handleContentStateContent(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: state ID required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	start := parseFlagInt(args, "--start", 0)
	limit := parseFlagInt(args, "--limit", 25)
	result, err := services.ContentStates.GetContentByState(context.Background(), args[0], start, limit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}

func handleInlineTaskGet(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: task ID required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	task, err := services.InlineTasks.Get(context.Background(), args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(task)
}

func handleInlineTaskList(args []string) {
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	start := parseFlagInt(args, "--start", 0)
	limit := parseFlagInt(args, "--limit", 25)
	result, err := services.InlineTasks.GetAll(context.Background(), start, limit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}

func handleInlineTaskContent(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: content ID required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	start := parseFlagInt(args, "--start", 0)
	limit := parseFlagInt(args, "--limit", 25)
	result, err := services.InlineTasks.GetByContent(context.Background(), args[0], start, limit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}

func handleInlineTaskUpdate(args []string) {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Error: content ID and task ID required\n")
		os.Exit(1)
	}
	status := parseFlag(args, "--status")
	if status == "" {
		fmt.Fprintf(os.Stderr, "Error: --status required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	result, err := services.InlineTasks.Update(context.Background(), args[0], args[1], status)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}

func handleRelationGet(args []string) {
	relationName := parseFlag(args, "--relation")
	sourceType := parseFlag(args, "--source-type")
	sourceKey := parseFlag(args, "--source-key")
	targetType := parseFlag(args, "--target-type")
	targetKey := parseFlag(args, "--target-key")
	if relationName == "" || sourceType == "" || sourceKey == "" || targetType == "" || targetKey == "" {
		fmt.Fprintf(os.Stderr, "Error: --relation, --source-type, --source-key, --target-type, --target-key required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	expand := parseFlag(args, "--expand")
	var expandList []string
	if expand != "" {
		expandList = split(expand, ",")
	}
	relation, err := services.Relations.Get(context.Background(), relationName, sourceType, sourceKey, targetType, targetKey, expandList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(relation)
}

func handleRelationCreate(args []string) {
	relationName := parseFlag(args, "--relation")
	sourceType := parseFlag(args, "--source-type")
	sourceKey := parseFlag(args, "--source-key")
	targetType := parseFlag(args, "--target-type")
	targetKey := parseFlag(args, "--target-key")
	if relationName == "" || sourceType == "" || sourceKey == "" || targetType == "" || targetKey == "" {
		fmt.Fprintf(os.Stderr, "Error: --relation, --source-type, --source-key, --target-type, --target-key required\n")
		os.Exit(1)
	}
	jsonPayload := parseFlag(args, "--json")
	var body interface{}
	if jsonPayload != "" {
		if err := json.Unmarshal([]byte(jsonPayload), &body); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
			os.Exit(1)
		}
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	result, err := services.Relations.Create(context.Background(), relationName, sourceType, sourceKey, targetType, targetKey, body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}

func handleRelationDelete(args []string) {
	relationName := parseFlag(args, "--relation")
	sourceType := parseFlag(args, "--source-type")
	sourceKey := parseFlag(args, "--source-key")
	targetType := parseFlag(args, "--target-type")
	targetKey := parseFlag(args, "--target-key")
	if relationName == "" || sourceType == "" || sourceKey == "" || targetType == "" || targetKey == "" {
		fmt.Fprintf(os.Stderr, "Error: --relation, --source-type, --source-key, --target-type, --target-key required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	err := services.Relations.Delete(context.Background(), relationName, sourceType, sourceKey, targetType, targetKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Relation %s deleted successfully\n", relationName)
}

func handleLongTaskGet(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: task ID required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	task, err := services.LongTasks.Get(context.Background(), args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(task)
}

func handleLongTaskList(args []string) {
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	start := parseFlagInt(args, "--start", 0)
	limit := parseFlagInt(args, "--limit", 25)
	result, err := services.LongTasks.GetAll(context.Background(), start, limit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}

func handleSystemStatus(args []string) {
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	status, err := services.System.GetStatus(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(status)
}

func handleBlueprintList(args []string) {
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	start := parseFlagInt(args, "--start", 0)
	limit := parseFlagInt(args, "--limit", 25)
	expand := parseFlag(args, "--expand")
	var expandList []string
	if expand != "" {
		expandList = split(expand, ",")
	}
	result, err := services.Blueprints.GetAll(context.Background(), start, limit, expandList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}

func handleBlueprintGet(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: blueprint ID required\n")
		os.Exit(1)
	}
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	expand := parseFlag(args, "--expand")
	var expandList []string
	if expand != "" {
		expandList = split(expand, ",")
	}
	bp, err := services.Blueprints.Get(context.Background(), args[0], expandList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(bp)
}

func handleHealthCheckGet(args []string) {
	configPath := parseFlag(args, "--config")
	services := getServices(configPath)
	result, err := services.HealthCheck.Get(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	outputJSON(result)
}
