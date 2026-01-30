package asc

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func printAppInfosTable(resp *AppInfosResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tApp Store State\tState\tAge Rating\tKids Age Band")
	for _, info := range resp.Data {
		attrs := info.Attributes
		fmt.Fprintf(
			w,
			"%s\t%s\t%s\t%s\t%s\n",
			info.ID,
			appInfoAttrString(attrs, "appStoreState"),
			appInfoAttrString(attrs, "state"),
			appInfoAttrString(attrs, "appStoreAgeRating"),
			appInfoAttrString(attrs, "kidsAgeBand"),
		)
	}
	return w.Flush()
}

func printAppInfosMarkdown(resp *AppInfosResponse) error {
	fmt.Fprintln(os.Stdout, "| ID | App Store State | State | Age Rating | Kids Age Band |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- | --- |")
	for _, info := range resp.Data {
		attrs := info.Attributes
		fmt.Fprintf(
			os.Stdout,
			"| %s | %s | %s | %s | %s |\n",
			escapeMarkdown(info.ID),
			escapeMarkdown(appInfoAttrString(attrs, "appStoreState")),
			escapeMarkdown(appInfoAttrString(attrs, "state")),
			escapeMarkdown(appInfoAttrString(attrs, "appStoreAgeRating")),
			escapeMarkdown(appInfoAttrString(attrs, "kidsAgeBand")),
		)
	}
	return nil
}

func appInfoAttrString(attrs AppInfoAttributes, key string) string {
	if attrs == nil {
		return ""
	}
	value, ok := attrs[key]
	if !ok || value == nil {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return typed
	default:
		return fmt.Sprintf("%v", typed)
	}
}
