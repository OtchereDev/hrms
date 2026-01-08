# Attendance & Shifts Module - Designer Documentation

## Overview

UI/UX specifications, visual design guidelines, and interaction patterns for attendance tracking and shift management.

---

## Design System

### Color Palette

**Attendance Status Colors**:
- Present: `#2ECC71` (Green)
- Absent: `#E74C3C` (Red)
- On Leave: `#3498DB` (Blue)
- Half Day: `#FFD43B` (Yellow)
- Work From Home: `#9B59B6` (Purple)

**Shift Status Colors**:
- Active: `#2ECC71` (Green)
- Inactive: `#95A5A6` (Gray)

**Check-in/Out**:
- Check-in: `#2ECC71` (Green)
- Check-out: `#E74C3C` (Red)

---

## Attendance Form

### Desktop Layout

```
┌─── Attendance ────────────────────────────────┐
│ HR-ATT-2024-00001                   [Actions] │
│ ● Present                                     │
└───────────────────────────────────────────────┘

┌─── Attendance Details ────────────────────────┐
│ ┌──────────────────┬───────────────────────┐ │
│ │ Naming Series    │ Status *              │ │
│ │ HR-ATT-.YYYY.-   │ [●  Present      ▼]   │ │
│ │                  │                       │ │
│ │ Employee *       │ Company               │ │
│ │ John Doe ▼       │ Acme Corp             │ │
│ │                  │                       │ │
│ │ Employee Name    │ Department            │ │
│ │ John Doe         │ Engineering           │ │
│ │                  │                       │ │
│ │ Attendance Date* │ Shift                 │ │
│ │ 📅 Jan 15, 2024  │ 8-Hour Shift ▼        │ │
│ │                  │                       │ │
│ │ Working Hours    │                       │ │
│ │ 8.5 hours        │                       │ │
│ └──────────────────┴───────────────────────┘ │
└───────────────────────────────────────────────┘

┌─── Timing Details (if shift) ─────────────────┐
│ ┌──────────────────┬───────────────────────┐ │
│ │ In Time          │ Out Time              │ │
│ │ 🕐 09:05 AM      │ 🕐 06:30 PM           │ │
│ │                  │                       │ │
│ │ ☑ Late Entry     │ ☐ Early Exit          │ │
│ │ 5 mins late      │                       │ │
│ └──────────────────┴───────────────────────┘ │
└───────────────────────────────────────────────┘
```

**Status Dropdown** (with icons):
```
┌──────────────────────────────┐
│ ● Present                    │
│ ● Absent                     │
│ ● On Leave                   │
│ ◐ Half Day                   │
│ 🏠 Work From Home             │
└──────────────────────────────┘
```

### Mobile Attendance Entry

```
┌────────────────────────────┐
│ ← Mark Attendance      ✓   │
├────────────────────────────┤
│                            │
│ 📅 Date                    │
│ Jan 15, 2024               │
│                            │
│ Status *                   │
│ ● Present                  │
│                            │
│ ┌────────────────────────┐│
│ │ ● Present              ││
│ │ ● Absent               ││
│ │ ● On Leave             ││
│ │ ◐ Half Day             ││
│ │ 🏠 Work From Home       ││
│ └────────────────────────┘│
│                            │
│ [Submit]                   │
└────────────────────────────┘
```

---

## Check-in Interface

### Mobile Quick Check-in

**Floating Action Button**:
```
┌────────────────────────────┐
│                            │
│   [Attendance Dashboard]   │
│                            │
│                            │
│                            │
│                         ┌─┐│
│                         │●││ ← FAB
│                         └─┘│
│                       Check │
│                         In  │
└────────────────────────────┘
```

**Button States**:

**Checked Out** (Ready to check in):
```
┌─────────────────┐
│   🟢 Check In   │  ← Green, pulse animation
└─────────────────┘
```

**Checked In** (Ready to check out):
```
┌─────────────────┐
│   🔴 Check Out  │  ← Red, steady
└─────────────────┘
```

### Check-in Confirmation

```
┌─────────────────────────────┐
│   ✓ Checked In              │
│                             │
│   🕐 09:05 AM               │
│   📍 Office Building A      │
│                             │
│   Shift: 8-Hour Shift       │
│   Status: On Time           │
│                             │
│   [View Details]            │
└─────────────────────────────┘
```

**Late Check-in Warning**:
```
┌─────────────────────────────┐
│   ⚠ Late Check-in           │
│                             │
│   Expected: 09:00 AM        │
│   Actual:   09:15 AM        │
│   Late by:  15 minutes      │
│                             │
│   [Check In Anyway]         │
│   [Cancel]                  │
└─────────────────────────────┘
```

### Location Map View

```
┌─────────────────────────────┐
│   Check-in Location         │
├─────────────────────────────┤
│                             │
│   [Interactive Map]         │
│   📍 Your Location          │
│   ⭕ Office Boundary        │
│      (200m radius)          │
│                             │
│   Distance: 45m             │
│   ✓ Within range            │
│                             │
│   [Confirm Location]        │
└─────────────────────────────┘
```

---

## Attendance Calendar

### Monthly Calendar View

```
┌────── January 2024 ────────┐
│         <  Today  >        │
├────────────────────────────┤
│ S  M  T  W  T  F  S        │
│ ─  1  2  3  4  5  6        │
│ 7  🟢 🟢 🟢 🟢 🟢 13       │
│ 14 🟢 🔴 🟢 🟢 🟢 20       │
│ 21 🟢 🔵 🔵 🟢 🟢 27       │
│ 28 🟢 🟢 🟢 ─  ─  ─        │
│                            │
│ Legend:                    │
│ 🟢 Present  🔴 Absent      │
│ 🔵 On Leave 🟡 Half Day   │
│ 🟣 WFH      ⚪ Not Marked  │
└────────────────────────────┘
```

**Calendar Day States**:
- Present: Green filled circle
- Absent: Red filled circle
- On Leave: Blue filled circle
- Half Day: Yellow half-filled circle
- WFH: Purple filled circle
- Not Marked: Gray outlined circle
- Weekend: Light gray background
- Holiday: Red background with text
- Today: Blue border

### Calendar Hover Card

```
┌──────────────────┐
│ Jan 15, 2024     │
├──────────────────┤
│ Status: Present  │
│ In: 09:05 AM     │
│ Out: 06:30 PM    │
│ Hours: 8.5       │
│                  │
│ [View Details →] │
└──────────────────┘
```

---

## Shift Management

### Shift Type Card

```
┌─────────────────────────────────┐
│ 8-Hour Shift            [Edit]  │
├─────────────────────────────────┤
│ 🕐 09:00 AM → 06:00 PM          │
│ Duration: 9 hours               │
│                                 │
│ ☑ Auto Attendance               │
│ ⏰ Grace: 60 mins               │
│ 📍 Geofencing: Enabled          │
│                                 │
│ Active Employees: 45            │
│                                 │
│ [View Assignments →]            │
└─────────────────────────────────┘
```

### Shift Assignment Timeline

```
┌─── Employee: John Doe ──────────┐
│                                 │
│ January 2024                    │
│ ────────────────────────────────│
│ 1-7   │ ████████ 8-Hour Shift  │
│ 8-14  │ ████████ Night Shift   │
│ 15-21 │ ████████ 8-Hour Shift  │
│ 22-31 │ ████████ Weekend Shift │
│                                 │
│ [Edit Schedule]                 │
└─────────────────────────────────┘
```

### Shift Assignment Tool

**Employee Selection Grid**:
```
┌─── Select Employees ─────────────────────────┐
│ ┌───────────────────────────────────────┐   │
│ │ ☑ All (25)                            │   │
│ ├───────────────────────────────────────┤   │
│ │ ☑ John Doe      Engineering    L3     │   │
│ │ ☑ Jane Smith    Sales          L4     │   │
│ │ ☐ Bob Johnson   Marketing      L2     │   │
│ │ ☑ Alice Brown   Engineering    L3     │   │
│ │ ☐ Charlie Davis HR             L2     │   │
│ └───────────────────────────────────────┘   │
│                                              │
│ 3 selected     [Assign Shift to Selected]   │
└──────────────────────────────────────────────┘
```

**Bulk Assignment Progress**:
```
┌─── Assigning Shifts ─────────────────┐
│                                      │
│  Assigning to 25 employees...        │
│                                      │
│  ████████████████░░░░ 80%            │
│                                      │
│  Completed: 20 / 25                  │
│  Failed: 1                           │
│                                      │
│  [View Details]                      │
└──────────────────────────────────────┘
```

---

## Mobile App Design

### Attendance Dashboard

```
┌─────────────────────────────────┐
│ ☰  Attendance              🔔   │
├─────────────────────────────────┤
│                                 │
│ ┌─────────────────────────────┐│
│ │ Today's Status              ││
│ │                             ││
│ │      🟢 Present             ││
│ │                             ││
│ │ In:  09:05 AM               ││
│ │ Out: Not yet checked out    ││
│ │                             ││
│ │ ┌───────────────────────┐  ││
│ │ │   🔴 Check Out        │  ││
│ │ └───────────────────────┘  ││
│ └─────────────────────────────┘│
│                                 │
│ This Week                       │
│ ┌─────────────────────────────┐│
│ │ M T W T F S S               ││
│ │ 🟢🟢🟢🟢🟢⚪⚪              ││
│ └─────────────────────────────┘│
│                                 │
│ Recent Activity                 │
│ ┌─────────────────────────────┐│
│ │ Check-in    09:05 AM  Today ││
│ │ Check-out   06:30 PM  Jan14 ││
│ │ Check-in    09:00 AM  Jan14 ││
│ └─────────────────────────────┘│
│                                 │
│ [View Calendar]                 │
└─────────────────────────────────┘
```

### Shift Schedule View

```
┌─────────────────────────────────┐
│ ← My Shifts                     │
├─────────────────────────────────┤
│ This Week                       │
│ ┌─────────────────────────────┐│
│ │ Monday, Jan 15              ││
│ │ 🌅 8-Hour Shift             ││
│ │ 09:00 AM - 06:00 PM         ││
│ └─────────────────────────────┘│
│ ┌─────────────────────────────┐│
│ │ Tuesday, Jan 16             ││
│ │ 🌅 8-Hour Shift             ││
│ │ 09:00 AM - 06:00 PM         ││
│ └─────────────────────────────┘│
│ ┌─────────────────────────────┐│
│ │ Wednesday, Jan 17           ││
│ │ 🌙 Night Shift              ││
│ │ 09:00 PM - 06:00 AM         ││
│ └─────────────────────────────┘│
│                                 │
│ [Request Shift Change]          │
└─────────────────────────────────┘
```

---

## User Flows

### Check-in Flow

```
User opens app
     ↓
Dashboard shows check-in button
     ↓
User taps "Check In"
     ↓
[Request location permission] (if first time)
     ↓
Get current location
     ↓
[Validate distance]
     ↓
[If out of range]
  → Show error
  → Offer "Check in anyway"
     ↓
[If in range]
  → Show confirmation
  → Location on map
  → Detected shift
     ↓
User confirms
     ↓
Create Employee Checkin
     ↓
[If late]
  → Show warning notification
     ↓
Show success
  → Update dashboard
  → Change button to "Check Out"
```

### Shift Assignment Flow

```
Manager opens Shift Assignment Tool
     ↓
Select action: "Assign Shift"
     ↓
Choose shift type and dates
     ↓
Apply filters (department, designation)
     ↓
Click "Get Employees"
     → System fetches eligible employees
     → Filters out conflicting assignments
     ↓
Employee grid displayed
     ↓
Select employees (checkbox)
     ↓
Click "Assign Shift to Selected"
     ↓
[If > 30 employees]
  → Show progress dialog
  → Queue background job
  → Real-time progress updates
     ↓
[If ≤ 30 employees]
  → Process immediately
     ↓
Show success summary
  → "25 assignments created"
  → "2 failed (view details)"
```

---

## Responsive Design

### Breakpoints

- Mobile: < 768px
- Tablet: 768px - 1024px
- Desktop: > 1024px

### Mobile Adaptations

**Attendance Form**:
- Single column layout
- Large touch targets (48x48px)
- Bottom sheet for status selection
- Native date picker

**Calendar**:
- Swipe to navigate months
- Pinch to zoom
- Tap day for details

**Shift Grid**:
- Horizontal scroll for timeline
- Card view instead of grid

---

## Animations

**Check-in Button**:
```css
.checkin-button {
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% {
    box-shadow: 0 0 0 0 rgba(46, 204, 113, 0.7);
  }
  50% {
    box-shadow: 0 0 0 10px rgba(46, 204, 113, 0);
  }
}
```

**Status Change**:
```css
.status-badge {
  transition: all 300ms ease-in-out;
}

.status-badge.changed {
  animation: pop 400ms cubic-bezier(0.68, -0.55, 0.265, 1.55);
}
```

**Calendar Day Hover**:
```css
.calendar-day:hover {
  transform: scale(1.1);
  transition: transform 200ms ease;
}
```

---

## Empty States

### No Attendance Marked

```
┌──────────────────────────────┐
│         📅                   │
│                              │
│   No Attendance Marked       │
│                              │
│  You haven't marked          │
│  attendance for this week.   │
│                              │
│  [Mark Attendance]           │
└──────────────────────────────┘
```

### No Shift Assigned

```
┌──────────────────────────────┐
│         ⏰                    │
│                              │
│   No Shift Assigned          │
│                              │
│  You don't have any shifts   │
│  assigned. Contact your      │
│  manager.                    │
│                              │
│  [Contact Manager]           │
└──────────────────────────────┘
```

---

## Accessibility

### Color Blindness Support

**Status Indicators** (not color-only):
- Present: Green + ✓ icon
- Absent: Red + ✗ icon
- On Leave: Blue + ✈ icon
- Half Day: Yellow + ◐ icon
- WFH: Purple + 🏠 icon

### Screen Reader

- Status announcements: "Attendance marked as Present"
- Calendar dates: "Monday, January 15, 2024. Present. 8.5 hours worked"
- Check-in: "Checked in at 9:05 AM"

### Keyboard Navigation

- Tab through form fields
- Enter to submit check-in
- Arrow keys to navigate calendar
- Space to toggle checkboxes

---

## Component Library

### Attendance Status Badge

**Sizes**: Small (24px), Medium (32px), Large (48px)

**Variants**:
```html
<badge type="present" size="medium">Present</badge>
<badge type="absent" size="medium">Absent</badge>
<badge type="on-leave" size="medium">On Leave</badge>
```

### Shift Card

**Anatomy**:
- Header: Shift name + action button
- Body: Time range + duration
- Meta: Auto-attendance status, geofencing
- Footer: Employee count + CTA

### Check-in Button

**States**:
- Default (check-in)
- Active (check-out)
- Loading
- Disabled
- Error

---

This designer documentation provides complete visual specifications for the Attendance & Shifts module.
