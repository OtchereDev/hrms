# Recruitment Module - Backend Documentation

## Overview

Manages the complete recruitment pipeline from job requisitions through candidate selection and employee creation, including applicant tracking, interviews, and job offers.

---

## Core Doctypes

### 1. Job Requisition

**Purpose**: Formal request to fill positions

**Auto-naming**: `HR-HIREQ-`

**Submittable**: No

#### Key Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| designation | Link | Yes | Position to fill |
| department | Link | No | Requesting department |
| no_of_positions | Int | Yes | Number of openings |
| expected_compensation | Currency | Yes | Expected salary |
| company | Link | Yes | Company |
| status | Select | Yes | Pending/Open & Approved/Rejected/Filled/On Hold/Cancelled |
| requested_by | Link | Yes | Employee making request |
| posting_date | Date | Yes | Request date |
| expected_by | Date | No | Target fill date |
| completed_on | Date | No | Actual fill date |
| time_to_fill | Duration | No | Calculated duration |
| description | Text Editor | Yes | Job description |
| reason_for_requesting | Text | No | Justification |

#### Business Logic

**Validation**:
```python
def validate_duplicates(self):
    # Prevent duplicate requisitions
    duplicate = frappe.db.exists("Job Requisition", {
        "designation": self.designation,
        "department": self.department or "",
        "requested_by": self.requested_by,
        "status": ("not in", ["Cancelled", "Filled"]),
        "name": ("!=", self.name)
    })

    if duplicate:
        frappe.throw(_("Requisition already exists for {0}").format(self.designation))
```

**Time to Fill Calculation**:
```python
def set_time_to_fill(self):
    if self.status == "Filled" and self.completed_on:
        delta = datetime.combine(self.completed_on, datetime.min.time()) - \
                datetime.combine(self.posting_date, datetime.min.time())
        self.time_to_fill = int(delta.total_seconds())
```

**Key Methods**:
```python
@frappe.whitelist()
def make_job_opening(source_name):
    # Creates Job Opening from requisition
    requisition = frappe.get_doc("Job Requisition", source_name)

    job_opening = frappe.new_doc("Job Opening")
    job_opening.job_title = requisition.designation
    job_opening.designation = requisition.designation
    job_opening.company = requisition.company
    job_opening.department = requisition.department
    job_opening.vacancies = requisition.no_of_positions
    job_opening.description = requisition.description
    job_opening.job_requisition = requisition.name

    return job_opening

@frappe.whitelist()
def associate_job_opening(requisition, job_opening):
    # Links existing Job Opening to requisition
    job_opening_doc = frappe.get_doc("Job Opening", job_opening)
    job_opening_doc.job_requisition = requisition
    job_opening_doc.save()
```

---

### 2. Job Opening

**Purpose**: Public/internal job posting

**Auto-naming**: `HR-OPN-.YYYY.-.####`

**Submittable**: No

#### Key Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| job_title | Data | Yes | Position title |
| designation | Link | Yes | Designation |
| company | Link | Yes | Hiring company |
| status | Select | No | Open/Closed |
| posted_on | Datetime | No | Posting timestamp |
| closes_on | Date | No | Auto-close date |
| closed_on | Date | No | Actual close date |
| department | Link | No | Department |
| employment_type | Link | No | Full-time/Contract/etc |
| location | Link | No | Branch/location |
| staffing_plan | Link | No | Associated plan |
| vacancies | Int | No | Number of positions |
| publish | Check | No | Publish on website |
| route | Data | No | URL slug |
| description | Text Editor | No | Job details |
| lower_range | Currency | No | Min salary |
| upper_range | Currency | No | Max salary |
| publish_salary_range | Check | No | Show salary publicly |

#### Business Logic

**Validation**:
```python
def validate(self):
    if not self.route:
        self.route = f"jobs/{frappe.scrub(self.company)}/{frappe.scrub(self.job_title)}"

    validate_current_vacancies()

def validate_current_vacancies(self):
    if not self.staffing_plan:
        return

    # Check staffing plan limits
    plan_details = frappe.db.get_value("Staffing Plan Detail", {
        "parent": self.staffing_plan,
        "designation": self.designation
    }, ["vacancies", "from_date", "to_date"], as_dict=True)

    if plan_details:
        # Count existing openings + current employees
        existing_openings = frappe.db.count("Job Opening", {
            "designation": self.designation,
            "company": self.company,
            "status": "Open",
            "name": ("!=", self.name)
        })

        current_employees = frappe.db.count("Employee", {
            "designation": self.designation,
            "company": self.company,
            "status": "Active"
        })

        total = existing_openings + current_employees + (self.vacancies or 0)

        if total > plan_details.vacancies:
            frappe.throw(_("Total positions ({0}) exceed staffing plan ({1})").format(
                total, plan_details.vacancies
            ))
```

**Auto-Close**:
```python
def close_expired_job_openings():
    # Scheduled daily
    expired = frappe.get_all("Job Opening", {
        "status": "Open",
        "closes_on": ("<", today())
    })

    for opening in expired:
        doc = frappe.get_doc("Job Opening", opening.name)
        doc.status = "Closed"
        doc.closed_on = today()
        doc.save()
```

**On Update**:
```python
def on_update(self):
    if self.status == "Closed":
        update_job_requisition_status()

def update_job_requisition_status(self):
    if self.job_requisition:
        requisition = frappe.get_doc("Job Requisition", self.job_requisition)
        requisition.status = "Filled"
        requisition.completed_on = self.closed_on or today()
        requisition.save()
```

---

### 3. Job Applicant

**Purpose**: Candidate records and application tracking

**Auto-naming**: By `email_id` (with suffix if duplicate)

**Submittable**: No

#### Key Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| applicant_name | Data | Yes | Full name |
| email_id | Data | Yes | Email (unique per application) |
| phone_number | Data | No | Contact number |
| job_title | Link | No | Job Opening applied for |
| designation | Link | No | Target designation |
| status | Select | Yes | Open/Replied/Rejected/Hold/Accepted |
| source | Link | No | Application source |
| source_name | Link | No | Employee (if referral) |
| employee_referral | Link | No | Linked referral |
| applicant_rating | Rating | No | Recruiter rating |
| cover_letter | Text | No | Cover letter |
| resume_attachment | Attach | No | Resume file |
| resume_link | Data | No | Online resume URL |
| lower_range | Currency | No | Min salary expectation |
| upper_range | Currency | No | Max salary expectation |

#### Business Logic

**Auto-naming**:
```python
def autoname(self):
    # Use email as name
    self.name = self.email_id

    # Handle duplicates (same person, different job)
    if frappe.db.exists("Job Applicant", self.name):
        count = frappe.db.count("Job Applicant", {
            "name": ("like", f"{self.email_id}%")
        })
        self.name = f"{self.email_id}-{count}"
```

**Validation**:
```python
def validate(self):
    if not self.applicant_name:
        self.applicant_name = self.email_id.split("@")[0]

    validate_job_opening_closed()

def validate_job_opening_closed(self):
    if self.job_title:
        status = frappe.db.get_value("Job Opening", self.job_title, "status")
        if status == "Closed":
            frappe.throw(_("Job Opening is closed"))
```

**Employee Referral Sync**:
```python
def on_update(self):
    update_employee_referral_status()

def update_employee_referral_status(self):
    if self.employee_referral:
        referral = frappe.get_doc("Employee Referral", self.employee_referral)

        # Map statuses
        status_map = {
            "Open": "In Process",
            "Replied": "In Process",
            "Hold": "In Process",
            "Accepted": "Accepted",
            "Rejected": "Rejected"
        }

        referral.status = status_map.get(self.status, "In Process")
        referral.save()
```

**Key Methods**:
```python
@frappe.whitelist()
def create_interview(self):
    # Creates Interview document
    interview = frappe.new_doc("Interview")
    interview.job_applicant = self.name
    interview.designation = self.designation
    interview.resume_link = self.resume_link
    return interview

def get_interview_details(self):
    # Returns interview summary for dashboard
    interviews = frappe.get_all("Interview", {
        "job_applicant": self.name,
        "docstatus": ("<", 2)
    }, ["*"])

    for interview in interviews:
        interview.avg_rating = get_average_rating(interview.name)

    return interviews
```

---

### 4. Job Offer

**Purpose**: Formal offer letter to candidates

**Auto-naming**: `HR-OFF-.YYYY.-.#####`

**Submittable**: Yes

#### Key Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| job_applicant | Link | Yes | Selected candidate |
| applicant_name | Data | Yes | Applicant name |
| designation | Link | Yes | Offered position |
| company | Link | Yes | Offering company |
| status | Select | No | Awaiting Response/Accepted/Rejected/Cancelled |
| offer_date | Date | Yes | Offer date |
| offer_terms | Table | No | Terms and conditions |
| terms | Text Editor | No | Additional terms |

#### Business Logic

**Validation**:
```python
def validate(self):
    validate_vacancies()
    validate_duplicate_offer()

def validate_vacancies(self):
    # Check staffing plan if enabled
    hr_settings = frappe.get_single("HR Settings")
    if not hr_settings.check_vacancies:
        return

    # Get staffing plan for designation
    staffing_plan = frappe.db.get_value("Staffing Plan Detail", {
        "designation": self.designation,
        "parent": ("is", "set")
    }, ["parent", "vacancies"], as_dict=True)

    if staffing_plan:
        # Count offers + employees
        offer_count = frappe.db.count("Job Offer", {
            "designation": self.designation,
            "company": self.company,
            "status": ("in", ["Awaiting Response", "Accepted"]),
            "docstatus": 1,
            "name": ("!=", self.name)
        })

        employee_count = frappe.db.count("Employee", {
            "designation": self.designation,
            "company": self.company,
            "status": "Active"
        })

        if offer_count + employee_count >= staffing_plan.vacancies:
            frappe.throw(_("No vacancies available"))

def validate_duplicate_offer(self):
    duplicate = frappe.db.exists("Job Offer", {
        "job_applicant": self.job_applicant,
        "status": ("!=", "Rejected"),
        "docstatus": ("<", 2),
        "name": ("!=", self.name)
    })

    if duplicate:
        frappe.throw(_("Offer already exists for {0}").format(self.job_applicant))
```

**On Change**:
```python
def on_change(self):
    update_job_applicant_status()

def update_job_applicant_status(self):
    if self.status in ["Accepted", "Rejected"]:
        frappe.db.set_value("Job Applicant", self.job_applicant, "status", self.status)
```

**Employee Creation**:
```python
@frappe.whitelist()
def make_employee(source_name):
    # Maps offer to employee
    offer = frappe.get_doc("Job Offer", source_name)
    applicant = frappe.get_doc("Job Applicant", offer.job_applicant)

    employee = frappe.new_doc("Employee")
    employee.first_name = applicant.applicant_name
    employee.personal_email = applicant.email_id
    employee.designation = offer.designation
    employee.company = offer.company
    employee.scheduled_confirmation_date = offer.offer_date
    employee.job_applicant = offer.job_applicant

    return employee
```

---

### 5. Interview

**Purpose**: Interview scheduling and evaluation (part of ERPNext)

#### Key Fields

| Field | Type | Description |
|-------|------|-------------|
| interview_round | Link | Interview stage |
| job_applicant | Link | Candidate |
| designation | Link | Position |
| scheduled_on | Datetime | Interview date/time |
| from_time | Time | Start time |
| to_time | Time | End time |
| interviewers | Table | Panel members |
| skill_assessment | Table | Skills being assessed |
| status | Select | Pending/Cleared/Rejected |
| average_rating | Float | Calculated avg |

#### Business Logic

**Average Rating**:
```python
def calculate_average_rating(self):
    if not self.skill_assessment:
        return 0

    total = sum(skill.rating for skill in self.skill_assessment if skill.rating)
    count = len([s for s in self.skill_assessment if s.rating])

    return total / count if count else 0
```

---

## Supporting Doctypes

### Job Applicant Source

**Purpose**: Track application channels

**Fields**: `source_name` (unique)

### Offer Term

**Purpose**: Master list of offer terms

**Fields**: `offer_term` (unique)

### Job Offer Term (Child Table)

**Fields**:
- `offer_term`: Link to Offer Term
- `value`: Small Text (term details)

### Job Offer Term Template

**Purpose**: Pre-defined offer term sets

**Fields**:
- `title`: Template name
- `offer_terms`: Table of terms

---

## Web Forms

### Job Application

**Route**: `/job_application`

**Purpose**: Public application form

**Fields**:
- job_title (read-only from URL)
- applicant_name (required)
- email_id (required, validated)
- phone_number
- country
- cover_letter
- resume_link (validated URL)
- Currency fields for salary expectations

**Client Script**:
```javascript
// Validate resume link is URL
frappe.web_form.validate = () => {
    let resume = frappe.web_form.get_value("resume_link");
    if (resume && !isValidURL(resume)) {
        frappe.throw("Please enter a valid URL for resume");
    }
}
```

---

## Scheduled Jobs

### Close Expired Job Openings (Daily)

```python
def close_expired_job_openings():
    expired = frappe.get_all("Job Opening",
        filters={
            "status": "Open",
            "closes_on": ("<", today())
        }
    )

    for opening in expired:
        doc = frappe.get_doc("Job Opening", opening.name)
        doc.status = "Closed"
        doc.closed_on = today()
        doc.save()
```

---

## Integration Points

### Employee Creation Hook

```python
# In employee_master.py
def update_job_applicant_and_offer(self, method=None):
    if self.job_applicant:
        # Update applicant
        frappe.db.set_value("Job Applicant", self.job_applicant, "status", "Accepted")

        # Update offer
        job_offer = frappe.db.get_value("Job Offer", {
            "job_applicant": self.job_applicant,
            "docstatus": 1
        }, "name")

        if job_offer:
            frappe.db.set_value("Job Offer", job_offer, "status", "Accepted")
```

### Staffing Plan Integration

```python
def validate_against_staffing_plan(designation, company, date):
    # Check if position is planned
    plan_detail = frappe.db.get_value("Staffing Plan Detail", {
        "designation": designation,
        "from_date": ("<=", date),
        "to_date": (">=", date)
    }, ["parent", "vacancies"], as_dict=True)

    if not plan_detail:
        frappe.msgprint(_("No staffing plan found for {0}").format(designation))
        return False

    # Count current usage
    current_count = get_current_headcount(designation, company)

    return current_count < plan_detail.vacancies
```

---

## Workflow

### Complete Recruitment Pipeline

```
1. Job Requisition Created
   → Status: Pending

2. Approved by Manager
   → Status: Open & Approved

3. Job Opening Created
   → Links to requisition
   → Published on website (optional)

4. Applications Received
   → Job Applicants created
   → Status: Open

5. Screening & Interviews
   → Status: Replied (contacted)
   → Interviews scheduled
   → Ratings collected

6. Candidate Selected
   → Status: Accepted

7. Job Offer Created
   → Terms populated from template
   → Submit for approval

8. Offer Accepted
   → Status: Accepted

9. Employee Created
   → From Job Offer
   → Links applicant and offer
   → Updates all statuses

10. Job Opening Closed
    → Status: Closed

11. Requisition Marked Filled
    → Status: Filled
    → Time to fill calculated
```

---

## Formulas

### Time to Fill

```python
time_to_fill_seconds = (completed_on - posting_date).total_seconds()
time_to_fill_days = time_to_fill_seconds / 86400
```

### Applicant to Hire Percentage

```python
total_applicants = count(Job Applicant)
total_hired = count(Job Applicant where status = "Accepted")
percentage = (total_hired / total_applicants) * 100
```

### Offer Acceptance Rate

```python
total_offers = count(Job Offer where docstatus = 1)
accepted_offers = count(Job Offer where status = "Accepted")
acceptance_rate = (accepted_offers / total_offers) * 100
```

---

This backend documentation provides complete implementation details for the Recruitment module.
