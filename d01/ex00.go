package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

type XMLRecipe struct {
	XMLName xml.Name  `xml:"recipes"`
	Cakes   []XMLCake `xml:"cake"`
}

type XMLCake struct {
	Name        string    `xml:"name"`
	StoveTime   string    `xml:"stovetime"`
	Ingredients []XMLItem `xml:"ingredients>item"`
}

type XMLItem struct {
	ItemName  string `xml:"itemname"`
	ItemCount string `xml:"itemcount"`
	ItemUnit  string `xml:"itemunit,omitempty"`
}

type JSONRecipe struct {
	Cakes []JSONCake `json:"cake"`
}

type JSONCake struct {
	Name        string     `json:"name"`
	Time        string     `json:"time"`
	Ingredients []JSONItem `json:"ingredients"`
}

type JSONItem struct {
	IngredientName  string `json:"ingredient_name"`
	IngredientCount string `json:"ingredient_count"`
	IngredientUnit  string `json:"ingredient_unit,omitempty"`
}

func readXML(filename string) (*XMLRecipe, error) {
	xmlFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer xmlFile.Close()

	data, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return nil, err
	}

	var recipe XMLRecipe
	err = xml.Unmarshal(data, &recipe)
	if err != nil {
		return nil, err
	}

	return &recipe, nil
}

func readJSON(filename string) (*JSONRecipe, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	data, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var recipe JSONRecipe
	err = json.Unmarshal(data, &recipe)
	if err != nil {
		return nil, err
	}

	return &recipe, nil
}

func convertToXML(recipe *JSONRecipe) (*XMLRecipe, error) {
	xmlRecipe := &XMLRecipe{}

	for _, jcake := range recipe.Cakes {
		xmlCake := XMLCake{
			Name:      jcake.Name,
			StoveTime: jcake.Time,
		}

		for _, jitem := range jcake.Ingredients {
			xmlItem := XMLItem{
				ItemName:  jitem.IngredientName,
				ItemCount: jitem.IngredientCount,
				ItemUnit:  jitem.IngredientUnit,
			}
			xmlCake.Ingredients = append(xmlCake.Ingredients, xmlItem)
		}

		xmlRecipe.Cakes = append(xmlRecipe.Cakes, xmlCake)
	}

	return xmlRecipe, nil
}

func main() {
	xmlFileName := flag.String("f", "", "Path to the XML or JSON file")
	flag.Parse()

	if *xmlFileName == "" {
		fmt.Println("Please provide the path to the XML or JSON file using -f flag.")
		return
	}

	if !fileExists(*xmlFileName) {
		fmt.Println("The specified file does not exist.")
		return
	}

	var recipe interface{}
	var err error

	if isXML(*xmlFileName) {
		recipe, err = readXML(*xmlFileName)
	} else if isJSON(*xmlFileName) {
		recipe, err = readJSON(*xmlFileName)
	} else {
		fmt.Println("Invalid file format. Supported formats are XML and JSON.")
		return
	}

	if err != nil {
		fmt.Printf("Error while reading the file: %v\n", err)
		return
	}

	prettyPrint(recipe)

	if isJSON(*xmlFileName) {
		// Convert JSON to XML
		if jsonRecipe, ok := recipe.(*JSONRecipe); ok {
			xmlRecipe, err := convertToXML(jsonRecipe)
			if err != nil {
				fmt.Printf("Error while converting JSON to XML: %v\n", err)
				return
			}
			prettyPrint(xmlRecipe)
		} else {
			fmt.Println("Invalid JSON format.")
		}
	}
}

func prettyPrint(data interface{}) {
	prettyJSON, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Printf("Error while printing the data: %v\n", err)
		return
	}
	fmt.Println(string(prettyJSON))
}

func isXML(filename string) bool {
	return len(filename) >= 4 && filename[len(filename)-4:] == ".xml"
}

func isJSON(filename string) bool {
	return len(filename) >= 5 && filename[len(filename)-5:] == ".json"
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
