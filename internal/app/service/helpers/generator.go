package helpers

import (
	"crypto/rand"
	"erp/internal/domain"
	_logger "erp/pkg/logger"
	"fmt"
	"strconv"
	"strings"

	"github.com/sethvargo/go-password/password"
)

type Generator interface {
	GeneratePassword() string
	GenerateCode() string
	GenerateSN(template string, startNum int, count int) ([]string, error)
	GenerateCodeAutoIncrement(template string, start int64) (string, error)
}

type generatorHelper struct {
	logger     _logger.Logger
	emitErrLog func(err error, options ..._logger.OptionLog)
}

func NewGeneratorHelper(
	logger _logger.Logger,
) Generator {
	return &generatorHelper{
		logger: logger,
		emitErrLog: func(err error, options ..._logger.OptionLog) {
			logger.LogError(err, append(options, _logger.OptionsLog.WithOperation("generator"),
				_logger.OptionsLog.WithFilename("generator"))...)
		},
	}
}

func (h *generatorHelper) GeneratePassword() string {
	res, err := password.Generate(10, 3, 0, false, false)
	if err != nil {
		h.emitErrLog(err, _logger.OptionsLog.WithMethod("GeneratePassword"))
	}
	return res
}

func (h *generatorHelper) GenerateCode() string {
	n := 4
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%X", b)
}
func (h *generatorHelper) GenerateSN(template string, startNum int, count int) ([]string, error) {
	var snList []string

	// Find the position of the placeholder '########' in the template
	parts := strings.Split(template, "#")
	prefix := parts[0]
	suffixLength := len(template) - len(prefix) // Calculate length of the number part (########)
	for i := startNum + 1; i <= startNum+count; i++ {
		// Convert the number to string
		numStr := strconv.Itoa(i)
		// Calculate how many leading zeros are needed
		leadingZeros := suffixLength - len(numStr)
		if leadingZeros < 0 {
			return []string{}, domain.OVERFLOW_SN
		}
		fmt.Println("LEADING ZEROS", leadingZeros)
		zeroPadding := strings.Repeat("0", leadingZeros) // Create the zero padding
		// Create the full serial number by appending the padded number to the prefix
		sn := prefix + zeroPadding + numStr
		snList = append(snList, sn)
	}

	return snList, nil
}

func (h *generatorHelper) GenerateCodeAutoIncrement(template string, start int64) (string, error) {
	parts := strings.Split(template, "#")
	prefix := parts[0]
	suffixLength := len(template) - len(prefix) // Calculate length of the number part (########)
	numStr := strconv.Itoa(int(start + 1))
	leadingZeros := suffixLength - len(numStr)
	if leadingZeros < 0 {
		return "", domain.OVERFLOW_SN
	}
	zeroPadding := strings.Repeat("0", leadingZeros) // Create the zero padding
	// Create the full serial number by appending the padded number to the prefix
	res := prefix + zeroPadding + numStr

	return res,nil
}
