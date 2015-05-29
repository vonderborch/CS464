package main

//////////// LIBRARIES
import (
    "encoding/json"
    "fmt"
    "net"
    "io/ioutil"
    "time"
    "math/rand"
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
    Army ARMY
    SelfMessages []MSG
}

//////////// DEFINE GLOBAL VARIABLES
var (
    Ready bool                       // whether the officer is ready for combat!
    OfficerName string              // the name of this Officer
    CurrentState int                // the current event state for the Officer
    TakeCut bool                    // whether we should be taking a cut
    
    Objective ARMY                  // the objective army count
    Remainder ARMY                  // the remainder army count
    
    Cut CUT                         // the cut information
    
    Self OFFICER                    // information about this officer
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
    
    HasObjectiveArmy bool           // the objective army for this officer
    HasRemainderArmy bool           // the remainder army for this officer
    
    LastMessageStr string           // the last message message
)

//////////// MAIN FUNCTION
func main() {
    TakeCut = false
    Ready = false
    HasObjectiveArmy = false
    HasRemainderArmy = false
    //////// SEED RAND AND DETERMINE WHO THIS OFFICER IS
    rand.Seed(time.Now().UTC().UnixNano())
    if (len(os.Args) == 1) {
        panic("Usage is: go run Officer.go [OfficerName]\nValid Officers: Brutus, Operachorus, Pompus")
    } else {
        i := os.Args[1]
        if i == "Brutus" {
            OfficerName = i
        } else if i == "Operachorus" {
            OfficerName = i
        } else if i == "Pompus" {
            OfficerName = i
        } else {
            panic("Usage is: go run Officer.go [OfficerName]\nValid Officers: Brutus, Operachorus, Pompus")
        }
    }
    
    CurrentState = 0
    //////// SETUP OFFICERS
    fmt.Println(OfficerName+" is setting up his tables...")
    SetupOfficers()
    //////// SETUP NETWORKING
    fmt.Println(OfficerName+" is choosing his messenger's routes...")
    SetupNetworking()
    //////// SETUP ARMY
    Self.Army = ARMY {
        Catapults: 1+rand.Intn(100),
        Archers: 1+rand.Intn(100),
        Cavalry: 1+rand.Intn(100),
        Spearmen: 1+rand.Intn(100),
        Infantry: 1+rand.Intn(100),
    }
    PrintArmy(Self.Army, OfficerName+"'s Base Army")
    
    //////// MAIN LOOP
    fmt.Println(OfficerName+" is running his camp...")
    go CutLoop()
    for {
        if Ready == true {
            break
        }
        if OfficerName == "Brutus" {
            HandleMessage(GetResponse(BrutusPort))
        } else if OfficerName == "Operachorus" {
            HandleMessage(GetResponse(OperachorusPort))
        } else if OfficerName == "Pompus" {
            HandleMessage(GetResponse(PompusPort))
        }
    }
}

//////////// FUNCTIONS
//// Prints the number of troops a officer has at his disposal
func PrintArmy(army ARMY, title string) {
    fmt.Printf("%s\n", title);
    fmt.Printf("Catapults: %d\n", army.Catapults)
    fmt.Printf("Archers: %d\n", army.Archers)
    fmt.Printf("Cavalry: %d\n", army.Cavalry)
    fmt.Printf("Spearmen: %d\n", army.Spearmen)
    fmt.Printf("Infantry: %d\n", army.Infantry)
}
//// Handles a received message
func HandleMessage(message MSG) {
    fmt.Println("Received message!")
    msg := MSG {
        Army: message.Army,
        To: "Caesar",
        From: OfficerName,
        Message: "",
        State: CurrentState,
        Messages: nil,
    }
    LastMessageStr = ""
    if message.Message == "Dispositions" {
        Caesar.State = message.State
        msg.Army = Self.Army
        msg.Message = "CurrentDispositions"
        SendMessage(CaesarPort, msg)
        fmt.Println("Sent disposition to Caesar")
    } else if message.Message == "TakeCut" {
        TakeCut = true
        fmt.Printf("Asked to make cut!")
        Cut.Army = Self.Army
        go CutLoop()
    } else if message.Message == "CutDone" {
        TakeCut = false
        fmt.Println("Cut done!")
        msg.Army = Cut.Army
        msg.Messages = Cut.SelfMessages
        SendMessage(CutPort, msg)
    } else if message.Message == "SendTroopsToBrutus" {
        LastMessageStr = "SendTroopsToBrutus"
        msg.To = "Brutus"
        msg.Message = "ReceiveTroops"
        Self.Army = ARMY {
            Catapults: Self.Army.Catapults+message.Army.Catapults,
            Archers: Self.Army.Archers+message.Army.Archers,
            Cavalry: Self.Army.Cavalry+message.Army.Cavalry,
            Spearmen: Self.Army.Spearmen+message.Army.Spearmen,
            Infantry: Self.Army.Infantry+message.Army.Infantry,
        }
        fmt.Println("Send troops to brutus...")
        t := SendGetResponse(BrutusPort, GetOwnPort(), msg)
        if (t.To != "") { }
        message.Message = "Sent"
        SendMessage(CaesarPort, message)
        CurrentState++
    } else if message.Message == "SendTroopsToOperachorus" {
        LastMessageStr = "SendTroopsToOperachorus"
        msg.To = "Operachorus"
        msg.Message = "ReceiveTroops"
        Self.Army = ARMY {
            Catapults: Self.Army.Catapults+message.Army.Catapults,
            Archers: Self.Army.Archers+message.Army.Archers,
            Cavalry: Self.Army.Cavalry+message.Army.Cavalry,
            Spearmen: Self.Army.Spearmen+message.Army.Spearmen,
            Infantry: Self.Army.Infantry+message.Army.Infantry,
        }
        fmt.Println("Send troops to operachorus...")
        t := SendGetResponse(OperachorusPort, GetOwnPort(), msg)
        if (t.To != "") { }
        message.Message = "Sent"
        SendMessage(CaesarPort, message)
        CurrentState++
    } else if message.Message == "SendTroopsToPompus" {
        LastMessageStr = "SendTroopsToPompus"
        msg.To = "Pompus"
        msg.Message = "ReceiveTroops"
        Self.Army = ARMY {
            Catapults: Self.Army.Catapults+message.Army.Catapults,
            Archers: Self.Army.Archers+message.Army.Archers,
            Cavalry: Self.Army.Cavalry+message.Army.Cavalry,
            Spearmen: Self.Army.Spearmen+message.Army.Spearmen,
            Infantry: Self.Army.Infantry+message.Army.Infantry,
        }
        fmt.Println("Send troops to pompus...")
        t := SendGetResponse(PompusPort, GetOwnPort(), msg)
        if (t.To != "") { }
        message.Message = "Sent"
        SendMessage(CaesarPort, message)
        CurrentState++ 
    } else if message.Message == "SendTroopsToCaesar" {
        LastMessageStr = "SendTroopsToCaesar"
        msg.To = "Caesar"
        msg.Message = "ReceiveTroops"
        Self.Army = ARMY {
            Catapults: Self.Army.Catapults+message.Army.Catapults,
            Archers: Self.Army.Archers+message.Army.Archers,
            Cavalry: Self.Army.Cavalry+message.Army.Cavalry,
            Spearmen: Self.Army.Spearmen+message.Army.Spearmen,
            Infantry: Self.Army.Infantry+message.Army.Infantry,
        }
        fmt.Println("Send troops to caesar...")
        t := SendGetResponse(CaesarPort, GetOwnPort(), msg)
        if (t.To != "") { }
        message.Message = "Sent"
        SendMessage(CaesarPort, message)
        CurrentState++
    } else if message.Message == "ReceiveTroops" {
        message.To = message.From
        message.From = OfficerName
        message.Message = "Received"
        Self.Army = ARMY {
            Catapults: Self.Army.Catapults-message.Army.Catapults,
            Archers: Self.Army.Archers-message.Army.Archers,
            Cavalry: Self.Army.Cavalry-message.Army.Cavalry,
            Spearmen: Self.Army.Spearmen-message.Army.Spearmen,
            Infantry: Self.Army.Infantry-message.Army.Infantry,
        }
        fmt.Println("Received Troops")
        SendMessage(GetPortFromName(message.To), message)
        
        if (message.To == "Brutus") {
            Brutus.State = message.State
        } else if (message.To == "Operachorus") {
            Operachorus.State = message.State
        } else if (message.To == "Pompus") {
            Pompus.State = message.State
        }
        CurrentState++
    } else if message.Message == "ObjectiveTroops" {
        Objective = message.Army
        HasObjectiveArmy = true
    } else if message.Message == "RemainderTroops" {
        Remainder = message.Army
        HasRemainderArmy = true
    } else if message.Message == "Ready!" {
        Ready = true
    }
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
//// Gets this Officer's port
func GetOwnPort() *net.UDPConn {
    if OfficerName == "Brutus" {
        return BrutusPort
    } else if OfficerName == "Operachorus" {
        return OperachorusPort
    } else if OfficerName == "Pompus" {
        return PompusPort
    }
    return nil
}
//// Handle cut taking
func HandleCut() {
    for {
        if TakeCut == false {
            if (IsEmpty(Cut.Army) == true) {
                Cut.Army = Self.Army
            }
        
            msg := MSG {
                Army: Cut.Army,
                To: "Caesar",
                From: OfficerName,
                Message: "Cut",
                State: CurrentState,
                Messages: Cut.SelfMessages,
            }
            SendMessage(CutPort, msg)
            fmt.Println("Sent a snapshot!")
            break
        }
        time.Sleep(time.Millisecond * 15)
    }
}
//// The loop to run to check if we need to ever take a cut
func CutLoop() {
    for {
        if Ready == true {
            break
        }
        if OfficerName == "Brutus" {
            HandleMessage(GetResponse(BrutusCutPort))
        } else if OfficerName == "Operachorus" {
            HandleMessage(GetResponse(OperachorusCutPort))
        } else if OfficerName == "Pompus" {
            HandleMessage(GetResponse(PompusCutPort))
        }
    }
}
//// Checks if the army is empty
func IsEmpty(army ARMY) bool {
    if army.Catapults == 0 && army.Archers == 0 && army.Cavalry == 0 && army.Spearmen == 0 && army.Infantry == 0 {
        return true
    }
    return false
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
    
    // check if a cut!
    if TakeCut == true {
        Cut.SelfMessages = append(Cut.SelfMessages, res)
    }
    
    return res
}
//// Sends the message and returns the response
func SendGetResponse(conn *net.UDPConn, list *net.UDPConn, message MSG) MSG {
    SendMessage(conn, message)
    return GetResponse(list)
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
    Caesar = OFFICER{
        Credentials: caesarBase,
        Army: army,
        State: 0,
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
}
//// Sets up the ports for the officers
func SetupNetworking() {
    fmt.Println("Caesar Address: " + Caesar.Credentials.Address)
    fmt.Println("Pompus Address: " + Pompus.Credentials.Address)
    fmt.Println("Brutus Address: " + Brutus.Credentials.Address)
    fmt.Println("Operachorus Address: " + Operachorus.Credentials.Address)
    // setup cut port
    cutaddr, err := net.ResolveUDPAddr("udp", Caesar.Credentials.Address+":9000")
    if err != nil {
        panic(err)
    }
    cutport, err := net.DialUDP("udp", nil, cutaddr)
    if err != nil {
        panic(err)
    }
    CutPort = cutport
    // setup caesar port
    caesaraddr, err := net.ResolveUDPAddr("udp", Caesar.Credentials.Address+":9001")
    if err != nil {
        panic(err)
    }
    caesarport, err := net.DialUDP("udp", nil, caesaraddr)
    if err != nil {
        panic(err)
    }
    CaesarPort = caesarport
    // setup pompus port
    if OfficerName == "Pompus" {
        pompusaddr, err := net.ResolveUDPAddr("udp", ":9002")
        if err != nil {
            panic(err)
        }
        pompusport, err := net.ListenUDP("udp", pompusaddr)
        if err != nil {
            panic(err)
        }
        PompusPort = pompusport
        pompusaddr2, err := net.ResolveUDPAddr("udp", ":9012")
        if err != nil {
            panic(err)
        }
        pompusport2, err := net.ListenUDP("udp", pompusaddr2)
        if err != nil {
            panic(err)
        }
        PompusCutPort = pompusport2
    } else {
        pompusaddr, err := net.ResolveUDPAddr("udp", Pompus.Credentials.Address+":9002")
        if err != nil {
            panic(err)
        }
        pompusport, err := net.DialUDP("udp", nil, pompusaddr)
        if err != nil {
            panic(err)
        }
        PompusPort = pompusport
        pompusaddr2, err := net.ResolveUDPAddr("udp", Pompus.Credentials.Address+":9012")
        if err != nil {
            panic(err)
        }
        pompusport2, err := net.DialUDP("udp", nil, pompusaddr2)
        if err != nil {
            panic(err)
        }
        PompusCutPort = pompusport2
    }
    // setup operachorus port
    if OfficerName == "Operachorus" {
        operachorusaddr, err := net.ResolveUDPAddr("udp", ":9003")
        if err != nil {
            panic(err)
        }
        operachorusport, err := net.ListenUDP("udp", operachorusaddr)
        if err != nil {
            panic(err)
        }
        OperachorusPort = operachorusport
        operachorusaddr2, err := net.ResolveUDPAddr("udp", ":9013")
        if err != nil {
            panic(err)
        }
        operachorusport2, err := net.ListenUDP("udp", operachorusaddr2)
        if err != nil {
            panic(err)
        }
        OperachorusCutPort = operachorusport2
    } else {
        operachorusaddr, err := net.ResolveUDPAddr("udp", Operachorus.Credentials.Address+":9003")
        if err != nil {
            panic(err)
        }
        operachorusport, err := net.DialUDP("udp", nil, operachorusaddr)
        if err != nil {
            panic(err)
        }
        OperachorusPort = operachorusport
        operachorusaddr2, err := net.ResolveUDPAddr("udp", Operachorus.Credentials.Address+":9013")
        if err != nil {
            panic(err)
        }
        operachorusport2, err := net.DialUDP("udp", nil, operachorusaddr2)
        if err != nil {
            panic(err)
        }
        OperachorusCutPort = operachorusport2
    }
    // setup brutus port
    if OfficerName == "Brutus" {
        brutusaddr, err := net.ResolveUDPAddr("udp", ":9004")
        if err != nil {
            panic(err)
        }
        brutusport, err := net.ListenUDP("udp", brutusaddr)
        if err != nil {
            panic(err)
        }
        BrutusPort = brutusport
        brutusaddr2, err := net.ResolveUDPAddr("udp", ":9014")
        if err != nil {
            panic(err)
        }
        brutusport2, err := net.ListenUDP("udp", brutusaddr2)
        if err != nil {
            panic(err)
        }
        BrutusCutPort = brutusport2
    } else {
        brutusaddr, err := net.ResolveUDPAddr("udp", Brutus.Credentials.Address+":9004")
        if err != nil {
            panic(err)
        }
        brutusport, err := net.DialUDP("udp", nil, brutusaddr)
        if err != nil {
            panic(err)
        }
        BrutusPort = brutusport
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
}
