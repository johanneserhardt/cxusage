# cxusage - Codex Usage Tracker

A beautiful CLI tool for analyzing OpenAI Codex CLI usage data from local files, inspired by ccusage for Claude Code. Track your token usage, costs, and usage patterns with gorgeous terminal output and real-time monitoring.

## ✨ Features

- 📊 **Daily Reports** - Beautiful daily usage tables with colors and formatting
- 📅 **Monthly Analysis** - Comprehensive monthly breakdowns with summaries
- 🔄 **Live Monitoring** - Real-time 5-hour billing block tracking with gorgeous dashboard
- 💰 **Cost Calculation** - Accurate cost calculation using current OpenAI pricing
- 🎯 **5-Hour Blocks** - Track usage in Claude-style 5-hour billing periods
- 📤 **Multiple Output Formats** - Beautiful tables or JSON output
- 🎨 **Beautiful Output** - Professional formatting matching ccusage's visual style
- 💾 **Local File Reading** - No API key needed - reads directly from ~/.codex/

## 🚀 Quick Start

**cxusage** reads usage data directly from your local Codex CLI installation.

### Prerequisites
1. **OpenAI Codex CLI** installed and configured
2. Used Codex CLI at least once (to generate usage data)
3. cxusage will automatically find and read from `~/.codex/` directory

### Basic Usage
```bash
# Validate your setup
cxusage validate

# See beautiful demo output
cxusage demo

# Daily usage report
cxusage daily

# 🔥 AMAZING: Live monitoring dashboard!
cxusage blocks --live
```

## 📦 Installation

```bash
go install github.com/johanneserhardt/cxusage/cmd/cxusage@latest
```

Or build from source:

```bash
git clone https://github.com/johanneserhardt/cxusage
cd cxusage
./scripts/build.sh
```

## ⚙️ Configuration

### Configuration File (Optional)

Create a config file at `~/.config/cxusage.yaml`:

```yaml
log_level: "info"
local_logging: true
logs_dir: "logs"
codex_path: "/custom/path/to/codex"  # Optional custom Codex directory
```

## 📋 Commands

### Daily Reports
```bash
# Last 7 days (default)
cxusage daily

# Last 30 days
cxusage daily 30

# Specific date range
cxusage daily --start-date 2024-01-01 --end-date 2024-01-31

# JSON output
cxusage daily --output json
```

### Monthly Reports
```bash
# Last 3 months (default)
cxusage monthly

# Last 6 months
cxusage monthly 6

# JSON output
cxusage monthly --output json
```

### 🔥 Live Monitoring (Best Feature!)
```bash
# Live dashboard with real-time updates
cxusage blocks --live

# Live monitoring with token limit warnings
cxusage blocks --live --token-limit 50000

# Custom refresh rate (2 seconds)
cxusage blocks --live --refresh-interval 2
```

### 5-Hour Blocks
```bash
# Show recent billing blocks
cxusage blocks

# Show only active block
cxusage blocks --active

# Show blocks from last 7 days
cxusage blocks --recent --recent-days 7
```

### Utility Commands
```bash
# Validate Codex CLI setup
cxusage validate

# Show version info
cxusage version

# See demo of beautiful output
cxusage demo
```

## 🎨 Beautiful Output

cxusage provides gorgeous terminal output with:

- **Color-coded tables** with professional formatting
- **Progress bars** with visual indicators
- **Real-time dashboard** with live updates
- **Smart number formatting** with thousand separators
- **Status indicators** and warnings
- **Professional borders** and spacing

## 🤖 Supported Models

Cost calculation support for all OpenAI models:

- **GPT-4 Family**: gpt-4, gpt-4-turbo, gpt-4o, gpt-4o-mini
- **GPT-3.5 Family**: gpt-3.5-turbo, gpt-3.5-turbo-16k
- **Legacy Models**: text-davinci-003, code-davinci-002
- **Embedding Models**: text-embedding-3-small, text-embedding-3-large
- **Fine-tuned Models**: Automatic detection and pricing

## 📊 Live Monitoring Dashboard

The `cxusage blocks --live` command provides a stunning real-time dashboard featuring:

- **🟢 SESSION** - Progress through current 5-hour block
- **🔥 USAGE** - Current token usage with burn rate
- **📈 PROJECTION** - Projected usage with limit warnings
- **⚙️ MODELS** - Active models being used
- **Real-time updates** every second
- **Visual progress bars** with color coding
- **Smart alerts** like "WILL EXCEED LIMIT"

## 🔧 Global Flags

- `--output, -o` - Output format: table (default) or json
- `--log-level` - Log level: debug, info, warn, error
- `--offline` - Use local logs only (legacy flag, always local now)

## 💡 Examples

```bash
# Beautiful daily report with colors
cxusage daily

# Live monitoring with gorgeous dashboard
cxusage blocks --live

# Monthly analysis with formatting
cxusage monthly

# Check your setup
cxusage validate

# See what cxusage can do
cxusage demo
```

## 🔮 How It Works

1. **Codex CLI** stores usage data locally in `~/.codex/`
2. **cxusage** reads these local files (no API key needed!)
3. Aggregates data into reports and billing blocks
4. Displays with beautiful formatting and colors
5. Live mode monitors files for real-time updates

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -am 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by [ccusage](https://github.com/ryoppippi/ccusage) for Claude Code by @ryoppippi
- Built with [Cobra](https://github.com/spf13/cobra) CLI framework
- Uses [tablewriter](https://github.com/olekukonko/tablewriter) for beautiful tables
- Colors powered by [fatih/color](https://github.com/fatih/color)

---

**Made with ❤️ for the Codex CLI community**
