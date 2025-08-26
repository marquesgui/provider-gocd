package gocd_test

import (
	"encoding/json"
	"fmt"
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
				ID:    ptr.ToPtr("id"),
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
