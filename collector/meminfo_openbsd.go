// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build openbsd
// +build !nomeminfo

package collector

import (
	"fmt"
)

/*
#include <sys/param.h>
#include <sys/types.h>
#include <sys/sysctl.h>

int
sysctl_uvmexp(struct uvmexp *uvmexp)
{
        static int uvmexp_mib[] = {CTL_VM, VM_UVMEXP};
        size_t sz = sizeof(struct uvmexp);

        if(sysctl(uvmexp_mib, 2, uvmexp, &sz, NULL, 0) < 0)
                return -1;

        return 0;
}

*/
import "C"

func (c *meminfoCollector) getMemInfo() (map[string]float64, error) {
	var uvmexp C.struct_uvmexp

	if _, err := C.sysctl_uvmexp(&uvmexp); err != nil {
		return nil, fmt.Errorf("sysctl CTL_VM VM_UVMEXP failed: %v", err)
	}

	ps := uvmexp.pagesize

	// see uvm(9)
	return map[string]float64{
		"active_bytes":                  float64(ps * uvmexp.active),
		"cache_bytes":                   float64(ps * uvmexp.vnodepages),
		"free_bytes":                    float64(ps * uvmexp.free),
		"inactive_bytes":                float64(ps * uvmexp.inactive),
		"swap_size_bytes":               float64(ps * uvmexp.swpages),
		"swap_used_bytes":               float64(ps * uvmexp.swpgonly),
		"swapped_in_pages_bytes_total":  float64(ps * uvmexp.pgswapin),
		"swapped_out_pages_bytes_total": float64(ps * uvmexp.pgswapout),
		"total_bytes":                   float64(ps * uvmexp.npages),
		"wired_bytes":                   float64(ps * uvmexp.wired),
	}, nil
}
