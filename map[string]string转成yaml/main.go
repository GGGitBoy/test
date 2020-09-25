package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"sigs.k8s.io/yaml"
)

// "k8s.io/apimachinery/pkg/util/yaml"

var row = map[string]string{
	// "id":                 "2",
	// "contact.name.first": "John",
	// "contact.name.last":  "Doe",
	// "contact.email":      "example@gmail.com",
	// "contact.info.me":    "classified",
	// "devices[0].sss":     "mobile",
	// "devices[1].aaa":     "laptop",
	"aaa.vvvv[0]": "oooo",
	"aaa.vvvv[1]": "pppp",
}

var test = map[string]interface{}{
	"id": "2",
	"contact": map[string]interface{}{
		"name": map[string]interface{}{
			"first": "John",
			"last":  "Doe",
		},
		"email": "example@gmail.com",
		"info": map[string]interface{}{
			"me": "classified",
		},
	},
	"devices": []map[string]interface{}{
		{
			"sss": "mobile",
		},
		{
			"aaa": "laptop",
		},
	},
	"aaa": map[string]interface{}{
		"vvvv": []string{"oooo", "pppp"},
	},
}

// "id":2,"contact.name.first":"John","contact.name.last":"Doe",
// "id":2,"contact": {"name": {"first": "John","last": "Doe"},

func main() {

	answer := make(map[string]interface{})
	for k, v := range row {
		answer = get(k, v, answer)
	}

	// fmt.Printf("%v", answer)

	b, err := json.Marshal(answer)
	if err != nil {
		fmt.Println("json.Marshal failed:", err)
		return
	}

	fmt.Println("\n\nb:", string(b))
	////////////
	j := []byte(string(b))
	y, err := yaml.JSONToYAML(j)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Println(string(y))
	/* Output:
	   name: John
	   age: 30
	*/
	j2, err := yaml.YAMLToJSON(y)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Println(string(j2))

}

// func get(key, value string, answer map[string]interface{}) map[string]interface{} {
// 	index1 := strings.Index(key, ".")
// 	index22 := strings.Index(key, "[")
// 	index44 := strings.Index(key, "]")

// 	if index22 > 0 && index44 > 0 {
// 		var arr []map[string]interface{}
// 		fmt.Println(key[:index22])   // devices
// 		fmt.Println(key[index44+1:]) // .sss
// 		test := make(map[string]interface{})
// 		if _, ok := answer[key[:index22]]; ok {
// 			if data, ok2 := answer[key[:index22]].([]map[string]interface{}); ok2 {
// 				arr = data
// 			}
// 		}

// 		arr = append(arr, get(key[index44+1:], value, test))
// 		answer[key[:index22]] = arr

// 	} else if index1 > 0 {
// 		test := make(map[string]interface{})
// 		if _, ok := answer[key[:index1]]; ok {
// 			if data, ok2 := answer[key[:index1]].(map[string]interface{}); ok2 {
// 				test = data
// 			}
// 		}
// 		answer[key[:index1]] = get(key[index1+1:], value, test)
// 	} else if index1 == 0 {
// 		answer[key[index1+1:]] = value
// 	} else {
// 		answer[key] = value
// 	}

// 	return answer
// }

func get(key, value string, answer map[string]interface{}) map[string]interface{} {
	fmt.Printf("%s %s %v\n", key, value, answer)
	indexdot := strings.Index(key, ".")
	indexLeft := strings.Index(key, "[")
	indexRight := strings.Index(key, "]")

	if indexdot > 0 && indexLeft < 0 { // .方法
		answerdot := make(map[string]interface{})
		if _, ok := answer[key[:indexdot]]; ok {
			if data, ok2 := answer[key[:indexdot]].(map[string]interface{}); ok2 {
				answerdot = data
			}
		}
		answer[key[:indexdot]] = get(key[indexdot+1:], value, answerdot)
	} else if indexdot < 0 && indexLeft > 0 { // 【方法
		if key[indexRight+1:] == "" {
			var list []string
			fmt.Printf("into %v\n", answer)
			// if ss, ok := answer[key[:indexLeft]]; ok {
			// 	fmt.Println("gag")
			// 	fmt.Println(ss)
			// if data, ok2 := answer[key[:indexLeft]].([]string); ok2 {
			// 	list = data
			// }
			// }
			// fmt.Println(key[:indexLeft])

			list = append(list, value)
			answer[key[:indexLeft]] = list
		} else {
			var arr []map[string]interface{}
			test := make(map[string]interface{})
			if _, ok := answer[key[:indexLeft]]; ok {
				if data, ok2 := answer[key[:indexLeft]].([]map[string]interface{}); ok2 {
					arr = data
				}
			}
			arr = append(arr, get(key[indexRight+2:], value, test))
			answer[key[:indexLeft]] = arr
		}
	} else if indexdot > 0 && indexLeft > 0 { // 判断大小

		if indexdot < indexLeft {
			answerdot := make(map[string]interface{})
			if _, ok := answer[key[:indexdot]]; ok {
				if data, ok2 := answer[key[:indexdot]].(map[string]interface{}); ok2 {
					answerdot = data
					// fmt.Println("=======")
					// fmt.Println(answerdot)
				}
			}
			answer[key[:indexdot]] = get(key[indexdot+1:], value, answerdot)
		} else if indexdot > indexLeft {
			if key[indexRight+1:] == "" {
				var list []string
				if _, ok := answer[key[:indexLeft]]; ok {
					if data, ok2 := answer[key[:indexLeft]].([]string); ok2 {
						list = data
					}
				}

				list = append(list, value)
				answer[key[:indexLeft]] = list
			} else {
				var arr []map[string]interface{}
				test := make(map[string]interface{})
				if _, ok := answer[key[:indexLeft]]; ok {
					if data, ok2 := answer[key[:indexLeft]].([]map[string]interface{}); ok2 {
						arr = data
					}
				}
				arr = append(arr, get(key[indexRight+2:], value, test))
				answer[key[:indexLeft]] = arr
			}
		}
	} else {
		answer[key] = value
	}

	// if index22 > 0 && index44 > 0 {
	// 	var arr []map[string]interface{}
	// 	test := make(map[string]interface{})
	// 	if _, ok := answer[key[:index22]]; ok {
	// 		if data, ok2 := answer[key[:index22]].([]map[string]interface{}); ok2 {
	// 			arr = data
	// 		}
	// 	}

	// 	arr = append(arr, get(key[index44+2:], value, test))
	// 	answer[key[:index22]] = arr

	// } else if index1 > 0 {
	// 	test := make(map[string]interface{})
	// 	if _, ok := answer[key[:index1]]; ok {
	// 		if data, ok2 := answer[key[:index1]].(map[string]interface{}); ok2 {
	// 			test = data
	// 		}
	// 	}
	// 	answer[key[:index1]] = get(key[index1+1:], value, test)
	// } else if index1 == 0 {
	// 	answer[key[index1+1:]] = value
	// } else {
	// 	answer[key] = value
	// }

	return answer
}
