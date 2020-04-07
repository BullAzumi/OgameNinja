/* Creaded by Bull
 * https://github.com/BullAzumi/Bull
 * donations by paypal 
 * hololo40@gmail.com
 * Thanks! Helps me a lot!
 */
 
/*
 * VERSION 1.34
 */

/* DESCRIPTION
 * Automatically sends expeditions when a slot or expedition slot is free.
 * You can decide whether to choose the ship automatically or by yourself.
 * If you don't have enough ships on the planet / moon, he will send the number of waves divided or set it to 0
 * notifies you of Telegram
 * It's possible you're targeting debris fields
 */

/*---------------------------------------------------------------------------------------------------------------------------*/

//######################################## SETTINGS START ########################################


home = "M:1:234:5"                                      //set your expo home
radius = 0                                              //set a radius, this will be flown around your system. If radius = 0 then we only fly in your system
rankOnePoints = 1234567890                              //how many points does the first place have
maxExpoSlotsUse = 8                                     //how many slots should we use for fly Expos?
maxDebrisSlots = 2										//how many slots should we use for mining Debris?
loop = true                                             //should we fly every round again?
expoTime = 1                                            //set duration
useWave = true                                          //should we let everyone fly on a coordinate and only then switch? true = yes / false = no
useTelegram = true                                      //should we notify you on telegram? true = yes / false = no
smallCargo = false                                      //would you fly with small cargos? true = yes / false = no
usePathfinder = true                                    //use pathfinder in automatically build
useReaper = true                                        //true = use reaoer in automatically build / false = use destroyer
HeavyFighter = true                                     //should add HeavyFighters? yes = true / no = false we would use 3.000 HeavyFighters
selfShips = false                                       //assemble the ships yourself or have them calculated automatically true = self / false = automatic
TeleID = TELEGRAM_CHAT_ID                               //you can exchange this for an ID for a possible second account (CloudHost only)
mineDebris = true                                       //should we mine debris fields? true = yes / false = no
endIt = {true:"21:30:00"}								//should we stop sending expedition at any time? (NinjaTime not OGameTime) mining will continue to work true = yes / false = no
expoDelay = [3,7]                                       //sends expeditions at intervals of x seconds
muchPaths = 1                                           //all debris larger than X pathfinder are removed

//Change if you are selfShipper!
ship = {LIGHTFIGHTER : 0 ,
        HEAVYFIGHTER : 0 ,
        CRUISER : 0 ,
        BATTLESHIP : 0 ,
        BATTLECRUISER : 0 ,
        BOMBER : 0 ,
        DESTROYER : 0 ,
        DEATHSTAR : 0 ,
        SMALLCARGO : 0 ,
        LARGECARGO : 0 ,
        COLONYSHIP : 0 ,
        RECYCLER : 0 ,
        ESPIONAGEPROBE : 0 ,
        REAPER : 0 ,
        PATHFINDER : 0}
//########################## !!!CHANGE IT ONLY IF YOU KNOW WHAT YOU ARE DOING!!! ###########################

hartDebris = false                                       //set this true do get hart Debris farm
LFandSC = false                                         //true = use light fighters and small cargos / false = use heavy fighters and large cargos

//######################################## SETTINGS END ########################################


speed, myGalaxy, SC, LC, expoSlots, expoInUse, freeSlots, reserved, minSecs, debrisSysDown, debrisSysUp, downSys, upSys, usedSlots, usedDebris = 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0
ships, tech, fleets = {},{},{}
send, StopSending = false, false
fleetIDList = []

//####################
//# debris functions #
//####################

func doDebris() {
	LogTelegram("D", "Call function sendMineDebris()")
	for {
        for info, amount in scanGala() {
            if onTheWayDebris(info.System(), amount) {
                LogTelegram("I", Dotify(info.ExpeditionDebris.Metal) + " metal " + Dotify(info.ExpeditionDebris.Crystal) + " crystal were found!")
                if usedDebris < maxDebrisSlots && freeSlots > 0 {
                    f, err = sendDebris(amount, info.System())
                    if err != nil {
                        LogTelegram("E", err)
                    }else {
                        LogTelegram("D", "Pathfinder were sendet")
                        fleetIDList += f
                        fleetsGet()
                        Sleep(Random(3,7)*1000)
                    }
                }else {
                    fleetsGet()
                }
            }
        }
        customSleep(DebrisInterval(), "doDebris")
    }
}

func DebrisInterval() {
	time = 999999999
	fleets, slots = GetFleets()
	for fleet in fleets {
		if fleet.Mission == EXPEDITION {
			if fleet.ArriveIn < 0 {
				time = Min(fleet.BackIn, time)
			}else {
				time = Min(fleet.ArriveIn, time)
			}
		}
	}
	if time == 999999999 {
		time = 0
	}
	return time
}

func scanGala() {
    LogTelegram("D", "Call function scanGala()")
	DebrisInfos = {}
	for i = debrisSysDown; i <= debrisSysUp; i++ {
		systemInfo, err = GalaxyInfos(myGalaxy, i)
		if err != nil {
			LogTelegram("E", err)
		}else {
			if systemInfo.ExpeditionDebris.PathfindersNeeded > muchPaths {
				LogTelegram("D", "Debris detected")
				LogTelegram("I", "Need " + Dotify(systemInfo.ExpeditionDebris.PathfindersNeeded) + " Pathfinder")
				DebrisInfos[systemInfo] = DebrisInfos[systemInfo]+systemInfo.ExpeditionDebris.PathfindersNeeded
			}
		}
	}
	return DebrisInfos
}

func onTheWayDebris(debrisSys, PathAmount) {
	fleets, slots = GetFleets()
	for fleet in fleets {
		if fleet.Mission == RECYCLEDEBRISFIELD && fleet.Destination == NewCoordinate(myGalaxy, debrisSys, 16, DEBRIS_TYPE) && fleet.Ships.Pathfinder == PathAmount{
			return false
		}
	}
    return true
}

func sendDebris(shipCounter, debrisSys) {
    LogTelegram("D", "Call function sendDebris()")
    sTo = NewCoordinate(myGalaxy, debrisSys, 16, DEBRIS_TYPE)
    fleet = NewFleet()
    fleet.SetOrigin(home)
    fleet.SetDestination(sTo)
    fleet.SetMission(RECYCLEDEBRISFIELD)
    fleet.AddShips(PATHFINDER, shipCounter)
    return fleet.SendNow()
}

//########################
//# expeditions function #
//########################

func doExpo() {
	send = false
	for usedSlots < maxExpoSlotsUse && expoInUse < expoSlots && freeSlots > 0 && !StopSending{
		if downSys <= upSys || downSys > upSys {
			newFleet, err = sendExpo()
			if err != nil {
				LogTelegram("E", "Fleet reports errors: " + err)
				fleetsGet()
				setMinSecs()
				send = false
				break
			}else{
				LogTelegram("I", "Fleet starts on the expedition! " + newFleet)
                fleetIDList += newFleet
				fleetsGet()
                send = true
				if !useWave {
					downSys++
				}
				Sleep(Random(expoDelay[0],expoDelay[1])*1000)
			}
		}
	}
}

func sendExpo() {
	LogTelegram("D", "Call function sendExpo()")  
	sTo = NewCoordinate(myGalaxy, downSys, 16, PLANET_TYPE)        
	fleet = NewFleet()
	fleet.SetOrigin(home)
	fleet.SetDestination(sTo)
	fleet.SetSpeed(HUNDRED_PERCENT)
	fleet.SetMission(EXPEDITION)
	fleet.SetShips(*ships)
	fleet.SetDuration(expoTime)
	return fleet.SendNow()
}

func doAll() {
	for {
		if !StopSending {
            fleetsGet()                                                             
            LogTelegram("D", "Expeditions available: " + (expoSlots - expoInUse))
            LogTelegram("D", "Slots available: " + freeSlots)
            if maxExpoSlotsUse - usedSlots > 0 {
                shipsSet()
            }
            doExpo()
            if useWave && send {
                downSys++
            }
            checkLoop()   
            setMinSecs()
            customSleep(minSecs, "doExpo")
		}
	}
}

func checkLoop() {
    if downSys == upSys + 1{
        if loop {
            LogTelegram("D", "Loop is active! We'll start again from the beginning!")
            downSys = GetCachedCelestial(home).Coordinate.System - radius
        }else {
            LogTelegram("I", "All done! We're out of here!")
            StopScript(__FILE__)
        }
    }
}

func enoughtShips(shipNames, shipAmounts) {
    eShips, sErr = GetCachedCelestial(home).GetShips()
    if sErr != nil {
        LogTelegram("E", sErr)
        StopScript(__FILE__)
    }else {
        if eShips.ByID(shipNames) < shipAmounts {
            LogTelegram("D", "Not enought " + ID2Str(shipNames))
            shipsMax = Floor(eShips.ByID(shipNames) / (maxExpoSlotsUse - usedSlots))
            LogTelegram("D", "We set it do " + shipsMax + " " + ID2Str(shipNames))
            return shipsMax
        }else {
            shipsMax = shipAmounts / (maxExpoSlotsUse - usedSlots)
            return shipsMax
        }
    }
}

func shipsSet() {
    LogTelegram("D", "Call function shipsSet()")
    if selfShips {
        manualSet()
    }
    if !selfShips && !hartDebris{
        autoSet()
    }
    if hartDebris {
        setHartDebris()
    }
}

func setHartDebris() {
    LogTelegram("D", "Hart Debris settings are adopted!")
    fleetsGet()
    ships = NewShipsInfos()
    if LFandSC {
        SC = enoughtShips(SMALLCARGO, SC * (maxExpoSlotsUse - usedSlots))
        LF = Ceil(enoughtShips(LIGHTFIGHTER, ((SC / 2) * (maxExpoSlotsUse - usedSlots))))
        LogTelegram("I", "We add " + Dotify(SC) + " SmallCargos.") 
        LogTelegram("I", "We add " + Dotify(LF) + " LightFighters.")
        ships.Set(SMALLCARGO, SC)
        ships.Set(LIGHTFIGHTER, LF)
        ships.Set(ESPIONAGEPROBE, enoughtShips(ESPIONAGEPROBE, 1 * (maxExpoSlotsUse - usedSlots)))
        if useReaper {
            ships.Set(REAPER, enoughtShips(REAPER, 1 * (maxExpoSlotsUse - usedSlots)))
        }else {
            ships.Set(DESTROYER, enoughtShips(DESTROYER, 1 * (maxExpoSlotsUse - usedSlots)))
        }
        if usePathfinder {
            ships.Set(PATHFINDER, enoughtShips(PATHFINDER, 1 * (maxExpoSlotsUse - usedSlots)))
        }
    }else {
        LC = enoughtShips(LARGECARGO, LC * (maxExpoSlotsUse - usedSlots))
        HF = Ceil(enoughtShips(HEAVYFIGHTER, ((LC / 1.5) * (maxExpoSlotsUse - usedSlots))))
        LogTelegram("I", "We add " + Dotify(LC) + " LargeCargos.")
        LogTelegram("I", "We add " + Dotify(HF) + " HeavyFighters.")
        ships.Set(HEAVYFIGHTER, HF)
        ships.Set(LARGECARGO, LC)
        ships.Set(ESPIONAGEPROBE, enoughtShips(ESPIONAGEPROBE, 1 * (maxExpoSlotsUse - usedSlots)))
        if useReaper {
            ships.Set(REAPER, enoughtShips(REAPER, 1 * (maxExpoSlotsUse - usedSlots)))
        }else {
            ships.Set(DESTROYER, enoughtShips(DESTROYER, 1 * (maxExpoSlotsUse - usedSlots)))
        }
        if usePathfinder {
            ships.Set(PATHFINDER, enoughtShips(PATHFINDER, 1 * (maxExpoSlotsUse - usedSlots)))
        }
    }
}

func manualSet() {
    LogTelegram("D", "Manual settings are adopted!")
    fleetsGet()
    ships = NewShipsInfos()
    for shipName, shipAmount in ship {
        if shipAmount != 0 { 
            amount = enoughtShips(shipName, shipAmount * (maxExpoSlotsUse - usedSlots))
            if amount > 0 {
                ships.Set(shipName, amount)
            }else {
                LogTelegram("E", "Not enought " + ID2Str(shipName))
            }
        }
    }
    LogTelegram("I", "The following ships have been imported: " + ships)
}

func autoSet() {
    LogTelegram("D", "Automatically settings are adopted!")
    fleetsGet()
    expoValue = 0
    ships = NewShipsInfos()
    switch rankOnePoints {
        case rankOnePoints >= 100000000:
            expoValue = 25000
        case rankOnePoints >= 75000000:
            expoValue = 21000
        case rankOnePoints >= 50000000:
            expoValue = 18000
        case rankOnePoints >= 25000000:
            expoValue = 15000
        case rankOnePoints >= 5000000:
            expoValue = 12000
        case rankOnePoints >= 1000000:
            expoValue = 9000
        case rankOnePoints >= 100000:
            expoValue = 6000
        case rankOnePoints >= 0:
            expoValue = 2400
    }
    if IsDiscoverer() {
        SC = Ceil((expoValue * speed * 1.5 * 2 * 200) / (5000 * (tech.HyperspaceTechnology * 0.05) + 1))
        LC = Ceil((expoValue * speed * 1.5 * 2 * 200) / (25000 * (tech.HyperspaceTechnology * 0.05) + 1))
    }else {
        SC = Ceil(expoValue / 20)
        LC = Ceil(expoValue / 60)
    }
    ships = NewShipsInfos()
    if smallCargo {
        tmp = enoughtShips(SMALLCARGO, SC * (maxExpoSlotsUse - usedSlots))
        ships.Set(SMALLCARGO, tmp)
        LogTelegram("I", "Add " + tmp + " small cargos")
    }else {
        tmp = enoughtShips(LARGECARGO, LC * (maxExpoSlotsUse - usedSlots))
        ships.Set(LARGECARGO, tmp)
        LogTelegram("I", "Add " + tmp + " large cargos")
    }
    ships.Set(ESPIONAGEPROBE, enoughtShips(ESPIONAGEPROBE, 1 * (maxExpoSlotsUse - usedSlots)))
    if useReaper {
        ships.Set(REAPER, enoughtShips(REAPER, 1 * (maxExpoSlotsUse - usedSlots)))
    }else {
        ships.Set(DESTROYER, enoughtShips(DESTROYER, 1 * (maxExpoSlotsUse - usedSlots)))
    }
    if usePathfinder {
        ships.Set(PATHFINDER, enoughtShips(PATHFINDER, 1 * (maxExpoSlotsUse - usedSlots)))
    }
    if HeavyFighter {
        ships.Set(HEAVYFIGHTER, enoughtShips(HEAVYFIGHTER, 3000 * (maxExpoSlotsUse - usedSlots)))
    }
}

//###################
//# other Functions #
//###################

func makeBreak(infos){
    for {
        for info in infos{
            customSleep(info, "takeABreak")
            if !breakBool {
                breakBool = true
            }else {
                breakBool = false
            }
        }
    }
}

func setMinSecs() {
    LogTelegram("D", "Call function setMinSecs()")
    minSecs = 999999999
    for fleet in fleets {
        if fleet.Mission == EXPEDITION {
			if fleet.ArriveIn < 0 {
				minSecs = Min(fleet.BackIn, minSecs)
			}else {
				minSecs = Min(fleet.ArriveIn, minSecs)
			}
        }
	}
	if minSecs == 999999999 {
		minSecs = 0
	}
}

func customSleep(sleepTime, from) {
    if sleepTime <= 0 {
        sleepTime = Random(5, 10)
    }
    LogTelegram("I", from + " wants me to wait " + ShortDur(sleepTime + 10))
    Sleep((sleepTime + 10) * 1000)
}

func fleetsGet() {
    fleets, slots = GetFleets()
    expoSlots = slots.ExpTotal
    expoInUse = slots.ExpInUse
	freeSlots = slots.Total - slots.InUse - reserved
    usedDebris, usedSlots = 0,0
    if len(fleetIDList) > 0 {
        for fleet in fleets {
            for fleetID in fleetIDList {
                if fleet.ID == fleetID.ID {
                    if fleet.Mission == EXPEDITION {
                        usedSlots++
                    }
                    if fleet.Mission == RECYCLEDEBRISFIELD && fleet.Destination.Position == 16 {
                        usedDebris++
                    }
                }
            }
        }
    }
    Put("ExposDone", howMuchExpos())
}

func errorHandler(){
    LogTelegram("D", "Call function errorHandler()")
    myPlanet, err = GetCelestial(home)
    sys = myPlanet.Coordinate.System
    if err != nil {
        LogTelegram("E", home + " is not one of your planet/moon")
        StopScript(__FILE__)
	}
	debrisSysDown = GetCachedCelestial(home).Coordinate.System - radius
	debrisSysUp = GetCachedCelestial(home).Coordinate.System + radius
    downSys = GetCachedCelestial(home).Coordinate.System - radius
    upSys = GetCachedCelestial(home).Coordinate.System + radius
    if downSys <= 0 {
        LogTelegram("D", "Radius exceeds system limit! Set it to 1 ")
		downSys = 1
		debrisSysDown = 1
    }
    if upSys >= 500 {
        LogTelegram("D", "Radius exceeds system limit! Set it to 499")
		upSys = 499
		debrisSysUp	= 499
    }
    if rankOnePoints == 0 || rankOnePoints == nil{
        LogTelegram("E", "We do not believe that the first place has " + rankOnePoints + " points. We at least set it at our level")
        rankOnePoints = GetCachedPlayer().Points
        LogTelegram("E", Dotify(rankOnePoints))
    }
    if maxExpoSlotsUse == 0 || maxExpoSlotsUse == nil{
        LogTelegram("E", "We are not allowed to use slots? Please set maxExpoSlotsUse at least to 1")
        StopScript(__FILE__)
    }
    if expoTime <= 0 || expoTime > 18 || expoTime == nil{
        LogTelegram("E", "we can not stay in expedition for " + expoTime + " hour")
        if expoTime <= 0 || expoTime == nil {
            LogTelegram("E", "we set it at least do 1 hour")
            expoTime = 1
        }else {
            LogTelegram("E", "we set it at least do 18 hour")
            expoTime = 18
        }
    }
    if useTelegram {
        if TeleID == nil || TeleID == 0 {
            LogTelegram("W", "your Telegram ID is not set!")
        }
    }
    if selfShips && hartDebris {
        LogTelegram("E", "it is not possible to have selfShips and hardDebris set true at the same time")
        StopScript(__FILE__)
    }
}

func sendTelegramIfActive(msg) {
	if useTelegram {
		SendTelegram(TeleID, msg)
	}
}

func LogTelegram(cat, message) {
    switch cat {
        case cat == "W" || cat == "w":
            LogWarn(message)
            sendTelegramIfActive("Warn [" + message + "]")
        case cat == "I" || cat == "i":
            LogInfo(message)
            sendTelegramIfActive("Info [" + message + "]")
        case cat == "E" || cat == "e":
            LogError(message)
            sendTelegramIfActive("Error [" + message + "]")
        case cat == "D" || cat == "d":
            LogDebug(message)
        default:
            LogError("no valid entry! Message is: " + message)
            sendTelegramIfActive("Error [No valid entry! Message is: " + message + "]")
    }
}

func boot() {
    LogTelegram("D", "Program starts. Please be patient!")
    errorHandler()
    Sleep(750)
    speed = GetServer().Settings.EconomySpeed
    LogTelegram("D", "Server Economy Speed are " + speed)
    Sleep(250)
    myGalaxy = GetCachedCelestial(home).Coordinate.Galaxy
    LogTelegram("D", "Galaxy was set to " + myGalaxy)
    Sleep(250)
    reserved = GetFleetSlotsReserved()
    LogTelegram("D", "reserved player slots were detected " + reserved)
    Sleep(250)
    tech = GetResearch()
    LogTelegram("D", "Research was received")
	Sleep(250)	
	fleetsGet()
	if maxExpoSlotsUse > expoSlots {
		LogTelegram("E", "your astrophysics is not sufficient for so many expeditions " + maxExpoSlotsUse)
		LogTelegram("W", "we set it to " + expoSlots)
		maxExpoSlotsUse = expoSlots
	}
    startShips, err = GetCachedCelestial(home).GetShips()
    LogTelegram("D", "read out ships, for possible calculations")
    Sleep(250)   
    if err != nil {
        LogTelegram("E", err)
        StopScript(__FILE__)
    }else {
        SC = Floor(startShips.SmallCargo / maxExpoSlotsUse)
        LC = Floor(startShips.LargeCargo / maxExpoSlotsUse)
    }
}

func howMuchExpos(){
    muchExpos = 0
    for fleet in fleetIDList {
        if fleet.Mission == EXPEDITION {
            muchExpos++
        }
    }
    return muchExpos
}

func doWork(){
    LogTelegram("D", "Call function doWork()")
    boot()
    info = []
	if mineDebris {
		ExecIn(0, doDebris)
	}
	ExecIn(0, doAll)
	for isBool, time in endIt {
		if isBool {
			ExecAt(time, func() {
                LogTelegram("I", "Now we are stop to sending Expeditions.")
                LogTelegram("I", "today we have completed " + howMuchExpos + " expeditions") 
				StopSending = true})
		}
	}
}
doWork()
<- OnQuitCh
