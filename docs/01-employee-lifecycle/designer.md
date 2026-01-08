# Employee Lifecycle Module - Designer Documentation

## Overview

This document describes the UI/UX design, layouts, visual elements, and user flows for the Employee Lifecycle module. It provides specifications for designers to understand and replicate the visual and interaction design.

---

## Design System

### Colors

**Status Colors**:
- Pending: `#FFA00A` (Orange)
- In Process: `#FFD43B` (Yellow)
- Completed: `#2ECC71` (Green)
- Cancelled: `#E74C3C` (Red)
- Draft: `#95A5A6` (Gray)

**Primary Actions**: `#2490EF` (Blue)

**Secondary Actions**: `#8D99A6` (Gray)

**Success**: `#2ECC71` (Green)

**Error**: `#E74C3C` (Red)

### Typography

**Headings**:
- Page Title: 28px, Semi-bold
- Section Heading: 16px, Semi-bold
- Field Labels: 13px, Medium
- Helper Text: 12px, Regular

**Body Text**: 14px, Regular

**Font Family**: System default (-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto)

### Spacing

**Standard Spacing Scale**:
- XS: 4px
- S: 8px
- M: 16px
- L: 24px
- XL: 32px
- XXL: 48px

**Form Spacing**:
- Between sections: 24px
- Between fields: 16px
- Field label to input: 8px
- Table row height: 36px

---

## Employee Onboarding

### Page Layout

**Header Section**:
```
┌─────────────────────────────────────────────────────────────┐
│ ← Employee Onboarding                              [Actions] │
│ HR-EMP-ONB-2024-00001                                        │
│ [Status Badge: Pending/In Process/Completed]                │
└─────────────────────────────────────────────────────────────┘
```

**Form Layout** (2-column responsive grid):

```
┌───────────────────────┬───────────────────────┐
│ Job Applicant *       │ Company               │
│ [Select dropdown]     │ [Read-only field]     │
│                       │                       │
│ Job Offer *           │ Boarding Status       │
│ [Select dropdown]     │ [Badge]               │
│                       │                       │
│ Employee Onboarding   │ Project               │
│ Template              │ [Link]                │
│ [Select dropdown]     │                       │
└───────────────────────┴───────────────────────┘

┌─── Employee Details ──────────────────────────┐
│                                               │
│ ┌─────────────────┬───────────────────────┐  │
│ │ Employee        │ Date of Joining *     │  │
│ │ [Link/Read-only]│ [Date picker]         │  │
│ │                 │                       │  │
│ │ Employee Name   │ Boarding Begins On *  │  │
│ │ [Read-only]     │ [Date picker]         │  │
│ │                 │                       │  │
│ │ Department      │                       │  │
│ │ [Auto-filled]   │                       │  │
│ │                 │                       │  │
│ │ Designation     │                       │  │
│ │ [Auto-filled]   │                       │  │
│ │                 │                       │  │
│ │ Employee Grade  │                       │  │
│ │ [Auto-filled]   │                       │  │
│ │                 │                       │  │
│ │ Holiday List    │                       │  │
│ │ [Optional]      │                       │  │
│ └─────────────────┴───────────────────────┘  │
└───────────────────────────────────────────────┘

┌─── Onboarding Activities ─────────────────────┐
│                                               │
│ [Table with editable rows]                    │
│ ┌──────┬──────┬──────┬────────┬────────┬──┐  │
│ │ Acti │ User │ Role │ Begin  │ Duration│✓│  │
│ │ vity │      │      │ On     │         │ │  │
│ ├──────┼──────┼──────┼────────┼────────┼──┤  │
│ │ ...  │ ...  │ ...  │ ...    │ ...    │☐│  │
│ └──────┴──────┴──────┴────────┴────────┴──┘  │
│                                               │
│ [+ Add Row]                                   │
│                                               │
│ ☐ Notify users by email                      │
└───────────────────────────────────────────────┘
```

### Visual Elements

**Status Badge**:
```
┌─────────────────┐
│ ● Pending       │  ← Orange circle + text
└─────────────────┘

┌─────────────────┐
│ ● In Process    │  ← Yellow circle + text
└─────────────────┘

┌─────────────────┐
│ ● Completed     │  ← Green circle + text with checkmark
└─────────────────┘
```

**Action Buttons** (top-right):
```
After Submit:
┌─────────┬────────────┬──────────┬──────────────────┐
│ Save ▼  │ Refresh    │ View ▼   │ Create ▼         │
└─────────┴────────────┴──────────┴──────────────────┘
                        │          │
                        │          └─ Employee (primary)
                        │
                        └─ Employee, Project, Task
```

**Activity Table Row**:
```
┌────────────────────────────────────────────────────────────────┐
│ Activity Name │ User      │ Role   │ Begin On │ Duration │ ☑  │
│ (editable)    │ (search)  │(search)│  (days)  │  (days)  │    │
│               │           │        │          │          │    │
│ [Text input]  │[Dropdown] │[Link]  │  [Int]   │  [Int]   │[Cb]│
└────────────────────────────────────────────────────────────────┘
     ↓ (expands on click)
┌────────────────────────────────────────────────────────────────┐
│ Description:                                                    │
│ [Rich text editor with formatting toolbar]                     │
└────────────────────────────────────────────────────────────────┘
```

### User Flows

**Creating New Onboarding**:

```
1. Click "New" → Employee Onboarding
        ↓
2. Select Job Applicant (autocomplete search)
   → Auto-fills: Employee Name, Employee (if exists)
        ↓
3. Select Job Offer (filtered by applicant)
   → Auto-fills: Department, Designation, Company
        ↓
4. (Optional) Select Template
   → Populates activities table
        ↓
5. Set dates: Date of Joining, Boarding Begins On
        ↓
6. Review/Edit activities
   - Adjust timelines
   - Assign users/roles
   - Mark required activities
        ↓
7. Check "Notify users by email" (optional)
        ↓
8. Save (Draft status)
        ↓
9. Review → Submit
   → Creates project and tasks
   → Status: Pending
        ↓
10. Track progress (status updates automatically)
        ↓
11. When ready: Click "Create Employee"
    → Opens pre-filled employee form
        ↓
12. Complete employee form → Save
    → Links back to onboarding
```

**Boarding Progress Visualization**:

List view shows progress:
```
┌─────────────────────────────────────────────────────────┐
│ ● John Doe              Jan 15, 2024    Engineering     │
│   HR-EMP-ONB-2024-00001 ───●─────○────○   In Process    │
│                            60% Complete                  │
└─────────────────────────────────────────────────────────┘
```

---

## Employee Separation

### Page Layout

Similar to onboarding with adaptations for exit process.

**Header**:
```
┌─────────────────────────────────────────────────────────────┐
│ ← Employee Separation                              [Actions] │
│ HR-EMP-SEP-2024-00001                                        │
│ [Status Badge]                                               │
└─────────────────────────────────────────────────────────────┘
```

**Form Layout**:
```
┌───────────────────────┬───────────────────────┐
│ Employee *            │ Company               │
│ [Select dropdown]     │ [Read-only]           │
│                       │                       │
│ Employee Name         │ Boarding Status       │
│ [Read-only]           │ [Badge]               │
│                       │                       │
│ Department            │ Resignation Letter    │
│ [Read-only]           │ Date [Read-only]      │
│                       │                       │
│ Designation           │ Boarding Begins On *  │
│ [Read-only]           │ [Date picker]         │
│                       │                       │
│ Employee Grade        │ Project               │
│ [Read-only]           │ [Link after submit]   │
└───────────────────────┴───────────────────────┘

┌─── Separation Activities ─────────────────────┐
│ Employee Separation Template                  │
│ [Select dropdown]                             │
│                                               │
│ [Activities Table - same as onboarding]       │
│                                               │
│ ☐ Notify users by email                      │
└───────────────────────────────────────────────┘

┌─── Exit Interview Summary ────────────────────┐
│ [Rich text editor]                            │
│                                               │
│ Space for exit interview notes and summary    │
└───────────────────────────────────────────────┘
```

### Visual Elements

**Warning Indicator** (if resignation date is close):
```
┌────────────────────────────────────────────┐
│ ⚠ Resignation date: Jan 30, 2024          │
│   14 days remaining                        │
└────────────────────────────────────────────┘
```

---

## Employee Transfer

### Page Layout

**Header**:
```
┌─────────────────────────────────────────────────────────────┐
│ ← Employee Transfer                                [Actions] │
│ HR-EMP-TRN-2024-00001                                        │
└─────────────────────────────────────────────────────────────┘
```

**Form Layout**:
```
┌───────────────────────┬───────────────────────┐
│ Employee *            │ Company               │
│ [Select: Active only] │ [Read-only]           │
│                       │                       │
│ Employee Name         │ New Company           │
│ [Read-only]           │ [Select - optional]   │
│                       │                       │
│ Transfer Date *       │ Department            │
│ [Date picker]         │ [Read-only, Bold]     │
└───────────────────────┴───────────────────────┘

┌─── Employee Transfer Details ─────────────────┐
│                                               │
│ [+ Add Employee Property] (button)            │
│                                               │
│ Transfer Details:                             │
│ ┌─────────────┬─────────────┬─────────────┐  │
│ │ Property    │ Current     │ New         │  │
│ ├─────────────┼─────────────┼─────────────┤  │
│ │ Department  │ Engineering │ Sales       │  │
│ │ Designation │ Developer   │ Manager     │  │
│ │ Branch      │ Mumbai      │ Delhi       │  │
│ └─────────────┴─────────────┴─────────────┘  │
│                                               │
│ ☐ Create New Employee ID                     │
│                                               │
│ New Employee ID: [Shows after submit]        │
└───────────────────────────────────────────────┘
```

### Add Property Dialog

**Modal Design**:
```
┌─────── Add Property ──────────────────────────┐
│                                               │
│ Property *                                    │
│ [Autocomplete dropdown]                       │
│ ▼ Department, Designation, Branch...          │
│                                               │
│ Current Value                                 │
│ [Read-only, Auto-filled]                      │
│ Engineering                                   │
│                                               │
│ New Value *                                   │
│ [Dynamic field based on property type]        │
│ [Dropdown for Department: Sales ▼]            │
│                                               │
│               [Cancel]  [Add Property]        │
└───────────────────────────────────────────────┘
```

**Property Selection** (autocomplete):
```
Type to search...
─────────────────────
Department
Designation
Branch
Grade
Employment Type
Reports To
─────────────────────
```

### Visual States

**New Employee ID Indicator**:
```
After submit with "Create New Employee ID" checked:

┌────────────────────────────────────────────┐
│ ✓ New Employee Created                     │
│   HR-EMP-00245                             │
│   [View Employee →]                        │
└────────────────────────────────────────────┘
```

---

## Employee Promotion

### Page Layout

Similar to transfer with CTC fields added.

**Form Layout**:
```
┌───────────────────────┬───────────────────────┐
│ Employee *            │ Promotion Date *      │
│ [Select: Active only] │ [Date picker]         │
│                       │                       │
│ Employee Name         │ Company               │
│ [Read-only]           │ [Read-only]           │
│                       │                       │
│ Department            │                       │
│ [Read-only]           │                       │
│                       │                       │
│ Salary Currency       │                       │
│ [Read-only]           │                       │
└───────────────────────┴───────────────────────┘

┌─── Employee Promotion Details ────────────────┐
│ Set the properties that should be updated in  │
│ the Employee master on promotion submission   │
│                                               │
│ [+ Add Employee Property]                     │
│                                               │
│ [Promotion Details Table]                     │
└───────────────────────────────────────────────┘

┌─── Salary Details ────────────────────────────┐
│ ┌─────────────────┬─────────────────────────┐│
│ │ Current CTC     │ Revised CTC             ││
│ │ ₹ 800,000.00    │ ₹ 1,000,000.00          ││
│ │                 │ [Currency input]        ││
│ └─────────────────┴─────────────────────────┘│
│                                               │
│ Increase: ₹ 200,000.00 (25%) ↑               │
└───────────────────────────────────────────────┘
```

### Visual Elements

**Promotion Indicator**:
```
┌────────────────────────────────────────────┐
│ ↗ Promotion                                │
│   Developer → Senior Developer              │
│   +25% CTC                                 │
└────────────────────────────────────────────┘
```

---

## Exit Interview

### Page Layout

**Form Layout**:
```
┌───────────────────────┬───────────────────────┐
│ Naming Series         │ Company *             │
│ HR-EXIT-INT-          │ [Select]              │
│                       │                       │
│ Employee *            │ Status *              │
│ [Select]              │ [Pending ▼]           │
│                       │                       │
│ Employee Name         │ Date                  │
│ [Read-only]           │ [Date picker]         │
│                       │  (required if         │
│ Email                 │   Scheduled)          │
│ [Read-only]           │                       │
└───────────────────────┴───────────────────────┘

┌─── Employee Details ──────────────────────────┐
│ (Collapsible section)                         │
│ ┌─────────────────┬─────────────────────────┐│
│ │ Department      │ Date of Joining         ││
│ │ Engineering     │ Jan 1, 2020             ││
│ │                 │                         ││
│ │ Designation     │ Relieving Date          ││
│ │ Developer       │ Jan 31, 2024            ││
│ │                 │                         ││
│ │ Reports To      │                         ││
│ │ John Manager    │                         ││
│ └─────────────────┴─────────────────────────┘│
└───────────────────────────────────────────────┘

┌─── Exit Questionnaire ────────────────────────┐
│ Reference Document Type                       │
│ [Select]                                      │
│                                               │
│ Reference Document Name                       │
│ [Dynamic Link]                                │
│                                               │
│ ☑ Questionnaire Email Sent                   │
│                                               │
│ [Send Exit Questionnaire] (button)           │
└───────────────────────────────────────────────┘

┌─── Interview Details ─────────────────────────┐
│ Interviewers *                                │
│ (required if status = Scheduled)              │
│ ┌───────────────────────────────────────────┐│
│ │ ☑ Jane Smith                              ││
│ │ ☑ John Manager                            ││
│ │ [+ Add Interviewer]                       ││
│ └───────────────────────────────────────────┘│
│                                               │
│ Interview Summary                             │
│ [Rich text editor]                            │
└───────────────────────────────────────────────┘

┌─── Final Decision ────────────────────────────┐
│ Employee Status * (required if Completed)     │
│ ○ Employee Retained                           │
│ ● Exit Confirmed                              │
└───────────────────────────────────────────────┘
```

### Status Progression Visual

```
Pending → Scheduled → Completed
  ○         ●           ○

[Status bar showing current state]
```

### Email Sent Indicator

```
┌────────────────────────────────────────────┐
│ ✉ Questionnaire Sent                       │
│   Sent to: john.doe@company.com            │
│   On: Jan 15, 2024 10:30 AM                │
└────────────────────────────────────────────┘
```

---

## Appointment Letter

### Page Layout

```
┌───────────────────────┬───────────────────────┐
│ Job Applicant *       │ Company *             │
│ [Select]              │ [Select]              │
│                       │                       │
│ Applicant Name        │ Appointment Date *    │
│ [Read-only]           │ [Date picker]         │
│                       │                       │
│                       │ Appointment Letter    │
│                       │ Template *            │
│                       │ [Select]              │
└───────────────────────┴───────────────────────┘

┌─── Body ──────────────────────────────────────┐
│ Introduction *                                │
│ [Long text area, auto-filled from template]  │
│ Dear [Applicant Name],                        │
│ We are pleased to...                          │
│                                               │
│                                               │
│ Terms *                                       │
│ ┌─────────────────┬───────────────────────┐  │
│ │ Offer Term      │ Value                 │  │
│ ├─────────────────┼───────────────────────┤  │
│ │ Position        │ Software Engineer     │  │
│ │ Start Date      │ February 1, 2024      │  │
│ │ Salary          │ ₹1,000,000 per annum  │  │
│ │ Probation       │ 6 months              │  │
│ └─────────────────┴───────────────────────┘  │
│                                               │
│                                               │
│ Closing Notes                                 │
│ [Text area]                                   │
│ We look forward to...                         │
└───────────────────────────────────────────────┘

┌─── Printing Details ──────────────────────────┐
│ (Collapsible)                                 │
│ ┌─────────────────┬───────────────────────┐  │
│ │ Letter Head     │ Print Heading         │  │
│ │ [Select]        │ [Select]              │  │
│ └─────────────────┴───────────────────────┘  │
└───────────────────────────────────────────────┘
```

### Print Preview

```
[Company Letterhead]

                    APPOINTMENT LETTER

Date: January 15, 2024

John Doe
123 Street
City, State

Dear John,

[Introduction text...]

Terms and Conditions:
1. Position: Software Engineer
2. Start Date: February 1, 2024
3. Salary: ₹1,000,000 per annum
...

[Closing notes...]

Sincerely,
HR Manager
Company Name
```

---

## Responsive Design

### Mobile View (< 768px)

- Single column layout
- Collapsible sections
- Bottom action buttons
- Simplified tables (cards instead of rows)

**Mobile Form Example**:
```
┌─────────────────────────────────┐
│ ← Employee Onboarding           │
│                                 │
│ ═══════════════════════════════ │
│                                 │
│ Job Applicant *                 │
│ [Select ▼]                      │
│                                 │
│ Employee Name                   │
│ [Read-only]                     │
│                                 │
│ Date of Joining *               │
│ [📅 Select date]                │
│                                 │
│ ▼ Activities (3)                │
│                                 │
│ [Save] [Submit]                 │
└─────────────────────────────────┘
```

---

## Accessibility

### ARIA Labels
- Form fields have descriptive labels
- Status indicators have aria-labels
- Buttons have clear action descriptions

### Keyboard Navigation
- Tab order follows logical flow
- Enter to submit forms
- Escape to close dialogs
- Arrow keys in dropdowns

### Color Contrast
- All text meets WCAG AA standards
- Minimum 4.5:1 contrast ratio
- Color is not the only indicator (icons + text)

---

## Animation & Transitions

**Form Transitions**:
- Fade in: 200ms
- Slide down (sections): 300ms
- Button hover: 150ms

**Status Changes**:
- Smooth color transition: 400ms
- Badge pulse on update: 600ms

**Loading States**:
- Spinner for async operations
- Progress bar for batch operations
- Skeleton screens for data loading

---

This designer documentation provides comprehensive visual specifications for implementing the Employee Lifecycle module UI.
