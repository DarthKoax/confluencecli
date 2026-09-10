package main

import (
	"fmt"
	"os"
)

func handleCommandWithHelp(commandName string, args []string, printHelp func(), actionHandlers map[string]func([]string), printActionHelp func(string)) {
	if len(args) == 0 {
		printHelp()
		return
	}

	if args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		if len(args) > 1 {
			printActionHelp(args[1])
		} else {
			printHelp()
		}
		return
	}

	action := args[0]
	actionArgs := args[1:]

	if len(actionArgs) > 0 && (actionArgs[0] == "help" || actionArgs[0] == "--help" || actionArgs[0] == "-h") {
		printActionHelp(action)
		return
	}

	if handler, ok := actionHandlers[action]; ok {
		handler(actionArgs)
	} else {
		fmt.Fprintf(os.Stderr, "Unknown %s action: %s\n\n", commandName, action)
		printHelp()
		os.Exit(1)
	}
}

func handleContent(args []string) {
	actionHandlers := map[string]func([]string){
		"get":          handleContentGet,
		"list":         handleContentList,
		"create":       handleContentCreate,
		"update":       handleContentUpdate,
		"delete":       handleContentDelete,
		"history":      handleContentHistory,
		"children":     handleContentChildren,
		"descendants":  handleContentDescendants,
		"versions":     handleContentVersions,
		"comments":     handleContentComments,
		"attachments":  handleContentAttachments,
	}
	handleCommandWithHelp("content", args, printContentHelp, actionHandlers, printContentActionHelp)
}

func handleSpace(args []string) {
	actionHandlers := map[string]func([]string){
		"list":    handleSpaceList,
		"get":     handleSpaceGet,
		"create":  handleSpaceCreate,
		"delete":  handleSpaceDelete,
		"content": handleSpaceContent,
	}
	handleCommandWithHelp("space", args, printSpaceHelp, actionHandlers, printSpaceActionHelp)
}

func handleSearch(args []string) {
	actionHandlers := map[string]func([]string){
		"cql":  handleSearchCQL,
		"user": handleSearchUser,
	}
	handleCommandWithHelp("search", args, printSearchHelp, actionHandlers, printSearchActionHelp)
}

func handleUser(args []string) {
	actionHandlers := map[string]func([]string){
		"current":   handleUserCurrent,
		"anonymous": handleUserAnonymous,
		"get":       handleUserGet,
		"unknown":   handleUserUnknown,
		"email":     handleUserEmail,
	}
	handleCommandWithHelp("user", args, printUserHelp, actionHandlers, printUserActionHelp)
}

func handleGroup(args []string) {
	actionHandlers := map[string]func([]string){
		"get":          handleGroupGet,
		"members":      handleGroupMembers,
		"add-member":   handleGroupAddMember,
		"remove-member": handleGroupRemoveMember,
	}
	handleCommandWithHelp("group", args, printGroupHelp, actionHandlers, printGroupActionHelp)
}

func handleSettings(args []string) {
	actionHandlers := map[string]func([]string){
		"systeminfo":   handleSettingsSystemInfo,
		"theme-get":    handleSettingsThemeGet,
		"theme-set":    handleSettingsThemeSet,
		"laf-get":      handleSettingsLookAndFeelGet,
		"laf-set":      handleSettingsLookAndFeelSet,
		"laf-reset":    handleSettingsLookAndFeelReset,
	}
	handleCommandWithHelp("settings", args, printSettingsHelp, actionHandlers, printSettingsActionHelp)
}

func handleAudit(args []string) {
	actionHandlers := map[string]func([]string){
		"get":       handleAuditGet,
		"since":     handleAuditSince,
		"retention": handleAuditRetention,
		"export":    handleAuditExport,
	}
	handleCommandWithHelp("audit", args, printAuditHelp, actionHandlers, printAuditActionHelp)
}

func handleTemplate(args []string) {
	actionHandlers := map[string]func([]string){
		"list-content":    handleTemplateListContent,
		"get-content":     handleTemplateGetContent,
		"create-content":  handleTemplateCreateContent,
		"update-content":  handleTemplateUpdateContent,
		"delete-content":  handleTemplateDeleteContent,
		"list-blueprint":  handleTemplateListBlueprint,
		"get-blueprint":   handleTemplateGetBlueprint,
	}
	handleCommandWithHelp("template", args, printTemplateHelp, actionHandlers, printTemplateActionHelp)
}

func handleContentState(args []string) {
	actionHandlers := map[string]func([]string){
		"list":    handleContentStateList,
		"get":     handleContentStateGet,
		"create":  handleContentStateCreate,
		"update":  handleContentStateUpdate,
		"delete":  handleContentStateDelete,
		"content": handleContentStateContent,
	}
	handleCommandWithHelp("contentstate", args, printContentStateHelp, actionHandlers, printContentStateActionHelp)
}

func handleInlineTask(args []string) {
	actionHandlers := map[string]func([]string){
		"get":     handleInlineTaskGet,
		"list":    handleInlineTaskList,
		"content": handleInlineTaskContent,
		"update":  handleInlineTaskUpdate,
	}
	handleCommandWithHelp("inlinetask", args, printInlineTaskHelp, actionHandlers, printInlineTaskActionHelp)
}

func handleRelation(args []string) {
	actionHandlers := map[string]func([]string){
		"get":    handleRelationGet,
		"create": handleRelationCreate,
		"delete": handleRelationDelete,
	}
	handleCommandWithHelp("relation", args, printRelationHelp, actionHandlers, printRelationActionHelp)
}

func handleLongTask(args []string) {
	actionHandlers := map[string]func([]string){
		"get":  handleLongTaskGet,
		"list": handleLongTaskList,
	}
	handleCommandWithHelp("longtask", args, printLongTaskHelp, actionHandlers, printLongTaskActionHelp)
}

func handleSystem(args []string) {
	actionHandlers := map[string]func([]string){
		"status": handleSystemStatus,
	}
	handleCommandWithHelp("system", args, printSystemHelp, actionHandlers, printSystemActionHelp)
}

func handleBlueprint(args []string) {
	actionHandlers := map[string]func([]string){
		"list": handleBlueprintList,
		"get":  handleBlueprintGet,
	}
	handleCommandWithHelp("blueprint", args, printBlueprintHelp, actionHandlers, printBlueprintActionHelp)
}

func handleHealthCheck(args []string) {
	actionHandlers := map[string]func([]string){
		"get": handleHealthCheckGet,
	}
	handleCommandWithHelp("healthcheck", args, printHealthCheckHelp, actionHandlers, printHealthCheckActionHelp)
}

func printContentActionHelp(action string) {
	switch action {
	default:
		fmt.Fprintf(os.Stderr, "Unknown content action: %s\n\n", action)
		printContentHelp()
		os.Exit(1)
	}
}

func printSpaceActionHelp(action string) {
	switch action {
	default:
		fmt.Fprintf(os.Stderr, "Unknown space action: %s\n\n", action)
		printSpaceHelp()
		os.Exit(1)
	}
}

func printSearchActionHelp(action string) {
	switch action {
	default:
		fmt.Fprintf(os.Stderr, "Unknown search action: %s\n\n", action)
		printSearchHelp()
		os.Exit(1)
	}
}

func printUserActionHelp(action string) {
	switch action {
	default:
		fmt.Fprintf(os.Stderr, "Unknown user action: %s\n\n", action)
		printUserHelp()
		os.Exit(1)
	}
}

func printGroupActionHelp(action string) {
	switch action {
	default:
		fmt.Fprintf(os.Stderr, "Unknown group action: %s\n\n", action)
		printGroupHelp()
		os.Exit(1)
	}
}

func printSettingsActionHelp(action string) {
	switch action {
	default:
		fmt.Fprintf(os.Stderr, "Unknown settings action: %s\n\n", action)
		printSettingsHelp()
		os.Exit(1)
	}
}

func printAuditActionHelp(action string) {
	switch action {
	default:
		fmt.Fprintf(os.Stderr, "Unknown audit action: %s\n\n", action)
		printAuditHelp()
		os.Exit(1)
	}
}

func printTemplateActionHelp(action string) {
	switch action {
	default:
		fmt.Fprintf(os.Stderr, "Unknown template action: %s\n\n", action)
		printTemplateHelp()
		os.Exit(1)
	}
}

func printContentStateActionHelp(action string) {
	switch action {
	default:
		fmt.Fprintf(os.Stderr, "Unknown contentstate action: %s\n\n", action)
		printContentStateHelp()
		os.Exit(1)
	}
}

func printInlineTaskActionHelp(action string) {
	switch action {
	default:
		fmt.Fprintf(os.Stderr, "Unknown inlinetask action: %s\n\n", action)
		printInlineTaskHelp()
		os.Exit(1)
	}
}

func printRelationActionHelp(action string) {
	switch action {
	default:
		fmt.Fprintf(os.Stderr, "Unknown relation action: %s\n\n", action)
		printRelationHelp()
		os.Exit(1)
	}
}

func printLongTaskActionHelp(action string) {
	switch action {
	default:
		fmt.Fprintf(os.Stderr, "Unknown longtask action: %s\n\n", action)
		printLongTaskHelp()
		os.Exit(1)
	}
}

func printSystemActionHelp(action string) {
	switch action {
	default:
		fmt.Fprintf(os.Stderr, "Unknown system action: %s\n\n", action)
		printSystemHelp()
		os.Exit(1)
	}
}

func printBlueprintActionHelp(action string) {
	switch action {
	default:
		fmt.Fprintf(os.Stderr, "Unknown blueprint action: %s\n\n", action)
		printBlueprintHelp()
		os.Exit(1)
	}
}

func printHealthCheckActionHelp(action string) {
	switch action {
	default:
		fmt.Fprintf(os.Stderr, "Unknown healthcheck action: %s\n\n", action)
		printHealthCheckHelp()
		os.Exit(1)
	}
}
