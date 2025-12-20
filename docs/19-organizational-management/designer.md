# Organizational Management Module - Designer Documentation

## Organizational Chart

```
┌─ Organization Chart ────────────────────────────────┐
│                                                     │
│                   [CEO]                             │
│                  Jane Doe                           │
│                     │                               │
│         ┌───────────┼───────────┐                  │
│         │           │           │                   │
│     [CTO]       [CFO]       [COO]                   │
│   Bob Wilson   Alice Brown  Tom Davis              │
│         │           │           │                   │
│    ┌────┼────┐     │      ┌────┼────┐             │
│    │    │    │     │      │    │    │              │
│  [Eng][QA][DevOps][Acc][Sales][Ops][HR]           │
│   5    3    2      4    8     6    3               │
│                                                     │
│ [Expand All] [Download] [Print]                    │
└─────────────────────────────────────────────────────┘
```

## Department Tree View

```
┌─ Departments ───────────────────────┐
│                                      │
│ ▼ Engineering                        │
│   ├─ Frontend Team (12)              │
│   ├─ Backend Team (18)               │
│   ├─ Mobile Team (8)                 │
│   └─ DevOps Team (5)                 │
│                                      │
│ ▼ Sales & Marketing                  │
│   ├─ Inside Sales (15)               │
│   ├─ Field Sales (22)                │
│   └─ Marketing (10)                  │
│                                      │
│ ▼ Operations                         │
│   ├─ Customer Support (20)           │
│   ├─ Logistics (8)                   │
│   └─ Quality Assurance (12)          │
│                                      │
│ ► Finance                            │
│ ► Human Resources                    │
│                                      │
│ [Add Department] [Reorganize]        │
└──────────────────────────────────────┘
```

## Branch Management

```
┌─ Branches ──────────────────────────┐
│                                      │
│ ┌─────────────────────────────────┐ │
│ │ 🏢 Headquarters                 │ │
│ │    San Francisco, CA            │ │
│ │    Employees: 150               │ │
│ │    Head: Jane Doe               │ │
│ │    [View Details]               │ │
│ └─────────────────────────────────┘ │
│                                      │
│ ┌─────────────────────────────────┐ │
│ │ 🏢 East Coast Office            │ │
│ │    New York, NY                 │ │
│ │    Employees: 75                │ │
│ │    Head: Bob Wilson             │ │
│ │    [View Details]               │ │
│ └─────────────────────────────────┘ │
│                                      │
│ ┌─────────────────────────────────┐ │
│ │ 🏢 APAC Office                  │ │
│ │    Singapore                    │ │
│ │    Employees: 45                │ │
│ │    Head: Alice Tan              │ │
│ │    [View Details]               │ │
│ └─────────────────────────────────┘ │
│                                      │
│ [Add Branch]                         │
└──────────────────────────────────────┘
```

## Designation & Grade Matrix

```
┌─ Career Ladder ─────────────────────┐
│                                      │
│ Grade  Designation       Salary Band │
│ L1     Junior Engineer   $40k-$60k   │
│ L2     Engineer          $60k-$85k   │
│ L3     Senior Engineer   $85k-$120k  │
│ L4     Staff Engineer    $120k-$160k │
│ M1     Engineering Mgr   $130k-$180k │
│ M2     Senior Mgr        $160k-$220k │
│ D1     Director          $200k-$300k │
│ E1     VP Engineering    $250k-$400k │
│                                      │
│ [Edit Grades] [Salary Structure]     │
└──────────────────────────────────────┘
```

## Color Scheme
- HQ/Main branch: Blue
- Regional offices: Green
- Department groups: Gray background
- Active employees: Bold text
- Reporting lines: Dotted connectors
