package gocd_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/marquesgui/provider-gocd/pkg/gocd"
	"github.com/marquesgui/provider-gocd/pkg/ptr"
)

func Test_JsonUnmarshal(t *testing.T) {
	j := `
	{
  	"_links" : {
  	  "self" : {
  	    "href" : "http://gocd-server.gocd.svc.cluster.local/go/api/admin/pipelines/new_pipeline"
  	  },
  	  "doc" : {
  	    "href" : "https://api.gocd.org/25.3.0/#pipeline-config"
  	  },
  	  "find" : {
  	    "href" : "http://gocd-server.gocd.svc.cluster.local/go/api/admin/pipelines/:pipeline_name"
  	  }
  	},
  	"label_template" : "${COUNT}",
  	"lock_behavior" : "lockOnFailure",
  	"name" : "new_pipeline",
  	"template" : null,
  	"group" : "sample",
  	"origin" : {
  	  "_links" : {
  	    "self" : {
  	      "href" : "http://gocd-server.gocd.svc.cluster.local/go/admin/config_xml"
  	    },
  	    "doc" : {
  	      "href" : "https://api.gocd.org/25.3.0/#get-configuration"
  	    }
  	  },
  	  "type" : "gocd"
  	},
  	"parameters" : [ ],
  	"environment_variables" : [ ],
  	"materials" : [ {
  	  "type" : "git",
  	  "attributes" : {
  	    "url" : "git@github.com:sample_repo/example.git",
  	    "destination" : "dest",
  	    "filter" : null,
  	    "invert_filter" : false,
  	    "name" : null,
  	    "auto_update" : true,
  	    "branch" : "main",
  	    "submodule_folder" : null,
  	    "shallow_clone" : false
  	  }
  	} ],
  	"stages" : [ {
  	  "name" : "defaultStage",
  	  "fetch_materials" : true,
  	  "clean_working_directory" : true,
  	  "never_cleanup_artifacts" : false,
  	  "approval" : {
  	    "type" : "success",
  	    "allow_only_on_success" : false,
  	    "authorization" : {
  	      "roles" : [ ],
  	      "users" : [ ]
  	    }
  	  },
  	  "environment_variables" : [ ],
  	  "jobs" : [ {
  	    "name" : "defaultJob",
  	    "run_instance_count" : null,
  	    "timeout" : "never",
  	    "environment_variables" : [ ],
  	    "resources" : [ ],
  	    "tasks" : [ {
  	      "type" : "exec",
  	      "attributes" : {
  	        "run_if" : [ "passed" ],
  	        "command" : "ls",
  	        "args" : ""
  	      }
  	    } ],
  	    "tabs" : [ ],
  	    "artifacts" : [ ]
  	  } ]
  	} ],
  	"tracking_tool" : null,
  	"timer" : null
	}
`
	var pc gocd.PipelineConfig
	err := json.Unmarshal([]byte(j), &pc)
	if err != nil {
		fmt.Print(fmt.Errorf("Deu ruim: %w", err))
		t.Fail()
	}
	fmt.Printf("object: %+v", pc)
}

func TestEquals(t *testing.T) {
	tests := []struct {
		name    string
		a       any
		b       any
		isEqual bool
	}{
		{
			name:    "stringEqual",
			a:       "teste",
			b:       "teste",
			isEqual: true,
		},
		{
			name:    "stringNotEqual",
			a:       "hi",
			b:       "by",
			isEqual: false,
		},
		{
			name:    "PipelineConfigLockBehaviorEqual",
			a:       gocd.PipelineConfigLockBehaviorUnlockWhenFinished,
			b:       gocd.PipelineConfigLockBehaviorUnlockWhenFinished,
			isEqual: true,
		},
		{
			name:    "PipelineConfigLockBehaviorNotEqual",
			a:       gocd.PipelineConfigLockBehaviorUnlockWhenFinished,
			b:       gocd.PipelineConfigLockBehaviorLockOnFailure,
			isEqual: false,
		},
		{
			name: "PipelineConfigOriginEqual",
			a: &gocd.PipelineConfigOrigin{
				Type:  ptr.ToPtr(gocd.PipelineConfigOriginTypeGoCD),
				ID:    ptr.ToPtr("id"),
				Links: nil,
			},
			b: &gocd.PipelineConfigOrigin{
				Type:  ptr.ToPtr(gocd.PipelineConfigOriginTypeGoCD),
				ID:    ptr.ToPtr("id"),
				Links: nil,
			},
			isEqual: true,
		},
		{
			name: "PipelineConfigOriginNotEqual",
			a: &gocd.PipelineConfigOrigin{
				Type:  ptr.ToPtr(gocd.PipelineConfigOriginTypeGoCD),
				ID:    ptr.ToPtr("id"),
				Links: nil,
			},
			b: &gocd.PipelineConfigOrigin{
				Type:  ptr.ToPtr(gocd.PipelineConfigOriginTypeGoCD),
				ID:    ptr.ToPtr("other-id"),
				Links: nil,
			},
			isEqual: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.isEqual != cmp.Equal(tt.a, tt.b) {
				t.Fail()
			}
		})
	}
}

func TestPipelineConfigsService_Get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", "etag")
		fmt.Fprintln(w, `{"name": "test-pipeline", "group": "sample"}`)
	}))
	defer ts.Close()

	client, _ := gocd.New(gocd.Config{BaseURL: ts.URL})
	pc, etag, err := client.PipelineConfigs().Get(context.Background(), "test-pipeline")

	if err != nil {
		t.Fatalf("PipelineConfigs.Get returned error: %v", err)
	}
	if etag != "etag" {
		t.Errorf("Expected etag 'etag', got %s", etag)
	}
	if *pc.Name != "test-pipeline" {
		t.Errorf("Expected name 'test-pipeline', got %s", *pc.Name)
	}
}

func TestPipelineConfigsService_Create(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", "new-etag")
		fmt.Fprintln(w, `{"name": "new-pipeline", "group": "sample"}`)
	}))
	defer ts.Close()

	client, _ := gocd.New(gocd.Config{BaseURL: ts.URL})
	pc, etag, err := client.PipelineConfigs().Create(context.Background(), &gocd.PipelineConfig{Name: ptr.ToPtr("new-pipeline"), Group: ptr.ToPtr("sample")})

	if err != nil {
		t.Fatalf("PipelineConfigs.Create returned error: %v", err)
	}
	if etag != "new-etag" {
		t.Errorf("Expected etag 'new-etag', got %s", etag)
	}
	if *pc.Name != "new-pipeline" {
		t.Errorf("Expected name 'new-pipeline', got %s", *pc.Name)
	}
}

func TestPipelineConfigsService_Update(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("If-Match") != "old-etag" {
			t.Errorf("Expected If-Match 'old-etag', got %s", r.Header.Get("If-Match"))
		}
		w.Header().Set("ETag", "new-etag")
		fmt.Fprintln(w, `{"name": "test-pipeline", "group": "sample"}`)
	}))
	defer ts.Close()

	client, _ := gocd.New(gocd.Config{BaseURL: ts.URL})
	_, etag, err := client.PipelineConfigs().Update(context.Background(), "old-etag", &gocd.PipelineConfig{Name: ptr.ToPtr("test-pipeline")})

	if err != nil {
		t.Fatalf("PipelineConfigs.Update returned error: %v", err)
	}
	if etag != "new-etag" {
		t.Errorf("Expected etag 'new-etag', got %s", etag)
	}
}

func TestPipelineConfigsService_Delete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, _ := gocd.New(gocd.Config{BaseURL: ts.URL})
	err := client.PipelineConfigs().Delete(context.Background(), "test-pipeline")

	if err != nil {
		t.Fatalf("PipelineConfigs.Delete returned error: %v", err)
	}
}
