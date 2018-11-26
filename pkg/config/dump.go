package config

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// DumpConfig try to dump config in proper formatted text
func DumpConfig(w io.Writer) {
	lock.Lock()
	defer lock.Unlock()

	tab := tabwriter.NewWriter(w, 0, 8, 0, '\t', 0)
	_, _ = fmt.Fprint(w, "Key\tDescription\tField\tValue\n")
	for key := range configs {
		d, ok := o.Get(key)
		_, _ = fmt.Fprintf(tab, "%s\t%s\t%v\t%v\n", key, configs[key], ok, d)
	}
	_ = tab.Flush()
}
