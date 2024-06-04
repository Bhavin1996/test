package main

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	Proposed  = "proposed"
	Approved  = "approved"
	Invested  = "invested"
	Disbursed = "disbursed"
)

const (
	SMTPServer   = "smtp.example.com"
	SMTPPort     = "587"
	SMTPUser     = "your-email@example.com"  // put dummy emails here
	SMTPPassword = "your-email-password"     // put dummy emails here
	FromEmail    = "your-email@example.com"  // put dummy emails here
)


type Loan struct {
	ID               string            `json:"id"`
	BorrowerID       string            `json:"borrower_id"`
	PrincipalAmount  float64           `json:"principal_amount"`
	Rate             float64           `json:"rate"`
	ROI              float64           `json:"roi"`
	AgreementLetter  string            `json:"agreement_letter"`
	State            string            `json:"state"`
	ApprovalInfo     *ApprovalInfo     `json:"approval_info,omitempty"`
	Investments      []Investment      `json:"investments,omitempty"`
	DisbursementInfo *DisbursementInfo `json:"disbursement_info,omitempty"`
}

type ApprovalInfo struct {
	PictureProof string    `json:"picture_proof"`
	EmployeeID   string    `json:"employee_id"`
	ApprovalDate time.Time `json:"approval_date"`
}

type Investment struct {
	InvestorID string  `json:"investor_id"`
	Amount     float64 `json:"amount"`
}

type DisbursementInfo struct {
	SignedAgreementLetter string    `json:"signed_agreement_letter"`
	EmployeeID            string    `json:"employee_id"`
	DisbursementDate      time.Time `json:"disbursement_date"`
}

var loans = make(map[string]*Loan)

func createLoan(c *gin.Context) {
	var loan Loan
	if err := c.BindJSON(&loan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	loan.ID = generateID()
	loan.State = Proposed
	loans[loan.ID] = &loan
	c.JSON(http.StatusOK, loan)
}

func approveLoan(c *gin.Context) {
	id := c.Param("id")
	loan, exists := loans[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "loan not found"})
		return
	}

	if loan.State != Proposed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "loan is not in proposed state"})
		return
	}

	var approvalInfo ApprovalInfo
	if err := c.BindJSON(&approvalInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loan.ApprovalInfo = &approvalInfo
	loan.State = Approved
	c.JSON(http.StatusOK, loan)
}

func investLoan(c *gin.Context) {
	id := c.Param("id")
	loan, exists := loans[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "loan not found"})
		return
	}

	if loan.State != Approved {
		c.JSON(http.StatusBadRequest, gin.H{"error": "loan is not in approved state"})
		return
	}

	var investment Investment
	if err := c.BindJSON(&investment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	totalInvested := 0.0
	for _, inv := range loan.Investments {
		totalInvested += inv.Amount
	}
	totalInvested += investment.Amount

	if totalInvested > loan.PrincipalAmount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "total invested amount cannot exceed the loan principal"})
		return
	}

	loan.Investments = append(loan.Investments, investment)

	if totalInvested == loan.PrincipalAmount {
		loan.State = Invested
		sendInvestorEmails(loan)
	}

	c.JSON(http.StatusOK, loan)
}

func disburseLoan(c *gin.Context) {
	id := c.Param("id")
	loan, exists := loans[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "loan not found"})
		return
	}

	if loan.State != Invested {
		c.JSON(http.StatusBadRequest, gin.H{"error": "loan is not in invested state"})
		return
	}

	var disbursementInfo DisbursementInfo
	if err := c.BindJSON(&disbursementInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loan.DisbursementInfo = &disbursementInfo
	loan.State = Disbursed
	c.JSON(http.StatusOK, loan)
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func sendInvestorEmails(loan *Loan) {
	for _, investment := range loan.Investments {
		toEmail := investment.InvestorID //  InvestorID is the email address
		subject := "Loan Investment Agreement"
		body := fmt.Sprintf("Dear Investor,\n\nThank you for your investment. Please find the agreement letter here: %s\n\nBest regards,\nLoan Service Team", loan.AgreementLetter)
		err := sendEmail(toEmail, subject, body)
		if err != nil {
			log.Printf("Failed to send email to %s: %v", toEmail, err)
		} else {
			log.Printf("Email sent to %s successfully", toEmail)
		}
	}
}

func sendEmail(toEmail, subject, body string) error {
	auth := smtp.PlainAuth("", SMTPUser, SMTPPassword, SMTPServer)
	msg := []byte("To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail(SMTPServer+":"+SMTPPort, auth, FromEmail, []string{toEmail}, msg)
	return err
}

func main() {
	r := gin.Default()

	r.POST("/loans", createLoan)
	r.POST("/loans/:id/approve", approveLoan)
	r.POST("/loans/:id/invest", investLoan)
	r.POST("/loans/:id/disburse", disburseLoan)

	r.Run(":8080")
}
