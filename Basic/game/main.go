package main

import "strings"

/**/
import "os"
import "bufio"
import "fmt"

/**/
type Player struct {
	place     string
	equipment []string
	//challenges map[string][]string
}

type Locations struct {
	ways    []string
	objects map[string][]string
	log     []string
	tasks   []string
	door    bool
	say     func() string
	rm      func()
}

func main() {
	/**/
	initGame()

	scanner := bufio.NewScanner(os.Stdin)
	var str string
	for {
		scanner.Scan()
		str = scanner.Text()

		answer := handleCommand(str)
		fmt.Println(answer)
	}
	/**/
}

var player Player
var room Locations
var kitchen Locations
var corridor Locations
var street Locations
var help Locations

var strct map[string]Locations
var put []string
var take []string
var furniture []string

func initGame() {
	player = Player{

		place:     "кухня",
		equipment: []string{},
	}

	room = Locations{
		ways: []string{"коридор"},
		objects: map[string][]string{
			"на столе: ": []string{"ключи", "конспекты"},
			"на стуле: ": []string{"рюкзак"},
		},
		//log: log_go, log_la
		log:   []string{"ты в своей комнате", "пустая комната"},
		tasks: []string{},
		door:  true,
		say: func() string {
			var lenth = 0
			var str = ""
			for _, value := range room.objects {
				lenth = lenth + len(value)
			}

			if lenth == 0 {
				return room.log[1]
			} else {
				return str
			}
		},
		rm: func() {

		},
	}

	kitchen = Locations{
		ways: []string{"коридор"},
		objects: map[string][]string{
			"на столе: ": []string{"чай"},
			"на стуле: ": []string{},
		},
		log:   []string{"кухня, ничего интересного", "ты находишься на кухне,"},
		tasks: []string{"собрать рюкзак", "идти в универ"},
		door:  true,
		say: func() string {
			return kitchen.log[1]
		},
		rm: func() {
			if len(player.equipment) >= 1 {
				kitchen.tasks = Remove(kitchen.tasks, "собрать рюкзак")
			}
			strct["кухня"] = kitchen
		},
	}

	corridor = Locations{
		ways: []string{"кухня", "комната", "улица"},
		objects: map[string][]string{
			"на столе: ": []string{},
			"на стуле: ": []string{},
		},
		log:   []string{"ничего интересного", "ты в коридоре"},
		tasks: []string{},
		door:  true,
		say: func() string {
			return corridor.log[1]
		},
		rm: func() {

		},
	}

	street = Locations{
		ways: []string{"домой"},
		objects: map[string][]string{
			"на столе: ": []string{},
			"на стуле: ": []string{},
		},
		log:   []string{"на улице весна", "ты на улице"},
		tasks: []string{},
		door:  false,
		say: func() string {
			return street.log[1]
		},
		rm: func() {

		},
	}

	strct = map[string]Locations{
		"комната": room,
		"коридор": corridor,
		"кухня":   kitchen,
		"улица":   street,
	}

	furniture = []string{"на столе: ", "на стуле: "}
	put = []string{"рюкзак"}
	take = []string{"чай", "конспекты", "ключи"}

}

func LookAround() string {

	help = strct[player.place]
	var loci []string
	var mission string
	dir := ". можно пройти - "
	var log = help.say()
	strct[player.place].rm()

	if len(strct[player.place].tasks) > 0 {
		mission = ", надо "
		if len(strct[player.place].tasks) == 1 {
			mission = mission + strct[player.place].tasks[0]
		} else if len(strct[player.place].tasks) == 2 {
			mission = mission + strings.Join(strct[player.place].tasks, " и ")
		} else if len(strct[player.place].tasks) > 2 {
			mission = mission + strings.Join(strct[player.place].tasks, ", ")
		}
	}

	for _, val := range furniture {
		if len(help.objects[val]) > 0 {
			if log == "" && len(loci) == 0 {
				loci = append(loci, val+strings.Join(help.objects[val], ", "))
			} else if log == "" {
				loci = append(loci, ", "+val+strings.Join(help.objects[val], ", "))
			} else {
				loci = append(loci, " "+val+strings.Join(help.objects[val], ", "))
			}
		}
	}
	site := strings.Join(loci, "")
	dir = dir + strings.Join(help.ways, ", ")

	return log + site + mission + dir
}

func Go(target string) string {

	dir := ". можно пройти - "
	help = strct[player.place]
	for _, item := range help.ways {
		if item == target {
			help = strct[target]
			if help.door == true {
				log := help.log[0]
				dir = dir + strings.Join(help.ways, ", ")
				player.place = target

				return log + dir
			} else {
				return "дверь закрыта"
			}
		}
	}
	return "нет пути в " + target

}

func PutOn(thing string, cmd string) string {

	help = strct[player.place]
	for key, val := range help.objects {
		for _, v := range val {
			if v == thing {
				if cmd == "надеть" {
					for _, item := range put {
						if item == thing {
							strct[player.place].objects[key] = Remove(strct[player.place].objects[key], thing)
							player.equipment = append(player.equipment, thing)

							return "вы надели: " + thing
						}
					}
				} else if cmd == "взять" {
					for _, item := range take {
						if item == thing {
							if len(player.equipment) > 0 {
								player.equipment = append(player.equipment, thing)
								strct[player.place].objects[key] = Remove(strct[player.place].objects[key], thing)

								return "предмет добавлен в инвентарь: " + thing
							}
							return "некуда класть"
						}
					}
					return "нельзя взять " + thing
				}
			}
		}
	}
	return "нет такого"
}

func Apply(this string, tothis string) string {

	var ok bool
	for _, item := range player.equipment {
		if item == this && tothis == "дверь" {
			street.door = true
			//Почему-то не работала констукция strct["улица"].door = true, пришлось сделать так как ниже
			strct["улица"] = street

			return "дверь открыта"
		} else if item == this {
			ok = true
		}
	}
	if ok {
		return "не к чему применить"
	} else {
		return "нет предмета в инвентаре - " + this
	}

}
func Remove(arr []string, del string) []string {
	arr_new := arr[:0]
	for _, item := range arr {
		if item != del {
			arr_new = append(arr_new, item)
		}
	}
	return arr_new
}
func handleCommand(command string) string {

	std := strings.Split(command, " ")
	//fmt.Println(std)

	if std[0] == "идти" {
		return Go(std[1])
	} else if std[0] == "осмотреться" {
		return LookAround()
	} else if std[0] == "взять" || std[0] == "надеть" {
		return PutOn(std[1], std[0])
	} else if std[0] == "применить" {
		return Apply(std[1], std[2])
	} else {
		return "неизвестная команда"
	}
}
