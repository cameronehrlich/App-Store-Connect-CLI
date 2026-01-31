package asc

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func printBackgroundAssetsTable(resp *BackgroundAssetsResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tAsset Pack Identifier\tArchived\tCreated Date")
	for _, item := range resp.Data {
		fmt.Fprintf(w, "%s\t%s\t%t\t%s\n",
			item.ID,
			compactWhitespace(item.Attributes.AssetPackIdentifier),
			item.Attributes.Archived,
			item.Attributes.CreatedDate,
		)
	}
	return w.Flush()
}

func printBackgroundAssetsMarkdown(resp *BackgroundAssetsResponse) error {
	fmt.Fprintln(os.Stdout, "| ID | Asset Pack Identifier | Archived | Created Date |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- |")
	for _, item := range resp.Data {
		fmt.Fprintf(os.Stdout, "| %s | %s | %t | %s |\n",
			escapeMarkdown(item.ID),
			escapeMarkdown(item.Attributes.AssetPackIdentifier),
			item.Attributes.Archived,
			escapeMarkdown(item.Attributes.CreatedDate),
		)
	}
	return nil
}

func printBackgroundAssetVersionsTable(resp *BackgroundAssetVersionsResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tVersion\tState\tPlatforms\tCreated Date")
	for _, item := range resp.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			item.ID,
			compactWhitespace(item.Attributes.Version),
			compactWhitespace(item.Attributes.State),
			formatPlatforms(item.Attributes.Platforms),
			item.Attributes.CreatedDate,
		)
	}
	return w.Flush()
}

func printBackgroundAssetVersionsMarkdown(resp *BackgroundAssetVersionsResponse) error {
	fmt.Fprintln(os.Stdout, "| ID | Version | State | Platforms | Created Date |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- | --- |")
	for _, item := range resp.Data {
		fmt.Fprintf(os.Stdout, "| %s | %s | %s | %s | %s |\n",
			escapeMarkdown(item.ID),
			escapeMarkdown(item.Attributes.Version),
			escapeMarkdown(item.Attributes.State),
			escapeMarkdown(formatPlatforms(item.Attributes.Platforms)),
			escapeMarkdown(item.Attributes.CreatedDate),
		)
	}
	return nil
}

func printBackgroundAssetUploadFilesTable(resp *BackgroundAssetUploadFilesResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tFile Name\tAsset Type\tFile Size\tState")
	for _, item := range resp.Data {
		state := ""
		if item.Attributes.AssetDeliveryState != nil && item.Attributes.AssetDeliveryState.State != nil {
			state = *item.Attributes.AssetDeliveryState.State
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n",
			item.ID,
			compactWhitespace(item.Attributes.FileName),
			string(item.Attributes.AssetType),
			item.Attributes.FileSize,
			state,
		)
	}
	return w.Flush()
}

func printBackgroundAssetUploadFilesMarkdown(resp *BackgroundAssetUploadFilesResponse) error {
	fmt.Fprintln(os.Stdout, "| ID | File Name | Asset Type | File Size | State |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- | --- |")
	for _, item := range resp.Data {
		state := ""
		if item.Attributes.AssetDeliveryState != nil && item.Attributes.AssetDeliveryState.State != nil {
			state = *item.Attributes.AssetDeliveryState.State
		}
		fmt.Fprintf(os.Stdout, "| %s | %s | %s | %d | %s |\n",
			escapeMarkdown(item.ID),
			escapeMarkdown(item.Attributes.FileName),
			escapeMarkdown(string(item.Attributes.AssetType)),
			item.Attributes.FileSize,
			escapeMarkdown(strings.TrimSpace(state)),
		)
	}
	return nil
}

func printBackgroundAssetVersionStateTable(id string, state string) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tState")
	fmt.Fprintf(w, "%s\t%s\n", id, state)
	return w.Flush()
}

func printBackgroundAssetVersionStateMarkdown(id string, state string) error {
	fmt.Fprintln(os.Stdout, "| ID | State |")
	fmt.Fprintln(os.Stdout, "| --- | --- |")
	fmt.Fprintf(os.Stdout, "| %s | %s |\n",
		escapeMarkdown(id),
		escapeMarkdown(state),
	)
	return nil
}

func printBackgroundAssetVersionAppStoreReleaseTable(resp *BackgroundAssetVersionAppStoreReleaseResponse) error {
	return printBackgroundAssetVersionStateTable(resp.Data.ID, resp.Data.Attributes.State)
}

func printBackgroundAssetVersionAppStoreReleaseMarkdown(resp *BackgroundAssetVersionAppStoreReleaseResponse) error {
	return printBackgroundAssetVersionStateMarkdown(resp.Data.ID, resp.Data.Attributes.State)
}

func printBackgroundAssetVersionExternalBetaReleaseTable(resp *BackgroundAssetVersionExternalBetaReleaseResponse) error {
	return printBackgroundAssetVersionStateTable(resp.Data.ID, resp.Data.Attributes.State)
}

func printBackgroundAssetVersionExternalBetaReleaseMarkdown(resp *BackgroundAssetVersionExternalBetaReleaseResponse) error {
	return printBackgroundAssetVersionStateMarkdown(resp.Data.ID, resp.Data.Attributes.State)
}

func printBackgroundAssetVersionInternalBetaReleaseTable(resp *BackgroundAssetVersionInternalBetaReleaseResponse) error {
	return printBackgroundAssetVersionStateTable(resp.Data.ID, resp.Data.Attributes.State)
}

func printBackgroundAssetVersionInternalBetaReleaseMarkdown(resp *BackgroundAssetVersionInternalBetaReleaseResponse) error {
	return printBackgroundAssetVersionStateMarkdown(resp.Data.ID, resp.Data.Attributes.State)
}
