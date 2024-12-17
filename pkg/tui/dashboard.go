package tui

import (
	"fmt"
	"github.com/echogy-io/echogy/pkg/stat"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	q "github.com/echogy-io/echogy/pkg/queue"
)

// Constants for layout and styling
const (
	defaultTableHeight = 10
	minHeaderWidth     = 80
)

// Style definitions
var (
	dashStyle = lipgloss.NewStyle().Padding(1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			BorderBottom(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: "#94A3B8", Dark: "#4A5568"})

	urlStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.AdaptiveColor{Light: "#2563EB", Dark: "#60A5FA"}).
			Width(45).
			Align(lipgloss.Left)

	statsStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.AdaptiveColor{Light: "#475569", Dark: "#94A3B8"}).
			Width(12).
			Align(lipgloss.Right)

	colMethodHeaderStyle = lipgloss.NewStyle().Align(lipgloss.Center).Bold(true)

	colNoStyle = lipgloss.NewStyle().Align(lipgloss.Left)

	colStatusStyle = lipgloss.NewStyle().Bold(true)

	colPathHeaderStyle = lipgloss.NewStyle().Bold(false)

	colPathStyle = lipgloss.NewStyle().Inherit(colPathHeaderStyle)

	colUseTimeHeaderStyle = lipgloss.NewStyle().Align(lipgloss.Right)

	colUseTimeStyle = lipgloss.NewStyle().
			Inherit(colUseTimeHeaderStyle).
			Bold(false).
			Foreground(lipgloss.AdaptiveColor{Light: "#475569", Dark: "#94A3B8"})

	qrStyle = lipgloss.NewStyle().
		Align(lipgloss.Top).
		Foreground(lipgloss.AdaptiveColor{Light: "#0F172A", Dark: "#F8FAFC"}).
		Padding(1)

	githubStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.AdaptiveColor{Light: "#2563EB", Dark: "#60A5FA"})

	descStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#334155", Dark: "#E2E8F0"}).
			PaddingTop(1)

	tableStyle = table.Styles{
		Header: lipgloss.NewStyle().Bold(true),
		Selected: lipgloss.NewStyle().
			Background(lipgloss.AdaptiveColor{Light: "#EDF2F7", Dark: "#1E293B"}).
			Foreground(lipgloss.AdaptiveColor{Light: "#2563EB", Dark: "#60A5FA"}).
			Bold(true),
		Cell: lipgloss.NewStyle(),
	}
)

// TableColumn defines the structure for table column configuration
type TableColumn struct {
	Title  string
	Weight float64
}

// RequestTable wraps the table model with additional configuration
type RequestTable struct {
	*table.Model
	columns []TableColumn
}

func (r *RequestTable) Init() tea.Cmd {
	return nil
}

func (r *RequestTable) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var m table.Model
	m, cmd = r.Model.Update(msg)
	r.Model = &m
	return r, cmd
}

// newRequestTable creates a new request table with predefined columns
func newRequestTable(width int) *RequestTable {
	columns := []TableColumn{
		{Title: colNoStyle.Render("#"), Weight: 0.05},                 // 5% of available width
		{Title: colMethodHeaderStyle.Render("Method"), Weight: 0.1},   // 10% of available width
		{Title: colStatusStyle.Render("Status"), Weight: 0.1},         // 10% of available width
		{Title: colPathHeaderStyle.Render("URI"), Weight: 0.65},       // 65% of available width
		{Title: colUseTimeHeaderStyle.Render("UseTime"), Weight: 0.1}, // 10% of available width
	}

	cols := make([]table.Column, len(columns))
	for i, col := range columns {
		cols[i] = table.Column{
			Title: col.Title,
			Width: int(col.Weight * float64(width)),
		}
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
		table.WithHeight(defaultTableHeight),
	)

	t.SetStyles(tableStyle)

	return &RequestTable{
		Model:   &t,
		columns: columns,
	}
}

// Dashboard represents the main TUI dashboard
type Dashboard struct {
	width      int
	height     int
	tunnelInfo TunnelInfo
	table      *RequestTable
	stat       *stat.Stat
	requests   *q.FixedQueue
}

// TunnelInfo holds information about the tunnel connection
type TunnelInfo struct {
	URL       string
	ExpiresIn time.Duration
	BytesRecv int64
	BytesSent int64
	ReqCount  int
	ResCount  int
}

// newDashboard creates a new dashboard instance
func newDashboard(queue *q.FixedQueue, stat *stat.Stat, tunnelAddr string, width, height int) *Dashboard {
	return &Dashboard{
		tunnelInfo: TunnelInfo{
			URL:       tunnelAddr,
			ExpiresIn: 10 * time.Minute,
		},
		width:    width,
		height:   height,
		table:    newRequestTable(width),
		stat:     stat,
		requests: queue,
	}
}

func (d *Dashboard) availableWidth() int {
	return d.width - 4
}

// updateTableWidth adjusts table column widths based on terminal width
func (d *Dashboard) updateTableWidth() {
	if d.width <= 0 {
		return
	}

	// Calculate available width (accounting for margins)
	aw := d.availableWidth()

	// Apply new column widths
	tableColumns := make([]table.Column, len(d.table.columns))
	for i, col := range d.table.columns {
		tableColumns[i] = table.Column{
			Title: col.Title, // Keep original styled title
			Width: int(float64(aw) * col.Weight),
		}
	}

	d.table.SetColumns(tableColumns)
}

// Init implements tea.Model
func (d *Dashboard) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (d *Dashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return d, tea.Quit
		}
	case tea.WindowSizeMsg:
		d.width = msg.Width
		d.height = msg.Height
		d.table.SetHeight(msg.Height - 8)
		d.updateTableWidth()
	}

	d.preUpdate()

	// Always update the table model
	var m table.Model
	m, cmd = d.table.Model.Update(msg)
	d.table.Model = &m
	return d, cmd
}

// renderHeader renders the header section with URLs and stats
func (d *Dashboard) renderHeader() string {
	// URLs section
	urls := []string{
		fmt.Sprintf("HTTP:  %s", d.tunnelInfo.URL),
		fmt.Sprintf("HTTPS: %s", d.tunnelInfo.URL),
	}

	// Stats section
	stats := []string{
		fmt.Sprintf("↓ %s", humanBytes(d.tunnelInfo.BytesRecv)),
		fmt.Sprintf("↑ %s", humanBytes(d.tunnelInfo.BytesSent)),
	}

	counts := []string{
		fmt.Sprintf("Req: %d", d.tunnelInfo.ReqCount),
		fmt.Sprintf("Res: %d", d.tunnelInfo.ResCount),
	}

	if d.width < minHeaderWidth {
		// 小屏幕：垂直布局
		leftURLS := lipgloss.JoinVertical(
			lipgloss.Left,
			urlStyle.Render(urls[0]),
			urlStyle.Render(urls[1]),
		)

		statsInfo := lipgloss.JoinHorizontal(
			lipgloss.Left,
			statsStyle.Render(stats[0]),
			statsStyle.Render(counts[0]),
		)

		countInfo := lipgloss.JoinHorizontal(
			lipgloss.Left,
			statsStyle.Render(stats[1]),
			statsStyle.Render(counts[1]),
		)

		return lipgloss.NewStyle().Inherit(headerStyle).Width(d.width - 2).Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				leftURLS,
				"\n",
				statsInfo,
				countInfo,
			),
		)
	}

	// 大屏幕：水平布局
	leftURLS := lipgloss.JoinVertical(
		lipgloss.Left,
		urlStyle.Render(urls[0]),
		urlStyle.Render(urls[1]),
	)

	rightStats := lipgloss.JoinVertical(
		lipgloss.Right,
		lipgloss.JoinHorizontal(
			lipgloss.Right,
			statsStyle.Render(stats[0]),
			statsStyle.Render(counts[0]),
		),
		lipgloss.JoinHorizontal(
			lipgloss.Right,
			statsStyle.Render(stats[1]),
			statsStyle.Render(counts[1]),
		),
	)

	return lipgloss.NewStyle().Inherit(headerStyle).Width(d.width - 2).Render(
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			leftURLS,
			rightStats,
		),
	)
}

// renderProjectInfo renders the project information section
func (d *Dashboard) renderProjectInfo() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		githubStyle.Render("Service: https://echogy.io"),
		descStyle.Render("Echogy - A lightweight and efficient SSH reverse proxy tool\n\n"+
			"Features:\n"+
			"• Easy to use HTTP/HTTPS tunnel\n"+
			"• Real-time traffic monitoring\n"+
			"• Beautiful TUI interface\n"+
			"• Cross-platform support\n\n"+
			"⭠ Scan QR code to access your tunnel"),
	)
}

// View implements tea.Model
func (d *Dashboard) View() string {
	head := d.renderHeader()

	var content string
	if d.requests.Len() == 0 {
		// Show QR code and project info when table is empty
		qrCode := generateQRCode(d.tunnelInfo.URL)
		if d.width > minHeaderWidth {
			content = lipgloss.JoinHorizontal(
				lipgloss.Center,
				qrStyle.Render(qrCode),
				lipgloss.NewStyle().PaddingRight(4).Render(d.renderProjectInfo()),
			)
		} else {
			content = lipgloss.JoinHorizontal(
				lipgloss.Center,
				qrStyle.Render(qrCode),
			)
		}
	} else {
		content = d.table.View()
	}

	return dashStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			head,
			content),
	)
}

func renderStatusCode(status int) string {
	s := strconv.Itoa(status)
	style := lipgloss.NewStyle().Inherit(colStatusStyle)
	switch {
	case status >= 500:
		return style.Foreground(lipgloss.AdaptiveColor{Light: "#E53E3E", Dark: "#FC8181"}).Render(s) // Bright Red
	case status >= 400:
		return style.Foreground(lipgloss.AdaptiveColor{Light: "#DD6B20", Dark: "#FBD38D"}).Render(s) // Bright Orange
	case status >= 300:
		return style.Foreground(lipgloss.AdaptiveColor{Light: "#3182CE", Dark: "#90CDF4"}).Render(s) // Bright Blue
	case status >= 200:
		return style.Foreground(lipgloss.AdaptiveColor{Light: "#38A169", Dark: "#9AE6B4"}).Render(s) // Bright Green
	default:
		return style.Foreground(lipgloss.AdaptiveColor{Light: "#718096", Dark: "#A0AEC0"}).Render(s) // Cool Gray
	}
}

func renderMethod(method string) string {
	style := lipgloss.NewStyle().Inherit(colMethodHeaderStyle)
	switch method {
	case "GET":
		return style.Foreground(lipgloss.AdaptiveColor{Light: "#3182CE", Dark: "#90CDF4"}).Render(method) // Bright Blue
	case "POST":
		return style.Foreground(lipgloss.AdaptiveColor{Light: "#38A169", Dark: "#9AE6B4"}).Render(method) // Bright Green
	case "PUT":
		return style.Foreground(lipgloss.AdaptiveColor{Light: "#805AD5", Dark: "#B794F4"}).Render(method) // Bright Purple
	case "DELETE":
		return style.Foreground(lipgloss.AdaptiveColor{Light: "#E53E3E", Dark: "#FC8181"}).Render(method) // Bright Red
	case "PATCH":
		return style.Foreground(lipgloss.AdaptiveColor{Light: "#DD6B20", Dark: "#FBD38D"}).Render(method) // Bright Orange
	default:
		return style.Foreground(lipgloss.AdaptiveColor{Light: "#718096", Dark: "#A0AEC0"}).Render(method) // Cool Gray
	}
}

func (d *Dashboard) preUpdate() {
	d.tunnelInfo.BytesRecv = d.stat.Receive
	d.tunnelInfo.BytesSent = d.stat.Send

	d.tunnelInfo.ReqCount = d.stat.Request
	d.tunnelInfo.ResCount = d.stat.Response

	l := min(d.requests.Len(), 32)

	items := d.requests.ReversedItems()

	rows := make([]table.Row, l)
	// Update table rows
	for i := 0; i < l; i++ {
		r := items[i].(*stat.RequestEntity)
		path := colPathStyle.Render(r.RequestURI)
		t := colUseTimeStyle.Render(humanMillis(r.UseTime))

		rows[i] = table.Row{
			colNoStyle.Render(strconv.Itoa(l - i)),
			renderMethod(r.Method),
			renderStatusCode(r.StatusCode),
			path,
			t,
		}
	}
	d.table.SetRows(rows)
}

// UpdateStats updates the tunnel information statistics
func (d *Dashboard) UpdateStats(info TunnelInfo) {
	d.tunnelInfo = info
}
