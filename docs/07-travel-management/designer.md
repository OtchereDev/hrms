# Travel Management Module - Designer Documentation

## Overview

UI/UX specifications for employee travel request management, itinerary planning, and cost estimation.

---

## Travel Request Form

### Layout

```
┌─────────────────────────────────────────────────────────────┐
│  Travel Request                               [Save ▼] [×] │
│  HR-TRQ-2024-00123                                          │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─ Travel Details ─────────────────────────────────────┐  │
│  │                                                       │  │
│  │  ┌──────────────────────┬──────────────────────────┐ │  │
│  │  │ Travel Type *        │ Purpose of Travel *      │ │  │
│  │  │ ● Domestic           │ 🔍 Conference            │ │  │
│  │  │ ○ International      │                          │ │  │
│  │  └──────────────────────┴──────────────────────────┘ │  │
│  │                                                       │  │
│  │  Travel Funding                                       │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ Require Full Funding                        ▼  │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │  Options: Require Full Funding, Fully Sponsored,      │  │
│  │           Partially Sponsored                         │  │
│  │                                                       │  │
│  │  Copy of Invitation/Announcement                      │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ 📎 [Upload File]                                │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  │  Details of Sponsor (Name, Location)                  │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ ABC Corp, New York, USA                         │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─ Employee Details ───────────────────────────────────┐  │
│  │                                                       │  │
│  │  ┌──────────────────────┬──────────────────────────┐ │  │
│  │  │ Employee *           │ Employee Name            │ │  │
│  │  │ ┌──────────────────┐ │ John Doe                 │ │  │
│  │  │ │ 🔍 Search...     │ │ [Auto-filled]            │ │  │
│  │  │ └──────────────────┘ │                          │ │  │
│  │  └──────────────────────┴──────────────────────────┘ │  │
│  │                                                       │  │
│  │  ┌──────────────────────┬──────────────────────────┐ │  │
│  │  │ Cell Number          │ Preferred Email          │ │  │
│  │  │ +1-555-0123          │ john@email.com           │ │  │
│  │  └──────────────────────┴──────────────────────────┘ │  │
│  │                                                       │  │
│  │  ┌──────────────────────┬──────────────────────────┐ │  │
│  │  │ Company *            │ Date of Birth            │ │  │
│  │  │ Acme Corp            │ 📅 Jan 15, 1990          │ │  │
│  │  └──────────────────────┴──────────────────────────┘ │  │
│  │                                                       │  │
│  │  ┌──────────────────────┬──────────────────────────┐ │  │
│  │  │ Personal ID Type     │ Personal ID Number       │ │  │
│  │  │ Driver's License     │ DL123456789              │ │  │
│  │  │                      │                          │ │  │
│  │  │ Passport Number      │                          │ │  │
│  │  │ P987654321           │                          │ │  │
│  │  └──────────────────────┴──────────────────────────┘ │  │
│  │                                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─ Travel Itinerary ───────────────────────────────────┐  │
│  │                                                       │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ From  │ To   │ Date   │ Arrival│ Hotel│ Mode   │ │  │
│  │  ├─────────────────────────────────────────────────┤ │  │
│  │  │ NYC   │ LA   │ Mar 10 │ Mar 10 │ ☑   │ ✈ Air │ │  │
│  │  │       │      │        │        │ 4-Star│       │ │  │
│  │  │       │      │        │        │       │       │ │  │
│  │  │ LA    │ SF   │ Mar 15 │ Mar 15 │ ☑   │ 🚗 Car│ │  │
│  │  │       │      │        │        │ Budget│       │ │  │
│  │  │       │      │        │        │       │       │ │  │
│  │  │ SF    │ NYC  │ Mar 20 │ Mar 20 │ ☐   │ ✈ Air │ │  │
│  │  │ [+ Add Row]                                     │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─ Costing Details ────────────────────────────────────┐  │
│  │                                                       │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ Item              │ Expense Type │ Estimated Cost││ │
│  │  ├─────────────────────────────────────────────────┤ │  │
│  │  │ Flight Tickets    │ Travel       │    $850.00    ││ │
│  │  │ Hotel (5 nights)  │ Lodging      │    $600.00    ││ │
│  │  │ Car Rental        │ Travel       │    $200.00    ││ │
│  │  │ Meals             │ Food         │    $150.00    ││ │
│  │  │ [+ Add Row]                                     │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ Total Estimated Cost:             $1,800.00     │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─ Event Details ──────────────────────────────────────┐  │
│  │  [Collapsible]                                        │  │
│  │                                                       │  │
│  │  Name of Organizer                                    │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ Global Tech Conference                          │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  │  Address of Organizer                                 │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ 123 Conference Center, Los Angeles, CA 90001    │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─ Description ─────────────────────────────────────────┐  │
│  │  [Collapsible]                                        │  │
│  │                                                       │  │
│  │  Any Other Details                                    │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ Attending annual tech conference to present     │ │  │
│  │  │ research paper on AI innovations...             │ │  │
│  │  │                                                 │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─ Accounting Dimensions ──────────────────────────────┐  │
│  │                                                       │  │
│  │  ┌──────────────────────┬──────────────────────────┐ │  │
│  │  │ Cost Center          │ Project                  │ │  │
│  │  │ [Optional]           │ [Optional]               │ │  │
│  │  └──────────────────────┴──────────────────────────┘ │  │
│  │                                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Design Elements

**Travel Type Radio Buttons**:
```
● Domestic    ○ International
```
- Selected: Blue filled circle
- Unselected: Gray outline circle

**Itinerary Icons**:
- ✈ Air (Plane icon)
- 🚗 Car (Car icon)
- 🚂 Train (Train icon)
- 🚌 Bus (Bus icon)

**Lodging Checkbox**:
```
☑ Lodging Required
Hotel Type: [4-Star ▼]
```

### List View

```
┌─────────────────────────────────────────────────────────────┐
│  Travel Requests                          [+ New] [⚙] [⟳]  │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Filters:  [Travel Type ▼] [Purpose ▼] [Date Range]        │
│                                                             │
├─────┬────────────┬──────────┬────────────┬────────┬────────┤
│ ●   │ ID         │ Employee │ Purpose    │ Type   │ Date   │
├─────┼────────────┼──────────┼────────────┼────────┼────────┤
│ ☑   │ TRQ-001    │ John Doe │ Conference │ Intl   │ Mar 10 │
│ ☑   │ TRQ-002    │ Jane S   │ Training   │ Dom    │ Mar 15 │
│ ☑   │ TRQ-003    │ Bob J    │ Meeting    │ Dom    │ Mar 20 │
└─────┴────────────┴──────────┴────────────┴────────┴────────┘
```

### Mobile View (< 768px)

```
┌───────────────────────────────┐
│  ☰  Travel Request       [×] │
│  HR-TRQ-2024-00123            │
├───────────────────────────────┤
│                               │
│  ┌─ Summary ─────────────┐   │
│  │ Type: International   │   │
│  │ Purpose: Conference   │   │
│  │ Employee: John Doe    │   │
│  │                       │   │
│  │ Travel: Mar 10-20     │   │
│  │ Est Cost: $1,800      │   │
│  └───────────────────────┘   │
│                               │
│  [View Itinerary] ▼           │
│  [View Costing] ▼             │
│                               │
└───────────────────────────────┘
```

---

## Color Palette

**Travel Types**:
- Domestic: #2490EF (Blue)
- International: #9B59B6 (Purple)

**Modes of Travel**:
- Air: #3498DB (Sky Blue)
- Train: #E74C3C (Red)
- Car: #2ECC71 (Green)
- Bus: #F39C12 (Orange)

---

This designer documentation provides complete UI/UX specifications for the Travel Management module.
