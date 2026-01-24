package asc

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func formatPersonName(firstName, lastName string) string {
	first := strings.TrimSpace(firstName)
	last := strings.TrimSpace(lastName)
	switch {
	case first == "" && last == "":
		return ""
	case first == "":
		return last
	case last == "":
		return first
	default:
		return first + " " + last
	}
}

func formatUserUsername(attr UserAttributes) string {
	username := strings.TrimSpace(attr.Username)
	if username != "" {
		return username
	}
	return strings.TrimSpace(attr.Email)
}

func printUsersTable(resp *UsersResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tUsername\tName\tRoles\tAll Apps\tProvisioning")
	for _, item := range resp.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%t\t%t\n",
			item.ID,
			compactWhitespace(formatUserUsername(item.Attributes)),
			compactWhitespace(formatPersonName(item.Attributes.FirstName, item.Attributes.LastName)),
			compactWhitespace(strings.Join(item.Attributes.Roles, ",")),
			item.Attributes.AllAppsVisible,
			item.Attributes.ProvisioningAllowed,
		)
	}
	return w.Flush()
}

func printUsersMarkdown(resp *UsersResponse) error {
	fmt.Fprintln(os.Stdout, "| ID | Username | Name | Roles | All Apps | Provisioning |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- | --- | --- |")
	for _, item := range resp.Data {
		fmt.Fprintf(os.Stdout, "| %s | %s | %s | %s | %t | %t |\n",
			escapeMarkdown(item.ID),
			escapeMarkdown(formatUserUsername(item.Attributes)),
			escapeMarkdown(formatPersonName(item.Attributes.FirstName, item.Attributes.LastName)),
			escapeMarkdown(strings.Join(item.Attributes.Roles, ",")),
			item.Attributes.AllAppsVisible,
			item.Attributes.ProvisioningAllowed,
		)
	}
	return nil
}

func printUserInvitationsTable(resp *UserInvitationsResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tEmail\tName\tRoles\tAll Apps\tProvisioning\tExpires")
	for _, item := range resp.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%t\t%t\t%s\n",
			item.ID,
			compactWhitespace(item.Attributes.Email),
			compactWhitespace(formatPersonName(item.Attributes.FirstName, item.Attributes.LastName)),
			compactWhitespace(strings.Join(item.Attributes.Roles, ",")),
			item.Attributes.AllAppsVisible,
			item.Attributes.ProvisioningAllowed,
			compactWhitespace(item.Attributes.ExpirationDate),
		)
	}
	return w.Flush()
}

func printUserInvitationsMarkdown(resp *UserInvitationsResponse) error {
	fmt.Fprintln(os.Stdout, "| ID | Email | Name | Roles | All Apps | Provisioning | Expires |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- | --- | --- | --- |")
	for _, item := range resp.Data {
		fmt.Fprintf(os.Stdout, "| %s | %s | %s | %s | %t | %t | %s |\n",
			escapeMarkdown(item.ID),
			escapeMarkdown(item.Attributes.Email),
			escapeMarkdown(formatPersonName(item.Attributes.FirstName, item.Attributes.LastName)),
			escapeMarkdown(strings.Join(item.Attributes.Roles, ",")),
			item.Attributes.AllAppsVisible,
			item.Attributes.ProvisioningAllowed,
			escapeMarkdown(item.Attributes.ExpirationDate),
		)
	}
	return nil
}

func printUserDeleteResultTable(result *UserDeleteResult) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tDeleted")
	fmt.Fprintf(w, "%s\t%t\n", result.ID, result.Deleted)
	return w.Flush()
}

func printUserDeleteResultMarkdown(result *UserDeleteResult) error {
	fmt.Fprintln(os.Stdout, "| ID | Deleted |")
	fmt.Fprintln(os.Stdout, "| --- | --- |")
	fmt.Fprintf(os.Stdout, "| %s | %t |\n", escapeMarkdown(result.ID), result.Deleted)
	return nil
}

func printUserInvitationRevokeResultTable(result *UserInvitationRevokeResult) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tRevoked")
	fmt.Fprintf(w, "%s\t%t\n", result.ID, result.Revoked)
	return w.Flush()
}

func printUserInvitationRevokeResultMarkdown(result *UserInvitationRevokeResult) error {
	fmt.Fprintln(os.Stdout, "| ID | Revoked |")
	fmt.Fprintln(os.Stdout, "| --- | --- |")
	fmt.Fprintf(os.Stdout, "| %s | %t |\n", escapeMarkdown(result.ID), result.Revoked)
	return nil
}
