// Package names is client of http://code.503web.com/names, provide a stream of
// generated names.
package names

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/redforks/errors"
)

// Kind of names.
type Kind int

const (
	// Person names retrieve from http://code.503web.com/names/name
	Person Kind = iota

	// Product names retrieve from http://code.503web.com/names/product
	Product

	// Address names retrieve from http://code.503web.com/names/address
	Address

	// Firm names retrieve from http://code.503web.com/names/firm
	Firm
)

// How many names retrieved from name service in a batch
const batchSize = 1000

// pump out random names one by one, get names from names service in batches.
type pump struct {
	sync.Mutex

	// buf and cur shares the same storage, cur changes after each .Next() call,
	// cur restores its capacity using buf on re-fetch from names service.
	buf, cur []string

	// url of names service
	url string
}

// newPump create a name pump, retrieve names from name service at url.
func newPump(url string) *pump {
	buf := make([]string, 0, batchSize)
	return &pump{
		buf: buf,
		cur: buf,
		url: url,
	}
}

var httpClient = http.Client{
	Timeout: 10 * time.Second,
}

// retrieve names from name service. Must lock pump before call retrieve().
// retrieve returns an error if name service does not return any name.
func (p *pump) retrieve() error {
	p.cur = p.buf
	resp, err := httpClient.Get(p.url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	names := bytes.Split(content, []byte("\n"))
	if len(names) == 0 {
		return errors.Externalf("retrieve names from %s failed, no names returned", p.url)
	}

	for _, n := range names {
		p.cur = append(p.cur, string(n))
	}
	return nil
}

func (p *pump) Next() (name string, err error) {
	p.Lock()
	if len(p.cur) == 0 {
		if err = p.retrieve(); err != nil {
			return
		}
	}

	name = p.cur[0]
	p.cur = p.cur[1:]
	p.Unlock()
	return
}

var (
	person  = newPump("http://code.503web.com/names/name")
	product = newPump("http://code.503web.com/names/product")
	address = newPump("http://code.503web.com/names/address")
	firm    = newPump("http://code.503web.com/names/firm")
)

// NextPerson returns next random person name.
func NextPerson() (name string, err error) {
	return person.Next()
}

// NextProduct returns next random product name.
func NextProduct() (name string, err error) {
	return product.Next()
}

// NextAddress returns next random address name.
func NextAddress() (name string, err error) {
	return address.Next()
}

// NextFirm returns next random firm name.
func NextFirm() (name string, err error) {
	return firm.Next()
}

// Next returns next generated name of sepecific kind. It is multi goroutine safe.
func Next(kind Kind) (name string, err error) {
	switch kind {
	case Person:
		return NextPerson()
	case Product:
		return NextProduct()
	case Address:
		return NextAddress()
	case Firm:
		return NextFirm()
	}
	return "", errors.New("Unknown names kind")
}
