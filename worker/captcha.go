package worker

import (
	"errors"
	"strings"
)

type CaptchaSolver interface {
	recaptcha(pageUrl string, siteKey string, dataS string) (gRecaptchaResponse string, err error)
}

var captchaSolvers []CaptchaSolver
var captchaSolversCounter = 0

func CaptchaSolverInit(solvers []CaptchaSolver) {
	captchaSolvers = solvers
}

func isRecaptcha(body string) (is bool, siteKey string, dataS string, q string) {
	const keyword = `class="g-recaptcha"`
	kwPos := strings.Index(body, keyword)
	if kwPos < 0 {
		is = false
		return
	}

	const skStartKw = `data-sitekey="`
	const skEndKw = `"`
	skStart := strings.Index(body[kwPos:], skStartKw)
	if skStart < 0 {
		is = false
		return
	}
	skStart += kwPos + len(skStartKw)
	skEnd := strings.Index(body[skStart:], skEndKw)
	if skEnd < 0 {
		is = false
		return
	}
	is = true
	siteKey = body[skStart : skStart+skEnd]

	const dsStartKw = `data-s="`
	const dsEndKw = `"></div>`
	dsStart := strings.Index(body[kwPos:], dsStartKw)
	if dsStart < 0 {
		return
	}
	dsStart += kwPos + len(dsStartKw)
	dsEnd := strings.Index(body[dsStart:], dsEndKw)
	if dsEnd < 0 {
		return
	}
	dataS = body[dsStart : dsStart+dsEnd]

	const qStartKw = `name='q' value='`
	const qEndKw = `'>`
	qStart := strings.Index(body[kwPos:], qStartKw)
	if qStart < 0 {
		return
	}
	qStart += kwPos + len(qStartKw)
	qEnd := strings.Index(body[qStart:], qEndKw)
	if qEnd < 0 {
		return
	}
	q = body[qStart : qStart+qEnd]
	return
}

func solveRecaptcha(pageUrl string, siteKey string, dataS string) (gRecaptchaResponse string, err error) {
	if len(captchaSolvers) > 0 {
		gRecaptchaResponse, err = captchaSolvers[captchaSolversCounter%len(captchaSolvers)].recaptcha(pageUrl, siteKey, dataS)
		captchaSolversCounter++
	} else {
		err = errors.New("captcha: no solving service provided")
	}
	return
}
