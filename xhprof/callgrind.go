package xhprof

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var (
	formatSpecPattern    = regexp.MustCompile(`^# callgrind format$`)
	formatVersionPattern = regexp.MustCompile(`^version: 1$`)
	creatorPattern       = regexp.MustCompile(`^creator: .*$`)
	headerPattern        = regexp.MustCompile(`^(\w+):\s*([[:graph:]]+)$`)
	costsPattern         = regexp.MustCompile(`^(?:\d+\s*)+$`)
	positionPattern      = regexp.MustCompile(`^(fl|fi|fn|cfi|cfl|cfn)=\s*(?:\((\d+)\))?\s*([[:graph:]]*)?$`)
	callsPattern         = regexp.MustCompile(`^calls=\s*(\d+)\s*(\d+\s*)+$`)
	emptyPattern         = regexp.MustCompile(`^\s*$`)
)

func ParseCallgrind(rd io.Reader) (*PairCallMap, error) {
	p := NewCallgrindParser(rd)
	return p.parseFile()
}

type CallgrindParser struct {
	scanner   *bufio.Scanner
	headers   map[string]string
	positions map[string]string
	pcMap     *PairCallMap
	lastFn    string
	lastCfn   string
}

func NewCallgrindParser(rd io.Reader) *CallgrindParser {
	p := new(CallgrindParser)
	p.scanner = bufio.NewScanner(rd)
	p.headers = make(map[string]string)
	p.positions = make(map[string]string)
	p.pcMap = NewPairCallMap()
	p.pcMap.NewPairCall("main()")
	p.pcMap.M["main()"].Count = 1

	return p
}

func (p *CallgrindParser) setHeader(k, v string) {
	p.headers[k] = v
}

func (p *CallgrindParser) getOrSetPosName(kind, num, posName string) (name string, err error) {
	name = posName
	if num == "" && name == "" {
		err = errors.New("A position must be defined either with a name or a reference number")
		return
	}

	if name == "" {
		var ok bool
		name, ok = p.positions["fn:"+num]
		if !ok {
			err = errors.New("Position referenced without being defined")
		}
	} else {
		if name == "{main}" {
			name = "main()"
		}

		p.positions["fn:"+num] = name
	}

	return
}

func (p *CallgrindParser) readLine() (text string, eof bool, err error) {
	ok := p.scanner.Scan()
	for ok {
		text = p.scanner.Text()
		if !emptyPattern.MatchString(text) {
			break
		}

		ok = p.scanner.Scan()
	}

	err = p.scanner.Err()
	if !ok {
		eof = true
	}

	return
}

func (p *CallgrindParser) parseFile() (pcMap *PairCallMap, err error) {
	var text string
	var eof bool
	text, eof, err = p.readLine()
	if eof || err != nil {
		return
	}

	if formatSpecPattern.MatchString(text) {
		text, eof, err = p.readLine()
		if eof || err != nil {
			return
		}
	}

	if formatVersionPattern.MatchString(text) {
		text, eof, err = p.readLine()
		if eof || err != nil {
			return
		}
	}

	if creatorPattern.MatchString(text) {
		text, eof, err = p.readLine()
		if eof || err != nil {
			return
		}
	}

	err = p.parsePartData()
	if err != nil {
		return
	}

	if sum, ok := p.headers["summary"]; ok && p.pcMap.M["main()"].WallTime == 0 {
		var wt float64
		wt, err = strconv.ParseFloat(sum, 32)
		if err != nil {
			return
		}

		p.pcMap.M["main()"].WallTime = float32(wt)
	}

	pcMap = p.pcMap

	return
}

func (p *CallgrindParser) parsePartData() (err error) {
	eof := false
	text := p.scanner.Text()
	for !eof && err == nil {
		if headerPattern.MatchString(text) {
			err = p.parseHeader()
		} else if positionPattern.MatchString(text) {
			err = p.parsePosition()
		} else if callsPattern.MatchString(text) {
			err = p.parseCalls()
		} else if costsPattern.MatchString(text) {
			err = p.parseCosts(false)
		} else {
			err = errors.New("PartData is not valid: " + text)
		}

		if err != nil {
			break
		}

		text, eof, err = p.readLine()
	}

	return
}

func (p *CallgrindParser) parseHeader() (err error) {
	text := p.scanner.Text()
	submatches := headerPattern.FindStringSubmatch(text)
	k := strings.TrimSpace(submatches[1])
	v := strings.TrimSpace(submatches[2])

	if k == "events" && v != "Time" {
		err = errors.New("Only Time event is supported")
	} else {
		p.setHeader(submatches[1], submatches[2])
	}

	return
}

func (p *CallgrindParser) parsePosition() (err error) {
	text := p.scanner.Text()
	submatches := positionPattern.FindStringSubmatch(text)
	posType := strings.TrimSpace(submatches[1])
	posNum := strings.TrimSpace(submatches[2])
	posName := strings.TrimSpace(submatches[3])

	// Ignore file information
	if posType != "fn" && posType != "cfn" {
		return
	}

	posName, err = p.getOrSetPosName(posType, posNum, posName)

	if posType == "fn" {
		p.lastFn = posName
		p.lastCfn = ""
	} else if posType == "cfn" {
		p.lastCfn = posName
	}

	if p.lastFn != "" && p.lastCfn != "" {
		p.pcMap.NewPairCall(pairName(p.lastFn, p.lastCfn))
	}

	return nil
}

func (p *CallgrindParser) parseCalls() (err error) {
	text := p.scanner.Text()
	submatches := callsPattern.FindStringSubmatch(text)
	count, err := strconv.Atoi(strings.TrimSpace(submatches[1]))
	if err != nil {
		return
	}

	if p.lastCfn == "" {
		return errors.New("Calls expression encountered without called function being defined")
	}

	p.pcMap.M[pairName(p.lastFn, p.lastCfn)].Count += count
	eof := false
	text, eof, err = p.readLine()
	if eof || err != nil {
		return errors.New("Expected inclusive cost of function call after calls expression")
	}

	if !costsPattern.MatchString(text) {
		return errors.New("Expected inclusive cost of function call after calls expression")
	}

	err = p.parseCosts(true)

	return
}

func (p *CallgrindParser) parseCosts(callCosts bool) (err error) {
	if !callCosts && !(p.lastFn == "main()" && p.lastCfn == "") {
		return
	}

	text := p.scanner.Text()
	match := costsPattern.FindString(text)
	cost, err := strconv.ParseFloat(strings.TrimSpace(strings.Split(match, " ")[1]), 32)
	if err != nil {
		return
	}

	if p.lastFn == "" {
		err = errors.New("Costs expression encountered without function being defined")
		return
	}

	p.pcMap.M[pairName(p.lastFn, p.lastCfn)].WallTime += float32(cost)
	if p.lastFn == "main()" && p.lastCfn != "" {
		p.pcMap.M["main()"].WallTime += float32(cost)
	}

	return
}
