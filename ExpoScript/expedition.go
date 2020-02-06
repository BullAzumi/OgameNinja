/* Creaded by Bull
 * https://github.com/BullAzumi/Bull
 */
 
/*
 * VERSION 1.24
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
sendFrom = 233                                          //system lower limit
sendTo = 235                                            //system upper limit
rankOnePoints = 1234567890                              //how many points does the first place have
smallCargo = false                                      //would you fly with small cargos? true = yes / false = no
usePathfinder = true                                    //Use pathfinder in automatically build
HeavyFighter = true                                     //should add HeavyFighters? yes = true / no = false we would use 3.000 HeavyFighters
maxSlotsUse = 1                                         //how many slots should we use?
loop = false                                            //should we fly every round again?
expoTime = 1                                            //set duration
selfShips = false                                       //assemble the ships yourself or have them calculated automatically true = self / false = automatic
useWave = true                                          //should we let everyone fly on a coordinate and only then switch? true = yes / false = no
ignoreFleets = false                                    //ignores fleets are on move true = yes / false = no
useTelegram = true                                      //should we notify you on telegram? true = yes / false = no

//!!!ONLY WORKS IF useWave IS TRUE!!!
mineDebris = true                                       //should we mine debris fields? true = yes / false = no


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
sendAll = false                                          //true = all fleets (set up below) you have are divided by the maxSlots value
                                                         //false = eine kalkulierte flotte die du unten eingestellt hast werden gesendet

LFandSC = true                                           //true = use light fighters and small cargos / false = use heavy fighters and large cargos

//######################################## SETTINGS END ########################################


speed = GetServer().Settings.EconomySpeed
myGalaxy = 0
SC = 0
LC = 0
ships = {}
expoSlots = 0
expoInUse = 0
freeSlots = 0
reserved = GetFleetSlotsReserved()
fleets = {}
send = true
minSecs = 0
debrisSys = 0
downSys = 0
upSys = 0
usedSlots = 0
tech = GetResearch()

//debris functions
func sendMineDebris() {
    LogTelegram("D", "Call function sendMineDebris()")
    if mineDebris && send && useWave {
        debris = scanGala()
        fleetsGet()
        if debris.PathfindersNeeded > 0 {
            LogTelegram("I", debris.Metal + " metal" + debris.Crystal + " crystal were found!")
            for {
                if usedSlots < maxSlotsUse && freeSlots > 0 {
                    f, err = sendDebris(debris.PathfindersNeeded)
                    if err != nil {
                        LogTelegram("E", err)
                        break
                    }else {
                        LogTelegram("I", "Debris in work!")
                        break
                    }
                }else {
                    fleetsGet()
                }
            }
        }
    }
}

func scanGala() {
    LogTelegram("D", "Call function scanGala()")
    systemInfo, err = GalaxyInfos(myGalaxy, debrisSys)
    
    if err != nil {
        LogTelegram("E", err)
    }
    if systemInfo.ExpeditionDebris.PathfindersNeeded > 0 {
        LogTelegram("D", "Debris detected")
        LogTelegram("I", "Need " + systemInfo.ExpeditionDebris.PathfindersNeeded + " Pathfinder")
        return systemInfo.ExpeditionDebris
    }else {
        LogTelegram("D", "No Debris detected")
        return systemInfo.ExpeditionDebris 
    }
}

func sendDebris(shipCounter) {
    LogTelegram("D", "Call function sendDebris()")
    
    sTo = NewCoordinate(myGalaxy, debrisSys, 16, DEBRIS_TYPE)
    fleet = NewFleet()
    fleet.SetOrigin(home)
    fleet.SetDestination(sTo)
    fleet.SetMission(RECYCLEDEBRISFIELD)
    fleet.AddShips(PATHFINDER, shipCounter)
    return fleet.SendNow()
}

//expeditions function
func doExpo() {
    send = false
    for usedSlots < maxSlotsUse{
        for {
            if expoInUse < expoSlots && freeSlots > 0 {
                break;
            }
            fleetsGet()
        }
        if downSys <= upSys || downSys > upSys {
            newFleet, err = sendExpo()
            if err != nil {
                LogTelegram("E", "Fleet reports errors: " + err)
                send = false
                continue
            }else{
                LogTelegram("I", "Fleet starts on the expedition! " + newFleet)
                if minSecs < newFleet.BackIn {
                    minSecs = newFleet.BackIn
                }
                fleetsGet()
                send = true
                if !useWave {
                    debrisSys = downSys
                    downSys = setSys()
                }
                Sleep(Random(3,7)*1000)
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
    LogTelegram("I", "Coords set to " + sTo )
    fleet.SetSpeed(HUNDRED_PERCENT)
    fleet.SetMission(EXPEDITION)
    fleet.SetShips(*ships)
    fleet.SetDuration(expoTime)
    return fleet.SendNow()
}

func setMinSecs() {
    LogTelegram("D", "Call function setMinSecs()")
    for fleet in fleets {
        if fleet.Mission == EXPEDITION {
            if minSecs < fleet.BackIn {
                minSecs = fleet.BackIn
            }
        }
    }
}

func setSys() {
    if sendFrom < sendTo {
        return downSys++
    }else {
        return downSys--
    }
}

func checkLoop() {
    if downSys == upSys + 1 {
            if loop {
                LogTelegram("I", "What a feeling! My mens and womens and myself think a lot about these expeditions!")
                LogTelegram("I", "But there is still so much unknown in these galaxies.")
                LogTelegram("I", "We have to explore the sectors again! Time was running out! On another round!")
                downSys = sendFrom
            }else {
                LogTelegram("I", "It was an honor to explore this unknown vastness of the galaxies!")
                LogTelegram("I", "But our troops need to recover! The ships have to be brought up to scratch again.")
                LogTelegram("I", "We thank the entire team for this excellent work and say goodbye for now!")
                StopScript(__FILE__)
            }
        }
}

func customSleep(sleepTime) {
    if sleepTime <= 0 {
        sleepTime = Random(5, 10)
    }
    LogTelegram("I", "Wait " + ShortDur(sleepTime + 10))
    Sleep((sleepTime + 10) * 1000)
}

func fleetsGet() {
    fleets, slots = GetFleets()
    expoSlots = slots.ExpTotal
    expoInUse = slots.ExpInUse
    freeSlots = slots.Total - slots.InUse - reserved
    usedSlots = 0
    for fleet in fleets {
        if fleet.Destination == NewCoordinate(myGalaxy, debrisSys, 16, DEBRIS_TYPE) || fleet.Mission == EXPEDITION {
            usedSlots++
        }
    }
}

func errorHandler(){
    LogTelegram("D", "Call function errorHandler()")

    myPlanet, err = GetCelestial(home)
    sys = myPlanet.Coordinate.System
    if err != nil {
        LogTelegram("E", home + " is not one of your planet/moon")
        StopScript(__FILE__)
    }
    if sendFrom == 0 || sendFrom == nil {
        LogTelegram("E", sendFrom + " is not a system that we can find")
        StopScript(__FILE__)
    }
    if sys > 499{
        if sendTo > 550 || sendTo == nil{
            LogTelegram("E", sendTo + " is not a system that we can find")
            StopScript(__FILE__)
        }
    }else {
        if sendTo > 499 || sendTo == nil{
            LogTelegram("E", sendTo + " is not a system that we can find")
            StopScript(__FILE__)
        }
    }
    if rankOnePoints == 0 || rankOnePoints == nil{
        LogTelegram("E", "We do not believe that the first place has " + rankOnePoints + " points. We at least set it at our level")
        rankOnePoints = GetCachedPlayer().Points
        LogTelegram("E", Dotify(rankOnePoints))
    }
    if maxSlotsUse == 0 || maxSlotsUse == nil{
        LogTelegram("E", "We are not allowed to use slots? Please set maxSlotsUse at least to 1")
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
        if TELEGRAM_CHAT_ID == nil || TELEGRAM_CHAT_ID == 0 {
            LogTelegram("W", "your Telegram ID is not set!")
        }
    }
}

func startExpo() {
    t = 0
    fleets, slots = GetFleets()
    for fleet in fleets {
        if t < fleet.BackIn {
            t = fleet.BackIn
        }
    }
    return t
}

func GalaxyGet() {
    LogTelegram("D", "Call function GalaxyGet()")
    myPlanet, err = GetCelestial(home)
    myGalaxy = myPlanet.Coordinate.Galaxy
    LogTelegram("D", "Galaxy " + myGalaxy + " was set Sir!")
}

func enoughtShips(shipNames, shipAmounts) {
    celt, cErr = GetCelestial(home)
    if cErr != nil {
        LogTelegram("E", cErr)
        StopScript(__FILE__)
    }
    eShips, sErr = celt.GetShips()
    if sErr != nil {
        LogTelegram("E", sErr)
        StopScript(__FILE__)
    }else {
        if sendAll && hartDebris {
            shipsMax = Floor(eShips.ByID(shipNames) / maxSlotsUse)
            LogTelegram("D", "Send " + Dotify(shipsMax) + " " + ID2Str(shipNames) + " per wave")
            return shipsMax
        }
    if eShips.ByID(shipNames) < shipAmounts {
            LogTelegram("W", "Not enought " + ID2Str(shipNames))
            shipsMax = Floor(eShips.ByID(shipNames) / maxSlotsUse)
            LogTelegram("W", "We set it do " + shipsMax + " " + ID2Str(shipNames))
            return shipsMax
        }else {
            shipsMax = shipAmounts / maxSlotsUse
            return shipsMax
        }
    }
}

func shipsSet() {
    LogTelegram("D", "Call function shipsSet()")
    if selfShips {
        LogTelegram("I", "Manual settings are adopted!")
        ships = NewShipsInfos()
        for shipName, shipAmount in ship {
            if shipAmount != 0 {
                amount = enoughtShips(shipName, shipAmount * maxSlotsUse)
                if amount > 0 {
                    ships.Set(shipName, amount)
                }else {
                    LogTelegram("E", "Not enought " + ID2Str(shipName))
                }
            }
        }
        LogTelegram("I", "The following ships have been imported: " + ships)
    }else {
        LogTelegram("I", "Automatically settings are adopted!")
        
        expoValue = 0

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
            SC = Ceil((expoValue * 2 * 200) / (5000 * (tech.HyperspaceTechnology * 0.05) + 1))
            LC = Ceil((expoValue * 2 * 200) / (25000 * (tech.HyperspaceTechnology * 0.05) + 1))
        }

        ships = NewShipsInfos()
        if smallCargo {
            tmp = enoughtShips(SMALLCARGO, SC * maxSlotsUse)
            ships.Set(SMALLCARGO, tmp)
            LogTelegram("I", "Add " + tmp + " small cargos")
        }else {
            tmp = enoughtShips(LARGECARGO, LC * maxSlotsUse)
            ships.Set(LARGECARGO, tmp)
            LogTelegram("I", "Add " + tmp + " large cargos")
        }
        ships.Set(ESPIONAGEPROBE, enoughtShips(ESPIONAGEPROBE, 1 * maxSlotsUse))
        ships.Set(DESTROYER, enoughtShips(DESTROYER, 1 * maxSlotsUse))
        if usePathfinder {
            ships.Set(PATHFINDER, enoughtShips(PATHFINDER, 1 * maxSlotsUse))
        }
        if HeavyFighter {
            ships.Set(HEAVYFIGHTER, enoughtShips(HEAVYFIGHTER, 3000 * maxSlotsUse))
        }
        if ships.ByID(SMALLCARGO) <= 0 && ships.ByID(LARGECARGO) <= 0 {
            LogTelegram("E", "No Caros set!")
            StopScript(__FILE__)
        }
    }
}

func setHartDebris() {
    LogTelegram("D", "Call function setHartDebris()")

    ships = NewShipsInfos()
    switch hartDebris {
    case LFandSC && sendAll:
        ships.Set(LIGHTFIGHTER, enoughtShips(LIGHTFIGHTER, 0))
        ships.Set(SMALLCARGO, enoughtShips(SMALLCARGO, 0))
    case !LFandSC && sendAll:
        ships.Set(HEAVYFIGHTER, enoughtShips(HEAVYFIGHTER, 0))
        ships.Set(LARGECARGO, enoughtShips(LARGECARGO, 0))
    case LFandSC && !sendAll:
        SC = enoughtShips(SMALLCARGO, 999999 * maxSlotsUse)
        LogTelegram("I", "We add " + Dotify(SC) + " SmallCargos.") 
        ships.Set(SMALLCARGO, SC)
        LF = Ceil(enoughtShips(LIGHTFIGHTER, ((SC / 2) * maxSlotsUse)))
        ships.Set(LIGHTFIGHTER, LF)
        LogTelegram("I", "We add " + Dotify(LF) + " LightFighters.")
        ships.Set(ESPIONAGEPROBE, enoughtShips(ESPIONAGEPROBE, 1 * maxSlotsUse))
        ships.Set(DESTROYER, enoughtShips(DESTROYER, 1 * maxSlotsUse))
        if usePathfinder {
            ships.Set(PATHFINDER, enoughtShips(PATHFINDER, 1 * maxSlotsUse))
        }
    case !LFandSC && !sendAll:
        LC = enoughtShips(LARGECARGO, 999999 * maxSlotsUse)
        ships.Set(LARGECARGO, LC)
        LogTelegram("I", "We add " + Dotify(LC) + " LargeCargos.") 
        HF = Ceil(enoughtShips(HEAVYFIGHTER, ((LC / 1.5) * maxSlotsUse)))
        ships.Set(HEAVYFIGHTER, HF)
        LogTelegram("I", "We add " + Dotify(HF) + " HeavyFighters.")
        ships.Set(ESPIONAGEPROBE, enoughtShips(ESPIONAGEPROBE, 1 * maxSlotsUse))
        ships.Set(DESTROYER, enoughtShips(DESTROYER, 1 * maxSlotsUse))
        if usePathfinder {
            ships.Set(PATHFINDER, enoughtShips(PATHFINDER, 1 * maxSlotsUse))
        }
    }
}

func LogTelegram(cat, message) {
    switch cat {
        case cat == "W" || cat == "w":
            LogWarn(message)
            if useTelegram {
                SendTelegram(TELEGRAM_CHAT_ID, "Warn [" + message + "]")
            }
        case cat == "I" || cat == "i":
            LogInfo(message)
            if useTelegram {
                SendTelegram(TELEGRAM_CHAT_ID, "Info [" + message + "]")
            }
        case cat == "E" || cat == "e":
            LogError(message)
            if useTelegram {
                SendTelegram(TELEGRAM_CHAT_ID, "Error [" + message + "]")
            }
        case cat == "D" || cat == "d":
            LogDebug(message)
        default:
            LogError("no valid entry! Message is: " + message)
            if useTelegram {
                SendTelegram(TELEGRAM_CHAT_ID, "Error [No valid entry! Message is: " + message + "]")
            }
    }
}

func doWork(){
    LogTelegram("D", "Call function doWork()")
    errorHandler()                                                              //checking any value error

    downSys = sendFrom
    upSys = sendTo

    if !ignoreFleets {
        LogTelegram("I", "Fleets are on the way! Wait until they are back!")
        LogTelegram("D",  "If you don't want this! Set ignoreFleets = true")
        customSleep(startExpo())                                                //checking fleets on going and wait until they are back 
    }

    GalaxyGet()                                                                 //get the Galaxy form home
    
    if hartDebris {
        setHartDebris()                                                         //set Ships automatically or manual for get the hartest debris
    }else {
        shipsSet()                                                              //set Ships automatically or manual
    }
    
    for {
        fleetsGet()                                                             //get slots

        setMinSecs()                                                            //set time
        
        fleetsGet()                                                             //calc free slots again
        LogTelegram("D", "Expeditions available: " + (expoSlots - expoInUse))
        LogTelegram("D", "Slots available: " + freeSlots)
        
        doExpo()                                                                //send expo
        
        checkLoop()                                                             //checking is loop when he is finished
        
        if useWave && send {                                                    //wave tactics
            debrisSys = downSys
            downSys = setSys() 
        }

        customSleep(minSecs)                                                    //set sleeper
        
        sendMineDebris()                                                        //after sleep checking debris
    }
}
doWork()
