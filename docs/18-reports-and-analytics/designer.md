# Reports and Analytics Module - Designer Documentation

## HR Analytics Dashboard

```
┌─ HR Analytics ──────────────────────────────────────┐
│ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐   │
│ │ Total       │ │ New Hires   │ │ Attrition   │   │
│ │ Employees   │ │ This Month  │ │ Rate        │   │
│ │    250      │ │     12      │ │   8.5%      │   │
│ └─────────────┘ └─────────────┘ └─────────────┘   │
│                                                     │
│ Headcount Trend (Last 12 Months)                   │
│ ┌───────────────────────────────────────────────┐ │
│ │     ██                                        │ │
│ │   ████                                        │ │
│ │ ██████  ██      ██    ██                      │ │
│ │ ██████████  ██████  ████  ██                  │ │
│ │ ████████████████████████████                  │ │
│ └───────────────────────────────────────────────┘ │
│ J F M A M J J A S O N D                            │
│                                                     │
│ Department Distribution                             │
│ Engineering:  ██████████████░░ 45%                 │
│ Sales:        ████████░░░░░░░ 25%                 │
│ Marketing:    ████░░░░░░░░░░░ 15%                 │
│ HR:           ██░░░░░░░░░░░░░  8%                 │
│ Finance:      ██░░░░░░░░░░░░░  7%                 │
└─────────────────────────────────────────────────────┘
```

## Report Builder Interface

```
┌─ Custom Report Builder ─────────────────┐
│ Report Name: [Active Employees]         │
│                                          │
│ Source: Employee                         │
│                                          │
│ Columns:                                 │
│ ☑ Employee Name                         │
│ ☑ Department                            │
│ ☑ Designation                           │
│ ☑ Date of Joining                       │
│ ☐ Salary                                │
│                                          │
│ Filters:                                 │
│ Status = Active                          │
│ Company = ABC Corp                       │
│ [+ Add Filter]                           │
│                                          │
│ Group By: Department                     │
│ Sort By: Date of Joining (Desc)         │
│                                          │
│ [Preview] [Export Excel] [Save Report]   │
└──────────────────────────────────────────┘
```

## Salary Register Report

```
┌─ Salary Register - January 2024 ───────────────────┐
│ Filters: All Departments | All Grades              │
│                                                     │
│ ┌─────────────────────────────────────────────────┐│
│ │ Name      Dept    Gross    Ded.    Net    Bank ││
│ │ John Doe  Eng    $5,000   $500  $4,500  ✓ Paid ││
│ │ Jane Smith Sales  $4,500   $450  $4,050  ✓ Paid ││
│ │ Bob Wilson IT     $6,000   $600  $5,400  ⏳ Pend││
│ │ ...                                             ││
│ └─────────────────────────────────────────────────┘│
│                                                     │
│ Total: 250 employees                                │
│ Gross Pay: $1,250,000                              │
│ Deductions: $125,000                               │
│ Net Pay: $1,125,000                                │
│                                                     │
│ [Export PDF] [Export Excel] [Email]                │
└─────────────────────────────────────────────────────┘
```

## Attendance Summary Chart

```
┌─ Attendance Summary - January 2024 ─────┐
│                                          │
│ ┌──────────────────────────────────────┐│
│ │                                      ││
│ │  Present:    22 days  ████████  88% ││
│ │  Leave:       2 days  ██░░░░░░   8% ││
│ │  Half Day:    1 day   █░░░░░░░   4% ││
│ │                                      ││
│ └──────────────────────────────────────┘│
│                                          │
│ Late Arrivals: 3                         │
│ Early Exits: 1                           │
│ Overtime: 8 hours                        │
│                                          │
│ [Download Report]                        │
└──────────────────────────────────────────┘
```

## Color Scheme
- Positive metrics: Green
- Negative metrics: Red
- Neutral: Blue
- Warnings: Orange
- Charts: Multi-color gradient
