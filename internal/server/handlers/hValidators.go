package handlers

import "net/http"

func validateUpdateMetric(mType string, name string, value string) (errCodePair, bool) {
	if !isTypeValid(mType) {
		return errCodePair{
			code:  http.StatusBadRequest,
			descr: "parsing counter type error\n",
		}, false
	}

	if code, ok := isNameValid(name); !ok {
		if code == http.StatusNotFound {
			return errCodePair{
				code:  code,
				descr: "metric isn't found\n",
			}, false
		} else {
			return errCodePair{
				code:  code,
				descr: "parsing counter name error\n",
			}, false
		}

	}

	if !isValueValid(value) {
		return errCodePair{
			code:  http.StatusBadRequest,
			descr: "parsing counter value error\n",
		}, false
	}

	return errCodePair{}, true
}

func validateGet(mType string, name string) (errCodePair, bool) {
	if !isTypeValid(mType) {
		return errCodePair{
			code:  http.StatusBadRequest,
			descr: "parsing counter type error\n",
		}, false
	}

	if code, ok := isNameValid(name); !ok {
		if code == http.StatusNotFound {
			return errCodePair{
				code:  code,
				descr: "metric isn't found\n",
			}, false
		} else {
			return errCodePair{
				code:  code,
				descr: "parsing counter name error\n",
			}, false
		}

	}

	return errCodePair{}, true
}
