package selector

import (
	//"errors"
	"fmt"
	//"os"
	//"os/exec"
	"math"
	"math/rand"
	"sync"
	"strconv"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/proto"

	"github.com/dchest/siphash"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/local/connection"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/registry"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/crossconnect"
	//v1alpha1 "github.com/networkservicemesh/networkservicemesh/k8s/pkg/apis/networkservice/v1alpha1"
	//"github.com/networkservicemesh/networkservicemesh/k8s/pkg/networkservice/clientset/versioned"
	//"github.com/networkservicemesh/networkservicemesh/k8s/pkg/networkservice/informers/externalversions"
	"k8s.io/client-go/rest"
	//"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"github.com/skydive-project/skydive/logging"

	//"github.com/networkservicemesh/networkservicemesh/k8s/cmd/crossconnect-monitor"
)

//import "math/rand"

type maglevSelector struct {
	sync.Mutex
	nseNbr uint64                             // Number of NSE pods (i.e., backend in the pool)
	Lk_Size uint64                             // size of the lookup table (number of requests)
	permutation [][]uint64                         // used to compute hash value per pod
	
	LookupTable map[uint64]int                  // Lookup table considering ReqId instead of index
	nseList []*registry.NetworkServiceEndpoint // list of nse candidates

	maglev map[string]uint64
}

// Batch variable, to check whether the request is the first on the batch or not
var FirstRequestInBatch bool = true

// Number of requests in the batch, corresponding to the size of Lookup table
var BatchSize uint64 = 2 // to be varied later

// Loockup traking table for Maglev algorithm
//var LoockupTable map[int]int // ReqId <-> NSE pod id

// Should be a global variable in Model
var RequestIdPerNetworkService map[string]map[string]uint64

// Lookup table for each NS to save decision
var LookupTablePerNetworkService map[string]map[uint64]int

//func NewmaglevSelector() Selector {
//	return &maglevSelector{
//		nseNbr: uint64,
//		Lk_Size: uint64,
//		permutation: make([][]uint64),
//		LookupTable: make(map[uint64]uint64),
//		nseList: make([]*registry.NetworkServiceEndpoint),
//		maglev: make(map[string]uint64),
//	}
//}
var MetricsPerEndpoint = map[string]map[string]*crossconnect.Metrics{}
//var CurrentServerMetrics = crossconnect.Metrics{}


const nsmResource = "networkservicemanagers"
var managers = map[string]string{}

var informerStopper chan struct{}





// Get Number of requests in the same batch requesting the same network service secure-intranet-connectivty for the current deployment
/*func SetNumberofRequests_InTheBatch() {

	//cmd := exec.Command("bash", "/go/src/github.com/networkservicemesh/networkservicemesh/scripts/nsc_pods_number.sh")
	//logrus.Infof("echo cmd %v ",cmd)

	var BatchSize_NS string
    BatchSize_NS = os.Getenv("NSC_PER_NS")
	logrus.Infof(" get BatchSize_NS %s ",BatchSize_NS)

    // convert the batch size to integer
    ConvertBatchSize, err := strconv.Atoi(BatchSize_NS)
	if err != nil {
			// handle error
	}
	logrus.Infof(" ConvertBatchSize %d",ConvertBatchSize)
	BatchSize =	uint64(ConvertBatchSize)

	var CLUSTER_PREFIX string
	CLUSTER_PREFIX = os.Getenv("CLUSTER_RULES_PREFIX")
	logrus.Infof(" obtained BatchSize %d CLUSTER_RULES_PREFIX %s ",BatchSize, CLUSTER_PREFIX)

	

}*/

func NewmaglevSelector() Selector {
	CreateMaglevBatchMaps()
	return &maglevSelector{
		LookupTable: make(map[uint64]int),
		maglev: make(map[string]uint64),
	}
}

// to be set at the begining, or for each incoming batch !!! (in the model for example, but only once)
func CreateMaglevBatchMaps() {
	RequestIdPerNetworkService = make (map[string]map[string]uint64)
	LookupTablePerNetworkService = make (map[string]map[uint64]int)
}
// Hashfunction 1 TBD
func HashFunction1(pod int) int {

	// var podName string = Name[pod]
	// TBD
	var hash1 int = rand.Intn(10) // find a function of the name!?, now I put it random
	return hash1
}

/*func SaveMetricsPerIntefaceState(metrics map[string]string, Interface string) {
	MetricsPerEndpoint[Interface] = metrics
}*/
// Hashfunction 2 TBD
func HashFunction2(pod int) int {

	// var podName string = Name[pod]
	// TBD
	var hash2 int = rand.Intn(10) // find a function of the name!?, now I put it random
	return hash2
}

// Compute Has value of pod
func ComputeHashValue(pod int, LookupSize int, nextIndex int) int {
	var hashValue int = 0

	var offset int = HashFunction1(pod) % LookupSize
	var skip int = HashFunction1(pod)%(LookupSize-1) + 1

	hashValue = (offset + nextIndex*skip) % LookupSize

	// Add also network performance metrics??

	return hashValue
}

// Method2: Use existing hash function
func (mg *maglevSelector) hashKey(obj string) uint64 {
	return siphash.Hash(0xdeadbabe, 0, []byte(obj))
}



// networkServiceEndpoints []*registry.NetworkServiceEndpoint: To be latter filled in the SelectEndpoint function: nseCandidates
func (mg *maglevSelector) CreateMaglev(nseCandidates []*registry.NetworkServiceEndpoint) error {

	//logrus.Infof("create Magglev for len(nse) %d", len(nseCandidates))
	//mg.Lock()
	//logrus.Infof("mg lock ")
	//defer mg.Unlock()
	//logrus.Infof("mg unlock ")

	//logrus.Infof("compute lookup table size ")
	mg.Lk_Size =  BatchSize //uint64(len(mg.maglev)) //BatchSize  // just to test   // set the lookup table size to requests nbr
	n := uint64(len(nseCandidates))
	//logrus.Infof("computed nse number  = %d ", n)
	// Useless, because Maglev should share the nbr of nse between all flows!!

	/*if mg.Lk_Size < n {
		logrus.Infof("error Number of nseCandidates is greater than lookup table %d !!", n)
		return errors.New("Number of nseCandidates is greater than lookup table")
	}*/

	//logrus.Infof("set nseList to nseCandidates ")
	mg.nseList = nseCandidates
	//m.nodeList = make([]*registry.NetworkServiceEndpoint, n)
	//logrus.Infof("copy nseCandidates in nseList")
	copy(mg.nseList, nseCandidates) // Copy to avoid modifying orinal input afterwards
	mg.nseNbr = n
	//logrus.Infof("mg.nseNbr  = %d ", mg.nseNbr)
	//logrus.Infof("end creation!")
	return nil
}

func (mg *maglevSelector) ComputeLoadPerpods(nse *registry.NetworkServiceEndpoint) float64 {

	var load float64

	load = rand.Float64() /// to be set according to performance metrics on NSM
	return load
}

// Compute hashValue, considering load balancing, Ref ConsistentHashing paper 2019

func (mg *maglevSelector) ConsistentHashing() {

	mg.permutation =  make([][]uint64, mg.Lk_Size)
	if len(mg.nseList) == 0 {
		return
	}
	// compute the total load of pods: mu_p

	var mu_p float64 = 0

	for i := 0; i < len(mg.nseList); i++ {
		// compute
		mu_p += mg.ComputeLoadPerpods(mg.nseList[i])

	}
	mu_p = (mu_p / float64(len(mg.nseList)))

	// compute fairness function
	var fairness = float64(0)
	for i := 0; i < len(mg.nseList); i++ {

		var load = mg.ComputeLoadPerpods(mg.nseList[i])

		fairness += math.Pow((load - fairness), 2)
	}
	fairness = math.Sqrt((fairness / (float64(len(mg.nseList)))))

	// How to use fairnes to compute permutation[][]??!!

}

// Compute Has value of pod
func (mg *maglevSelector) ComputeHashValues() {

	//mg.permutation = nil
	mg.permutation =  make([][]uint64, mg.nseNbr)
	logrus.Infof("len(mg.nseList) %d vs mg.nseNbr %d vs mg.Lk_Size %d ", len(mg.nseList), mg.nseNbr, mg.Lk_Size)
	if len(mg.nseList) == 0 {
		return
	}
	
	
	//sort.Strings(mg.nseList)
	for i := 0; i < len(mg.nseList); i++ {

		bNseName := []byte(mg.nseList[i].GetName())
		logrus.Infof("bNseName %s for index %d vs Lk_size %d ", mg.nseList[i].GetName(), i, mg.Lk_Size)
	
		offset := siphash.Hash(0xdeadbabe, 0, bNseName) % mg.Lk_Size
		div := (mg.Lk_Size - 1)
		if mg.Lk_Size == 1 {
			div = 1
		}
		skip := (siphash.Hash(0xdeadbeef, 0, bNseName) % (div)) + 1
		logrus.Infof("offset %d skip %d mg.Lk_Size %d ", offset, skip, mg.Lk_Size)

		mg.permutation[i] =  make([]uint64, mg.Lk_Size)
		var j int 
		//uint64
		for j = 0; j < int(mg.Lk_Size); j++ {
			//iRow[j] = (offset + uint64(j)*skip) % mg.Lk_Size
			//logrus.Infof("len (permutation) = %d ", len(mg.permutation))
			logrus.Infof("len (permutation[i]) = %d ", len(mg.permutation[i]))
			mg.permutation[i][j] = (offset + uint64(j)*skip) % mg.Lk_Size
			logrus.Infof("permutation for i %d j %d = %d ", i, j, mg.permutation[i][j])
		}

	}

	// Add also network performance metrics??

}

func ComputehashValuePerRequest(ReqId uint64) uint64 {

	// b := []byte(strconv.FormatInt(num, 10))
	bReqName := []byte(strconv.FormatInt(int64(ReqId),10))

	//bReqName := []byte(ReqId)
	
	// we consider pnly hash function without mod M
	ReqHashValue := siphash.Hash(0xdeadbabe, 0, bReqName) 
	//% mg.Lk_Size
	//skip := (siphash.Hash(0xdeadbeef, 0, bReqName) % (mg.Lk_Size - 1)) + 1
	//logrus.Infof("offset %d skip %d ", offset, skip)
	return ReqHashValue
}

func getK8SConfig() (*rest.Config, error) {
	// check if CRD is installed
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}

	logging.GetLogger().Debugf("Unable to get in K8S cluster config, trying with KUBECONFIG env")

	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
	return kubeconfig.ClientConfig()

}

func monitorCrossConnects(address string, continuousMonitor bool) {
	var err error
	logrus.Infof("Starting CrossConnections Monitor on %s", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		logrus.Errorf("failure to communicate with the socket %s with error: %+v", address, err)
		return
	}
	logrus.Infof("conn.Close() on %s conn %v ", address, conn)
	defer conn.Close()
	logrus.Infof("get dataplaneClient on %s", address)
	dataplaneClient := crossconnect.NewMonitorCrossConnectClient(conn)
	logrus.Infof(" dataplaneClient obtained! get stream")
	// Looping indefinitely or until grpc returns an error indicating the other end closed connection.
	stream, err := dataplaneClient.MonitorCrossConnects(context.Background(), &empty.Empty{})
	logrus.Infof("get stream %v on address %s  ", stream, address)
	if err != nil {
		logrus.Warningf("Error: %+v.", err)
		logrus.Infof("Error getting stream: %+v.", err)
		return
	}
	logrus.Infof("stream obtained! get TextMarshaler")
	t := proto.TextMarshaler{}
	logrus.Infof("proto.TextMarshaler on %s", address)
	for {
		event, err := stream.Recv()
		if err != nil {
			logrus.Errorf("Error: %+v.", err)
			return
		}
		data := fmt.Sprintf("event type \u001b[31m*** %s\n\u001b[0m", event.Type)
		logrus.Infof("event type in Maglev \u001b[31m*** %s\n\u001b[0m", event.Type)
		logrus.Infof("event metrics in Maglev %v ", event.Metrics)
		data += fmt.Sprintf("address in Maglev \u001b[31m*** %s\n\u001b[0m", address)
		logrus.Infof("address in Maglev \u001b[31m*** %s\n\u001b[0m", address)
		logrus.Infof(" event CrossConnects  in Maglev %v ",  event.CrossConnects)
		for _, cc := range event.CrossConnects {
			if cc != nil {
				data += fmt.Sprintf("TextMarshaler Maglev of event crossConnects \u001b[32m%s\n\u001b[0m", t.Text(cc))
				// replace address by the local endpoint name to get from crossConnect struct, and log file
				MetricsPerEndpoint[address] = event.Metrics
				
			}
		}
		println(data)
		if !continuousMonitor {
			logrus.Infof("Monitoring of server: %s. is complete...", address)
			delete(managers, address)
			return
		}
		logrus.Infof("continuousMonitor in Maglev ")
	}
}
// Compute the weight of a pod based on the NSM performance metrics
func SaveCrossConnectMetrics() {

	logrus.Infof("start SaveCrossConnectMetrics ")

	// test with random now, but to be updates by metrics later
	//weight := 1 + rand.Float64() *(99)
	//weight := float64(-1.0)

	//min + rand.Float64() * (max - min)

	config, err := getK8SConfig()
	if err != nil {
		return 
	}

	logging.GetLogger().Debugf("NSM: getting NSM client")
	logging.GetLogger().Debugf("NSM: getting config %v ", config)
	// Initialize clientset

	
	/*logging.GetLogger().Debugf("NSM: getting v1alpha1 %v ", v1alpha1.SchemeGroupVersion.WithResource("networkservicemanagers"))
	nsmClientSet, err := versioned.NewForConfig(config)

	if err != nil {
		logging.GetLogger().Errorf("Unable to initialize the NSM probe: %v", err)
		return 
	}
	logging.GetLogger().Debugf("NSM: getting nsmClientSet %v ", nsmClientSet)
	factory := externalversions.NewSharedInformerFactory(nsmClientSet, 0)
	logging.GetLogger().Debugf("NSM: getting factory %v ", factory)
	genericInformer, err := factory.ForResource(v1alpha1.SchemeGroupVersion.WithResource(nsmResource))
	if err != nil {
		logging.GetLogger().Errorf("Unable to create the K8S cache factory: %v", err)
		return 
	}
	logging.GetLogger().Debugf("NSM: getting genericInformer %v ", genericInformer)
	//monitorClient := crossconnect.NewMonitorCrossConnectClient(conn)
	//stream, err := monitorClient.monitorCrossConnects(nsm.Status.URL)

	informer := genericInformer.Informer()
	informerStopper := make(chan struct{})
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			nsm := obj.(*v1alpha1.NetworkServiceManager)
			logging.GetLogger().Infof("New NSMgr Added: %v", nsm)
			//go p.monitorCrossConnects(nsm.Status.URL)
			//go monitorCrossConnects(nsm.Status.URL, true)
		},
	})
	go informer.Run(informerStopper)*/

	// get the weight from the crossConnectMetrics

	// save the it on pod weights

	// return 
  
}
// Compute the weight of a pod based on the NSM performance metrics
func GetPodWeight(endpoint string) float64{

	logrus.Infof("GetPodWeight for endpoint %s ", endpoint)

	// test with random now, but to be updates by metrics later
	weight := 1 + rand.Float64() *(99)
	// weight := float64(-1.0)


	cumulMetrics := float64(0.0000)
	// Should Get only MetricsPerEndpoint[endpoint]??
	for k, v := range MetricsPerEndpoint { 
		logrus.Infof("nsm address[%s] value[%s]\n", k, v)
		for metric, Value := range v {
			logrus.Infof("metric[%s] Value[%v]\n", metric, Value)
			for m, val_metric := range Value.Metrics {
				
				logrus.Infof("metric[%s] val_metric[%v]\n", m, val_metric)
				rx_bytes, err := strconv.ParseFloat(val_metric,4) //strconv.Atoi(val_metric)
				if err == nil {
					fmt.Println(rx_bytes)
				}
				cumulMetrics = float64(cumulMetrics + rx_bytes)
				
			}
		}
	}

	// Weight depending on metrics ??!
	weight = float64(cumulMetrics)

	// get the weight from the cumulMetrics in crossConnectMetrics


	return weight

}

// function Popoluate maglev hashing lookup table, considering weights
func (mg *maglevSelector) PopulateMaglevHashing(requestConnection *connection.Connection, nseNbr int, ns *registry.NetworkService) {

	// ##### Variable to save global pod selection per batch
	var LookupIdxPerRequestId = map[string]uint64{}
	//var ok bool = true

	// dynamic table Next

	var podID int
	var next int
	var entry uint64
	Next := make([]uint64, mg.nseNbr)
	//lookupTable := make([]uint64, mg.Lk_Size)
	lookupTable := make(map[uint64]int, mg.Lk_Size)

	//Intialization of tables
	//logrus.Infof("Intialization of tables")
	for podID = 0; podID < int(mg.nseNbr); podID++ {
		Next[podID] = 0
	}
	// logrus.Infof("Next table %v ",Next)
	for entry = 0; entry < mg.Lk_Size; entry++ {
		// initialVal := -1
		lookupTable[entry] = int(-1)
		// logrus.Infof("initial entry %d lookupTable[entry] %d ", entry, lookupTable[entry])
	}

	logrus.Infof("compute permutations for mg.nseNbr %d ", mg.nseNbr)
	var n uint64
	n = 0
	MinDisance := make([]float64, mg.nseNbr)
	for podID = 0; podID < int(mg.nseNbr); podID++ {
		MinDisance[podID] = 90000000000000000.0  // choose a max value
	}
	for { // true
		for next = 0; next < int(mg.Lk_Size); next++ {
			n = 0
			for podID = 0; podID < int(mg.nseNbr); podID++ {
				//logrus.Infof("compute permutation for entry %d  and PODID %d Next[podID]  %d ", entry, podID, Next[podID])
				//logrus.Infof("Next table before %v ",Next)
				
				entry = mg.permutation[podID][next]
				weight := GetPodWeight(mg.nseList[podID].GetName())

				//logrus.Infof("Next table after weight %v ",Next)
				// compute hash value of request n
				ReqHashValue := ComputehashValuePerRequest(n)

				// Compute the distance
				// Get la partie entiere of the difference
				IntDifference := uint64((entry - ReqHashValue))
				distance := (float64((entry - ReqHashValue)- IntDifference)/weight)

				logrus.Infof("compute distance %f weight %f  IntDifference %d ", distance, weight, IntDifference)
				logrus.Infof("lookupTable[entry]  %d entry %d podID %d ", lookupTable[entry] , entry, podID)
				if lookupTable[entry] >= 0 {
					// keep the entry with the min distance
					if MinDisance[podID] > distance {
						
						MinDisance[podID] = distance
						entry = mg.permutation[podID][next]
						logrus.Infof("set lookupTable becuas distance for entry %d to podId %d ", entry, podID)
						//Next[podID] = (Next[podID]+1)
						//lookupTable[entry] = podID
					} 
				} 
				logrus.Infof("set lookupTable for entry %d to podId %d ", entry, podID)
				lookupTable[entry] = podID
				

				if requestConnection.GetId() == "" {
					logrus.Infof("requestConnection.Id cannot be empty: %v", requestConnection)
				}
				// else {
				requestName := strconv.FormatInt(int64(entry+1), 10)
				LookupIdxPerRequestId[requestName] = entry
				//logrus.Infof("save requestName %s for entry %d Next[podID]  %d ", requestName, entry, Next[podID] )
				//}
				
				//LookupIdxPerRequestId[requestConnection.GetId()] = entry
				
			}
			n++
		
			if (n == mg.Lk_Size) {
				break
			}
				
		}
		if (n == mg.Lk_Size) || (podID == int(mg.nseNbr)) {
			//logrus.Infof("save RequestIdPerNetworkService for ns %s ",ns.GetName())
			mg.LookupTable = lookupTable
			logrus.Infof("lookupTable %v ",lookupTable)
			RequestIdPerNetworkService[ns.GetName()] = LookupIdxPerRequestId
			//logrus.Infof("set LookupTablePerNetworkService  %v", LookupTablePerNetworkService)
			LookupTablePerNetworkService[ns.GetName()] = mg.LookupTable
			//logrus.Infof("table fullfilled, return ")
			return
		}
	}

}

// function Popoluate maglev hashing lookup table
func (mg *maglevSelector) PopulateMaglevHashing_v1(requestConnection *connection.Connection, nseNbr int, ns *registry.NetworkService) {

	// ##### Variable to save global pod selection per batch
	var LookupIdxPerRequestId = map[string]uint64{}
	//var ok bool = true

	// dynamic table Next

	var podID int
	var entry uint64
	Next := make([]uint64, mg.nseNbr)
	//lookupTable := make([]uint64, mg.Lk_Size)
	lookupTable := make(map[uint64]int, mg.Lk_Size)

	//Intialization of tables
	//logrus.Infof("Intialization of tables")
	for podID = 0; podID < int(mg.nseNbr); podID++ {
		Next[podID] = 0
	}
	for entry = 0; entry < mg.Lk_Size; entry++ {
		//initialVal := -1
		lookupTable[entry] = int(-1)
		//logrus.Infof("initial entry %d lookupTable[entry] %d ", entry, lookupTable[entry])
	}

	//logrus.Infof("compute permutations for mg.nseNbr %d ", mg.nseNbr)
	var n uint64
	for { // true
		for podID = 0; podID < int(mg.nseNbr); podID++ {
			//logrus.Infof("compute permutation for entry %d  and PODID %d Next[podID]  %d ", entry, podID, Next[podID])
			entry = mg.permutation[podID][Next[podID]]
			//logrus.Infof("entry %d lookupTable[entry] %d ", entry, lookupTable[entry])
			for lookupTable[entry] >= 0 {
				Next[podID] = Next[podID] + 1
				if Next[podID]>= mg.Lk_Size {
					break
				}
				logrus.Infof("Next[podID] %d entry %d lookupTable[entry] %d ", Next[podID], entry, lookupTable[entry])
				entry = mg.permutation[podID][Next[podID]]
			}
			// logrus.Infof("set entry %d in lookuptable to podID %d ", entry, podID)
			lookupTable[entry] = podID
			// uint64(podID)
			

			// logrus.Infof("set LookupIdxPerRequestId  for reqid %s len(LookupIdxPerRequestId) %d ", requestConnection.GetId(), len(LookupIdxPerRequestId))
			
			if requestConnection.GetId() == "" {
				logrus.Infof("requestConnection.Id cannot be empty: %v", requestConnection)
			}
			// else {
			requestName := strconv.FormatInt(int64(entry+1), 10)
			LookupIdxPerRequestId[requestName] = entry
			logrus.Infof("save requestName %s for entry %d ", requestName, entry)
			//}
			
			//LookupIdxPerRequestId[requestConnection.GetId()] = entry
			

			Next[podID] = Next[podID] + 1
			n++
			if n == mg.Lk_Size {
				mg.LookupTable = lookupTable
				RequestIdPerNetworkService[ns.GetName()] = LookupIdxPerRequestId
				//logrus.Infof("set LookupTablePerNetworkService  %v", LookupTablePerNetworkService)
				LookupTablePerNetworkService[ns.GetName()] = mg.LookupTable
				//logrus.Infof("table fullfilled, return ")
				return
			}
		}
	}

}


// Useless, maglev already created, and lookup table is set in main algorithm
// include saved pod selection to the new maglev struct
func (mg *maglevSelector) SetMaglevSelection(requestConnection *connection.Connection, nseNbr int, LookupSize int, ns *registry.NetworkService) {

	mg.LookupTable = LookupTablePerNetworkService[ns.GetName()]
}

// Use Loockup table for pod selection in maglev[] table in the struct

func (mg *maglevSelector) SelectEndpoint(requestConnection *connection.Connection, ns *registry.NetworkService, networkServiceEndpoints []*registry.NetworkServiceEndpoint) *registry.NetworkServiceEndpoint {
	logrus.Infof("start SelectEndpoint for ns.name %s and requestConnection.getid %s ", ns.GetName(), requestConnection.GetId())

	logrus.Infof("networkServiceEndpoints list %v ", networkServiceEndpoints)

	// set the batchSize: Idea: Get that from the yaml file of vpn-gatewaynsc: The number of requests per batch=Number of replicas, for each nsc deployment
	
	// Set the BatchSize
	//SetNumberofRequests_InTheBatch()
	
	
	if mg == nil {
		logrus.Infof("return mg is nil ")
		return nil
	}
	if len(networkServiceEndpoints) == 0 {
		logrus.Infof("return networkServiceEndpoints is empty ")
		return nil
	}
	mg.Lock()
	defer mg.Unlock()

	// ************ crossConnectMetrics use ************* //
	
	
	// ************ end of crossConnectMetrics use ************* //

	//for endpointId = 0; endpointId < int(mg.nseNbr); endpointId++ {
		
	//SaveCrossConnectMetrics();

	//}

	// increment the number of requests per NS, which will be used as size of lookup table in maglev
	mg.maglev[ns.GetName()] =  mg.maglev[ns.GetName()] + 1

	logrus.Infof("number of requests in Lk of Maglev %d ", len(mg.maglev))

	// ####### create maglev slector #######
	mg.CreateMaglev(networkServiceEndpoints)

	var requestIdx uint64 = 0
	//var LookupIdxPerRequestId = map[string]int{}
	//var ok bool = true
	logrus.Infof("test if it is ok ")
	LookupIdxPerRequestId, ok := RequestIdPerNetworkService[ns.GetName()]

	if (ok) {

		logrus.Infof("Network service already existing ")
		fmt.Println("Network service already existing")
		// var requestName string = requestConnection.GetId()

		// requestIdx, ok := RequestIdPerNetworkService[ns.GetName()][requestConnection.GetId()]

		IdxPerRequestId, ok2 := LookupIdxPerRequestId[requestConnection.GetId()]
		ConvertReqId, err := strconv.Atoi(requestConnection.GetId())
		if err != nil {
			// handle error
		}
		logrus.Infof(" ConvertReqId %d",ConvertReqId)
		if (ok2) && RequestIdPerNetworkService[ns.GetName()] != nil && (ConvertReqId == int(IdxPerRequestId)) {
			// get request Idx
			requestIdx = IdxPerRequestId
			//LookupIdxPerRequestId[requestConnection.GetId()]
			logrus.Infof("requestIdx %d requestId %s ", requestIdx, requestConnection.GetId())
			mg.LookupTable = LookupTablePerNetworkService[ns.GetName()]

		} else {
			// Compute hash values in permutation[][]
			logrus.Infof("Compute for request consistent hashing with Maglev ")
			mg.ComputeHashValues()
			// Compute Maglev consistent hashing decision
			mg.PopulateMaglevHashing(requestConnection, len(networkServiceEndpoints), ns)
			requestIdx = LookupIdxPerRequestId[requestConnection.GetId()]
			logrus.Infof("requestIdx %d requestId %s ", requestIdx, requestConnection.GetId())
		}
	} else { // add the network service to map

		// Compute hash values in permutation[][]
		logrus.Infof("Compute consistent hashing with Maglev ")

		mg.ComputeHashValues()
		// Compute Maglev consistent hashing decision
		//logrus.Infof("PopulateMaglevHashing ")
		mg.PopulateMaglevHashing(requestConnection, len(networkServiceEndpoints), ns)

		//logrus.Infof("get requestIdx ")
		requestIdx = RequestIdPerNetworkService[ns.GetName()][requestConnection.GetId()]
		logrus.Infof("get requestIdx %d requestId %s RequestIdPerNetworkService %v ", requestIdx, requestConnection.GetId(), RequestIdPerNetworkService)

	}

	// Get selected NSE pod index
	logrus.Infof("get index for requestIdx %d requestName %s ", requestIdx, requestConnection.GetId())
	idx := LookupTablePerNetworkService[ns.GetName()][requestIdx]
	// mg.LookupTable[requestIdx] // RequestIde is string, we need index!!
	logrus.Infof("get endpoint for idx %d vs len(nsEndpoints) %d ", idx, len(networkServiceEndpoints))

	endpoint := networkServiceEndpoints[idx]

	if endpoint == nil {
		logrus.Infof("selected endpoint nil %v for idx %d ", endpoint, idx)
		return nil
	}
	//mg.maglev[ns.GetName()] = mg.maglev[ns.GetName()] + 1
	logrus.Infof("maglev selected %v with index %d ", endpoint, idx)
	
	logrus.Infof("maglev selected endpoint name %s", endpoint.GetName())
	return endpoint
}
