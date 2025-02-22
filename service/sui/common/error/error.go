package error

import (
	"regexp"
	"strconv"
	"strings"
)

var ErrorMapping = map[int]string{
	257:    "Sy zero deposit",
	258:    "Sy insufficient sharesOut",
	259:    "Sy zero redeem",
	260:    "Sy insufficient amountOut",
	513:    "Interest fee rate too high",
	514:    "Reward fee rate too high",
	515:    "Factory zero expiry divisor",
	516:    "Factory invalid expiry",
	517:    "Factory invalid yt amount",
	518:    "Py contract exists",
	519:    "Mismatch yt pt tokens",
	520:    "Factory yc expired",
	521:    "Factory yc not expired",
	528:    "Invalid py state",
	769:    "Market scalar root below zero",
	770:    "Market pt expired",
	771:    "Market ln fee rate too high",
	772:    "Market initial anchor too low",
	773:    "Market factory reserve fee too high",
	774:    "Market exists",
	775:    "Market scalar root is zero",
	776:    "Market pt amount is zero",
	777:    "Market sy amount is zero",
	784:    "Market expired",
	785:    "Market liquidity too low",
	786:    "Market exchange rate negative",
	787:    "Market proportion too high",
	788:    "Market proportion cannot be one",
	789:    "Market exchange rate cannot be one",
	790:    "Insufficient liquidity in the pool, please try a smaller amount.",
	791:    "Market burn sy amount is zero",
	792:    "Market burn pt amount is zero",
	793:    "Insufficient liquidity in the pool, please try a smaller amount.",
	800:    "Market rate scalar negative",
	801:    "Market insufficient sy for swap",
	802:    "Repay sy in exceeds expected sy in",
	803:    "Market insufficient sy in for swap yt",
	804:    "Swapped sy borrowed amount not equal",
	805:    "Market cap exceeded",
	806:    "Invalid repay",
	807:    "Register sy invalid sender",
	808:    "Sy not supported",
	809:    "Register sy type already registered",
	816:    "Register sy type not registered",
	817:    "Sy insufficient repay",
	818:    "Factory invalid py",
	819:    "Invalid py amount",
	820:    "Market insufficient pt in for mint lp",
	821:    "Market invalid py state",
	822:    "Market invalid market position",
	823:    "Market lp amount is zero",
	824:    "Market insufficient lp for burn",
	825:    "Market insufficient yt balance swap",
	832:    "Invalid flash loan position",
	833:    "Create market invalid sender",
	834:    "Invalid epoch",
	835:    "Swap exact yt amount mismatch",
	836:    "Insufficient lp output",
	837:    "Price fluctuation too large",
	1025:   "Acl invalid permission",
	1026:   "Acl role already exists",
	1027:   "Acl role not exists",
	1028:   "Version mismatch error",
	1029:   "Update config invalid sender",
	1030:   "Withdraw from treasury invalid sender",
	1031:   "Invalid yt approx out",
	1032:   "Invalid sy approx out",
	1033:   "Wrong slippage tolerance",
	65537:  "Denominator error",
	65542:  "Abort code on calculation result is negative",
	131074: "The quotient value would be too large to be held in a u128",
	131075: "The multiplied value would be too large to be held in a u128",
	65540:  "A division by zero was encountered",
	131077: "The computed ratio when converting to a FixedPoint64 would be unrepresentable",
}

type ErrorResponse struct {
	Error  string `json:"error"`
	Detail string `json:"detail"`
}

func GetErrorMessage(errorCode int, errorString string) string {
	if msg, ok := ErrorMapping[errorCode]; ok {
		return msg
	}
	return errorString
}

func ParseErrorMessage(errorString string) ErrorResponse {
	if strings.Contains(errorString, "OUT_OF_GAS") {
		return ErrorResponse{
			Error:  "Insufficient liquidity in the pool.",
			Detail: "",
		}
	}

	re := regexp.MustCompile(`[^\d]*(\d+)\)`)
	matches := re.FindStringSubmatch(errorString)

	var errorCode int
	if len(matches) > 1 {
		errorCode, _ = strconv.Atoi(matches[1])
	} else if len(matches) > 0 {
		errorCode64, _ := strconv.ParseInt(matches[0], 16, 64)
		errorCode = int(errorCode64)
	}

	detail := ""
	if errorCode == 790 || errorCode == 793 {
		detail = "To ensure the capital efficiency of the liquidity pool, Nemo's flash swap is utilized when selling YT, which requires higher liquidity. You can try swapping again later or reduce the selling amount."
	}

	err := errorString
	if errorCode != 0 {
		err = GetErrorMessage(errorCode, errorString)
	}

	return ErrorResponse{
		Error:  err,
		Detail: detail,
	}
}
