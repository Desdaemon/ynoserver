package main

import (
	"errors"
	"strconv"
	"strings"
)

func (h *Hub) handleM(msg []string, sender *Client) (err error) {
	if len(msg) != 3 {
		return err
	}
	//check if the coordinates are valid
	x, errconv := strconv.Atoi(msg[1])
	if errconv != nil || x < 0 {
		return errconv
	}
	y, errconv := strconv.Atoi(msg[2])
	if errconv != nil || y < 0 {
		return errconv
	}
	sender.x = x
	sender.y = y

	if msg[0] == "m" {
		if sender.syncCoords {
			checkHubConditions(h, sender, "coords", "")
		}
		h.broadcast([]byte("m" + paramDelimStr + strconv.Itoa(sender.id) + paramDelimStr + msg[1] + paramDelimStr + msg[2])) //user %id% moved to x y
	} else {
		checkHubConditions(h, sender, "teleport", "")
	}

	return nil
}

func (h *Hub) handleF(msg []string, sender *Client) (err error) {
	if len(msg) != 2 {
		return err
	}
	//check if direction is valid
	facing, errconv := strconv.Atoi(msg[1])
	if errconv != nil || facing < 0 || facing > 3 {
		return errconv
	}
	sender.facing = facing
	h.broadcast([]byte("f" + paramDelimStr + strconv.Itoa(sender.id) + paramDelimStr + msg[1])) //user %id% facing changed to f

	return nil
}

func (h *Hub) handleSpd(msg []string, sender *Client) (err error) {
	if len(msg) != 2 {
		return err
	}
	spd, errconv := strconv.Atoi(msg[1])
	if errconv != nil {
		return errconv
	}
	if spd < 0 || spd > 10 { //something's not right
		return errconv
	}
	sender.spd = spd
	h.broadcast([]byte("spd" + paramDelimStr + strconv.Itoa(sender.id) + paramDelimStr + msg[1]))

	return nil
}

func (h *Hub) handleSpr(msg []string, sender *Client) (err error) {
	if len(msg) != 3 {
		return err
	}
	if !isValidSprite(msg[1]) {
		return err
	}
	/*if config.gameName == "2kki" {
		if !strings.Contains(msg[1], "syujinkou") && !strings.Contains(msg[1], "effect") && !strings.Contains(msg[1], "yukihitsuji_game") && !strings.Contains(msg[1], "zenmaigaharaten_kisekae") && !strings.Contains(msg[1], "主人公") {
			return err
		}
		if strings.Contains(msg[1], "zenmaigaharaten_kisekae") && h.roomName != "176" {
			return err
		}
	}*/
	index, errconv := strconv.Atoi(msg[2])
	if errconv != nil || index < 0 {
		return errconv
	}
	sender.spriteName = msg[1]
	sender.spriteIndex = index
	h.broadcast([]byte("spr" + paramDelimStr + strconv.Itoa(sender.id) + paramDelimStr + msg[1] + paramDelimStr + msg[2]))

	return nil
}

func (h *Hub) handleFl(msg []string, sender *Client) (err error) {
	if len(msg) != 6 {
		return err
	}
	red, errconv := strconv.Atoi(msg[1])
	if errconv != nil || red < 0 || red > 255 {
		return errconv
	}
	green, errconv := strconv.Atoi(msg[2])
	if errconv != nil || green < 0 || green > 255 {
		return errconv
	}
	blue, errconv := strconv.Atoi(msg[3])
	if errconv != nil || blue < 0 || blue > 255 {
		return errconv
	}
	power, errconv := strconv.Atoi(msg[4])
	if errconv != nil || power < 0 {
		return errconv
	}
	frames, errconv := strconv.Atoi(msg[5])
	if errconv != nil || frames < 0 {
		return errconv
	}
	if msg[0] == "rfl" {
		sender.flash[0] = red
		sender.flash[1] = green
		sender.flash[2] = blue
		sender.flash[3] = power
		sender.flash[4] = frames
		sender.repeatingFlash = true
	}
	h.broadcast([]byte(msg[0] + paramDelimStr + strconv.Itoa(sender.id) + paramDelimStr + msg[1] + paramDelimStr + msg[2] + paramDelimStr + msg[3] + paramDelimStr + msg[4] + paramDelimStr + msg[5]))

	return nil
}

func (h *Hub) handleRrfl(msg []string, sender *Client) (err error) {
	sender.repeatingFlash = false
	for i := 0; i < 5; i++ {
		sender.flash[i] = 0
	}
	h.broadcast([]byte("rrfl" + paramDelimStr + strconv.Itoa(sender.id)))

	return nil
}

func (h *Hub) handleT(msg []string, sender *Client) (err error) {
	if len(msg) != 5 {
		return err
	}
	red, errconv := strconv.Atoi(msg[1])
	if errconv != nil || red < 0 || red > 255 {
		return errconv
	}
	green, errconv := strconv.Atoi(msg[2])
	if errconv != nil || green < 0 || green > 255 {
		return errconv
	}
	blue, errconv := strconv.Atoi(msg[3])
	if errconv != nil || blue < 0 || blue > 255 {
		return errconv
	}
	gray, errconv := strconv.Atoi(msg[4])
	if errconv != nil || red < 0 || gray > 255 {
		return errconv
	}
	sender.tone[0] = red
	sender.tone[1] = green
	sender.tone[2] = blue
	sender.tone[3] = gray
	h.broadcast([]byte("t" + paramDelimStr + strconv.Itoa(sender.id) + paramDelimStr + msg[1] + paramDelimStr + msg[2] + paramDelimStr + msg[3] + paramDelimStr + msg[4]))

	return nil
}

func (h *Hub) handleSys(msg []string, sender *Client) (err error) {
	if len(msg) != 2 {
		return err
	}
	if !isValidSystem(msg[1], false) {
		return err
	}
	sender.systemName = msg[1]
	h.broadcast([]byte("sys" + paramDelimStr + strconv.Itoa(sender.id) + paramDelimStr + msg[1]))

	return nil
}

func (h *Hub) handleSe(msg []string, sender *Client) (err error) {
	if len(msg) != 5 || msg[1] == "" {
		return err
	}
	if !isValidSound(msg[1]) {
		return err
	}
	volume, errconv := strconv.Atoi(msg[2])
	if errconv != nil || volume < 0 || volume > 100 {
		return errconv
	}
	tempo, errconv := strconv.Atoi(msg[3])
	if errconv != nil || tempo < 10 || tempo > 400 {
		return errconv
	}
	balance, errconv := strconv.Atoi(msg[4])
	if errconv != nil || balance < 0 || balance > 100 {
		return errconv
	}
	h.broadcast([]byte("se" + paramDelimStr + strconv.Itoa(sender.id) + paramDelimStr + msg[1] + paramDelimStr + msg[2] + paramDelimStr + msg[3] + paramDelimStr + msg[4]))

	return nil
}

func (h *Hub) handleP(msg []string, sender *Client) (err error) {
	isShow := msg[0] == "ap"
	msgLength := 18
	if isShow {
		msgLength = msgLength + 2
	}
	if len(msg) != msgLength {
		return err
	}

	if isShow {
		checkHubConditions(h, sender, "picture", msg[17])
		if !isValidPicName(msg[17]) {
			return err
		}
	}

	picId, errconv := strconv.Atoi(msg[1])
	if errconv != nil || picId < 1 {
		return errconv
	}

	positionX, errconv := strconv.Atoi(msg[2])
	if errconv != nil {
		return errconv
	}
	positionY, errconv := strconv.Atoi(msg[3])
	if errconv != nil {
		return errconv
	}
	mapX, errconv := strconv.Atoi(msg[4])
	if errconv != nil {
		return errconv
	}
	mapY, errconv := strconv.Atoi(msg[5])
	if errconv != nil {
		return errconv
	}
	panX, errconv := strconv.Atoi(msg[6])
	if errconv != nil {
		return errconv
	}
	panY, errconv := strconv.Atoi(msg[7])
	if errconv != nil {
		return errconv
	}

	magnify, errconv := strconv.Atoi(msg[8])
	if errconv != nil || magnify < 0 {
		return errconv
	}
	topTrans, errconv := strconv.Atoi(msg[9])
	if errconv != nil || topTrans < 0 {
		return errconv
	}
	bottomTrans, errconv := strconv.Atoi(msg[10])
	if errconv != nil || bottomTrans < 0 {
		return errconv
	}

	red, errconv := strconv.Atoi(msg[11])
	if errconv != nil || red < 0 || red > 200 {
		return errconv
	}
	green, errconv := strconv.Atoi(msg[12])
	if errconv != nil || green < 0 || green > 200 {
		return errconv
	}
	blue, errconv := strconv.Atoi(msg[13])
	if errconv != nil || blue < 0 || blue > 200 {
		return errconv
	}
	saturation, errconv := strconv.Atoi(msg[14])
	if errconv != nil || saturation < 0 || saturation > 200 {
		return errconv
	}

	effectMode, errconv := strconv.Atoi(msg[15])
	if errconv != nil || effectMode < 0 {
		return errconv
	}
	effectPower, errconv := strconv.Atoi(msg[16])
	if errconv != nil {
		return errconv
	}

	var pic *Picture
	if isShow {
		picName := msg[17]
		if picName == "" {
			return err
		}

		useTransparentColorBin, errconv := strconv.Atoi(msg[18])
		if errconv != nil || useTransparentColorBin < 0 || useTransparentColorBin > 1 {
			return errconv
		}

		fixedToMapBin, errconv := strconv.Atoi(msg[19])
		if errconv != nil || fixedToMapBin < 0 || fixedToMapBin > 1 {
			return errconv
		}

		pic = &Picture{
			name:                picName,
			useTransparentColor: useTransparentColorBin == 1,
			fixedToMap:          fixedToMapBin == 1,
		}

		if _, found := sender.pictures[picId]; found {
			_, rpErr := h.processMsg("rp"+paramDelimStr+msg[1], sender)
			if rpErr != nil {
				return rpErr
			}
		}
	} else {
		if _, found := sender.pictures[picId]; found {
			duration, errconv := strconv.Atoi(msg[17])
			if errconv != nil || duration < 0 {
				return errconv
			}

			pic = sender.pictures[picId]
		} else {
			return nil
		}
	}

	pic.positionX = positionX
	pic.positionY = positionY
	pic.mapX = mapX
	pic.mapY = mapY
	pic.panX = panX
	pic.panY = panY
	pic.magnify = magnify
	pic.topTrans = topTrans
	pic.bottomTrans = bottomTrans
	pic.red = red
	pic.blue = blue
	pic.green = green
	pic.saturation = saturation
	pic.effectMode = effectMode
	pic.effectPower = effectPower

	message := msg[0] + paramDelimStr + strconv.Itoa(sender.id) + paramDelimStr + msg[1] + paramDelimStr + msg[2] + paramDelimStr + msg[3] + paramDelimStr + msg[4] + paramDelimStr + msg[5] + paramDelimStr + msg[6] + paramDelimStr + msg[7] + paramDelimStr + msg[8] + paramDelimStr + msg[9] + paramDelimStr + msg[10] + paramDelimStr + msg[11] + paramDelimStr + msg[12] + paramDelimStr + msg[13] + paramDelimStr + msg[14] + paramDelimStr + msg[15] + paramDelimStr + msg[16] + paramDelimStr + msg[17]
	if isShow {
		message += paramDelimStr + msg[18] + paramDelimStr + msg[19]
	}
	h.broadcast([]byte(message))
	sender.pictures[picId] = pic

	return nil
}

func (h *Hub) handleRp(msg []string, sender *Client) (err error) {
	if len(msg) != 2 {
		return err
	}
	picId, errconv := strconv.Atoi(msg[1])
	if errconv != nil || picId < 1 {
		return errconv
	}
	h.broadcast([]byte("rp" + paramDelimStr + strconv.Itoa(sender.id) + paramDelimStr + msg[1]))
	delete(sender.pictures, picId)

	return nil
}

func (h *Hub) handleSay(msg []string, sender *Client) (err error) {
	msgLength := 2
	if msg[0] == "gsay" {
		msgLength++
	}
	if len(msg) != msgLength {
		return err
	}
	msgContents := strings.TrimSpace(msg[1])
	if sender.name == "" || sender.systemName == "" || msgContents == "" || len(msgContents) > 150 {
		return err
	}
	switch msg[0] {
	case "gsay":
		enableLocBin, errconv := strconv.Atoi(msg[2])
		if errconv != nil || enableLocBin < 0 || enableLocBin > 1 {
			return errconv
		}

		mapId := "0000"
		prevMapId := "0000"
		var prevLocations string
		x := -1
		y := -1

		if enableLocBin == 1 {
			mapId = sender.mapId
			prevMapId = sender.prevMapId
			prevLocations = sender.prevLocations
			x = sender.x
			y = sender.y
		}

		var accountBin int
		if sender.account {
			accountBin = 1
		}

		globalBroadcast([]byte("gsay" + paramDelimStr + sender.uuid + paramDelimStr + sender.name + paramDelimStr + sender.systemName + paramDelimStr + strconv.Itoa(sender.rank) + paramDelimStr + strconv.Itoa(accountBin) + paramDelimStr + sender.badge + paramDelimStr + mapId + paramDelimStr + prevMapId + paramDelimStr + prevLocations + paramDelimStr + strconv.Itoa(x) + paramDelimStr + strconv.Itoa(y) + paramDelimStr + msgContents))
	case "psay":
		partyId, err := readPlayerPartyId(sender.uuid)
		if err != nil {
			return err
		}
		if partyId == 0 {
			return errors.New("player not in a party")
		}
		partyMemberUuids, err := readPartyMemberUuids(partyId)
		if err != nil {
			return err
		}
		for _, uuid := range partyMemberUuids {
			if client, ok := allClients[uuid]; ok {
				client.send <- []byte("psay" + paramDelimStr + sender.uuid + paramDelimStr + msgContents)
			}
		}
	default:
		h.broadcast([]byte("say" + paramDelimStr + strconv.Itoa(sender.id) + paramDelimStr + msgContents))
	}

	return nil
}

func (h *Hub) handleName(msg []string, sender *Client) (err error) {
	if sender.name != "" || len(msg) != 2 || !isOkString(msg[1]) || len(msg[1]) > 12 {
		return err
	}
	sender.name = msg[1]
	h.broadcast([]byte("name" + paramDelimStr + strconv.Itoa(sender.id) + paramDelimStr + sender.name))

	return nil
}

func (h *Hub) handleSs(msg []string, sender *Client) (err error) {
	if len(msg) != 3 {
		return err
	}
	switchId, errconv := strconv.Atoi(msg[1])
	if errconv != nil {
		return errconv
	}
	valueBin, errconv := strconv.Atoi(msg[2])
	if errconv != nil || valueBin < 0 || valueBin > 1 {
		return errconv
	}
	var value bool
	if valueBin == 1 {
		value = true
	}
	sender.switchCache[switchId] = value
	if switchId == 1430 && config.gameName == "2kki" {
		if value {
			sender.send <- []byte("sv" + paramDelimStr + "88" + paramDelimStr + "0")
		}
	} else {
		if len(sender.hub.minigameConfigs) > 0 {
			for m, minigame := range sender.hub.minigameConfigs {
				if minigame.SwitchId == switchId && minigame.SwitchValue == value && sender.minigameScores[m] < sender.varCache[minigame.VarId] {
					tryWritePlayerMinigameScore(sender.uuid, minigame.MinigameId, sender.varCache[minigame.VarId])
				}
			}
		}

		for _, c := range h.conditions {
			validVars := !c.VarTrigger
			if c.VarTrigger {
				if c.VarId > 0 {
					if value, ok := sender.varCache[c.VarId]; ok {
						if validVar, _ := c.checkVar(c.VarId, value); validVar {
							validVars = true
						}
					}
				} else if len(c.VarIds) > 0 {
					validVars = true
					for _, vId := range c.VarIds {
						if value, ok := sender.varCache[vId]; ok {
							if validVar, _ := c.checkVar(vId, value); !validVar {
								validVars = false
								break
							}
						} else {
							validVars = false
							break
						}
					}
				} else {
					validVars = true
				}
			}

			if validVars {
				if switchId == c.SwitchId {
					if valid, _ := c.checkSwitch(switchId, value); valid {
						if c.VarTrigger || (c.VarId == 0 && len(c.VarIds) == 0) {
							if !c.TimeTrial {
								if checkConditionCoords(c, sender) {
									success, err := tryWritePlayerTag(sender.uuid, c.ConditionId)
									if err != nil {
										return err
									}
									if success {
										sender.send <- []byte("b")
									}
								}
							} else if config.gameName == "2kki" {
								sender.send <- []byte("ss" + paramDelimStr + "1430" + paramDelimStr + "0")
							}
						} else {
							varId := c.VarId
							if len(c.VarIds) > 0 {
								varId = c.VarIds[0]
							}
							sender.send <- []byte("sv" + paramDelimStr + strconv.Itoa(varId) + paramDelimStr + "0")
						}
					}
				} else if len(c.SwitchIds) > 0 {
					if valid, s := c.checkSwitch(switchId, value); valid {
						if s == len(c.SwitchIds)-1 {
							if c.VarTrigger || (c.VarId == 0 && len(c.VarIds) == 0) {
								if !c.TimeTrial {
									if checkConditionCoords(c, sender) {
										success, err := tryWritePlayerTag(sender.uuid, c.ConditionId)
										if err != nil {
											return err
										}
										if success {
											sender.send <- []byte("b")
										}
									}
								} else if config.gameName == "2kki" {
									sender.send <- []byte("ss" + paramDelimStr + "1430" + paramDelimStr + "0")
								}
							} else {
								varId := c.VarId
								if len(c.VarIds) > 0 {
									varId = c.VarIds[0]
								}
								sender.send <- []byte("sv" + paramDelimStr + strconv.Itoa(varId) + paramDelimStr + "0")
							}
						} else {
							sender.send <- []byte("ss" + paramDelimStr + strconv.Itoa(c.SwitchIds[s+1]) + paramDelimStr + "0")
						}
					}
				}
			}
		}
	}

	return nil
}

func (h *Hub) handleSv(msg []string, sender *Client) (err error) {
	if len(msg) != 3 {
		return err
	}
	varId, errconv := strconv.Atoi(msg[1])
	if errconv != nil {
		return errconv
	}
	value, errconv := strconv.Atoi(msg[2])
	if errconv != nil {
		return errconv
	}
	sender.varCache[varId] = value
	if varId == 88 && config.gameName == "2kki" {
		for _, c := range h.conditions {
			if c.TimeTrial && value < 3600 {
				if checkConditionCoords(c, sender) {
					mapId, _ := strconv.Atoi(h.roomName)
					success, err := tryWritePlayerTimeTrial(sender.uuid, mapId, value)
					if err != nil {
						return err
					}
					if success {
						sender.send <- []byte("b")
					}
				}
			}
		}
	} else {
		if len(sender.hub.minigameConfigs) > 0 {
			for m, minigame := range sender.hub.minigameConfigs {
				if minigame.VarId == varId && sender.minigameScores[m] < value {
					if minigame.SwitchId > 0 {
						sender.send <- []byte("ss" + paramDelimStr + strconv.Itoa(minigame.SwitchId) + paramDelimStr + "0")
					} else {
						tryWritePlayerMinigameScore(sender.uuid, minigame.MinigameId, value)
					}
				}
			}
		}

		for _, c := range h.conditions {
			validSwitches := c.VarTrigger
			if !c.VarTrigger {
				if c.SwitchId > 0 {
					if value, ok := sender.switchCache[c.SwitchId]; ok {
						if validSwitch, _ := c.checkSwitch(c.SwitchId, value); validSwitch {
							validSwitches = true
						}
					}
				} else if len(c.SwitchIds) > 0 {
					validSwitches = true
					for _, sId := range c.SwitchIds {
						if value, ok := sender.switchCache[sId]; ok {
							if validSwitch, _ := c.checkSwitch(sId, value); !validSwitch {
								validSwitches = false
								break
							}
						} else {
							validSwitches = false
							break
						}
					}
				} else {
					validSwitches = true
				}
			}

			if validSwitches {
				if varId == c.VarId {
					if valid, _ := c.checkVar(varId, value); valid {
						if !c.VarTrigger || (c.SwitchId == 0 && len(c.SwitchIds) == 0) {
							if !c.TimeTrial {
								if checkConditionCoords(c, sender) {
									success, err := tryWritePlayerTag(sender.uuid, c.ConditionId)
									if err != nil {
										return err
									}
									if success {
										sender.send <- []byte("b")
									}
								}
							} else if config.gameName == "2kki" {
								sender.send <- []byte("ss" + paramDelimStr + "1430" + paramDelimStr + "0")
							}
						} else {
							switchId := c.SwitchId
							if len(c.SwitchIds) > 0 {
								switchId = c.SwitchIds[0]
							}
							sender.send <- []byte("ss" + paramDelimStr + strconv.Itoa(switchId) + paramDelimStr + "0")
						}
					}
				} else if len(c.VarIds) > 0 {
					if valid, v := c.checkVar(varId, value); valid {
						if v == len(c.VarIds)-1 {
							if !c.VarTrigger || (c.SwitchId == 0 && len(c.SwitchIds) == 0) {
								if !c.TimeTrial {
									if checkConditionCoords(c, sender) {
										success, err := tryWritePlayerTag(sender.uuid, c.ConditionId)
										if err != nil {
											return err
										}
										if success {
											sender.send <- []byte("b")
										}
									}
								} else if config.gameName == "2kki" {
									sender.send <- []byte("ss" + paramDelimStr + "1430" + paramDelimStr + "0")
								}
							} else {
								switchId := c.SwitchId
								if len(c.SwitchIds) > 0 {
									switchId = c.SwitchIds[0]
								}
								sender.send <- []byte("ss" + paramDelimStr + strconv.Itoa(switchId) + paramDelimStr + "0")
							}
						} else {
							sender.send <- []byte("sv" + paramDelimStr + strconv.Itoa(c.VarIds[v+1]) + paramDelimStr + "0")
						}
					}
				}
			}
		}
	}

	return nil
}

func (h *Hub) handleSev(msg []string, sender *Client) (err error) {
	if len(msg) != 3 {
		return err
	}
	actionBin, errconv := strconv.Atoi(msg[2])
	if errconv != nil || actionBin < 0 || actionBin > 1 {
		return errconv
	}
	triggerType := "event"
	if actionBin == 1 {
		triggerType = "eventAction"
	}
	checkHubConditions(h, sender, triggerType, msg[1])

	return nil
}