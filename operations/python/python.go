package main

import (
	"github.com/sbinet/go-python"
	"log"
)

func initialize() {
	err := python.Initialize()
	if err != nil {
		log.Fatal(err)
	}

}
//
//func PyRotate(self, args *python.PyObject) *python.PyObject {
//}

func CreatePythonModule() (*python.PyObject, error) {

	functions := make([]python.PyMethodDef, 0)



	module, err := python.Py_InitModule("lightsd", functions)



	return module, err
}


func TestPython() {
	module := python.PyImport_ImportModule("plugin")
	if module == nil {
		log.Fatal("lightsd.py not fonud")
	}
	defer module.DecRef()


}

