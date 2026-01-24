package asc

import (
	"fmt"
	"os"
	"text/tabwriter"
)

// InAppPurchaseDeleteResult represents CLI output for IAP deletions.
type InAppPurchaseDeleteResult struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

func printInAppPurchasesTable(resp *InAppPurchasesV2Response) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tName\tProduct ID\tType\tState")
	for _, item := range resp.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			item.ID,
			compactWhitespace(item.Attributes.Name),
			item.Attributes.ProductID,
			item.Attributes.InAppPurchaseType,
			item.Attributes.State,
		)
	}
	return w.Flush()
}

func printInAppPurchasesMarkdown(resp *InAppPurchasesV2Response) error {
	fmt.Fprintln(os.Stdout, "| ID | Name | Product ID | Type | State |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- | --- |")
	for _, item := range resp.Data {
		fmt.Fprintf(os.Stdout, "| %s | %s | %s | %s | %s |\n",
			escapeMarkdown(item.ID),
			escapeMarkdown(item.Attributes.Name),
			escapeMarkdown(item.Attributes.ProductID),
			escapeMarkdown(item.Attributes.InAppPurchaseType),
			escapeMarkdown(item.Attributes.State),
		)
	}
	return nil
}

func printInAppPurchaseLocalizationsTable(resp *InAppPurchaseLocalizationsResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tLocale\tName\tDescription")
	for _, item := range resp.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			item.ID,
			item.Attributes.Locale,
			compactWhitespace(item.Attributes.Name),
			compactWhitespace(item.Attributes.Description),
		)
	}
	return w.Flush()
}

func printInAppPurchaseLocalizationsMarkdown(resp *InAppPurchaseLocalizationsResponse) error {
	fmt.Fprintln(os.Stdout, "| ID | Locale | Name | Description |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- |")
	for _, item := range resp.Data {
		fmt.Fprintf(os.Stdout, "| %s | %s | %s | %s |\n",
			escapeMarkdown(item.ID),
			escapeMarkdown(item.Attributes.Locale),
			escapeMarkdown(item.Attributes.Name),
			escapeMarkdown(item.Attributes.Description),
		)
	}
	return nil
}

func printInAppPurchaseDeleteResultTable(result *InAppPurchaseDeleteResult) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tDeleted")
	fmt.Fprintf(w, "%s\t%t\n", result.ID, result.Deleted)
	return w.Flush()
}

func printInAppPurchaseDeleteResultMarkdown(result *InAppPurchaseDeleteResult) error {
	fmt.Fprintln(os.Stdout, "| ID | Deleted |")
	fmt.Fprintln(os.Stdout, "| --- | --- |")
	fmt.Fprintf(os.Stdout, "| %s | %t |\n",
		escapeMarkdown(result.ID),
		result.Deleted,
	)
	return nil
}
