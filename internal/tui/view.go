package tui

import (
	"fmt"
	"strings"
)

func (m Model) View() string {
	if m.Loading {
		return separatorStyle.Render("═══════════════════════════════════════") + "\n" +
			headerStyle.Render("  Loading proxies...")
	}

	if m.Err != nil {
		return separatorStyle.Render("═══════════════════════════════════════") + "\n" +
			headerStyle.Render("  Error") + "\n" +
			fmt.Sprintf("  %v\n", m.Err) +
			helpStyle.Render("  Press [r] retry, [q] quit")
	}

	if len(m.Groups) == 0 {
		return separatorStyle.Render("═══════════════════════════════════════") + "\n" +
			headerStyle.Render("  No proxy groups found") + "\n" +
			helpStyle.Render("  Press [r] refresh, [q] quit")
	}

	var s string

	for i, group := range m.Groups {
		proxy, ok := m.Proxies[group]
		if !ok {
			continue
		}

		var groupLabel string
		if i == m.CurrentIdx {
			groupLabel = selectedGroupStyle.Render("● " + group)
		} else {
			groupLabel = normalStyle.Render("○ " + group)
		}
		s += groupLabel + "\n"

		if i == m.CurrentIdx {
			visibleProxies := proxy.All
			if len(proxy.All) > maxVisibleProxies {
				startIdx := m.ViewportOffset
				if startIdx < 0 {
					startIdx = 0
				}
				endIdx := startIdx + maxVisibleProxies
				if endIdx > len(proxy.All) {
					endIdx = len(proxy.All)
				}
				visibleProxies = proxy.All[startIdx:endIdx]
			}

			for j, p := range visibleProxies {
				actualIdx := j + m.ViewportOffset
				var line string
				if actualIdx == m.Cursor && p == proxy.Now {
					line = cursorStyle.Render(">● ") + activeProxyStyle.Render(p)
				} else if actualIdx == m.Cursor {
					line = cursorStyle.Render(">  ") + p
				} else if p == proxy.Now {
					line = " ● " + activeProxyStyle.Render(p)
				} else {
					line = "   " + normalStyle.Render(p)
				}
				s += line + "\n"
			}

			if len(proxy.All) > maxVisibleProxies {
				s += helpStyle.Render("  " + strings.Repeat("-", 20) + fmt.Sprintf(" %d/%d ", len(proxy.All), len(proxy.All))) + "\n"
			}
		}
	}

	s += separatorStyle.Render("═══════════════════════════════════════") + "\n"
	s += helpStyle.Render(" [←h]Prev [→l]Next  [↑k]↑ [↓j]↓  [Ent]Select  [r]Reload  [q]Quit")

	return s
}
