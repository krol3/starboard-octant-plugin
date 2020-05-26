package view

import (
	"fmt"
	"sort"
	"strconv"

	starboard "github.com/aquasecurity/starboard/pkg/apis/aquasecurity/v1alpha1"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func NewPolarisReport(report *starboard.ConfigAuditReport) component.Component {
	if report == nil {
		return component.NewText("No report. Run kubectl starboard polairs")
	}

	flexLayout := component.NewFlexLayout("")
	flexLayout.AddSections(component.FlexLayoutSection{
		{
			Width: component.WidthFull,
			View:  component.NewMarkdownText("#### Config Audit Reports"),
		},
		{
			Width: component.WidthThird,
			View:  NewReportSummary(report.GetCreationTimestamp().Time),
		},
		{
			Width: component.WidthThird,
			View:  NewScannerSummary(report.Report.Scanner),
		},
		{
			Width: component.WidthThird,
			View:  createChecksSummary(report.Report),
		},
	})

	var items []component.FlexLayoutItem
	items = append(items, component.FlexLayoutItem{
		Width: component.WidthThird,
		View:  createCardComponent("Pod Template", report.Report.PodChecks),
	})

	var containerNames []string

	for containerName := range report.Report.ContainerChecks {
		containerNames = append(containerNames, containerName)
	}

	sort.Strings(containerNames)

	for _, containerName := range containerNames {
		items = append(items, component.FlexLayoutItem{
			Width: component.WidthThird,
			View:  createCardComponent(fmt.Sprintf("Container %s", containerName), report.Report.ContainerChecks[containerName]),
		})
	}

	flexLayout.AddSections(items)

	return flexLayout
}

// Deprecated
func createSummary(report *starboard.ConfigAudit) (s *component.Summary) {
	s = component.NewSummary("")
	s.AddSection("Pod Template", createChecksTable(report.PodChecks))
	for ccn, cc := range report.ContainerChecks {
		s.AddSection(fmt.Sprintf("Container %s", ccn), createChecksTable(cc))
	}
	return
}

func createCardComponent(title string, checks []starboard.Check) (x *component.Card) {
	x = component.NewCard(component.TitleFromString(title))
	x.SetBody(createChecksTable(checks))
	return x
}

func createChecksTable(checks []starboard.Check) component.Component {
	table := component.NewTableWithRows(
		"", "There are no checks!",
		component.NewTableCols("Success", "ID", "Severity", "Category"),
		[]component.TableRow{})

	for _, c := range checks {
		checkID := c.ID

		tr := component.TableRow{
			"Success":  component.NewText(strconv.FormatBool(c.Success)),
			"ID":       component.NewText(checkID),
			"Severity": component.NewText(fmt.Sprintf("%s", c.Severity)),
			"Category": component.NewText(c.Category),
		}
		table.Add(tr)
	}

	return table
}

// TODO Implement summary
func createChecksSummary(report starboard.ConfigAudit) (c *component.Summary) {
	c = component.NewSummary("Summary")

	sections := []component.SummarySection{
		{Header: "error ", Content: component.NewText(strconv.Itoa(30))},
		{Header: "warning ", Content: component.NewText(strconv.Itoa(10))},
	}
	c.Add(sections...)
	return
}