package Config

import (
	"bibletool/basic"
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
)

// write to Config file
func Store(data interface{}, path *basic.OSPaths) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(data)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile(path.Configpath, buffer.Bytes(), 0600)
	if err != nil {
		fmt.Println(err)
	}
}

// load from Config file
func Load(data interface{}, path *basic.OSPaths) {
	if _, err := os.Stat(path.Configpath); err == nil {

		raw, err := ioutil.ReadFile(path.Configpath)
		if err != nil {
			fmt.Println(err)
		}
		buffer := bytes.NewBuffer(raw)
		dec := gob.NewDecoder(buffer)
		err = dec.Decode(data)
		if err != nil {
			fmt.Println(err)
		}
	}
}
