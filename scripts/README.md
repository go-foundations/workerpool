# 🚀 Auto-Create GitHub Issues Scripts

This directory contains scripts to automatically create GitHub issues from the documentation files for the Resource-Aware Scheduling implementation.

## 📋 Available Scripts

### **Python Script (Recommended)**
- **File**: `create_issues.py`
- **Dependencies**: PyGithub
- **Features**: Creates milestones and issues with proper formatting

### **Shell Script**
- **File**: `create_issues.sh`
- **Dependencies**: GitHub CLI (`gh`)
- **Features**: Creates milestones and issues with proper formatting

**Note**: The `create_all_issues.sh` script was used to create the initial set of issues and has been removed since it's no longer needed.

## 🛠️ Setup Instructions

### **Option 1: Python Script (Recommended)**

1. **Install Python dependencies:**
   ```bash
   pip install -r requirements.txt
   ```

2. **Set GitHub token:**
   ```bash
   export GITHUB_TOKEN=your_github_personal_access_token
   ```
   
   **To get a token:**
   - Go to GitHub Settings → Developer settings → Personal access tokens
   - Generate new token with `repo` scope
   - Copy the token and set it as environment variable

3. **Run the script:**
   ```bash
   cd scripts
   python create_issues.py
   ```

### **Option 2: GitHub CLI Script**

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

3. **Run the script:**
   ```bash
   cd scripts
   chmod +x create_issues.sh
   ./create_issues.sh
   ```

## 📊 What Gets Created

### **Milestones (7 total)**
- Phase 1: Foundation
- Phase 2: Basic Resource-Aware Scheduling
- Phase 3: Learning and Profiling
- Phase 4: Advanced Scheduling Algorithms
- Phase 5: Observability and Monitoring
- Phase 6: Cloud Integration
- Phase 7: Testing and Documentation

### **Issues (6 created, 23 remaining)**
- **Phase 1**: 3 issues ✅ (Foundation)
- **Phase 2**: 3 issues ✅ (Basic Scheduling)
- **Phase 3**: 3 issues (Learning) - Ready for creation
- **Phase 4**: 3 issues (Advanced Algorithms) - Ready for creation
- **Phase 5**: 3 issues (Observability) - Ready for creation
- **Phase 6**: 3 issues (Cloud Integration) - Ready for creation
- **Phase 7**: 2 issues (Testing & Docs) - Ready for creation

## 🏷️ Issue Labels

Each issue gets automatically labeled with:
- **Phase labels**: `phase1`, `phase2`, etc.
- **Feature labels**: `foundation`, `scheduling`, `resources`, etc.
- **Type labels**: `enhancement`
- **Priority labels**: Based on phase importance

## 🔧 Customization

### **Modify Repository**
Edit the `REPO_NAME` variable in the script:
```python
REPO_NAME = "your-username/your-repo"
```

### **Add More Issues**
The Python script can be extended to create the remaining 23 issues for Phases 3-7. Currently 6 issues are created (Phases 1-2).

### **Modify Issue Content**
Edit the issue body text in the script to customize descriptions, requirements, and acceptance criteria.

## 🚨 Troubleshooting

### **Common Issues**

1. **"GITHUB_TOKEN not set"**
   - Ensure you've set the environment variable
   - Check token has `repo` scope

2. **"Repository not found"**
   - Verify repository name is correct
   - Ensure token has access to the repository

3. **"Milestone already exists"**
   - Script will reuse existing milestones
   - No duplicate milestones will be created

4. **"Rate limit exceeded"**
   - GitHub has rate limits for API calls
   - Wait a few minutes and try again

### **Debug Mode**

Add debug logging to the Python script:
```python
import logging
logging.basicConfig(level=logging.DEBUG)
```

## 📈 Current Status & Next Steps

### **✅ What's Already Done:**
- **6 issues created** for Phases 1-2 (Foundation & Basic Scheduling)
- **7 milestones created** for all phases
- **12 custom labels** for proper categorization
- **Repository ready** for community development

### **🔄 Next Steps:**

1. **Review created issues** on GitHub
2. **Assign issues** to team members
3. **Set priorities** and due dates
4. **Begin implementation** starting with Phase 1
5. **Create remaining issues** for Phases 3-7 when ready
6. **Update progress** as issues are completed

## 🔗 Useful Links

- [GitHub Issues API](https://docs.github.com/en/rest/issues)
- [PyGithub Documentation](https://pygithub.readthedocs.io/)
- [GitHub CLI Documentation](https://cli.github.com/)
- [Personal Access Tokens](https://github.com/settings/tokens)

---

**Note**: These scripts create a comprehensive set of issues for the Resource-Aware Scheduling implementation. Each issue includes detailed requirements, acceptance criteria, and implementation notes to guide development.
