# Employee Grievances Module - Designer Documentation

## Grievance Submission Form

```
┌─ Submit Grievance ──────────────────┐
│ ☐ Submit Anonymously                │
│                                      │
│ Grievance Type:                      │
│ ○ Harassment                        │
│ ○ Discrimination                    │
│ ○ Workload Issues                   │
│ ● Compensation                      │
│ ○ Other                             │
│                                      │
│ Against: [Select Person/Dept]       │
│                                      │
│ Description:                         │
│ ┌─────────────────────────────────┐ │
│ │ Unfair treatment regarding...   │ │
│ │                                 │ │
│ └─────────────────────────────────┘ │
│                                      │
│ Attachments: [Upload Files]         │
│                                      │
│ [Submit Confidentially]              │
└──────────────────────────────────────┘
```

## HR Dashboard View

```
┌─ Grievances Dashboard ──────────────┐
│ ┌─────────────────────────────────┐ │
│ │ 🔴 CRITICAL - Harassment        │ │
│ │    Employee #1234               │ │
│ │    Raised: 2 days ago           │ │
│ │    Status: Under Investigation  │ │
│ │    [View Details]               │ │
│ └─────────────────────────────────┘ │
│                                      │
│ ┌─────────────────────────────────┐ │
│ │ 🟡 MEDIUM - Workload            │ │
│ │    John Doe                     │ │
│ │    Raised: 5 days ago           │ │
│ │    Status: Open                 │ │
│ │    [Assign Investigator]        │ │
│ └─────────────────────────────────┘ │
│                                      │
│ Stats:                               │
│ Open: 12 | Investigating: 8         │
│ Resolved: 45 | Avg Resolution: 12d  │
└──────────────────────────────────────┘
```

## Resolution Tracking

```
┌─ Grievance #GRV-2024-0042 ──────────┐
│ Status: Under Investigation          │
│                                      │
│ Timeline:                            │
│ ✓ Submitted      - Jan 15, 2024     │
│ ✓ Acknowledged   - Jan 15, 2024     │
│ ✓ Assigned       - Jan 16, 2024     │
│ ⏳ Investigation  - In Progress      │
│ ⏳ Resolution     - Pending          │
│                                      │
│ [Add Update] [Close Grievance]       │
└──────────────────────────────────────┘
```

## Color Scheme
- Critical/Harassment: Red
- High Priority: Orange
- Medium: Yellow
- Low: Blue
- Resolved: Green
- Confidential: Lock icon overlay
