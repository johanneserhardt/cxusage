# cx - Codex Usage Tracker

A gorgeous CLI tool for analyzing OpenAI Codex CLI usage data from local files, inspired by ccusage for Claude Code. Track your token usage, costs, and usage patterns with stunning responsive terminal output and real-time live monitoring.

## 📋 Releases

👉 Check releases page 📖 **[Releases](releases/README.md)**

## 🎨 Beautiful Output

```
┌──────────────────────────────────────┐
│ Codex CLI Token Usage Report - Daily │
└──────────────────────────────────────┘

┌──────────────┬──────────────────────┬────────────┬────────────┬──────────────┬──────────────┬──────────────┬────────────┐
│ Date         │ Models               │ Input      │ Output     │ Cache Create │ Cache Read   │ Total Tokens │ Cost (USD) │
├──────────────┼──────────────────────┼────────────┼────────────┼──────────────┼──────────────┼──────────────┼────────────┤
│ 2025-09-02   │ - gpt-4o, gpt-3.5... │ 2,000      │ 500        │ 0            │ 0            │ 2,500        │ $0.0003    │
│ TOTAL        │                      │ 2,000      │ 500        │ 0            │ 0            │ 2,500        │ $0.0003    │
└──────────────┴──────────────────────┴────────────┴────────────┴──────────────┴──────────────┴──────────────┴────────────┘
```

*Responsive tables with beautiful Unicode borders and theme-adaptive colors*

## ✨ Features

- 📊 **Responsive Tables** - Beautiful adaptive layouts that fit any terminal size
- 🔄 **Live Monitoring** - Real-time 5-hour billing block tracking with stunning dashboard
- 🎨 **Theme-Adaptive** - Gorgeous output that works on any terminal theme
- 💰 **Precise Cost Tracking** - Accurate calculation with 4-decimal precision for small amounts
- 🎯 **5-Hour Blocks** - Track usage in Claude-style billing periods with live projections
- 📤 **Multiple Output Formats** - Beautiful tables or JSON output
- 💾 **Local File Reading** - No API key needed - reads directly from ~/.codex/

## 🚀 Quick Start

### Prerequisites
1. **OpenAI Codex CLI** installed and configured
2. Used Codex CLI at least once (to generate usage data)
3. cx will automatically find and read from `~/.codex/` directory

### Basic Usage
```bash
# Validate your setup
cx validate

# See beautiful demo output
cx demo

# Daily usage report
cx daily

# 🔥 AMAZING: Live monitoring dashboard!
cx blocks --live
```

## 📦 Installation

### 🍺 Homebrew (macOS/Linux) - Recommended
The easiest way to install cx on macOS or Linux:

```bash
# Add tap and install
brew tap johanneserhardt/tap
brew install cxusage

# Then use immediately:
cx demo
cx blocks --live
```

**Benefits:**
- ✅ **Automatic platform detection** (Intel/ARM, macOS/Linux)
- ✅ **Secure installation** with SHA256 verification
- ✅ **Both commands available**: `cx` (short) and `cxusage` (full)
- ✅ **Professional tap integration** with discoverable formulas
- ✅ **Automatic updates** when new releases are published

### 🚀 Download Binary (Windows/Manual Install)
For Windows users or manual installation:

**[📥 Download from GitHub Releases](https://github.com/johanneserhardt/cxusage/releases/latest)**

Available for:
- **Linux** (AMD64, ARM64)
- **macOS** (Intel, Apple Silicon)
- **Windows** (AMD64)

After download:
```bash
# Make executable (Linux/macOS)
chmod +x cx-*
sudo mv cx-* /usr/local/bin/cx

# Windows: Add to PATH or use directly
```

### 🛠️ Go Install
```bash
go install github.com/johanneserhardt/cxusage/cmd/cxusage@latest
# Then use as: cx [command]
```

### 🔨 Build from Source
```bash
git clone https://github.com/johanneserhardt/cxusage
cd cxusage

# Build and install automatically
./scripts/install.sh

# Or just build
./scripts/build.sh
```

## ⚙️ Configuration

### Configuration File (Optional)

Create a config file at `~/.config/cxusage.yaml`:

```yaml
log_level: "warn"
local_logging: true
logs_dir: "logs"
codex_path: "/custom/path/to/codex"  # Optional custom Codex directory
```

## 📋 Commands

### Daily Reports
```bash
# Last 7 days (default)
cx daily

# Last 30 days
cx daily 30

# Specific date range
cx daily --start-date 2024-01-01 --end-date 2024-01-31

# JSON output
cx daily --output json
```

### Monthly Reports
```bash
# Last 3 months (default)
cx monthly

# Last 6 months
cx monthly 6

# JSON output
cx monthly --output json
```

### 🔥 Live Monitoring (Best Feature!)
```bash
# Live dashboard with real-time updates
cx blocks --live

# Live monitoring with token limit warnings
cx blocks --live --token-limit 50000

# Custom refresh rate (2 seconds)
cx blocks --live --refresh-interval 2
```

### 5-Hour Blocks
```bash
# Show recent billing blocks
cx blocks

# Show only active block
cx blocks --active

# Show blocks from last 7 days
cx blocks --recent --recent-days 7
```

### Utility Commands
```bash
# Validate Codex CLI setup
cx validate

# Show version info
cx version

# See demo of beautiful output
cx demo
```

## 🎨 Visual Features

cx provides absolutely stunning terminal output with:

- **Responsive Unicode tables** that adapt to any screen size
- **Theme-adaptive colors** that work perfectly on dark/light terminals
- **Real-time live dashboard** with gorgeous progress bars
- **Professional typography** with proper spacing and alignment
- **Smart number formatting** with thousand separators
- **Precise cost tracking** with 4-decimal precision for small amounts
- **Status indicators** and visual warnings

## 🤖 Supported Models

Cost calculation support for all OpenAI models:

- **GPT-4 Family**: gpt-4, gpt-4-turbo, gpt-4o, gpt-4o-mini
- **GPT-3.5 Family**: gpt-3.5-turbo, gpt-3.5-turbo-16k
- **Legacy Models**: text-davinci-003, code-davinci-002
- **Embedding Models**: text-embedding-3-small, text-embedding-3-large
- **Fine-tuned Models**: Automatic detection and pricing

## 📊 Live Monitoring Dashboard

The `cx blocks --live` command provides a stunning real-time dashboard featuring:

- **🟢 SESSION** - Progress through current 5-hour block with visual timeline
- **🔥 USAGE** - Current token usage with live burn rate tracking
- **📈 PROJECTION** - Projected usage with limit warnings ("WILL EXCEED LIMIT")
- **⚙️ MODELS** - Active models being used in current session
- **Real-time updates** every second with smooth animations
- **Visual progress bars** with color coding (green → yellow → red)
- **Smart alerts** and professional status indicators

## 🔧 Global Flags

- `--output, -o` - Output format: table (default) or json
- `--log-level` - Log level: debug, info, warn, error

## 🛠️ Troubleshooting

### Homebrew Installation Issues

**Problem:** Getting "example.invalid" download errors when installing via Homebrew
```
Error: cxusage: Failed to download resource "cxusage (0.0.0)"
Download failed: https://example.invalid/cxusage-darwin-arm64
```

**Solution:** Homebrew is using a cached version of the old formula. Refresh the tap:
```bash
# Remove and re-add the tap to clear cache
brew untap johanneserhardt/tap
brew tap johanneserhardt/tap

# Or force update the tap
brew tap johanneserhardt/tap --force
brew update

# Then install
brew install cxusage
```

### No Codex Usage Data Found

**Problem:** cx shows "No Codex CLI usage data found"

**Solutions:**
```bash
# Check if Codex CLI is properly set up
cx validate

# Verify Codex directory exists
ls -la ~/.codex/

# Use Codex CLI first to generate data, then run cx
# After using Codex CLI, try again:
cx daily
cx blocks --live
```

### Live Dashboard Shows "Waiting for Activity"

**Problem:** `cx blocks --live` shows waiting message despite recent Codex usage

**Reasons:**
- **No current active block** - Usage was in a previous 5-hour window
- **Codex usage too old** - Activity more than 1 hour ago
- **Different timezone** - Block times might seem off

**Solutions:**
```bash
# Check all blocks to see your usage history
cx blocks

# Use Codex CLI now to create activity in current block
# Then try live monitoring again:
cx blocks --live
```

### Terminal Output Issues

**Problem:** Tables look broken or colors don't display properly

**Solutions:**
```bash
# Check terminal compatibility
echo $TERM

# Try without colors if needed
NO_COLOR=1 cx daily

# Ensure terminal supports Unicode
cx demo  # Should show beautiful tables
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -am 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by [ccusage](https://github.com/ryoppippi/ccusage)
---

**Made with ❤️ for the OpenAI community**
