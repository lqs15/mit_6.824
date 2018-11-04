package mapreduce

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
)

func doReduce(
	jobName string, // the name of the whole MapReduce job
	reduceTask int, // which reduce task this is
	outFile string, // write the output here
	nMap int, // the number of map tasks that were run ("M" in the paper)
	reduceF func(key string, values []string) string,
) {
	//
	// doReduce manages one reduce task: it should read the intermediate
	// files for the task, sort the intermediate key/value pairs by key,
	// call the user-defined reduce function (reduceF) for each key, and
	// write reduceF's output to disk.
	// 读中间文件
	// You'll need to read one intermediate file from each map task;
	// reduceName(jobName, m, reduceTask) yields the file
	// name from map task m.
	//
	// Your doMap() encoded the key/value pairs in the intermediate
	// files, so you will need to decode them. If you used JSON, you can
	// read and decode by creating a decoder and repeatedly calling
	// .Decode(&kv) on it until it returns an error.
	//
	// You may find the first example in the golang sort package
	// documentation useful.
	//
	// reduceF() is the application's reduce function. You should
	// call it once per distinct key, with a slice of all the values
	// for that key. reduceF() returns the reduced value for that key.
	//
	// You should write the reduce output as JSON encoded KeyValue
	// objects to the file named outFile. We require you to use JSON
	// because that is what the merger than combines the output
	// from all the reduce tasks expects. There is nothing special about
	// JSON -- it is just the marshalling format we chose to use. Your
	// output code will look something like this:
	//
	// enc := json.NewEncoder(file)
	// for key := ... {
	// 	enc.Encode(KeyValue{key, reduceF(...)})
	// }
	// file.Close()
	//
	// Your code here (Part I).
	//
	//fmt.Println(jobName, reduceTask, outFile, nMap)
	//duration := time.Duration(2) * time.Second
	//time.Sleep(duration)
	//NOTE::这里有个问题，sequential的时候，nMap传过来为1，reduceTask 却为0
	//在创建文件*_map.go 时，nmap为0， reduce task却为1， 导致输入和输出的文件不一样
	//FIX：这里是我的理解的问题，题目中nMap, 我理解 成的被写的文件，实际上应该是，有多少要写入的文件
	for nmap := 0; nmap < nMap; nmap++ {
		readfile := reduceName(jobName, nmap, reduceTask)
		readfileHandler, err := os.OpenFile(readfile, os.O_RDONLY, 0664)
		defer readfileHandler.Close()
		content, _ := ioutil.ReadAll(readfileHandler)
		if err != nil {
			panic(err)
		}
		//直接写入文件
		//从文件中读取，花了一点功夫在GO的语法上面
		var kvList []KeyValue

		json.Unmarshal(content, &kvList)
		sort.Slice(kvList, func(a, b int) bool {
			i, err := strconv.ParseInt(kvList[a].Key, 10, 32)
			j, err := strconv.ParseInt(kvList[b].Key, 10, 32)
			if err != nil {
				panic(err)
			}
			return i < j
		})
		writeTofileHandler, err := os.OpenFile(outFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
		defer writeTofileHandler.Close()
		if err != nil {
			panic(err)
		}

		encoderHandler := json.NewEncoder(writeTofileHandler)
		for _, val := range kvList {
			var valList []string
			valList = append(valList, val.Value)
			encoderHandler.Encode(KeyValue{val.Key, reduceF(val.Key, valList)})
		}
	}
}
