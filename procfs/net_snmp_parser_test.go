package procfs

import (
	"bytes"
	"fmt"
	"path"
	"testing"
)

type NetSnmpTestCase struct {
	name            string
	procfsRoot      string
	primeProcfsRoot string
	wantNetSnmp     *NetSnmp
	wantError       error
	wantInfoChanged bool
}

var netSnmpTestDataDir = path.Join(PROCFS_TESTDATA_ROOT, "net", "snmp")

var netSnmpIndexName = []string{
	"NET_SNMP_IP_FORWARDING",
	"NET_SNMP_IP_DEFAULT_TTL",
	"NET_SNMP_IP_IN_RECEIVES",
	"NET_SNMP_IP_IN_HDR_ERRORS",
	"NET_SNMP_IP_IN_ADDR_ERRORS",
	"NET_SNMP_IP_FORW_DATAGRAMS",
	"NET_SNMP_IP_IN_UNKNOWN_PROTOS",
	"NET_SNMP_IP_IN_DISCARDS",
	"NET_SNMP_IP_IN_DELIVERS",
	"NET_SNMP_IP_OUT_REQUESTS",
	"NET_SNMP_IP_OUT_DISCARDS",
	"NET_SNMP_IP_OUT_NO_ROUTES",
	"NET_SNMP_IP_REASM_TIMEOUT",
	"NET_SNMP_IP_REASM_REQDS",
	"NET_SNMP_IP_REASM_OKS",
	"NET_SNMP_IP_REASM_FAILS",
	"NET_SNMP_IP_FRAG_OKS",
	"NET_SNMP_IP_FRAG_FAILS",
	"NET_SNMP_IP_FRAG_CREATES",
	"NET_SNMP_ICMP_IN_MSGS",
	"NET_SNMP_ICMP_IN_ERRORS",
	"NET_SNMP_ICMP_IN_CSUM_ERRORS",
	"NET_SNMP_ICMP_IN_DEST_UNREACHS",
	"NET_SNMP_ICMP_IN_TIME_EXCDS",
	"NET_SNMP_ICMP_IN_PARM_PROBS",
	"NET_SNMP_ICMP_IN_SRC_QUENCHS",
	"NET_SNMP_ICMP_IN_REDIRECTS",
	"NET_SNMP_ICMP_IN_ECHOS",
	"NET_SNMP_ICMP_IN_ECHO_REPS",
	"NET_SNMP_ICMP_IN_TIMESTAMPS",
	"NET_SNMP_ICMP_IN_TIMESTAMP_REPS",
	"NET_SNMP_ICMP_IN_ADDR_MASKS",
	"NET_SNMP_ICMP_IN_ADDR_MASK_REPS",
	"NET_SNMP_ICMP_OUT_MSGS",
	"NET_SNMP_ICMP_OUT_ERRORS",
	"NET_SNMP_ICMP_OUT_DEST_UNREACHS",
	"NET_SNMP_ICMP_OUT_TIME_EXCDS",
	"NET_SNMP_ICMP_OUT_PARM_PROBS",
	"NET_SNMP_ICMP_OUT_SRC_QUENCHS",
	"NET_SNMP_ICMP_OUT_REDIRECTS",
	"NET_SNMP_ICMP_OUT_ECHOS",
	"NET_SNMP_ICMP_OUT_ECHO_REPS",
	"NET_SNMP_ICMP_OUT_TIMESTAMPS",
	"NET_SNMP_ICMP_OUT_TIMESTAMP_REPS",
	"NET_SNMP_ICMP_OUT_ADDR_MASKS",
	"NET_SNMP_ICMP_OUT_ADDR_MASK_REPS",
	"NET_SNMP_ICMPMSG_IN_TYPE3",
	"NET_SNMP_ICMPMSG_OUT_TYPE3",
	"NET_SNMP_TCP_RTO_ALGORITHM",
	"NET_SNMP_TCP_RTO_MIN",
	"NET_SNMP_TCP_RTO_MAX",
	"NET_SNMP_TCP_MAX_CONN",
	"NET_SNMP_TCP_ACTIVE_OPENS",
	"NET_SNMP_TCP_PASSIVE_OPENS",
	"NET_SNMP_TCP_ATTEMPT_FAILS",
	"NET_SNMP_TCP_ESTAB_RESETS",
	"NET_SNMP_TCP_CURR_ESTAB",
	"NET_SNMP_TCP_IN_SEGS",
	"NET_SNMP_TCP_OUT_SEGS",
	"NET_SNMP_TCP_RETRANS_SEGS",
	"NET_SNMP_TCP_IN_ERRS",
	"NET_SNMP_TCP_OUT_RSTS",
	"NET_SNMP_TCP_IN_CSUM_ERRORS",
	"NET_SNMP_UDP_IN_DATAGRAMS",
	"NET_SNMP_UDP_NO_PORTS",
	"NET_SNMP_UDP_IN_ERRORS",
	"NET_SNMP_UDP_OUT_DATAGRAMS",
	"NET_SNMP_UDP_RCVBUF_ERRORS",
	"NET_SNMP_UDP_SNDBUF_ERRORS",
	"NET_SNMP_UDP_IN_CSUM_ERRORS",
	"NET_SNMP_UDP_IGNORED_MULTI",
	"NET_SNMP_UDP_MEM_ERRORS",
	"NET_SNMP_UDPLITE_IN_DATAGRAMS",
	"NET_SNMP_UDPLITE_NO_PORTS",
	"NET_SNMP_UDPLITE_IN_ERRORS",
	"NET_SNMP_UDPLITE_OUT_DATAGRAMS",
	"NET_SNMP_UDPLITE_RCVBUF_ERRORS",
	"NET_SNMP_UDPLITE_SNDBUF_ERRORS",
	"NET_SNMP_UDPLITE_IN_CSUM_ERRORS",
	"NET_SNMP_UDPLITE_IGNORED_MULTI",
	"NET_SNMP_UDPLITE_MEM_ERRORS",
}

func testNetSnmpAsUInt32(i int32) uint32 {
	return uint32(i)
}

func testNetSnmpParser(tc *NetSnmpTestCase, t *testing.T) {
	t.Logf(`
name=%q
procfsRoot=%q
primeProcfsRoot=%q
`,
		tc.name, tc.procfsRoot, tc.primeProcfsRoot,
	)

	var netSnmp *NetSnmp
	if tc.primeProcfsRoot != "" {
		primeNetSnmp := NewNetSnmp(tc.primeProcfsRoot)
		err := primeNetSnmp.Parse()
		if err != nil {
			t.Fatal(err)
		}
		netSnmp = primeNetSnmp.Clone(true)
		if tc.procfsRoot != "" {
			netSnmp.path = NetSnmpPath(tc.procfsRoot)
		}
	} else {
		netSnmp = NewNetSnmp(tc.procfsRoot)
	}

	err := netSnmp.Parse()
	if tc.wantError != nil {
		if err == nil || tc.wantError.Error() != err.Error() {
			t.Fatalf("want: %v error, got: %v", tc.wantError, err)
		}
		return
	}
	if err != nil {
		t.Fatal(err)
	}

	wantNetSnmp := tc.wantNetSnmp
	if len(wantNetSnmp.Values) != len(netSnmp.Values) {
		t.Fatalf(
			"\nlen(Values): want: %d, got: %d",
			len(wantNetSnmp.Values), len(netSnmp.Values),
		)
	}

	errBuf := &bytes.Buffer{}
	for i, wantValue := range wantNetSnmp.Values {
		gotValue := netSnmp.Values[i]
		if wantValue != gotValue {
			fmt.Fprintf(
				errBuf,
				"\nValues[%d (%s)]: want: %d, got: %d",
				i, netSnmpIndexName[i], wantValue, gotValue,
			)
		}
	}
	if tc.wantInfoChanged && netSnmp.InfoChanged == nil {
		fmt.Fprintf(
			errBuf,
			"\nInfoChanged expected: got: %v",
			netSnmp.InfoChanged,
		)
	} else if !tc.wantInfoChanged && netSnmp.InfoChanged != nil {
		fmt.Fprintf(
			errBuf,
			"\nunexpected InfoChanged: %q",
			netSnmp.InfoChanged,
		)
	}

	if errBuf.Len() > 0 {
		t.Fatal(errBuf.String())
	}
}

func TestNetSnmpParser(t *testing.T) {
	netSnmpTestReference := NewNetSnmp(path.Join(netSnmpTestDataDir, "reference"))
	err := netSnmpTestReference.Parse()
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []*NetSnmpTestCase{
		{
			name:       "field_mapping",
			procfsRoot: path.Join(netSnmpTestDataDir, "field_mapping"),
			wantNetSnmp: &NetSnmp{
				Values: []uint32{
					1000, 1001, 1002, 1003, 1004, 1005, 1006, 1007, 1008, 1009, 1010, 1011, 1012, 1013, 1014, 1015, 1016, 1017, 1018,
					3000, 3001, 3002, 3003, 3004, 3005, 3006, 3007, 3008, 3009, 3010, 3011, 3012, 3013, 3014, 3015, 3016, 3017, 3018, 3019, 3020, 3021, 3022, 3023, 3024, 3025, 3026,
					5000, 5001,
					7000, 7001, 7002, testNetSnmpAsUInt32(-7003), 7004, 7005, 7006, 7007, 7008, 7009, 7010, 7011, 7012, 7013, 7014,
					9000, 9001, 9002, 9003, 9004, 9005, 9006, 9007, 9008,
					11000, 11001, 11002, 11003, 11004, 11005, 11006, 11007, 11008,
				},
			},
		},
		{
			name:            "reuse",
			procfsRoot:      path.Join(netSnmpTestDataDir, "field_mapping"),
			primeProcfsRoot: path.Join(netSnmpTestDataDir, "reference"),
			wantNetSnmp: &NetSnmp{
				Values: []uint32{
					1000, 1001, 1002, 1003, 1004, 1005, 1006, 1007, 1008, 1009, 1010, 1011, 1012, 1013, 1014, 1015, 1016, 1017, 1018,
					3000, 3001, 3002, 3003, 3004, 3005, 3006, 3007, 3008, 3009, 3010, 3011, 3012, 3013, 3014, 3015, 3016, 3017, 3018, 3019, 3020, 3021, 3022, 3023, 3024, 3025, 3026,
					5000, 5001,
					7000, 7001, 7002, testNetSnmpAsUInt32(-7003), 7004, 7005, 7006, 7007, 7008, 7009, 7010, 7011, 7012, 7013, 7014,
					9000, 9001, 9002, 9003, 9004, 9005, 9006, 9007, 9008,
					11000, 11001, 11002, 11003, 11004, 11005, 11006, 11007, 11008,
				},
			},
		},
		{
			name:            "new_info",
			procfsRoot:      path.Join(netSnmpTestDataDir, "field_mapping"),
			primeProcfsRoot: path.Join(netSnmpTestDataDir, "no_icmp_msg"),
			wantNetSnmp: &NetSnmp{
				Values: []uint32{
					1000, 1001, 1002, 1003, 1004, 1005, 1006, 1007, 1008, 1009, 1010, 1011, 1012, 1013, 1014, 1015, 1016, 1017, 1018,
					3000, 3001, 3002, 3003, 3004, 3005, 3006, 3007, 3008, 3009, 3010, 3011, 3012, 3013, 3014, 3015, 3016, 3017, 3018, 3019, 3020, 3021, 3022, 3023, 3024, 3025, 3026,
					5000, 5001,
					7000, 7001, 7002, testNetSnmpAsUInt32(-7003), 7004, 7005, 7006, 7007, 7008, 7009, 7010, 7011, 7012, 7013, 7014,
					9000, 9001, 9002, 9003, 9004, 9005, 9006, 9007, 9008,
					11000, 11001, 11002, 11003, 11004, 11005, 11006, 11007, 11008,
				},
			},
			wantInfoChanged: true,
		},
		{
			name:            "partial_new_info",
			procfsRoot:      path.Join(netSnmpTestDataDir, "field_mapping"),
			primeProcfsRoot: path.Join(netSnmpTestDataDir, "partial_icmp_msg"),
			wantNetSnmp: &NetSnmp{
				Values: []uint32{
					1000, 1001, 1002, 1003, 1004, 1005, 1006, 1007, 1008, 1009, 1010, 1011, 1012, 1013, 1014, 1015, 1016, 1017, 1018,
					3000, 3001, 3002, 3003, 3004, 3005, 3006, 3007, 3008, 3009, 3010, 3011, 3012, 3013, 3014, 3015, 3016, 3017, 3018, 3019, 3020, 3021, 3022, 3023, 3024, 3025, 3026,
					5000, 5001,
					7000, 7001, 7002, testNetSnmpAsUInt32(-7003), 7004, 7005, 7006, 7007, 7008, 7009, 7010, 7011, 7012, 7013, 7014,
					9000, 9001, 9002, 9003, 9004, 9005, 9006, 9007, 9008,
					11000, 11001, 11002, 11003, 11004, 11005, 11006, 11007, 11008,
				},
			},
			wantInfoChanged: true,
		},
	} {
		t.Run(
			tc.name,
			func(t *testing.T) { testNetSnmpParser(tc, t) },
		)
	}
}
