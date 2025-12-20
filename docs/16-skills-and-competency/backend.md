# Skills and Competency Module - Backend Documentation

## Overview

Manage employee skills, competencies, training needs, and skill gap analysis.

## Core Doctype: Skill

### Fields
- `skill_name`: Skill identifier
- `description`: Skill definition
- `skill_category`: Technical/Soft/Domain/Language

## Core Doctype: Employee Skill Map

### Fields
- `employee`: Employee reference
- `skills`: Child table with:
  - `skill`: Skill reference
  - `proficiency`: Beginner/Intermediate/Advanced/Expert
  - `evaluation_date`: When assessed
  - `evaluator`: Who assessed

### Business Logic
- Self-assessment vs manager assessment
- Skill certification tracking
- Proficiency scoring (0-100)
- Skill decay over time

## Core Doctype: Training Program

### Fields
- `training_program_name`: Program identifier
- `description`: Program details
- `trainer_name`, `trainer_email`
- `skills_covered`: Child table of skills
- `duration_hours`: Total duration

## Core Doctype: Training Event

### Fields
- `training_program`: Program reference
- `event_name`: Specific event instance
- `start_date`, `end_date`
- `location`: Physical/Virtual
- `attendees`: Child table with employees
- `training_status`: Scheduled/Completed/Cancelled

### Business Logic
- Attendance tracking
- Pre/post assessment
- Certificate generation
- Cost allocation

## Core Doctype: Training Result

### Fields
- `training_event`: Event reference
- `employee`: Trainee reference
- `grade`: Pass/Fail/Grade
- `comments`: Feedback
- `certificate_issued`: Boolean

## Skill Gap Analysis

### Logic
- Compare required skills (job role) vs current skills
- Identify gaps
- Recommend training programs
- Track skill improvement over time

## Key Rules
- Skills must be validated by manager
- Training attendance mandatory for identified gaps
- Certification expiry tracking
- Skill proficiency updated post-training
- Career progression based on competencies

## Gap Analysis Logic
```python
def analyze_skill_gap(employee):
    required_skills = get_role_required_skills(employee.designation)
    current_skills = get_employee_skills(employee)
    gaps = []
    for required in required_skills:
        current = find_skill(current_skills, required.skill)
        if not current or current.proficiency < required.min_proficiency:
            gaps.append({
                'skill': required.skill,
                'required_level': required.min_proficiency,
                'current_level': current.proficiency if current else 0
            })
    return gaps
```
