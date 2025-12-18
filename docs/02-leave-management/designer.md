# Leave Management Module - Designer Documentation

## Overview

This document provides UI/UX specifications, visual design guidelines, and user flow diagrams for the Leave Management module.

---

## Design System

### Color Palette

**Leave Balance Indicators**:
- High Balance (>5 days): `#2ECC71` (Green)
- Medium Balance (2-5 days): `#FFD43B` (Yellow)
- Low Balance (<2 days): `#E74C3C` (Red)

**Status Colors**:
- Open: `#FFA00A` (Orange)
- Approved: `#2ECC71` (Green)
- Rejected: `#E74C3C` (Red)
- Cancelled: `#95A5A6` (Gray)

### Typography

**Leave Balance Display**: 24px, Bold
**Field Labels**: 13px, Medium
**Helper Text**: 12px, Regular
**Error Messages**: 13px, Medium, Red

---

## Leave Application Form

### Desktop Layout

```
┌─────── Leave Application ─────────────────────────────────┐
│ HR-LAP-2024-00001                                [Actions] │
│ ● Open                                                     │
└────────────────────────────────────────────────────────────┘

┌─── Employee Details ──────────────────────────────────────┐
│ ┌──────────────────────┬──────────────────────────────┐  │
│ │ Employee *           │ Leave Type *                 │  │
│ │ John Doe             │ [Privilege Leave ▼]          │  │
│ │ [Auto-filled]        │                              │  │
│ │                      │ Leave Balance: 12.5 days     │  │
│ │ Employee Name        │ [Bold, color-coded]          │  │
│ │ John Doe             │                              │  │
│ │                      │                              │  │
│ │ Department           │ Company                      │  │
│ │ Engineering          │ Acme Corp                    │  │
│ └──────────────────────┴──────────────────────────────┘  │
└────────────────────────────────────────────────────────────┘

┌─── Leave Period ──────────────────────────────────────────┐
│ ┌──────────────────────┬──────────────────────────────┐  │
│ │ From Date *          │ To Date *                    │  │
│ │ [📅 Jan 15, 2024]    │ [📅 Jan 19, 2024]            │  │
│ │                      │                              │  │
│ │ ☐ Half Day           │ Half Day Date                │  │
│ │                      │ [Date picker - disabled]     │  │
│ │                      │                              │  │
│ │ Total Leave Days     │                              │  │
│ │ 5.0 days             │                              │  │
│ │ [Large, Bold]        │                              │  │
│ └──────────────────────┴──────────────────────────────┘  │
└────────────────────────────────────────────────────────────┘

┌─── Details ───────────────────────────────────────────────┐
│ Description (Reason for Leave)                            │
│ ┌──────────────────────────────────────────────────────┐ │
│ │ Family vacation                                      │ │
│ │                                                      │ │
│ │                                                      │ │
│ └──────────────────────────────────────────────────────┘ │
└────────────────────────────────────────────────────────────┘

┌─── Approval ──────────────────────────────────────────────┐
│ ┌──────────────────────┬──────────────────────────────┐  │
│ │ Leave Approver       │ Status                       │  │
│ │ Jane Manager         │ [Open ▼]                     │  │
│ │                      │                              │  │
│ │ Posting Date         │ ☑ Follow via Email           │  │
│ │ Jan 10, 2024         │                              │  │
│ └──────────────────────┴──────────────────────────────┘  │
└────────────────────────────────────────────────────────────┘
```

### Leave Balance Dashboard (Sidebar)

```
┌─── Leave Balance ─────────────────────┐
│                                       │
│ Privilege Leave                       │
│ 12.5 / 20 days                        │
│ ████████████░░░░░░░░ 62%              │
│                                       │
│ Sick Leave                            │
│ 8 / 10 days                           │
│ ████████████████░░░░ 80%              │
│                                       │
│ Casual Leave                          │
│ 1.5 / 12 days     ⚠ Low              │
│ ███░░░░░░░░░░░░░░░░░ 12%              │
│                                       │
│ [View All Allocations →]             │
└───────────────────────────────────────┘
```

**Visual Specifications**:
- Card padding: 16px
- Leave type font: 14px, Semi-bold
- Balance font: 20px, Bold
- Progress bar height: 8px
- Progress bar radius: 4px
- Warning icon: 16px, Yellow

### Half Day Indicator

When half day is checked:
```
┌──────────────────────────────────────┐
│ ◐ Half Day Leave                     │
│   Select the date for half day leave │
│                                      │
│   Half Day Date *                    │
│   [📅 Jan 15, 2024]                  │
│                                      │
│   Total: 0.5 days                    │
└──────────────────────────────────────┘
```

### Balance Warning

When exceeding balance:
```
┌──────────────────────────────────────┐
│ ⚠ Warning                             │
│                                      │
│ You have 1.5 days of Casual Leave.   │
│ You are applying for 3.0 days.       │
│                                      │
│ This will result in negative balance │
│ of -1.5 days.                        │
│                                      │
│ [Cancel]  [Apply Anyway]             │
└──────────────────────────────────────┘
```

---

## Mobile App Design

### Leave Dashboard

**Screen Layout**:
```
┌─────────────────────────────────┐
│ ☰  Leave                    🔔  │
├─────────────────────────────────┤
│                                 │
│ ┌─────────┬─────────┬─────────┐│
│ │ PL      │ SL      │ CL      ││
│ │ 12.5/20 │ 8/10    │ 1.5/12  ││
│ │ ████▓   │ ████▓   │ █▓▓▓    ││
│ └─────────┴─────────┴─────────┘│
│                                 │
│ ┌─────────────────────────────┐│
│ │ + Request a Leave           ││
│ └─────────────────────────────┘│
│                                 │
│ Recent Leaves                   │
│ ┌─────────────────────────────┐│
│ │ Privilege Leave             ││
│ │ Jan 15 - Jan 19 (5 days)    ││
│ │ ● Approved                  ││
│ └─────────────────────────────┘│
│ ┌─────────────────────────────┐│
│ │ Sick Leave                  ││
│ │ Dec 20 - Dec 21 (2 days)    ││
│ │ ● Approved                  ││
│ └─────────────────────────────┘│
│                                 │
│ Upcoming Holidays               │
│ ┌─────────────────────────────┐│
│ │ 🎉 Republic Day             ││
│ │ Jan 26, 2024                ││
│ └─────────────────────────────┘│
└─────────────────────────────────┘
```

**Balance Card Specifications**:
- Card width: 30% (3 per row)
- Height: 80px
- Border radius: 8px
- Background: Gradient based on utilization
- Shadow: 0 2px 4px rgba(0,0,0,0.1)

### Leave Application Form (Mobile)

```
┌─────────────────────────────────┐
│ ← Request Leave              ✓  │
├─────────────────────────────────┤
│                                 │
│ Leave Type *                    │
│ Privilege Leave (12.5 days) ▼   │
│                                 │
│ From Date *                     │
│ 📅 Jan 15, 2024                 │
│                                 │
│ To Date *                       │
│ 📅 Jan 19, 2024                 │
│                                 │
│ ◐ Half Day                      │
│   ☐ Yes                         │
│                                 │
│ ┌─────────────────────────────┐│
│ │ Total Leave Days            ││
│ │ 5.0 days                    ││
│ └─────────────────────────────┘│
│                                 │
│ Reason                          │
│ ┌─────────────────────────────┐│
│ │ Family vacation             ││
│ │                             ││
│ │                             ││
│ └─────────────────────────────┘│
│                                 │
│ ┌─────────────────────────────┐│
│ │ Submit Request              ││
│ └─────────────────────────────┘│
│                                 │
└─────────────────────────────────┘
```

**Mobile Specifications**:
- Touch targets: Minimum 44x44px
- Input fields: 48px height
- Padding: 16px horizontal, 12px vertical
- Submit button: 48px height, full width
- Font sizes: 16px body, 14px labels

### Leave Calendar View

```
┌─────────────────────────────────┐
│ January 2024            < >     │
├─────────────────────────────────┤
│ S  M  T  W  T  F  S             │
│                                 │
│ ─  1  2  3  4  5  6             │
│ 7  8  9  10 11 12 13            │
│ 14 🟢 🟢 🟢 🟢 🟢 20            │
│ 21 22 23 24 25 🔴 27            │
│ 28 29 30 31 ─  ─  ─             │
│                                 │
│ ┌─────────────────────────────┐│
│ │ 🟢 Leave Applied            ││
│ │ 🔴 Holiday                  ││
│ │ ⚪ Available                ││
│ └─────────────────────────────┘│
└─────────────────────────────────┘
```

**Calendar Colors**:
- Leave days: `#2ECC71` (Green circle)
- Holidays: `#E74C3C` (Red circle)
- Weekends: `#F0F0F0` (Light gray background)
- Today: `#2490EF` (Blue outline)

---

## List View Design

### Web List View

```
┌─────── Leave Applications ─────────────────────────────────┐
│ [+ New] [My Leaves] [Team Leaves] [⚙ Filters] [↻ Refresh] │
├────────────────────────────────────────────────────────────┤
│ ☑ Employee      Leave Type    From → To      Days  Status │
├────────────────────────────────────────────────────────────┤
│ ☐ John Doe     Privilege     Jan 15→19      5.0   ● Appr. │
│   Engineering  Leave                                       │
├────────────────────────────────────────────────────────────┤
│ ☐ Jane Smith   Sick Leave    Jan 10→11      2.0   ● Open  │
│   Sales                                                    │
├────────────────────────────────────────────────────────────┤
│ ☐ Bob Johnson  Casual        Jan 8 (Half)   0.5   ● Appr. │
│   Marketing    Leave                                       │
└────────────────────────────────────────────────────────────┘
```

**List Specifications**:
- Row height: 56px
- Checkbox: 18x18px
- Status badge: 6px circle + text
- Hover effect: Background #F8F9FA
- Selected row: Background #E3F2FD

### Mobile List View

```
┌─────────────────────────────────┐
│ [My Leaves] [Team Leaves]       │
├─────────────────────────────────┤
│ ┌─────────────────────────────┐│
│ │ Privilege Leave             ││
│ │ Jan 15 - Jan 19, 2024       ││
│ │ 5.0 days   ● Approved       ││
│ │ Family vacation             ││
│ └─────────────────────────────┘│
│                                 │
│ ┌─────────────────────────────┐│
│ │ Sick Leave                  ││
│ │ Jan 10 - Jan 11, 2024       ││
│ │ 2.0 days   ● Open           ││
│ │ Medical appointment         ││
│ └─────────────────────────────┘│
│                                 │
│ ┌─────────────────────────────┐│
│ │ Casual Leave (Half Day)     ││
│ │ Jan 8, 2024                 ││
│ │ 0.5 days   ● Approved       ││
│ │ Personal work               ││
│ └─────────────────────────────┘│
└─────────────────────────────────┘
```

**Card Specifications**:
- Margin: 8px vertical
- Padding: 16px
- Border-radius: 8px
- Shadow: 0 1px 3px rgba(0,0,0,0.12)
- Tap area: Full card

---

## User Flows

### Apply for Leave Flow

```
Start
  ↓
[Dashboard] → Click "Request Leave"
  ↓
[Form Page]
  ↓
Select Leave Type
  → Auto-fetch balance
  → Show balance indicator
  ↓
Select From Date
  ↓
Select To Date
  → Auto-calculate days
  → Show total
  ↓
(Optional) Check Half Day
  → Show half day date picker
  → Recalculate days
  ↓
Enter Reason
  ↓
Review Summary:
  ┌────────────────────────┐
  │ Leave Type: PL         │
  │ Period: Jan 15-19      │
  │ Days: 5.0              │
  │ Balance: 12.5 → 7.5    │
  └────────────────────────┘
  ↓
Click "Submit"
  ↓
[If balance sufficient]
  → Show success
  → Navigate to list
  ↓
[If balance insufficient]
  → Show warning dialog
  → [Cancel] or [Apply Anyway]
```

### Approval Flow (Manager)

```
Manager receives notification
  ↓
[Email/Push Notification]
"John Doe applied for 5 days leave"
  ↓
Click notification
  ↓
[Leave Application Page]
  → View employee details
  → View leave reason
  → Check team calendar
  → Check balance
  ↓
Make Decision
  ↓
[Approve]              [Reject]
  ↓                      ↓
Select status        Enter reason
  ↓                      ↓
Submit               Submit
  ↓                      ↓
Employee notified    Employee notified
```

---

## Responsive Design

### Breakpoints

- Mobile: < 768px (single column)
- Tablet: 768px - 1024px (2 columns)
- Desktop: > 1024px (full layout)

### Mobile Adaptations

1. **Form Fields**: Full width, stacked vertically
2. **Balance Dashboard**: Horizontal scroll cards
3. **Date Pickers**: Native mobile date picker
4. **Buttons**: Full width, bottom of form

---

## Accessibility

### Color Contrast
- Text on white: Minimum 4.5:1
- Status badges: Include icons, not just color
- Balance indicators: Text + color

### Keyboard Navigation
- Tab order: Top to bottom, left to right
- Enter: Submit form
- Escape: Close dialogs
- Space: Toggle checkboxes

### Screen Reader
- ARIA labels on all form fields
- Status announcements on changes
- Balance updates announced

---

## Animation & Transitions

**Page Transitions**: 300ms ease-in-out
**Balance Updates**: Smooth number counting (500ms)
**Status Badge Change**: Fade + scale (200ms)
**Error Shake**: Horizontal shake (400ms)

**Balance Bar Animation**:
```css
.balance-bar {
  transition: width 600ms cubic-bezier(0.4, 0.0, 0.2, 1);
}
```

---

## Empty States

### No Leave Balance

```
┌──────────────────────────────────┐
│         📭                       │
│                                  │
│    No Leave Balance              │
│                                  │
│  You don't have any leave        │
│  allocations yet. Contact HR.    │
│                                  │
│  [Contact HR]                    │
└──────────────────────────────────┘
```

### No Leave Applications

```
┌──────────────────────────────────┐
│         📝                       │
│                                  │
│  No Leave Applications           │
│                                  │
│  You haven't applied for any     │
│  leaves yet.                     │
│                                  │
│  [Request Leave]                 │
└──────────────────────────────────┘
```

---

## Error States

### Insufficient Balance

```
┌──────────────────────────────────┐
│ ⚠ Insufficient Leave Balance     │
│                                  │
│ You have 1.5 days available.     │
│ You're requesting 3.0 days.      │
│                                  │
│ Would you like to proceed with   │
│ negative balance?                │
│                                  │
│ [Go Back] [Proceed]              │
└──────────────────────────────────┘
```

### Overlapping Leave

```
┌──────────────────────────────────┐
│ ❌ Leave Already Exists           │
│                                  │
│ You have already applied for     │
│ leave during this period:        │
│                                  │
│ Jan 15-17: Privilege Leave       │
│                                  │
│ [View Application] [Close]       │
└──────────────────────────────────┘
```

---

This designer documentation provides complete visual specifications for implementing the Leave Management module UI/UX.
