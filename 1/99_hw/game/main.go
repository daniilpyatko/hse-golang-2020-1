package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// PickMethod -- Способ подобрать предмет
type PickMethod struct {
	Name          string
	SuccesfulText string
}

// Inventory -- Инвентарь
type Inventory struct {
	Name         string
	Pick         *PickMethod
	PickFunc     func()
	Dependencies []string
	Unsuccessful []string
}

//Place -- Места
type Place struct {
	Name             string
	PresentInventory []Inventory
}

// Quest -- квесты, которые нужно выполнить
type Quest struct {
	Name string
	// Зависимости -- названия объектов, которые должны быть в инвентаре, чтобы пройти квест
	Dependencies []string
	// Функция, которая вызывается, когда выполнен квест
	Complete func() string
}

// Location -- Локации
type Location struct {
	Name              string
	Introduction      string
	Objective         string
	ReasonLocked      string
	Greeting          string
	IsRoom            bool
	Unlocked          bool
	Places            []Place
	AdjacentLocations []*Location
}

var CurrentLocation *Location
var CurrentInventory []Inventory

var AllQuests map[string]*Quest
var AllQuestsArray []Quest
var AllLocations map[string]*Location
var AllLocationsArray []Location
var AllPickMethods map[string]*PickMethod
var AllPickMethodsArray []PickMethod

func main() {
	initGame()
	reader := bufio.NewReader(os.Stdin)
	text := ""
	for text != "stop" {
		text, _ = reader.ReadString('\n')
		text = text[:len(text)-2]
		fmt.Println(handleCommand(text))
	}
}

func initGame() {
	AllPickMethodsArray = []PickMethod{
		PickMethod{
			Name:          "взять",
			SuccesfulText: "предмет добавлен в инвентарь",
		},
		PickMethod{
			Name:          "надеть",
			SuccesfulText: "вы надели",
		},
	}

	AllPickMethods = map[string]*PickMethod{
		"взять":  &AllPickMethodsArray[0],
		"надеть": &AllPickMethodsArray[1],
	}
	AllLocationsArray = []Location{
		Location{
			Name:         "кухня",
			Introduction: "ты находишься на кухне",
			Objective:    "надо собрать рюкзак и идти в универ",
			Greeting:     "кухня, ничего интересного",
			IsRoom:       true,
			Unlocked:     true,
			Places: []Place{
				Place{
					Name: "столе",
					PresentInventory: []Inventory{
						Inventory{
							Name:     "чай",
							PickFunc: func() {},
						},
					},
				},
			},
		},
		Location{
			Name:         "комната",
			Introduction: "",
			Objective:    "",
			Greeting:     "ты в своей комнате",
			IsRoom:       true,
			Unlocked:     true,
			Places: []Place{
				Place{
					Name: "столе",
					PresentInventory: []Inventory{
						Inventory{
							Name: "ключи",
							Pick: AllPickMethods["взять"],
							Dependencies: []string{
								"рюкзак",
							},
							Unsuccessful: []string{
								"некуда класть",
							},
							PickFunc: func() {},
						},
						Inventory{
							Name: "конспекты",
							Pick: AllPickMethods["взять"],
							Dependencies: []string{
								"рюкзак",
							},
							Unsuccessful: []string{
								"некуда класть",
							},
							PickFunc: func() {},
						},
					},
				},
				Place{
					Name: "стуле",
					PresentInventory: []Inventory{
						Inventory{
							Name:     "рюкзак",
							Pick:     AllPickMethods["надеть"],
							PickFunc: func() {},
						},
					},
				},
			},
		},
		Location{
			Name:     "коридор",
			Greeting: "ничего интересного",
			Unlocked: true,
			IsRoom:   false,
		},
		Location{
			Name:         "улица",
			ReasonLocked: "дверь закрыта",
			Greeting:     "на улице весна",
			Unlocked:     false,
			IsRoom:       false,
		},
		Location{
			Name:     "домой",
			Unlocked: true,
			IsRoom:   false,
		},
	}
	AllLocations = map[string]*Location{
		"кухня":   &AllLocationsArray[0],
		"комната": &AllLocationsArray[1],
		"коридор": &AllLocationsArray[2],
		"улица":   &AllLocationsArray[3],
		"домой":   &AllLocationsArray[4],
	}
	AllLocations["комната"].Places[1].PresentInventory[0].PickFunc = func() {
		AllLocations["кухня"].Objective = "надо идти в универ"
	}
	AllLocations["кухня"].AdjacentLocations = append(AllLocations["кухня"].AdjacentLocations, AllLocations["коридор"])
	AllLocations["коридор"].AdjacentLocations = append(AllLocations["коридор"].AdjacentLocations, AllLocations["кухня"])
	AllLocations["комната"].AdjacentLocations = append(AllLocations["комната"].AdjacentLocations, AllLocations["коридор"])
	AllLocations["коридор"].AdjacentLocations = append(AllLocations["коридор"].AdjacentLocations, AllLocations["комната"])
	AllLocations["коридор"].AdjacentLocations = append(AllLocations["коридор"].AdjacentLocations, AllLocations["улица"])
	AllLocations["улица"].AdjacentLocations = append(AllLocations["улица"].AdjacentLocations, AllLocations["домой"])

	AllQuestsArray = append(AllQuestsArray, Quest{
		Name: "дверь",
		Dependencies: []string{
			"ключи",
		},
		Complete: func() string {
			AllLocations["улица"].Unlocked = true
			return "дверь открыта"
		},
	})
	AllQuests = map[string]*Quest{
		"дверь": &AllQuestsArray[0],
	}
	CurrentInventory = CurrentInventory[:0]
	CurrentLocation = AllLocations["кухня"]
}

func getAvailableInventory() string {
	if !CurrentLocation.IsRoom {
		return ""
	}
	Ar := []string{}
	IsPlace := []bool{}
	Res := ""
	IsEmpty := true
	for _, CurPlace := range CurrentLocation.Places {
		Ar = append(Ar, CurPlace.Name)
		IsPlace = append(IsPlace, true)
		for _, CurInventory := range CurPlace.PresentInventory {
			Ar = append(Ar, CurInventory.Name)
			IsPlace = append(IsPlace, false)
			IsEmpty = false
		}
	}
	ArFixed := []string{}
	IsPlaceFixed := []bool{}
	for i := 0; i < len(Ar); {
		if i+1 < len(Ar) && !IsPlace[i+1] {
			ArFixed = append(ArFixed, Ar[i])
			IsPlaceFixed = append(IsPlaceFixed, IsPlace[i])
			for i+1 < len(Ar) && !IsPlace[i+1] {
				ArFixed = append(ArFixed, Ar[i+1])
				IsPlaceFixed = append(IsPlaceFixed, IsPlace[i+1])
				i++
			}
		}
		i++
	}
	for i := 0; i < len(ArFixed); i++ {
		if IsPlaceFixed[i] {
			Res += "на " + ArFixed[i] + ": "
		} else {
			Res += ArFixed[i]
			if i != len(ArFixed)-1 {
				Res += ", "
			}
		}
	}
	if IsEmpty {
		return "пустая " + CurrentLocation.Name
	}
	return Res
}

func getAdjacentLocations() string {
	Res := "можно пройти - "
	for index, CurLocation := range CurrentLocation.AdjacentLocations {
		Res += CurLocation.Name
		if index != len(CurrentLocation.AdjacentLocations)-1 {
			Res += ", "
		}
	}
	return Res
}

func handleCheck() string {
	Res := ""
	if CurrentLocation.Introduction != "" {
		Res += CurrentLocation.Introduction + ", "
	}
	Res += getAvailableInventory()
	if CurrentLocation.Objective != "" {
		Res += ", "
		Res += CurrentLocation.Objective + ". "
	} else {
		if Res != "" {
			Res += ". "
		}
	}
	Res += getAdjacentLocations()
	return Res
}

func handleGo(Where string) string {
	Reachable := false
	var NewLocation *Location
	for _, CurLocation := range CurrentLocation.AdjacentLocations {
		if CurLocation.Name == Where {
			Reachable = true
			NewLocation = CurLocation
			break
		}
	}
	if !Reachable {
		return "нет пути в " + Where
	}
	if !NewLocation.Unlocked {
		return NewLocation.ReasonLocked
	}
	CurrentLocation = NewLocation
	Res := CurrentLocation.Greeting + ". "
	Res += getAdjacentLocations()
	return Res
}

func isPick(Candidate string) bool {
	for _, CurPick := range AllPickMethods {
		if Candidate == CurPick.Name {
			return true
		}
	}
	return false
}

func deleteElement(a []Inventory, index int) []Inventory {
	copy(a[index:], a[index+1:])
	a = a[:len(a)-1]
	return a
}

func tryPick(TryPick string, TryInventory string) string {
	Found := false
	var FoundInventory int
	var FoundPlace int

	for index, CurPlace := range CurrentLocation.Places {
		for index1, CurInventory := range CurPlace.PresentInventory {
			if TryInventory == CurInventory.Name {
				Found = true
				FoundInventory = index1
				FoundPlace = index
				if CurInventory.Pick.Name != TryPick {
					panic("Wrong pick method applied")
				}
			}
		}
	}
	if !Found {
		return "нет такого"
	}
	FoundInventoryRef := &CurrentLocation.Places[FoundPlace].PresentInventory[FoundInventory]
	for index, CurNeeded := range FoundInventoryRef.Dependencies {
		FoundDependency := false
		for _, CurInv := range CurrentInventory {
			if CurInv.Name == CurNeeded {
				FoundDependency = true
			}
		}
		if !FoundDependency {
			return FoundInventoryRef.Unsuccessful[index]
		}
	}
	FoundInventoryRef.PickFunc()
	CurrentInventory = append(CurrentInventory, CurrentLocation.Places[FoundPlace].PresentInventory[FoundInventory])
	CurrentLocation.Places[FoundPlace].PresentInventory = deleteElement(CurrentLocation.Places[FoundPlace].PresentInventory, FoundInventory)
	Res := CurrentInventory[len(CurrentInventory)-1].Pick.SuccesfulText + ": " + TryInventory
	return Res
}

func tryApply(What string, Where string) string {
	FoundInventory := false
	for _, CurInventory := range CurrentInventory {
		if CurInventory.Name == What {
			FoundInventory = true
		}
	}
	if !FoundInventory {
		return "нет предмета в инвентаре - " + What
	}
	_, present := AllQuests[Where]
	if !present {
		return "не к чему применить"
	}
	for _, NeedInventory := range AllQuests[Where].Dependencies {
		CurNeededInventory := false
		for _, CurInventory := range CurrentInventory {
			if CurInventory.Name == NeedInventory {
				CurNeededInventory = true
			}
		}
		if !CurNeededInventory {
			panic("Not everything present")
		}
	}
	Res := AllQuests[Where].Complete()
	return Res
}

func handleCommand(command string) string {
	commands := strings.Fields(command)
	if commands[0] == "осмотреться" {
		return handleCheck()
	} else if commands[0] == "идти" {
		return handleGo(commands[1])
	} else if isPick(commands[0]) {
		return tryPick(commands[0], commands[1])
	} else if commands[0] == "применить" {
		return tryApply(commands[1], commands[2])
	} else {
		return "неизвестная команда"
	}
}
