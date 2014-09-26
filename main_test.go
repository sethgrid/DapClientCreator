package main

import (
	"testing"

	_ "log"
)

func TestParseMetaData(t *testing.T) {
	data := `
	[{
        "description": "Custom Endpoint GET Example",
        "location": "api/v1/custom/endpoint",
        "method": "GET",
        "notes": "Used as an example endpoint. Only returns a ping response.",
        "primary": "id",
        "properties": {
            "sample_property": {
                "description": "example property",
                "type": "varchar"
            }
        },
        "required": [],
        "responseSchema": {
            "sample_property": "varchar"
        },
        "title": "",
        "type": "mysql"
    }]`

	m, err := ParseMeta([]byte(data))
	if err != nil {
		t.Fatalf("error parsing meta file for single entry: %s", err)
	}
	// currently assuming that if one field was populated, the all the fields will be populated
	if m[0].Description != "Custom Endpoint GET Example" {
		t.Error("error getting description. got %v, want %v", m[0].Description, "Custom Endpoint GET Example")
	}
}

func TestParseMetaMulti(t *testing.T) {
	meta := `
	[{
        "description": "Custom Endpoint GET Example",
        "location": "api/v1/custom/endpoint",
        "method": "GET",
        "notes": "Used as an example endpoint. Only returns a ping response.",
        "primary": "id",
        "properties": {
            "sample_property": {
                "description": "example property",
                "type": "varchar"
            }
        },
        "required": [],
        "responseSchema": {
            "sample_property": "varchar"
        },
        "title": "",
        "type": "mysql"
    },
    {
        "description": "Custom Endpoint POST Example",
        "location": "api/v1/custom/endpoint",
        "method": "POST",
        "notes": "Used as an example endpoint. Will accept data.",
        "primary": "id",
        "properties": {
            "sample_property": {
                "description": "example property",
                "type": "varchar"
            }
        },
        "required": [],
        "responseSchema": {
            "InsertedID": "int",
            "Message": "string"
        },
        "title": "",
        "type": "mysql"
    }
]`

	metas, err := ParseMeta([]byte(meta))
	if err != nil {
		t.Error("ParseMeta should not error ", err)
	}
	if len(metas) != 2 {
		t.Errorf("Wrong number of metas found. Got %d, Want %d", len(metas), 2)
	}
}
