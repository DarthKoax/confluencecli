package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darthkoax/confluencecli/internal/client"
	"github.com/darthkoax/confluencecli/internal/config"
)

func testClient(t *testing.T, handler http.Handler) (*client.Client, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)
	cfg := &config.Config{
		Confluence: config.ConfluenceConfig{
			BaseURL:  server.URL,
			APIToken: "test-token",
			Timeout:  30,
		},
		Methods: config.MethodsConfig{
			AllowGet:    true,
			AllowPost:   true,
			AllowPut:    true,
			AllowDelete: true,
		},
		AdminMode: true,
		Endpoints: config.DefaultEndpoints(),
	}
	c, err := client.New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	return c, server
}

func TestContentService_Get(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/content/12345" {
			t.Errorf("Path = %v, want /rest/api/content/12345", r.URL.Path)
		}
		json.NewEncoder(w).Encode(Content{ID: "12345", Title: "Test Page"})
	}))
	defer server.Close()

	svc := NewContentService(c)
	content, err := svc.Get(context.Background(), "12345", nil)
	if err != nil {
		t.Fatal(err)
	}
	if content.ID != "12345" {
		t.Errorf("ID = %v, want 12345", content.ID)
	}
}

func TestContentService_Create(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %v, want POST", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Content{ID: "12346", Title: "New Page"})
	}))
	defer server.Close()

	svc := NewContentService(c)
	content, err := svc.Create(context.Background(), &Content{Title: "New Page"})
	if err != nil {
		t.Fatal(err)
	}
	if content.ID != "12346" {
		t.Errorf("ID = %v, want 12346", content.ID)
	}
}

func TestContentService_EndpointDisabled(t *testing.T) {
	_, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Should not reach server when endpoint is disabled")
	}))
	defer server.Close()

	cfg := &config.Config{
		Confluence: config.ConfluenceConfig{
			BaseURL:  server.URL,
			APIToken: "test-token",
			Timeout:  30,
		},
		Methods: config.MethodsConfig{
			AllowGet:  true,
			AllowPost: true,
		},
		Endpoints: config.DefaultEndpoints(),
	}
	cfg.Endpoints.Content = false

	c2, err := client.New(cfg)
	if err != nil {
		t.Fatal(err)
	}

	svc := NewContentService(c2)
	_, err = svc.Get(context.Background(), "12345", nil)
	if err == nil {
		t.Error("Expected error for disabled content endpoint, got nil")
	}
}

func TestSpaceService_GetAll(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(SpaceResult{
			Results: []Space{
				{Key: "DEV", Name: "Development"},
				{Key: "QA", Name: "Quality Assurance"},
			},
		})
	}))
	defer server.Close()

	svc := NewSpaceService(c)
	result, err := svc.GetAll(context.Background(), "", "", 0, 25, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Results) != 2 {
		t.Errorf("len(Results) = %v, want 2", len(result.Results))
	}
}

func TestSpaceService_Get(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(Space{Key: "DEV", Name: "Development"})
	}))
	defer server.Close()

	svc := NewSpaceService(c)
	space, err := svc.Get(context.Background(), "DEV", nil)
	if err != nil {
		t.Fatal(err)
	}
	if space.Key != "DEV" {
		t.Errorf("Key = %v, want DEV", space.Key)
	}
}

func TestSpaceService_Create_NoAdminRequired(t *testing.T) {
	_, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %v, want POST", r.Method)
		}
		json.NewEncoder(w).Encode(Space{Key: "NEW", Name: "New Space"})
	}))
	defer server.Close()

	cfg := &config.Config{
		Confluence: config.ConfluenceConfig{
			BaseURL:  server.URL,
			APIToken: "test-token",
			Timeout:  30,
		},
		Methods: config.MethodsConfig{
			AllowGet:  true,
			AllowPost: true,
		},
		AdminMode: false,
		Endpoints: config.DefaultEndpoints(),
	}
	c2, _ := client.New(cfg)

	svc := NewSpaceService(c2)
	space, err := svc.Create(context.Background(), &Space{Key: "NEW", Name: "New Space"})
	if err != nil {
		t.Errorf("Expected no error for space create without admin mode, got %v", err)
	}
	if space.Key != "NEW" {
		t.Errorf("Key = %v, want NEW", space.Key)
	}
}

func TestSearchService_Search(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(SearchResult{
			TotalSize: 1,
			Results: []SearchResultItem{
				{Title: "Test Page"},
			},
		})
	}))
	defer server.Close()

	svc := NewSearchService(c)
	result, err := svc.Search(context.Background(), "type=page", 0, 25, nil)
	if err != nil {
		t.Fatal(err)
	}
	if result.TotalSize != 1 {
		t.Errorf("TotalSize = %v, want 1", result.TotalSize)
	}
}

func TestUserService_GetCurrent(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(User{Username: "admin", DisplayName: "Admin User"})
	}))
	defer server.Close()

	svc := NewUserService(c)
	user, err := svc.GetCurrent(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if user.Username != "admin" {
		t.Errorf("Username = %v, want admin", user.Username)
	}
}

func TestGroupService_Get(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(Group{Name: "developers"})
	}))
	defer server.Close()

	svc := NewGroupService(c)
	group, err := svc.Get(context.Background(), "developers", nil)
	if err != nil {
		t.Fatal(err)
	}
	if group.Name != "developers" {
		t.Errorf("Name = %v, want developers", group.Name)
	}
}

func TestGroupService_AddMember_NoAdminRequired(t *testing.T) {
	_, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %v, want POST", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	cfg := &config.Config{
		Confluence: config.ConfluenceConfig{
			BaseURL:  server.URL,
			APIToken: "test-token",
			Timeout:  30,
		},
		Methods:   config.MethodsConfig{AllowGet: true, AllowPost: true},
		AdminMode: false,
		Endpoints: config.DefaultEndpoints(),
	}
	c2, _ := client.New(cfg)

	svc := NewGroupService(c2)
	err := svc.AddMember(context.Background(), "developers", &User{Username: "john"})
	if err != nil {
		t.Errorf("Expected no error for group add member without admin mode, got %v", err)
	}
}

func TestSettingsService_GetSystemInfo_AdminRequired(t *testing.T) {
	_, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Should not reach server when admin mode disabled")
	}))
	defer server.Close()

	cfg := &config.Config{
		Confluence: config.ConfluenceConfig{
			BaseURL:  server.URL,
			APIToken: "test-token",
			Timeout:  30,
		},
		Methods:   config.MethodsConfig{AllowGet: true},
		AdminMode: false,
		Endpoints: config.DefaultEndpoints(),
	}
	c2, _ := client.New(cfg)

	svc := NewSettingsService(c2)
	_, err := svc.GetSystemInfo(context.Background())
	if err == nil {
		t.Error("Expected error for admin operation without admin mode, got nil")
	}
}

func TestAuditService_Get_AdminRequired(t *testing.T) {
	_, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Should not reach server when admin mode disabled")
	}))
	defer server.Close()

	cfg := &config.Config{
		Confluence: config.ConfluenceConfig{
			BaseURL:  server.URL,
			APIToken: "test-token",
			Timeout:  30,
		},
		Methods:   config.MethodsConfig{AllowGet: true},
		AdminMode: false,
		Endpoints: config.DefaultEndpoints(),
	}
	c2, _ := client.New(cfg)

	svc := NewAuditService(c2)
	_, err := svc.Get(context.Background(), "", "", "", 0, 25)
	if err == nil {
		t.Error("Expected error for admin operation without admin mode, got nil")
	}
}

func TestSystemService_GetStatus_AdminRequired(t *testing.T) {
	_, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Should not reach server when admin mode disabled")
	}))
	defer server.Close()

	cfg := &config.Config{
		Confluence: config.ConfluenceConfig{
			BaseURL:  server.URL,
			APIToken: "test-token",
			Timeout:  30,
		},
		Methods:   config.MethodsConfig{AllowGet: true},
		AdminMode: false,
		Endpoints: config.DefaultEndpoints(),
	}
	c2, _ := client.New(cfg)

	svc := NewSystemService(c2)
	_, err := svc.GetStatus(context.Background())
	if err == nil {
		t.Error("Expected error for admin operation without admin mode, got nil")
	}
}

func TestHealthCheckService_Get(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(HealthCheckResponse{
			State: "HEALTHY",
			Checks: []HealthCheckResult{
				{Name: "database", Passed: true},
			},
		})
	}))
	defer server.Close()

	svc := NewHealthCheckService(c)
	result, err := svc.Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if result.State != "HEALTHY" {
		t.Errorf("State = %v, want HEALTHY", result.State)
	}
}

func TestContentStateService_GetAll(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(ContentStateResult{
			Results: []ContentState{
				{ID: "1", Name: "Draft"},
				{ID: "2", Name: "Published"},
			},
		})
	}))
	defer server.Close()

	svc := NewContentStateService(c)
	result, err := svc.GetAll(context.Background(), 0, 25)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Results) != 2 {
		t.Errorf("len(Results) = %v, want 2", len(result.Results))
	}
}

func TestInlineTaskService_Get(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(InlineTask{ID: "task-1", Status: "incomplete"})
	}))
	defer server.Close()

	svc := NewInlineTaskService(c)
	task, err := svc.Get(context.Background(), "task-1")
	if err != nil {
		t.Fatal(err)
	}
	if task.ID != "task-1" {
		t.Errorf("ID = %v, want task-1", task.ID)
	}
}

func TestLongTaskService_Get(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(LongTask{ID: "lt-1", Finished: true})
	}))
	defer server.Close()

	svc := NewLongTaskService(c)
	task, err := svc.Get(context.Background(), "lt-1")
	if err != nil {
		t.Fatal(err)
	}
	if task.ID != "lt-1" {
		t.Errorf("ID = %v, want lt-1", task.ID)
	}
}

func TestTemplateService_GetContentTemplate(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(ContentTemplate{TemplateID: "tmpl-1", Name: "Meeting Notes"})
	}))
	defer server.Close()

	svc := NewTemplateService(c)
	tmpl, err := svc.GetContentTemplate(context.Background(), "tmpl-1", nil)
	if err != nil {
		t.Fatal(err)
	}
	if tmpl.Name != "Meeting Notes" {
		t.Errorf("Name = %v, want Meeting Notes", tmpl.Name)
	}
}

func TestBlueprintService_Get(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(ContentBlueprint{ContentBlueprintID: "bp-1", Name: "Page Blueprint"})
	}))
	defer server.Close()

	svc := NewBlueprintService(c)
	bp, err := svc.Get(context.Background(), "bp-1", nil)
	if err != nil {
		t.Fatal(err)
	}
	if bp.Name != "Page Blueprint" {
		t.Errorf("Name = %v, want Page Blueprint", bp.Name)
	}
}

func TestContentService_GetAll(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(ContentResult{
			Results: []Content{
				{ID: "1", Title: "Page 1"},
				{ID: "2", Title: "Page 2"},
			},
			Size: 2,
		})
	}))
	defer server.Close()

	svc := NewContentService(c)
	result, err := svc.GetAll(context.Background(), "page", "ds", "", 0, 25, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Results) != 2 {
		t.Errorf("len(Results) = %v, want 2", len(result.Results))
	}
}

func TestContentService_Update(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Method = %v, want PUT", r.Method)
		}
		json.NewEncoder(w).Encode(Content{ID: "12345", Title: "Updated Page"})
	}))
	defer server.Close()

	svc := NewContentService(c)
	content, err := svc.Update(context.Background(), "12345", &Content{Title: "Updated Page"})
	if err != nil {
		t.Fatal(err)
	}
	if content.Title != "Updated Page" {
		t.Errorf("Title = %v, want Updated Page", content.Title)
	}
}

func TestContentService_Delete(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Method = %v, want DELETE", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	svc := NewContentService(c)
	err := svc.Delete(context.Background(), "12345")
	if err != nil {
		t.Fatal(err)
	}
}

func TestContentService_GetHistory(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(ContentHistory{Latest: true})
	}))
	defer server.Close()

	svc := NewContentService(c)
	history, err := svc.GetHistory(context.Background(), "12345", nil)
	if err != nil {
		t.Fatal(err)
	}
	if !history.Latest {
		t.Error("Expected Latest to be true")
	}
}

func TestContentService_GetChildren(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(ContentResult{Size: 0})
	}))
	defer server.Close()

	svc := NewContentService(c)
	children, err := svc.GetChildren(context.Background(), "12345", nil)
	if err != nil {
		t.Fatal(err)
	}
	if children == nil {
		t.Error("Expected children to be non-nil")
	}
}

func TestContentService_GetDescendants(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(ContentResult{Size: 0})
	}))
	defer server.Close()

	svc := NewContentService(c)
	descendants, err := svc.GetDescendants(context.Background(), "12345", nil)
	if err != nil {
		t.Fatal(err)
	}
	if descendants == nil {
		t.Error("Expected descendants to be non-nil")
	}
}

func TestContentService_GetVersions(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(ContentResult{Size: 0})
	}))
	defer server.Close()

	svc := NewContentService(c)
	result, err := svc.GetVersions(context.Background(), "12345", 0, 25)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Error("Expected result to be non-nil")
	}
}

func TestContentService_GetComments(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(ContentResult{Size: 0})
	}))
	defer server.Close()

	svc := NewContentService(c)
	result, err := svc.GetComments(context.Background(), "12345")
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Error("Expected result to be non-nil")
	}
}

func TestContentService_GetAttachments(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(ContentResult{Size: 1})
	}))
	defer server.Close()

	svc := NewContentService(c)
	result, err := svc.GetAttachments(context.Background(), "12345")
	if err != nil {
		t.Fatal(err)
	}
	if result.Size != 1 {
		t.Errorf("Size = %v, want 1", result.Size)
	}
}

func TestSpaceService_Create(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %v, want POST", r.Method)
		}
		json.NewEncoder(w).Encode(Space{Key: "NEW", Name: "New Space"})
	}))
	defer server.Close()

	svc := NewSpaceService(c)
	space, err := svc.Create(context.Background(), &Space{Key: "NEW", Name: "New Space"})
	if err != nil {
		t.Fatal(err)
	}
	if space.Key != "NEW" {
		t.Errorf("Key = %v, want NEW", space.Key)
	}
}

func TestSpaceService_Delete(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Method = %v, want DELETE", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	svc := NewSpaceService(c)
	err := svc.Delete(context.Background(), "TEST")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSpaceService_GetContent(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(SpaceContent{
			Page: &ContentResult{Size: 1, Results: []Content{{ID: "1", Title: "Test"}}},
		})
	}))
	defer server.Close()

	svc := NewSpaceService(c)
	result, err := svc.GetContent(context.Background(), "ds", "", 0, 25, nil)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Error("Expected result to be non-nil")
	}
	if result.Page == nil {
		t.Error("Expected Page to be non-nil")
	}
}

func TestSearchService_SearchUser(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(SearchResult{TotalSize: 0})
	}))
	defer server.Close()

	svc := NewSearchService(c)
	result, err := svc.SearchUser(context.Background(), "type=user", 0, 25)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Error("Expected result to be non-nil")
	}
}

func TestUserService_GetAnonymous(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(User{Username: "anonymous", DisplayName: "Anonymous"})
	}))
	defer server.Close()

	svc := NewUserService(c)
	user, err := svc.GetAnonymous(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if user.Username != "anonymous" {
		t.Errorf("Username = %v, want anonymous", user.Username)
	}
}

func TestUserService_Get(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(User{Username: "admin", DisplayName: "Admin"})
	}))
	defer server.Close()

	svc := NewUserService(c)
	user, err := svc.Get(context.Background(), "admin", "", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if user.Username != "admin" {
		t.Errorf("Username = %v, want admin", user.Username)
	}
}

func TestUserService_GetUnknown(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(User{Username: "unknown", DisplayName: "Unknown"})
	}))
	defer server.Close()

	svc := NewUserService(c)
	user, err := svc.GetUnknown(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if user.Username != "unknown" {
		t.Errorf("Username = %v, want unknown", user.Username)
	}
}

func TestUserService_GetUserByEmail(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(User{Username: "admin", Email: "admin@example.com"})
	}))
	defer server.Close()

	svc := NewUserService(c)
	user, err := svc.GetUserByEmail(context.Background(), "admin@example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	if user.Email != "admin@example.com" {
		t.Errorf("Email = %v, want admin@example.com", user.Email)
	}
}

func TestGroupService_GetMembers(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(GroupMembers{
			Results: []User{{Username: "user1"}, {Username: "user2"}},
			Size:    2,
		})
	}))
	defer server.Close()

	svc := NewGroupService(c)
	members, err := svc.GetMembers(context.Background(), "developers", 0, 25)
	if err != nil {
		t.Fatal(err)
	}
	if len(members.Results) != 2 {
		t.Errorf("len(Results) = %v, want 2", len(members.Results))
	}
}

func TestGroupService_AddMember(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %v, want POST", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	svc := NewGroupService(c)
	err := svc.AddMember(context.Background(), "developers", &User{Username: "john"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGroupService_RemoveMember(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Method = %v, want DELETE", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	svc := NewGroupService(c)
	err := svc.RemoveMember(context.Background(), "developers", "john")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSettingsService_GetSystemInfo(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(SystemInfo{CloudID: "test-cloud"})
	}))
	defer server.Close()

	svc := NewSettingsService(c)
	info, err := svc.GetSystemInfo(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if info.CloudID != "test-cloud" {
		t.Errorf("CloudID = %v, want test-cloud", info.CloudID)
	}
}

func TestSettingsService_GetTheme(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(Theme{ThemeKey: "default", Name: "Default Theme"})
	}))
	defer server.Close()

	svc := NewSettingsService(c)
	theme, err := svc.GetTheme(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if theme.ThemeKey != "default" {
		t.Errorf("ThemeKey = %v, want default", theme.ThemeKey)
	}
}

func TestSettingsService_SetTheme(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Method = %v, want PUT", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	svc := NewSettingsService(c)
	err := svc.SetTheme(context.Background(), "modern")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSettingsService_GetLookAndFeel(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(LookAndFeel{})
	}))
	defer server.Close()

	svc := NewSettingsService(c)
	laf, err := svc.GetLookAndFeel(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	if laf == nil {
		t.Error("Expected laf to be non-nil")
	}
}

func TestSettingsService_SetLookAndFeel(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Method = %v, want PUT", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	svc := NewSettingsService(c)
	err := svc.SetLookAndFeel(context.Background(), &LookAndFeel{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestSettingsService_ResetLookAndFeel(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Method = %v, want DELETE", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	svc := NewSettingsService(c)
	err := svc.ResetLookAndFeel(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAuditService_Get(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(AuditResult{Size: 0})
	}))
	defer server.Close()

	svc := NewAuditService(c)
	result, err := svc.Get(context.Background(), "", "", "", 0, 25)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Error("Expected result to be non-nil")
	}
}

func TestAuditService_GetSince(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(AuditResult{Size: 0})
	}))
	defer server.Close()

	svc := NewAuditService(c)
	result, err := svc.GetSince(context.Background(), "7d", "", 0, 25)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Error("Expected result to be non-nil")
	}
}

func TestAuditService_GetRetentionPeriod(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{"number": 30, "units": "days"})
	}))
	defer server.Close()

	svc := NewAuditService(c)
	result, err := svc.GetRetentionPeriod(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Error("Expected result to be non-nil")
	}
}

func TestAuditService_SetRetentionPeriod(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Method = %v, want PUT", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	svc := NewAuditService(c)
	err := svc.SetRetentionPeriod(context.Background(), 30, "days")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAuditService_Export(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{"exportId": "123"})
	}))
	defer server.Close()

	svc := NewAuditService(c)
	result, err := svc.Export(context.Background(), "2024-01-01", "2024-12-31")
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Error("Expected result to be non-nil")
	}
}

func TestTemplateService_GetContentTemplates(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(ContentTemplateResult{Size: 0})
	}))
	defer server.Close()

	svc := NewTemplateService(c)
	result, err := svc.GetContentTemplates(context.Background(), "", 0, 25, nil)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Error("Expected result to be non-nil")
	}
}

func TestTemplateService_CreateContentTemplate(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %v, want POST", r.Method)
		}
		json.NewEncoder(w).Encode(ContentTemplate{TemplateID: "tmpl-1", Name: "New Template"})
	}))
	defer server.Close()

	svc := NewTemplateService(c)
	result, err := svc.CreateContentTemplate(context.Background(), &ContentTemplate{Name: "New Template"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Name != "New Template" {
		t.Errorf("Name = %v, want New Template", result.Name)
	}
}

func TestTemplateService_UpdateContentTemplate(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Method = %v, want PUT", r.Method)
		}
		json.NewEncoder(w).Encode(ContentTemplate{TemplateID: "tmpl-1", Name: "Updated Template"})
	}))
	defer server.Close()

	svc := NewTemplateService(c)
	result, err := svc.UpdateContentTemplate(context.Background(), &ContentTemplate{TemplateID: "tmpl-1", Name: "Updated Template"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Name != "Updated Template" {
		t.Errorf("Name = %v, want Updated Template", result.Name)
	}
}

func TestTemplateService_RemoveContentTemplate(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Method = %v, want DELETE", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	svc := NewTemplateService(c)
	err := svc.RemoveContentTemplate(context.Background(), "tmpl-1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestTemplateService_GetBlueprintTemplates(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(BlueprintTemplateResult{Size: 0})
	}))
	defer server.Close()

	svc := NewTemplateService(c)
	result, err := svc.GetBlueprintTemplates(context.Background(), "", 0, 25, nil)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Error("Expected result to be non-nil")
	}
}

func TestTemplateService_GetBlueprintTemplate(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(BlueprintTemplate{ContentBlueprintID: "bp-1", Name: "Blueprint"})
	}))
	defer server.Close()

	svc := NewTemplateService(c)
	result, err := svc.GetBlueprintTemplate(context.Background(), "bp-1", nil)
	if err != nil {
		t.Fatal(err)
	}
	if result.Name != "Blueprint" {
		t.Errorf("Name = %v, want Blueprint", result.Name)
	}
}

func TestContentStateService_Get(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(ContentState{ID: "1", Name: "Draft"})
	}))
	defer server.Close()

	svc := NewContentStateService(c)
	state, err := svc.Get(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
	if state.Name != "Draft" {
		t.Errorf("Name = %v, want Draft", state.Name)
	}
}

func TestContentStateService_Create(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %v, want POST", r.Method)
		}
		json.NewEncoder(w).Encode(ContentState{ID: "1", Name: "New State"})
	}))
	defer server.Close()

	svc := NewContentStateService(c)
	result, err := svc.Create(context.Background(), &ContentState{Name: "New State"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Name != "New State" {
		t.Errorf("Name = %v, want New State", result.Name)
	}
}

func TestContentStateService_Update(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Method = %v, want PUT", r.Method)
		}
		json.NewEncoder(w).Encode(ContentState{ID: "1", Name: "Updated State"})
	}))
	defer server.Close()

	svc := NewContentStateService(c)
	result, err := svc.Update(context.Background(), "1", &ContentState{Name: "Updated State"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Name != "Updated State" {
		t.Errorf("Name = %v, want Updated State", result.Name)
	}
}

func TestContentStateService_Delete(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Method = %v, want DELETE", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	svc := NewContentStateService(c)
	err := svc.Delete(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestContentStateService_GetContentByState(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(ContentResult{Size: 0})
	}))
	defer server.Close()

	svc := NewContentStateService(c)
	result, err := svc.GetContentByState(context.Background(), "1", 0, 25)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Error("Expected result to be non-nil")
	}
}

func TestInlineTaskService_GetAll(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(InlineTaskResult{Size: 0})
	}))
	defer server.Close()

	svc := NewInlineTaskService(c)
	result, err := svc.GetAll(context.Background(), 0, 25)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Error("Expected result to be non-nil")
	}
}

func TestInlineTaskService_GetByContent(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(InlineTaskResult{Size: 0})
	}))
	defer server.Close()

	svc := NewInlineTaskService(c)
	result, err := svc.GetByContent(context.Background(), "12345", 0, 25)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Error("Expected result to be non-nil")
	}
}

func TestInlineTaskService_Update(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Method = %v, want PUT", r.Method)
		}
		json.NewEncoder(w).Encode(InlineTask{ID: "task-1", Status: "complete"})
	}))
	defer server.Close()

	svc := NewInlineTaskService(c)
	result, err := svc.Update(context.Background(), "12345", "task-1", "complete")
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != "complete" {
		t.Errorf("Status = %v, want complete", result.Status)
	}
}

func TestRelationService_Get(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(Relation{Name: "favorite"})
	}))
	defer server.Close()

	svc := NewRelationService(c)
	relation, err := svc.Get(context.Background(), "favorite", "user", "admin", "content", "12345", nil)
	if err != nil {
		t.Fatal(err)
	}
	if relation.Name != "favorite" {
		t.Errorf("Name = %v, want favorite", relation.Name)
	}
}

func TestRelationService_Create(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %v, want POST", r.Method)
		}
		json.NewEncoder(w).Encode(Relation{Name: "favorite"})
	}))
	defer server.Close()

	svc := NewRelationService(c)
	result, err := svc.Create(context.Background(), "favorite", "user", "admin", "content", "12345", nil)
	if err != nil {
		t.Fatal(err)
	}
	if result.Name != "favorite" {
		t.Errorf("Name = %v, want favorite", result.Name)
	}
}

func TestRelationService_Delete(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Method = %v, want DELETE", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	svc := NewRelationService(c)
	err := svc.Delete(context.Background(), "favorite", "user", "admin", "content", "12345")
	if err != nil {
		t.Fatal(err)
	}
}

func TestLongTaskService_GetAll(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(LongTaskResult{Size: 0})
	}))
	defer server.Close()

	svc := NewLongTaskService(c)
	result, err := svc.GetAll(context.Background(), 0, 25)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Error("Expected result to be non-nil")
	}
}

func TestSystemService_GetStatus(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(SystemStatus{State: "RUNNING"})
	}))
	defer server.Close()

	svc := NewSystemService(c)
	status, err := svc.GetStatus(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if status.State != "RUNNING" {
		t.Errorf("State = %v, want RUNNING", status.State)
	}
}

func TestBlueprintService_GetAll(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(ContentBlueprintResult{Size: 0})
	}))
	defer server.Close()

	svc := NewBlueprintService(c)
	result, err := svc.GetAll(context.Background(), 0, 25, nil)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Error("Expected result to be non-nil")
	}
}

func TestNewServices(t *testing.T) {
	c, server := testClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()

	services := NewServices(c)
	if services == nil {
		t.Error("Expected services to be non-nil")
	}
	if services.Content == nil {
		t.Error("Expected Content service to be non-nil")
	}
	if services.Spaces == nil {
		t.Error("Expected Spaces service to be non-nil")
	}
	if services.Search == nil {
		t.Error("Expected Search service to be non-nil")
	}
	if services.Users == nil {
		t.Error("Expected Users service to be non-nil")
	}
	if services.Groups == nil {
		t.Error("Expected Groups service to be non-nil")
	}
	if services.Settings == nil {
		t.Error("Expected Settings service to be non-nil")
	}
	if services.Audit == nil {
		t.Error("Expected Audit service to be non-nil")
	}
	if services.Templates == nil {
		t.Error("Expected Templates service to be non-nil")
	}
	if services.ContentStates == nil {
		t.Error("Expected ContentStates service to be non-nil")
	}
	if services.InlineTasks == nil {
		t.Error("Expected InlineTasks service to be non-nil")
	}
	if services.Relations == nil {
		t.Error("Expected Relations service to be non-nil")
	}
	if services.LongTasks == nil {
		t.Error("Expected LongTasks service to be non-nil")
	}
	if services.System == nil {
		t.Error("Expected System service to be non-nil")
	}
	if services.Blueprints == nil {
		t.Error("Expected Blueprints service to be non-nil")
	}
	if services.HealthCheck == nil {
		t.Error("Expected HealthCheck service to be non-nil")
	}
}

func TestAllServices_EndpointDisabled(t *testing.T) {
	endpointTests := []struct {
		name     string
		endpoint string
		testFn   func(*client.Client) error
	}{
		{"Content.Get", "content", func(c *client.Client) error { _, err := NewContentService(c).Get(context.Background(), "1", nil); return err }},
		{"Content.GetAll", "content", func(c *client.Client) error { _, err := NewContentService(c).GetAll(context.Background(), "", "", "", 0, 25, nil); return err }},
		{"Spaces.GetAll", "spaces", func(c *client.Client) error { _, err := NewSpaceService(c).GetAll(context.Background(), "", "", 0, 25, nil); return err }},
		{"Spaces.Get", "spaces", func(c *client.Client) error { _, err := NewSpaceService(c).Get(context.Background(), "ds", nil); return err }},
		{"Search.Search", "search", func(c *client.Client) error { _, err := NewSearchService(c).Search(context.Background(), "type=page", 0, 25, nil); return err }},
		{"Search.SearchUser", "search", func(c *client.Client) error { _, err := NewSearchService(c).SearchUser(context.Background(), "type=user", 0, 25); return err }},
		{"Users.GetCurrent", "users", func(c *client.Client) error { _, err := NewUserService(c).GetCurrent(context.Background(), nil); return err }},
		{"Users.GetAnonymous", "users", func(c *client.Client) error { _, err := NewUserService(c).GetAnonymous(context.Background(), nil); return err }},
		{"Users.Get", "users", func(c *client.Client) error { _, err := NewUserService(c).Get(context.Background(), "admin", "", "", nil); return err }},
		{"Users.GetUnknown", "users", func(c *client.Client) error { _, err := NewUserService(c).GetUnknown(context.Background(), nil); return err }},
		{"Users.GetUserByEmail", "users", func(c *client.Client) error { _, err := NewUserService(c).GetUserByEmail(context.Background(), "test@test.com", nil); return err }},
		{"Groups.Get", "groups", func(c *client.Client) error { _, err := NewGroupService(c).Get(context.Background(), "test", nil); return err }},
		{"Groups.GetMembers", "groups", func(c *client.Client) error { _, err := NewGroupService(c).GetMembers(context.Background(), "test", 0, 25); return err }},
		{"Settings.GetSystemInfo", "settings", func(c *client.Client) error { _, err := NewSettingsService(c).GetSystemInfo(context.Background()); return err }},
		{"Audit.Get", "audit", func(c *client.Client) error { _, err := NewAuditService(c).Get(context.Background(), "", "", "", 0, 25); return err }},
		{"Templates.GetContentTemplates", "templates", func(c *client.Client) error { _, err := NewTemplateService(c).GetContentTemplates(context.Background(), "", 0, 25, nil); return err }},
		{"ContentStates.GetAll", "contentstates", func(c *client.Client) error { _, err := NewContentStateService(c).GetAll(context.Background(), 0, 25); return err }},
		{"InlineTasks.Get", "inlinetasks", func(c *client.Client) error { _, err := NewInlineTaskService(c).Get(context.Background(), "1"); return err }},
		{"InlineTasks.GetAll", "inlinetasks", func(c *client.Client) error { _, err := NewInlineTaskService(c).GetAll(context.Background(), 0, 25); return err }},
		{"Relations.Get", "relations", func(c *client.Client) error { _, err := NewRelationService(c).Get(context.Background(), "fav", "user", "a", "content", "1", nil); return err }},
		{"LongTasks.Get", "longtasks", func(c *client.Client) error { _, err := NewLongTaskService(c).Get(context.Background(), "1"); return err }},
		{"LongTasks.GetAll", "longtasks", func(c *client.Client) error { _, err := NewLongTaskService(c).GetAll(context.Background(), 0, 25); return err }},
		{"System.GetStatus", "system", func(c *client.Client) error { _, err := NewSystemService(c).GetStatus(context.Background()); return err }},
		{"Blueprints.Get", "blueprints", func(c *client.Client) error { _, err := NewBlueprintService(c).Get(context.Background(), "1", nil); return err }},
		{"Blueprints.GetAll", "blueprints", func(c *client.Client) error { _, err := NewBlueprintService(c).GetAll(context.Background(), 0, 25, nil); return err }},
		{"HealthCheck.Get", "healthcheck", func(c *client.Client) error { _, err := NewHealthCheckService(c).Get(context.Background()); return err }},
	}

	for _, tt := range endpointTests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Error("Should not reach server when endpoint is disabled")
			}))
			defer server.Close()

			cfg := &config.Config{
				Confluence: config.ConfluenceConfig{
					BaseURL:  server.URL,
					APIToken: "test-token",
					Timeout:  30,
				},
				Methods: config.MethodsConfig{
					AllowGet:  true,
					AllowPost: true,
				},
				AdminMode: true,
				Endpoints: config.DefaultEndpoints(),
			}

			switch tt.endpoint {
			case "content":
				cfg.Endpoints.Content = false
			case "spaces":
				cfg.Endpoints.Spaces = false
			case "search":
				cfg.Endpoints.Search = false
			case "users":
				cfg.Endpoints.Users = false
			case "groups":
				cfg.Endpoints.Groups = false
			case "settings":
				cfg.Endpoints.Settings = false
			case "audit":
				cfg.Endpoints.Audit = false
			case "templates":
				cfg.Endpoints.Templates = false
			case "contentstates":
				cfg.Endpoints.ContentStates = false
			case "inlinetasks":
				cfg.Endpoints.InlineTasks = false
			case "relations":
				cfg.Endpoints.Relations = false
			case "longtasks":
				cfg.Endpoints.LongTasks = false
			case "system":
				cfg.Endpoints.System = false
			case "blueprints":
				cfg.Endpoints.Blueprints = false
			case "healthcheck":
				cfg.Endpoints.HealthCheck = false
			}

			c, err := client.New(cfg)
			if err != nil {
				t.Fatal(err)
			}

			err = tt.testFn(c)
			if err == nil {
				t.Errorf("%s expected error for disabled endpoint %s, got nil", tt.name, tt.endpoint)
			}
		})
	}
}
