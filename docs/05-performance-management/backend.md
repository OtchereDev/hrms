# Performance Management Module - Backend Documentation

## Overview

Business logic, validations, and data processing for the complete performance management system including appraisal cycles, goal tracking, performance reviews, and feedback management.

---

## Core Doctypes

### 1. Appraisal Cycle

**Purpose**: Manages periodic performance review cycles for employees.

#### Fields

**Basic Information**:
- `cycle_name` (Data, required): Unique name for the cycle
- `company` (Link, required): Company
- `status` (Select): In Progress, Completed
- `kra_evaluation_method` (Select): Automatic (Goal-based), Manual Rating

**Timeline**:
- `start_date` (Date, required)
- `end_date` (Date, required)

**Employee Selection Filters**:
- `department` (Link): Optional filter
- `branch` (Link): Optional filter
- `designation` (Link): Optional filter

**Appraisees (Child Table)**:
- `employee` (Link)
- `employee_name` (Data)
- `designation` (Link)
- `department` (Link)
- `branch` (Link)
- `appraisal_template` (Link)

**Scoring Configuration**:
- `calculate_final_score_based_on_formula` (Check)
- `final_score_formula` (Small Text): Custom formula for final score calculation

#### Business Logic

**Validation**:
```python
def validate(self):
    self.validate_from_to_dates("start_date", "end_date")
    self.validate_evaluation_method_change()

def validate_evaluation_method_change(self):
    if self.is_new():
        return

    if self.has_value_changed("kra_evaluation_method") and self.check_if_appraisals_exist():
        frappe.throw(
            "Evaluation Method cannot be changed as there are existing appraisals created for this cycle",
            title="Not Allowed"
        )
```

**Set Employees**:
```python
@frappe.whitelist()
def set_employees(self):
    """Pull employees in appraisee list based on selected filters"""
    employees = self.get_employees_for_appraisal()
    appraisal_templates = self.get_appraisal_template_map()

    if employees:
        self.set("appraisees", [])
        template_missing = False

        for data in employees:
            if not appraisal_templates.get(data.designation):
                template_missing = True

            self.append("appraisees", {
                "employee": data.name,
                "employee_name": data.employee_name,
                "branch": data.branch,
                "designation": data.designation,
                "department": data.department,
                "appraisal_template": appraisal_templates.get(data.designation)
            })

        if template_missing:
            self.show_missing_template_message()
    else:
        frappe.msgprint("No employees found for the selected criteria")

    return self

def get_employees_for_appraisal(self):
    filters = {
        "status": "Active",
        "company": self.company
    }
    if self.department:
        filters["department"] = self.department
    if self.branch:
        filters["branch"] = self.branch
    if self.designation:
        filters["designation"] = self.designation

    return frappe.db.get_all("Employee", filters=filters, fields=[...])

def get_appraisal_template_map(self):
    """Maps designation to appraisal template"""
    designations = frappe.get_all("Designation", fields=["name", "appraisal_template"])
    appraisal_templates = {}

    for entry in designations:
        appraisal_templates[entry.name] = entry.appraisal_template

    return appraisal_templates
```

**Create Appraisals**:
```python
@frappe.whitelist()
def create_appraisals(self):
    self.check_permission("write")

    if not self.appraisees:
        frappe.throw("Please select employees to create appraisals for",
                    title="No Employees Selected")

    if not all(appraisee.appraisal_template for appraisee in self.appraisees):
        self.show_missing_template_message(raise_exception=True)

    if len(self.appraisees) > 30:
        # Queue for background processing
        frappe.enqueue(
            create_appraisals_for_cycle,
            queue="long",
            timeout=600,
            appraisal_cycle=self
        )
        frappe.msgprint("Appraisal creation is queued. It may take a few minutes.",
                       alert=True, indicator="blue")
    else:
        create_appraisals_for_cycle(self, publish_progress=True)
        self.reload()

def create_appraisals_for_cycle(appraisal_cycle, publish_progress=False):
    """Creates appraisals for employees in the appraisee list"""
    count = 0

    for employee in appraisal_cycle.appraisees:
        try:
            appraisal = frappe.get_doc({
                "doctype": "Appraisal",
                "company": appraisal_cycle.company,
                "appraisal_template": employee.appraisal_template,
                "employee": employee.employee,
                "appraisal_cycle": appraisal_cycle.name
            })

            appraisal.rate_goals_manually = 1 if appraisal_cycle.kra_evaluation_method == "Manual Rating" else 0
            appraisal.set_kras_and_rating_criteria()
            appraisal.insert()

            if publish_progress:
                count += 1
                frappe.publish_progress(
                    count * 100 / len(appraisal_cycle.appraisees),
                    title="Creating Appraisals..."
                )
        except frappe.DuplicateEntryError:
            # already exists
            pass
```

**Complete Cycle**:
```python
@frappe.whitelist()
def complete_cycle(self):
    self.check_permission("write")

    draft_appraisals = frappe.db.count("Appraisal", {
        "appraisal_cycle": self.name,
        "docstatus": 0
    })

    if draft_appraisals:
        frappe.throw(
            f"{draft_appraisals} Appraisal(s) are not submitted yet. "
            "Please submit them before marking the cycle as Completed",
            title="Unsubmitted Appraisals"
        )

    self.status = "Completed"
    self.save()
```

**Cycle Summary**:
```python
@frappe.whitelist()
def get_appraisal_cycle_summary(cycle_name):
    summary = {}

    summary["appraisees"] = frappe.db.count("Appraisal", {
        "appraisal_cycle": cycle_name,
        "docstatus": ("!=", 2)
    })

    summary["self_appraisal_pending"] = frappe.db.count("Appraisal", {
        "appraisal_cycle": cycle_name,
        "docstatus": 0,
        "self_score": 0
    })

    summary["goals_missing"] = get_employees_without_goals(cycle_name)
    summary["feedback_missing"] = get_employees_without_feedback(cycle_name)

    return summary

def get_employees_without_goals(cycle_name):
    # Get employees who haven't set any goals
    ...

def get_employees_without_feedback(cycle_name):
    # Get employees who haven't received any feedback
    ...
```

---

### 2. Appraisal Template

**Purpose**: Defines standard KRAs and rating criteria for performance appraisals.

#### Fields

**Basic Information**:
- `template_name` (Data, required)
- `description` (Small Text)

**Goals/KRAs (Child Table)**:
- `key_result_area` (Data): Name of KRA
- `per_weightage` (Float): Weightage percentage

**Rating Criteria (Child Table)**:
- `criteria` (Data): Rating criteria name
- `per_weightage` (Float): Weightage percentage

#### Business Logic

**Validation**:
```python
def validate(self):
    self.validate_total_weightage("goals", "KRAs")
    self.validate_total_weightage("rating_criteria", "Criteria")

# From AppraisalMixin
def validate_total_weightage(self, table_name, label):
    total_weightage = sum(flt(entry.per_weightage) for entry in self.get(table_name))

    if total_weightage and flt(total_weightage, 2) != 100.0:
        frappe.throw(
            f"Total weightage for all {label} must add up to 100. Currently, it is {total_weightage}%",
            title="Incorrect Weightage Allocation"
        )
```

---

### 3. Appraisal

**Purpose**: Individual employee performance review document.

#### Fields

**Basic Information**:
- `employee` (Link, required)
- `employee_name` (Data, read-only)
- `appraisal_cycle` (Link, required)
- `appraisal_template` (Link, required)
- `company` (Link, required)
- `start_date` (Date)
- `end_date` (Date)
- `rate_goals_manually` (Check): Whether to rate goals manually or auto-calculate from goal progress

**KRAs (Child Table - auto-calculated mode)**:
- `kra` (Data): Key Result Area
- `per_weightage` (Float): Weightage percentage
- `goal_completion` (Percent): Average progress of goals under this KRA
- `goal_score` (Float): Calculated score = goal_completion × per_weightage / 100

**Goals (Child Table - manual rating mode)**:
- `kra` (Data): Key Result Area
- `per_weightage` (Float): Weightage percentage
- `score` (Rating): Manual score (1-5 stars)
- `score_earned` (Float): Calculated = score × per_weightage / 100

**Self Ratings (Child Table)**:
- `criteria` (Data): Rating criteria
- `per_weightage` (Float): Weightage percentage
- `rating` (Rating): Employee's self-rating (1-5)

**Scores**:
- `total_score` (Float): Goal/KRA score
- `goal_score_percentage` (Percent): Only in auto mode
- `self_score` (Float): Self-appraisal score
- `avg_feedback_score` (Float): Average of all peer/manager feedback
- `final_score` (Float): Final calculated score

#### Business Logic

**Validation**:
```python
def validate(self):
    self.set_kra_evaluation_method()
    validate_active_employee(self.employee)
    validate_active_appraisal_cycle(self.appraisal_cycle)

    self.validate_duplicate()
    self.validate_total_weightage("appraisal_kra", "KRAs")
    self.validate_total_weightage("self_ratings", "Self Ratings")

    self.set_goal_score()
    self.calculate_self_appraisal_score()
    self.calculate_avg_feedback_score()
    self.calculate_final_score()

def validate_duplicate(self):
    # Check for duplicate appraisal for same employee in same cycle or overlapping period
    duplicate = frappe.qb.from_(Appraisal).select(...).where(
        (employee == self.employee) &
        (docstatus != 2) &
        (name != self.name) &
        (
            (appraisal_cycle == self.appraisal_cycle) |
            # Check for date overlaps
            ...
        )
    ).run()

    if duplicate:
        frappe.throw("Appraisal already exists for this employee and cycle/period",
                    exc=frappe.DuplicateEntryError)
```

**Set Template**:
```python
@frappe.whitelist()
def set_appraisal_template(self):
    """Sets appraisal template from Appraisee table in Cycle"""
    if not self.appraisal_cycle:
        return

    appraisal_template = frappe.db.get_value("Appraisee", {
        "employee": self.employee,
        "parent": self.appraisal_cycle
    }, "appraisal_template")

    if appraisal_template:
        self.appraisal_template = appraisal_template
        self.set_kras_and_rating_criteria()

@frappe.whitelist()
def set_kras_and_rating_criteria(self):
    if not self.appraisal_template:
        return

    self.set("appraisal_kra", [])
    self.set("self_ratings", [])
    self.set("goals", [])

    template = frappe.get_doc("Appraisal Template", self.appraisal_template)

    # Populate KRAs or Goals based on evaluation method
    for entry in template.goals:
        table_name = "goals" if self.rate_goals_manually else "appraisal_kra"

        self.append(table_name, {
            "kra": entry.key_result_area,
            "per_weightage": entry.per_weightage
        })

    # Populate self-rating criteria
    for entry in template.rating_criteria:
        self.append("self_ratings", {
            "criteria": entry.criteria,
            "per_weightage": entry.per_weightage
        })

    return self
```

**Goal Score Calculation (Auto mode)**:
```python
def set_goal_score(self, update=False):
    """Calculates KRA scores based on average goal progress"""
    for kra in self.appraisal_kra:
        # Get average progress of all top-level goals for this KRA
        Goal = frappe.qb.DocType("Goal")
        avg_goal_completion = (
            frappe.qb.from_(Goal)
            .select(Avg(Goal.progress).as_("avg_goal_completion"))
            .where(
                (Goal.kra == kra.kra) &
                (Goal.employee == self.employee) &
                # archived goals should not contribute
                (Goal.status != "Archived") &
                # only top-level goals
                ((Goal.parent_goal == "") | (Goal.parent_goal.isnull())) &
                (Goal.appraisal_cycle == self.appraisal_cycle)
            )
        ).run()[0][0]

        kra.goal_completion = flt(avg_goal_completion, kra.precision("goal_completion"))
        kra.goal_score = flt(kra.goal_completion * kra.per_weightage / 100,
                            kra.precision("goal_score"))

        if update:
            kra.db_update()

    self.calculate_total_score()

    if update:
        self.calculate_final_score()
        self.db_update()

    return self
```

**Total Score Calculation**:
```python
def calculate_total_score(self):
    total_weightage, total, goal_score_percentage = 0, 0, 0
    number_of_stars = 5

    if self.rate_goals_manually:
        # Manual rating mode
        for entry in self.goals:
            if flt(entry.score) > flt(number_of_stars):
                frappe.throw(f"Row {entry.idx}: Goal Score cannot be greater than {number_of_stars}")

            entry.score_earned = flt(entry.score) * flt(entry.per_weightage) / 100
            total += flt(entry.score_earned)
            total_weightage += flt(entry.per_weightage)
    else:
        # Auto-calculated from goal progress
        for entry in self.appraisal_kra:
            goal_score_percentage += flt(entry.goal_score)
            total_weightage += flt(entry.per_weightage)

        self.goal_score_percentage = flt(goal_score_percentage,
                                        self.precision("goal_score_percentage"))
        # Convert goal score percentage to total score out of 5
        total = flt(goal_score_percentage) / 20

    if total_weightage and flt(total_weightage, 2) != 100.0:
        frappe.throw(
            f"Total weightage for all KRAs/Goals must add up to 100. Currently, it is {total_weightage}%",
            title="Incorrect Weightage Allocation"
        )

    self.total_score = flt(total, self.precision("total_score"))
```

**Self Appraisal Score**:
```python
def calculate_self_appraisal_score(self):
    total = 0
    number_of_stars = 5

    for entry in self.self_ratings:
        score = flt(entry.rating) * flt(number_of_stars) * flt(entry.per_weightage / 100)
        total += flt(score)

    self.self_score = flt(total, self.precision("self_score"))
```

**Average Feedback Score**:
```python
def calculate_avg_feedback_score(self, update=False):
    avg_feedback_score = frappe.qb.avg(
        "Employee Performance Feedback",
        "total_score",
        {
            "employee": self.employee,
            "appraisal": self.name,
            "docstatus": 1
        }
    )

    self.avg_feedback_score = flt(avg_feedback_score,
                                  self.precision("avg_feedback_score"))

    if update:
        self.calculate_final_score()
        self.db_update()
```

**Final Score Calculation**:
```python
def calculate_final_score(self):
    final_score = 0
    appraisal_cycle_doc = frappe.get_cached_doc("Appraisal Cycle", self.appraisal_cycle)

    formula = appraisal_cycle_doc.final_score_formula
    based_on_formula = appraisal_cycle_doc.calculate_final_score_based_on_formula

    if based_on_formula:
        # Custom formula evaluation
        employee_doc = frappe.get_cached_doc("Employee", self.employee)
        data = {
            "goal_score": flt(self.total_score),
            "average_feedback_score": flt(self.avg_feedback_score),
            "self_appraisal_score": flt(self.self_score)
        }
        data.update(appraisal_cycle_doc.as_dict())
        data.update(employee_doc.as_dict())
        data.update(self.as_dict())

        sanitized_formula = sanitize_expression(formula)
        final_score = frappe.safe_eval(sanitized_formula, data)
    else:
        # Default: Average of all three scores
        final_score = (flt(self.total_score) + flt(self.avg_feedback_score) + flt(self.self_score)) / 3

    self.final_score = flt(final_score, self.precision("final_score"))
```

**Add Feedback**:
```python
@frappe.whitelist()
def add_feedback(self, feedback, feedback_ratings):
    """Creates a new Employee Performance Feedback document"""
    feedback_doc = frappe.get_doc({
        "doctype": "Employee Performance Feedback",
        "appraisal": self.name,
        "employee": self.employee,
        "added_on": now(),
        "feedback": feedback,
        "reviewer": frappe.db.get_value("Employee", {"user_id": frappe.session.user})
    })

    for entry in feedback_ratings:
        feedback_doc.append("feedback_ratings", {
            "criteria": entry.get("criteria"),
            "rating": entry.get("rating"),
            "per_weightage": entry.get("per_weightage")
        })

    feedback_doc.submit()
    return feedback_doc
```

---

### 4. Goal

**Purpose**: Hierarchical goal management with parent-child relationships (tree structure).

#### Fields

**Basic Information**:
- `goal_name` (Data, required)
- `description` (Small Text)
- `employee` (Link, required)
- `employee_name` (Data, read-only)
- `company` (Link, required)
- `appraisal_cycle` (Link): Optional link to appraisal cycle
- `kra` (Data): Key Result Area
- `status` (Select): Pending, In Progress, Completed, Archived, Closed

**Hierarchy**:
- `is_group` (Check): Whether this goal has child goals
- `parent_goal` (Link): Parent goal (for nested goals)
- `lft`, `rgt` (Int): Nested set fields for tree structure

**Progress**:
- `progress` (Percent): 0-100%
- `start_date` (Date)
- `end_date` (Date)

#### Business Logic

**Validation**:
```python
def validate(self):
    if self.appraisal_cycle:
        validate_active_appraisal_cycle(self.appraisal_cycle)

    validate_active_employee(self.employee)
    self.validate_parent_fields()
    self.validate_from_to_dates(self.start_date, self.end_date)
    self.validate_progress()
    self.set_status()

def validate_parent_fields(self):
    """Validates child goal matches parent in employee, KRA, and cycle"""
    if not self.parent_goal:
        return

    parent_details = frappe.db.get_value("Goal", self.parent_goal,
                                        ["employee", "kra", "appraisal_cycle"],
                                        as_dict=True)

    if self.employee != parent_details.employee:
        frappe.throw("Goal should be owned by the same employee as its parent goal.",
                    title="Not Allowed")

    if self.kra != parent_details.kra:
        frappe.throw("Goal should be aligned with the same KRA as its parent goal.",
                    title="Not Allowed")

    if self.appraisal_cycle != parent_details.appraisal_cycle:
        frappe.throw("Goal should belong to the same Appraisal Cycle as its parent goal.",
                    title="Not Allowed")

def validate_progress(self):
    if flt(self.progress) > 100:
        frappe.throw("Goal progress percentage cannot be more than 100.")

def set_status(self):
    """Auto-set status based on progress"""
    if self.status in ["Archived", "Closed"]:
        return

    if flt(self.progress) == 0:
        self.status = "Pending"
    elif flt(self.progress) == 100:
        self.status = "Completed"
    elif flt(self.progress) < 100:
        self.status = "In Progress"
```

**Update Hierarchy**:
```python
def on_update(self):
    NestedSet.on_update(self)

    doc_before_save = self.get_doc_before_save()

    if doc_before_save:
        self.update_kra_in_child_goals(doc_before_save)

        if doc_before_save.parent_goal != self.parent_goal:
            # parent goal changed, update progress of old parent
            self.update_parent_progress(doc_before_save.parent_goal)

    self.update_parent_progress()
    self.update_goal_progress_in_appraisal()

def update_kra_in_child_goals(self, doc_before_save):
    """Aligns children's KRA to parent goal's KRA if parent goal's KRA is changed"""
    if doc_before_save.kra != self.kra and self.is_group:
        Goal = frappe.qb.DocType("Goal")
        (frappe.qb.update(Goal)
         .set(Goal.kra, self.kra)
         .where(Goal.parent_goal == self.name)).run()

        frappe.msgprint("KRA updated for all child goals.", alert=True, indicator="green")

def update_parent_progress(self, old_parent=None):
    """Updates parent goal's progress based on average of children"""
    parent_goal = old_parent or self.parent_goal

    if not parent_goal:
        return

    Goal = frappe.qb.DocType("Goal")
    avg_goal_completion = (
        frappe.qb.from_(Goal)
        .select(Avg(Goal.progress).as_("avg_goal_completion"))
        .where(
            (Goal.parent_goal == parent_goal) &
            (Goal.employee == self.employee) &
            # archived goals should not contribute
            (Goal.status != "Archived")
        )
    ).run()[0][0]

    parent_goal_doc = frappe.get_doc("Goal", parent_goal)
    parent_goal_doc.progress = flt(avg_goal_completion,
                                   parent_goal_doc.precision("progress"))
    parent_goal_doc.ignore_permissions = True
    parent_goal_doc.ignore_mandatory = True
    parent_goal_doc.save()

def update_goal_progress_in_appraisal(self):
    """Updates KRA scores in linked appraisal"""
    if not self.appraisal_cycle:
        return

    appraisal = frappe.db.get_value("Appraisal", {
        "employee": self.employee,
        "appraisal_cycle": self.appraisal_cycle
    })

    if appraisal:
        appraisal = frappe.get_doc("Appraisal", appraisal)
        appraisal.set_goal_score(update=True)
```

**Tree Operations**:
```python
@frappe.whitelist()
def get_children(doctype, parent, is_root=False, **filters):
    """Get child goals for tree view"""
    Goal = frappe.qb.DocType(doctype)

    query = (
        frappe.qb.from_(Goal)
        .select(
            Goal.name.as_("value"),
            Goal.goal_name.as_("title"),
            Goal.is_group.as_("expandable"),
            Goal.status,
            Goal.employee,
            Goal.progress,
            Goal.kra
        )
        .where(Goal.status != "Archived")
    )

    # Apply filters
    if filters.get("employee"):
        query = query.where(Goal.employee == filters.get("employee"))

    if filters.get("appraisal_cycle"):
        query = query.where(Goal.appraisal_cycle == filters.get("appraisal_cycle"))

    # Get top-level or child goals
    if filters.get("goal"):
        query = query.where(Goal.parent_goal == filters.get("goal"))
    elif parent and not is_root:
        query = query.where(Goal.parent_goal == parent)
    else:
        query = query.where(ifnull(Goal.parent_goal, "") == "")

    goals = query.orderby(Goal.employee, Goal.kra).run(as_dict=True)

    # Add completion counts for group goals
    for goal in goals:
        if goal.expandable:
            total_goals = frappe.db.count("Goal", dict(parent_goal=goal.value))
            if total_goals:
                completed = frappe.db.count("Goal", {
                    "parent_goal": goal.value,
                    "status": "Completed"
                }) or 0
                goal["completion_count"] = f"{completed} of {total_goals} Completed"

    return goals

@frappe.whitelist()
def update_progress(progress, goal):
    """Update goal progress"""
    goal = frappe.get_doc("Goal", goal)
    goal.progress = progress
    goal.flags.ignore_mandatory = True
    goal.save()
    return goal

@frappe.whitelist()
def update_status(status, goals):
    """Bulk update goal status"""
    if isinstance(goals, str):
        import json
        goals = json.loads(goals)

    for goal in goals:
        goal_doc = frappe.get_doc("Goal", goal)
        goal_doc.status = status
        if status == "Completed":
            goal_doc.progress = 100
        goal_doc.flags.ignore_mandatory = True
        goal_doc.save()

    return goals
```

---

### 5. Employee Performance Feedback

**Purpose**: Allows peers and managers to provide feedback on employee performance.

#### Fields

**Basic Information**:
- `employee` (Link, required): Employee being reviewed
- `appraisal` (Link, required): Linked appraisal
- `appraisal_cycle` (Link, read-only): From appraisal
- `reviewer` (Link, required): Employee providing feedback
- `reviewer_name`, `reviewer_designation` (Data, read-only)
- `added_on` (Datetime)
- `feedback` (Text): Written feedback

**Feedback Ratings (Child Table)**:
- `criteria` (Data): Rating criteria
- `per_weightage` (Float): Weightage percentage
- `rating` (Rating): Score (1-5)

**Score**:
- `total_score` (Float): Calculated total score

#### Business Logic

**Validation**:
```python
def validate(self):
    validate_active_appraisal_cycle(self.appraisal_cycle)
    self.validate_employee()
    self.validate_appraisal()
    self.validate_total_weightage("feedback_ratings", "Feedback Ratings")
    self.set_total_score()

def validate_employee(self):
    """Prevent self-feedback"""
    if self.employee == self.reviewer:
        frappe.throw(
            "Employees cannot give feedback to themselves. Use Self Appraisal instead",
            title="Not Allowed"
        )

    validate_active_employee(self.employee)
    validate_active_employee(self.reviewer)

def validate_appraisal(self):
    """Ensure appraisal belongs to employee"""
    employee = frappe.db.get_value("Appraisal", self.appraisal, "employee")

    if employee != self.employee:
        frappe.throw(f"Appraisal {self.appraisal} does not belong to Employee {self.employee}")
```

**Score Calculation**:
```python
def set_total_score(self):
    total = 0
    for entry in self.feedback_ratings:
        score = flt(entry.rating) * 5 * flt(entry.per_weightage / 100)
        total += flt(score)

    self.total_score = flt(total, self.precision("total_score"))
```

**Update Appraisal**:
```python
def on_submit(self):
    self.update_avg_feedback_score_in_appraisal()

def on_cancel(self):
    self.update_avg_feedback_score_in_appraisal()

def update_avg_feedback_score_in_appraisal(self):
    if not self.appraisal:
        return

    appraisal = frappe.get_doc("Appraisal", self.appraisal)
    appraisal.calculate_avg_feedback_score(update=True)
```

**Set Criteria from Template**:
```python
@frappe.whitelist()
def set_feedback_criteria(self):
    if not self.appraisal:
        return

    template = frappe.db.get_value("Appraisal", self.appraisal, "appraisal_template")
    template = frappe.get_doc("Appraisal Template", template)

    self.set("feedback_ratings", [])
    for entry in template.rating_criteria:
        self.append("feedback_ratings", {
            "criteria": entry.criteria,
            "per_weightage": entry.per_weightage
        })

    return self
```

---

## Whitelisted APIs

### Get Feedback History
```python
@frappe.whitelist()
def get_feedback_history(employee, appraisal):
    """Returns all feedback with rating distribution"""
    data = {}

    data.feedback_history = frappe.get_list(
        "Employee Performance Feedback",
        filters={"employee": employee, "appraisal": appraisal, "docstatus": 1},
        fields=["feedback", "reviewer", "reviewer_name", "reviewer_designation",
               "added_on", "total_score", "name"],
        order_by="added_on desc"
    )

    # Calculate percentage of reviews per rating (1-5)
    reviews_per_rating = []
    feedback_count = frappe.db.count("Employee Performance Feedback", {...})

    for i in range(1, 6):
        count = frappe.db.count("Employee Performance Feedback", {
            "appraisal": appraisal,
            "employee": employee,
            "total_score": ("between", [i, i + 0.99]),
            "docstatus": 1
        })

        percent = flt((count / feedback_count) * 100, 0) if feedback_count else 0
        reviews_per_rating.append(percent)

    data.reviews_per_rating = reviews_per_rating
    data.avg_feedback_score = frappe.db.get_value("Appraisal", appraisal, "avg_feedback_score")

    return data
```

---

## Key Business Rules

1. **Appraisal Cycle Status**: Once marked as "Completed", cannot create/modify appraisals
2. **Duplicate Prevention**: One appraisal per employee per cycle
3. **Weightage Validation**: All weightages must total 100%
4. **Goal Hierarchy**: Child goals must belong to same employee, KRA, and cycle as parent
5. **Parent Progress**: Auto-calculated as average of child goal progress
6. **KRA Alignment**: Changing parent goal's KRA updates all children
7. **Self-Feedback**: Employees cannot give feedback to themselves
8. **Score Calculation**: Final score can use custom formula or default average
9. **Goal Status**: Auto-set based on progress (0% = Pending, 100% = Completed, else In Progress)
10. **Archived Goals**: Excluded from progress calculations

---

This backend documentation provides complete business logic for the Performance Management module.
