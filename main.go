package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/cpf2021-gif/grpclient/utils"
)

var (
	args        string
	addr        string
	service     string
	function    string
	serviceList []string
	scv2func    map[string][]string = make(map[string][]string)
)

func main() {
	// TODO: add config option
	getAddr()

	// get all services
	getServiceList()

	// get all functions
	getFunctionList()

	// generate form
	form := huh.NewForm(
		generateGroup()...,
	)

	// select service and function
	err := form.Run()
	if err != nil {
		log.Fatalln(err)
	}

	// call grpcurl
	cmd := exec.Command("grpcurl", "-plaintext", "-d", args, addr, function)
	output, err := cmd.Output()

	if err != nil {
		log.Fatalln(err)
	}

	// print output
	// TODO: beautify output
	fmt.Print(string(output))
}

func getAddr() {
	list := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Input your address").
				Value(&addr),
		),
	)

	err := list.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func getServiceList() {
	cmd := exec.Command("grpcurl", "-plaintext", addr, "list")
	out, err := cmd.Output()
	if err != nil {
		log.Fatalln(err)
	}

	// filter service
	serviceList = strings.Split(string(out), "\n")
	serviceList = utils.Filter(serviceList, utils.NotEmpty)
	serviceList = utils.Filter(serviceList, func(s string) bool {
		return !strings.HasPrefix(s, "grpc")
	})
}

func getFunctionList() {
	for _, scv := range serviceList {
		// grpcurl -plaintext [addr] describe [service]
		/*
			pb.Greeter is a service:
			service Greeter {
				  rpc SayHello ( .pb.HelloRequest ) returns ( .pb.HelloResponse );
			}
		*/
		cmd := exec.Command("grpcurl", "-plaintext", addr, "describe", scv)
		out, err := cmd.Output()
		if err != nil {
			log.Fatalln(err)
		}

		funcs, err := utils.GetFunction(string(out))
		funcNames := make([]string, 0)
		for _, f := range funcs {
			funcNames = append(funcNames, scv+"."+f.Name)
		}

		if err != nil {
			log.Fatalln(err)
		}
		scv2func[scv] = funcNames
	}
}

func generateGroup() []*huh.Group {
	var groups []*huh.Group

	groups = append(groups,
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose your Service").
				Options(
					huh.NewOptions(serviceList...)...,
				).
				Value(&service),
		))

	for _, scv := range serviceList {
		groups = append(groups, huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose your Function").
				Options(
					huh.NewOptions(scv2func[scv]...)...,
				).
				Value(&function),
		).WithHideFunc(func() bool {
			return scv != service
		}))
	}

	// TODO: funcs -> args
	groups = append(groups,
		huh.NewGroup(
			huh.NewInput().
				Title("input your args").
				Value(&args),
		).WithHideFunc(func() bool {
			return function == ""
		}))

	return groups
}
