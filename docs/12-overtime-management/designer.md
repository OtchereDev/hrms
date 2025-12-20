# Overtime Management Module - Designer Documentation

## Overtime Request Form

```
┌─ Overtime Request ──────────────────┐
│ Employee: John Doe                  │
│ Period: Jan 15 - Jan 20, 2024       │
│                                      │
│ Total Hours: 12.5                   │
│ Type: ○ Paid  ● Compensatory Off    │
│                                      │
│ Reason:                              │
│ ┌─────────────────────────────────┐ │
│ │ Critical product launch support │ │
│ └─────────────────────────────────┘ │
│                                      │
│ [Submit Request]                     │
└──────────────────────────────────────┘
```

## Approval Workflow

```
┌─ Pending Overtime Requests ─────────┐
│ ┌─────────────────────────────────┐ │
│ │ John Doe - 12.5 hrs             │ │
│ │ Jan 15-20, 2024                 │ │
│ │ Type: Compensatory              │ │
│ │ [✓ Approve] [✗ Reject]          │ │
│ └─────────────────────────────────┘ │
└──────────────────────────────────────┘
```

## Color Scheme
- Approved: Green
- Pending: Orange
- Rejected: Red
- Weekend overtime: Purple highlight
