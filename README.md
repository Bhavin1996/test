# Steps to run the code and other information

## Note ( All values are dummy here please add valid values to run the code properly )

### In this we are directly creating all the structures and their mappings.
### Recommended way is to have a segregated utils folder

## Steps for SMTP configuration 
### Note - Make sure to add demo emails don't put anything here as mails are associtaed with active phone numbers keep it private
- ### Replace SMTPServer, SMTPPort, SMTPUser, SMTPPassword, and FromEmail with your actual SMTP server details.
- ### Use Gmail service for testing, but remember to enable "less secure apps" or use an app-specific password.

## Run the main code
### Use curl or postman for following configuration
- ### Endpoint: POST /loans
- ### Request body {
  "borrower_id": "12345",
  "principal_amount": 1000.00,
  "rate": 5.0,
  "roi": 10.0,
  "agreement_letter": "http://example.com/agreement.pdf"
}

- ### Endpoint: POST /loans/:id/approve ( replace id as actual loan_id returned by the above request )
- ### Request Body  {
  "picture_proof": "http://example.com/proof.jpg",
  "employee_id": "emp123",
  "approval_date": "2023-01-01T00:00:00Z"
}

- ### Endpoint: POST /loans/:id/invest (id - loan_id)
- ### Request Body {
  "investor_id": "investor1@example.com",
  "amount": 500.00
}

- ### Endpoint POST /loans/:id/disburse
- ### Request Body {
  "signed_agreement_letter": "http://example.com/signed_agreement.pdf",
  "employee_id": "emp456",
  "disbursement_date": "2023-01-02T00:00:00Z"
}

### Ensure all values are valid like mail address and all
### Once loan reaches the "invested" state then an email will be received by the receipant 
