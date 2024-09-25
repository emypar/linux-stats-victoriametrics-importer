#! /usr/bin/env python3

from dataclasses import dataclass, field
from typing import List

# JSON serialize-able NetSnmp, matching profcs/net_snmp_parser.go:


NET_SNMP_IP_FORWARDING = 0
NET_SNMP_IP_DEFAULT_TTL = 1
NET_SNMP_IP_IN_RECEIVES = 2
NET_SNMP_IP_IN_HDR_ERRORS = 3
NET_SNMP_IP_IN_ADDR_ERRORS = 4
NET_SNMP_IP_FORW_DATAGRAMS = 5
NET_SNMP_IP_IN_UNKNOWN_PROTOS = 6
NET_SNMP_IP_IN_DISCARDS = 7
NET_SNMP_IP_IN_DELIVERS = 8
NET_SNMP_IP_OUT_REQUESTS = 9
NET_SNMP_IP_OUT_DISCARDS = 10
NET_SNMP_IP_OUT_NO_ROUTES = 11
NET_SNMP_IP_REASM_TIMEOUT = 12
NET_SNMP_IP_REASM_REQDS = 13
NET_SNMP_IP_REASM_OKS = 14
NET_SNMP_IP_REASM_FAILS = 15
NET_SNMP_IP_FRAG_OKS = 16
NET_SNMP_IP_FRAG_FAILS = 17
NET_SNMP_IP_FRAG_CREATES = 18
NET_SNMP_ICMP_IN_MSGS = 19
NET_SNMP_ICMP_IN_ERRORS = 20
NET_SNMP_ICMP_IN_CSUM_ERRORS = 21
NET_SNMP_ICMP_IN_DEST_UNREACHS = 22
NET_SNMP_ICMP_IN_TIME_EXCDS = 23
NET_SNMP_ICMP_IN_PARM_PROBS = 24
NET_SNMP_ICMP_IN_SRC_QUENCHS = 25
NET_SNMP_ICMP_IN_REDIRECTS = 26
NET_SNMP_ICMP_IN_ECHOS = 27
NET_SNMP_ICMP_IN_ECHO_REPS = 28
NET_SNMP_ICMP_IN_TIMESTAMPS = 29
NET_SNMP_ICMP_IN_TIMESTAMP_REPS = 30
NET_SNMP_ICMP_IN_ADDR_MASKS = 31
NET_SNMP_ICMP_IN_ADDR_MASK_REPS = 32
NET_SNMP_ICMP_OUT_MSGS = 33
NET_SNMP_ICMP_OUT_ERRORS = 34
NET_SNMP_ICMP_OUT_DEST_UNREACHS = 35
NET_SNMP_ICMP_OUT_TIME_EXCDS = 36
NET_SNMP_ICMP_OUT_PARM_PROBS = 37
NET_SNMP_ICMP_OUT_SRC_QUENCHS = 38
NET_SNMP_ICMP_OUT_REDIRECTS = 39
NET_SNMP_ICMP_OUT_ECHOS = 40
NET_SNMP_ICMP_OUT_ECHO_REPS = 41
NET_SNMP_ICMP_OUT_TIMESTAMPS = 42
NET_SNMP_ICMP_OUT_TIMESTAMP_REPS = 43
NET_SNMP_ICMP_OUT_ADDR_MASKS = 44
NET_SNMP_ICMP_OUT_ADDR_MASK_REPS = 45
NET_SNMP_ICMPMSG_IN_TYPE3 = 46
NET_SNMP_ICMPMSG_OUT_TYPE3 = 47
NET_SNMP_TCP_RTO_ALGORITHM = 48
NET_SNMP_TCP_RTO_MIN = 49
NET_SNMP_TCP_RTO_MAX = 50
NET_SNMP_TCP_MAX_CONN = 51
NET_SNMP_TCP_ACTIVE_OPENS = 52
NET_SNMP_TCP_PASSIVE_OPENS = 53
NET_SNMP_TCP_ATTEMPT_FAILS = 54
NET_SNMP_TCP_ESTAB_RESETS = 55
NET_SNMP_TCP_CURR_ESTAB = 56
NET_SNMP_TCP_IN_SEGS = 57
NET_SNMP_TCP_OUT_SEGS = 58
NET_SNMP_TCP_RETRANS_SEGS = 59
NET_SNMP_TCP_IN_ERRS = 60
NET_SNMP_TCP_OUT_RSTS = 61
NET_SNMP_TCP_IN_CSUM_ERRORS = 62
NET_SNMP_UDP_IN_DATAGRAMS = 63
NET_SNMP_UDP_NO_PORTS = 64
NET_SNMP_UDP_IN_ERRORS = 65
NET_SNMP_UDP_OUT_DATAGRAMS = 66
NET_SNMP_UDP_RCVBUF_ERRORS = 67
NET_SNMP_UDP_SNDBUF_ERRORS = 68
NET_SNMP_UDP_IN_CSUM_ERRORS = 69
NET_SNMP_UDP_IGNORED_MULTI = 70
NET_SNMP_UDP_MEM_ERRORS = 71
NET_SNMP_UDPLITE_IN_DATAGRAMS = 72
NET_SNMP_UDPLITE_NO_PORTS = 73
NET_SNMP_UDPLITE_IN_ERRORS = 74
NET_SNMP_UDPLITE_OUT_DATAGRAMS = 75
NET_SNMP_UDPLITE_RCVBUF_ERRORS = 76
NET_SNMP_UDPLITE_SNDBUF_ERRORS = 77
NET_SNMP_UDPLITE_IN_CSUM_ERRORS = 78
NET_SNMP_UDPLITE_IGNORED_MULTI = 79
NET_SNMP_UDPLITE_MEM_ERRORS = 80

NET_SNMP_NUM_VALUES = 81

NetSnmpValueMayBeNegative = {
    NET_SNMP_TCP_MAX_CONN,
}


@dataclass
class NetSnmp:
    Values: List[int] = field(default_factory=lambda: [0] * NET_SNMP_NUM_VALUES)