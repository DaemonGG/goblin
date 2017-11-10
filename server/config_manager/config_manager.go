package config_manager

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"goblin"
)

type recipe struct {
	tag              string   // json:"tag"
	topics           []string // json:"topics"
	background_fetch bool     // json:"background_fetch"
}

type topic struct {
	key            string    // json:"key"
	filters        []*Filter // json:"filters"
	connector_type string    // aws, gc, azure...
	// fields_selected []string // json:"fields_selected"
}

type Filter struct {
	name   string
	values []string
}

const (
	// TODO: Change to use gflag
	dataRecipesConfigFileName         = "data_recipes.config.json"
	dataRecipesConfigFilePath         = "../../config/"
	dataRecipesConfigFilePollInterval = 1 // second
)

type ExecPlan func(query_json string) chan goblin.JsonQueryReply

var (
	// data_recipes.config.json is not loaded incessantly, it will only be
	// loaded when the service tag cannot be found in local map.
	data_recipes_lst_modification_time time.Time

	// Number of recipes described in data_recipes.config.json
	num_recipes int

	// A map used to find the execution method accoding the data recipe
	// tag.
	tag_2_exec_map map[string]ExecPlan

	// Record the number of continuous load data recipe config failures.
	load_data_recipe_fail_times int
)

// TODO: Need a cron implemented. or use opensource one.
func ScheduleReloadDataRecipeConfig() {

}

// Load the data recipe config file. For any newly modified tags(recipes),
// generate execution plan (type: ExecPlan) for this tag.
func loadDataRecipeConfig() {
	data_recipe_config_path :=
		dataRecipesConfigFilePath + dataRecipesConfigFileName
	config_file_info, _ := os.Stat(data_recipe_config_path)
	if config_file_info.ModTime().After(
		data_recipes_lst_modification_time) {
		// This file has been changed since last load, so reload.
		var recipes []recipe
		if data, err := ioutil.ReadFile(data_recipe_config_path); err == nil {
			json.Unmarshal(data, &recipes)
		} else {
			load_data_recipe_fail_times += 1
			return
		}
	}
}

// Generate ExecPlans for each data recipe.
func loadDataRecipesInternal(recipes []recipe) {
	for _, one_recipe := range recipes {

	}
}

func Solution(tag string) (ExecPlan, error) {
	// Try to reload this config in case this file is changed.
	//loadDataRecipeConfig()
	if solution, ok := tag_2_exec_map[tag]; ok {
		return solution, nil
	} else {
		log.Printf("%s is not a valid tag\n", tag)
		return nil, errors.New("Invalid Tag.")
	}
}
