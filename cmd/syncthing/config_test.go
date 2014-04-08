package main

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestDefaultValues(t *testing.T) {
	expected := OptionsConfiguration{
		ListenAddress:      []string{":22000"},
		GUIEnabled:         true,
		GUIAddress:         "127.0.0.1:8080",
		GlobalAnnServer:    "announce.syncthing.net:22025",
		GlobalAnnEnabled:   true,
		LocalAnnEnabled:    true,
		ParallelRequests:   16,
		MaxSendKbps:        0,
		RescanIntervalS:    60,
		ReconnectIntervalS: 60,
		MaxChangeKbps:      1000,
		StartBrowser:       true,
	}

	cfg, err := readConfigXML(bytes.NewReader(nil))
	if err != io.EOF {
		t.Error(err)
	}

	if !reflect.DeepEqual(cfg.Options, expected) {
		t.Errorf("Default config differs;\n  E: %#v\n  A: %#v", expected, cfg.Options)
	}
}

func TestNodeConfig(t *testing.T) {
	v1data := []byte(`
<configuration version="1">
    <repository id="test" directory="~/Sync">
        <node id="node1" name="node one">
            <address>a</address>
        </node>
        <node id="node2" name="node two">
            <address>b</address>
        </node>
    </repository>
    <options>
        <readOnly>true</readOnly>
    </options>
</configuration>
`)

	v2data := []byte(`
<configuration version="2">
    <repository id="test" directory="~/Sync" ro="true">
        <node id="node1"/>
        <node id="node2"/>
    </repository>
    <node id="node1" name="node one">
        <address>a</address>
    </node>
    <node id="node2" name="node two">
        <address>b</address>
    </node>
</configuration>
`)

	for i, data := range [][]byte{v1data, v2data} {
		cfg, err := readConfigXML(bytes.NewReader(data))
		if err != nil {
			t.Error(err)
		}

		expectedRepos := []RepositoryConfiguration{
			{
				ID:        "test",
				Directory: "~/Sync",
				Nodes:     []NodeConfiguration{{NodeID: "node1"}, {NodeID: "node2"}},
				ReadOnly:  true,
			},
		}
		expectedNodes := []NodeConfiguration{
			{
				NodeID:    "node1",
				Name:      "node one",
				Addresses: []string{"a"},
			},
			{
				NodeID:    "node2",
				Name:      "node two",
				Addresses: []string{"b"},
			},
		}
		expectedNodeIDs := []string{"node1", "node2"}

		if cfg.Version != 2 {
			t.Errorf("%d: Incorrect version %d != 2", i, cfg.Version)
		}
		if !reflect.DeepEqual(cfg.Repositories, expectedRepos) {
			t.Errorf("%d: Incorrect Repositories\n  A: %#v\n  E: %#v", i, cfg.Repositories, expectedRepos)
		}
		if !reflect.DeepEqual(cfg.Nodes, expectedNodes) {
			t.Errorf("%d: Incorrect Nodes\n  A: %#v\n  E: %#v", i, cfg.Nodes, expectedNodes)
		}
		if !reflect.DeepEqual(cfg.Repositories[0].NodeIDs(), expectedNodeIDs) {
			t.Errorf("%d: Incorrect NodeIDs\n  A: %#v\n  E: %#v", i, cfg.Repositories[0].NodeIDs(), expectedNodeIDs)
		}
	}
}

func TestNoListenAddress(t *testing.T) {
	data := []byte(`<configuration version="1">
    <repository directory="~/Sync">
        <node id="..." name="...">
            <address>dynamic</address>
        </node>
    </repository>
    <options>
        <listenAddress></listenAddress>
    </options>
</configuration>
`)

	cfg, err := readConfigXML(bytes.NewReader(data))
	if err != nil {
		t.Error(err)
	}

	expected := []string{""}
	if !reflect.DeepEqual(cfg.Options.ListenAddress, expected) {
		t.Errorf("Unexpected ListenAddress %#v", cfg.Options.ListenAddress)
	}
}

func TestOverriddenValues(t *testing.T) {
	data := []byte(`<configuration version="2">
    <repository directory="~/Sync">
        <node id="..." name="...">
            <address>dynamic</address>
        </node>
    </repository>
    <options>
       <listenAddress>:23000</listenAddress>
        <allowDelete>false</allowDelete>
        <guiEnabled>false</guiEnabled>
        <guiAddress>125.2.2.2:8080</guiAddress>
        <globalAnnounceServer>syncthing.nym.se:22025</globalAnnounceServer>
        <globalAnnounceEnabled>false</globalAnnounceEnabled>
        <localAnnounceEnabled>false</localAnnounceEnabled>
        <parallelRequests>32</parallelRequests>
        <maxSendKbps>1234</maxSendKbps>
        <rescanIntervalS>600</rescanIntervalS>
        <reconnectionIntervalS>6000</reconnectionIntervalS>
        <maxChangeKbps>2345</maxChangeKbps>
        <startBrowser>false</startBrowser>
    </options>
</configuration>
`)

	expected := OptionsConfiguration{
		ListenAddress:      []string{":23000"},
		GUIEnabled:         false,
		GUIAddress:         "125.2.2.2:8080",
		GlobalAnnServer:    "syncthing.nym.se:22025",
		GlobalAnnEnabled:   false,
		LocalAnnEnabled:    false,
		ParallelRequests:   32,
		MaxSendKbps:        1234,
		RescanIntervalS:    600,
		ReconnectIntervalS: 6000,
		MaxChangeKbps:      2345,
		StartBrowser:       false,
	}

	cfg, err := readConfigXML(bytes.NewReader(data))
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(cfg.Options, expected) {
		t.Errorf("Overridden config differs;\n  E: %#v\n  A: %#v", expected, cfg.Options)
	}
}
