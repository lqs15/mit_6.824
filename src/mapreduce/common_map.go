package mapreduce

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"os"
)

func doMap(
	jobName string, // the name of the MapReduce job
	mapTask int, // which map task this is
	inFile string,
	nReduce int, // the number of reduce task that will be run ("R" in the paper)
	mapF func(filename string, contents string) []KeyValue,
) {
	//
	// doMap manages one map task: it should read one of the input files
	// (inFile), call the user-defined map function (mapF) for that file's
	// contents, and partition mapF's output into nReduce intermediate files.
	//
	// There is one intermediate file per reduce task. The file name
	// includes both the map task number and the reduce task number. Use
	// the filename generated by reduceName(jobName, mapTask, r)
	// as the intermediate file for reduce task r. Call ihash() (see
	// below) on each key, mod nReduce, to pick r for a key/value pair.
	//
	// 中间文件名: map task number+R task number  can be generated by reduceName
	// ihash 挑选reduce number
	//
	// mapF() is the map function provided by the application. The first
	// argument should be the input file name, though the map function
	// typically ignores it. The second argument should be the entire
	// input file contents. mapF() returns a slice containing the
	// key/value pairs for reduce; see common.go for the definition of
	// KeyValue.
	//
	// Look at Go's ioutil and os packages for functions to read
	// and write files.
	//
	// Coming up with a scheme for how to format the key/value pairs on
	// disk can be tricky, especially when taking into account that both
	// keys and values could contain newlines, quotes, and any other
	// character you can think of.
	//
	// One format often used for serializing data to a byte stream that the
	// other end can correctly reconstruct is JSON. You are not required to
	// use JSON, but as the output of the reduce tasks *must* be JSON,
	// familiarizing yourself with it here may prove useful. You can write
	// out a data structure as a JSON string to a file using the commented
	// code below. The corresponding decoding functions can be found in
	// common_reduce.go.
	//
	//   enc := json.NewEncoder(file)
	//   for _, kv := ... {
	//     err := enc.Encode(&kv)
	//
	// Remember to close the file after you have written all the values!
	//
	// Your code here (Part I).
	//

	//func doMap(
	//	jobName string, // the name of the MapReduce job
	//	mapTask int, // which map task this is
	//	inFile string,
	//	nReduce int, // the number of reduce task that will be run ("R" in the paper)
	//	mapF func(filename string, contents string) []KeyValue,
	//) {
	//读入文件
	//fmt.Println(jobName, mapTask, inFile, nReduce)
	//duration := time.Duration(2) * time.Second
	//time.Sleep(duration)
	content, err := ioutil.ReadFile(inFile)
	if err != nil {
		fmt.Printf("File read failed: ")
		return
	}
	//开始mapF
	keyValueRet := mapF(inFile, string(content[:]))
	//写入文件
	//keyValueLength := len(keyValueRet)
	filenameToWriteList := make([]string, nReduce)
	filenameWritehander := make([]*os.File, nReduce)
	for i := 0; i < nReduce; i++ {
		filenameToWrite := reduceName(jobName, mapTask, i)
		filenameToWriteList[i] = filenameToWrite
		f, err := os.OpenFile(filenameToWrite, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		filenameWritehander[i] = f
		defer f.Close()
	}
	type oneKeyValList []KeyValue
	var kvList = make([]oneKeyValList, nReduce)

	for _, val := range keyValueRet {
		fileind := ihash(string(val.Key)) % nReduce
		kvList[fileind] = append(kvList[fileind], val)
	}

	for ind, val := range kvList {
		valInJson, err := json.MarshalIndent(val, "", "\t")
		if err != nil {
			panic(err)
		}
		filenameWritehander[ind].Write(valInJson)
	}

}

func ihash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32() & 0x7fffffff)
}
