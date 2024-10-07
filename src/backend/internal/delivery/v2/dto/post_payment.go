package dto

type Amount struct {
    Value    string `json:"value" binding:"required"`
    Currency string `json:"currency" binding:"required"`
}

type Confirmation struct {
    Type      string `json:"type" binding:"required"`
    ReturnURL string `json:"return_url" binding:"required"`
}

type PaymentRequest struct {
    Amount       Amount       `json:"amount" binding:"required"`
    Confirmation Confirmation `json:"confirmation" binding:"required"`
    Capture      bool         `json:"capture" binding:"required"`
    Description  string       `json:"description" binding:"required"`
}

type PaymentResponse struct {
    ID          string `json:"id"`
    Status      string `json:"status"`
    Paid        bool `json:"paid"`
    Amount      struct {
        Value    string `json:"value"`
        Currency string `json:"currency"`
    } `json:"amount"`
    Confirmation struct {
        Type           string `json:"type"`
        ReturnURL      string `json:"return_url"`
        ConfirmationURL string `json:"confirmation_url"`
    } `json:"confirmation"`
    CreatedAt   string `json:"created_at"`
    Description string `json:"description"`
    Metadata    struct{} `json:"metadata"`
    PaymentMethod struct {
        Type  string `json:"type"`
        ID    string `json:"id"`
        Saved bool `json:"saved"`
    } `json:"payment_method"`
    Recipient struct {
        AccountID string `json:"account_id"`
        GatewayID string `json:"gateway_id"`
    } `json:"recipient"`
    Refundable bool `json:"refundable"`
    Test       bool `json:"test"`
}
