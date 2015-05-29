package main

//////////// LIBRARIES
import (
    "encoding/json"
    "fmt"
    "net"
    "io/ioutil"
    "time"
    "math/rand"
    "strconv"
    "os"
)

//////////// STRUCTURES
//// Army structure
type ARMY struct {
    Catapults int
    Archers int
    Cavalry int
    Spearmen int
    Infantry int
}
//// Officer info structure
type OFFICERBASE struct {
    Name string
    Address string
}
//// Officer structure
type OFFICER struct {
    Credentials OFFICERBASE
    Army ARMY
    State int
}
//// Message structure
type MSG struct {
    Army ARMY
    To string
    From string
    Message string
    State int
    Messages []MSG
}
//// Officer cut structure
type OFFICERCUT struct {
    Army ARMY
    Messages []MSG
}
//// Cut structure
type CUT struct {
    Brutus OFFICERCUT
    Operachorus OFFICERCUT
    Pompus OFFICERCUT
    State int
    CaesarArmy ARMY
    CaesarMessages []MSG
    HasBrutus bool
    HasOperachorus bool
    HasPompus bool
}

//////////// DEFINE GLOBAL VARIABLES
var (
    CutState int                    // the event state to take the cut at, randomly choosen
    CurrentState int                // the current event state for Caesar
    
    Cut CUT                         // the cut information
    
    Totals ARMY                     // the total size of Caesar's forces
    Objective ARMY                  // the approximate amount of troops each officer should have
    Remainder ARMY                  // the extra troops available to play around with
    
    Caesar OFFICER                  // information for Caesar
    Brutus OFFICER                  // information for Brutus
    Operachorus OFFICER             // information for Operachorus
    Pompus OFFICER                  // information for Pompus
        
    CaesarPort *net.UDPConn         // the port to connect to Caesar
    BrutusPort *net.UDPConn         // the port to connect to Brutus
    OperachorusPort *net.UDPConn    // the port to connect to Operachorus
    PompusPort *net.UDPConn         // the port to connect to Pompus
    CutPort *net.UDPConn            // the port to connect to for Cuts
    BrutusCutPort *net.UDPConn      // the port to connect to for Brutus' cuts
    OperachorusCutPort *net.UDPConn // the port to connect to for Operachorus' cuts
    PompusCutPort *net.UDPConn      // the port to connect to for Pompus' cuts
)

//////////// MAIN FUNCTION
func main() {
    //////// SEED RAND AND DETERMINE WHEN TO TAKE THE CUT
    fmt.Println("Caesar is determining when to take a cut...")
    rand.Seed(time.Now().UTC().UnixNano())
    if (len(os.Args) == 1) {
        CutState = rand.Intn(6)
    } else {
        i, err := strconv.Atoi(os.Args[1])
        if err != nil {
            CutState = rand.Intn(6)
        } else {
            CutState = i
        }
    }
    CurrentState = 0
    
    //////// SETUP OFFICERS
    fmt.Println("Caesar is setting up his tables...")
    SetupOfficers()
    //////// SETUP NETWORKING
    fmt.Println("Caesar is choosing his messenger's routes...")
    SetupNetworking()
    //////// GATHER BASIC ARMY INFO
    fmt.Printf("Caesar is gathering basic intelligence on his armies...")
    GatherDispositions()
    PrintOfficer(Caesar, 0)
    PrintOfficer(Brutus, 0)
    PrintOfficer(Operachorus, 0)
    PrintOfficer(Pompus, 0)
    //hasDoneCatapults := false
    //hasDoneArchers := false
    //hasDoneCavalry := false
    //hasDoneSpearmen := false
    //hasDoneInfantry := false
    doneBrutus := false
    doneOperachorus := false
    donePompus := false
    doneCaesar := false
    
    //////// MAIN LOOP
    fmt.Println("Caesar is distributing his troops...")
    for {
        ////// check if we need to take a cut
        if (CurrentState == CutState) {
            go CutRoutine()
        }
        
        ////// REDISTRIBUTE TROOPS
        if IsDistributed(Brutus.Army) == false {
            //// Brutus
            fmt.Println("Brutus!")
            doneBrutus = DistributeArmy(Brutus, Operachorus, Pompus, Caesar, OperachorusPort, PompusPort, CaesarPort)
        } else if IsDistributed(Operachorus.Army) == false {
            //// Operachorus
            fmt.Println("Operachorus!")
            doneOperachorus = DistributeArmy(Operachorus, Brutus, Pompus, Caesar, BrutusPort, PompusPort, CaesarPort)
        } else if IsDistributed(Pompus.Army) == false {
            //// Pompus
            fmt.Println("Pompus!")
            donePompus = DistributeArmy(Pompus, Operachorus, Brutus, Caesar, OperachorusPort, BrutusPort, CaesarPort)
        } else if IsDistributed(Caesar.Army) == false {
            //// Caesar
            fmt.Println("Caesar!")
            doneCaesar = DistributeArmy(Caesar, Operachorus, Brutus, Pompus, OperachorusPort, BrutusPort, PompusPort)
        } else {
            fmt.Println("Done!")
            break
        }
        
        AskDispositions()
        
        if doneOperachorus && donePompus && doneBrutus && doneCaesar {
            fmt.Println("Done!")
            break
        }
    }

    fmt.Println("The Legions are ready!")
    SayReady()
    PrintTotals()
    PrintOfficer(Caesar, 0)
    PrintOfficer(Brutus, 0)
    PrintOfficer(Operachorus, 0)
    PrintOfficer(Pompus, 0)
}

//////////// FUNCTIONS
func DistributeArmy(mainOfficer OFFICER, officerA OFFICER, officerB OFFICER, officerC OFFICER, officerAPort *net.UDPConn, officerBPort *net.UDPConn, officerCPort *net.UDPConn) bool {
    message := MSG {
        Army: ARMY {
                    Catapults: 0,
                    Archers: 0,
                    Cavalry: 0,
                    Spearmen: 0,
                    Infantry: 0,
                },
        From: "Caesar",
        Message: "SendTroopsTo",
        State: CurrentState,
        Messages: nil,
    }
    var name string
    var troopstr string
    var port *net.UDPConn
    ////// DISTRIBUTION
    if IsEvelyDistributed(mainOfficer.Army.Catapults, Objective.Catapults, Remainder.Catapults) == false {
        //// Catapults
        main := HasSpareUnits(mainOfficer.Army.Catapults, Objective.Catapults, Remainder.Catapults)
        a := HasSpareUnits(officerA.Army.Catapults, Objective.Catapults, Remainder.Catapults)
        b := HasSpareUnits(officerB.Army.Catapults, Objective.Catapults, Remainder.Catapults)
        c := HasSpareUnits(officerC.Army.Catapults, Objective.Catapults, Remainder.Catapults)
        
        if main < 0 {
            // Main Officer has too many Catapults!
            if (a > 0 || (a == 0 && Remainder.Catapults > 1 && b <= 0 && c <= 0)) {
                message.To = mainOfficer.Credentials.Name
                message.Army.Catapults = main
                port = GetPortFromName(mainOfficer.Credentials.Name)
                name = officerA.Credentials.Name
                troopstr = strconv.Itoa(-main)
            } else if (b > 0 || (b == 0 && Remainder.Catapults > 1 && a <= 0 && c <= 0)) {
                message.To = mainOfficer.Credentials.Name
                message.Army.Catapults = main
                port = GetPortFromName(mainOfficer.Credentials.Name)
                name = officerB.Credentials.Name
                troopstr = strconv.Itoa(-main)
            } else if (c > 0 || (c == 0 && Remainder.Catapults > 1 && b <= 0 && a <= 0)) {
                message.To = mainOfficer.Credentials.Name
                message.Army.Catapults = main
                port = GetPortFromName(mainOfficer.Credentials.Name)
                name = officerC.Credentials.Name
                troopstr = strconv.Itoa(-main)
            }
        } else if (a < 0) {
            // Officer A has too many Catapults!
            message.To = officerA.Credentials.Name
            message.Army.Catapults = a
            port = GetPortFromName(officerA.Credentials.Name)
            name = mainOfficer.Credentials.Name
            troopstr = strconv.Itoa(-a)
        } else if (b < 0) {
            // Officer B has too many Catapults!
            message.To = officerB.Credentials.Name
            message.Army.Catapults = b
            port = GetPortFromName(officerB.Credentials.Name)
            name = mainOfficer.Credentials.Name
            troopstr = strconv.Itoa(-b)
        } else if (c < 0) {
            // Officer C has too many Catapults!
            message.To = officerC.Credentials.Name
            message.Army.Catapults = c
            port = GetPortFromName(officerC.Credentials.Name)
            name = mainOfficer.Credentials.Name
            troopstr = strconv.Itoa(-c)
        }
        troopstr += " Catapults"
    } else if IsEvelyDistributed(mainOfficer.Army.Archers, Objective.Archers, Remainder.Archers) == false {
        //// Archers
        main := HasSpareUnits(mainOfficer.Army.Archers, Objective.Archers, Remainder.Archers)
        a := HasSpareUnits(officerA.Army.Archers, Objective.Archers, Remainder.Archers)
        b := HasSpareUnits(officerB.Army.Archers, Objective.Archers, Remainder.Archers)
        c := HasSpareUnits(officerC.Army.Archers, Objective.Archers, Remainder.Archers)
        
        if main < 0 {
            // Main Officer has too many Archers!
            if (a > 0 || (a == 0 && Remainder.Archers > 1 && b <= 0 && c <= 0)) {
                message.To = mainOfficer.Credentials.Name
                message.Army.Archers = main
                port = GetPortFromName(mainOfficer.Credentials.Name)
                name = officerA.Credentials.Name
                troopstr = strconv.Itoa(-main)
            } else if (b > 0 || (b == 0 && Remainder.Archers > 1 && a <= 0 && c <= 0)) {
                message.To = mainOfficer.Credentials.Name
                message.Army.Archers = main
                port = GetPortFromName(mainOfficer.Credentials.Name)
                name = officerB.Credentials.Name
                troopstr = strconv.Itoa(-main)
            } else if (c > 0 || (c == 0 && Remainder.Archers > 1 && b <= 0 && a <= 0)) {
                message.To = mainOfficer.Credentials.Name
                message.Army.Archers = main
                port = GetPortFromName(mainOfficer.Credentials.Name)
                name = officerC.Credentials.Name
                troopstr = strconv.Itoa(-main)
            }
        } else if (a < 0) {
            // Officer A has too many Archers!
            message.To = officerA.Credentials.Name
            message.Army.Archers = a
            port = GetPortFromName(officerA.Credentials.Name)
            name = mainOfficer.Credentials.Name
            troopstr = strconv.Itoa(-a)
        } else if (b != 0) {
            // Officer B has too many Archers!
            message.To = officerB.Credentials.Name
            message.Army.Archers = b
            port = GetPortFromName(officerB.Credentials.Name)
            name = mainOfficer.Credentials.Name
            troopstr = strconv.Itoa(-b)
        } else if (c < 0) {
            // Officer C has too many Archers!
            message.To = officerC.Credentials.Name
            message.Army.Archers = c
            port = GetPortFromName(officerC.Credentials.Name)
            name = mainOfficer.Credentials.Name
            troopstr = strconv.Itoa(-c)
        }
        troopstr += " Archers"
    } else if IsEvelyDistributed(mainOfficer.Army.Cavalry, Objective.Cavalry, Remainder.Cavalry) == false {
        //// Cavalry
        main := HasSpareUnits(mainOfficer.Army.Cavalry, Objective.Cavalry, Remainder.Cavalry)
        a := HasSpareUnits(officerA.Army.Cavalry, Objective.Cavalry, Remainder.Cavalry)
        b := HasSpareUnits(officerB.Army.Cavalry, Objective.Cavalry, Remainder.Cavalry)
        c := HasSpareUnits(officerC.Army.Cavalry, Objective.Cavalry, Remainder.Cavalry)
        
        if main < 0 {
            // Main Officer has too many Cavalry!
            if (a > 0 || (a == 0 && Remainder.Cavalry > 1 && b <= 0 && c <= 0)) {
                message.To = mainOfficer.Credentials.Name
                message.Army.Cavalry = main
                port = GetPortFromName(mainOfficer.Credentials.Name)
                name = officerA.Credentials.Name
                troopstr = strconv.Itoa(-main)
            } else if (b > 0 || (b == 0 && Remainder.Cavalry > 1 && a <= 0 && c <= 0)) {
                message.To = mainOfficer.Credentials.Name
                message.Army.Cavalry = main
                port = GetPortFromName(mainOfficer.Credentials.Name)
                name = officerB.Credentials.Name
                troopstr = strconv.Itoa(-main)
            } else if (c > 0 || (c == 0 && Remainder.Cavalry > 1 && b <= 0 && a <= 0)) {
                message.To = mainOfficer.Credentials.Name
                message.Army.Cavalry = main
                port = GetPortFromName(mainOfficer.Credentials.Name)
                name = officerC.Credentials.Name
                troopstr = strconv.Itoa(-main)
            }
        } else if (a < 0) {
            // Officer A has too many Cavalry!
            message.To = officerA.Credentials.Name
            message.Army.Cavalry = a
            port = GetPortFromName(officerA.Credentials.Name)
            name = mainOfficer.Credentials.Name
            troopstr = strconv.Itoa(-a)
        } else if (b < 0) {
            // Officer B has too many Cavalry!
            message.To = officerB.Credentials.Name
            message.Army.Cavalry = b
            port = GetPortFromName(officerB.Credentials.Name)
            name = mainOfficer.Credentials.Name
            troopstr = strconv.Itoa(-b)
        } else if (c < 0) {
            // Officer C has too many Cavalry!
            message.To = officerC.Credentials.Name
            message.Army.Cavalry = c
            port = GetPortFromName(officerC.Credentials.Name)
            name = mainOfficer.Credentials.Name
            troopstr = strconv.Itoa(-c)
        }
        troopstr += " Cavalry"
    } else if IsEvelyDistributed(mainOfficer.Army.Spearmen, Objective.Spearmen, Remainder.Spearmen) == false {
        //// Spearmen
        main := HasSpareUnits(mainOfficer.Army.Spearmen, Objective.Spearmen, Remainder.Spearmen)
        a := HasSpareUnits(officerA.Army.Spearmen, Objective.Spearmen, Remainder.Spearmen)
        b := HasSpareUnits(officerB.Army.Spearmen, Objective.Spearmen, Remainder.Spearmen)
        c := HasSpareUnits(officerC.Army.Spearmen, Objective.Spearmen, Remainder.Spearmen)
        
        if main < 0 {
            // Main Officer has too many Spearmen!
            if (a > 0 || (a == 0 && Remainder.Spearmen > 1 && b <= 0 && c <= 0)) {
                message.To = mainOfficer.Credentials.Name
                message.Army.Spearmen = main
                port = GetPortFromName(mainOfficer.Credentials.Name)
                name = officerA.Credentials.Name
                troopstr = strconv.Itoa(-main)
            } else if (b > 0 || (b == 0 && Remainder.Spearmen > 1 && a <= 0 && c <= 0)) {
                message.To = mainOfficer.Credentials.Name
                message.Army.Spearmen = main
                port = GetPortFromName(mainOfficer.Credentials.Name)
                name = officerB.Credentials.Name
                troopstr = strconv.Itoa(-main)
            } else if (c > 0 || (c == 0 && Remainder.Spearmen > 1 && b <= 0 && a <= 0)) {
                message.To = mainOfficer.Credentials.Name
                message.Army.Spearmen = main
                port = GetPortFromName(mainOfficer.Credentials.Name)
                name = officerC.Credentials.Name
                troopstr = strconv.Itoa(-main)
            }
        } else if (a < 0) {
            // Officer A has too many Spearmen!
            message.To = officerA.Credentials.Name
            message.Army.Spearmen = a
            port = GetPortFromName(officerA.Credentials.Name)
            name = mainOfficer.Credentials.Name
            troopstr = strconv.Itoa(-a)
        } else if (b < 0) {
            // Officer B has too many Spearmen!
            message.To = officerB.Credentials.Name
            message.Army.Spearmen = b
            port = GetPortFromName(officerB.Credentials.Name)
            name = mainOfficer.Credentials.Name
            troopstr = strconv.Itoa(-b)
        } else if (c < 0) {
            // Officer C has too many Spearmen!
            message.To = officerC.Credentials.Name
            message.Army.Spearmen = c
            port = GetPortFromName(officerC.Credentials.Name)
            name = mainOfficer.Credentials.Name
            troopstr = strconv.Itoa(-c)
        }
        troopstr += " Spearmen"
    } else if IsEvelyDistributed(mainOfficer.Army.Infantry, Objective.Infantry, Remainder.Infantry) == false {
        //// Infantry
        main := HasSpareUnits(mainOfficer.Army.Infantry, Objective.Infantry, Remainder.Infantry)
        a := HasSpareUnits(officerA.Army.Infantry, Objective.Infantry, Remainder.Infantry)
        b := HasSpareUnits(officerB.Army.Infantry, Objective.Infantry, Remainder.Infantry)
        c := HasSpareUnits(officerC.Army.Infantry, Objective.Infantry, Remainder.Infantry)
        
        if main < 0 {
            // Main Officer has too many Infantry!
            if (a > 0 || (a == 0 && Remainder.Infantry > 1 && b <= 0 && c <= 0)) {
                message.To = mainOfficer.Credentials.Name
                message.Army.Infantry = main
                port = GetPortFromName(mainOfficer.Credentials.Name)
                name = officerA.Credentials.Name
                troopstr = strconv.Itoa(-main)
            } else if (b > 0 || (b == 0 && Remainder.Infantry > 1 && a <= 0 && c <= 0)) {
                message.To = mainOfficer.Credentials.Name
                message.Army.Infantry = main
                port = GetPortFromName(mainOfficer.Credentials.Name)
                name = officerB.Credentials.Name
                troopstr = strconv.Itoa(-main)
            } else if (c > 0 || (c == 0 && Remainder.Infantry > 1 && b <= 0 && a <= 0)) {
                message.To = mainOfficer.Credentials.Name
                message.Army.Infantry = main
                port = GetPortFromName(mainOfficer.Credentials.Name)
                name = officerC.Credentials.Name
                troopstr = strconv.Itoa(-main)
            }
        } else if (a < 0) {
            // Officer A has too many Infantry!
            message.To = officerA.Credentials.Name
            message.Army.Infantry = a
            port = GetPortFromName(officerA.Credentials.Name)
            name = mainOfficer.Credentials.Name
            troopstr = strconv.Itoa(-a)
        } else if (b < 0) {
            // Officer B has too many Infantry!
            message.To = officerB.Credentials.Name
            message.Army.Infantry = b
            port = GetPortFromName(officerB.Credentials.Name)
            name = mainOfficer.Credentials.Name
            troopstr = strconv.Itoa(-b)
        } else if (c < 0) {
            // Officer C has too many Infantry!
            message.To = officerC.Credentials.Name
            message.Army.Infantry = c
            port = GetPortFromName(officerC.Credentials.Name)
            name = mainOfficer.Credentials.Name
            troopstr = strconv.Itoa(-c)
        }
        troopstr += " Infantry"
    }
    
    if (port != nil) {
        message.Message += name
        fmt.Println("Send "+troopstr+" to "+name+" from "+message.To)
        if GetNameFromPort(port) != "Caesar" {
            x := SendGetResponse(port, CaesarPort, message)
            if (x.To != "") {
                if (x.Message == "ReceiveTroops") {
                    Caesar.Army = ARMY {
                        Catapults: Caesar.Army.Catapults-message.Army.Catapults,
                        Archers: Caesar.Army.Archers-message.Army.Archers,
                        Cavalry: Caesar.Army.Cavalry-message.Army.Cavalry,
                        Spearmen: Caesar.Army.Spearmen-message.Army.Spearmen,
                        Infantry: Caesar.Army.Infantry-message.Army.Infantry,
                    }
                }
                CurrentState++
                return false
            }
        } else {
            message.Message = "ReceiveTroops"
            port = GetPortFromName(name)
            Caesar.Army = ARMY {
                Catapults: Caesar.Army.Catapults+message.Army.Catapults,
                Archers: Caesar.Army.Archers+message.Army.Archers,
                Cavalry: Caesar.Army.Cavalry+message.Army.Cavalry,
                Spearmen: Caesar.Army.Spearmen+message.Army.Spearmen,
                Infantry: Caesar.Army.Infantry+message.Army.Infantry,
            }
            CurrentState++
            x := SendGetResponse(port, CaesarPort, message)
            if (x.To != "") {
                CurrentState++
                return false
            }
            y := GetResponse(CaesarPort)
            if (y.To != "") {
                CurrentState++
                return false
            }
        }
    }
    return true
}
//// Get port from name
func GetPortFromName(name string) *net.UDPConn {
    if name == "Brutus" {
        return BrutusPort
    } else if name == "Operachorus" {
        return OperachorusPort
    } else if name == "Pompus" {
        return PompusPort
    } else if name == "Caesar" {
        return CaesarPort
    }
    return nil
}
//// Get name from port
func GetNameFromPort(port *net.UDPConn) string {
    if (port == BrutusPort) {
        return "Brutus"
    } else if (port == OperachorusPort) {
        return "Operachorus"
    } else if (port == PompusPort) {
        return "Pompus"
    } else if (port == CaesarPort) {
        return "Caesar"
    }
    return ""
}
//// Checks if an army is distributed
func IsDistributed(army ARMY) bool {
    // check Catapults
    if IsEvelyDistributed(army.Catapults, Objective.Catapults, Remainder.Catapults) == false {
        //fmt.Println("Cfalse")
        return false
    }
    // check Archers
    if IsEvelyDistributed(army.Archers, Objective.Archers, Remainder.Archers) == false {
        //fmt.Println("Afalse")
        return false
    }
    // check Cavalry
    if IsEvelyDistributed(army.Cavalry, Objective.Cavalry, Remainder.Cavalry) == false {
        //fmt.Println("Cvfalse")
        return false
    }
    // check Spearmen
    if IsEvelyDistributed(army.Spearmen, Objective.Spearmen, Remainder.Spearmen) == false {
        //fmt.Println("Sfalse")
        return false
    }
    // check Infantry
    if IsEvelyDistributed(army.Infantry, Objective.Infantry, Remainder.Infantry) == false {
        //fmt.Println("Ifalse")
        return false
    }

    return true
}
//// Checks if a army unit type is evenly distributed or not
func IsEvelyDistributed(count int, objectivecount int, remaindercount int) bool {
    if count != objectivecount {
        diff := objectivecount - count
        if diff == -1 && remaindercount != 0 {
            return true
        } else {
            return false
        }
    }
    return true
}
//// Checks if an army has spare units of the type pass
func HasSpareUnits(count int, objectivecount int, remaindercount int) int {
    if count != objectivecount {
        diff := objectivecount - count
        if remaindercount == 2 && diff == -2 {
            diff = -1
        }
        if remaindercount == 3 && diff == -3 {
            diff = -1
        }
        if remaindercount == 3 && diff == -2 {
            diff = -1
        }
        return diff
    }
    
    return 0;
}
//// Sends the message to the address specified, returns true if it succeeds
func SendMessage(conn *net.UDPConn, message MSG) bool {
    msg, err := json.Marshal(message)
    if err != nil {
        panic(err)
    }
    
    n, err := conn.Write(msg)
    if err != nil {
        panic(err)
    }
    
    // check if a cut!
    if CurrentState == CutState {
        Cut.CaesarMessages = append(Cut.CaesarMessages, message)
    }
    
    if n == 0 {
        return false
    }
    return true
}
//// Gets (and returns) message sent to the specified port
func GetResponse(list *net.UDPConn) MSG {
    mesbuf := make([]byte, 4096)
    nread, _, err := list.ReadFromUDP(mesbuf)
    if err != nil {
        panic(err)
    }
    mesbuf = mesbuf[:nread]
    
    var res MSG
    err = json.Unmarshal(mesbuf, &res)
    if err != nil {
        panic(err)
    }
    
    return res
}
//// Sends the message and returns the response
func SendGetResponse(conn *net.UDPConn, list *net.UDPConn, message MSG) MSG {
    SendMessage(conn, message)
    return GetResponse(list)
}
//// Prints the total number of troops available
func PrintTotals() {
    fmt.Println("TOTAL TROOPS AVAILABLE")
    fmt.Printf("Catapults: %d, each officer should have ~%dR%d\n", Totals.Catapults, Objective.Catapults, Remainder.Catapults)
    fmt.Printf("Archers: %d, each officer should have ~%dR%d\n", Totals.Archers, Objective.Archers, Remainder.Archers)
    fmt.Printf("Cavalry: %d, each officer should have ~%dR%d\n", Totals.Cavalry, Objective.Cavalry, Remainder.Cavalry)
    fmt.Printf("Spearmen: %d, each officer should have ~%dR%d\n", Totals.Spearmen, Objective.Spearmen, Remainder.Spearmen)
    fmt.Printf("Infantry: %d, each officer should have ~%dR%d\n", Totals.Infantry, Objective.Infantry, Remainder.Infantry)
}
//// Prints the number of troops each officer has at his disposal
func PrintOfficer(officer OFFICER, status int) {
    fmt.Printf("%s's Troops\n", officer.Credentials.Name);
    fmt.Printf("Catapults: %d\n", officer.Army.Catapults)
    fmt.Printf("Archers: %d\n", officer.Army.Archers)
    fmt.Printf("Cavalry: %d\n", officer.Army.Cavalry)
    fmt.Printf("Spearmen: %d\n", officer.Army.Spearmen)
    fmt.Printf("Infantry: %d\n", officer.Army.Infantry)
    
    if status == 1 {
        fmt.Printf("State: %d\n", officer.State)
    }
}
//// Setsup all the officers
func SetupOfficers() {
    // setup default nodes
    var caesarBase OFFICERBASE
    var brutusBase OFFICERBASE
    var operachorusBase OFFICERBASE
    var pompusBase OFFICERBASE
    // read settings from files
    file, err := ioutil.ReadFile("setcaesar")
    if err != nil {
        panic(err)
    } else {
        erra := json.Unmarshal(file, &caesarBase)
        if erra != nil {
            caesarBase = OFFICERBASE{
                Name: "Caesar",
                Address: "10.1.1.5",
            }
        }
    }
    file2, err2 := ioutil.ReadFile("setpompus")
    if err2 != nil {
        panic(err2)
    } else {
        err2a := json.Unmarshal(file2, &pompusBase)
        if err2a != nil {
            pompusBase = OFFICERBASE{
                Name: "Pompus",
                Address: "10.1.1.3",
            }
        }
    }
    file3, err3 := ioutil.ReadFile("setoperachorus")
    if err3 != nil {
        panic(err3)
    } else {
        err3a := json.Unmarshal(file3, &operachorusBase)
        if err3a != nil {
            operachorusBase = OFFICERBASE{
                Name: "Operachorus",
                Address: "10.1.1.4",
            }
        }
    }
    file4, err4 := ioutil.ReadFile("setbrutus")
    if err4 != nil {
        panic(err4)
    } else {
        err4a := json.Unmarshal(file4, &brutusBase)
        if err4a != nil {
            brutusBase = OFFICERBASE{
                Name: "Brutus",
                Address: "10.1.1.2",
            }
        }
    }
    
    // create officers
    var army ARMY
    army = ARMY{
        Catapults: 0,
        Archers: 0,
        Cavalry: 0,
        Spearmen: 0,
        Infantry: 0,
    }
    Brutus = OFFICER{
        Credentials: brutusBase,
        Army: army,
        State: 0,
    }
    Operachorus = OFFICER{
        Credentials: operachorusBase,
        Army: army,
        State: 0,
    }
    Pompus = OFFICER{
        Credentials: pompusBase,
        Army: army,
        State: 0,
    }
    cArmy := ARMY {
        Catapults: 1+rand.Intn(100),
        Archers: 1+rand.Intn(100),
        Cavalry: 1+rand.Intn(100),
        Spearmen: 1+rand.Intn(100),
        Infantry: 1+rand.Intn(100),
    }
    Caesar = OFFICER{
        Credentials: caesarBase,
        Army: cArmy,
        State: 0,
    }
}
//// Sets up the ports for the officers
func SetupNetworking() {
    fmt.Println("Caesar Address: " + Caesar.Credentials.Address)
    fmt.Println("Pompus Address: " + Pompus.Credentials.Address)
    fmt.Println("Brutus Address: " + Brutus.Credentials.Address)
    fmt.Println("Operachorus Address: " + Operachorus.Credentials.Address)
    //// SETUP PERSON PORTS (and main cut port)
    // setup cut port
    cutaddr, err := net.ResolveUDPAddr("udp", ":9000")
    if err != nil {
        panic(err)
    }
    cutport, err := net.ListenUDP("udp", cutaddr)
    if err != nil {
        panic(err)
    }
    CutPort = cutport
    // setup caesar port
    caesaraddr, err := net.ResolveUDPAddr("udp", ":9001")
    if err != nil {
        panic(err)
    }
    caesarport, err := net.ListenUDP("udp", caesaraddr)
    if err != nil {
        panic(err)
    }
    CaesarPort = caesarport
    // setup pompus port
    pompusaddr, err := net.ResolveUDPAddr("udp", Pompus.Credentials.Address+":9002")
    if err != nil {
        panic(err)
    }
    pompusport, err := net.DialUDP("udp", nil, pompusaddr)
    if err != nil {
        panic(err)
    }
    PompusPort = pompusport
    // setup operachorus port
    operachorusaddr, err := net.ResolveUDPAddr("udp", Operachorus.Credentials.Address+":9003")
    if err != nil {
        panic(err)
    }
    operachorusport, err := net.DialUDP("udp", nil, operachorusaddr)
    if err != nil {
        panic(err)
    }
    OperachorusPort = operachorusport
    // setup brutus port
    brutusaddr, err := net.ResolveUDPAddr("udp", Brutus.Credentials.Address+":9004")
    if err != nil {
        panic(err)
    }
    brutusport, err := net.DialUDP("udp", nil, brutusaddr)
    if err != nil {
        panic(err)
    }
    BrutusPort = brutusport
    
    //// SETUP CUT PORTS
    // setup pompus port
    pompusaddr2, err := net.ResolveUDPAddr("udp", Pompus.Credentials.Address+":9012")
    if err != nil {
        panic(err)
    }
    pompusport2, err := net.DialUDP("udp", nil, pompusaddr2)
    if err != nil {
        panic(err)
    }
    PompusCutPort = pompusport2
    // setup operachorus port
    operachorusaddr2, err := net.ResolveUDPAddr("udp", Operachorus.Credentials.Address+":9013")
    if err != nil {
        panic(err)
    }
    operachorusport2, err := net.DialUDP("udp", nil, operachorusaddr2)
    if err != nil {
        panic(err)
    }
    OperachorusCutPort = operachorusport2
    // setup brutus port
    brutusaddr2, err := net.ResolveUDPAddr("udp", Brutus.Credentials.Address+":9014")
    if err != nil {
        panic(err)
    }
    brutusport2, err := net.DialUDP("udp", nil, brutusaddr2)
    if err != nil {
        panic(err)
    }
    BrutusCutPort = brutusport2
}
//// Gather Disposition Intel
func GatherDispositions() {
    message := MSG {
        Army: ARMY {
                    Catapults: 0,
                    Archers: 0,
                    Cavalry: 0,
                    Spearmen: 0,
                    Infantry: 0,
                },
        To: "Pompus",
        From: "Caesar",
        Message: "Dispositions",
        State: 0,
        Messages: nil,
    }
    fmt.Printf("asking Pompus...")
    pompusmsg := SendGetResponse(PompusPort, CaesarPort, message)
    Pompus.Army = pompusmsg.Army
    fmt.Printf("asking Operachorus...")
    message.To = "Operachorus"
    operachorusmsg := SendGetResponse(OperachorusPort, CaesarPort, message)
    Operachorus.Army = operachorusmsg.Army
    fmt.Printf("asking Brutus...")
    message.To = "Brutus"
    brutusmsg := SendGetResponse(BrutusPort, CaesarPort, message)
    Brutus.Army = brutusmsg.Army
    fmt.Printf("dispositions known!\n")
    Totals = ARMY{
        Catapults: Caesar.Army.Catapults+Pompus.Army.Catapults+Operachorus.Army.Catapults+Brutus.Army.Catapults,
        Archers: Caesar.Army.Archers+Pompus.Army.Archers+Operachorus.Army.Archers+Brutus.Army.Archers,
        Cavalry: Caesar.Army.Cavalry+Pompus.Army.Cavalry+Operachorus.Army.Cavalry+Brutus.Army.Cavalry,
        Spearmen: Caesar.Army.Spearmen+Pompus.Army.Spearmen+Operachorus.Army.Spearmen+Brutus.Army.Spearmen,
        Infantry: Caesar.Army.Infantry+Pompus.Army.Infantry+Operachorus.Army.Infantry+Brutus.Army.Infantry,
    }
    Objective = ARMY{
        Catapults: Totals.Catapults / 4,
        Archers: Totals.Archers / 4,
        Cavalry: Totals.Cavalry / 4,
        Spearmen: Totals.Spearmen / 4,
        Infantry: Totals.Infantry / 4,
    }
    Remainder = ARMY{
        Catapults: Totals.Catapults % 4,
        Archers: Totals.Archers % 4,
        Cavalry: Totals.Cavalry % 4,
        Spearmen: Totals.Spearmen % 4,
        Infantry: Totals.Infantry % 4,
    }
    
    // send objective message
    msg := MSG {
        Army: Objective,
        To: "Pompus",
        From: "Caesar",
        Message: "ObjectiveTroops",
        State: 0,
        Messages: nil,
    }
    SendMessage(PompusPort, msg)
    message.To = "Operachorus"
    SendMessage(OperachorusPort, msg)
    message.To = "Brutus"
    SendMessage(BrutusPort, msg)
    // send remainder message
    message.Army = Remainder
    message.To = "Pompus"
    SendMessage(PompusPort, msg)
    message.To = "Operachorus"
    SendMessage(OperachorusPort, msg)
    message.To = "Brutus"
    SendMessage(BrutusPort, msg)
    
    PrintTotals()
}
//// ask for updated dispositions
func AskDispositions() {
    message := MSG {
        Army: ARMY {
                    Catapults: 0,
                    Archers: 0,
                    Cavalry: 0,
                    Spearmen: 0,
                    Infantry: 0,
                },
        To: "Pompus",
        From: "Caesar",
        Message: "Dispositions",
        State: CurrentState,
        Messages: nil,
    }
    fmt.Printf("\nAsking for updated dispositions...asking Pompus...")
    pompusmsg := SendGetResponse(PompusPort, CaesarPort, message)
    Pompus.Army = pompusmsg.Army
    fmt.Printf("asking Operachorus...")
    message.To = "Operachorus"
    operachorusmsg := SendGetResponse(OperachorusPort, CaesarPort, message)
    Operachorus.Army = operachorusmsg.Army
    fmt.Printf("asking Brutus...")
    message.To = "Brutus"
    brutusmsg := SendGetResponse(BrutusPort, CaesarPort, message)
    Brutus.Army = brutusmsg.Army
    fmt.Printf("updated dispositions known!\n")
    PrintOfficer(Caesar, 0)
    PrintOfficer(Brutus, 0)
    PrintOfficer(Operachorus, 0)
    PrintOfficer(Pompus, 0)
    time.Sleep(time.Millisecond * 500)
}
//// what to do when the Legions are Ready
func SayReady() {
    message := MSG {
        Army: ARMY {
                    Catapults: 0,
                    Archers: 0,
                    Cavalry: 0,
                    Spearmen: 0,
                    Infantry: 0,
                },
        To: "Pompus",
        From: "Caesar",
        Message: "Ready!",
        State: CurrentState,
        Messages: nil,
    }
    SendMessage(PompusPort, message)
    message.To = "Operachorus"
    SendMessage(OperachorusPort, message)
    message.To = "Brutus"
    SendMessage(BrutusPort, message)
}
//// the handler for cut messages
func HandleCutMessage(message MSG) {
    if message.From == "Brutus" {
        Cut.Brutus = OFFICERCUT {
            Army: message.Army,
            Messages: message.Messages,
        }
        Cut.HasBrutus = true
    } else if message.From == "Operachorus" {
        Cut.Operachorus = OFFICERCUT {
            Army: message.Army,
            Messages: message.Messages,
        }
        Cut.HasOperachorus = true
    } else if message.From == "Pompus" {
        Cut.Pompus = OFFICERCUT {
            Army: message.Army,
            Messages: message.Messages,
        }
        Cut.HasPompus = true
    }
}
//// The routine to thread for cuts
func CutRoutine() {
    fmt.Println("Taking a cut!")
    Cut.CaesarArmy = Caesar.Army
    Cut.State = CurrentState
    message := MSG {
        Army: ARMY {
                    Catapults: 0,
                    Archers: 0,
                    Cavalry: 0,
                    Spearmen: 0,
                    Infantry: 0,
                },
        To: "Brutus",
        From: "Caesar",
        Message: "TakeCut",
        State: CurrentState,
        Messages: nil,
    }
    SendMessage(BrutusCutPort, message)
    message.To = "Operachorus"
    SendMessage(OperachorusCutPort, message)
    message.To = "Pompus"
    SendMessage(PompusCutPort, message)
    
    for {
        if Cut.State != CurrentState { break }
        time.Sleep(time.Millisecond * 5)
    }
    
    message.Message = "CutDone"
    HandleCutMessage(SendGetResponse(BrutusCutPort, CutPort, message))
    HandleCutMessage(SendGetResponse(OperachorusCutPort, CutPort, message))
    HandleCutMessage(SendGetResponse(PompusCutPort, CutPort, message))
    
    fmt.Println("Creating a snapshot file...")
    // assemble cut information into a string
    var out string
    out = "CUT AT STATE " + strconv.Itoa(Cut.State) + "\n"
    out += "THEORETICAL TOTALS" + "\n"
    out += "Catapults: " + strconv.Itoa(Totals.Catapults) + "\n"
    out += "Archers: " + strconv.Itoa(Totals.Archers) + "\n"
    out += "Cavalry: " + strconv.Itoa(Totals.Catapults) + "\n"
    out += "Spearmen: " + strconv.Itoa(Totals.Spearmen) + "\n"
    out += "Infantry: " + strconv.Itoa(Totals.Infantry) + "\n\n"
    out += "--------------------------------------------\n"
    out += "--------------------------------------------\n"
    out += "CAESAR\n"
    out += "Catapults: " + strconv.Itoa(Cut.CaesarArmy.Catapults) + "\n"
    out += "Archers: " + strconv.Itoa(Cut.CaesarArmy.Archers) + "\n"
    out += "Cavalry: " + strconv.Itoa(Cut.CaesarArmy.Cavalry) + "\n"
    out += "Spearmen: " + strconv.Itoa(Cut.CaesarArmy.Spearmen) + "\n"
    out += "Infantry: " + strconv.Itoa(Cut.CaesarArmy.Infantry) + "\n\n"
    out += "Messages..."
    for i := 0;  i<len(Cut.CaesarMessages); i++ {
        out += "Message #" + strconv.Itoa(i + 1)
        out += "\nArmy { Catapults: " + strconv.Itoa(Cut.CaesarMessages[i].Army.Catapults) + ", Archers: " + strconv.Itoa(Cut.CaesarMessages[i].Army.Archers) + ", Cavalry: " + strconv.Itoa(Cut.CaesarMessages[i].Army.Cavalry) + ", Spearmen: " + strconv.Itoa(Cut.CaesarMessages[i].Army.Spearmen) + ", Infantry: " + strconv.Itoa(Cut.CaesarMessages[i].Army.Infantry) + " }\n"
        out += "To: " + Cut.CaesarMessages[i].To + "\n"
        out += "From: " + Cut.CaesarMessages[i].From + "\n"
        out += "Message: " + Cut.CaesarMessages[i].Message + "\n"
        out += "State: " + strconv.Itoa(Cut.CaesarMessages[i].State) + "\n\n"
    }
    out += "\n--------------------------------------------\n"
    out += "--------------------------------------------\n"
    out += "BRUTUS\n"
    out += "Catapults: " + strconv.Itoa(Cut.Brutus.Army.Catapults) + "\n"
    out += "Archers: " + strconv.Itoa(Cut.Brutus.Army.Archers) + "\n"
    out += "Cavalry: " + strconv.Itoa(Cut.Brutus.Army.Cavalry) + "\n"
    out += "Spearmen: " + strconv.Itoa(Cut.Brutus.Army.Spearmen) + "\n"
    out += "Infantry: " + strconv.Itoa(Cut.Brutus.Army.Infantry) + "\n\n"
    out += "Messages..."
    for i := 0;  i<len(Cut.Brutus.Messages); i++ {
        out += "Message #" + strconv.Itoa(i + 1)
        out += "\nArmy { Catapults: " + strconv.Itoa(Cut.Brutus.Messages[i].Army.Catapults) + ", Archers: " + strconv.Itoa(Cut.Brutus.Messages[i].Army.Archers) + ", Cavalry: " + strconv.Itoa(Cut.Brutus.Messages[i].Army.Cavalry) + ", Spearmen: " + strconv.Itoa(Cut.Brutus.Messages[i].Army.Spearmen) + ", Infantry: " + strconv.Itoa(Cut.Brutus.Messages[i].Army.Infantry) + " }\n"
        out += "To: " + Cut.Brutus.Messages[i].To + "\n"
        out += "From: " + Cut.Brutus.Messages[i].From + "\n"
        out += "Message: " + Cut.Brutus.Messages[i].Message + "\n"
        out += "State: " + strconv.Itoa(Cut.Brutus.Messages[i].State) + "\n\n"
    }
    out += "\n--------------------------------------------\n"
    out += "--------------------------------------------\n"
    out += "OPERACHORUS\n"
    out += "Catapults: " + strconv.Itoa(Cut.Operachorus.Army.Catapults) + "\n"
    out += "Archers: " + strconv.Itoa(Cut.Operachorus.Army.Archers) + "\n"
    out += "Cavalry: " + strconv.Itoa(Cut.Operachorus.Army.Cavalry) + "\n"
    out += "Spearmen: " + strconv.Itoa(Cut.Operachorus.Army.Spearmen) + "\n"
    out += "Infantry: " + strconv.Itoa(Cut.Operachorus.Army.Infantry) + "\n\n"
    out += "Messages..."
    for i := 0;  i<len(Cut.Operachorus.Messages); i++ {
        out += "Message #" + strconv.Itoa(i + 1)
        out += "\nArmy { Catapults: " + strconv.Itoa(Cut.Operachorus.Messages[i].Army.Catapults) + ", Archers: " + strconv.Itoa(Cut.Operachorus.Messages[i].Army.Archers) + ", Cavalry: " + strconv.Itoa(Cut.Operachorus.Messages[i].Army.Cavalry) + ", Spearmen: " + strconv.Itoa(Cut.Operachorus.Messages[i].Army.Spearmen) + ", Infantry: " + strconv.Itoa(Cut.Operachorus.Messages[i].Army.Infantry) + " }\n"
        out += "To: " + Cut.Operachorus.Messages[i].To + "\n"
        out += "From: " + Cut.Operachorus.Messages[i].From + "\n"
        out += "Message: " + Cut.Operachorus.Messages[i].Message + "\n"
        out += "State: " + strconv.Itoa(Cut.Operachorus.Messages[i].State) + "\n\n"
    }
    out += "\n--------------------------------------------\n"
    out += "--------------------------------------------\n"
    out += "POMPUS\n"
    out += "Catapults: " + strconv.Itoa(Cut.Pompus.Army.Catapults) + "\n"
    out += "Archers: " + strconv.Itoa(Cut.Pompus.Army.Archers) + "\n"
    out += "Cavalry: " + strconv.Itoa(Cut.Pompus.Army.Cavalry) + "\n"
    out += "Spearmen: " + strconv.Itoa(Cut.Pompus.Army.Spearmen) + "\n"
    out += "Infantry: " + strconv.Itoa(Cut.Pompus.Army.Infantry) + "\n\n"
    out += "Messages..."
    for i := 0;  i<len(Cut.Pompus.Messages); i++ {
        out += "Message #" + strconv.Itoa(i + 1)
        out += "\nArmy { Catapults: " + strconv.Itoa(Cut.Pompus.Messages[i].Army.Catapults) + ", Archers: " + strconv.Itoa(Cut.Pompus.Messages[i].Army.Archers) + ", Cavalry: " + strconv.Itoa(Cut.Pompus.Messages[i].Army.Cavalry) + ", Spearmen: " + strconv.Itoa(Cut.Pompus.Messages[i].Army.Spearmen) + ", Infantry: " + strconv.Itoa(Cut.Pompus.Messages[i].Army.Infantry) + " }\n"
        out += "To: " + Cut.Pompus.Messages[i].To + "\n"
        out += "From: " + Cut.Pompus.Messages[i].From + "\n"
        out += "Message: " + Cut.Pompus.Messages[i].Message + "\n"
        out += "State: " + strconv.Itoa(Cut.Pompus.Messages[i].State) + "\n\n"
    }
    
    // Print Cut information!
    fmt.Println("WRITING CUT INFORMATION TO FILE! (file = 'cutfile.txt')")
    err := ioutil.WriteFile("cutfile.txt", []byte(out), 0644)
    if err != nil {
        panic(err)
    }
    fmt.Println("CUT INFORMATION WRITTEN TO FILE! (file = 'cutfile.txt')")
    
    
    fmt.Println("CUT INFO BEGIN")
    fmt.Println("CUT INFO END")
}
