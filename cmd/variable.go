/*
Package cmd Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// getVariableCmd represents the variable command
var getVariableCmd = &cobra.Command{
	Use:   "variable",
	Short: "Get Variables",
	Long: `Get Code Stream Variables by name, project or by id - e.g:

Get by ID
	cs-cli get variable --id 6b7936d3-a19d-4298-897a-65e9dc6620c8
	
Get by Name
	cs-cli get variable --name my-variable
	
Get by Project
	cs-cli get variable --project production`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := ensureTargetConnection(); err != nil {
			log.Fatalln(err)
		}

		response, err := getVariable(id, name, project, exportPath)
		if err != nil {
			log.Fatalln("Unable to get Code Stream Variables: ", err)
		}
		var resultCount = len(response)
		if resultCount == 0 {
			// No results
			log.Warnln("No results found")
		} else if resultCount == 1 {
			// Print the single result
			if exportPath != "" {
				exportVariable(response[0], exportPath)
			}
			PrettyPrint(response[0])
		} else {
			// Print result table
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Id", "Name", "Project", "Type", "Description"})
			for _, c := range response {
				table.Append([]string{c.ID, c.Name, c.Project, c.Type, c.Description})
			}
			table.Render()
		}
	},
}

// getVariableCmd represents the variable command
var createVariableCmd = &cobra.Command{
	Use:   "variable",
	Short: "Create a Variable",
	Long:  `Create a Variable`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := ensureTargetConnection(); err != nil {
			log.Fatalln(err)
		}

		if importPath != "" { // If we are importing a file
			variables := importVariables(importPath)
			for _, value := range variables {
				if project != "" { // If the project is specified update the object
					value.Project = project
				}
				createResponse, err := createVariable(value.Name, value.Description, value.Type, value.Project, value.Value)
				if err != nil {
					log.Println("Unable to create Code Stream Variable: ", err)
				} else {
					log.Println("Created variable", createResponse.Name, "in", createResponse.Project)
				}
			}
		} else {
			createResponse, err := createVariable(name, description, typename, project, value)
			if err != nil {
				log.Println("Unable to create Code Stream Variable: ", err)
			}
			PrettyPrint(createResponse)
		}
	},
}

// updateVariableCmd represents the variable command
var updateVariableCmd = &cobra.Command{
	Use:   "variable",
	Short: "Update a Variable",
	Long:  `Update a Variable`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := ensureTargetConnection(); err != nil {
			log.Fatalln(err)
		}

		if importPath != "" { // If we are importing a file
			variables := importVariables(importPath)
			for _, value := range variables {
				exisitingVariable, err := getVariable("", value.Name, value.Project, "")
				if err != nil {
					log.Println("Update failed - unable to find existing Code Stream Variable", value.Name, "in", value.Project)
				} else {
					_, err := updateVariable(exisitingVariable[0].ID, value.Name, value.Description, value.Type, value.Value)
					if err != nil {
						log.Println("Unable to update Code Stream Variable: ", err)
					} else {
						log.Println("Updated variable", value.Name)
					}
				}
			}
		} else { // Else we are updating using flags
			updateResponse, err := updateVariable(id, name, description, typename, value)
			if err != nil {
				log.Println("Unable to update Code Stream Variable: ", err)
			}
			log.Println("Updated variable", updateResponse.Name)
		}
	},
}

// deleteVariableCmd represents the executions command
var deleteVariableCmd = &cobra.Command{
	Use:   "variable",
	Short: "Delete a Variable",
	Long:  `Delete a Variable.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := ensureTargetConnection(); err != nil {
			log.Fatalln(err)
		}

		response, err := deleteVariable(id)
		if err != nil {
			log.Fatalln("Unable to delete variable: ", err)
		}
		log.Println("Variable with id " + response.ID + " deleted")
	},
}

func init() {
	// Get Variable
	getCmd.AddCommand(getVariableCmd)
	getVariableCmd.Flags().StringVarP(&name, "name", "n", "", "List variable with name")
	getVariableCmd.Flags().StringVarP(&project, "project", "p", "", "List variables in project")
	getVariableCmd.Flags().StringVarP(&id, "id", "i", "", "List variables by id")
	getVariableCmd.Flags().StringVarP(&exportPath, "exportPath", "", "", "Path to export objects - relative or absolute location")
	// Create Variable
	createCmd.AddCommand(createVariableCmd)
	createVariableCmd.Flags().StringVarP(&name, "name", "n", "", "The name of the variable to create")
	createVariableCmd.Flags().StringVarP(&typename, "type", "t", "", "The type of the variable to create REGULAR|SECRET|RESTRICTED")
	createVariableCmd.Flags().StringVarP(&project, "project", "p", "", "The project in which to create the variable")
	createVariableCmd.Flags().StringVarP(&value, "value", "v", "", "The value of the variable to create")
	createVariableCmd.Flags().StringVarP(&description, "description", "d", "", "The description of the variable to create")
	createVariableCmd.Flags().StringVarP(&importPath, "importpath", "i", "", "Path to a YAML file with the variables to import")

	// Update Variable
	updateCmd.AddCommand(updateVariableCmd)
	updateVariableCmd.Flags().StringVarP(&id, "id", "i", "", "ID of the variable to update")
	updateVariableCmd.Flags().StringVarP(&name, "name", "n", "", "Update the name of the variable")
	updateVariableCmd.Flags().StringVarP(&typename, "type", "t", "", "Update the type of the variable REGULAR|SECRET|RESTRICTED")
	updateVariableCmd.Flags().StringVarP(&value, "value", "v", "", "Update the value of the variable ")
	updateVariableCmd.Flags().StringVarP(&description, "description", "d", "", "Update the description of the variable")
	updateVariableCmd.Flags().StringVarP(&importPath, "importpath", "", "", "Path to a YAML file with the variables to import")
	//updateVariableCmd.MarkFlagRequired("id")
	// Delete Variable
	deleteCmd.AddCommand(deleteVariableCmd)
	deleteVariableCmd.Flags().StringVarP(&id, "id", "i", "", "Delete variable by id")
	deleteVariableCmd.MarkFlagRequired("id")
}
