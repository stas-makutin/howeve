package zwave

import "testing"

func TestZWaveUtils(t *testing.T) {

	dataFrame := DataRequest([]byte{1, 2, 3})

	t.Run("Valid request data frame", func(t *testing.T) {
		if res, pos := ValidateDataFrame(dataFrame); res != FrameOK || pos != len(dataFrame) {
			t.Errorf("The request data frame is invalid. Validation result: %v, pos: %d", res, pos)
		}
	})

	dataFrame = DataResponse([]byte{4, 5, 6})

	t.Run("Valid response data frame", func(t *testing.T) {
		if res, pos := ValidateDataFrame(dataFrame); res != FrameOK || pos != len(dataFrame) {
			t.Errorf("The response data frame is invalid. Validation result: %v, pos: %d", res, pos)
		}
	})

	t.Run("Incomplete data frame", func(t *testing.T) {
		if res, _ := ValidateDataFrame(dataFrame[0 : len(dataFrame)-1]); res != FrameIncomplete {
			t.Errorf("The data frame checksum must be invalid. Instead the validation result is: %v", res)
		}
	})

	dataFrame[len(dataFrame)-1] -= 1
	t.Run("Invalid data frame checksum", func(t *testing.T) {
		if res, _ := ValidateDataFrame(dataFrame); res != FrameWrongChecksum {
			t.Errorf("The data frame checksum must be invalid. Instead the validation result is: %v", res)
		}
	})

	dataFrame[1] = 0
	t.Run("Invalid data frame length", func(t *testing.T) {
		if res, _ := ValidateDataFrame(dataFrame); res != FrameWrongLength {
			t.Errorf("The data frame checksum must be invalid. Instead the validation result is: %v", res)
		}
	})
}
