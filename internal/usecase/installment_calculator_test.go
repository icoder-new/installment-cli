package usecase_test

import (
	"errors"
	"testing"

	"github.com/icoder-new/installment-cli/internal/domain"
	"github.com/icoder-new/installment-cli/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockSMSSender struct {
	mock.Mock
}

func (m *MockSMSSender) SendSMS(phoneNumber string, message string) error {
	args := m.Called(phoneNumber, message)
	return args.Error(0)
}

func TestInstallmentCalculator_CalculateInstallment(t *testing.T) {
	tests := []struct {
		name           string
		product        domain.Product
		setupMocks     func(*MockSMSSender)
		expectedResult float64
		expectError    bool
		errorMessage   string
	}{
		{
			name: "Smartphone 3 months - no interest",
			product: domain.Product{
				Type:         domain.Smartphone,
				Price:        1000,
				PhoneNumber:  "+992001002005",
				PeriodMonths: 3,
			},
			setupMocks: func(m *MockSMSSender) {
				m.On("SendSMS", "+992001002005", mock.Anything).Return(nil)
			},
			expectedResult: 1000,
			expectError:    false,
		},
		{
			name: "Computer 6 months - with interest",
			product: domain.Product{
				Type:         domain.Computer,
				Price:        3000,
				PhoneNumber:  "+992001002005",
				PeriodMonths: 6,
			},
			setupMocks: func(m *MockSMSSender) {
				m.On("SendSMS", "+992001002005", mock.Anything).Return(nil)
			},
			expectedResult: 3120,
			expectError:    false,
		},
		{
			name: "TV 12 months - with interest",
			product: domain.Product{
				Type:         domain.TV,
				Price:        2000,
				PhoneNumber:  "+992001002005",
				PeriodMonths: 12,
			},
			setupMocks: func(m *MockSMSSender) {
				m.On("SendSMS", "+992001002005", mock.Anything).Return(nil)
			},
			expectedResult: 2300,
			expectError:    false,
		},
		{
			name: "Invalid period - too short",
			product: domain.Product{
				Type:         domain.Smartphone,
				Price:        1000,
				PhoneNumber:  "+992001002005",
				PeriodMonths: 2,
			},
			setupMocks:     func(m *MockSMSSender) {},
			expectedResult: 0,
			expectError:    true,
			errorMessage:   "неверный срок рассрочки: для Смартфон допустимый срок от 3 до 9 месяцев",
		},
		{
			name: "Invalid period - not in allowed values",
			product: domain.Product{
				Type:         domain.Smartphone,
				Price:        1000,
				PhoneNumber:  "+992001002005",
				PeriodMonths: 4,
			},
			setupMocks:     func(m *MockSMSSender) {},
			expectedResult: 0,
			expectError:    true,
			errorMessage:   "неверный срок рассрочки: допустимые значения: [3 6 9 12 18 24]",
		},
		{
			name: "SMS send failure",
			product: domain.Product{
				Type:         domain.Smartphone,
				Price:        1000,
				PhoneNumber:  "+992001002005",
				PeriodMonths: 3,
			},
			setupMocks: func(m *MockSMSSender) {
				m.On("SendSMS", "+992001002005", mock.Anything).Return(errors.New("sms service unavailable"))
			},
			expectedResult: 0,
			expectError:    true,
			errorMessage:   "не удалось отправить SMS: sms service unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSMS := new(MockSMSSender)
			tt.setupMocks(mockSMS)

			calculator := usecase.NewInstallmentCalculator(mockSMS)

			result, err := calculator.CalculateInstallment(tt.product)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
			} else {
				assert.NoError(t, err)
				assert.InDelta(t, tt.expectedResult, result, 0.0001, "Expected result to be within 0.0001 of %v, got %v", tt.expectedResult, result)
			}

			mockSMS.AssertExpectations(t)
		})
	}
}

func TestInterestCalculation(t *testing.T) {
	tests := []struct {
		name           string
		product        domain.Product
		expectedAmount float64
	}{
		{
			name: "Smartphone 3 months - no interest",
			product: domain.Product{
				Type:         domain.Smartphone,
				Price:        1000,
				PhoneNumber:  "+992001002005",
				PeriodMonths: 3,
			},
			expectedAmount: 1000,
		},
		{
			name: "Smartphone 6 months - 3% interest",
			product: domain.Product{
				Type:         domain.Smartphone,
				Price:        1000,
				PhoneNumber:  "+992001002005",
				PeriodMonths: 6,
			},
			expectedAmount: 1030,
		},
		{
			name: "Computer 12 months - 12% interest (3 periods of 4%)",
			product: domain.Product{
				Type:         domain.Computer,
				Price:        2000,
				PhoneNumber:  "+992001002005",
				PeriodMonths: 12,
			},
			expectedAmount: 2240,
		},
		{
			name: "TV 18 months - 25% interest (5 periods of 5%)",
			product: domain.Product{
				Type:         domain.TV,
				Price:        3000,
				PhoneNumber:  "+992001002005",
				PeriodMonths: 18,
			},
			expectedAmount: 3750,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSMS := new(MockSMSSender)
			mockSMS.On("SendSMS", "+992001002005", mock.Anything).Return(nil)

			calculator := usecase.NewInstallmentCalculator(mockSMS)

			result, err := calculator.CalculateInstallment(tt.product)
			require.NoError(t, err)
			assert.InDelta(t, tt.expectedAmount, result, 0.01, "Expected amount to be within 0.01 of %v, got %v", tt.expectedAmount, result)

			mockSMS.AssertExpectations(t)
		})
	}
}

func TestProduct_Validate(t *testing.T) {
	tests := []struct {
		name        string
		product     domain.Product
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid smartphone",
			product: domain.Product{
				Type:         domain.Smartphone,
				Price:        1000,
				PhoneNumber:  "+992001002005",
				PeriodMonths: 6,
			},
			expectError: false,
		},
		{
			name: "Invalid price",
			product: domain.Product{
				Type:         domain.Smartphone,
				Price:        0,
				PhoneNumber:  "+992001002005",
				PeriodMonths: 6,
			},
			expectError: true,
			errorMsg:    "цена должна быть больше 0",
		},
		{
			name: "Missing phone number",
			product: domain.Product{
				Type:         domain.Smartphone,
				Price:        1000,
				PhoneNumber:  "",
				PeriodMonths: 6,
			},
			expectError: true,
			errorMsg:    "необходимо указать номер телефона",
		},
		{
			name: "Invalid period for computer - too long",
			product: domain.Product{
				Type:         domain.Computer,
				Price:        3000,
				PhoneNumber:  "+992001002005",
				PeriodMonths: 15,
			},
			expectError: true,
			errorMsg:    "неверный срок рассрочки: для Компьютер допустимый срок от 3 до 12 месяцев",
		},
		{
			name: "Invalid period - not in allowed values",
			product: domain.Product{
				Type:         domain.Computer,
				Price:        3000,
				PhoneNumber:  "+992001002005",
				PeriodMonths: 7,
			},
			expectError: true,
			errorMsg:    "неверный срок рассрочки: допустимые значения: [3 6 9 12 18 24]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.product.Validate()
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
