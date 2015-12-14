/**
 * The MIT License (MIT)
 *
 * Copyright (c) 2015 Edward Kim <edward@webconn.me>
 * Copyright (c) 2015 Jane Lee <jane@webconn.me>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package gpio

import (
	"os"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	IN = "in"
	OUT = "out"
)

const (
	LOW = 0
	HIGH = 1
)

type Gpio struct {
	Pin int
	Dir string
}

func (g *Gpio) Open() error {

	err := ioutil.WriteFile("/sys/class/gpio/export", []byte(strconv.Itoa(g.Pin)), 0200)
	if err != nil {
		return err
	}

	dirFile := fmt.Sprintf("/sys/class/gpio/gpio%d/direction", g.Pin)
	_, err = os.Stat(dirFile)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dirFile, []byte(g.Dir), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (g *Gpio) Close() error {
	err := ioutil.WriteFile("/sys/class/gpio/unexport", []byte(strconv.Itoa(g.Pin)), 0200)
	if err != nil {
		return err
	}
	return nil
}

func (g *Gpio) Out(value int) error {
	valueFile := fmt.Sprintf("/sys/class/gpio/gpio%d/value", g.Pin)
	_, err := os.Stat(valueFile)
	if err != nil {
		return err
	}

	var v []byte
	if value == 0 {
		v = []byte("0")
	} else {
		v = []byte("1")
	}
	err = ioutil.WriteFile(valueFile, []byte(v), 0644)
	if err != nil {
		return err
	}
	
	return nil
}

func (g *Gpio) In() (int, error) {
	valueFile := fmt.Sprintf("/sys/class/gpio/gpio%d/value", g.Pin)
	_, err := os.Stat(valueFile)
	if err != nil {
		return -1, err
	}

	buf, err := ioutil.ReadFile(valueFile)
	if err != nil {
		return -1, err
	}

	i, err := strconv.Atoi(strings.Trim(string(buf),"\n"))
	if err != nil {
		return -1, err
	}

	return i, nil
}
