#!/usr/bin/env python3
"""
Create GitHub Issues from markdown files using GitHub CLI
Reads .github/ISSUES_PHASE*.md files and creates corresponding GitHub issues
Skips already existing issues for Phases 1-2
"""

import os
import sys
import re
import glob
import subprocess
from pathlib import Path

def parse_issue_file(file_path):
    """Parse a markdown issue file and extract issue information"""
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Split into individual issues
    issues = []
    issue_sections = content.split('---')
    
    for section in issue_sections:
        if not section.strip():
            continue
            
        # Extract issue number and title
        title_match = re.search(r'## Issue #(\d+): (.+)', section)
        if not title_match:
            continue
            
        issue_num = int(title_match.group(1))
        title = title_match.group(2).strip()
        
        # Extract metadata
        labels_match = re.search(r'\*\*Labels\*\*: (.+)', section)
        labels = [label.strip() for label in labels_match.group(1).split(',')] if labels_match else []
        
        milestone_match = re.search(r'\*\*Milestone\*\*: (.+)', section)
        milestone = milestone_match.group(1).strip() if milestone_match else None
        
        # Map milestone names to actual milestone names
        if milestone:
            if "Phase 1" in milestone:
                milestone = "Phase 1"
            elif "Phase 2" in milestone:
                milestone = "Phase 2"
            elif "Phase 3" in milestone:
                milestone = "Phase 3"
            elif "Phase 4" in milestone:
                milestone = "Phase 4"
            elif "Phase 5" in milestone:
                milestone = "Phase 5"
            elif "Phase 6" in milestone:
                milestone = "Phase 6"
            elif "Phase 7" in milestone:
                milestone = "Phase 7"
        
        # Extract body content (everything after the metadata section)
        body_start = section.find('### Description')
        if body_start == -1:
            continue
            
        body = section[body_start:].strip()
        
        issues.append({
            'number': issue_num,
            'title': title,
            'body': body,
            'labels': labels,
            'milestone': milestone
        })
    
    return issues

def create_issue(issue_data):
    """Create a GitHub issue using GitHub CLI"""
    try:
        # Create temporary file for issue body
        body_file = f"/tmp/issue_{issue_data['number']}.md"
        with open(body_file, 'w', encoding='utf-8') as f:
            f.write(issue_data['body'])
        
        # Build gh command
        cmd = [
            'gh', 'issue', 'create',
            '--repo', 'go-foundations/workerpool',
            '--title', issue_data['title'],
            '--body-file', body_file
        ]
        
        # Add labels
        for label in issue_data['labels']:
            cmd.extend(['--label', label])
        
        # Add milestone
        if issue_data['milestone']:
            cmd.extend(['--milestone', issue_data['milestone']])
        
        # Execute command
        result = subprocess.run(cmd, capture_output=True, text=True)
        
        if result.returncode == 0:
            print(f"✅ Created Issue #{issue_data['number']}: {issue_data['title']}")
            # Clean up temp file
            os.remove(body_file)
            return True
        else:
            print(f"❌ Failed to create Issue #{issue_data['number']}: {result.stderr}")
            # Clean up temp file
            os.remove(body_file)
            return False
            
    except Exception as e:
        print(f"❌ Error creating Issue #{issue_data['number']}: {e}")
        return False

def main():
    print("🚀 Creating GitHub Issues from markdown files using GitHub CLI")
    print("Repository: go-foundations/workerpool")
    print("")
    
    # Check if gh is available
    try:
        subprocess.run(['gh', '--version'], capture_output=True, check=True)
    except (subprocess.CalledProcessError, FileNotFoundError):
        print("❌ GitHub CLI (gh) not found. Please install it first.")
        print("   macOS: brew install gh")
        print("   Linux: sudo apt install gh")
        print("   Windows: winget install GitHub.cli")
        sys.exit(1)
    
    # Check if authenticated
    try:
        result = subprocess.run(['gh', 'auth', 'status'], capture_output=True, text=True)
        if result.returncode != 0:
            print("❌ Not authenticated with GitHub CLI. Please run 'gh auth login' first.")
            sys.exit(1)
    except Exception as e:
        print(f"❌ Error checking authentication: {e}")
        sys.exit(1)
    
    # Find all issue files
    issue_files = glob.glob(".github/ISSUES_PHASE*.md")
    issue_files.sort()
    
    if not issue_files:
        print("❌ No issue files found in .github/ directory")
        sys.exit(1)
    
    print(f"📁 Found {len(issue_files)} issue files:")
    for file_path in issue_files:
        print(f"   - {file_path}")
    print("")
    
    # Parse all issues
    all_issues = []
    for file_path in issue_files:
        print(f"📖 Parsing {file_path}...")
        issues = parse_issue_file(file_path)
        all_issues.extend(issues)
        print(f"   Found {len(issues)} issues")
    
    print(f"\n📋 Total issues found: {len(all_issues)}")
    
    # Filter out already existing issues (Phases 1-2, issues #10-#15)
    existing_issue_numbers = [10, 11, 12, 13, 14, 15]
    new_issues = [issue for issue in all_issues if issue['number'] not in existing_issue_numbers]
    
    print(f"📋 Issues already exist: {len(existing_issue_numbers)} (Phases 1-2)")
    print(f"📋 New issues to create: {len(new_issues)} (Phases 3-7)")
    print("")
    
    # Create new issues
    created_count = 0
    for issue_data in new_issues:
        if create_issue(issue_data):
            created_count += 1
    
    print(f"\n🎉 Successfully created {created_count} out of {len(new_issues)} new issues!")
    print("")
    print("🔗 View issues at: https://github.com/go-foundations/workerpool/issues")
    print("📊 View milestones at: https://github.com/go-foundations/workerpool/milestones")

if __name__ == "__main__":
    main()
