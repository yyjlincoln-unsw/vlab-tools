package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const PROFILE_PATH = "/home/z5423219/local-public/internals/profiles.json"

var PROFILES map[string]Profile = (func() map[string]Profile {
	var profiles map[string]Profile = map[string]Profile{}
	content, err := ioutil.ReadFile(PROFILE_PATH)
	if err != nil {
		ErrorOutput("Error: Could not read profiles: read: %v", err)
		fmt.Printf("\n")
		return profiles
	}
	if err := json.Unmarshal(content, &profiles); err != nil {
		ErrorOutput("Error: could not read profiles: unmarshal: %v", err)
		fmt.Printf("\n")
		return profiles
	}
	return profiles
})()

type Profile struct {
	Name  string
	Class string
}

func GetProfile(key string) *Profile {
	if val, ok := PROFILES[key]; ok {
		return &val
	}
	return nil
}

func CheckProfileForCourseInfo() *CourseInformation {
	profile := GetProfile(GetCurrentUser())
	if profile == nil {
		return nil
	}
	class := profile.Class
	for _, v := range AUTOTEST_MAP {
		if v.Identifier == class {
			fmt.Printf("The current course for %v is %v (%v)\n", profile.Name, v.CourseName, v.Identifier)
			return v
		}
	}
	return nil
}
