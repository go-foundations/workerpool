# 🚀 Auto-Create GitHub Issues from Markdown

This directory contains the essential script to automatically create GitHub issues from markdown documentation files.

## 📋 Available Script

### **Main Script**
- **File**: `create_github_issues.py`
- **Dependencies**: GitHub CLI (`gh`)
- **Features**: 
  - Reads `.github/ISSUES_PHASE*.md` files
  - Parses issue content, labels, and milestones
  - Creates GitHub issues automatically
  - Skips already existing issues
  - Handles milestone mapping automatically

## 🛠️ Setup Instructions

### **Prerequisites**

1. **Install GitHub CLI:**
   ```bash
   # macOS
   brew install gh
   
   # Linux
   sudo apt install gh
   
   # Windows
   winget install GitHub.cli
   ```

2. **Authenticate:**
   ```bash
   gh auth login
   ```

3. **Ensure you have the required labels and milestones** (the script will create them if missing)

### **Usage**

```bash
# Run the script to create issues from markdown files
python scripts/create_github_issues.py
```

## 📊 What Gets Created

### **From Markdown Files:**
- **Issue titles** and descriptions
- **Labels** for categorization
- **Milestone assignments**
- **Full issue content** with requirements, acceptance criteria, and implementation notes

### **Current Status:**
- **20 issues created** across all 7 phases
- **7 milestones** with proper descriptions
- **40+ custom labels** for comprehensive categorization

## 🔧 How It Works

1. **Scans `.github/` directory** for `ISSUES_PHASE*.md` files
2. **Parses each file** to extract issue information
3. **Maps milestone names** to actual GitHub milestone names
4. **Creates issues** using GitHub CLI
5. **Skips existing issues** to avoid duplicates

## 📁 File Structure

```
.github/
├── ISSUES_PHASE1.md    # Phase 1 issues (Foundation)
├── ISSUES_PHASE2.md    # Phase 2 issues (Basic Scheduling)
├── ISSUES_PHASE3.md    # Phase 3 issues (Learning & Profiling)
├── ISSUES_PHASE4.md    # Phase 4 issues (Advanced Algorithms)
├── ISSUES_PHASE5.md    # Phase 5 issues (Observability)
├── ISSUES_PHASE6.md    # Phase 6 issues (Cloud Integration)
└── ISSUES_PHASE7.md    # Phase 7 issues (Testing & Docs)

scripts/
├── create_github_issues.py    # Main script
├── requirements.txt            # Python dependencies
└── README.md                  # This file
```

## 🏷️ Issue Labels

The script automatically creates and applies labels including:
- **Phase labels**: `phase1`, `phase2`, `phase3`, etc.
- **Feature labels**: `foundation`, `scheduling`, `resources`, `learning`, etc.
- **Type labels**: `enhancement`, `bug`, `documentation`
- **Technical labels**: `algorithms`, `observability`, `cloud`, etc.

## 🚨 Troubleshooting

### **Common Issues**

1. **"GitHub CLI not found"**
   - Install GitHub CLI first
   - Ensure `gh` is in your PATH

2. **"Not authenticated"**
   - Run `gh auth login` to authenticate

3. **"Label not found"**
   - The script will create missing labels automatically
   - Check if you have permission to create labels

4. **"Milestone not found"**
   - Ensure milestones exist in the repository
   - Check milestone names match exactly

### **Debug Mode**

Add debug logging to the script:
```python
import logging
logging.basicConfig(level=logging.DEBUG)
```

## 📈 Next Steps

After running the script:

1. **Review created issues** on GitHub
2. **Assign issues** to team members
3. **Set priorities** and due dates
4. **Begin implementation** starting with Phase 1
5. **Update markdown files** and re-run script for new issues

## 🔗 Useful Links

- [GitHub CLI Documentation](https://cli.github.com/)
- [GitHub Issues API](https://docs.github.com/en/rest/issues)
- [Repository Issues](https://github.com/go-foundations/workerpool/issues)
- [Repository Milestones](https://github.com/go-foundations/workerpool/milestones)

---

**Note**: This script provides a clean, maintainable way to convert markdown documentation into actionable GitHub issues. All issue content is stored in markdown files, making it easy to update and version control.
