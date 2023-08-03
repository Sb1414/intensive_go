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

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func isXML(filename string) bool {
	return len(filename) >= 4 && filename[len(filename)-4:] == ".xml"
}

func isJSON(filename string) bool {
	return len(filename) >= 5 && filename[len(filename)-5:] == ".json"
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
	xmlRecipe := &XMLRecipe{} // Use '=' here instead of ':='

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

func convertToJSON(recipe *XMLRecipe) (*JSONRecipe, error) {
	jsonRecipe := &JSONRecipe{} // Use '=' here instead of ':='

	for _, xcake := range recipe.Cakes {
		jsonCake := JSONCake{
			Name: xcake.Name,
			Time: xcake.StoveTime,
		}

		for _, xitem := range xcake.Ingredients {
			jsonItem := JSONItem{
				IngredientName:  xitem.ItemName,
				IngredientCount: xitem.ItemCount,
				IngredientUnit:  xitem.ItemUnit,
			}
			jsonCake.Ingredients = append(jsonCake.Ingredients, jsonItem)
		}

		jsonRecipe.Cakes = append(jsonRecipe.Cakes, jsonCake)
	}

	return jsonRecipe, nil
}

func writeXML(xmlRecipe *XMLRecipe, xmlFileName string) error {
	xmlFile, err := os.Create(xmlFileName)
	if err != nil {
		return err
	}
	defer xmlFile.Close()

	xmlData, err := xml.MarshalIndent(xmlRecipe, "", "    ")
	if err != nil {
		return err
	}

	_, err = xmlFile.Write(xmlData)
	if err != nil {
		return err
	}

	fmt.Printf("JSON file successfully converted to XML and saved as '%s'\n", xmlFileName)
	return nil
}

func writeJSON(jsonRecipe *JSONRecipe, jsonFileName string) error {
	jsonFile, err := os.Create(jsonFileName)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonData, err := json.MarshalIndent(jsonRecipe, "", "    ")
	if err != nil {
		return err
	}

	_, err = jsonFile.Write(jsonData)
	if err != nil {
		return err
	}

	fmt.Printf("XML file successfully converted to JSON and saved as '%s'\n", jsonFileName)
	return nil
}

func prettyPrint(data interface{}) {
	prettyJSON, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Printf("Error while printing the data: %v\n", err)
		return
	}
	fmt.Println(string(prettyJSON))
}

func compareDatabases(oldXML *XMLRecipe, oldJSON *JSONRecipe, newXML *XMLRecipe, newJSON *JSONRecipe) {
	// Create maps to store recipes based on their names
	oldRecipes := make(map[string]XMLCake)
	newRecipes := make(map[string]XMLCake)

	// Populate the oldRecipes map using the oldXML recipe data
	for _, cake := range oldXML.Cakes {
		oldRecipes[cake.Name] = cake
	}

	// Populate the newRecipes map using the newXML recipe data
	for _, cake := range newXML.Cakes {
		newRecipes[cake.Name] = cake
	}

	// Check for added and removed cakes
	for name, _ := range oldRecipes {
		if _, ok := newRecipes[name]; !ok {
			fmt.Printf("REMOVED cake \"%s\"\n", name)
		}
	}

	for name, _ := range newRecipes {
		if _, ok := oldRecipes[name]; !ok {
			fmt.Printf("ADDED cake \"%s\"\n", name)
		}
	}

	// Compare cooking times and ingredients for cakes present in both old and new databases
	for _, oldCake := range oldXML.Cakes {
		if newCake, ok := newRecipes[oldCake.Name]; ok {
			if oldCake.StoveTime != newCake.StoveTime {
				fmt.Printf("CHANGED cooking time for cake \"%s\" - \"%s\" instead of \"%s\"\n", oldCake.Name, newCake.StoveTime, oldCake.StoveTime)
			}

			oldIngredients := make(map[string]XMLItem)
			newIngredients := make(map[string]XMLItem)

			// Populate the oldIngredients map using the oldCake ingredients data
			for _, item := range oldCake.Ingredients {
				oldIngredients[item.ItemName] = item
			}

			// Populate the newIngredients map using the newCake ingredients data
			for _, item := range newCake.Ingredients {
				newIngredients[item.ItemName] = item
			}

			// Check for added and removed ingredients
			for name, _ := range oldIngredients {
				if _, ok := newIngredients[name]; !ok {
					fmt.Printf("REMOVED ingredient \"%s\" for cake \"%s\"\n", name, oldCake.Name)
				}
			}

			for name, _ := range newIngredients {
				if _, ok := oldIngredients[name]; !ok {
					fmt.Printf("ADDED ingredient \"%s\" for cake \"%s\"\n", name, oldCake.Name)
				}
			}

			// Compare ingredient unit and unit count for ingredients present in both old and new cakes
			for _, oldItem := range oldCake.Ingredients {
				if newItem, ok := newIngredients[oldItem.ItemName]; ok {
					if oldItem.ItemUnit != newItem.ItemUnit {
						fmt.Printf("CHANGED unit for ingredient \"%s\" for cake \"%s\" - \"%s\" instead of \"%s\"\n", oldItem.ItemName, oldCake.Name, newItem.ItemUnit, oldItem.ItemUnit)
					}
					if oldItem.ItemCount != newItem.ItemCount {
						fmt.Printf("CHANGED unit count for ingredient \"%s\" for cake \"%s\" - \"%s\" instead of \"%s\"\n", oldItem.ItemName, oldCake.Name, newItem.ItemCount, oldItem.ItemCount)
					}
				}
			}
		}
	}
}

func main() {
	oldFileName := flag.String("old", "", "Path to the old XML or JSON file")
	newFileName := flag.String("new", "", "Path to the new XML or JSON file")

	flag.Parse()

	if *oldFileName == "" {
		fmt.Println("Please provide the path to the old XML or JSON file using -old flag.")
		return
	}

	if *newFileName == "" {
		fmt.Println("Please provide the path to the new XML or JSON file using -new flag.")
		return
	}

	if !fileExists(*oldFileName) {
		fmt.Println("The specified file does not exist.")
		return
	}

	var oldXML *XMLRecipe
	var oldJSON *JSONRecipe

	if isXML(*oldFileName) {
		oldXML, err := readXML(*oldFileName)
		if err != nil {
			fmt.Printf("Error while reading the old XML file: %v\n", err)
			return
		}
		prettyPrint(oldXML)

		// Convert XML to JSON
		oldJSON, err = convertToJSON(oldXML)
		if err != nil {
			fmt.Printf("Error while converting old XML to JSON: %v\n", err)
			return
		}
	} else if isJSON(*oldFileName) {
		oldJSON, err := readJSON(*oldFileName)
		if err != nil {
			fmt.Printf("Error while reading the old JSON file: %v\n", err)
			return
		}
		prettyPrint(oldJSON)

		// Convert JSON to XML
		oldXML, err = convertToXML(oldJSON)
		if err != nil {
			fmt.Printf("Error while converting old JSON to XML: %v\n", err)
			return
		}
		prettyPrint(oldXML)
	} else {
		fmt.Println("Invalid old file format. Supported formats are XML and JSON.")
		return
	}

	// Read the new XML or JSON file
	if isXML(*newFileName) {
		newXML, err := readXML(*newFileName)
		if err != nil {
			fmt.Printf("Error while reading the new XML file: %v\n", err)
			return
		}
		prettyPrint(newXML)

		// Convert new XML to JSON
		newJSON, err := convertToJSON(newXML)
		if err != nil {
			fmt.Printf("Error while converting new XML to JSON: %v\n", err)
			return
		}

		// Compare the old and new databases
		compareDatabases(oldXML, oldJSON, newXML, newJSON)
	} else if isJSON(*newFileName) {
		newJSON, err := readJSON(*newFileName)
		if err != nil {
			fmt.Printf("Error while reading the new JSON file: %v\n", err)
			return
		}
		prettyPrint(newJSON)

		// Convert new JSON to XML
		newXML, err := convertToXML(newJSON)
		if err != nil {
			fmt.Printf("Error while converting new JSON to XML: %v\n", err)
			return
		}

		// Compare the old and new databases
		compareDatabases(oldXML, oldJSON, newXML, newJSON)
	} else {
		fmt.Println("Invalid new file format. Supported formats are XML and JSON.")
		return
	}
}
