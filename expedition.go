/* Creaded by Bull
 * https://github.com/BullAzumi/Bull
 */
 
/*
 * VERSION 1.13
 */
 

/* DESCRIPTION
 * Automatically sends expeditions when a slot or expedition slot is free.
 * You can decide whether to choose the ship automatically or by yourself.
 */

/*---------------------------------------------------------------------------------------------------------------------------*/

//######################################## SETTINGS START ########################################


home = "M:4:363:4"                                      //set your expo home
sendFrom = 353                                          //system lower limit
sendTo = 373                                            //system upper limit
rankOnePoints = 1342123587                              //how many points does the first place have
smallCargo = false                                      //would you fly with small cargos? true = yes / false = no
usePathfinder = true                                    //Use pathfinder in automatically build
HeavyFighter = true                                     //should add HeavyFighters? yes = true / no = false
maxSlotsUse = 7                                         //how many slots should we use?
loop = false                                            //should we fly every round again?
expoTime = 1                                            //set duration
selfShips = false                                       //assemble the ships yourself or have them calculated automatically true = self / false = automatic
useWave = true                                          //should we let everyone fly on a coordinate and only then switch? true = yes / false = no

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


settings = GetServer().Settings
speed = settings.EconomySpeed
myGalaxy = 0
SC = 0
LC = 0
tmpsend = 0
ships = {}

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

func sendExpo() {
    LogDebug("Call function sendExpo()")
    
    sTo = NewCoordinate(myGalaxy, sendFrom, 16, PLANET_TYPE)        
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

func doWork(){
    LogDebug("Call function doWork()")

    //Checking errors
    errorHandler()

    //Check if fleets are still on the way when starting the script because of SaveFleet Script
    waitTime = startExpo()

    LogInfo("Wait until all fleets are back" + ShortDur(waitTime + 10))
    Sleep((waitTime + 10) * 1000)

    //get Galaxy and set ships
    GalaxyGet()
    shipsSet()
    
    for {
        //get slots
        fleets, slots = GetFleets()
        reserved = GetFleetSlotsReserved()
        //set time
        minSecs = 0
        for fleet in fleets {
            if fleet.Mission == EXPEDITION {
                if minSecs < fleet.BackIn {
                    minSecs = fleet.BackIn
                }
            }
        }
        //calc free slots
        expoSlots = slots.ExpTotal
        expoInUse = slots.ExpInUse
        freeSlots = slots.Total - slots.InUse - reserved
        LogDebug("Expeditions available: " + expoSlots - expoInUse)
        LogDebug("Slots available: " + freeSlots)
        
        Sleep(4000)
        //send fleet
        for expoInUse < maxSlotsUse{
            for {
                if expoInUse < expoSlots && freeSlots > 0 {
                    break;
                }
                expoSlots = slots.ExpTotal
                expoInUse = slots.ExpInUse
                freeSlots = slots.Total - slots.InUse - reserved
            }
            if sendFrom <= sendTo {
                newFleet, err = sendExpo()
                if err != nil {
                    LogError("Fleet reports errors: " + err)
                    contine
                }else{
                    LogInfo("Fleet starts on the expedition! " + newFleet)
                    if minSecs < newFleet.BackIn {
                        minSecs = newFleet.BackIn
                    }
                    expoInUse++
                    if !useWave {
                        tmpsend++
                        sendFrom++ 
                    }
                    Sleep(Random(3,7)*1000)
                }
            }
        }
            if sendFrom == sendTo + 1 {
                if loop {
                    LogInfo("What a feeling! My mens and womens and myself think a lot about these expeditions!")
                    LogInfo("But there is still so much unknown in these galaxies.")
                    LogInfo("We have to explore the sectors again! Time was running out! On another round!")
                    sendFrom -= tmpsend
                }else {
                    LogInfo("It was an honor to explore this unknown vastness of the galaxies!")
                    LogInfo("But our troops need to recover! The ships have to be brought up to scratch again.")
                    LogInfo("We thank the entire team for this excellent work and say goodbye for now!")
                    StopScript(__FILE__)
                }
            }        
        if useWave {
            tmpsend++
            sendFrom++ 
        } 

        LogInfo("Checking again in " + ShortDur(minSecs + 10))
        Sleep((minSecs + 10)*1000)
    }
}
doWork()
