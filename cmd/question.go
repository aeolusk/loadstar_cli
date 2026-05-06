package cmd

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// qRe matches OPEN_QUESTIONS lines:
//
//	[Q1], [Q1 DEFERRED], [Q1 CONFIRMED], [Q1 CONFIRMED some-ref]
var qRe = regexp.MustCompile(`^\s*-\s*\[Q(\d+)(?:\s+(DEFERRED|CONFIRMED)(?:\s+([\w.\-]+))?)?\]\s*(.*)$`)

type qEntry struct {
	Address  string
	QID      string
	State    string
	Ref      string // decision file ref (for CONFIRMED with Decision file)
	Question string
}

func scanQuestions(loadstarBase string) []qEntry {
	wpAddrs := collectAllWaypoints(loadstarBase)
	var entries []qEntry

	for _, addr := range wpAddrs {
		wpFile := addressToFilePath(loadstarBase, addr)
		data, err := os.ReadFile(wpFile)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			m := qRe.FindStringSubmatch(line)
			if m == nil {
				continue
			}
			state := "OPEN"
			switch strings.TrimSpace(m[2]) {
			case "DEFERRED":
				state = "DEFERRED"
			case "CONFIRMED":
				state = "CONFIRMED"
			}
			entries = append(entries, qEntry{
				Address:  addr,
				QID:      "Q" + m[1],
				State:    state,
				Ref:      strings.TrimSpace(m[3]),
				Question: strings.TrimSpace(m[4]),
			})
		}
	}
	return entries
}

var withAll bool

var questionCmd = &cobra.Command{
	Use:   "question [FILTER]",
	Short: "List OPEN/DEFERRED questions from all WayPoints",
	Long: `Scan all WayPoint files and list OPEN and DEFERRED OPEN_QUESTIONS.
Use --all to include CONFIRMED items as well.

Optional FILTER matches address or question text (case-insensitive substring).

To confirm a question (mark as CONFIRMED), edit the WayPoint file directly or
use the UI Questions panel — there is no write subcommand by design.

Examples:
  loadstar question
  loadstar question --all
  loadstar question M://root/maintenance`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		loadstarBase := fs.AvcsPath("")
		entries := scanQuestions(loadstarBase)

		filter := ""
		if len(args) == 1 {
			filter = strings.ToLower(args[0])
		}

		var visible []qEntry
		for _, e := range entries {
			if !withAll && e.State == "CONFIRMED" {
				continue
			}
			if filter != "" &&
				!strings.Contains(strings.ToLower(e.Address), filter) &&
				!strings.Contains(strings.ToLower(e.Question), filter) {
				continue
			}
			visible = append(visible, e)
		}

		if len(visible) == 0 {
			fmt.Println("no questions found")
			return
		}

		stateOrder := map[string]int{"OPEN": 0, "DEFERRED": 1, "CONFIRMED": 2}
		sort.Slice(visible, func(i, j int) bool {
			oi, oj := stateOrder[visible[i].State], stateOrder[visible[j].State]
			if oi != oj {
				return oi < oj
			}
			if visible[i].Address != visible[j].Address {
				return visible[i].Address < visible[j].Address
			}
			return visible[i].QID < visible[j].QID
		})

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ADDRESS\tQID\tSTATE\tREF\tQUESTION")
		fmt.Fprintln(w, "-------\t---\t-----\t---\t--------")
		for _, e := range visible {
			q := e.Question
			if len(q) > 70 {
				q = q[:70] + "..."
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", e.Address, e.QID, e.State, e.Ref, q)
		}
		w.Flush()
		fmt.Printf("\n%d question(s)\n", len(visible))
	},
}

var questionStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show OPEN_QUESTIONS statistics",
	Run: func(cmd *cobra.Command, args []string) {
		loadstarBase := fs.AvcsPath("")
		entries := scanQuestions(loadstarBase)

		counts := map[string]int{"OPEN": 0, "DEFERRED": 0, "CONFIRMED": 0}
		for _, e := range entries {
			counts[e.State]++
		}

		total := counts["OPEN"] + counts["DEFERRED"] + counts["CONFIRMED"]
		fmt.Printf("OPEN:      %d\n", counts["OPEN"])
		fmt.Printf("DEFERRED:  %d\n", counts["DEFERRED"])
		fmt.Printf("CONFIRMED: %d\n", counts["CONFIRMED"])
		fmt.Printf("TOTAL:     %d\n", total)
	},
}

func init() {
	questionCmd.Flags().BoolVar(&withAll, "all", false, "include CONFIRMED items")
	questionCmd.AddCommand(questionStatsCmd)
}
