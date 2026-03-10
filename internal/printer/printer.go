package printer

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"evm-profiler/internal/model"

	"github.com/fatih/color"
)

// PrintProfile outputs a visually appealing terminal report for the AddressProfile.
func PrintProfile(p *model.AddressProfile) {
	fmt.Println()
	titleColor := color.New(color.FgCyan, color.Bold)
	titleColor.Printf("🔍 EVM Address Profile Report: %s\n", p.Address)
	fmt.Println(strings.Repeat("=", 60))

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	// Address Type
	addrType := "Externally Owned Account (EOA)"
	if p.IsContract {
		addrType = color.MagentaString("Smart Contract")
	} else {
		addrType = color.BlueString(addrType)
	}
	fmt.Fprintf(w, "%s\t%s\n", color.WhiteString("Type:"), addrType)

	// Balance
	balanceStr := fmt.Sprintf("%.4f ETH", p.BalanceEther)
	if p.BalanceEther != nil {
		if bigFloatCmp, _ := p.BalanceEther.Float64(); bigFloatCmp > 100 {
			balanceStr = color.YellowString(balanceStr)
		} else {
			balanceStr = color.GreenString(balanceStr)
		}
	}
	fmt.Fprintf(w, "%s\t%s\n", color.WhiteString("Balance:"), balanceStr)

	// Activity
	activityStr := p.ActivityLabel
	if p.TransactionCount > 0 {
		activityStr += fmt.Sprintf(" (~%d txs)", p.TransactionCount)
	}
	fmt.Fprintf(w, "%s\t%s\n", color.WhiteString("Activity:"), activityStr)

	// Initial Funder
	funder := p.Funder
	if funder == "" {
		funder = color.HiBlackString("Unknown / Not Found")
	}
	fmt.Fprintf(w, "%s\t%s\n", color.WhiteString("Funder:"), funder)

	// Security Risks
	riskStr := color.GreenString("✅ No Known Risks")
	if len(p.RiskLabels) > 0 {
		riskStr = color.RedString("🚨 Risks Detected: %s", strings.Join(p.RiskLabels, ", "))
	}
	fmt.Fprintf(w, "%s\t%s\n", color.WhiteString("Security:"), riskStr)

	w.Flush()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()
}
