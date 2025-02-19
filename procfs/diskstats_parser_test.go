package procfs

import (
	"bytes"
	"fmt"
	"path"
	"testing"
)

type DiskstatsTestCase struct {
	name                     string
	procfsRoot               string
	primeDiskstats           *Diskstats
	disableJiffiesToMillisec bool
	wantDiskstats            *Diskstats
	wantError                error
}

var diskstatsIndexName = []string{
	"DISKSTATS_NUM_READS_COMPLETED",
	"DISKSTATS_NUM_READS_MERGED",
	"DISKSTATS_NUM_READ_SECTORS",
	"DISKSTATS_READ_MILLISEC",
	"DISKSTATS_NUM_WRITES_COMPLETED",
	"DISKSTATS_NUM_WRITES_MERGED",
	"DISKSTATS_NUM_WRITE_SECTORS",
	"DISKSTATS_WRITE_MILLISEC",
	"DISKSTATS_NUM_IO_IN_PROGRESS",
	"DISKSTATS_IO_MILLISEC",
	"DISKSTATS_IO_WEIGTHED_MILLISEC",
	"DISKSTATS_NUM_DISCARDS_COMPLETED",
	"DISKSTATS_NUM_DISCARDS_MERGED",
	"DISKSTATS_NUM_DISCARD_SECTORS",
	"DISKSTATS_DISCARD_MILLISEC",
	"DISKSTATS_NUM_FLUSH_REQUESTS",
	"DISKSTATS_FLUSH_MILLISEC",
}

var diskstatsTestDataDir = path.Join(PROCFS_TESTDATA_ROOT, "diskstats")

func testDiskstatsParser(tc *DiskstatsTestCase, t *testing.T) {
	t.Logf(`
name=%q
procfsRoot=%q
primeDiskstats=%v
disableJiffiesToMillisec=%v
`,
		tc.name, tc.procfsRoot, (tc.primeDiskstats != nil), tc.disableJiffiesToMillisec,
	)

	var diskstats *Diskstats

	if tc.primeDiskstats != nil {
		diskstats = tc.primeDiskstats.Clone(true)
		diskstats.path = DiskstatsPath(tc.procfsRoot)
	} else {
		diskstats = NewDiskstats(tc.procfsRoot)
		if tc.disableJiffiesToMillisec {
			diskstats.jiffiesToMillisec = 0
		}
	}

	err := diskstats.Parse()
	if tc.wantError != nil {
		if err == nil || tc.wantError.Error() != err.Error() {
			t.Fatalf("want: %v error, got: %v", tc.wantError, err)
		}
		return
	}
	if err != nil {
		t.Fatal(err)
	}

	wantDiskstats := tc.wantDiskstats
	if wantDiskstats == nil {
		return
	}

	diffBuf := &bytes.Buffer{}
	for devMajMin, wantDevInfo := range wantDiskstats.DevInfoMap {
		gotDevInfo := diskstats.DevInfoMap[devMajMin]
		if gotDevInfo == nil {
			fmt.Fprintf(
				diffBuf,
				"\nDevInfoMap[%q]: missing device", devMajMin,
			)
		} else {
			if wantDevInfo.Name != gotDevInfo.Name {
				fmt.Fprintf(
					diffBuf,
					"\\nDevInfoMap[%q].Name: want: %q, got: %q",
					devMajMin, wantDevInfo.Name, gotDevInfo.Name,
				)

			}

			wantDevStats, gotDevStats := wantDevInfo.Stats, gotDevInfo.Stats
			if len(wantDevStats) != len(gotDevStats) {
				fmt.Fprintf(
					diffBuf,
					"\nlen(DevInfoMap[%q].Stats): want: %d, got: %d",
					devMajMin, len(wantDevStats), len(gotDevStats),
				)
			} else {
				for i, wantStat := range wantDevStats {
					gotStat := gotDevStats[i]
					if wantStat != gotStat {
						fmt.Fprintf(
							diffBuf,
							"\nDevInfoMap[%q].Stats[%d (%s)]: want: %d, got: %d",
							devMajMin, i, diskstatsIndexName[i], wantStat, gotStat,
						)
					}
				}
			}
		}
	}
	for devMajMin := range diskstats.DevInfoMap {
		if wantDiskstats.DevInfoMap[devMajMin] == nil {
			fmt.Fprintf(
				diffBuf,
				"\nDevInfoMap[%q]: unexpected device", devMajMin,
			)
		}
	}

	if diffBuf.Len() > 0 {
		t.Fatal(diffBuf.String())
	}
}

func TestDiskstatsParser(t *testing.T) {
	for _, tc := range []*DiskstatsTestCase{
		{
			name:       "field_mapping",
			procfsRoot: path.Join(diskstatsTestDataDir, "field_mapping"),
			wantDiskstats: &Diskstats{
				DevInfoMap: map[string]*DiskstatsDevInfo{
					"0:0": {
						Name:  "disk0",
						Stats: []uint32{1001, 1002, 1003, 1004, 1005, 1006, 1007, 1008, 1009, 1010, 1011, 1012, 1013, 1014, 1015, 1016, 1017},
					},
					"1:1": {
						Name:  "disk1",
						Stats: []uint32{2001, 2002, 2003, 2004, 2005, 2006, 2007, 2008, 2009, 2010, 2011, 2012, 2013, 2014, 2015, 2016, 2017},
					},
				},
			},
			disableJiffiesToMillisec: true,
		},
		{
			name:       "reuse",
			procfsRoot: path.Join(diskstatsTestDataDir, "field_mapping"),
			primeDiskstats: &Diskstats{
				DevInfoMap: map[string]*DiskstatsDevInfo{
					"0:0": {
						Name:    "disk0",
						Stats:   make([]uint32, 17),
						scanNum: 42,
					},
					"1:1": {
						Name:    "disk1",
						Stats:   make([]uint32, 17),
						scanNum: 42,
					},
				},
				scanNum:           42,
				jiffiesToMillisec: 0,
				fieldsInJiffies:   diskstatsFieldsInJiffies,
			},
			wantDiskstats: &Diskstats{
				DevInfoMap: map[string]*DiskstatsDevInfo{
					"0:0": {
						Name:  "disk0",
						Stats: []uint32{1001, 1002, 1003, 1004, 1005, 1006, 1007, 1008, 1009, 1010, 1011, 1012, 1013, 1014, 1015, 1016, 1017},
					},
					"1:1": {
						Name:  "disk1",
						Stats: []uint32{2001, 2002, 2003, 2004, 2005, 2006, 2007, 2008, 2009, 2010, 2011, 2012, 2013, 2014, 2015, 2016, 2017},
					},
				},
			},
		},
		{
			name:       "remove_dev",
			procfsRoot: path.Join(diskstatsTestDataDir, "field_mapping"),
			primeDiskstats: &Diskstats{
				DevInfoMap: map[string]*DiskstatsDevInfo{
					"0:0": {
						Name:    "disk0",
						Stats:   make([]uint32, 17),
						scanNum: 42,
					},
					"1:1": {
						Name:    "disk1",
						Stats:   make([]uint32, 17),
						scanNum: 42,
					},
					"255:255": {
						Name:    "removed",
						Stats:   make([]uint32, 17),
						scanNum: 42,
					},
				},
				scanNum:           42,
				jiffiesToMillisec: 0,
				fieldsInJiffies:   diskstatsFieldsInJiffies,
			},
			wantDiskstats: &Diskstats{
				DevInfoMap: map[string]*DiskstatsDevInfo{
					"0:0": {
						Name:  "disk0",
						Stats: []uint32{1001, 1002, 1003, 1004, 1005, 1006, 1007, 1008, 1009, 1010, 1011, 1012, 1013, 1014, 1015, 1016, 1017},
					},
					"1:1": {
						Name:  "disk1",
						Stats: []uint32{2001, 2002, 2003, 2004, 2005, 2006, 2007, 2008, 2009, 2010, 2011, 2012, 2013, 2014, 2015, 2016, 2017},
					},
				},
			},
		},
		{
			name:       "jiffies",
			procfsRoot: path.Join(diskstatsTestDataDir, "field_mapping"),
			primeDiskstats: &Diskstats{
				DevInfoMap: map[string]*DiskstatsDevInfo{
					"0:0": {
						Name:    "disk0",
						Stats:   make([]uint32, 17),
						scanNum: 42,
					},
					"1:1": {
						Name:    "disk1",
						Stats:   make([]uint32, 17),
						scanNum: 42,
					},
				},
				scanNum:           42,
				jiffiesToMillisec: 10,
				fieldsInJiffies:   diskstatsFieldsInJiffies,
			},
			wantDiskstats: &Diskstats{
				DevInfoMap: map[string]*DiskstatsDevInfo{
					"0:0": {
						Name:  "disk0",
						Stats: []uint32{1001, 1002, 1003, 1004, 1005, 1006, 1007, 1008, 1009, 10100, 1011, 1012, 1013, 1014, 1015, 1016, 1017},
					},
					"1:1": {
						Name:  "disk1",
						Stats: []uint32{2001, 2002, 2003, 2004, 2005, 2006, 2007, 2008, 2009, 20100, 2011, 2012, 2013, 2014, 2015, 2016, 2017},
					},
				},
			},
		},
	} {
		t.Run(
			tc.name,
			func(t *testing.T) { testDiskstatsParser(tc, t) },
		)
	}
}
