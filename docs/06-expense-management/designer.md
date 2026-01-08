# Expense Management Module - Designer Documentation

## Overview

UI/UX specifications and visual design patterns for employee expense claims, approvals, and payment management.

---

## Design System

### Color Palette

**Status Colors**:
```
Draft:            #9B9B9B (Gray)
Submitted:        #2490EF (Blue)
Approved:         #98D85B (Green)
Rejected:         #F56B6B (Red)
Paid:             #98D85B (Green)
Unpaid:           #FFA00A (Orange)
```

### Typography

**Form Headers**: 18px, Semi-bold
**Section Titles**: 14px, Semi-bold
**Labels**: 12px, Regular
**Values**: 14px, Regular
**Totals**: 16px, Bold

---

## Expense Claim Form

### Layout

```
┌─────────────────────────────────────────────────────────────┐
│  Expense Claim                                [Submit ▼] [×]│
│  EXP-CLAIM-2024-00123                                       │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Status: ● Unpaid                                           │
│  Approval: ● Approved                                       │
│                                                             │
│  ┌─ Claim Details ──────────────────────────────────────┐  │
│  │                                                       │  │
│  │  ┌──────────────────────┬──────────────────────────┐ │  │
│  │  │ Employee *           │ Employee Name            │ │  │
│  │  │ ┌──────────────────┐ │ John Doe                 │ │  │
│  │  │ │ 🔍 Search...     │ │ [Auto-filled]            │ │  │
│  │  │ └──────────────────┘ │                          │ │  │
│  │  └──────────────────────┴──────────────────────────┘ │  │
│  │                                                       │  │
│  │  ┌──────────────────────┬──────────────────────────┐ │  │
│  │  │ Company *            │ Posting Date *           │ │  │
│  │  │ Acme Corp            │ 📅 Jan 15, 2024          │ │  │
│  │  └──────────────────────┴──────────────────────────┘ │  │
│  │                                                       │  │
│  │  ┌──────────────────────┬──────────────────────────┐ │  │
│  │  │ Expense Approver     │ Department               │ │  │
│  │  │ 🔍 Manager Name      │ Engineering              │ │  │
│  │  └──────────────────────┴──────────────────────────┘ │  │
│  │                                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─ Project Tracking ───────────────────────────────────┐  │
│  │  ┌──────────────────────┬──────────────────────────┐ │  │
│  │  │ Project              │ Task                     │ │  │
│  │  │ [Optional]           │ [Optional - filtered]    │ │  │
│  │  │                      │                          │ │  │
│  │  │ Cost Center          │ Department               │ │  │
│  │  │ [Optional]           │ [Optional]               │ │  │
│  │  └──────────────────────┴──────────────────────────┘ │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─ Expenses ───────────────────────────────────────────┐  │
│  │                                                       │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ Date  │ Type     │ Amount │ Sanctioned│ Desc.  │ │  │
│  │  ├─────────────────────────────────────────────────┤ │  │
│  │  │ Jan 10│ Travel   │  $350  │   $350    │ Flight │ │  │
│  │  │ Jan 11│ Hotel    │  $120  │   $120    │ 2 nts  │ │  │
│  │  │ Jan 12│ Meals    │   $75  │    $50    │ Food   │ │  │
│  │  │ Jan 13│ Taxi     │   $45  │    $45    │ Uber   │ │  │
│  │  │ [+ Add Row]                                     │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ Total Claimed Amount:              $590         │ │  │
│  │  │ Total Sanctioned Amount:           $565         │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─ Taxes and Charges ──────────────────────────────────┐  │
│  │                                                       │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ Account          │ Rate  │ Amount    │ Total    │ │  │
│  │  ├─────────────────────────────────────────────────┤ │  │
│  │  │ Service Tax      │ 10%   │  $56.50   │ $621.50  │ │  │
│  │  │ [+ Add Row]                                     │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ Total Taxes and Charges:           $56.50       │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─ Advances ───────────────────────────────────────────┐  │
│  │                                                       │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ Advance│ Date    │ Paid   │ Unclaimed│ Allocated││ │
│  │  ├─────────────────────────────────────────────────┤ │  │
│  │  │ ADV-001│ Jan 5   │  $500  │   $500   │  $500    ││ │
│  │  │ [Auto-populated from employee advances]         │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ Total Advance Amount:              $500         │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─ Payment Summary ────────────────────────────────────┐  │
│  │                                                       │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │                                                 │ │  │
│  │  │  Total Sanctioned Amount:         $565.00       │ │  │
│  │  │  Total Taxes and Charges:          $56.50       │ │  │
│  │  │  Total Advance Amount:           -$500.00       │ │  │
│  │  │  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │ │  │
│  │  │  Grand Total (To be Paid):        $121.50       │ │  │
│  │  │                                                 │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ ☐ Paid                                          │ │  │
│  │  │                                                 │ │  │
│  │  │ Mode of Payment                                 │ │  │
│  │  │ [Cash/Bank ▼]                                   │ │  │
│  │  │                                                 │ │  │
│  │  │ Payable Account *                               │ │  │
│  │  │ [Account ▼]                                     │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
└─────────────────────────────────────────────────────────────┘

[Create ▼]  [Payment]  [View ▼]
```

### Status Indicators

**Approval Status Badge**:
```
┌─ Approval Status ────┐
│ ● Approved           │  Green badge
│ ● Rejected           │  Red badge
│ ● Draft              │  Gray badge
└──────────────────────┘
```

**Payment Status Badge**:
```
┌─ Payment Status ─────┐
│ ● Paid               │  Green badge
│ ● Unpaid             │  Orange badge
│ ● Draft              │  Gray badge
└──────────────────────┘
```

### Expense Detail Row

```
┌───────────────────────────────────────────────────────────┐
│ Date        │ Expense Type  │ Description                 │
│ 📅 Jan 10  │ 🔍 Travel ▼  │ Flight to NY                │
│                                                           │
│ Amount *    │ Sanctioned    │ Account                     │
│ $ 350.00    │ $ 350.00      │ Travel Expenses - T         │
│                                                           │
│ Cost Center │ Project       │                             │
│ [Optional]  │ [Optional]    │                       [🗑] │
└───────────────────────────────────────────────────────────┘
```

### Multi-currency Display

When currency ≠ company currency:
```
┌─ Amounts ────────────────────────────────────┐
│  Currency: USD   Exchange Rate: 1.0          │
│                                              │
│  Total Sanctioned Amount (USD):   $565.00   │
│  Total Sanctioned Amount (INR):  ₹46,123    │
│                                              │
│  Grand Total (USD):               $121.50   │
│  Grand Total (INR):               ₹9,924    │
└──────────────────────────────────────────────┘
```

### List View

```
┌─────────────────────────────────────────────────────────────┐
│  Expense Claims                           [+ New] [⚙] [⟳]  │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Filters:  [Employee ▼] [Status ▼] [Date Range]            │
│                                                             │
├─────┬──────────────┬──────────┬──────────┬────────┬────────┤
│ ●   │ ID           │ Employee │ Amount   │ Status │ Approval│
├─────┼──────────────┼──────────┼──────────┼────────┼────────┤
│ ☑   │ EXP-001      │ John Doe │  $565.00 │ ● Paid│ ●Approved│
│ ☑   │ EXP-002      │ Jane S   │  $235.50 │●Unpaid│ ●Approved│
│ ☑   │ EXP-003      │ Bob J    │  $890.00 │ ○Draft│ ○Draft  │
└─────┴──────────────┴──────────┴──────────┴────────┴────────┘
```

**Status Indicators**:
- ● Green: Paid/Approved
- ● Orange: Unpaid
- ● Red: Rejected
- ○ Gray: Draft

---

## Custom Buttons and Actions

### Create Payment

When status = Unpaid and submitted:
```
┌─────────────────────────────────┐
│  Create ▼                       │
│  ├─ Payment                     │
│  └─ [Opens Payment Entry form] │
└─────────────────────────────────┘
```

### View Actions

```
┌─────────────────────────────────┐
│  View ▼                         │
│  ├─ Accounting Ledger           │
│  ├─ Bank Entries                │
│  └─ Exchange Gain/Loss Journals │
└─────────────────────────────────┘
```

---

## Mobile View (< 768px)

```
┌───────────────────────────────┐
│  ☰  Expense Claim        [×] │
│  EXP-CLAIM-2024-00123         │
├───────────────────────────────┤
│                               │
│  Status: ● Unpaid             │
│  Approval: ● Approved         │
│                               │
│  ┌─ Summary ─────────────┐   │
│  │ Employee: John Doe    │   │
│  │ Date: Jan 15, 2024    │   │
│  │                       │   │
│  │ Total: $565.00        │   │
│  │ Advances: -$500.00    │   │
│  │ ━━━━━━━━━━━━━━━━━━  │   │
│  │ To Pay: $121.50       │   │
│  └───────────────────────┘   │
│                               │
│  [View Expenses] ▼            │
│  [View Advances] ▼            │
│                               │
│  [Create Payment]             │
│                               │
└───────────────────────────────┘
```

---

## Validation Messages

**Self-Approval Prevented**:
```
┌────────────────────────────────────┐
│ ⚠ Not Allowed                     │
│                                    │
│ Self-approval for Expense Claims   │
│ is not allowed                     │
│                                    │
│ [OK]                               │
└────────────────────────────────────┘
```

**Payable Account Required**:
```
┌────────────────────────────────────┐
│ ⚠ Validation Error                │
│                                    │
│ Payable Account is mandatory to    │
│ submit an Expense Claim            │
│                                    │
│ [OK]                               │
└────────────────────────────────────┘
```

**Sanctioned > Claimed**:
```
┌────────────────────────────────────┐
│ ⚠ Invalid Amount                  │
│                                    │
│ Row 3: Sanctioned Amount cannot    │
│ be greater than Claimed Amount     │
│                                    │
│ [OK]                               │
└────────────────────────────────────┘
```

---

## Animation and Transitions

**Row Calculations**: Update on blur with 200ms fade
**Total Updates**: Smooth number counting animation
**Status Change**: Color transition 300ms
**Form Field Focus**: 200ms ease
**Button Hover**: 100ms ease

---

This designer documentation provides complete UI/UX specifications for the Expense Management module.
