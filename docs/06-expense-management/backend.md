# Expense Management Module - Backend Documentation

## Overview

Business logic, validations, and data processing for employee expense claims, approvals, payments, and reimbursements with full accounting integration.

---

## Core Doctypes

### 1. Expense Claim

**Purpose**: Employee expense reimbursement requests with approval workflow and payment tracking.

#### Fields

**Basic Information**:
- `employee` (Link, required)
- `employee_name` (Data, read-only)
- `company` (Link, required)
- `posting_date` (Date, default: today)
- `expense_approver` (Link): Department approver
- `approval_status` (Select): Draft, Approved, Rejected, Cancelled
- `status` (Select): Draft, Submitted, Paid, Unpaid, Rejected, Cancelled

**Project Tracking**:
- `project` (Link)
- `task` (Link)
- `cost_center` (Link)
- `department` (Link)

**Expenses (Child Table - Expense Claim Detail)**:
- `expense_date` (Date)
- `expense_type` (Link to Expense Claim Type)
- `default_account` (Link to Account)
- `amount` (Currency): Claimed amount
- `sanctioned_amount` (Currency): Approved amount
- `description` (Text)
- `cost_center` (Link)
- `project` (Link)

**Advances (Child Table - Expense Claim Advance)**:
- `employee_advance` (Link)
- `posting_date` (Date)
- `advance_account` (Link)
- `advance_paid` (Currency)
- `unclaimed_amount` (Currency)
- `allocated_amount` (Currency): Amount to deduct from claim
- `exchange_rate` (Float)

**Taxes (Child Table - Expense Taxes and Charges)**:
- `account_head` (Link)
- `description` (Text)
- `rate` (Percent)
- `tax_amount` (Currency)
- `total` (Currency)

**Totals**:
- `total_claimed_amount` (Currency)
- `total_sanctioned_amount` (Currency)
- `total_taxes_and_charges` (Currency)
- `total_advance_amount` (Currency)
- `grand_total` (Currency): sanctioned + taxes - advances

**Payment**:
- `is_paid` (Check)
- `mode_of_payment` (Link)
- `payable_account` (Link, required when submitting)
- `total_amount_reimbursed` (Currency): Paid amount

**Multi-currency**:
- `currency` (Link)
- `exchange_rate` (Float)
- `base_*` fields: Amounts in company currency

#### Business Logic

**Validation**:
```python
def validate(self):
    validate_active_employee(self.employee)
    set_employee_name(self)
    self.validate_sanctioned_amount()
    self.calculate_total_amount()
    self.validate_advances()
    self.set_expense_account(validate=True)
    self.calculate_taxes()
    self.set_status()
    self.validate_company_and_department()

def validate_company_and_department(self):
    if self.department:
        company = frappe.db.get_value("Department", self.department, "company")
        if company and self.company != company:
            frappe.throw(f"Department {self.department} does not belong to company: {self.company}")

def validate_sanctioned_amount(self):
    # Ensure sanctioned amount not greater than claimed amount
    for expense in self.expenses:
        if flt(expense.sanctioned_amount) > flt(expense.amount):
            frappe.throw(f"Sanctioned Amount cannot be greater than Claimed Amount")
```

**Prevent Self-Approval**:
```python
def validate_for_self_approval(self):
    self_expense_approval_not_allowed = frappe.db.get_single_value(
        "HR Settings", "prevent_self_expense_approval"
    )
    employee_user = frappe.db.get_value("Employee", self.employee, "user_id")

    if (self_expense_approval_not_allowed and
        employee_user == frappe.session.user and
        not get_workflow_name("Expense Claim")):
        frappe.throw("Self-approval for Expense Claims is not allowed")
```

**Status Calculation**:
```python
def set_status(self, update=False):
    status = {"0": "Draft", "1": "Submitted", "2": "Cancelled"}[str(self.docstatus or 0)]

    if self.docstatus == 1:
        if self.approval_status == "Approved":
            if (self.is_paid or
                (flt(self.total_sanctioned_amount) > 0 and
                 (flt(self.grand_total) == flt(self.total_amount_reimbursed) or
                  flt(self.grand_total) == 0))):
                status = "Paid"
            elif flt(self.total_sanctioned_amount) > 0:
                status = "Unpaid"
        elif self.approval_status == "Rejected":
            status = "Rejected"

    if update:
        self.db_set("status", status)
    else:
        self.status = status
```

**Calculate Totals**:
```python
def calculate_total_amount(self):
    self.total_claimed_amount = sum(flt(d.amount) for d in self.expenses)
    self.total_sanctioned_amount = sum(flt(d.sanctioned_amount) for d in self.expenses)

    # Convert to base currency
    for expense in self.expenses:
        expense.base_amount = flt(expense.amount) * flt(self.exchange_rate)
        expense.base_sanctioned_amount = flt(expense.sanctioned_amount) * flt(self.exchange_rate)

def calculate_taxes(self):
    self.total_taxes_and_charges = 0

    for tax in self.taxes:
        if tax.rate:
            tax.tax_amount = flt(self.total_sanctioned_amount) * flt(tax.rate) / 100
        tax.total = flt(self.total_sanctioned_amount) + flt(tax.tax_amount)
        tax.base_tax_amount = flt(tax.tax_amount) * flt(self.exchange_rate)
        tax.base_total = flt(tax.total) * flt(self.exchange_rate)

        self.total_taxes_and_charges += flt(tax.tax_amount)

    self.base_total_taxes_and_charges = flt(self.total_taxes_and_charges) * flt(self.exchange_rate)

def calculate_grand_total(self):
    self.total_advance_amount = sum(flt(d.allocated_amount) for d in self.advances)

    self.grand_total = (flt(self.total_sanctioned_amount) +
                       flt(self.total_taxes_and_charges) -
                       flt(self.total_advance_amount))

    self.base_grand_total = flt(self.grand_total) * flt(self.exchange_rate)
```

**Advance Allocation**:
```python
def validate_advances(self):
    """Auto-allocate employee advances against claim"""
    allocated = 0

    for advance in self.advances:
        unclaimed = flt(advance.unclaimed_amount) - flt(advance.return_amount)

        if allocated < self.total_sanctioned_amount:
            if (self.total_sanctioned_amount - allocated) >= unclaimed:
                advance.allocated_amount = unclaimed
            else:
                advance.allocated_amount = self.total_sanctioned_amount - allocated

            allocated += advance.allocated_amount

        advance.base_allocated_amount = flt(advance.allocated_amount) * flt(self.exchange_rate)
```

**Accounting Entries**:
```python
def on_submit(self):
    if self.approval_status == "Draft":
        frappe.throw("Approval Status must be 'Approved' or 'Rejected'")

    self.update_task_and_project()
    self.make_gl_entries()
    update_reimbursed_amount(self)
    self.update_claimed_amount_in_employee_advance()
    self.create_exchange_gain_loss_je()

def make_gl_entries(self, cancel=False):
    if flt(self.total_sanctioned_amount) > 0:
        gl_entries = self.get_gl_entries()
        make_gl_entries(gl_entries, cancel)

def get_gl_entries(self):
    gl_entry = []

    # Payable entry (Credit)
    if self.grand_total:
        gl_entry.append({
            "account": self.payable_account,
            "credit": self.base_grand_total,
            "credit_in_account_currency": self.grand_total,
            "against": ",".join([d.default_account for d in self.expenses]),
            "party_type": "Employee",
            "party": self.employee,
            "cost_center": self.cost_center,
            "project": self.project
        })

    # Expense entries (Debit)
    for expense in self.expenses:
        gl_entry.append({
            "account": expense.default_account,
            "debit": expense.base_sanctioned_amount,
            "debit_in_account_currency": expense.sanctioned_amount,
            "against": self.employee,
            "cost_center": expense.cost_center or self.cost_center,
            "project": expense.project or self.project
        })

    # Advance entries (Credit) - deduct from payable
    for advance in self.advances:
        if advance.allocated_amount:
            gl_entry.append({
                "account": advance.advance_account,
                "credit": advance.base_allocated_amount,
                "credit_in_account_currency": advance.allocated_amount,
                "against": ",".join([d.default_account for d in self.expenses]),
                "party_type": "Employee",
                "party": self.employee,
                "advance_voucher_type": "Employee Advance",
                "advance_voucher_no": advance.employee_advance
            })

    # Tax entries (Debit)
    for tax in self.taxes:
        gl_entry.append({
            "account": tax.account_head,
            "debit": tax.base_tax_amount,
            "debit_in_account_currency": tax.tax_amount,
            "against": self.employee,
            "cost_center": self.cost_center
        })

    # If paid immediately (Credit payment account)
    if self.is_paid and self.grand_total:
        payment_account = get_bank_cash_account(self.mode_of_payment, self.company).get("account")
        gl_entry.append({
            "account": payment_account,
            "credit": self.base_grand_total,
            "credit_in_account_currency": self.grand_total,
            "against": self.employee
        })

    return gl_entry
```

**Update Project/Task**:
```python
def update_task_and_project(self):
    if self.task:
        task = frappe.get_doc("Task", self.task)

        # Sum all submitted expense claims for this task
        ExpenseClaim = frappe.qb.DocType("Expense Claim")
        task.total_expense_claim = (
            frappe.qb.from_(ExpenseClaim)
            .select(Sum(ExpenseClaim.total_sanctioned_amount))
            .where(
                (ExpenseClaim.docstatus == 1) &
                (ExpenseClaim.project == self.project) &
                (ExpenseClaim.task == self.task)
            )
        ).run()[0][0]

        task.save()
    elif self.project:
        frappe.get_doc("Project", self.project).update_project()
```

**Update Advance Claimed Amount**:
```python
def update_claimed_amount_in_employee_advance(self):
    for advance in self.advances:
        frappe.get_doc("Employee Advance", advance.employee_advance).update_claimed_amount()
```

**Notifications**:
```python
def after_insert(self):
    self.notify_approver()

def notify_approver(self):
    if self.expense_approver:
        parent_doc = frappe.get_doc("Employee", self.expense_approver)
        parent_doc.notify({
            "subject": f"New Expense Claim {self.name} awaiting approval",
            "message": f"Employee {self.employee_name} has submitted an expense claim for {self.grand_total}"
        })
```

---

### 2. Expense Claim Type

**Purpose**: Categorize expenses with default accounting configuration.

#### Fields

- `expense_type` (Data, required): e.g., "Travel", "Food", "Accommodation"
- `description` (Text)
- `accounts` (Child Table - Expense Claim Account):
  - `company` (Link)
  - `default_account` (Link to Account)

#### Business Logic

```python
@frappe.whitelist()
def get_expense_claim_account_and_cost_center(expense_claim_type, company):
    """Returns default account and cost center for expense type"""
    account = frappe.db.get_value(
        "Expense Claim Account",
        {"parent": expense_claim_type, "company": company},
        "default_account"
    )

    cost_center = frappe.get_cached_value("Company", company, "cost_center")

    return {
        "account": account,
        "cost_center": cost_center
    }
```

---

## Whitelisted APIs

### Get Advances
```python
@frappe.whitelist()
def get_advances(expense_claim, advance_id=None):
    """Returns unclaimed employee advances for the employee"""
    if isinstance(expense_claim, str):
        expense_claim = frappe.get_doc("Expense Claim", expense_claim)

    filters = {
        "docstatus": 1,
        "employee": expense_claim.employee,
        "paid_amount": (">", 0),
        "status": ("not in", ["Claimed", "Returned", "Partly Claimed and Returned"])
    }

    if advance_id:
        filters["name"] = advance_id

    advances = frappe.get_all(
        "Employee Advance",
        filters=filters,
        fields=["name as employee_advance", "posting_date", "advance_account",
               "paid_amount as advance_paid", "claimed_amount", "return_amount"]
    )

    for advance in advances:
        advance.unclaimed_amount = (flt(advance.advance_paid) -
                                   flt(advance.claimed_amount) -
                                   flt(advance.return_amount))

    return advances
```

### Make Payment Entry
```python
@frappe.whitelist()
def make_bank_entry(dt, dn):
    """Creates Journal Entry for expense claim payment"""
    expense_claim = frappe.get_doc(dt, dn)

    je = frappe.new_doc("Journal Entry")
    je.voucher_type = "Bank Entry"
    je.company = expense_claim.company
    je.posting_date = today()

    # Debit payable account
    je.append("accounts", {
        "account": expense_claim.payable_account,
        "debit_in_account_currency": expense_claim.grand_total,
        "party_type": "Employee",
        "party": expense_claim.employee,
        "reference_type": "Expense Claim",
        "reference_name": expense_claim.name
    })

    # Credit bank account
    je.append("accounts", {
        "account": get_bank_cash_account(expense_claim.mode_of_payment, expense_claim.company).get("account"),
        "credit_in_account_currency": expense_claim.grand_total
    })

    return je.as_dict()
```

### Update Reimbursed Amount
```python
def update_reimbursed_amount(expense_claim):
    """Updates total_amount_reimbursed from linked payment entries"""
    amount_reimbursed = frappe.db.sql("""
        SELECT sum(debit_in_account_currency)
        FROM `tabJournal Entry Account`
        WHERE reference_type = 'Expense Claim'
        AND reference_name = %s
        AND docstatus = 1
        AND account = %s
    """, (expense_claim.name, expense_claim.payable_account))[0][0]

    expense_claim.db_set("total_amount_reimbursed", flt(amount_reimbursed))
    expense_claim.set_status(update=True)
```

---

## Key Business Rules

1. **Approval Required**: Must be approved/rejected before submission
2. **No Self-Approval**: Employees cannot approve own claims (configurable)
3. **Payable Account Required**: Must be set before submission
4. **Sanctioned vs Claimed**: Sanctioned amount cannot exceed claimed amount
5. **Advance Auto-allocation**: Employee advances auto-deducted from payout
6. **Multi-currency Support**: Expenses can be in different currencies
7. **Status Automation**: Status auto-updated based on payment and approval
8. **Accounting Integration**: Full GL entries created on submission
9. **Project Tracking**: Updates project/task expense totals
10. **Exchange Gain/Loss**: Auto-creates journal entries for currency differences

---

This backend documentation provides complete business logic for the Expense Management module.
