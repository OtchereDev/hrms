# Salary Withholding Module - Designer Documentation

## Employee Advance Form

```
┌─ Employee Advance ──────────────────┐
│ Employee: John Doe                  │
│ Purpose: Medical Emergency          │
│ Advance Amount: $5,000              │
│                                      │
│ Repayment:                           │
│ Monthly Deduction: $500             │
│ Start From: March 2024              │
│ Duration: 10 months                 │
│                                      │
│ Outstanding: $5,000                 │
│ [Submit Request]                     │
└──────────────────────────────────────┘
```

## Repayment Schedule View

```
┌─ Repayment Schedule ────────────────┐
│ ┌─────────────────────────────────┐ │
│ │ Month      Amount    Status     │ │
│ │ Mar 2024   $500      ✓ Paid     │ │
│ │ Apr 2024   $500      ⏳ Pending │ │
│ │ May 2024   $500      ⏳ Pending │ │
│ │ ...                              │ │
│ └─────────────────────────────────┘ │
│                                      │
│ Total: $5,000                       │
│ Paid: $500                          │
│ Outstanding: $4,500                 │
└──────────────────────────────────────┘
```

## Salary Slip Deduction Section

```
┌─ Deductions ────────────────────────┐
│ Professional Tax:          $200     │
│ Employee PF:             $1,800     │
│ Loan Repayment:            $500     │
│ Advance Recovery:          $500     │
│ ──────────────────────────────────  │
│ Total Deductions:        $3,000     │
└──────────────────────────────────────┘
```

## Color Scheme
- Active deductions: Orange
- Completed: Green
- Overdue: Red
