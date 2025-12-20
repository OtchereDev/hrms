# Performance Management Module - Designer Documentation

## Overview

UI/UX specifications, visual design patterns, and user interface guidelines for the complete performance management system including appraisal cycles, goals, performance reviews, and feedback.

---

## Design System

### Color Palette

**Status Colors**:
```
Pending:          #FFA00A (Orange)
In Progress:      #2490EF (Blue)
Completed:        #98D85B (Green)
Archived:         #9B9B9B (Gray)
Closed:           #4A5568 (Dark Gray)
Not Started:      #FFD66B (Yellow)
```

**Score Indicators**:
```
Excellent (4.5-5.0):  #98D85B (Green)
Good (3.5-4.4):       #2490EF (Blue)
Average (2.5-3.4):    #FFA00A (Orange)
Poor (< 2.5):         #F56B6B (Red)
```

### Typography

**Headers**: 18px, Semi-bold, #2C3E50
**Section Titles**: 14px, Semi-bold, #4A5568
**Labels**: 12px, Regular, #718096
**Values**: 14px, Regular, #2D3748
**Scores**: 24px, Bold, color-coded by performance

### Spacing

**Section Gaps**: 24px
**Field Gaps**: 16px
**Card Padding**: 20px
**Progress Bar Height**: 8px
**Rating Star Size**: 20px

---

## Appraisal Cycle

### Form Layout

```
┌─────────────────────────────────────────────────────────────┐
│  Appraisal Cycle                              [Save ▼] [×] │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─ Cycle Details ──────────────────────────────────────┐  │
│  │                                                       │  │
│  │  ┌──────────────────────┬──────────────────────────┐ │  │
│  │  │ Cycle Name *         │ Status *                 │ │  │
│  │  │ Q1 2024 Performance  │ ┌──────────────────────┐ │ │  │
│  │  │ Review               │ │ In Progress      ▼  │ │ │  │
│  │  │                      │ └──────────────────────┘ │ │  │
│  │  │                      │ ● In Progress            │ │  │
│  │  └──────────────────────┴──────────────────────────┘ │  │
│  │                                                       │  │
│  │  ┌──────────────────────┬──────────────────────────┐ │  │
│  │  │ Company *            │ KRA Evaluation Method *  │ │  │
│  │  │ ┌──────────────────┐ │ ● Automatic (Goal-based)│ │  │
│  │  │ │ Acme Corp    ▼  │ │ ○ Manual Rating         │ │  │
│  │  │ └──────────────────┘ │                         │ │  │
│  │  └──────────────────────┴──────────────────────────┘ │  │
│  │                                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─ Timeline ───────────────────────────────────────────┐  │
│  │                                                       │  │
│  │  ┌──────────────────────┬──────────────────────────┐ │  │
│  │  │ Start Date *         │ End Date *               │ │  │
│  │  │ 📅 Jan 1, 2024       │ 📅 Mar 31, 2024          │ │  │
│  │  └──────────────────────┴──────────────────────────┘ │  │
│  │                                                       │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ ⏱ 90 days │ Started 30 days ago │ 60 days left  │ │  │
│  │  │ ████████████░░░░░░░░░░░░░░░░░░░░░  33%          │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─ Employee Selection Filters ─────────────────────────┐  │
│  │                                                       │  │
│  │  ┌──────────────┬──────────────┬──────────────────┐  │  │
│  │  │ Department   │ Branch       │ Designation      │  │  │
│  │  │ [Optional ▼] │ [Optional ▼] │ [Optional ▼]     │  │  │
│  │  └──────────────┴──────────────┴──────────────────┘  │  │
│  │                                                       │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │            [Get Employees]                      │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─ Appraisees ─────────────────────────────────────────┐  │
│  │                                                       │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ Employee  │ Designation│ Department │ Template  │ │  │
│  │  ├─────────────────────────────────────────────────┤ │  │
│  │  │ John Doe  │ Engineer  │ Engineering│ Standard  │ │  │
│  │  │ Jane Smith│ Manager   │ Sales      │ Manager   │ │  │
│  │  │ Bob Jones │ Designer  │ Design     │ Creative  │ │  │
│  │  │ [Auto-populated after clicking Get Employees]   │ │  │
│  │  │ [+ Add Row]                                     │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  │  ⓘ 45 employees selected                             │  │
│  │                                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─ Final Score Configuration ──────────────────────────┐  │
│  │                                                       │  │
│  │  ☑ Calculate Final Score Based on Formula            │  │
│  │                                                       │  │
│  │  Final Score Formula                                  │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ (goal_score * 0.5) +                            │ │  │
│  │  │ (average_feedback_score * 0.3) +                │ │  │
│  │  │ (self_appraisal_score * 0.2)                    │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  │  ⓘ Available variables:                               │  │
│  │  • goal_score                                         │  │
│  │  • average_feedback_score                             │  │
│  │  • self_appraisal_score                               │  │
│  │  • All fields from Employee, Appraisal Cycle, Appraisal│ │
│  │                                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Dashboard Indicators

**Cycle Progress Summary**:
```
┌─────────────────────────────────────────────────────────────┐
│  ⓘ Appraisees: 45                                          │
│  ⓘ Self Appraisal Pending: 12                              │
│  ⓘ Employees without Feedback: 8                           │
│  ⓘ Employees without Goals: 5                              │
└─────────────────────────────────────────────────────────────┘
```
- Background: Light blue (#EBF8FF)
- Icons: Blue (#2490EF)
- Text: Dark gray (#2D3748)

### Action Buttons

**Primary Actions** (varies by status):

Not Started + No Appraisals:
```
┌─────────────────────────────────┐
│  [🎯 Create Appraisals]         │
└─────────────────────────────────┘
```

In Progress:
```
┌─────────────────────────────────┐
│  [✓ Mark as Completed]          │
└─────────────────────────────────┘
```

**Secondary Actions**:
```
┌─────────────────────────────────┐
│  [📊 View Goals]                │
│  [📝 Create Appraisals]         │
└─────────────────────────────────┘
```

---

## Appraisal

### Form Layout

```
┌─────────────────────────────────────────────────────────────┐
│  Appraisal - John Doe                         [Submit ▼] [×]│
│  [Employee Photo]                                           │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─ Appraisal Details ──────────────────────────────────┐  │
│  │                                                       │  │
│  │  ┌──────────────────────┬──────────────────────────┐ │  │
│  │  │ Employee *           │ Employee Name            │ │  │
│  │  │ ┌──────────────────┐ │ John Doe                 │ │  │
│  │  │ │ 🔍 Search...     │ │ [Read-only]              │ │  │
│  │  │ └──────────────────┘ │                          │ │  │
│  │  └──────────────────────┴──────────────────────────┘ │  │
│  │                                                       │  │
│  │  ┌──────────────────────┬──────────────────────────┐ │  │
│  │  │ Appraisal Cycle *    │ Company *                │ │  │
│  │  │ Q1 2024 Review       │ Acme Corp                │ │  │
│  │  └──────────────────────┴──────────────────────────┘ │  │
│  │                                                       │  │
│  │  ┌──────────────────────┬──────────────────────────┐ │  │
│  │  │ Start Date           │ End Date                 │ │  │
│  │  │ 📅 Jan 1, 2024       │ 📅 Mar 31, 2024          │ │  │
│  │  └──────────────────────┴──────────────────────────┘ │  │
│  │                                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─ Template Configuration ─────────────────────────────┐  │
│  │                                                       │  │
│  │  Appraisal Template *                                 │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ Software Engineer Template                  ▼  │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │  [Auto-populated from cycle]                          │  │
│  │                                                       │  │
│  │  ☐ Rate Goals Manually                                │  │
│  │  ⓘ Check this to manually rate goals instead of      │  │
│  │    auto-calculating from goal progress                │  │
│  │                                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─ Goals/KRAs ─────────────────────────────────────────┐  │
│  │  [Auto-Calculated Mode - shown when unchecked]        │  │
│  │                                                       │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ Key Result Area│ Weight│ Completion│ Score      │ │  │
│  │  ├─────────────────────────────────────────────────┤ │  │
│  │  │ Delivery       │  40%  │           │            │ │  │
│  │  │                │       │ ████████  │            │ │  │
│  │  │                │       │   75%     │   30.00    │ │  │
│  │  ├─────────────────────────────────────────────────┤ │  │
│  │  │ Quality        │  30%  │ ████████  │   27.00    │ │  │
│  │  ├─────────────────────────────────────────────────┤ │  │
│  │  │ Learning       │  30%  │ ██████    │   18.00    │ │  │
│  │  │                │       │   60%     │            │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ Goal Score Percentage:      75%                 │ │  │
│  │  │ Total Score:             3.75 / 5.00            │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  │  [Manual Rating Mode - shown when checked]            │  │
│  │                                                       │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ KRA         │ Weight│ Rating    │ Score Earned  │ │  │
│  │  ├─────────────────────────────────────────────────┤ │  │
│  │  │ Delivery    │  40%  │ ★★★★☆     │    3.20       │ │  │
│  │  │ Quality     │  30%  │ ★★★★★     │    3.00       │ │  │
│  │  │ Learning    │  30%  │ ★★★☆☆     │    1.80       │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ Total Score:             4.00 / 5.00            │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─ Self Appraisal ─────────────────────────────────────┐  │
│  │                                                       │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ Criteria        │ Weight│ My Rating             │ │  │
│  │  ├─────────────────────────────────────────────────┤ │  │
│  │  │ Communication   │  30%  │ ★★★★☆                 │ │  │
│  │  │ Team Work       │  30%  │ ★★★★★                 │ │  │
│  │  │ Initiative      │  40%  │ ★★★★☆                 │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ Self Score:              4.30 / 5.00            │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─ Performance Summary ────────────────────────────────┐  │
│  │                                                       │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │                                                 │ │  │
│  │  │  Goal Score:             3.75 / 5.00            │ │  │
│  │  │  Self Score:             4.30 / 5.00            │ │  │
│  │  │  Feedback Score:         4.10 / 5.00            │ │  │
│  │  │                                                 │ │  │
│  │  │  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │ │  │
│  │  │                                                 │ │  │
│  │  │  FINAL SCORE           4.05 / 5.00              │ │  │
│  │  │  ★★★★☆                                         │ │  │
│  │  │                                                 │ │  │
│  │  │  ⓘ Excellent Performance                       │ │  │
│  │  │                                                 │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─ Feedback History ───────────────────────────────────┐  │
│  │  [Dynamic HTML Section]                               │  │
│  │                                                       │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ Average Feedback: ★★★★☆ 4.1/5                  │ │  │
│  │  │                                                 │ │  │
│  │  │ Rating Distribution:                            │ │  │
│  │  │ 5 ★ ████████████████████ 60%                    │ │  │
│  │  │ 4 ★ ████████████ 40%                            │ │  │
│  │  │ 3 ★          0%                                 │ │  │
│  │  │ 2 ★          0%                                 │ │  │
│  │  │ 1 ★          0%                                 │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ Feedback from: Sarah Manager (Manager)          │ │  │
│  │  │ Date: Jan 28, 2024                              │ │  │
│  │  ├─────────────────────────────────────────────────┤ │  │
│  │  │ Overall Rating: ★★★★★ 5.0/5                    │ │  │
│  │  │                                                 │ │  │
│  │  │ "John has shown exceptional performance this    │ │  │
│  │  │ quarter. His delivery quality and teamwork      │ │  │
│  │  │ have been outstanding..."                       │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ Feedback from: Bob Tech Lead (Tech Lead)        │ │  │
│  │  │ Date: Jan 25, 2024                              │ │  │
│  │  │ Overall Rating: ★★★★☆ 4.2/5                    │ │  │
│  │  │ ...                                             │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │            [+ Add Feedback]                     │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
└─────────────────────────────────────────────────────────────┘

[📊 View Goals]
```

### Performance Chart

**KRA Score Visualization**:
```
┌─── Goal vs Maximum Scores ───────────────────────────────┐
│                                                          │
│  100 │                                                   │
│      │  ████                                             │
│   75 │  ████     ████                                    │
│      │  ████     ████     ████                           │
│   50 │  ████     ████     ████  ████                     │
│      │  ████     ████     ████  ████                     │
│   25 │  ████ ▓▓▓ ████ ▓▓▓ ████  ████ ▓▓▓                │
│      │  ████ ▓▓▓ ████ ▓▓▓ ████  ████ ▓▓▓                │
│    0 └──────────────────────────────────────────────────│
│       Delivery Quality Learning Team Work Communication │
│                                                          │
│  ████ Maximum Score (Weightage)   ▓▓▓ Score Obtained    │
│                                                          │
└──────────────────────────────────────────────────────────┘
```
- Bar height represents percentage
- Blue bars: Maximum possible score (weightage)
- Green bars: Actual score obtained
- Hover shows exact values

### Add Feedback Dialog

```
┌──────────────────────────────────────────────┐
│  Add Feedback                           [×] │
├──────────────────────────────────────────────┤
│                                              │
│  Feedback *                                  │
│  ┌────────────────────────────────────────┐ │
│  │ B I U ≡ • ⊕ ⊞ 🔗 ⚙ @mention          │ │
│  ├────────────────────────────────────────┤ │
│  │                                        │ │
│  │ John has consistently demonstrated... │ │
│  │                                        │ │
│  │                                        │ │
│  │                                        │ │
│  └────────────────────────────────────────┘ │
│                                              │
│  Feedback Ratings *                          │
│  ┌────────────────────────────────────────┐ │
│  │ Criteria        │ Weight │ Rating      │ │
│  ├────────────────────────────────────────┤ │
│  │ Communication   │  30%   │ ★★★★★       │ │
│  │ Team Work       │  30%   │ ★★★★☆       │ │
│  │ Initiative      │  40%   │ ★★★★★       │ │
│  └────────────────────────────────────────┘ │
│  ⓘ Cannot add/remove rows - from template   │
│                                              │
│                                              │
│                        [Cancel] [Submit]     │
└──────────────────────────────────────────────┘
```

---

## Goal Management

### Tree View

**Route**: `/app/goal/view/tree`

```
┌─────────────────────────────────────────────────────────────┐
│  Goal Tree                                  [+ New] [⚙] [⟳]│
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Filters:  [Employee ▼] [Appraisal Cycle ▼] [Date Range]   │
│                                                             │
│  ┌───────────────────────────────────────────────────────┐ │
│  │                                                       │ │
│  │  📊 All Goals                                         │ │
│  │  │                                                    │ │
│  │  ├─ 📁 Delivery Excellence                           │ │
│  │  │  ├─ Progress: 60%  ████████████░░░░░░░░░░        │ │
│  │  │  ├─ 2 of 3 Completed                              │ │
│  │  │  │                                                 │ │
│  │  │  ├─ 📄 Complete Project Alpha                     │ │
│  │  │  │  │  100%  ████████████████████████  ✓ Completed│ │
│  │  │  │                                                 │ │
│  │  │  ├─ 📄 Reduce Bugs by 20%                         │ │
│  │  │  │  │   45%  █████████░░░░░░░░░░░  ⏸ In Progress │ │
│  │  │  │                                                 │ │
│  │  │  └─ 📄 Improve Response Time                      │ │
│  │  │     │   35%  ███████░░░░░░░░░░░░░  ⏸ In Progress │ │
│  │  │                                                    │ │
│  │  ├─ 📁 Team Collaboration                            │ │
│  │  │  ├─ Progress: 80%  ████████████████░░░░          │ │
│  │  │  ├─ 1 of 3 Completed                              │ │
│  │  │  │                                                 │ │
│  │  │  ├─ 📄 Mentor 2 juniors                           │ │
│  │  │  │  │  100%  ████████████████████████  ✓ Completed│ │
│  │  │  │                                                 │ │
│  │  │  ├─ 📄 Lead weekly standup                        │ │
│  │  │  │  │   90%  ██████████████████░░  ⏸ In Progress │ │
│  │  │  │                                                 │ │
│  │  │  └─ 📄 Cross-team documentation                   │ │
│  │  │     │   50%  ██████████░░░░░░░░░░  ⏸ In Progress │ │
│  │  │                                                    │ │
│  │  └─ 📁 Learning & Development                        │ │
│  │     ├─ Progress: 40%  ████████░░░░░░░░░░            │ │
│  │     ├─ 0 of 2 Completed                              │ │
│  │     │                                                 │ │
│  │     ├─ 📄 Complete AWS certification                 │ │
│  │     │  │    0%  ░░░░░░░░░░░░░░░░░░░░  ⊙ Pending    │ │
│  │     │                                                 │ │
│  │     └─ 📄 Read 3 tech books                          │ │
│  │        │   80%  ████████████████░░░░  ⏸ In Progress │ │
│  │                                                       │ │
│  └───────────────────────────────────────────────────────┘ │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Tree Node Design**:

Group Goal (Parent):
```
📁 Delivery Excellence
├─ Progress: 60%  ████████████░░░░░░░░░░
├─ 2 of 3 Completed
└─ [Click to expand/collapse]
```

Individual Goal (Leaf):
```
📄 Complete Project Alpha
└─ 100%  ████████████████████████  ✓ Completed
   [Click to edit progress]
```

**Status Icons**:
- ⊙ Pending (Orange circle)
- ⏸ In Progress (Blue pause icon)
- ✓ Completed (Green checkmark)
- ⊗ Closed (Gray X)
- 📦 Archived (Gray box)

### Form View

```
┌─────────────────────────────────────────────────────────────┐
│  Goal                                         [Save ▼] [×] │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─ Goal Details ───────────────────────────────────────┐  │
│  │                                                       │  │
│  │  ┌──────────────────────┬──────────────────────────┐ │  │
│  │  │ Goal Name *          │ Status                   │ │  │
│  │  │ Complete Project     │ ● In Progress            │ │  │
│  │  │ Alpha                │ [Auto-set from progress] │ │  │
│  │  └──────────────────────┴──────────────────────────┘ │  │
│  │                                                       │  │
│  │  ┌──────────────────────┬──────────────────────────┐ │  │
│  │  │ Employee *           │ Employee Name            │ │  │
│  │  │ ┌──────────────────┐ │ John Doe                 │ │  │
│  │  │ │ 🔍 Search...     │ │ [Read-only]              │ │  │
│  │  │ └──────────────────┘ │                          │ │  │
│  │  └──────────────────────┴──────────────────────────┘ │  │
│  │                                                       │  │
│  │  ┌──────────────────────┬──────────────────────────┐ │  │
│  │  │ Company *            │ Appraisal Cycle          │ │  │
│  │  │ Acme Corp            │ Q1 2024 Review           │ │  │
│  │  └──────────────────────┴──────────────────────────┘ │  │
│  │                                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─ Goal Structure ─────────────────────────────────────┐  │
│  │                                                       │  │
│  │  ☐ Is Group (Parent Goal)                            │  │
│  │  ⓘ Check if this goal will have sub-goals            │  │
│  │                                                       │  │
│  │  Parent Goal                                          │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ 📁 Delivery Excellence                      ▼  │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │  [Filtered to show only group goals for this employee]│ │
│  │                                                       │  │
│  │  Key Result Area (KRA)                                │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ Delivery                                    ▼  │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │  [Filtered from employee's appraisal KRAs]            │  │
│  │                                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─ Progress ───────────────────────────────────────────┐  │
│  │                                                       │  │
│  │  Progress (%)                                         │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ ████████████████░░░░░░░░░░░░░░  65%             │ │  │
│  │  │ ▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁               │ │  │
│  │  │ 0           25          50          75      100  │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  │  ⓘ Group goal's progress is auto-calculated based on │  │
│  │    the average progress of all child goals            │  │
│  │                                                       │  │
│  │  ┌──────────────────────┬──────────────────────────┐ │  │
│  │  │ Start Date           │ End Date                 │ │  │
│  │  │ 📅 Jan 1, 2024       │ 📅 Feb 15, 2024          │ │  │
│  │  └──────────────────────┴──────────────────────────┘ │  │
│  │                                                       │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ ⏱ 45 days │ 30 days elapsed │ 15 days remaining  │ │  │
│  │  │ ████████████████████░░░░░░░░  67% time elapsed  │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
│  ┌─ Description ────────────────────────────────────────┐  │
│  │                                                       │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ Successfully complete and deploy Project Alpha  │ │  │
│  │  │ with all required features and documentation... │ │  │
│  │  │                                                 │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
└─────────────────────────────────────────────────────────────┘

[Status ▼]  [Archive]  [Close]  [Unarchive]  [Reopen]
```

### Progress Update Dialog

```
┌──────────────────────────────────────┐
│  Update Goal Progress           [×] │
├──────────────────────────────────────┤
│                                      │
│  Goal Name                           │
│  Complete Project Alpha              │
│  [Read-only]                         │
│                                      │
│  Current Progress: 45%               │
│  ████████████░░░░░░░░░░░░░░          │
│                                      │
│  New Progress (%) *                  │
│  ┌──────────────────────────────┐   │
│  │ 65                       ▼  │   │
│  └──────────────────────────────┘   │
│                                      │
│  ┌──────────────────────────────┐   │
│  │ ████████████████░░░░░░░░░░░░ │   │
│  │ 0      25      50      75  100│   │
│  └──────────────────────────────┘   │
│                                      │
│                   [Cancel] [Update]  │
└──────────────────────────────────────┘
```

---

## Reports and Visualizations

### Appraisal Overview Report

```
┌─────────────────────────────────────────────────────────────┐
│  Appraisal Overview Report                    [Export ▼]   │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Filters:                                                   │
│  ┌──────────────┬──────────────┬──────────────────────┐    │
│  │ Cycle        │ Company      │ Department           │    │
│  │ Q1 2024  ▼  │ Acme Corp ▼ │ [Optional]       ▼  │    │
│  └──────────────┴──────────────┴──────────────────────┘    │
│                                          [Run Report]       │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Employee    │ Des.  │ Goal│ Self│ Feedback│ Final  │   │
│  ├─────────────────────────────────────────────────────┤   │
│  │ John Doe    │ Eng   │ 4.2 │ 4.0 │ 4.5    │ 4.23   │   │
│  │             │       │ ★★★★│ ★★★★│ ★★★★★ │ ★★★★☆ │   │
│  ├─────────────────────────────────────────────────────┤   │
│  │ Jane Smith  │ Mgr   │ 4.8 │ 4.5 │ 4.7    │ 4.67   │   │
│  │             │       │ ★★★★│ ★★★★│ ★★★★★ │ ★★★★★ │   │
│  ├─────────────────────────────────────────────────────┤   │
│  │ Bob Johnson │ Des   │ 3.5 │ 3.8 │ 4.0    │ 3.77   │   │
│  │             │       │ ★★★☆│ ★★★★│ ★★★★  │ ★★★★  │   │
│  ├─────────────────────────────────────────────────────┤   │
│  │ AVERAGE     │       │ 4.17│ 4.10│ 4.40   │ 4.22   │   │
│  │             │       │ ★★★★│ ★★★★│ ★★★★☆ │ ★★★★☆ │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  ┌─── Performance Distribution ─────────────────────────┐  │
│  │                                                       │  │
│  │  Excellent (4.5-5.0):  ██████████████ 35%            │  │
│  │  Good (3.5-4.4):       ████████████████████ 50%      │  │
│  │  Average (2.5-3.4):    ██████ 15%                    │  │
│  │  Poor (< 2.5):         0%                            │  │
│  │                                                       │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## Interaction States

### Button States

**Primary Action (Submit/Complete)**:
```
Normal:    [✓ Submit Appraisal]  - Green #98D85B
Hover:     [✓ Submit Appraisal]  - Darker Green #7CB842
Active:    [✓ Submit Appraisal]  - Pressed state
Disabled:  [✓ Submit Appraisal]  - Gray #9B9B9B
```

**Secondary Actions**:
```
Normal:    [📊 View Goals]  - White with border
Hover:     [📊 View Goals]  - Light gray background
Active:    [📊 View Goals]  - Darker gray background
```

### Progress Bars

**Normal**:
```
████████████░░░░░░░░░░░░░░  50%
```
- Filled: #2490EF (Blue)
- Empty: #E2E8F0 (Light Gray)
- Border: 1px solid #CBD5E0
- Border radius: 4px

**Completed (100%)**:
```
████████████████████████████  100%
```
- Filled: #98D85B (Green)

**Low Progress (< 25%)**:
```
███░░░░░░░░░░░░░░░░░░░░░░░░  15%
```
- Filled: #F56B6B (Red)

### Rating Stars

**Interactive (for input)**:
```
☆☆☆☆☆  (No rating)
★★★☆☆  (3/5 rating)
★★★★★  (5/5 rating)
```
- Empty: #E2E8F0
- Filled: #FFA00A (Orange/Gold)
- Hover: #FFB84D (Lighter Orange)
- Size: 20px × 20px

**Display-only**:
```
★★★★☆ 4.2/5.0
```
- Filled: #FFA00A
- Partial fill for decimal ratings
- Text color: #4A5568

---

## Mobile Responsiveness

### Mobile Appraisal View (< 768px)

```
┌───────────────────────────────┐
│  ☰  Appraisal            [×] │
│  John Doe                     │
├───────────────────────────────┤
│                               │
│  [Employee Photo]             │
│                               │
│  ┌─ Summary ─────────────┐   │
│  │ Final Score           │   │
│  │   4.05 / 5.00         │   │
│  │   ★★★★☆              │   │
│  └───────────────────────┘   │
│                               │
│  ┌─ Scores ──────────────┐   │
│  │ Goals:     3.75       │   │
│  │ Self:      4.30       │   │
│  │ Feedback:  4.10       │   │
│  └───────────────────────┘   │
│                               │
│  [View Details] ▼             │
│                               │
└───────────────────────────────┘
```

---

## Accessibility

- **Minimum Touch Target**: 44px × 44px
- **Color Contrast**: 4.5:1 for text, 3:1 for large text
- **Focus Indicators**: 2px blue outline
- **Screen Reader**: ARIA labels on all interactive elements
- **Keyboard Navigation**: Full Tab/Enter/Escape support

---

## Animation and Transitions

**Progress Bars**: Smooth fill animation over 500ms
**Score Updates**: Fade transition 300ms
**Tree Expand/Collapse**: Slide animation 200ms
**Button Hover**: 100ms ease
**Form Field Focus**: 200ms ease

---

This designer documentation provides complete UI/UX specifications for the Performance Management module.
