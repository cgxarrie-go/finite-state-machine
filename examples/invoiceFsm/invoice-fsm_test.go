package invoiceFsm

import (
	"testing"

	"github.com/cgxarrie/fsm-go/pkg/fsm"
	"github.com/stretchr/testify/assert"
)

func Test_ConfirmCommand(t *testing.T) {
	tests := []struct {
		name          string
		from          InvoiceState
		to            InvoiceState
		needSignature bool
		wantError     bool
	}{
		{
			name:          "draft.NoSignature",
			from:          draft,
			to:            waitingForApproval,
			needSignature: false,
			wantError:     false,
		},
		{
			name:          "draft.Signature",
			from:          draft,
			to:            waitingForApproval,
			needSignature: true,
			wantError:     false,
		},
		{
			name:          "waitingForApproval.NoSignature",
			from:          waitingForApproval,
			to:            waitingForApproval,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "waitingForApproval.Signature",
			from:          waitingForApproval,
			to:            waitingForApproval,
			needSignature: true,
			wantError:     true,
		},
		{
			name:          "waitingForPayment.NoSignature",
			from:          waitingForPayment,
			to:            waitingForPayment,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "waitingForPayment.Signature",
			from:          waitingForPayment,
			to:            waitingForPayment,
			needSignature: true,
			wantError:     true,
		},
		{
			name:          "waitingForsignature.NoSignature",
			from:          waitingForsignature,
			to:            waitingForsignature,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "waitingForsignature.Signature",
			from:          waitingForsignature,
			to:            waitingForsignature,
			needSignature: true,
			wantError:     true,
		},
		{
			name:          "rejected.NoSignature",
			from:          rejected,
			to:            rejected,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "rejected.Signature",
			from:          rejected,
			to:            rejected,
			needSignature: true,
			wantError:     true,
		},
		{
			name:          "completed.NoSignature",
			from:          completed,
			to:            completed,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "completed.Signature",
			from:          completed,
			to:            completed,
			needSignature: true,
			wantError:     true,
		},
		{
			name:          "abandoned.NoSignature",
			from:          abandoned,
			to:            abandoned,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "abandoned.Signature",
			from:          abandoned,
			to:            abandoned,
			needSignature: true,
			wantError:     true,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			inv := NewInvoice(test.needSignature)
			inv.SetState(fsm.State(test.from))
			sm := NewInvoiceStateMachine(&inv)
			err := sm.ExecuteCommand("confirm")
			if !test.wantError {
				assert.NoError(t, err)
				assert.Equal(t, fsm.State(test.to), inv.State())
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func Test_ReceiveSignatureCommand(t *testing.T) {
	tests := []struct {
		name          string
		from          InvoiceState
		to            InvoiceState
		needSignature bool
		wantError     bool
	}{
		{
			name:          "draft.NoSignature",
			from:          draft,
			to:            draft,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "draft.Signature",
			from:          draft,
			to:            draft,
			needSignature: true,
			wantError:     true,
		},
		{
			name:          "waitingForApproval.NoSignature",
			from:          waitingForApproval,
			to:            waitingForApproval,
			needSignature: false,
			wantError:     false,
		},
		{
			name:          "waitingForApproval.Signature",
			from:          waitingForApproval,
			to:            waitingForApproval,
			needSignature: true,
			wantError:     false,
		},
		{
			name:          "waitingForPayment.NoSignature",
			from:          waitingForPayment,
			to:            waitingForPayment,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "waitingForPayment.Signature",
			from:          waitingForPayment,
			to:            waitingForPayment,
			needSignature: true,
			wantError:     true,
		},
		{
			name:          "waitingForsignature.NoSignature",
			from:          waitingForsignature,
			to:            waitingForPayment,
			needSignature: false,
			wantError:     false,
		},
		{
			name:          "waitingForsignature.Signature",
			from:          waitingForsignature,
			to:            waitingForPayment,
			needSignature: true,
			wantError:     false,
		},
		{
			name:          "rejected.NoSignature",
			from:          rejected,
			to:            rejected,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "rejected.Signature",
			from:          rejected,
			to:            rejected,
			needSignature: true,
			wantError:     true,
		},
		{
			name:          "completed.NoSignature",
			from:          completed,
			to:            completed,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "completed.Signature",
			from:          completed,
			to:            completed,
			needSignature: true,
			wantError:     true,
		},
		{
			name:          "abandoned.NoSignature",
			from:          abandoned,
			to:            abandoned,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "abandoned.Signature",
			from:          abandoned,
			to:            abandoned,
			needSignature: true,
			wantError:     true,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			inv := NewInvoice(test.needSignature)
			inv.SetState(fsm.State(test.from))
			sm := NewInvoiceStateMachine(&inv)
			err := sm.ExecuteCommand("receiveSignature")
			if !test.wantError {
				assert.NoError(t, err)
				assert.Equal(t, fsm.State(test.to), inv.State())
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func Test_RejectCommand(t *testing.T) {
	tests := []struct {
		name          string
		from          InvoiceState
		to            InvoiceState
		needSignature bool
		wantError     bool
	}{
		{
			name:          "draft.NoSignature",
			from:          draft,
			to:            draft,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "draft.Signature",
			from:          draft,
			to:            draft,
			needSignature: true,
			wantError:     true,
		},
		{
			name:          "waitingForApproval.NoSignature",
			from:          waitingForApproval,
			to:            rejected,
			needSignature: false,
			wantError:     false,
		},
		{
			name:          "waitingForApproval.Signature",
			from:          waitingForApproval,
			to:            rejected,
			needSignature: true,
			wantError:     false,
		},
		{
			name:          "waitingForPayment.NoSignature",
			from:          waitingForPayment,
			to:            waitingForPayment,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "waitingForPayment.Signature",
			from:          waitingForPayment,
			to:            waitingForPayment,
			needSignature: true,
			wantError:     true,
		},
		{
			name:          "waitingForsignature.NoSignature",
			from:          waitingForsignature,
			to:            waitingForsignature,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "waitingForsignature.Signature",
			from:          waitingForsignature,
			to:            waitingForsignature,
			needSignature: true,
			wantError:     true,
		},
		{
			name:          "rejected.NoSignature",
			from:          rejected,
			to:            rejected,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "rejected.Signature",
			from:          rejected,
			to:            rejected,
			needSignature: true,
			wantError:     true,
		},
		{
			name:          "completed.NoSignature",
			from:          completed,
			to:            completed,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "completed.Signature",
			from:          completed,
			to:            completed,
			needSignature: true,
			wantError:     true,
		},
		{
			name:          "abandoned.NoSignature",
			from:          abandoned,
			to:            abandoned,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "abandoned.Signature",
			from:          abandoned,
			to:            abandoned,
			needSignature: true,
			wantError:     true,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			inv := NewInvoice(test.needSignature)
			inv.SetState(fsm.State(test.from))
			sm := NewInvoiceStateMachine(&inv)
			err := sm.ExecuteCommand("reject")
			if !test.wantError {
				assert.NoError(t, err)
				assert.Equal(t, fsm.State(test.to), inv.State())
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func Test_ApproveCommand(t *testing.T) {
	tests := []struct {
		name          string
		from          InvoiceState
		to            InvoiceState
		needSignature bool
		wantError     bool
	}{
		{
			name:          "draft.NoSignature",
			from:          draft,
			to:            draft,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "draft.Signature",
			from:          draft,
			to:            draft,
			needSignature: true,
			wantError:     true,
		},
		{
			name:          "waitingForApproval.NoSignature",
			from:          waitingForApproval,
			to:            waitingForPayment,
			needSignature: false,
			wantError:     false,
		},
		{
			name:          "waitingForApproval.Signature",
			from:          waitingForApproval,
			to:            waitingForsignature,
			needSignature: true,
			wantError:     false,
		},
		{
			name:          "waitingForPayment.NoSignature",
			from:          waitingForPayment,
			to:            waitingForPayment,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "waitingForPayment.Signature",
			from:          waitingForPayment,
			to:            waitingForPayment,
			needSignature: true,
			wantError:     true,
		},
		{
			name:          "waitingForsignature.NoSignature",
			from:          waitingForsignature,
			to:            waitingForsignature,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "waitingForsignature.Signature",
			from:          waitingForsignature,
			to:            waitingForsignature,
			needSignature: true,
			wantError:     true,
		},
		{
			name:          "rejected.NoSignature",
			from:          rejected,
			to:            rejected,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "rejected.Signature",
			from:          rejected,
			to:            rejected,
			needSignature: true,
			wantError:     true,
		},
		{
			name:          "completed.NoSignature",
			from:          completed,
			to:            completed,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "completed.Signature",
			from:          completed,
			to:            completed,
			needSignature: true,
			wantError:     true,
		},
		{
			name:          "abandoned.NoSignature",
			from:          abandoned,
			to:            abandoned,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "abandoned.Signature",
			from:          abandoned,
			to:            abandoned,
			needSignature: true,
			wantError:     true,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			inv := NewInvoice(test.needSignature)
			inv.SetState(fsm.State(test.from))
			sm := NewInvoiceStateMachine(&inv)
			err := sm.ExecuteCommand("approve")
			if !test.wantError {
				assert.NoError(t, err)
				assert.Equal(t, fsm.State(test.to), inv.State())
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func Test_PayCommand(t *testing.T) {
	tests := []struct {
		name          string
		from          InvoiceState
		to            InvoiceState
		needSignature bool
		wantError     bool
	}{
		{
			name:          "draft.NoSignature",
			from:          draft,
			to:            draft,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "draft.Signature",
			from:          draft,
			to:            draft,
			needSignature: true,
			wantError:     true,
		},
		{
			name:          "waitingForApproval.NoSignature",
			from:          waitingForApproval,
			to:            waitingForApproval,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "waitingForApproval.Signature",
			from:          waitingForApproval,
			to:            waitingForApproval,
			needSignature: true,
			wantError:     true,
		},
		{
			name:          "waitingForPayment.NoSignature",
			from:          waitingForPayment,
			to:            completed,
			needSignature: false,
			wantError:     false,
		},
		{
			name:          "waitingForPayment.Signature",
			from:          waitingForPayment,
			to:            completed,
			needSignature: true,
			wantError:     false,
		},
		{
			name:          "waitingForsignature.NoSignature",
			from:          waitingForsignature,
			to:            waitingForsignature,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "waitingForsignature.Signature",
			from:          waitingForsignature,
			to:            waitingForsignature,
			needSignature: true,
			wantError:     true,
		},
		{
			name:          "rejected.NoSignature",
			from:          rejected,
			to:            rejected,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "rejected.Signature",
			from:          rejected,
			to:            rejected,
			needSignature: true,
			wantError:     true,
		},
		{
			name:          "completed.NoSignature",
			from:          completed,
			to:            completed,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "completed.Signature",
			from:          completed,
			to:            completed,
			needSignature: true,
			wantError:     true,
		},
		{
			name:          "abandoned.NoSignature",
			from:          abandoned,
			to:            abandoned,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "abandoned.Signature",
			from:          abandoned,
			to:            abandoned,
			needSignature: true,
			wantError:     true,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			inv := NewInvoice(test.needSignature)
			inv.SetState(fsm.State(test.from))
			sm := NewInvoiceStateMachine(&inv)
			err := sm.ExecuteCommand("pay")
			if !test.wantError {
				assert.NoError(t, err)
				assert.Equal(t, fsm.State(test.to), inv.State())
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func Test_AbandonCommand(t *testing.T) {
	tests := []struct {
		name          string
		from          InvoiceState
		to            InvoiceState
		needSignature bool
		wantError     bool
	}{
		{
			name:          "draft.NoSignature",
			from:          draft,
			to:            abandoned,
			needSignature: false,
			wantError:     false,
		},
		{
			name:          "draft.Signature",
			from:          draft,
			to:            abandoned,
			needSignature: true,
			wantError:     false,
		},
		{
			name:          "waitingForApproval.NoSignature",
			from:          waitingForApproval,
			to:            abandoned,
			needSignature: false,
			wantError:     false,
		},
		{
			name:          "waitingForApproval.Signature",
			from:          waitingForApproval,
			to:            abandoned,
			needSignature: true,
			wantError:     false,
		},
		{
			name:          "waitingForPayment.NoSignature",
			from:          waitingForPayment,
			to:            abandoned,
			needSignature: false,
			wantError:     false,
		},
		{
			name:          "waitingForPayment.Signature",
			from:          waitingForPayment,
			to:            abandoned,
			needSignature: true,
			wantError:     false,
		},
		{
			name:          "waitingForsignature.NoSignature",
			from:          waitingForsignature,
			to:            abandoned,
			needSignature: false,
			wantError:     false,
		},
		{
			name:          "waitingForsignature.Signature",
			from:          waitingForsignature,
			to:            abandoned,
			needSignature: true,
			wantError:     false,
		},
		{
			name:          "rejected.NoSignature",
			from:          rejected,
			to:            rejected,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "rejected.Signature",
			from:          rejected,
			to:            rejected,
			needSignature: true,
			wantError:     true,
		},
		{
			name:          "completed.NoSignature",
			from:          completed,
			to:            completed,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "completed.Signature",
			from:          completed,
			to:            completed,
			needSignature: true,
			wantError:     true,
		},
		{
			name:          "abandoned.NoSignature",
			from:          abandoned,
			to:            abandoned,
			needSignature: false,
			wantError:     true,
		},
		{
			name:          "abandoned.Signature",
			from:          abandoned,
			to:            abandoned,
			needSignature: true,
			wantError:     true,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			inv := NewInvoice(test.needSignature)
			inv.SetState(fsm.State(test.from))
			sm := NewInvoiceStateMachine(&inv)
			err := sm.ExecuteCommand("abandon")
			if !test.wantError {
				assert.NoError(t, err)
				assert.Equal(t, fsm.State(test.to), inv.State())
			} else {
				assert.Error(t, err)
			}
		})
	}
}
