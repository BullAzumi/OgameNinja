/* Creaded by Bull
 * https://github.com/BullAzumi/Bull
 */
 
/*
 * VERSION 1.2
 */

/* DESCRIPTION
 * Automatically sends expeditions when a slot or expedition slot is free.
 * You can decide whether to choose the ship automatically or by yourself.
 */

/*---------------------------------------------------------------------------------------------------------------------------*/

//######################################## SETTINGS START ########################################


home = "M:4:363:4"                                      //set your expo home
sendFrom = 363                                          //system lower limit
sendTo = 364                                            //system upper limit
rankOnePoints = 1234567890                              //how many points does the first place have
smallCargo = false                                      //would you fly with small cargos? true = yes / false = no
usePathfinder = true                                    //Use pathfinder in automatically build
HeavyFighter = true                                     //should add HeavyFighters? yes = true / no = false
maxSlotsUse = 1                                         //how many slots should we use?
loop = false                                            //should we fly every round again?
expoTime = 1                                            //set duration
selfShips = false                                       //assemble the ships yourself or have them calculated automatically true = self / false = automatic
useWave = true                                          //should we let everyone fly on a coordinate and only then switch? true = yes / false = no
ignoreFleets = false                                    //ignores fleets are on move true = yes / false = no

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


//debris functions
func sendMineDebris() {
    LogDebug("Call function sendMineDebris()")
    if mineDebris && send && useWave {
        PathfinderNeeded = scanGala()
        fleetsGet()
        if PathfinderNeeded > 0 {
            for {
                if expoInUse < expoSlots && freeSlots > 0 {
                    f, err = sendDebris(PathfinderNeeded)
                    if err != nil {
                        LogError(err)
                    }else {
                        LogInfo("Debris in work!")
                    }
                }else {
                    fleetsGet()
                }
            }
        }
    }
}

func scanGala() {
    LogDebug("Call function scanGala()")
    systemInfo, err = GalaxyInfos(myGalaxy, debrisSys)
    
    if err != nil {
        LogDebug(err)
    }
    if systemInfo.ExpeditionDebris.PathfindersNeeded > 0 {
        LogDebug("Debris detected")
        LogDebug("Need " + systemInfo.ExpeditionDebris.PathfindersNeeded + " Pathfinder")
        return systemInfo.ExpeditionDebris.PathfindersNeeded
    }else {
        LogDebug("No Debris detected")
        return systemInfo.ExpeditionDebris.PathfindersNeeded
    }
}

func sendDebris(shipCounter) {
    LogDebug("Call function sendDebris()")
    
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
    for expoInUse < maxSlotsUse{
        for {
            if expoInUse < expoSlots && freeSlots > 0 {
                break;
            }
            fleetsGet()
        }
        if downSys <= upSys || downSys > upSys {
            newFleet, err = sendExpo()
            if err != nil {
                LogError("Fleet reports errors: " + err)
                send = false
                continue
            }else{
                LogInfo("Fleet starts on the expedition! " + newFleet)
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
    LogDebug("Call function sendExpo()")
    
    sTo = NewCoordinate(myGalaxy, downSys, 16, PLANET_TYPE)        
    fleet = NewFleet()
    fleet.SetOrigin(home)
    fleet.SetDestination(sTo)
    LogInfo("Coords set to " + sTo )
    fleet.SetSpeed(HUNDRED_PERCENT)
    fleet.SetMission(EXPEDITION)
    fleet.SetShips(*ships)
    fleet.SetDuration(expoTime)
    return fleet.SendNow()
}

func setMinSecs() {
    LogDebug("Call function setMinSecs()")
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
                LogInfo("What a feeling! My mens and womens and myself think a lot about these expeditions!")
                LogInfo("But there is still so much unknown in these galaxies.")
                LogInfo("We have to explore the sectors again! Time was running out! On another round!")
                downSys = sendFrom
            }else {
                LogInfo("It was an honor to explore this unknown vastness of the galaxies!")
                LogInfo("But our troops need to recover! The ships have to be brought up to scratch again.")
                LogInfo("We thank the entire team for this excellent work and say goodbye for now!")
                StopScript(__FILE__)
            }
        }
}

func customSleep(sleepTime) {
    if sleepTime <= 0 {
        sleepTime = Random(5, 10)
    }
    LogInfo("Wait " + ShortDur(sleepTime + 10))
    Sleep((sleepTime + 10) * 1000)
}

func fleetsGet() {
    fleets, slots = GetFleets()
    expoSlots = slots.ExpTotal
    expoInUse = slots.ExpInUse
    freeSlots = slots.Total - slots.InUse - reserved
}

func errorHandler(){
    LogDebug("Call function errorHandler()")

    myPlanet, err = GetCelestial(home)
    sys = myPlanet.Coordinate.System
    if err != nil {
        LogError(home + " is not one of your planet/moon")
        StopScript(__FILE__)
    }
    if sendFrom == 0 || sendFrom == nil {
        LogError(sendFrom + " is not a system that we can find")
        StopScript(__FILE__)
    }
    if sys > 499{
        if sendTo > 550 || sendTo == nil{
            LogError(sendTo + " is not a system that we can find")
            StopScript(__FILE__)
        }
    }else {
        if sendTo > 499 || sendTo == nil{
            LogError(sendTo + " is not a system that we can find")
            StopScript(__FILE__)
        }
    }
    if rankOnePoints == 0 || rankOnePoints == nil{
        LogError("We do not believe that the first place has " + rankOnePoints + " points. We at least set it at our level")
        rankOnePoints = GetCachedPlayer().Points
        LogError(Dotify(rankOnePoints))
    }
    if maxSlotsUse == 0 || maxSlotsUse == nil{
        LogError("We are not allowed to use slots? Please set maxSlotsUse at least to 1")
        StopScript(__FILE__)
    }
    if expoTime <= 0 || expoTime > 18 || expoTime == nil{
        LogError("we can not stay in expedition for " + expoTime + " hour")
        if expoTime <= 0 || expoTime == nil {
            LogError("we set it at least do 1 hour")
            expoTime = 1
        }else {
            LogError("we set it at least do 18 hour")
            expoTime = 18
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
    LogDebug("Call function GalaxyGet()")
    myPlanet, err = GetCelestial(home)
    myGalaxy = myPlanet.Coordinate.Galaxy
    LogInfo("Galaxy " + myGalaxy + " was set Sir!")
}

func shipsSet() {
    LogDebug("Call function shipsSet()")
    if selfShips {
        LogInfo("Manual settings are adopted!")
        Sleep(3000)
        ships = NewShipsInfos()
        for shipName, shipAmount in ship {        
            ships.Set(shipName, shipAmount)
        }
        LogDebug("The following ships have been imported:")
        LogDebug(ships)
    }else {
        LogInfo("Automatically settings are adopted!")
        Sleep(3000)

        switch rankOnePoints {
            case rankOnePoints >= 100000000:
                SC = 1250 * speed
                LC = 417 * speed
            case rankOnePoints >= 75000000:
                SC = 1050 * speed
                LC = 350 * speed
            case rankOnePoints >= 50000000:
                SC = 900 * speed
                LC = 300 * speed
            case rankOnePoints >= 25000000:
                SC = 750 * speed
                LC = 250 * speed
            case rankOnePoints >= 5000000:
                SC = 600 * speed
                LC = 200 * speed
            case rankOnePoints >= 1000000:
                SC = 450 * speed
                LC = 150 * speed
            case rankOnePoints >= 100000:
                SC = 300 * speed
                LC = 100 * speed
            case rankOnePoints >= 0:
                SC = 125 * speed
                LC = 42 * speed
        }

        ships = NewShipsInfos()
        if smallCargo {
            ships.Set(SMALLCARGO, SC)
            LogDebug("Add " + SC + " small cargos")
        }else {
            ships.Set(LARGECARGO, LC)
            LogDebug("Add " + LC + " large cargos")
        }
        ships.Set(ESPIONAGEPROBE, 1)
        ships.Set(DESTROYER, 1)
        if usePathfinder {
            ships.Set(PATHFINDER, 1)
        }
        if HeavyFighter {
            ships.Set(HEAVYFIGHTER, 3000)
        }
    }
}

func doWork(){
    LogDebug("Call function doWork()")
    errorHandler()                                                              //checking any value error

    downSys = sendFrom
    upSys = sendTo

    if !ignoreFleets {
        customSleep(startExpo())                                                //checking fleets on going and wait until they are back 
    }

    GalaxyGet()                                                                 //get the Galaxy form home
    shipsSet()                                                                  //set Ships automatically or manual
    
    for {
        fleetsGet()                                                             //get slots

        minSecs = setMinSecs()                                                  //set time
        
        fleetsGet()                                                             //calc free slots again
        LogDebug("Expeditions available: " + (expoSlots - expoInUse))
        LogDebug("Slots available: " + freeSlots)
        
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
