package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

type brokenRef struct {
	Source string // address of the file containing the reference
	Field  string // PARENT, CHILDREN, REFERENCE, or WAYPOINTS
	Target string // the broken address
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Check all element references for broken links",
	Long: `Scan all WayPoint and Map files, extract referenced addresses,
and verify that each referenced element exists on disk.

Reports any broken references (addresses pointing to non-existent files).

Examples:
  loadstar validate`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		loadstarBase := fs.AvcsPath("")
		var broken []brokenRef
		wpCount := 0
		mapCount := 0

		// Scan WAYPOINT files
		wpDir := filepath.Join(loadstarBase, "WAYPOINT")
		if files, err := os.ReadDir(wpDir); err == nil {
			for _, f := range files {
				if f.IsDir() || !strings.HasSuffix(f.Name(), ".md") {
					continue
				}
				wpCount++
				dotName := strings.TrimSuffix(f.Name(), ".md")
				path := strings.ReplaceAll(dotName, ".", "/")
				sourceAddr := "W://" + path

				content, err := os.ReadFile(filepath.Join(wpDir, f.Name()))
				if err != nil {
					continue
				}
				text := string(content)

				// Check PARENT
				if parent := extractSingleField(text, `(?m)^-\s*PARENT:\s*(.+)$`); parent != "" {
					if !elementExists(loadstarBase, parent) {
						broken = append(broken, brokenRef{sourceAddr, "PARENT", parent})
					}
				}

				// Check CHILDREN
				for _, child := range extractListField(text, `(?m)-\s*CHILDREN:\s*\[([^\]]*)\]`) {
					if !elementExists(loadstarBase, child) {
						broken = append(broken, brokenRef{sourceAddr, "CHILDREN", child})
					}
				}

				// Check REFERENCE
				for _, ref := range extractListField(text, `(?m)-\s*REFERENCE:\s*\[([^\]]*)\]`) {
					if !elementExists(loadstarBase, ref) {
						broken = append(broken, brokenRef{sourceAddr, "REFERENCE", ref})
					}
				}

				// Check CONFIRMED Q decision file references
				qConfirmedRe := regexp.MustCompile(`\[Q\d+\s+CONFIRMED\s+([\w.\-]+)\]`)
				for _, m := range qConfirmedRe.FindAllStringSubmatch(text, -1) {
					decisionRef := m[1]
					decisionFile := filepath.Join(loadstarBase, "DECISIONS", decisionRef+".md")
					if !fs.Exists(decisionFile) {
						broken = append(broken, brokenRef{sourceAddr, "OPEN_QUESTIONS", "DECISIONS/" + decisionRef + ".md"})
					}
				}
			}
		}

		// Scan MAP files
		mapDir := filepath.Join(loadstarBase, "MAP")
		if files, err := os.ReadDir(mapDir); err == nil {
			for _, f := range files {
				if f.IsDir() || !strings.HasSuffix(f.Name(), ".md") {
					continue
				}
				mapCount++
				dotName := strings.TrimSuffix(f.Name(), ".md")
				path := strings.ReplaceAll(dotName, ".", "/")
				sourceAddr := "M://" + path

				content, err := os.ReadFile(filepath.Join(mapDir, f.Name()))
				if err != nil {
					continue
				}

				for _, child := range extractMapWaypoints(string(content)) {
					if !elementExists(loadstarBase, child) {
						broken = append(broken, brokenRef{sourceAddr, "WAYPOINTS", child})
					}
				}
			}
		}

		if len(broken) == 0 {
			fmt.Printf("all references valid (%d waypoints, %d maps checked)\n", wpCount, mapCount)
			return
		}

		fmt.Printf("found %d broken reference(s):\n\n", len(broken))
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "SOURCE\tFIELD\tBROKEN_ADDRESS")
		fmt.Fprintln(w, "------\t-----\t--------------")
		for _, b := range broken {
			fmt.Fprintf(w, "%s\t%s\t%s\n", b.Source, b.Field, b.Target)
		}
		w.Flush()
	},
}

// elementExists checks if a LOADSTAR address has a corresponding file on disk.
func elementExists(loadstarBase, addrStr string) bool {
	addrStr = strings.TrimSpace(addrStr)
	if addrStr == "" {
		return true // empty is not broken
	}
	addr, err := svc.ParseAddress(addrStr)
	if err != nil {
		return false
	}
	filePath := addr.ToFilePath(loadstarBase)
	return fs.Exists(filePath)
}

// extractSingleField extracts a single address from a regex match.
func extractSingleField(content, pattern string) string {
	re := regexp.MustCompile(pattern)
	m := re.FindStringSubmatch(content)
	if len(m) < 2 {
		return ""
	}
	val := strings.TrimSpace(m[1])
	if val == "" || val == "(없음)" {
		return ""
	}
	return val
}

// extractListField extracts a comma-separated list of addresses from a bracket field.
func extractListField(content, pattern string) []string {
	re := regexp.MustCompile(pattern)
	m := re.FindStringSubmatch(content)
	if len(m) < 2 {
		return nil
	}
	raw := strings.TrimSpace(m[1])
	if raw == "" {
		return nil
	}
	var result []string
	for _, item := range strings.Split(raw, ",") {
		item = strings.TrimSpace(item)
		if item != "" && strings.Contains(item, "://") {
			result = append(result, item)
		}
	}
	return result
}

// extractMapWaypoints extracts addresses from a MAP's WAYPOINTS section.
func extractMapWaypoints(content string) []string {
	var result []string
	lines := strings.Split(content, "\n")
	inWaypoints := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "### WAYPOINTS" {
			inWaypoints = true
			continue
		}
		if inWaypoints {
			if strings.HasPrefix(trimmed, "###") || (trimmed == "" && len(result) > 0) {
				break
			}
			if strings.HasPrefix(trimmed, "- ") && strings.Contains(trimmed, "://") {
				addr := strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))
				result = append(result, addr)
			}
			if trimmed == "(없음)" {
				break
			}
		}
	}
	return result
}
