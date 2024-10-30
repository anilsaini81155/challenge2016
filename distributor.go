package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type PermissionType int

const (
	Include PermissionType = iota
	Exclude
)

type Permission struct {
	Type   PermissionType
	Region string
}

type Distributor struct {
	Name            string
	Permissions     []Permission
	SubDistributors []*Distributor
}

var distributors map[string]*Distributor
var cityData map[string]City

type City struct {
	Code         string
	ProvinceCode string
	CountryCode  string
	CityName     string
	ProvinceName string
	CountryName  string
}

func loadCities() {
	cityData = make(map[string]City)
	file, err := os.Open("cities.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, _ = reader.Read() // Skip header row

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		city := City{
			Code:         record[0],
			ProvinceCode: record[1],
			CountryCode:  record[2],
			CityName:     record[3],
			ProvinceName: record[4],
			CountryName:  record[5],
		}
		cityData[city.Code] = city
	}
}

func createDistributor(name string) *Distributor {
	if distributors == nil {
		distributors = make(map[string]*Distributor)
	}
	distributor := &Distributor{Name: name}
	distributors[name] = distributor
	return distributor
}

func addPermissions(distributor *Distributor, permissionType PermissionType, regions []string) {
	for _, region := range regions {
		distributor.Permissions = append(distributor.Permissions, Permission{
			Type:   permissionType,
			Region: strings.TrimSpace(region),
		})
	}
}

func linkDistributor(parentName, childName string) error {
	parent, parentExists := distributors[parentName]
	if !parentExists {
		return fmt.Errorf("parent distributor %s does not exist", parentName)
	}
	child, childExists := distributors[childName]
	if !childExists {
		return fmt.Errorf("child distributor %s does not exist", childName)
	}

	// Validate that the child permissions are within the parent permissions
	for _, childPerm := range child.Permissions {
		if childPerm.Type == Include && isRegionExcluded(parent, childPerm.Region) {
			return fmt.Errorf("cannot link %s to %s: %s is excluded in parent distributor", childName, parentName, childPerm.Region)
		}
	}

	parent.SubDistributors = append(parent.SubDistributors, child)
	fmt.Printf("Successfully linked %s to %s\n", childName, parentName)
	return nil
}

func isRegionExcluded(distributor *Distributor, region string) bool {
	for _, perm := range distributor.Permissions {
		if perm.Type == Exclude && perm.Region == region {
			return true
		}
	}
	return false
}

// showPermissionsWithInheritance recursively collects permissions, including from parent distributors
// func showPermissionsWithInheritance(distributor *Distributor) {
// 	fmt.Printf("Permissions for %s:\n", distributor.Name)
// 	effectivePermissions := gatherEffectivePermissions(distributor)
// 	for _, perm := range effectivePermissions {
// 		prefix := "INCLUDE"
// 		if perm.Type == Exclude {
// 			prefix = "EXCLUDE"
// 		}
// 		fmt.Printf("%s: %s\n", prefix, perm.Region)
// 	}
// }

func gatherEffectivePermissions(distributor *Distributor) []Permission {
	permissions := make([]Permission, 0)
	visited := make(map[string]PermissionType)

	// Recursive function to gather permissions from parent distributors only
	var gatherFromParents func(d *Distributor)
	gatherFromParents = func(d *Distributor) {
		for _, perm := range d.Permissions {
			// Apply exclusion if it overrides an existing include
			if existingType, exists := visited[perm.Region]; !exists || (perm.Type == Exclude && existingType == Include) {
				visited[perm.Region] = perm.Type
			}
		}
		// Continue up the parent chain (don't include children)
		for _, parent := range d.SubDistributors {
			gatherFromParents(parent)
		}
	}

	// Gather own and parent permissions (without children)
	gatherFromParents(distributor)

	// Convert map to Permission slice
	for region, permType := range visited {
		permissions = append(permissions, Permission{
			Type:   permType,
			Region: region,
		})
	}

	return permissions
}

func showEffectivePermissions(distributor *Distributor) {
	effectivePermissions := gatherEffectivePermissions(distributor)
	fmt.Printf("Permissions for %s:\n", distributor.Name)
	for _, perm := range effectivePermissions {
		prefix := "INCLUDE"
		if perm.Type == Exclude {
			prefix = "EXCLUDE"
		}
		fmt.Printf("%s: %s\n", prefix, perm.Region)
	}
}
